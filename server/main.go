package server

import (
	"log"
	"text/template"

	temp "github.com/Rubncal04/ksmanager/templates"

	"github.com/Rubncal04/ksmanager/config"
	"github.com/Rubncal04/ksmanager/db"
	"github.com/Rubncal04/ksmanager/handlers"
	"github.com/Rubncal04/ksmanager/middleware"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
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

	if err != nil {
		log.Println(err.Error())
	}

	db.RunMigrations(*database)

	cache := db.StartCacheDb(envVariables)

	renderer := &temp.TemplateRenderer{
		Templates: template.Must(template.ParseGlob("templates/*.html")),
	}

	e.Use(echoMidd.Logger())
	e.Use(echoMidd.Recover())
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(envVariables.SESSION_KEY))))

	e.Renderer = renderer
	e.Static("/assets", "assets")

	e.GET("/", handlers.Welcome())

	// Views
	e.GET("/user", handlers.ShowUser())
	e.GET("/church", handlers.ShowChurch())

	// User enpoints
	e.POST("/register", handlers.Register(database))
	e.POST("/login", handlers.Login(database))

	private := e.Group("/api/v1")
	private.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(envVariables.JWT_KEY),
	}))
	private.Use(handlers.NewCache(cache))

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

	// Categories endpoints
	private.GET("/categories", handlers.FindAllCategory(database), middleware.Authorization(RoleRoot, RolePastor, RoleTreasurer))
	private.POST("/categories", handlers.CreateCategory(database), middleware.Authorization(RoleRoot, RolePastor, RoleTreasurer))

	// Offerings endpoints
	private.GET("/worship-services/:worship_id/offerings", handlers.FindAllOffering(database), middleware.Authorization(RoleRoot, RolePastor, RoleTreasurer))
	private.POST("/worship-services/:worship_id/offerings", handlers.CreateOffering(database), middleware.Authorization(RoleRoot, RolePastor, RoleTreasurer))

	e.Logger.Fatal(e.Start(port))
}
