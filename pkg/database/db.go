package database

import (
	"context"
	"database/sql"
	"embed"
	"strings"

	_ "github.com/glebarez/go-sqlite"
	"github.com/google/uuid"
	"github.com/pressly/goose/v3"
	"github.com/zekrotja/hermans/pkg/model"
)

//go:embed migrations/*.sql
var migrationsFs embed.FS

type Database struct {
	conn *sql.DB
}

func New(file string) (*Database, error) {
	conn, err := sql.Open("sqlite", file)
	if err != nil {
		return nil, err
	}

	goose.SetBaseFS(migrationsFs)
	if err = goose.SetDialect("sqlite"); err != nil {
		return nil, err
	}
	if err = goose.Up(conn, "migrations"); err != nil {
		return nil, err
	}

	return &Database{conn: conn}, nil
}

func (t *Database) CreateOrderList(list *model.OrderList) error {
	_, err := t.conn.Exec(
		`INSERT INTO "OrderList" ("Id", "Created", "Deadline") VALUES (?, ?, ?);`,
		list.Id, list.Created, list.Deadline)
	return wrapErr(err)
}

func (t *Database) CreateOrder(orderListId string, order *model.Order) error {
	tx, err := t.conn.BeginTx(context.TODO(), nil)
	if err != nil {
		return wrapErr(err)
	}
	defer tx.Rollback()

	var drinkId sql.NullString
	if order.Drink != nil {
		drinkId.String = uuid.New().String()
		drinkId.Valid = true
		_, err = tx.Exec(
			`INSERT INTO "Drink" ("Id", "Name", "Size") VALUES (?, ?, ?);`,
			drinkId.String, order.Drink.Name, order.Drink.Size)
		if err != nil {
			return wrapErr(err)
		}
	}

	_, err = tx.Exec(
		`INSERT INTO "Order" ("Id", "Created", "Creator", "OrderListId", "DrinkId", "EditKey")
		 VALUES (?, ?, ?, ?, ?, ?);`,
		order.Id, order.Created, order.Creator, orderListId, drinkId, order.EditKey)
	if err != nil {
		return wrapErr(err)
	}

	for _, item := range order.StoreItems {
		_, err = tx.Exec(
			`INSERT INTO "OrderItems" ("OrderId", "StoreItemId") VALUES (?, ?);`,
			order.Id, item.Id)
		if err != nil {
			return wrapErr(err)
		}
		for _, variant := range item.Variants {
			_, err = tx.Exec(
				`INSERT INTO "StoreItemVariant" ("OrderId", "StoreItemId", "Variant") VALUES (?, ?, ?);`,
				order.Id, item.Id, variant)
			if err != nil {
				return wrapErr(err)
			}
		}
		for _, dip := range item.Dips {
			_, err = tx.Exec(
				`INSERT INTO "StoreItemDip" ("OrderId", "StoreItemId", "Dip") VALUES (?, ?, ?);`,
				order.Id, item.Id, dip)
			if err != nil {
				return wrapErr(err)
			}
		}
	}

	return wrapErr(tx.Commit())
}

func (t *Database) GetOrderList(orderListId string) (*model.OrderList, error) {
	var list model.OrderList
	var deadline sql.NullTime
	err := t.conn.QueryRow(`SELECT "Id", "Created", "Deadline" FROM "OrderList" WHERE "Id" = ?`, orderListId).
		Scan(&list.Id, &list.Created, &deadline)
	if err != nil {
		return nil, wrapErr(err)
	}
	if deadline.Valid {
		list.Deadline = &deadline.Time
	}
	return &list, nil
}

