package dto

import "time"

type RackAdd struct {
	WarehouseId int    `json:"warehouse_id" validate:"required,gt=0"`
	RackCode    string `json:"rack_code" validate:"required"`
	Description string `json:"description" validated:"required"`
}

type RackResponse struct {
	ID                int       `json:"id"`
	Code              string    `json:"rack_code"`
	Description       string    `json:"description"`
	WarehouseName     string    `json:"warehouse_name"`
	WarehouseLocation string    `json:"warehouse_location"`
	Created_At        time.Time `json:"create_at"`
	Updated_At        time.Time `json:"update_at"`
}
