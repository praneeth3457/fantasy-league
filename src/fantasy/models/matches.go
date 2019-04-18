package model

// Match :
type Match struct {
	MID         int    `json:"mid"`
	Team        string `json:"team"`
	Opposition  string `json:"opposition"`
	MatchDate   string `json:"matchDate"`
	Result      int    `json:"result"`
	IsCompleted int    `json:"isCompleted"`
}
