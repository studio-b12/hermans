package controller

import (
	"github.com/zekrotja/hermans/pkg/model"
)

type Database interface {
	CreateOrderList(list *model.OrderList) error
	GetOrderList(id string) (*model.OrderList, error)
	DeleteOrderList(id string) error
	CreateOrder(orderListId string, order *model.Order) error
	GetOrders(orderListId string) ([]*model.Order, error)
	GetOrder(orderListId, orderId string) (*model.Order, error)
	UpdateOrder(orderListId string, order *model.Order) error
	DeleteOrder(orderListId, orderId string) error
	ClearAllData() error //Debug
}
