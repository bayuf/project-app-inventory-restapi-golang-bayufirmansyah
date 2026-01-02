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
	MinStock   int
	Price      decimal.Decimal
	Created_At time.Time
	Updated_At time.Time
	Deleted_At time.Time
}

type ItemDetail struct {
	ID           uuid.UUID
	Name         string
	CategoryID   int
	CategoryName string
	RackID       int
	RackCode     string
	RackName     string
	Location     string
	Stock        int
	Price        decimal.Decimal
	Created_At   time.Time
	Updated_At   time.Time
}
