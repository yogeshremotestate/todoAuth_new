package main

import (
	"LearnGo-todoAuth/controllers"
	"LearnGo-todoAuth/initializers"
	"LearnGo-todoAuth/middleware"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func main() {

	r := gin.Default()
	r.POST("/notes/", middleware.Validate, controllers.NoteCreate)
	r.GET("/notes/", middleware.Validate, controllers.GetAllNote)
	r.GET("/notes/:id", middleware.Validate, controllers.GetNote)
	r.PUT("/notes/:id", middleware.Validate, controllers.UpdateNote)
	r.DELETE("/notes/:id", middleware.Validate, controllers.DeleteNote)

	r.POST("/user/signup/", controllers.SignUpUser)
	r.POST("/user/login/", controllers.LoginUser)
	r.Run()
}
