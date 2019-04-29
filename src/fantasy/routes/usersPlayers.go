package routes

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"fmt"
	"strings"
	"strconv"

	"db"
	model "models"
)

func CreateAvailability(w http.ResponseWriter, r *http.Request) {
	type Value []interface{}
	var (
		uid int
		response model.Response
		responseObj model.Response
	)

	err := database.Db.QueryRow("SELECT UID FROM users WHERE Username=@Username AND isStarted=0", sql.Named("Username", r.Header["User-Context"][0])).Scan(&uid)
	if err != nil {
		fmt.Errorf("Error in finding user")
		response = model.Response{Success: false, Message: "Unable to find user."}
		json.NewEncoder(w).Encode(response)
		return
	}

	playerRows, pErr := database.Db.Query("SELECT PID FROM playersTbl")
	if pErr != nil {
		fmt.Errorf("Error in finding players")
		response = model.Response{Success: false, Message: "Unable to find players."}
		json.NewEncoder(w).Encode(response)
		return
	}
	defer playerRows.Close()

	tx, _ := database.Db.Begin()
	//sqlStr := "INSERT INTO availabilityTbl(UID, PID) VALUES "
	//availabilities := []interface{}{}
	stmt, aErr := tx.Prepare("INSERT INTO availabilityTbl(UID, PID) VALUES (@UID, @PID)")
	for playerRows.Next() {
		var (
			pid int
		)
		scanerr := playerRows.Scan(&pid)
		if scanerr != nil {
			fmt.Errorf("Error in scanning playerRows")
			response = model.Response{Success: false, Message: "Error in scanning playerRows"}
			json.NewEncoder(w).Encode(response)
			return
		}
		_, aErr = stmt.Exec(sql.Named("UID", uid), sql.Named("PID", pid))
		if aErr != nil {
			log.Println("Rollback for inserting availability rows")
			tx.Rollback()
			responseObj = model.Response{Success: false, Message: aErr.Error()}
			json.NewEncoder(w).Encode(responseObj)
			return
		}
	}

	stmt, aErr = tx.Prepare("UPDATE users SET isStarted = @isStarted Where Username=@Username")
	_, aErr = stmt.Exec(sql.Named("isStarted", 1), sql.Named("Username", r.Header["User-Context"][0]))
		if aErr != nil {
			log.Println("Rollback for inserting users isStarted")
			tx.Rollback()
			responseObj = model.Response{Success: false, Message: aErr.Error()}
			json.NewEncoder(w).Encode(responseObj)
			return
		}
	
	tx.Commit()
	responseObj = model.Response{Success: true, Message: "Successfully created."}
	json.NewEncoder(w).Encode(responseObj)
}

// ReplaceSQL replaces the instance occurrence of any string pattern with an increasing $n based sequence
func ReplaceSQL(old, searchPattern string) string {
	tmpCount := strings.Count(old, searchPattern)
	for m := 1; m <= tmpCount; m++ {
	   old = strings.Replace(old, searchPattern, "$"+strconv.Itoa(m), 1)
	}
	return old
}

func GetAvailability(w http.ResponseWriter, r *http.Request) {
	var (
		user model.User
		response model.Response
		pids []int
		availabilities []model.Availability
	)
	_ = json.NewDecoder(r.Body).Decode(&user)
	rows, err := database.Db.Query("SELECT * FROM availabilityTbl WHERE UID=@UID", sql.Named("UID", user.UID))
	if err != nil {
		response = model.Response{Success: false, Message: "no rows in result set"}
		json.NewEncoder(w).Encode(response)
		return
	}

	defer rows.Close()
	for rows.Next() {
		var (
			pid			int
			aid			int
			uid			int
			available	int
		)
		scanerr := rows.Scan(&aid, &uid, &pid, &available)
		if scanerr != nil {
			response = model.Response{Success: false, Message: "no rows-2 in result set"}
			json.NewEncoder(w).Encode(response)
			return
		}

		pids = append(pids, pid)
		availabilities = append(availabilities, model.Availability{AID: aid, UID: uid, PID: pid, Name: "", Role: "", Team: "", Available: available})
	}
	params2 := make([]string, 0, len(pids))
	for i := range pids {
		params2 = append(params2, strconv.Itoa(pids[i]))
	}

	var some interface{} = strings.Join(params2, ", ")
	q2 := fmt.Sprintf(`SELECT * FROM playersTbl WHERE PID IN (%s) ORDER BY Name`, some)
	playerRows, playerserr := database.Db.Query(q2)
	if playerserr != nil {
		response = model.Response{Success: false, Message: playerserr.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}
	defer playerRows.Close()
	for playerRows.Next() {
		
		var (
			player model.Player
		)
		scanerr2 := playerRows.Scan(&player.PID, &player.Name, &player.Role, &player.Team)
		if scanerr2 != nil {
			response = model.Response{Success: false, Message: "no rows-3 in result set"}
			json.NewEncoder(w).Encode(response)
			return
		}

		for i := range availabilities {
			if availabilities[i].PID == player.PID {
				availabilities[i].Name = player.Name
				availabilities[i].Role = player.Role
				availabilities[i].Team = player.Team
				break
			}
		}
	}

	json.NewEncoder(w).Encode(model.AvailabilityResponse{Success: true, Message: availabilities})
}

