package db

import (
	"errors"
	"log"

	"github.com/Rubncal04/ksmanager/models"
	"gorm.io/gorm"
)

func RunMigrations(db PostgresRepo) {
	db.db.AutoMigrate(models.Country{})
	db.db.AutoMigrate(models.State{})
	db.db.AutoMigrate(models.City{})
	db.db.AutoMigrate(models.Church{})
	db.db.AutoMigrate(models.Member{})
	db.db.AutoMigrate(models.Role{})
	db.db.AutoMigrate(models.User{})
	role := createRootRole(db)
	addRoleIDToUsers(db, *role)
}

func addRoleIDToUsers(db PostgresRepo, role models.Role) error {
	db.db.Model(&models.User{}).Where("role_id IS NULL").Update("role_id", role.ID)

	return db.db.Exec("ALTER TABLE users ALTER COLUMN role_id SET NOT NULL").Error
}

func createRootRole(db PostgresRepo) *models.Role {
	role := models.Role{}
	existsRole := db.db.Where("name = ?", "root").First(&role)
	if existsRole.Error != nil {
		log.Println(existsRole.Error)
		return nil
	}

	if existsRole != nil && errors.Is(existsRole.Error, gorm.ErrRecordNotFound) {
		user := models.Role{
			Name: "root",
			Permissions: models.Permissions{
				CreateUser:   true,
				UpdateUser:   true,
				DeleteUser:   true,
				GetUser:      true,
				CreateMember: true,
				UpdateMember: true,
				DeleteMember: true,
				GetMember:    true,
				CreateChurch: true,
				UpdateChurch: true,
				DeleteChurch: true,
				GetChurch:    true,
				CreateRole:   true,
				UpdateRole:   true,
				DeleteRole:   true,
				GetRoles:     true,
			},
		}

		db.CreateRole(&user)
		return &user
	}

	return &role
}
