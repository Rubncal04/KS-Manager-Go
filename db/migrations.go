package db

import "github.com/Rubncal04/ksmanager/models"

func RunMigrations(db PostgresRepo) {
	db.db.AutoMigrate(models.Country{})
	db.db.AutoMigrate(models.State{})
	db.db.AutoMigrate(models.City{})
	db.db.AutoMigrate(models.Church{})
	db.db.AutoMigrate(models.Member{})
}
