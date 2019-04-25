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

//GetAllMatches :
func OtherMatchDetails(w http.ResponseWriter, r *http.Request) {
	var (
		other model.Other
		otherNew model.Other
		otherDetail model.OtherDetail
		response	model.Response 
	)
	_ = json.NewDecoder(r.Body).Decode(&other)

	row := database.Db.QueryRow("SELECT * FROM othersTbl WHERE UID=@UID AND MID=@MID", sql.Named("UID", other.UID), sql.Named("MID", other.MID))
	err := row.Scan(&otherNew.OID, &otherNew.Captain, &otherNew.MVBA, &otherNew.MVBO, &otherNew.MVAR, &otherNew.UID, &otherNew.MID)
	if err != nil {
		response = model.Response{Success: false, Message: "no rows in result set"}
		json.NewEncoder(w).Encode(response)
		return
	}

	players, playersErr := database.Db.Query("SELECT * FROM playersTbl WHERE PID IN (@Captain, @MVBA, @MVBO, @MVAR)", 
						sql.Named("Captain", otherNew.Captain), sql.Named("MVBA", otherNew.MVBA), sql.Named("MVBO", otherNew.MVBO), sql.Named("MVAR", otherNew.MVAR))
	if playersErr != nil {
		response = model.Response{Success: false, Message: playersErr.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	defer players.Close()
	for players.Next() {
		var (
			pid		int
			name	string
			role	string
			team	string
		)
		scanerr := players.Scan(&pid, &name, &role, &team)
		if scanerr != nil {
			json.NewEncoder(w).Encode(scanerr.Error())
			fmt.Errorf("Error in scanning answer query")
		}
		if pid == otherNew.Captain {
			otherDetail.Captain = model.Player{PID: pid, Name: name, Role: role, Team: team}
		} else if pid == otherNew.MVBA {
			otherDetail.MVBA = model.Player{PID: pid, Name: name, Role: role, Team: team}
		} else if pid == otherNew.MVBO {
			otherDetail.MVBO = model.Player{PID: pid, Name: name, Role: role, Team: team}
		} else if pid == otherNew.MVAR {
			otherDetail.MVAR = model.Player{PID: pid, Name: name, Role: role, Team: team}
		}
	}
	otherDetail.OID = otherNew.OID
	otherDetail.UID = otherNew.UID
	otherDetail.MID = otherNew.MID
	
	json.NewEncoder(w).Encode(model.Response3{Success: true, Message: otherDetail})

}