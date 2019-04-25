package model

// Match :
type Match struct {
	MID         int    `json:"mid"`
	Team        string `json:"team"`
	Opposition  string `json:"opposition"`
	MatchDate   string `json:"matchDate"`
	Result      *int    `json:"result"`
	IsCompleted int    `json:"isCompleted"`
	Status 		int	   `json:"status"`
}

type Other struct {
	OID	int `json:"oid"`
	Captain	int `json:"captain"`
	MVBA	int `json:"mvba"`
	MVBO	int `json:"mvbo"`
	MVAR	int `json:"mvar"`
	UID	int `json:"uid"`
	MID	int `json:"mid"`
}

type Player struct {
	PID	int `json:"pid"`
	Name	string `json:"name"`
	Role	string `json:"role"`
	Team	string `json:"team"`
}

type OtherDetail struct {
	OID	int `json:"oid"`
	Captain	Player `json:"captain"`
	MVBA	Player `json:"mvba"`
	MVBO	Player `json:"mvbo"`
	MVAR	Player `json:"mvar"`
	UID	int `json:"uid"`
	MID	int `json:"mid"`
}
