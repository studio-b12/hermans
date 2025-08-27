package model

import (
	"time"
)

type DrinkSize int

const (
	DrinkSizeSmall DrinkSize = 0
	DrinkSizeLarge DrinkSize = 1
)

type OrderList struct {
	Id       string     `json:"id"`
	Created  time.Time  `json:"created"`
	Orders   []*Order   `json:"orders"`
	Deadline *time.Time `json:"deadline,omitempty"`
}

type StoreItem struct {
	Id       string   `json:"id" validate:"required"`
	Variants []string `json:"variants" validate:"unique"`
	Dips     []string `json:"dips" validate:"unique"`
}

type Drink struct {
	Name string    `json:"name" validate:"required"`
	Size DrinkSize `json:"size" validate:"gte=0,lte=1"`
}

type Order struct {
	Id         string       `json:"id"`
	Created    time.Time    `json:"created"`
	Creator    string       `json:"creator" validate:"required"`
	StoreItems []*StoreItem `json:"store_items" validate:"required,min=1"`
	Drink      *Drink       `json:"drink"`
	EditKey    string       `json:"-"`
}
