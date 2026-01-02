package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type ItemAdd struct {
	Name       string          `json:"name" validate:"required"`
	CategoryId int             `json:"category_id" validate:"required,gt=0"`
	RackId     int             `json:"rack_id" validate:"required,gt=0"`
	Stock      int             `json:"stock" validate:"required,gt=0"`
	MinStock   int             `json:"min_stock" validate:"required,gt=0"`
	Price      decimal.Decimal `json:"price" validate:"required,decimal_gt_zero"`
}

type ItemResponse struct {
	ID         uuid.UUID       `json:"id"`
	Name       string          `json:"name"`
	CategoryId int             `json:"category_id"`
	Category   string          `json:"category"`
	RackID     int             `json:"rack_id"`
	Rack       string          `json:"rack"`
	Location   string          `json:"location"`
	Stock      int             `json:"stock"`
	Price      decimal.Decimal `json:"price"`
	Created_At time.Time       `json:"created_at"`
	Updated_At time.Time       `json:"updated_at"`
}
