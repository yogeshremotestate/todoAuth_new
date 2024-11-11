package initializers

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

func ConnectToDB() {
	dsn := ENV.DB_URL
	var err error
	DB, err = sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}
	fmt.Println("Database connected successfully")
}
