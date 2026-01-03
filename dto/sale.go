package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type NewSale struct {
	ItemID   uuid.UUID `json:"item_id" validate:"required"`
	Quantity int       `json:"quantity" validate:"required,gt=0"`
}

type SaleResponse struct {
	ID          uuid.UUID       `json:"id"`
	TotalAmount decimal.Decimal `json:"total_amount"`
	Status      string          `json:"status"`
	Created_At  time.Time       `json:"created_at"`
}

type SaleDetailResponse struct {
	ID          uuid.UUID       `json:"id"`
	ItemID      uuid.UUID       `json:"item_id"`
	ItemName    string          `json:"item_name"`
	Quantity    int             `json:"quantity"`
	Price       decimal.Decimal `json:"price"`
	TotalAmount decimal.Decimal `json:"total_amount"`
	Status      string          `json:"status"`
	Created_At  time.Time       `json:"created_at"`
}

type SalesUpdate struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	TotalAmount decimal.Decimal
	Status      string
}

type SalesItemsUpdate struct {
	ID       int
	SaleID   uuid.UUID
	ItemID   uuid.UUID
	Quantity int
	Price    decimal.Decimal
	SubTotal decimal.Decimal
}

type StockUpdateFromSale struct {
	ID    uuid.UUID
	Stock int
}
