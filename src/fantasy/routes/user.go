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
		username string
		responseObj model.Response
	)
	_ = json.NewDecoder(r.Body).Decode(&user)
	tx, _ := database.Db.Begin()
	passwrd = hashing.HashAndSalt([]byte(user.Password))

	//Verifying if security question exist
	secQErr := tx.QueryRow("SELECT Question_Name FROM security_questions WHERE QID=@QID", sql.Named("QID", user.QID)).Scan(&questionName)
	if secQErr != nil {
		fmt.Println("Invalid security question.", secQErr.Error())
		responseObj = model.Response{Success: false, Message: "Invalid security question."}
		json.NewEncoder(w).Encode(responseObj)
		return
	}

	//Verifying if security answer exist for that question
	secAErr := tx.QueryRow("SELECT AID FROM security_answers WHERE QID=@QID AND AID=@AID", sql.Named("QID", user.QID), sql.Named("AID", user.AID)).Scan(&questionAnswer)
	if secAErr != nil {
		fmt.Println("Invalid security answer.", secAErr.Error())
		responseObj = model.Response{Success: false, Message: "Invalid security answer."}
		json.NewEncoder(w).Encode(responseObj)
		return
	}

	//Verifying if username already exist
	usernameDoesntExist := tx.QueryRow("SELECT Username FROM users WHERE Username=@Username", sql.Named("Username", user.Username)).Scan(&username)

	//If everything good, then create a new user. Else throw error.
	if usernameDoesntExist != nil {
		stmt, err := tx.Prepare("INSERT INTO users(Name,Username,Password,QID,AID) VALUES(@Name, @Username, @Password, @QID, @AID)")
		_, err = stmt.Exec(sql.Named("Name", user.Name), sql.Named("Username", user.Username), sql.Named("Password", passwrd), sql.Named("QID", user.QID), sql.Named("AID", user.AID))
		if err != nil {
			log.Println("Rollback for create user")
			fmt.Println("Unable to create user", err.Error())
			tx.Rollback()
			responseObj = model.Response{Success: false, Message: "Error in creating user. Please try again with different username."}
			json.NewEncoder(w).Encode(responseObj)
			return
		}
		responseObj = model.Response{Success: true, Message: "Registered successfully."}
		json.NewEncoder(w).Encode(responseObj)
	} else {
		fmt.Println("Username exist")
		responseObj = model.Response{Success: false, Message: "Username exist. Please try another username."}
		json.NewEncoder(w).Encode(responseObj)
	}

	tx.Commit()
}

//VerifyUser :
func VerifyUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user model.User
	var username string
	var password string
	var uid int
	var responseObj model.Response
	tx, _ := database.Db.Begin()
	_ = json.NewDecoder(r.Body).Decode(&user)
	row := tx.QueryRow("SELECT UID, Username,Password FROM users WHERE Username=@Username", sql.Named("Username", user.Username))
	err := row.Scan(&uid, &username, &password)
	if err != nil {
		log.Println("Rollback for verify user")
		tx.Rollback()
		responseObj = model.Response{Success: false, Message: "Invalid username/password."}
		json.NewEncoder(w).Encode(responseObj)
		return
	}
	tx.Commit()
	var isValid = hashing.ComparePasswords(password, []byte(user.Password))

	if !isValid {
		responseObj = model.Response{Success: false, Message: "Invalid username/password."}
		json.NewEncoder(w).Encode(responseObj)
	} else {
		tokenString, err := authorization.GenerateJWT(user.Username)
		if err != nil {
			responseObj = model.Response{Success: false, Message: "Error generating token."}
			json.NewEncoder(w).Encode(responseObj)
		}
		json.NewEncoder(w).Encode(model.Response2{Success: true, Token: tokenString, UID: uid, Username: username})
	}

}

//GetAnswers :
func GetAnswers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var answers []model.Answer
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		json.NewEncoder(w).Encode(err.Error())
		fmt.Errorf("Error in converting string to int")
		return
	}
	rows, qerr := database.Db.Query("SELECT AID, Answer FROM security_answers WHERE QID=@QID", sql.Named("QID", id))
	if qerr != nil {
		json.NewEncoder(w).Encode(qerr.Error())
		fmt.Errorf("Error in query")
		return
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
			fmt.Errorf("Error in scanning answer query")
		}
		resultAnswer := model.Answer{AID: aid, Answer: ans}
		fmt.Println(resultAnswer)
		answers = append(answers, resultAnswer)
	}
	json.NewEncoder(w).Encode(answers)
}

//GetQuestions :
func GetQuestions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var questions []model.Question
	rows, qerr := database.Db.Query("SELECT * FROM security_questions")
	if qerr != nil {
		json.NewEncoder(w).Encode(qerr.Error())
		fmt.Errorf("Error in query")
		return
	}
	defer rows.Close()
	for rows.Next() {
		var (
			qid int
			ques string
		)
		scanerr := rows.Scan(&qid, &ques)
		if scanerr != nil {
			json.NewEncoder(w).Encode(scanerr.Error())
			fmt.Errorf("Error in scanning answer query")
		}
		resultQuestion := model.Question{QID: qid, QuestionName: ques}
		fmt.Println(resultQuestion)
		questions = append(questions, resultQuestion)
	}
	json.NewEncoder(w).Encode(questions)
}