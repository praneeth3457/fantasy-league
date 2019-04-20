package main

import (
	"log"
	"net/http"
	"os"
	"constant"
	"db"
	"routes"
	"authorization"
	"fmt"


	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	//Database connection
	database.DbConnect()
	defer database.Db.Close()
	port := ":" + os.Getenv("PORT")
	fmt.Println(port)
	//port := ":8000"
	//Init Router
	r := mux.NewRouter()

	//User routes
	r.HandleFunc("/api/user/create", routes.CreateUser).Methods("POST")
	r.HandleFunc("/api/user/verify", routes.VerifyUser).Methods("POST")
	r.HandleFunc("/api/user/getAnswers/{id}", routes.GetAnswers).Methods("GET")
	r.HandleFunc("/api/user/getQuestions", routes.GetQuestions).Methods("GET")
	//Match routes
	r.Handle("/api/match/create", authorization.IsAuthorized(routes.CreateMatch, constant.UserAdmin)).Methods("POST")

	// Solves Cross Origin Access Issue
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:4200","https://falcons-fantasy.herokuapp.com"},
	})
	handler := c.Handler(r)

	srv := &http.Server{
		Handler: handler,
		Addr:    ":" + os.Getenv("PORT"),
	}
	
	log.Fatal(srv.ListenAndServe())
	//log.Fatal(http.ListenAndServe(port, r))
}
