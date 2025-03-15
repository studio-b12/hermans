package model

import (
	"time"
)

type OrderList struct {
	Id      string    `json:"id"`
	Created time.Time `json:"created"`
	Orders  []Order   `json:"orders"`
}

type StoreItem struct {
	Id       string   `json:"id"`
	Variants []string `json:"variants"`
	Dips     []string `json:"dips"`
}

type Drink struct {
	Name string `json:"name"`
}

type Order struct {
	Id        string     `json:"id"`
	Created   time.Time  `json:"created"`
	StoreItem *StoreItem `json:"store_item"`
	Drink     *Drink     `json:"drink"`
}
