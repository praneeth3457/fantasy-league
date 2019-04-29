package routes

import (
	"encoding/json"
	"net/http"
	"fmt"
	"log"
	"db"
	model "models"
	"github.com/denisenkom/go-mssqldb"
)

//GetAllPlayers :
func GetAllPlayers(w http.ResponseWriter, r *http.Request) {
	var (
		players 	[]model.Player
		response model.Response
		response2 model.ResponsePlayers
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
	response2 = model.ResponsePlayers{Success: true, Message: players}
	json.NewEncoder(w).Encode(response2)
}


/*  === SavePlayersScore ===
	Insert Batting Scores
	Insert Bowling Scores
	Insert Fielding Scores
	Calculate Batting, Bowling, Fielding points
	Insert Points
	Get all IDs from the above updated tables
	Insert PlayerMatches table
*/
func SavePlayersScore(w http.ResponseWriter, r *http.Request) {
	var (
		playersScore model.PlayersScore
		response model.Response
		points []model.Points
		playersMatches []model.PlayerMatches
	)
	_ = json.NewDecoder(r.Body).Decode(&playersScore)

	tx, _ := database.Db.Begin()

	// Insert Batting Scores
	stmt, err := tx.Prepare(mssql.CopyIn("battingTbl", mssql.BulkOptions{}, "Runs", "Balls", "6s", "4s", "isDuck", "MID", "PID"))
	if err != nil {
		response = model.Response{Success: false, Message: "Error in preparing batting"}
		json.NewEncoder(w).Encode(response)
		return
	}
	for _, row := range playersScore.Batting {
		_, err = stmt.Exec(row.Runs, row.Balls, row.Sixes, row.Fours, row.IsDuck, playersScore.MID, row.PID)
		if err != nil {
			response = model.Response{Success: false, Message: "Error in executing batting"}
			tx.Rollback()
			json.NewEncoder(w).Encode(response)
			return
		}
		//Calculating Batting Points
		points = append(points, model.Points{Batting_pts: calculateBatting(row), MID: playersScore.MID, PID: row.PID})
	}
	_, err = stmt.Exec()
    if err != nil {
        log.Fatal(err)
    }
	err = stmt.Close()
	if err != nil {
		response = model.Response{Success: false, Message: "Error in closing batting"}
		tx.Rollback()
		json.NewEncoder(w).Encode(response)
		return
	}


	// Insert Bowling Scores
	stmt, err = tx.Prepare(mssql.CopyIn("bowlingTbl", mssql.BulkOptions{}, "Overs", "Maidens", "Runs", "Wickets", "MID", "PID"))
	if err != nil {
		response = model.Response{Success: false, Message: "Error in preparing bowling"}
		json.NewEncoder(w).Encode(response)
		return
	}
	for _, row := range playersScore.Bowling {
		_, err = stmt.Exec(row.Overs, row.Maidens, row.Runs, row.Wickets, playersScore.MID, row.PID)
		if err != nil {
			response = model.Response{Success: false, Message: "Error in executing bowling"}
			tx.Rollback()
			json.NewEncoder(w).Encode(response)
			return
		}

		//Calculating Bowling Points
		for index, pointsRow := range points {
			if pointsRow.PID == row.PID && pointsRow.MID == playersScore.MID {
				points[index].Bowling_pts = calculateBowling(row)
				break
			}
		}
	}
	_, err = stmt.Exec()
    if err != nil {
        log.Fatal(err)
    }
	err = stmt.Close()
	if err != nil {
		response = model.Response{Success: false, Message: "Error in closing bowling"}
		tx.Rollback()
		json.NewEncoder(w).Encode(response)
		return
	}


	// Insert Fielding Scores
	stmt, err = tx.Prepare(mssql.CopyIn("fieldingTbl", mssql.BulkOptions{}, "Catches", "Stumpings", "Runouts", "MID", "PID"))
	if err != nil {
		response = model.Response{Success: false, Message: "Error in preparing fielding"}
		json.NewEncoder(w).Encode(response)
		return
	}
	for _, row := range playersScore.Fielding {
		_, err = stmt.Exec(row.Catches, row.Stumpings, row.Runouts, playersScore.MID, row.PID)
		if err != nil {
			response = model.Response{Success: false, Message: "Error in executing fielding"}
			tx.Rollback()
			json.NewEncoder(w).Encode(response)
			return
		}
		//Calculating Fielding Points
		for index, pointsRow := range points {
			if pointsRow.PID == row.PID && pointsRow.MID == playersScore.MID {
				points[index].Fielding_pts = calculateFielding(row)
				break
			}
		}
	}
	_, err = stmt.Exec()
    if err != nil {
        log.Fatal(err)
    }
	err = stmt.Close()
	if err != nil {	
		response = model.Response{Success: false, Message: "Error in closing fielding"}
		tx.Rollback()
		json.NewEncoder(w).Encode(response)
		return
	}



	// Insert Points
	stmt, err = tx.Prepare(mssql.CopyIn("pointsTbl", mssql.BulkOptions{}, "Batting_pts", "Bowling_pts", "Fielding_pts", "Other_pts", "Total_pts", "MID", "PID"))
	if err != nil {
		response = model.Response{Success: false, Message: "Error in preparing points"}
		json.NewEncoder(w).Encode(response)
		return
	}
	for _, row := range points {
		Total_pts := row.Batting_pts + row.Bowling_pts + row.Fielding_pts + row.Other_pts
		_, err = stmt.Exec(row.Batting_pts, row.Bowling_pts, row.Fielding_pts, row.Other_pts, Total_pts, playersScore.MID, row.PID)
		if err != nil {
			response = model.Response{Success: false, Message: "Error in executing points"}
			tx.Rollback()
			json.NewEncoder(w).Encode(response)
			return
		}
	}
	_, err = stmt.Exec()
    if err != nil {
        log.Fatal(err)
    }
	err = stmt.Close()
	if err != nil {
		response = model.Response{Success: false, Message: "Error in closing points"}
		tx.Rollback()
		json.NewEncoder(w).Encode(response)
		return
	}


	//Getting all IDs from the above tables
	allIds, err2 := tx.Query("SELECT b.MID, b.PID, b.BAID, bo.BOID, f.FID, p.PTID FROM [battingTbl] b JOIN [bowlingTbl] bo ON b.MID = bo.MID AND b.PID = bo.PID JOIN [fieldingTbl] f ON b.MID = f.MID AND b.PID = f.PID JOIN [pointsTbl] p ON b.MID = p.MID AND b.PID = p.PID")
	if err2 != nil {
		response = model.Response{Success: false, Message: "Error in getting ids"}
		tx.Rollback()
		json.NewEncoder(w).Encode(response)
		return
	}
	defer allIds.Close()
	for allIds.Next() {
		var (
			playersMatch model.PlayerMatches
		)
		scanerr := allIds.Scan(&playersMatch.MID, &playersMatch.PID, &playersMatch.BAID, &playersMatch.BOID, &playersMatch.FID, &playersMatch.PTID)
		if scanerr != nil {
			json.NewEncoder(w).Encode(scanerr.Error())
			fmt.Errorf("Error in scanning answer query")
		}
		playersMatches = append(playersMatches, playersMatch)
	}


	//Insert PlayerMatches table
	stmt, err = tx.Prepare(mssql.CopyIn("playerMatchesTbl", mssql.BulkOptions{}, "MID", "PID", "BAID", "BOID", "FID", "PTID"))
	if err != nil {
		response = model.Response{Success: false, Message: "Error in preparing playerMatches"}
		json.NewEncoder(w).Encode(response)
		return
	}
	for _, row := range playersMatches {
		_, err = stmt.Exec(row.MID, row.PID, row.BAID, row.BOID, row.FID, row.PTID)
		if err != nil {
			response = model.Response{Success: false, Message: "Error in executing playerMatches"}
			tx.Rollback()
			json.NewEncoder(w).Encode(response)
			return
		}
	}
	_, err = stmt.Exec()
    if err != nil {
        log.Fatal(err)
    }
	err = stmt.Close()
	if err != nil {
		response = model.Response{Success: false, Message: "Error in closing playerMatches"}
		tx.Rollback()
		json.NewEncoder(w).Encode(response)
		return
	}



	//Commit all the queries
	err = tx.Commit()
	if err != nil {
		fmt.Println("fgjhklasdfghj")
		log.Fatal(err)
	}
	response = model.Response{Success: true, Message: "Successfully updated."}
	json.NewEncoder(w).Encode(response)
}

func calculateBatting(row model.Batting) int {
	var (
		score int = 0
	)

	//Each Run
	score += row.Runs

	//Runs (20-29), (30-49) & (50+)
	if row.Runs >= 20 && row.Runs < 30 {
		score += 5
	} else if row.Runs >= 30 && row.Runs < 50 {
		score += 10
	} else if row.Runs >= 50 {
		score += 15
	}

	//Sixes
	if row.Sixes > 0 {
		score += (row.Sixes * 2)
	}

	//Fours
	if row.Fours > 0 {
		score += row.Fours
	}

	//Duck
	if row.IsDuck == 1 {
		score -= 5
	}

	//Strike rates (100 - 119), (120 - 149) && (150+)
	if row.Balls > 5 {
		strikeRate  := (float64(row.Runs)/float64(row.Balls))*100
		if strikeRate >= 100 && strikeRate < 120 {
			score += 5
		} else if strikeRate >= 120 && strikeRate < 150 {
			score += 10
		} else if strikeRate >= 150 {
			score += 15
		}
	}

	return score
}


func calculateBowling(row model.Bowling) int {
	var (
		bowlingScore int = 0
	)

	//Each Wicket
	bowlingScore += row.Wickets * 10

	//Wickets 2, 3, 4, 5
	if row.Wickets == 2 {
		bowlingScore += 5
	} else if row.Wickets == 3 {
		bowlingScore += 10
	} else if row.Wickets == 4 {
		bowlingScore += 15
	} else if row.Wickets >= 5 {
		bowlingScore += 20
	}

	//Economy 
	if(row.Overs >= 1) {
		economy := float64(row.Runs)/row.Overs
		if economy < 4 {
			bowlingScore += 15
		} else if economy >= 4 && economy < 5 {
			bowlingScore += 10
		} else if economy >= 5 && economy < 6 {
			bowlingScore += 5
		} else if economy >= 8 {
			bowlingScore -= 5
		}
	}

	//Maiden over
	if row.Maidens > 0 {
		bowlingScore += row.Maidens * 5
	}

	return bowlingScore
}


func calculateFielding(row model.Fielding) int {
	var (
		fieldingScore int = 0
	)

	//Each Catch
	fieldingScore += row.Catches * 5

	//Each Stumping
	fieldingScore += row.Stumpings * 10

	//Each Runout
	fieldingScore += row.Runouts * 5

	return fieldingScore
}