package controller

import (
	"errors"
	"path/filepath"
	"slices"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/studio-b12/elk"
	"github.com/zekrotja/hermans/pkg/cache"
	"github.com/zekrotja/hermans/pkg/model"
	"github.com/zekrotja/hermans/pkg/scraper"
)

type Controller struct {
	db Database

	validator *validator.Validate

	scrapeCache *cache.LocalCache[*scraper.Data]
}

func New(cacheDir string, db Database) (*Controller, error) {
	scrapeDb, err := cache.OpenLocalCache[*scraper.Data](filepath.Join(cacheDir, "scrape_data.msgpack"))
	if err != nil {
		return nil, err
	}

	t := &Controller{
		db:          db,
		scrapeCache: scrapeDb,
		validator:   validator.New(validator.WithRequiredStructEnabled()),
	}
	return t, nil
}

func (t *Controller) Scrape() (*scraper.Data, error) {
	data, err := scraper.ScrapeAll()
	if err != nil {
		return nil, err
	}

	err = t.scrapeCache.Store(data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (t *Controller) GetScrapedData() (*scraper.Data, error) {
	data, err := t.scrapeCache.Load()
	if err != nil {
		return nil, err
	}

	if data == nil {
		data, err = t.Scrape()
		if err != nil {
			return nil, err
		}
	}

	surpriseCat := []*scraper.Category{
		{
			Id:   "__etc",
			Name: "Etc",
			Items: []*scraper.StoreItem{
				{
					Id:          "__surprise",
					Title:       "🎉 Überrasch mich 🎉",
					Description: "Die bestellende Person sucht sich etwas für dich aus 😎",
					Variants: []*scraper.Variant{
						{
							Name:        "vegetarisch",
							Description: "Vegetarisch",
						},
						{
							Name:        "ohne zwiebeln",
							Description: "one Zwiebeln (wenn vorhanden)",
						},
					},
				},
			},
		},
	}
	data.Categories = append(surpriseCat, data.Categories...)

	return data, nil
}

func (t *Controller) CreateOrderList(deadline *time.Time) (*model.OrderList, error) {
	list := model.OrderList{
		Id:       uuid.New().String(),
		Created:  time.Now(),
		Deadline: deadline,
	}
	err := t.db.CreateOrderList(&list)

	if err != nil {
		return nil, err
	}
	return &list, nil
}

// Debug
func (t *Controller) ClearAllData() error {
	return t.db.ClearAllData()
}

var ErrDeadlineExceeded = errors.New("deadline for this order list has been exceeded")

func (t *Controller) CreateOrder(orderListId string, order *model.Order) (*model.Order, error) {
	list, err := t.db.GetOrderList(orderListId)
	if err != nil {
		return nil, err
	}
	if list.Deadline != nil && time.Now().After(*list.Deadline) {
		return nil, ErrDeadlineExceeded
	}

	order.Id = uuid.New().String()
	order.Created = time.Now()
	order.EditKey = uuid.New().String()
	err = t.validator.Struct(order)
	if err != nil {
		return nil, err
	}

	item, ok, err := t.getStoreItem(order.StoreItem.Id)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, elk.NewErrorf(ErrInvalidStoreItem, "invalid store item ID: %s", order.StoreItem.Id)
	}

	var invalidVariants ListError
	for _, variant := range order.StoreItem.Variants {
		if !item.VariantsContain(variant) {
			invalidVariants = append(invalidVariants, variant)
		}
	}
	if len(invalidVariants) > 0 {
		return nil, elk.Wrap(ErrInvalidVariants, invalidVariants, "invalid variants")
	}

	var invalidDips ListError
	for _, dip := range order.StoreItem.Dips {
		if !slices.Contains(item.Dips, dip) {
			invalidDips = append(invalidDips, dip)
		}
	}
	if len(invalidDips) > 0 {
		return nil, elk.Wrap(ErrInvalidDips, invalidDips, "invalid dips")
	}

	err = t.db.CreateOrder(orderListId, order)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (t *Controller) GetOrders(orderListId string) (*model.OrderList, error) {
	orderList, err := t.db.GetOrderList(orderListId)
	if err != nil {
		return nil, err
	}

	orderList.Orders, err = t.db.GetOrders(orderListId)
	if err != nil {
		return nil, err
	}

	return orderList, nil
}

func (t *Controller) DeleteOrderList(orderListId string) error {
	err := t.db.DeleteOrderList(orderListId)
	if err != nil {
		return err
	}

	return nil
}

func (t *Controller) GetOrder(orderListId, orderId string) (*model.Order, error) {
	return t.db.GetOrder(orderListId, orderId)
}

// UpdateOrder bearbeitet eine Bestellung nach der Prüfung des geheimen Schlüssels.
func (t *Controller) UpdateOrder(orderListId, orderId, editKey string, updatedOrder *model.Order) (*model.Order, error) {
	order, err := t.db.GetOrder(orderListId, orderId)
	if err != nil {
		return nil, err
	}
	if order.EditKey != editKey {
		return nil, errors.New("invalid edit key: access denied")
	}

	order.Creator = updatedOrder.Creator
	order.StoreItem = updatedOrder.StoreItem
	order.Drink = updatedOrder.Drink

	if err := t.db.UpdateOrder(orderListId, order); err != nil {
		return nil, err
	}
	return order, nil
}

// DeleteOrder löscht eine Bestellung nach der Prüfung von dem geheimen Schlüssel.
func (t *Controller) DeleteOrder(orderListId, orderId, editKey string) error {
	order, err := t.db.GetOrder(orderListId, orderId)
	if err != nil {
		return err
	}
	if order.EditKey != editKey {
		return errors.New("invalid edit key: access denied")
	}
	return t.db.DeleteOrder(orderListId, orderId)
}

func (t *Controller) getStoreItem(id string) (si *scraper.StoreItem, ok bool, err error) {
	data, err := t.GetScrapedData()
	if err != nil {
		return nil, false, err
	}

	for _, category := range data.Categories {
		for _, si = range category.Items {
			if si.Id == id {
				return si, true, nil
			}
		}
	}

	return nil, false, nil
}
