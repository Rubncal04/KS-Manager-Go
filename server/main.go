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
	e.PUT("/churches/:id", handlers.Update(database))
	e.POST("/churches", handlers.Create(database))

	e.GET("/members/:id", handlers.FindOneMember(database))
	e.PUT("/members/:id", handlers.UpdateMember(database))
	e.DELETE("/members/:id", handlers.DeleteMember(database))
	e.GET("/members/churches/:church_id", handlers.IndexMembers(database))
	e.POST("/members/churches/:church_id", handlers.CreateMember(database))

	e.Logger.Fatal(e.Start(":1323"))
}
