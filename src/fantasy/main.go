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

	// Solves Cross Origin Access Issue
	c := cors.New(cors.Options{
		AllowedHeaders: []string{"Content-Type","Token","User-Context"},
		AllowedOrigins: []string{"http://localhost:4200","https://falcons-fantasy.herokuapp.com","http://falcons-fantasy.herokuapp.com"},
	})
	handler := c.Handler(r)

	//User routes
	r.HandleFunc("/api/user/create", routes.CreateUser).Methods("POST")
	r.HandleFunc("/api/user/verify", routes.VerifyUser).Methods("POST")
	r.HandleFunc("/api/user/getAnswers/{id}", routes.GetAnswers).Methods("GET")
	r.HandleFunc("/api/user/getQuestions", routes.GetQuestions).Methods("GET")
	//Match routes
	r.Handle("/api/match/create", authorization.IsAuthorized(routes.CreateMatch, constant.UserAdmin)).Methods("POST")
	r.Handle("/api/match/getAllMatches", authorization.IsAuthorized(routes.GetAllMatches, constant.UserAny)).Methods("GET")
	r.Handle("/api/match/otherMatchDetails", authorization.IsAuthorized(routes.OtherMatchDetails, constant.UserAny)).Methods("POST")
	r.Handle("/api/match/saveOtherMatchDetails", authorization.IsAuthorized(routes.SaveOtherMatchDetails, constant.UserAny)).Methods("POST")
	r.Handle("/api/match/getAllMatchPoints", authorization.IsAuthorized(routes.GetAllMatchPoints, constant.UserAny)).Methods("POST")
	r.Handle("/api/match/getAllUserMatchPoints", authorization.IsAuthorized(routes.GetAllUserMatchPoints, constant.UserAny)).Methods("GET")
	//usersPlayers routes
	r.Handle("/api/usersPlayers/createAvailability", authorization.IsAuthorized(routes.CreateAvailability, constant.UserAuthor)).Methods("GET")
	r.Handle("/api/usersPlayers/getAvailability", authorization.IsAuthorized(routes.GetAvailability, constant.UserAny)).Methods("POST")
	//Players routes
	r.Handle("/api/player/getAllPlayers", authorization.IsAuthorized(routes.GetAllPlayers, constant.UserAny)).Methods("GET")
	r.Handle("/api/player/savePlayersScores", authorization.IsAuthorized(routes.SavePlayersScore, constant.UserAdmin)).Methods("POST")

	srv := &http.Server{
		Handler: handler,
		Addr:    ":" + os.Getenv("PORT"),
	}
	
	log.Fatal(srv.ListenAndServe())
	//log.Fatal(http.ListenAndServe(port, r))
}
