package api

import (
	"github.com/zekrotja/hermans/pkg/model"
	"github.com/zekrotja/hermans/pkg/scraper"
)

type Controller interface {
	CreateOrder(orderListId string, order *model.Order) (*model.Order, error)
	CreateOrderList() (*model.OrderList, error)
	GetOrders(orderListId string) ([]*model.Order, error)
	GetScrapedData() (*scraper.Data, error)
	Scrape() (*scraper.Data, error)
}
