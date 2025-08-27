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
	GetOrder(orderListId, orderId string) (*model.Order, error)
	UpdateOrder(orderListId string, order *model.Order) error
	DeleteOrder(orderListId, orderId string) error
	ClearAllData() error //debug
	//Feedback\\
	CreateFeedback(feedback *model.Feedback) error
	GetAllFeedback() ([]*model.Feedback, error)
}
