package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

type Permissions struct {
	CreateUser   bool `json:"create_user"`
	UpdateUser   bool `json:"update_user"`
	DeleteUser   bool `json:"delete_user"`
	GetUser      bool `json:"get_user"`
	CreateMember bool `json:"create_member"`
	UpdateMember bool `json:"update_member"`
	DeleteMember bool `json:"delete_member"`
	GetMember    bool `json:"get_member"`
	CreateChurch bool `json:"create_church"`
	UpdateChurch bool `json:"update_church"`
	DeleteChurch bool `json:"delete_church"`
	GetChurch    bool `json:"get_church"`
	CreateRole   bool `json:"create_role"`
	UpdateRole   bool `json:"update_role"`
	DeleteRole   bool `json:"delete_role"`
	GetRoles     bool `json:"get_roles"`
}

func (p Permissions) Value() (driver.Value, error) {
	return json.Marshal(p)
}

func (p *Permissions) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, p)
}

type Role struct {
	gorm.Model
	Name        string `gorm:"type:varchar(255)"`
	Users       []User
	Permissions Permissions `gorm:"type:jsonb"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
