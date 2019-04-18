package routes

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"db"
	model "models"
)

//CreateMatch :
func CreateMatch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var match model.Match
	_ = json.NewDecoder(r.Body).Decode(&match)

	tx, _ := database.Db.Begin()
	stmt, err := tx.Prepare("INSERT INTO matchesTbl(Team,Opposition,Match_Date,IsCompleted) VALUES(@Team, @Opposition, @MatchDate, @IsCompleted)")
	_, err = stmt.Exec(sql.Named("Team", match.Team), sql.Named("Opposition", match.Opposition), sql.Named("MatchDate", match.MatchDate), sql.Named("IsCompleted", match.IsCompleted))
	if err != nil {
		log.Println("Rollback for create match")
		tx.Rollback()
		json.NewEncoder(w).Encode(err)
		log.Fatal(err)
	}
	tx.Commit()
	json.NewEncoder(w).Encode(true)
}
