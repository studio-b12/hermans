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
	err = goose.SetDialect("sqlite")
	if err != nil {
		return nil, err
	}

	err = goose.Up(conn, "migrations")
	if err != nil {
		return nil, err
	}

	t := &Database{conn: conn}
	return t, nil
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
		`INSERT INTO "Order" ("Id", "Created", "Creator", "OrderListId", "StoreItemId", "DrinkId", "EditKey")
		 VALUES (?, ?, ?, ?, ?, ?, ?);`,
		order.Id, order.Created, order.Creator, orderListId, order.StoreItem.Id, drinkId, order.EditKey) // Der Wert für EditKey wird hier übergeben
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
	var deadline sql.NullTime // Spezieller Typ für optionale Zeitstempel

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

// Einzelne Orders Anzeigen
func (t *Database) GetOrder(orderListId, orderId string) (*model.Order, error) {
	var (
		order   model.Order
		editKey *string
	)
	order.StoreItem = new(model.StoreItem)

	err := t.conn.QueryRow(
		`SELECT "Id", "Created", "Creator", "StoreItemId", "EditKey" FROM "Order" WHERE "OrderListId" = ? AND "Id" = ?`,
		orderListId, orderId).
		Scan(&order.Id, &order.Created, &order.Creator, &order.StoreItem.Id, &editKey)

	if err != nil {
		return nil, wrapErr(err)
	}

	if editKey != nil {
		order.EditKey = *editKey
	}

	return &order, nil
}

// Order Updaten
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
	if _, err = tx.Exec(`DELETE FROM "Drink" WHERE "Id" = (SELECT "DrinkId" FROM "Order" WHERE "Id" = ?)`, order.Id); err != nil {
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

	if err != nil {
		return nil, wrapErr(err)
	}
	defer rows.Close()

	var orders []*model.Order
	for rows.Next() {
		var (
			order     model.Order
			drinkName *string
			drinkSize *model.DrinkSize
			variants  *string
			dips      *string
		)
		order.StoreItem = new(model.StoreItem)

		err = rows.Scan(&order.Id, &order.Created, &order.Creator, &order.StoreItem.Id, &drinkName, &drinkSize, &variants, &dips)
		if err != nil {
			return nil, wrapErr(err)
		}

		if drinkName != nil {
			order.Drink = &model.Drink{
				Name: *drinkName,
				Size: *drinkSize,
			}
		}
		if variants != nil {
			order.StoreItem.Variants = strings.Split(*variants, ",")
		}
		if dips != nil {
			order.StoreItem.Dips = strings.Split(*dips, ",")
		}
		orders = append(orders, &order)
	}

	return orders, nil
}

func (t *Database) DeleteOrderList(orderListId string) error {
	_, err := t.conn.Exec(`DELETE FROM "OrderList" WHERE "Id" = ?`, orderListId)
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}
	return err
}

func (t *Database) DeleteOrder(orderListId, orderId string) error {
	res, err := t.conn.Exec(`DELETE FROM "Order" WHERE "Id" = ? AND "OrderListId" = ?`, orderId, orderListId)
	if err != nil {
		return wrapErr(err)
	}

	rowsAffected, err := res.RowsAffected()
	if err == nil && rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return wrapErr(err)
}

//DEBUGGING

// alle Daten löschen
func (t *Database) ClearAllData() error {
	// Beginne eine Transaktion, um sicherzustellen, dass alles oder nichts gelöscht wird
	tx, err := t.conn.BeginTx(context.TODO(), nil)
	if err != nil {
		return wrapErr(err)
	}
	defer tx.Rollback() // Bricht ab, wenn etwas schiefgeht

	// Lösche zuerst die abhängigen Daten (in der richtigen Reihenfolge)
	if _, err := tx.Exec(`DELETE FROM "StoreItemDip";`); err != nil {
		return wrapErr(err)
	}
	if _, err := tx.Exec(`DELETE FROM "StoreItemVariant";`); err != nil {
		return wrapErr(err)
	}
	if _, err := tx.Exec(`DELETE FROM "Drink";`); err != nil {
		return wrapErr(err)
	}
	if _, err := tx.Exec(`DELETE FROM "Order";`); err != nil {
		return wrapErr(err)
	}
	if _, err := tx.Exec(`DELETE FROM "OrderList";`); err != nil {
		return wrapErr(err)
	}

	// Wenn alles gut gegangen ist, bestätige die Änderungen
	return wrapErr(tx.Commit())
}
