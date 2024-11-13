package models

import (
	"time"

	"gorm.io/gorm"
)

type Offering struct {
	gorm.Model
	ChurchId         int
	Church           Church `gorm:"foreignKey:ChurchId;notnull"`
	CategoryId       int
	Category         Category `gorm:"foreignKey:CategoryId;notnull"`
	WorshipServiceId int
	WorshipService   WorshipService `gorm:"foreignKey:WorshipServiceId;notnull"`
	Name             string         `gorm:"type:varchar(255);notnull"`
	Value            int            `gorm:"notnull"`
	Date             time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
