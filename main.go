package main

import (
	"log"
	"net/http"
	"os"

	"andalalin/controllers"
	"andalalin/initializers"
	"andalalin/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	server              *gin.Engine
	AuthController      controllers.AuthController
	AuthRouteController routes.AuthRouteController

	UserController      controllers.UserController
	UserRouteController routes.UserRouteController

	AndalalinController      controllers.AndalalinController
	AndalalinRouteController routes.AndalalinRouteController

	SurveyController      controllers.SurveyController
	SurveyRouteController routes.SurveyRouteController

	DataMasterController      controllers.DataMasterControler
	DataMasterRouteController routes.DataMasterRouteController
)

func init() {
	config, err := initializers.LoadConfig()
	if err != nil {
		log.Fatal("Could not load environment variables", err)
	}

	initializers.ConnectDB(&config)

	AuthController = controllers.NewAuthController(initializers.DB)
	AuthRouteController = routes.NewAuthRouteController(AuthController)

	UserController = controllers.NewUserController(initializers.DB)
	UserRouteController = routes.NewRouteUserController(UserController)

	AndalalinController = controllers.NewAndalalinController(initializers.DB)
	AndalalinRouteController = routes.NewRouteAndalalinController(AndalalinController)

	SurveyController = controllers.NewSurveyController(initializers.DB)
	SurveyRouteController = routes.NewSurveyRouteController(SurveyController)

	DataMasterController = controllers.NewDataMasterControler(initializers.DB)
	DataMasterRouteController = routes.NewDataMasterRouteController(DataMasterController)

	os.Setenv("GIN_MODE", "release")
	gin.SetMode(gin.ReleaseMode)
	server = gin.New()
}

func main() {
	config, err := initializers.LoadConfig()
	if err != nil {
		log.Fatal("Could not load environment variables", err)
	}

	corsConfig := cors.DefaultConfig()

	corsConfig.AllowOrigins = []string{"*"}
	corsConfig.AllowCredentials = true

	corsConfig.AllowHeaders = []string{"*"}
	corsConfig.AllowMethods = []string{"*"}

	server.Use(cors.New(corsConfig))

	router := server.Group("/v1")
	router.GET("/healthchecker", func(ctx *gin.Context) {
		message := "Welcome to andalalin"
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": message})
	})

	router.Use(func(c *gin.Context) {
		// Start streaming on startup
		controllers.GetStartUp(initializers.DB)
	})

	AuthRouteController.AuthRoute(router)
	UserRouteController.UserRoute(router)
	AndalalinRouteController.AndalalainRoute(router)
	SurveyRouteController.SurveyRoute(router)
	DataMasterRouteController.DataMasterRoute(router)

	log.Fatal(server.Run(":" + config.ServerPort))
}
