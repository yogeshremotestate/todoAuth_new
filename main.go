package main

import (
	"LearnGo-todoAuth/controllers"
	_ "LearnGo-todoAuth/docs"
	"LearnGo-todoAuth/initializers"
	Log "LearnGo-todoAuth/log"
	"LearnGo-todoAuth/middleware"
	"log"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title TODO APIs
// @version 1.0
// @description Testing Swagger APIs.
// @termsOfService http://google.com/
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @securityDefinitions.apiKey BearerAuth
// @in header
// @name Authorization
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8081
// @BasePath /api/v1
// @schemes http
func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func main() {

	if err := Log.InitializeLogger(); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	// defer middleware.Logger.Sync()
	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.Use(Log.LoggerMiddleware())
	noteRoutes := r.Group("/notes", middleware.AuthValidate)
	{
		noteRoutes.POST("/", controllers.NoteCreate)
		noteRoutes.GET("/", controllers.GetAllNote)
		noteRoutes.POST("/upload", controllers.UploadExcel)

		userSpecific := noteRoutes.Group("/:id", middleware.VerifyUserNote)
		{
			userSpecific.GET("/", controllers.GetNote)
			userSpecific.PUT("/", controllers.UpdateNote)
			userSpecific.DELETE("/", controllers.DeleteNote)
		}

	}

	userRoutes := r.Group("/user")
	{
		userRoutes.POST("/signup", controllers.SignUpUser)
		userRoutes.POST("/login", controllers.LoginUser)
	}
	r.Run()
}
