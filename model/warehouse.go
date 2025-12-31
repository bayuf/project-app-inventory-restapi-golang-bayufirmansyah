package model

import "time"

type Warehouse struct {
	ID         int
	Name       string
	Location   string
	IsActive   bool
	Created_at time.Time
	Updated_at time.Time
}
