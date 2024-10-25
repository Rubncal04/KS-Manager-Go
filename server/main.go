package server

import (
	"log"
	"net/http"

	"github.com/Rubncal04/ksmanager/config"
	"github.com/Rubncal04/ksmanager/db"
	"github.com/Rubncal04/ksmanager/handlers"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	echoMidd "github.com/labstack/echo/v4/middleware"
)

func StartServer() {
	e := echo.New()
	envVariables := config.GetVariables()
	port := envVariables.PORT

	database, err := db.NewPostgresRepo(envVariables)
	db.RunMigrations(*database)
	if err != nil {
		log.Println(err.Error())
	}

	e.Use(echoMidd.Logger())
	e.Use(echoMidd.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	// User enpoints
	e.POST("/register", handlers.Register(database))
	e.POST("/login", handlers.Login(database))

	private := e.Group("/api/v1")
	private.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(envVariables.JWT_KEY),
	}))

	// Church endpoints
	private.GET("/churches", handlers.Index(database))
	private.GET("/churches/:id", handlers.FindOne(database))
	private.PUT("/churches/:id", handlers.Update(database))
	private.POST("/churches", handlers.Create(database))

	// Member endpoints
	private.GET("/members/:id", handlers.FindOneMember(database))
	private.PUT("/members/:id", handlers.UpdateMember(database))
	private.DELETE("/members/:id", handlers.DeleteMember(database))
	private.GET("/members/churches/:church_id", handlers.IndexMembers(database))
	private.POST("/members/churches/:church_id", handlers.CreateMember(database))

	e.Logger.Fatal(e.Start(port))
}
