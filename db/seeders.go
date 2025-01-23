package db

import (
	"fmt"

	"github.com/Rubncal04/ksmanager/models"
	"gorm.io/gorm"
)

func Seeders(db PostgresRepo) {
	CountrySeed(db)
	StateSeed(db)
	CitySeed(db)
	ChurchSeed(db)
	RoleSeed(db)
}

func RoleSeed(db PostgresRepo) {
	roles := []models.Role{
		{
			Name: "assistant",
			Permissions: models.Permissions{
				CreateUser:     false,
				UpdateUser:     false,
				DeleteUser:     false,
				GetUser:        true,
				CreateMember:   false,
				UpdateMember:   false,
				DeleteMember:   false,
				GetMember:      true,
				CreateChurch:   false,
				UpdateChurch:   false,
				DeleteChurch:   false,
				GetChurch:      true,
				CreateRole:     false,
				UpdateRole:     false,
				DeleteRole:     false,
				GetRoles:       false,
				CreateOffering: false,
				UpdateOffering: false,
				DeleteOffering: false,
				GetOffering:    false,
			},
		},
		{
			Name: "pastor",
			Permissions: models.Permissions{
				CreateUser:     true,
				UpdateUser:     true,
				DeleteUser:     true,
				GetUser:        true,
				CreateMember:   true,
				UpdateMember:   true,
				DeleteMember:   true,
				GetMember:      true,
				CreateChurch:   true,
				UpdateChurch:   true,
				DeleteChurch:   true,
				GetChurch:      true,
				CreateRole:     true,
				UpdateRole:     true,
				DeleteRole:     true,
				GetRoles:       true,
				CreateOffering: true,
				UpdateOffering: true,
				DeleteOffering: true,
				GetOffering:    true,
			},
		},
		{
			Name: "treasure",
			Permissions: models.Permissions{
				CreateUser:     true,
				UpdateUser:     true,
				DeleteUser:     true,
				GetUser:        true,
				CreateMember:   false,
				UpdateMember:   false,
				DeleteMember:   false,
				GetMember:      true,
				CreateChurch:   false,
				UpdateChurch:   false,
				DeleteChurch:   false,
				GetChurch:      true,
				CreateRole:     true,
				UpdateRole:     false,
				DeleteRole:     false,
				GetRoles:       true,
				CreateOffering: true,
				UpdateOffering: true,
				DeleteOffering: true,
				GetOffering:    true,
			},
		},
		{
			Name: "secretary",
			Permissions: models.Permissions{
				CreateUser:     false,
				UpdateUser:     false,
				DeleteUser:     false,
				GetUser:        true,
				CreateMember:   false,
				UpdateMember:   false,
				DeleteMember:   false,
				GetMember:      true,
				CreateChurch:   false,
				UpdateChurch:   false,
				DeleteChurch:   false,
				GetChurch:      true,
				CreateRole:     false,
				UpdateRole:     false,
				DeleteRole:     false,
				GetRoles:       false,
				CreateOffering: false,
				UpdateOffering: false,
				DeleteOffering: false,
				GetOffering:    false,
			},
		},
	}

	for _, role := range roles {
		var existingRole models.Role
		err := db.db.Where("name = ?", role.Name).First(&existingRole).Error

		if err == nil {
			fmt.Printf("Skipping role '%s' already exists.\n", role.Name)
			continue
		}

		if err != gorm.ErrRecordNotFound {
			fmt.Printf("Role '%s' not found: %v\n", role.Name, err)
			continue
		}

		if err := db.db.Create(&role).Error; err != nil {
			fmt.Printf("Role '%s' couldn't be created: %v\n", role.Name, err)
		} else {
			fmt.Printf("Role '%s' was created successfully.\n", role.Name)
		}
	}
}

