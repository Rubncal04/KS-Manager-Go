package server

import (
	"log"
	"net/http"

	"github.com/Rubncal04/ksmanager/db"
	"github.com/Rubncal04/ksmanager/handlers"
	"github.com/labstack/echo/v4"
)

func StartServer() {
	e := echo.New()

	database, err := db.NewPostgresRepo()
	db.RunMigrations(*database)
	if err != nil {
		log.Println(err.Error())
	}

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.GET("/churches", handlers.Index(database))
	e.GET("/churches/:id", handlers.FindOne(database))
	e.POST("/churches", handlers.Create(database))

	e.Logger.Fatal(e.Start(":1323"))
}
