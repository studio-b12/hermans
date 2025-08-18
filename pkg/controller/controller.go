package controller

import (
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

	if data != nil {
		return data, nil
	}

	return t.Scrape()
}

func (t *Controller) CreateOrderList() (*model.OrderList, error) {
	list := model.OrderList{
		Id:      uuid.New().String(),
		Created: time.Now(),
	}

	err := t.db.CreateOrderList(&list)
	if err != nil {
		return nil, err
	}

	return &list, nil
}

func (t *Controller) CreateOrder(orderListId string, order *model.Order) (*model.Order, error) {
	order.Id = uuid.New().String()
	order.Created = time.Now()

	err := t.validator.Struct(order)
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
