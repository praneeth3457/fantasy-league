package database

import (
	"database/sql"
	//"fmt"
	"log"
)

var (
	server   = "fantasyfalcons.cylbbo3xafqn.us-east-2.rds.amazonaws.com"
	port     = 1433
	user     = "fantasy123"
	password = "falcons123"
	database = "fantasyLeague"
	// Db : Global Db vairable
	Db *sql.DB
)

// DbConnect :
func DbConnect() {
	var err error

	//connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
	//	server, user, password, port, database)
	Db, err = sql.Open("sqlserver", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Error creating connection pool: " + err.Error())
	}
	log.Printf("Connecteddsdad!\n")
}
