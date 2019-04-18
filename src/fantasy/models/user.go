package model

// User :
type User struct {
	UID      int
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
	AID      int    `json:"aid"`
	QID      int    `json:"qid"`
	Role     int    `json:"role"`
}

// Answer :
type Answer struct {
	AID    int    `json:"aid"`
	Answer string `json:"answer"`
	QID    int    `json:"qid"`
}
