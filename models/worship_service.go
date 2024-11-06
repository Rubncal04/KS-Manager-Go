package models

import (
	"time"

	"gorm.io/gorm"
)

type WeekDay string

const (
	Monday    WeekDay = "Monday"
	Tuesday   WeekDay = "Tuesday"
	Wednesday WeekDay = "Wednesday"
	Thursday  WeekDay = "Thursday"
	Friday    WeekDay = "Friday"
	Saturday  WeekDay = "Saturday"
	Sunday    WeekDay = "Sunday"
)

type WorshipService struct {
	gorm.Model
	ChurchId  int     `gorm:"notnull"`
	Church    Church  `gorm:"foreignKey:ChurchId"`
	Name      string  `gorm:"type:varchar(255);notnull"`
	Day       WeekDay `gorm:"type:varchar(255)"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (w *WorshipService) IsValidWeekday() bool {
	switch w.Day {
	case Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday:
		return true
	}

	return false
}
