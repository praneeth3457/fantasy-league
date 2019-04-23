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
		//sqlStr += "(?, ?),"
		//availabilities = append(availabilities, uid, pid)
	}
	//trim the last ,
	//sqlStr = strings.TrimSuffix(sqlStr, ",")
	//Replacing ? with $n for postgres
	//sqlStr = ReplaceSQL(sqlStr, "?")
	//stmt, aErr := tx.Prepare(sqlStr)
	//_, aErr = stmt.Exec(availabilities...)
	//_, aErr = stmt.Exec(10001, 40000001)
	// if aErr != nil {
	// 	log.Println("Rollback for inserting availability rows")
	// 	tx.Rollback()
	// 	responseObj = model.Response{Success: false, Message: aErr.Error()}
	// 	json.NewEncoder(w).Encode(responseObj)
	// 	return
	// }
	//sql.Named("Team", match.Team), sql.Named("Team", match.Team)

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