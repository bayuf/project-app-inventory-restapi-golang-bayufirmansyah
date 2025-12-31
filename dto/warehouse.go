package dto

import "time"

type Warehouse struct {
	ID       int    `validate:"required,gt=0"`
	Name     string `validate:"required"`
	Location string `validate:"required"`
}

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
