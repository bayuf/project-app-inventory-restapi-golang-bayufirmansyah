package dto

import "time"

type WarehouseAdd struct {
	Name     string `validate:"required"`
	Location string `validate:"required"`
}

type WarehouseResponse struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Location   string    `json:"location"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
}
