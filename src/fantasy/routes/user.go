package routes

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	authorization "authorization"
	"db"
	"hashing"
	model "models"

	"github.com/gorilla/mux"
)

//CreateUser :
func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var (
		user           model.User
		passwrd        string
		questionName   string
		questionAnswer int
		isQAnsExist    = false
	)
	_ = json.NewDecoder(r.Body).Decode(&user)
	tx, _ := database.Db.Begin()
	passwrd = hashing.HashAndSalt([]byte(user.Password))

	//Verifying if security question exist
	secQErr := tx.QueryRow("SELECT Question_Name FROM security_questions WHERE QID=@QID", sql.Named("QID", user.QID)).Scan(&questionName)
	if secQErr != nil {
		fmt.Println("Oops! something went wrong", secQErr.Error())
	}

	//Verifying if security answer exist for that question
	ansids, secAErr := tx.Query("SELECT AID FROM security_answers WHERE QID=@QID", sql.Named("QID", user.QID))
	if secAErr != nil {
		fmt.Println("Oops! something went wrong", secQErr.Error())
	}
	defer ansids.Close()
	for ansids.Next() {
		err := ansids.Scan(&questionAnswer)
		if err != nil {
			log.Fatal(err)
		}
		if questionAnswer == user.AID {
			isQAnsExist = true
			break
		}
	}

	//If everything good, then create a new user. Else throw error.
	if isQAnsExist {
		stmt, err := tx.Prepare("INSERT INTO users(Name,Username,Password,QID,AID) VALUES(@Name, @Username, @Password, @QID, @AID)")
		_, err = stmt.Exec(sql.Named("Name", user.Name), sql.Named("Username", user.Username), sql.Named("Password", passwrd), sql.Named("QID", user.QID), sql.Named("AID", user.AID))
		if err != nil {
			log.Println("Rollback for create user")
			tx.Rollback()
			json.NewEncoder(w).Encode(err)
			log.Fatal(err)
		}
		json.NewEncoder(w).Encode(true)
	} else {
		fmt.Println("Security answer is invalid", secQErr.Error())
	}

	tx.Commit()
}

//VerifyUser :
func VerifyUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("Success :)")
	// var user model.User
	// var username string
	// var password string
	// tx, _ := database.Db.Begin()
	// _ = json.NewDecoder(r.Body).Decode(&user)
	// err := tx.QueryRow("SELECT Username,Password FROM users WHERE Username=@Username", sql.Named("Username", user.Username)).Scan(&username, &password)
	// var isValid = hashing.ComparePasswords(password, []byte(user.Password))
	// if err != nil {
	// 	log.Println("Rollback for create user")
	// 	tx.Rollback()
	// 	json.NewEncoder(w).Encode(err)
	// 	log.Fatal(err)
	// }

	// tx.Commit()
	// if !isValid {
	// 	json.NewEncoder(w).Encode(isValid)
	// } else {
	// 	tokenString, err := authorization.GenerateJWT(user.Username)
	// 	if err != nil {
	// 		fmt.Println("Error generating token")
	// 		json.NewEncoder(w).Encode("Error generating token")
	// 		log.Fatal("Error generating token")
	// 	}
	// 	fmt.Println(tokenString)
	// 	json.NewEncoder(w).Encode(tokenString)
	// }

}

//GetAnswers :
func GetAnswers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// type answer struct {
	// 	qid int
	// 	ans string
	// }
	var answers []model.Answer
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		log.Fatal("Error in converting string to int")
	}
	rows, qerr := database.Db.Query("SELECT AID, Answer FROM security_answers WHERE QID=@QID", sql.Named("QID", id))
	if qerr != nil {
		json.NewEncoder(w).Encode(qerr.Error())
		log.Fatal("Error in query")
	}
	defer rows.Close()
	for rows.Next() {
		var (
			aid int
			ans string
		)
		scanerr := rows.Scan(&aid, &ans)
		if scanerr != nil {
			json.NewEncoder(w).Encode(scanerr.Error())
			log.Fatal("Error in scanning answer query")
		}
		resultAnswer := model.Answer{AID: aid, Answer: ans}
		fmt.Println(resultAnswer)
		answers = append(answers, resultAnswer)
	}
	json.NewEncoder(w).Encode(answers)
}
