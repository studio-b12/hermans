package model

import "time"

type CreateListPayload struct {
	Deadline *time.Time `json:"deadline"`
}

type UpdateOrderPayload struct {
	Order
	EditKey string `json:"editKey"`
}

type DeleteOrderPayload struct {
	EditKey string `json:"editKey"`
}
