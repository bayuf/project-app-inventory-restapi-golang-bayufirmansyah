package model

import "time"

type Rack struct {
	ID                int
	WarehouseID       int
	WarehouseName     string
	WarehouseLocation string
	Code              string
	Description       string
	Created_At        time.Time
	Updated_At        time.Time
}
