package routes

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"fmt"

	"db"
	model "models"
)

//CreateMatch :
func CreateMatch(w http.ResponseWriter, r *http.Request) {
	var match model.Match
	_ = json.NewDecoder(r.Body).Decode(&match)

	tx, _ := database.Db.Begin()
	stmt, err := tx.Prepare("INSERT INTO matchesTbl(Team,Opposition,Match_Date,IsCompleted) VALUES(@Team, @Opposition, @MatchDate, @IsCompleted)")
	_, err = stmt.Exec(sql.Named("Team", match.Team), sql.Named("Opposition", match.Opposition), sql.Named("MatchDate", match.MatchDate), sql.Named("IsCompleted", match.IsCompleted))
	if err != nil {
		log.Println("Rollback for create match")
		tx.Rollback()
		json.NewEncoder(w).Encode(err)
	}
	tx.Commit()
	json.NewEncoder(w).Encode(true)
}

//GetAllMatches :
func GetAllMatches(w http.ResponseWriter, r *http.Request) {
	var (
		matches 	[]model.Match
		response	model.Response 
	)
	rows, err := database.Db.Query("SELECT * FROM matchesTbl")
	if err != nil {
		fmt.Errorf("Error in getAllMatches results")
		response = model.Response{Success: false, Message: "Unable to fetch matches list."}
		json.NewEncoder(w).Encode(response)
		defer rows.Close()
		return
	}
	defer rows.Close()
	for rows.Next() {
		var (
			mid int
			team string
			opposition string
			matchDate string
			result *int
			isCompleted int
			status int
		)
		scanerr := rows.Scan(&mid, &team, &opposition, &matchDate, &result, &isCompleted, &status)
		if scanerr != nil {
			fmt.Errorf("Error in scanning answer query")
			response = model.Response{Success: false, Message: scanerr.Error()}
			json.NewEncoder(w).Encode(response)
			return
		}
		resultMatch := model.Match{MID: mid, Team: team, Opposition: opposition, MatchDate: matchDate, Result: result, IsCompleted: isCompleted, Status: status}
		matches = append(matches, resultMatch)
	}
	json.NewEncoder(w).Encode(matches)
}
