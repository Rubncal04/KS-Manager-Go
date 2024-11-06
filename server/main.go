package server

import (
	"log"
	"net/http"

	"github.com/Rubncal04/ksmanager/config"
	"github.com/Rubncal04/ksmanager/db"
	"github.com/Rubncal04/ksmanager/handlers"
	"github.com/Rubncal04/ksmanager/middleware"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	echoMidd "github.com/labstack/echo/v4/middleware"
)

const (
	RoleRoot      = "root"
	RolePastor    = "pastor"
	RoleSecretary = "secretary"
	RoleTreasurer = "treasurer"
	RoleAssistant = "assistant"
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
	private.GET("/churches", handlers.Index(database), middleware.Authorization(RoleRoot, RolePastor, RoleSecretary, RoleTreasurer))
	private.GET("/churches/:id", handlers.FindOne(database), middleware.Authorization(RoleRoot, RolePastor, RoleSecretary, RoleTreasurer))
	private.PUT("/churches/:id", handlers.Update(database), middleware.Authorization(RoleRoot, RolePastor))
	private.POST("/churches", handlers.Create(database), middleware.Authorization(RoleRoot, RolePastor))

	// Member endpoints
	private.GET("/members/:id", handlers.FindOneMember(database), middleware.Authorization(RoleRoot, RolePastor, RoleSecretary, RoleAssistant))
	private.PUT("/members/:id", handlers.UpdateMember(database), middleware.Authorization(RoleRoot, RolePastor, RoleSecretary))
	private.DELETE("/members/:id", handlers.DeleteMember(database), middleware.Authorization(RoleRoot, RolePastor, RoleSecretary))
	private.GET("/members/churches/:church_id", handlers.IndexMembers(database), middleware.Authorization(RoleRoot, RolePastor, RoleSecretary, RoleAssistant))
	private.POST("/members/churches/:church_id", handlers.CreateMember(database), middleware.Authorization(RoleRoot, RolePastor, RoleSecretary))

	// Roles endpoints
	private.GET("/roles", handlers.GetAllRoles(database), middleware.Authorization(RoleRoot, RolePastor))
	private.POST("/roles", handlers.CreateRole(database), middleware.Authorization(RoleRoot))
	private.DELETE("/roles/:id", handlers.DeleteRole(database), middleware.Authorization(RoleRoot))

	// Worship services endpoints
	private.GET("/worship-services/:id", handlers.FindOneWorshipService(database), middleware.Authorization(RoleRoot, RolePastor, RoleSecretary, RoleTreasurer))
	private.GET("/worship-services", handlers.IndexWorshipService(database), middleware.Authorization(RoleRoot, RolePastor, RoleSecretary, RoleTreasurer))
	private.POST("/worship-services", handlers.CreateWorshipService(database), middleware.Authorization(RoleRoot, RolePastor, RoleSecretary))
	private.PUT("/worship-services/:id", handlers.UpdateWorshipService(database), middleware.Authorization(RoleRoot, RolePastor, RoleSecretary))
	private.DELETE("/worship-services/:id", handlers.DeleteWorshipService(database), middleware.Authorization(RoleRoot, RolePastor, RoleSecretary))

	e.Logger.Fatal(e.Start(port))
}
