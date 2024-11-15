package main

import (
	"LearnGo-todoAuth/controllers"
	"LearnGo-todoAuth/initializers"
	"LearnGo-todoAuth/middleware"
	"log"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func main() {

	if err := middleware.InitializeLogger(); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer middleware.Logger.Sync()
	r := gin.Default()
	r.Use(middleware.LoggerMiddleware())
	noteRoutes := r.Group("/notes", middleware.AuthValidate)
	{
		noteRoutes.POST("/", controllers.NoteCreate)
		noteRoutes.GET("/", controllers.GetAllNote)
		noteRoutes.GET("/:id", controllers.GetNote)
		noteRoutes.PUT("/:id", controllers.UpdateNote)
		noteRoutes.DELETE("/:id", controllers.DeleteNote)
		noteRoutes.POST("/upload", controllers.UploadExcel)
	}

	userRoutes := r.Group("/user")
	{
		userRoutes.POST("/signup", controllers.SignUpUser)
		userRoutes.POST("/login", controllers.LoginUser)
	}
	r.Run()
}
