package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Item struct {
	ID         uuid.UUID
	Name       string
	CategoryID int
	RackID     int
	Stock      int
	Price      decimal.Decimal
	Create_At  time.Time
	Update_At  time.Time
	Deleted_At time.Time
}
