package main

import (
	"LearnGo-todoAuth/initializers"
)

func init() {

	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func main() {

	// initializers.DB.AutoMigrate(&models.Note{})
	// initializers.DB.Migrator().AddColumn(&models.Note{}, "UserID")

	// NOTE: 1. other way to do it you can create funtions which have migration calls and have only one main function
	// add new migration funtion to main and run it, the old table will not be effected and new one will create/update
}
