package api

import (
	"time"

	"github.com/zekrotja/hermans/pkg/model"
	"github.com/zekrotja/hermans/pkg/scraper"
)

type Controller interface {
	Scrape() (*scraper.Data, error)
	GetScrapedData() (*scraper.Data, error)
	CreateOrderList(deadline *time.Time) (*model.OrderList, error)
	GetOrders(orderListId string) (*model.OrderList, error)
	DeleteOrderList(orderListId string) error
	CreateOrder(orderListId string, order *model.Order) (*model.Order, error)
	UpdateOrder(orderListId, orderId, editKey string, updatedOrder *model.Order) (*model.Order, error)
	DeleteOrder(orderListId, orderId, editKey string) error
	GetOrder(orderListId, orderId string) (*model.Order, error)
	GetOrderList(orderListId string) (*model.OrderList, error)
	ClearAllData() error
}
