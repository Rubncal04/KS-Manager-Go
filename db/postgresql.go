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

func NewPostgresRepo(variables *config.EnvVariables) (*PostgresRepo, error) {
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

func (db *PostgresRepo) FindAllMembers(id string) ([]models.Member, error) {
	var members []models.Member

	err := db.db.Where("church_id = ?", id).Find(&members).Error

	return members, err
}

func (db *PostgresRepo) FindOneMember(id string) (*models.Member, error) {
	var member models.Member
	result := db.db.First(&member, id)

	if result.Error != nil {
		log.Println(result.Error)
		return nil, result.Error
	}

	return &member, nil
}

func (db *PostgresRepo) CreateMember(member *models.Member) (*models.Member, error) {
	result := db.db.Create(&member)

	if result.Error != nil {
		log.Println(result.Error)
		return nil, result.Error
	}

	return member, nil
}

func (db *PostgresRepo) UpdateMember(id string, fields *models.Member) (*models.Member, error) {
	var update models.Member

	if err := db.db.First(&update, id).Error; err != nil {
		return nil, err
	}

	if fields.Name != "" {
		update.Name = fields.Name
	}

	if fields.LastName != "" {
		update.LastName = fields.LastName
	}

	if fields.Email != "" {
		update.Email = fields.Email
	}

	if fields.IdentificationNumber != "" {
		update.IdentificationNumber = fields.IdentificationNumber
	}

	if fields.Address != "" {
		update.Address = fields.Address
	}

	if fields.Birthday != "" {
		update.Birthday = fields.Birthday
	}

	if fields.BaptizedBy != "" {
		update.BaptizedBy = fields.BaptizedBy
	}

	if fields.BaptizedOn != "" {
		update.BaptizedOn = fields.BaptizedOn
	}

	if fields.HolySpiritOn != "" {
		update.HolySpiritOn = fields.HolySpiritOn
	}

	if fields.Position != "" {
		update.Position = fields.Position
	}

	if fields.NumChildren != 0 {
		update.NumChildren = fields.NumChildren
	}

	if len(fields.ChildrenNames) > 0 {
		update.ChildrenNames = fields.ChildrenNames
	}

	if fields.PartnerName != "" {
		update.PartnerName = fields.PartnerName
	}

	if fields.Degree != "" {
		update.Degree = fields.Degree
	}

	if fields.Profession != "" {
		update.Profession = fields.Profession
	}

	if err := db.db.Save(&update).Error; err != nil {
		return nil, err
	}

	return &update, nil
}

func (db *PostgresRepo) DeleteMember(id string) (string, error) {
	var member *models.Member
	if err := db.db.Delete(&member, id).Error; err != nil {
		return "", err
	}

	return id, nil
}

func (db *PostgresRepo) CreateUser(user *models.User) (*models.User, error) {
	newUser := db.db.Create(&user)

	if newUser.Error != nil {
		log.Println(newUser.Error)
		return nil, newUser.Error
	}

	return user, nil
}

func (db *PostgresRepo) FindUser(user *models.User) (*models.User, error) {
	result := db.db.First(&user)

	if result.Error != nil {
		log.Println(result.Error)
		return nil, result.Error
	}

	return user, nil
}

func (db *PostgresRepo) FindUserBy(user *models.User, field, property string) (*models.User, error) {
	result := db.db.Where(property+" = ?", field).First(&user)

	if result.Error != nil {
		log.Println(result.Error)
		return nil, result.Error
	}

	return user, nil
}

func (db *PostgresRepo) CreateRole(role *models.Role) (*models.Role, error) {
	newRole := db.db.Create(&role)

	if newRole.Error != nil {
		log.Println(newRole.Error)
		return nil, newRole.Error
	}

	return role, nil
}

func (db *PostgresRepo) FindAllRoles() ([]models.Role, error) {
	var roles []models.Role

	result := db.db.Find(&roles)

	if result.Error != nil {
		log.Println(result.Error)
		return nil, result.Error
	}

	return roles, nil
}

func (db *PostgresRepo) FindOneRoleById(id string) (*models.Role, error) {
	var role models.Role

	result := db.db.First(&role, id)

	if result.Error != nil {
		log.Println(result.Error)
		return nil, result.Error
	}

	return &role, nil
}

func (db *PostgresRepo) FindRoleByName(name string) (*models.Role, error) {
	var role models.Role

	err := db.db.Where("name = ?", name).First(&role)
	if err.Error != nil {
		log.Println(err.Error)
		return nil, err.Error
	}

	return &role, nil
}

func (db *PostgresRepo) DeleteRole(id string) (string, error) {
	var role *models.Role

	if err := db.db.Delete(&role, id).Error; err != nil {
		return "", err
	}

	return id, nil
}

func (db *PostgresRepo) CreateWorshipService(wService *models.WorshipService) (*models.WorshipService, error) {
	result := db.db.Create(&wService)

	if result.Error != nil {
		log.Println(result.Error)
		return nil, result.Error
	}

	return wService, nil
}

func (db *PostgresRepo) FindAllWorship(id string) ([]models.WorshipService, error) {
	var worship []models.WorshipService

	err := db.db.Where("church_id = ?", id).Find(&worship).Error

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return worship, err
}

func (db *PostgresRepo) UpdateWorship(id string, fields *models.WorshipService) (*models.WorshipService, error) {
	var update models.WorshipService

	if err := db.db.First(&update, id).Error; err != nil {
		log.Println(err)
		return nil, err
	}

	if fields.Name != "" {
		update.Name = fields.Name
	}

	if fields.Day != "" {
		update.Day = fields.Day
	}

	if err := db.db.Save(&update).Error; err != nil {
		log.Println(err)
		return nil, err
	}

	return &update, nil
}

func (db *PostgresRepo) DeleteWorship(id string) (string, error) {
	var worship models.WorshipService

	if err := db.db.Delete(&worship, id).Error; err != nil {
		return "", err
	}

	return id, nil
}

func (db *PostgresRepo) FindWorshipByID(id string) (*models.WorshipService, error) {
	var worship models.WorshipService

	if err := db.db.First(&worship, id).Error; err != nil {
		log.Println(err)
		return nil, err
	}

	return &worship, nil
}
