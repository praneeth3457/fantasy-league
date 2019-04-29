package model

type Batting struct {
	BAID int `json:"baid"`
	Runs int `json:"runs"`
	Balls int `json:"balls"`
	Sixes int `json:"sixes"`
	Fours int `json:"fours"`
	IsDuck int `json:"isDuck"`
	MID int `json:"mid"`
	PID int `json:"pid"`
}

type Bowling struct {
	BOID int `json:"boid"`
	Overs float64 `json:"overs"`
	Maidens int `json:"maidens"`
	Runs int `json:"runs"`
	Wickets int `json:"wickets"`
	MID int `json:"mid"`
	PID int `json:"pid"`
}

type Fielding struct {
	FID int `json:"fid"`
	Catches int `json:"catches"`
	Stumpings int `json:"stumpings"`
	Runouts int `json:"runouts"`
	MID int `json:"mid"`
	PID int `json:"pid"`
}

type Points struct {
	PTID int `json:"baid"`
	Batting_pts int `json:"battingsPts"`
	Bowling_pts int `json:"bowlingPts"`
	Fielding_pts int `json:"fieldingPts"`
	Other_pts int `json:"otherPts"`
	Total_pts int `json:"totalPts"`
	MID int `json:"mid"`
	PID int `json:"pid"`
}

type PlayerMatches struct {
	PMID int `json:"pmid"`
	MID int `json:"mid"`
	PID int `json:"pid"`
	BAID int `json:"baid"`
	BOID int `json:"boid"`
	FID int `json:"fid"`
	PTID int `json:"baid"`
}

type PlayersScore struct {
	MID int `json:"mid"`
	Batting []Batting `json:"batting"`
	Bowling []Bowling `json:"bowling"`
	Fielding []Fielding `json:"fielding"`
}