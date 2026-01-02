package dto

import (
	"github.com/shopspring/decimal"
)

type ItemAdd struct {
	Name       string          `json:"name" validate:"required"`
	CategoryId int             `json:"category_id" validate:"required,gt=0"`
	RackId     int             `json:"rack_id" validate:"required,gt=0"`
	Stock      int             `json:"stock" validate:"required,gt=0"`
	Price      decimal.Decimal `json:"price" validate:"required,gt=0"`
}
