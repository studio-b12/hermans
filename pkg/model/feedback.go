package model

import "time"

type Feedback struct {
	Id        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Type      string    `json:"type" validate:"required"`
	Message   string    `json:"message" validate:"required"`
	Page      string    `json:"page" validate:"required"`
}