func ChurchSeed(db PostgresRepo) {
	var country models.Country
	var state models.State
	var city models.City
	err := db.db.Where("name = ?", "Colombia").First(&country).Error
	if err != nil {
		fmt.Println("Error finding country, state or city.")
	}
	err = db.db.Where("name = ?", "Atlantico").First(&state).Error
	if err != nil {
		fmt.Println("Error finding country, state or city.")
	}
	err = db.db.Where("name = ?", "Barranquilla").First(&city).Error

	if err != nil {
		fmt.Println("Error finding country, state or city.")
	}

	churches := []models.Church{
		{
			Name:      "IPUC El Bosque central",
			Address:   "Cra 6A # 73C - 257",
			CountryId: int(country.ID),
			StateId:   int(state.ID),
			CityId:    int(city.ID),
		},
		{
			Name:      "IPUC Las Americas",
			Address:   "Calle 54 C N 4-11",
			CountryId: int(country.ID),
			StateId:   int(state.ID),
			CityId:    int(city.ID),
		},
		{
			Name:      "IPUC Central de Barranquilla",
			Address:   "Cl. 57 #43-131",
			CountryId: int(country.ID),
			StateId:   int(state.ID),
			CityId:    int(city.ID),
		},
	}

	for _, church := range churches {
		var existingChurch models.Church
		err := db.db.Where("name = ?", church.Name).First(&existingChurch).Error

		if err == nil {
			fmt.Printf("Skipping church '%s' already exists.\n", church.Name)
			continue
		}

		if err != gorm.ErrRecordNotFound {
			fmt.Printf("Church '%s' not found: %v\n", church.Name, err)
			continue
		}

		if err := db.db.Create(&church).Error; err != nil {
			fmt.Printf("Church '%s' couldn't be created: %v\n", church.Name, err)
		} else {
			fmt.Printf("Church '%s' was created successfully.\n", church.Name)
		}
	}
}

func CountrySeed(db PostgresRepo) {
	country := models.Country{
		Name: "Colombia",
	}

	var existingCountry models.Country
	err := db.db.Where("name = ?", "Colombia").First(&existingCountry).Error
	if err == nil {
		fmt.Printf("Skipping country '%s' already exists.\n", country.Name)
		return
	}

	if err != gorm.ErrRecordNotFound {
		fmt.Printf("Country '%s' not found: %v\n", country.Name, err)
	}

	if err := db.db.Create(&country).Error; err != nil {
		fmt.Printf("Country '%s' couldn't be created: %v\n", country.Name, err)
	} else {
		fmt.Printf("Country '%s' was created successfully.\n", country.Name)
	}
}

func StateSeed(db PostgresRepo) {
	var country models.Country
	err := db.db.Where("name = ?", "Colombia").First(&country).Error

	if err == nil {
		fmt.Printf("Error finding country '%s'.\n", country.Name)
	}

	state := models.State{
		Name:      "Atlantico",
		CountryId: int(country.ID),
	}

	err = db.db.Where("name = ?", "Atlantico").First(&state).Error
	if err == nil {
		fmt.Printf("Skipping state '%s' already exists.\n", state.Name)
		return
	}

	if err != gorm.ErrRecordNotFound {
		fmt.Printf("State '%s' not found: %v\n", state.Name, err)
	}

	if err := db.db.Create(&state).Error; err != nil {
		fmt.Printf("State '%s' couldn't be created: %v\n", state.Name, err)
	} else {
		fmt.Printf("State '%s' was created successfully.\n", state.Name)
	}
}

func CitySeed(db PostgresRepo) {
	var state models.State
	err := db.db.Where("name = ?", "Atlantico").First(&state).Error

	if err == nil {
		fmt.Printf("Error finding state '%s'.\n", state.Name)
	}

	city := models.City{
		Name:    "Barranquilla",
		StateId: int(state.ID),
	}

	err = db.db.Where("name = ?", "Barranquilla").First(&city).Error
	if err == nil {
		fmt.Printf("Skipping city '%s' already exists.\n", city.Name)
		return
	}

	if err != gorm.ErrRecordNotFound {
		fmt.Printf("City '%s' not found: %v\n", city.Name, err)
	}

	if err := db.db.Create(&city).Error; err != nil {
		fmt.Printf("City '%s' couldn't be created: %v\n", city.Name, err)
	} else {
		fmt.Printf("City '%s' was created successfully.\n", city.Name)
	}
}
