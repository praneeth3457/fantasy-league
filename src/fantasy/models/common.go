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
}