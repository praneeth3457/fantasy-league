package model

// Others :
type Availability struct {
	AID    		int    `json:"aid"`
	UID    		int    `json:"uid"`
	PID    		int    `json:"pid"`
	Name		string	`json:"name"`
	Role		string	`json:"role"`
	Team		string	`json:"team"`
	Available   int    `json:"available"`
}

type PostAvailability struct {
	UID    		int    `json:"uid"`
	PID    		int    `json:"pid"`
}

type AvailabilityResponse struct {
	Success bool `json:"success"`
	Message []Availability `json:"message"`
}