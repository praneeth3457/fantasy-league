package routes

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"fmt"
	"strconv"

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
			player model.Player
		)
		scanerr := players.Scan(&pid, &name, &role, &team)
		if scanerr != nil {
			json.NewEncoder(w).Encode(scanerr.Error())
			fmt.Errorf("Error in scanning answer query")
		}
		if pid == otherNew.Captain {
			player = model.Player{PID: pid, Name: name, Role: role, Team: team, Type: "CAPTAIN"}
		} else if pid == otherNew.MVBA {
			player = model.Player{PID: pid, Name: name, Role: role, Team: team, Type: "MV BATSMAN"}
		} else if pid == otherNew.MVBO {
			player = model.Player{PID: pid, Name: name, Role: role, Team: team, Type: "MV BOWLER"}
		} else if pid == otherNew.MVAR {
			player = model.Player{PID: pid, Name: name, Role: role, Team: team, Type: "MV FIELDER"}
		}

		otherDetail.Players = append(otherDetail.Players, player)
	}
	otherDetail.OID = otherNew.OID
	otherDetail.UID = otherNew.UID
	otherDetail.MID = otherNew.MID
	
	json.NewEncoder(w).Encode(model.Response3{Success: true, Message: otherDetail})

}

/*
	Checking if all players are unique
	Checking if match is in active mode(edit mode)
	Checking if other row exist
	If exist modify the existing other row
	Also update the players availability by reducing the available field & by adding the removed player
	Else create new other row
	Also, update the players availability by reducing the available field
*/
func SaveOtherMatchDetails(w http.ResponseWriter, r *http.Request) {
	var (
		other model.Other
		otherDb model.Other
		response model.Response
		status int
	)
	_ = json.NewDecoder(r.Body).Decode(&other)

	//Checking if all players are unique
	if other.Captain == other.MVBA || other.Captain == other.MVBO || other.Captain == other.MVAR || other.MVBA == other.MVBO || other.MVBA == other.MVAR || other.MVBO == other.MVAR{
		response = model.Response{Success: false, Message: "All players should be unique"}
		json.NewEncoder(w).Encode(response)
		return
	}

	//Checking if match is in active mode(edit mode)
	err := database.Db.QueryRow("SELECT Status FROM matchesTbl WHERE MID=@MID", sql.Named("MID", other.MID)).Scan(&status)
	if err != nil {
		response = model.Response{Success: false, Message: "Invalid match."}
		json.NewEncoder(w).Encode(response)
		return
	}
	if status == 2 {
		response = model.Response{Success: false, Message: "Cannot make changes. Match is locked."}
		json.NewEncoder(w).Encode(response)
		return
	} 
	if status == 3 {
		response = model.Response{Success: false, Message: "Cannot make changes. Match is completed."}
		json.NewEncoder(w).Encode(response)
		return
	}

	//Checking if other row exist
	tx, _ := database.Db.Begin()
	err2 := tx.QueryRow("SELECT OID, Captain, MVBA, MVBO, MVAR FROM othersTbl WHERE MID=@MID AND UID=@UID", sql.Named("MID", other.MID), sql.Named("UID", other.UID)).Scan(&otherDb.OID, &otherDb.Captain, &otherDb.MVBA, &otherDb.MVBO, &otherDb.MVAR)
	if err2 != nil {
		//Else create new other row
		stmt, err3 := tx.Prepare("INSERT INTO othersTbl(Captain,MVBA,MVBO,MVAR,UID,MID) VALUES(@Captain, @MVBA, @MVBO, @MVAR, @UID, @MID)")
		_, err3 = stmt.Exec(sql.Named("Captain", other.Captain), sql.Named("MVBA", other.MVBA), sql.Named("MVBO", other.MVBO), sql.Named("MVAR", other.MVAR), sql.Named("UID", other.UID), sql.Named("MID", other.MID))
		if err3 != nil {
			tx.Rollback()
			response = model.Response{Success: false, Message: "Error in saving data."}
			json.NewEncoder(w).Encode(response)
			return
		}

		//Also, update the players availability by reducing the available field
		ids := strconv.Itoa(other.Captain) + `, ` + strconv.Itoa(other.MVBA) + `, ` + strconv.Itoa(other.MVBO) + `, ` + strconv.Itoa(other.MVAR)
		updateAvailability := fmt.Sprintf(`UPDATE availabilityTbl SET Available = Available - 1 WHERE UID=@UID AND PID IN (%s)`, ids)
		stmt, err4 := tx.Prepare(updateAvailability)
		_, err4 = stmt.Exec(sql.Named("UID", other.UID))
		if err4 != nil {
			tx.Rollback()
			response = model.Response{Success: false, Message: err4.Error()}
			json.NewEncoder(w).Encode(response)
			return
		}
		tx.Commit()
		response = model.Response{Success: true, Message: "Successfully saved."}
		json.NewEncoder(w).Encode(response)
		return
	}

	//Also update the players availability by adding the removed player
	ids2 := strconv.Itoa(otherDb.Captain) + `, ` + strconv.Itoa(otherDb.MVBA) + `, ` + strconv.Itoa(otherDb.MVBO) + `, ` + strconv.Itoa(otherDb.MVAR)
	updateAvailability := fmt.Sprintf(`UPDATE availabilityTbl SET Available = Available + 1 WHERE UID=@UID AND PID IN (%s)`, ids2)
	stmt, err5 := tx.Prepare(updateAvailability)
	_, err5 = stmt.Exec(sql.Named("UID", other.UID))
	if err5 != nil {
		tx.Rollback()
		log.Println("err5" + err5.Error())
		response = model.Response{Success: false, Message: err5.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	//If exist modify the existing other row
	stmt, err6 := tx.Prepare("UPDATE othersTbl SET Captain=@Captain, MVBA=@MVBA, MVBO=@MVBO, MVAR=@MVAR WHERE OID=@OID")
	_, err6 = stmt.Exec(sql.Named("Captain", other.Captain), sql.Named("MVBA", other.MVBA), sql.Named("MVBO", other.MVBO), sql.Named("MVAR", other.MVAR), sql.Named("OID", otherDb.OID))
	if err6 != nil {
		tx.Rollback()
		log.Println("err6" + err6.Error())
		response = model.Response{Success: false, Message: err6.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	//Also update the players availability by reducing the available field
	ids3 := strconv.Itoa(other.Captain) + `, ` + strconv.Itoa(other.MVBA) + `, ` + strconv.Itoa(other.MVBO) + `, ` + strconv.Itoa(other.MVAR)
	updateAvailability3 := fmt.Sprintf(`UPDATE availabilityTbl SET Available = Available - 1 WHERE UID=@UID AND PID IN (%s)`, ids3)
	stmt, err7 := tx.Prepare(updateAvailability3)
	_, err7 = stmt.Exec(sql.Named("UID", other.UID))
	if err7 != nil {
		tx.Rollback()
		log.Println("err7" + err7.Error())
		response = model.Response{Success: false, Message: err7.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	tx.Commit()
	response = model.Response{Success: true, Message: "Successfully saved."}
	json.NewEncoder(w).Encode(response)
}