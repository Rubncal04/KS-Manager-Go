package models

import (
	"time"

	"gorm.io/gorm"
)

type Church struct {
	gorm.Model
	ID        uint
	CityId    int `gorm:"notnull"`
	CountryId int `gorm:"notnull"`
	StateId   int `gorm:"notnull"`
	Members   []Member
	Name      string
	Address   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
