package model

import "time"

type Category struct {
	ID          int
	Name        string
	Description string
	Created_At  time.Time
	Updated_At  time.Time
	Deleted_At  time.Time
}
