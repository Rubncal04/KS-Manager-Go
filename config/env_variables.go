package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type EnvVariables struct {
	DatabaseUrl string
	PORT        string
	DB_HOST     string
	DB_USER     string
	DB_PASSWORD string
	DB_PORT     string
	SSL_MODE    string
	TIME_ZONE   string
	DB_NAME     string
	JWT_KEY     string
	SESSION_KEY string
}

func GetVariables() *EnvVariables {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return &EnvVariables{
		DatabaseUrl: os.Getenv("DATABASE_URL"),
		PORT:        os.Getenv("PORT"),
		DB_HOST:     os.Getenv("DB_HOST"),
		DB_USER:     os.Getenv("DB_USER"),
		DB_PASSWORD: os.Getenv("DB_PASSWORD"),
		DB_PORT:     os.Getenv("DB_PORT"),
		SSL_MODE:    os.Getenv("SSL_MODE"),
		TIME_ZONE:   os.Getenv("TIME_ZONE"),
		DB_NAME:     os.Getenv("DB_NAME"),
		JWT_KEY:     os.Getenv("JWT_KEY"),
		SESSION_KEY: os.Getenv("SESSION_KEY"),
	}
}
