package database

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"strings"

	_ "github.com/glebarez/go-sqlite"
	"github.com/google/uuid"
	"github.com/pressly/goose/v3"
	"github.com/zekrotja/hermans/pkg/model"
)

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
		`INSERT INTO "OrderList" ("Id", "Created") VALUES (?, ?);`,
		list.Id, list.Created)
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
		`INSERT INTO "Order" ("Id", "Created", "Creator", "OrderListId", "StoreItemId", "DrinkId", "EditKey")
		 VALUES (?, ?, ?, ?, ?, ?, ?);`,
		order.Id, order.Created, order.Creator, orderListId, order.StoreItem.Id, drinkId, order.EditKey)
	if err != nil {
		return wrapErr(err)
	}
	for _, variant := range order.StoreItem.Variants {
		_, err = tx.Exec(
			`INSERT INTO "StoreItemVariant" ("OrderId", "Variant") VALUES (?, ?);`,
			order.Id, variant)
		if err != nil {
			return wrapErr(err)
		}
	}
	for _, dip := range order.StoreItem.Dips {
		_, err = tx.Exec(
			`INSERT INTO "StoreItemDip" ("OrderId", "Dip") VALUES (?, ?);`,
			order.Id, dip)
		if err != nil {
			return wrapErr(err)
		}
	}
	return wrapErr(tx.Commit())
}

func (t *Database) GetOrderList(orderListId string) (*model.OrderList, error) {
	var list model.OrderList
	err := t.conn.QueryRow(`SELECT "Id", "Created" FROM "OrderList" WHERE "Id" = ?`, orderListId).
		Scan(&list.Id, &list.Created)
	if err != nil {
		return nil, wrapErr(err)
	}
	return &list, nil
}

func (t *Database) GetOrders(orderListId string) ([]*model.Order, error) {
	rows, err := t.conn.Query(`
		SELECT "Order"."Id", "Created", "Creator", "StoreItemId", "Drink"."Name", "Drink"."Size",
			   GROUP_CONCAT(DISTINCT "StoreItemVariant"."Variant") AS "variants",
			   GROUP_CONCAT(DISTINCT "StoreItemDip"."Dip") AS "dips"
		FROM "Order"
		LEFT JOIN "Drink" ON "Drink"."Id" = "Order"."DrinkId"
		LEFT JOIN "StoreItemVariant" ON "StoreItemVariant"."OrderId" = "Order"."Id"
		LEFT JOIN "StoreItemDip" ON "StoreItemDip"."OrderId" = "Order"."Id"
		WHERE "Order"."OrderListId" = ?
		GROUP BY "Order"."Id"`,
		orderListId)
	if errors.Is(err, sql.ErrNoRows) {
		return []*model.Order{}, nil
	}
	if err != nil {
		return nil, wrapErr(err)
	}
	var orders []*model.Order
	for rows.Next() {
		var (
			order     model.Order
			drinkName sql.NullString
			drinkSize sql.NullInt64
			storeItem model.StoreItem
			variants  sql.NullString
			dips      sql.NullString
		)
		err = rows.Scan(&order.Id, &order.Created, &order.Creator, &storeItem.Id, &drinkName, &drinkSize, &variants, &dips)
		if err != nil {
			return nil, wrapErr(err)
		}
		order.StoreItem = &storeItem
		if drinkName.Valid {
			order.Drink = &model.Drink{
				Name: drinkName.String,
				Size: model.DrinkSize(drinkSize.Int64),
			}
		}
		if variants.Valid {
			order.StoreItem.Variants = strings.Split(variants.String, ",")
		}
		if dips.Valid {
			order.StoreItem.Dips = strings.Split(dips.String, ",")
		}
		orders = append(orders, &order)
	}
	return orders, nil
}

func (t *Database) DeleteOrderList(orderListId string) error {
	_, err := t.conn.Exec(`DELETE FROM "OrderList" WHERE "Id" = ?`, orderListId)
	return wrapErr(err)
}

func (t *Database) GetOrder(orderListId, orderId string) (*model.Order, error) {
	var order model.Order
	order.StoreItem = new(model.StoreItem)
	err := t.conn.QueryRow(
		`SELECT "Id", "Created", "Creator", "StoreItemId", "EditKey" FROM "Order" WHERE "OrderListId" = ? AND "Id" = ?`,
		orderListId, orderId).
		Scan(&order.Id, &order.Created, &order.Creator, &order.StoreItem.Id, &order.EditKey)
	if err != nil {
		return nil, wrapErr(err)
	}
	return &order, nil
}

func (t *Database) UpdateOrder(orderListId string, order *model.Order) error {
	tx, err := t.conn.BeginTx(context.TODO(), nil)
	if err != nil {
		return wrapErr(err)
	}
	defer tx.Rollback()

	if _, err = tx.Exec(`DELETE FROM "StoreItemVariant" WHERE "OrderId" = ?`, order.Id); err != nil {
		return wrapErr(err)
	}
	if _, err = tx.Exec(`DELETE FROM "StoreItemDip" WHERE "OrderId" = ?`, order.Id); err != nil {
		return wrapErr(err)
	}

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
		`UPDATE "Order" SET "Creator" = ?, "StoreItemId" = ?, "DrinkId" = ? WHERE "Id" = ? AND "OrderListId" = ?`,
		order.Creator, order.StoreItem.Id, drinkId, order.Id, orderListId)
	if err != nil {
		return wrapErr(err)
	}

	for _, variant := range order.StoreItem.Variants {
		_, err = tx.Exec(`INSERT INTO "StoreItemVariant" ("OrderId", "Variant") VALUES (?, ?);`, order.Id, variant)
		if err != nil {
			return wrapErr(err)
		}
	}
	for _, dip := range order.StoreItem.Dips {
		_, err = tx.Exec(`INSERT INTO "StoreItemDip" ("OrderId", "Dip") VALUES (?, ?);`, order.Id, dip)
		if err != nil {
			return wrapErr(err)
		}
	}

	return wrapErr(tx.Commit())
}

func (t *Database) DeleteOrder(orderListId, orderId string) error {
	_, err := t.conn.Exec(`DELETE FROM "Order" WHERE "Id" = ? AND "OrderListId" = ?`, orderId, orderListId)
	return wrapErr(err)
}
