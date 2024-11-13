package models

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	Name      string `gorm:"type:varchar(255)"`
	Offerings []Offering
	CreatedAt time.Time
	UpdatedAt time.Time
}
