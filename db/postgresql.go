package db

import (
	"fmt"
	"log"
	"os"

	"github.com/Rubncal04/ksmanager/config"
	"github.com/Rubncal04/ksmanager/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresRepo struct {
	db *gorm.DB
}

func NewPostgresRepo() (*PostgresRepo, error) {
	var variables = config.GetVariables()
	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=%v TimeZone=%v",
		variables.DB_HOST, variables.DB_USER, variables.DB_PASSWORD, variables.DB_NAME, variables.DB_PORT,
		variables.SSL_MODE, variables.TIME_ZONE,
	)

	dbpool, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
		return nil, err
	}

	log.Println("Starting database")

	return &PostgresRepo{
		db: dbpool,
	}, nil
}

func (db *PostgresRepo) FindAllChurches() ([]models.Church, error) {
	var churches []models.Church
	result := db.db.Find(&churches)

	if result.Error != nil {
		log.Println(result.Error)
		return nil, result.Error
	}

	return churches, nil
}

func (db *PostgresRepo) FindOneChurch(id string) (*models.Church, error) {
	var church models.Church
	result := db.db.First(&church, id)

	if result.Error != nil {
		log.Println(result.Error)
		return nil, result.Error
	}

	return &church, nil
}

func (db *PostgresRepo) CreateChurch(church *models.Church) (*models.Church, error) {
	result := db.db.Create(&church)

	if result.Error != nil {
		log.Println(result.Error)
		return nil, result.Error
	}

	return church, nil
}

func (db *PostgresRepo) UpdateChurch(id string, fields *models.Church) (*models.Church, error) {
	var update models.Church
	db.db.First(&update, id)

	if fields.Name != "" {
		update.Name = fields.Name
	}

	if fields.Address != "" {
		update.Address = fields.Address
	}

	if fields.CityId != 0 {
		update.CityId = fields.CityId
	}

	if fields.CountryId != 0 {
		update.CountryId = fields.CountryId
	}

	if fields.StateId != 0 {
		update.StateId = fields.StateId
	}

	db.db.Save(&update)

	return &update, nil
}
