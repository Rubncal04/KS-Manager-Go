package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ChurchId  int `gorm:"notnull"`
	RoleId    int
	Role      Role   `gorm:"foreignKey:RoleId"`
	Name      string `gorm:"type:varchar(255)"`
	Email     string `gorm:"uniqueIndex;type:varchar(255)"`
	UserName  string `gorm:"uniqueIndex;type:varchar(255);notnull"`
	Password  string `gorm:"type:varchar(255);notnull"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	hashedPassword, err := hashPassword(u.Password)
	if err != nil {
		return err
	}
	u.Password = hashedPassword
	return
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
