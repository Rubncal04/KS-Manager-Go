package models

import (
	"database/sql/driver"
	"encoding/json"

	"gorm.io/gorm"
)

type StringArray []string

// Scan implementa la interfaz sql.Scanner
func (sa *StringArray) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}

	return json.Unmarshal(bytes, sa)
}

// Value implementa la interfaz driver.Valuer
func (sa StringArray) Value() (driver.Value, error) {
	if sa == nil {
		return nil, nil
	}
	return json.Marshal(sa)
}

type Member struct {
	ID                   uint
	ChurchId             int `gorm:"notnull"`
	Name                 string
	LastName             string
	IdentificationNumber string `gorm:"uniqueIndex"`
	Address              string
	Email                string
	Birthday             string
	BaptizedBy           string
	BaptizedOn           string
	HolySpiritOn         string
	Position             string
	NumChildren          int
	ChildrenNames        StringArray `gorm:"type:jsonb"`
	PartnerName          string
	Degree               string
	Profession           string
}

func (m *Member) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ChildrenNames == nil {
		m.ChildrenNames = StringArray{}
	}
	return
}
