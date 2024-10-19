package models

import (
	"time"

	"gorm.io/gorm"
)

type State struct {
	gorm.Model
	ID        uint
	CountryId int `gorm:"notnull"`
	Name      string
	Cities    []City
	CreatedAt time.Time
	UpdatedAt time.Time
}
