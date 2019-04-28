package model

//
type Response struct {
	Success bool `json:"success"`
	Message string `json:"message"`
}

type Response2 struct {
	Success bool `json:"success"`
	Token string `json:"token"`
	UID int `json:"uid"`
	Username string `json:"username"`
	Role int `json:"role"`
	IsStarted int `json:"isStarted"`
}

type Response3 struct {
	Success bool `json:"success"`
	Message OtherDetail `json:"message"`
}

type ResponsePlayers struct {
	Success bool `json:"success"`
	Message []Player `json:"message"`
}