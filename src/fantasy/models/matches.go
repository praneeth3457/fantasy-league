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
	Type	string `json:"type"`
}

type OtherDetail struct {
	OID	int `json:"oid"`
	Players []Player `json:"players"`
	UID	int `json:"uid"`
	MID	int `json:"mid"`
}

type Points2 struct {
	Role string `json:"role"`
	PTID int `json:"ptid"`
	Batting_pts int `json:"battingsPts"`
	Bowling_pts int `json:"bowlingPts"`
	Fielding_pts int `json:"fieldingPts"`
	Other_pts int `json:"otherPts"`
	Total_pts int `json:"totalPts"`
	MID int `json:"mid"`
	PID int `json:"pid"`
}
type TotalPoints struct {
	Batting_pts int `json:"battingsPts"`
	Bowling_pts int `json:"bowlingPts"`
	Fielding_pts int `json:"fieldingPts"`
	Total_pts int `json:"totalPts"`
}

type GetPoints struct {
	MID int `json:"mid"`
	Opposition string `json:"opposition"`
	MatchDate string `json:"matchDate"`
	Points []Points2 `json:"points"`
	TotalPoints TotalPoints `json:"totalPoints"`
}

type ID struct {
	UID int `json:"uid"`
}

type MatchStatus struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Value int `json:"value"`
	MID int `json:"mid"`
}
