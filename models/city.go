package models

import (
	"time"

	"gorm.io/gorm"
)

type City struct {
	gorm.Model
	ID        uint
	StateId   int `gorm:"notnull"`
	Name      string
	Churches  []Church
	CreatedAt time.Time
	UpdatedAt time.Time
}
