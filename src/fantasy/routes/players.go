package routes

import (
	"encoding/json"
	"net/http"
	"fmt"

	"db"
	model "models"
)

//GetAllMatches :
func GetAllPlayers(w http.ResponseWriter, r *http.Request) {
	var (
		players 	[]model.Player
		response model.Response
	)
	rows, err := database.Db.Query("SELECT * FROM playersTbl")
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
			player model.Player
		)
		scanerr := rows.Scan(&player.PID, &player.Name, &player.Role, &player.Team)
		if scanerr != nil {
			fmt.Errorf("Error in scanning answer query")
			response = model.Response{Success: false, Message: scanerr.Error()}
			json.NewEncoder(w).Encode(response)
			return
		}
		resultMatch := model.Player{PID: player.PID, Name: player.Name, Role: player.Role, Team: player.Team}
		players = append(players, resultMatch)
	}
	json.NewEncoder(w).Encode(players)
}