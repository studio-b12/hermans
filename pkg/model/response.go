package model

import "time"

type CreateOrderResponse struct {
	Id         string       `json:"id"`
	Created    time.Time    `json:"created"`
	Creator    string       `json:"creator"`
	StoreItems []*StoreItem `json:"store_items"`
	Drink      *Drink       `json:"drink"`
	EditKey    string       `json:"editKey"`
}

type GetOrderListResponse struct {
	Id       string     `json:"id"`
	Created  time.Time  `json:"created"`
	Deadline *time.Time `json:"deadline"`
	Orders   []*Order   `json:"orders"`
}
