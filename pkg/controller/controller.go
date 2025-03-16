package controller

import (
	"path/filepath"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/zekrotja/hermans/pkg/cache"
	"github.com/zekrotja/hermans/pkg/model"
	"github.com/zekrotja/hermans/pkg/scraper"
)

type Controller struct {
	db Database

	scrapeCache *cache.LocalCache[*scraper.Data]
	validator   *validator.Validate
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

	err = t.db.CreateOrder(orderListId, order)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (t *Controller) GetOrders(orderListId string) ([]*model.Order, error) {
	orders, err := t.db.GetOrders(orderListId)
	if err != nil {
		return nil, err
	}

	return orders, nil
}
