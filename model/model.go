package model

import "time"

type Model struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Created_At time.Time `json:"create_at"`
	Updated_At time.Time `json:"updated_at"`
	Deleted_At time.Time `json:"deleted_at"`
}
