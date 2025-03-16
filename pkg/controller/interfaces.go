package controller

import (
	"github.com/zekrotja/hermans/pkg/model"
)

type Database interface {
	CreateOrderList(list *model.OrderList) error
	CreateOrder(orderListId string, order *model.Order) error
	GetOrderList(orderListId string) (*model.OrderList, error)
	GetOrders(orderListId string) ([]*model.Order, error)
	DeleteOrderList(orderListId string) error
	DeleteOrder(orderId string) error
}
