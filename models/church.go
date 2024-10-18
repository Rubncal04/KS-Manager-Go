package models

import (
	"time"
)

type Church struct {
	ID        uint
	Name      string
	Address   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
