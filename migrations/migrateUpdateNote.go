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

// 	initializers.DB.Migrator().AddColumn(&models.Note{}, "UserID")
// }