func (t *Database) GetOrders(orderListId string) ([]*model.Order, error) {
	rows, err := t.conn.Query(`
        SELECT o.Id, o.Created, o.Creator, o.EditKey, d.Name, d.Size 
        FROM "Order" o 
        LEFT JOIN "Drink" d ON d.Id = o.DrinkId 
        WHERE o.OrderListId = ?`, orderListId)
	if err != nil {
		return nil, wrapErr(err)
	}
	defer rows.Close()

	ordersMap := make(map[string]*model.Order)
	var orderIDs []interface{}

	for rows.Next() {
		var order model.Order
		var drinkName sql.NullString
		var drinkSize sql.NullInt64
		if err := rows.Scan(&order.Id, &order.Created, &order.Creator, &order.EditKey, &drinkName, &drinkSize); err != nil {
			return nil, wrapErr(err)
		}
		if drinkName.Valid {
			order.Drink = &model.Drink{Name: drinkName.String, Size: model.DrinkSize(drinkSize.Int64)}
		}
		order.StoreItems = []*model.StoreItem{}
		ordersMap[order.Id] = &order
		orderIDs = append(orderIDs, order.Id)
	}
	if err = rows.Err(); err != nil {
		return nil, wrapErr(err)
	}
	if len(orderIDs) == 0 {
		return []*model.Order{}, nil
	}

	query := `
        SELECT oi.OrderId, oi.StoreItemId, 
               GROUP_CONCAT(DISTINCT sv.Variant) as variants, 
               GROUP_CONCAT(DISTINCT sd.Dip) as dips
        FROM OrderItems oi
        LEFT JOIN StoreItemVariant sv ON sv.OrderId = oi.OrderId AND sv.StoreItemId = oi.StoreItemId
        LEFT JOIN StoreItemDip sd ON sd.OrderId = oi.OrderId AND sd.StoreItemId = oi.StoreItemId
        WHERE oi.OrderId IN (?` + strings.Repeat(",?", len(orderIDs)-1) + `)
        GROUP BY oi.OrderId, oi.StoreItemId`

	itemRows, err := t.conn.Query(query, orderIDs...)
	if err != nil {
		return nil, wrapErr(err)
	}
	defer itemRows.Close()

	for itemRows.Next() {
		var orderId, storeItemId string
		var variants, dips sql.NullString
		if err := itemRows.Scan(&orderId, &storeItemId, &variants, &dips); err != nil {
			return nil, wrapErr(err)
		}
		item := &model.StoreItem{Id: storeItemId}
		if variants.Valid {
			item.Variants = strings.Split(variants.String, ",")
		}
		if dips.Valid {
			item.Dips = strings.Split(dips.String, ",")
		}
		if order, ok := ordersMap[orderId]; ok {
			order.StoreItems = append(order.StoreItems, item)
		}
	}
	if err = itemRows.Err(); err != nil {
		return nil, wrapErr(err)
	}

	finalOrders := make([]*model.Order, 0, len(ordersMap))
	for _, order := range ordersMap {
		finalOrders = append(finalOrders, order)
	}
	return finalOrders, nil
}

func (t *Database) GetOrder(orderListId, orderId string) (*model.Order, error) {
	var (
		order     model.Order
		editKey   sql.NullString
		drinkName sql.NullString
		drinkSize sql.NullInt64
	)

	err := t.conn.QueryRow(`
		SELECT o.Id, o.Created, o.Creator, o.EditKey, d.Name, d.Size 
		FROM "Order" o 
		LEFT JOIN "Drink" d ON d.Id = o.DrinkId 
		WHERE o.OrderListId = ? AND o.Id = ?`, orderListId, orderId).
		Scan(&order.Id, &order.Created, &order.Creator, &editKey, &drinkName, &drinkSize)

	if err != nil {
		return nil, wrapErr(err)
	}

	if editKey.Valid {
		order.EditKey = editKey.String
	}
	if drinkName.Valid {
		order.Drink = &model.Drink{Name: drinkName.String, Size: model.DrinkSize(drinkSize.Int64)}
	}

	rows, err := t.conn.Query(`
		SELECT oi.StoreItemId, 
			   GROUP_CONCAT(DISTINCT sv.Variant) as variants, 
			   GROUP_CONCAT(DISTINCT sd.Dip) as dips
		FROM OrderItems oi
		LEFT JOIN StoreItemVariant sv ON sv.OrderId = oi.OrderId AND sv.StoreItemId = oi.StoreItemId
		LEFT JOIN StoreItemDip sd ON sd.OrderId = oi.OrderId AND sd.StoreItemId = oi.StoreItemId
		WHERE oi.OrderId = ?
		GROUP BY oi.StoreItemId`, orderId)
	if err != nil {
		return nil, wrapErr(err)
	}
	defer rows.Close()

	for rows.Next() {
		var storeItemId string
		var variants, dips sql.NullString
		if err := rows.Scan(&storeItemId, &variants, &dips); err != nil {
			return nil, wrapErr(err)
		}

		item := &model.StoreItem{Id: storeItemId}
		if variants.Valid {
			item.Variants = strings.Split(variants.String, ",")
		}
		if dips.Valid {
			item.Dips = strings.Split(dips.String, ",")
		}
		order.StoreItems = append(order.StoreItems, item)
	}

	return &order, nil
}

