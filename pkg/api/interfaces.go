package api

import (
	"time"

	"github.com/zekrotja/hermans/pkg/model"
	"github.com/zekrotja/hermans/pkg/scraper"
)

type Controller interface {
	GetScrapedData() (*scraper.Data, error)
	CreateOrderList(deadline *time.Time) (*model.OrderList, error)
	GetOrderList(orderListId string) (*model.OrderList, error)
	DeleteOrderList(orderListId string) error
	CreateOrder(orderListId string, order *model.Order) (*model.Order, error)
	UpdateOrder(orderListId, orderId, editKey string, updatedOrder *model.Order) (*model.Order, error)
	DeleteOrder(orderListId, orderId, editKey string) error
	GetOrders(orderListId string) ([]*model.Order, error)
	GetOrder(orderListId, orderId string) (*model.Order, error)
	ClearAllData() error
}
