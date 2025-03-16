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
		`INSERT INTO "OrderList" ("Id", "Created") VALUES (?, ?);`,
		list.Id, list.Created)
	return err
}

func (t *Database) CreateOrder(orderListId string, order *model.Order) error {
	tx, err := t.conn.BeginTx(context.TODO(), nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var drinkId string

	if order.Drink != nil {
		drinkId = uuid.New().String()
		_, err = tx.Exec(
			`INSERT INTO "Drink" ("Id", "Name", "Size") VALUES (?, ?, ?);`,
			drinkId, order.Drink.Name, order.Drink.Size)
		if err != nil {
			return err
		}
	}

	_, err = tx.Exec(
		`INSERT INTO "Order" ("Id", "Created", "OrderListId", "StoreItemId", "DrinkId") 
		VALUES (?, ?, ?, ?, ?);`,
		order.Id, order.Created, orderListId, order.StoreItem.Id, drinkId)
	if err != nil {
		return err
	}

	for _, variant := range order.StoreItem.Variants {
		_, err = tx.Exec(
			`INSERT INTO "StoreItemVariant" ("OrderId", "Variant") VALUES (?, ?);`,
			order.Id, variant)
		if err != nil {
			return err
		}
	}

	for _, dip := range order.StoreItem.Dips {
		_, err = tx.Exec(
			`INSERT INTO "StoreItemDip" ("OrderId", "Dip") VALUES (?, ?);`,
			order.Id, dip)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (t *Database) GetOrders(orderListId string) ([]*model.Order, error) {

	rows, err := t.conn.Query(`
		SELECT "Id", "Created", "StoreItemId", "Name", "Size", "variants",
			GROUP_CONCAT("StoreItemDip"."Dip") AS "dips" FROM (
				SELECT "Order"."Id", "Created", "StoreItemId", "Drink"."Name", "Drink"."Size",
					GROUP_CONCAT("StoreItemVariant"."Variant") AS "variants"
				FROM "Order"
				LEFT JOIN "Drink" ON "Drink"."Id" = "Order"."DrinkId"
				LEFT JOIN "StoreItemVariant" ON "StoreItemVariant"."OrderId" = "Order"."Id"
				WHERE "Order"."OrderListId" = ?
		)
		LEFT JOIN "StoreItemDip" ON "StoreItemDip"."OrderId" = "Id"`,
		orderListId)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var orders []*model.Order
	for rows.Next() {
		var (
			order     model.Order
			drinkName *string
			drinkSize *model.DrinkSize
			storeItem model.StoreItem
			variants  *string
			dips      *string
		)
		err = rows.Scan(&order.Id, &order.Created, &storeItem.Id, &drinkName, &drinkSize, &variants, &dips)
		if err != nil {
			return nil, err
		}
		order.StoreItem = &storeItem
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

func (t *Database) DeleteOrder(orderId string) error {
	_, err := t.conn.Exec(`DELETE FROM "Order" WHERE "Id" = ?`, orderId)
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}
	return err
}
