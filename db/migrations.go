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
	Seeders(db)
	db.db.AutoMigrate(models.WorshipService{})
	db.db.AutoMigrate(models.Category{})
	db.db.AutoMigrate(models.Offering{})
}

func addRoleIDToUsers(db PostgresRepo, role models.Role) error {
	if role.ID != 0 {
		db.db.Model(&models.User{}).Where("role_id IS NULL").Update("role_id", role.ID)

		return db.db.Exec("ALTER TABLE users ALTER COLUMN role_id SET NOT NULL").Error
	}

	return nil
}

func createRootRole(db PostgresRepo) *models.Role {
	role := models.Role{}
	existsRole := db.db.Where("name = ?", "root").First(&role)

	if errors.Is(existsRole.Error, gorm.ErrRecordNotFound) {
		role = models.Role{
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

		db.CreateRole(&role)
		if _, err := db.CreateRole(&role); err != nil {
			log.Fatalf("Failed to create root role: %v", err)
			return nil
		}
	}

	return &role
}
