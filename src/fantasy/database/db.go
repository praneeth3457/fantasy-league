package database

import (
	"database/sql"
	"fmt"
	//"os"
	"log"
)

var (
	server   = "fantasyfalcons.cylbbo3xafqn.us-east-2.rds.amazonaws.com"
	port     = 1433
	user     = "fantasy123"
	password = "falcons123"
	database = "fantasyLeague"
	sslca = "rds-combined-ca-bundle.pem"
	// Db : Global Db vairable
	Db *sql.DB
)

// DbConnect :
func DbConnect() {
	var err error

	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;sslca=%s;",
	 	server, user, password, port, database, sslca)
	//connString := "fantasy123:falcons123@fantasyfalcons.cylbbo3xafqn.us-east-2.rds.amazonaws.com:1433/fantasyleague?ssl=true&sslrootcert=rds-combined-ca-bundle.pem&sslmode=require"
	Db, err = sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatal("Error creating connection pool: " + err.Error())
	}
	log.Printf("Connecteddsdad!\n")
}
