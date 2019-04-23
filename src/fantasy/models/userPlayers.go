package model

// Others :
type Availability struct {
	AID    		int    `json:"aid"`
	UID    		int    `json:"uid"`
	PID    		int    `json:"pid"`
	Available   int    `json:"available"`
}

type PostAvailability struct {
	UID    		int    `json:"uid"`
	PID    		int    `json:"pid"`
}