package models

import (
	"time"

	"gorm.io/gorm"
)

type Country struct {
	gorm.Model
	ID        uint
	Name      string
	States    []State
	CreatedAt time.Time
	UpdatedAt time.Time
}