func (t *Database) UpdateOrder(orderListId string, order *model.Order) error {
	tx, err := t.conn.BeginTx(context.TODO(), nil)
	if err != nil {
		return wrapErr(err)
	}
	defer tx.Rollback()

	if _, err = tx.Exec(`DELETE FROM "OrderItems" WHERE "OrderId" = ?`, order.Id); err != nil {
		return wrapErr(err)
	}
	if _, err = tx.Exec(`DELETE FROM "StoreItemVariant" WHERE "OrderId" = ?`, order.Id); err != nil {
		return wrapErr(err)
	}
	if _, err = tx.Exec(`DELETE FROM "StoreItemDip" WHERE "OrderId" = ?`, order.Id); err != nil {
		return wrapErr(err)
	}
	if _, err = tx.Exec(`DELETE FROM "Drink" WHERE "Id" = (SELECT "DrinkId" FROM "Order" WHERE "Id" = ?)`, order.Id); err != nil {
		return wrapErr(err)
	}

	var drinkId sql.NullString
	if order.Drink != nil {
		drinkId.String = uuid.New().String()
		drinkId.Valid = true
		_, err = tx.Exec(`INSERT INTO "Drink" ("Id", "Name", "Size") VALUES (?, ?, ?);`, drinkId.String, order.Drink.Name, order.Drink.Size)
		if err != nil {
			return wrapErr(err)
		}
	}

	_, err = tx.Exec(
		`UPDATE "Order" SET "Creator" = ?, "DrinkId" = ? WHERE "Id" = ? AND "OrderListId" = ?`,
		order.Creator, drinkId, order.Id, orderListId)
	if err != nil {
		return wrapErr(err)
	}

	for _, item := range order.StoreItems {
		_, err = tx.Exec(`INSERT INTO "OrderItems" ("OrderId", "StoreItemId") VALUES (?, ?);`, order.Id, item.Id)
		if err != nil {
			return wrapErr(err)
		}
		for _, variant := range item.Variants {
			_, err = tx.Exec(`INSERT INTO "StoreItemVariant" ("OrderId", "StoreItemId", "Variant") VALUES (?, ?, ?);`, order.Id, item.Id, variant)
			if err != nil {
				return wrapErr(err)
			}
		}
		for _, dip := range item.Dips {
			_, err = tx.Exec(`INSERT INTO "StoreItemDip" ("OrderId", "StoreItemId", "Dip") VALUES (?, ?, ?);`, order.Id, item.Id, dip)
			if err != nil {
				return wrapErr(err)
			}
		}
	}

	return wrapErr(tx.Commit())
}

func (t *Database) DeleteOrderList(orderListId string) error {
	_, err := t.conn.Exec(`DELETE FROM "OrderList" WHERE "Id" = ?`, orderListId)
	return wrapErr(err)
}

func (t *Database) DeleteOrder(orderListId, orderId string) error {
	_, err := t.conn.Exec(`DELETE FROM "Order" WHERE "Id" = ? AND "OrderListId" = ?`, orderId, orderListId)
	return wrapErr(err)
}

func (t *Database) ClearAllData() error {
	tx, err := t.conn.BeginTx(context.TODO(), nil)
	if err != nil {
		return wrapErr(err)
	}
	defer tx.Rollback()

	tables := []string{"OrderItems", "StoreItemDip", "StoreItemVariant", "Drink", "Order", "OrderList"}
	for _, tbl := range tables {
		if _, err := tx.Exec(`DELETE FROM "` + tbl + `";`); err != nil {
			return wrapErr(err)
		}
	}

	return wrapErr(tx.Commit())
}

//Feedback\\

func (t *Database) CreateFeedback(feedback *model.Feedback) error {
	_, err := t.conn.Exec(
		`INSERT INTO "Feedback" ("Id", "Timestamp", "Type", "Message", "Page") VALUES (?, ?, ?, ?, ?);`,
		feedback.Id, feedback.Timestamp, feedback.Type, feedback.Message, feedback.Page)
	return wrapErr(err)
}

func (t *Database) GetAllFeedback() ([]*model.Feedback, error) {
	rows, err := t.conn.Query(`SELECT "Id", "Timestamp", "Type", "Message", "Page" FROM "Feedback" ORDER BY "Timestamp" DESC`)
	if err != nil {
		return nil, wrapErr(err)
	}
	defer rows.Close()

	var feedbacks []*model.Feedback
	for rows.Next() {
		var fb model.Feedback
		if err := rows.Scan(&fb.Id, &fb.Timestamp, &fb.Type, &fb.Message, &fb.Page); err != nil {
			return nil, wrapErr(err)
		}
		feedbacks = append(feedbacks, &fb)
	}
	return feedbacks, nil
}
