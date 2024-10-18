package db

import "github.com/Rubncal04/ksmanager/models"

func RunMigrations(db PostgresRepo) {
	db.db.AutoMigrate(models.Church{})
}
