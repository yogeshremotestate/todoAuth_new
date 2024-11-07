package main

import (
	"LearnGo-todoAuth/initializers"
	// "LearnGo-todoAuth/models"
)

func init() {

	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

// func main() {

// 	initializers.DB.AutoMigrate(&models.User{})
// }
