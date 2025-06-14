package handlers

import "time" //time.Now().Unix() - may be => can be replaced by INTEGER

type DefaultResponse struct {
	Type    string `json:"type"`    // Error | Data | Message
	Message string `json:"message"` // Message
}

type CreatePayload struct { //POST, URL: projectId=int;
	Name string `json:"name"`
}

type CreateResponse struct { //POST
	ID          int       `json:"id"`
	ProjectID   int       `json:"projectId"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Priority    int       `json:"priority"`
	Removed     bool      `json:"removed"`   //false
	CreatedAt   time.Time `json:"createdAt"` //time.Now().Unix() - may be
}

type UpdatePayload struct { //PATCH, URL: id=int & projectId=int //check is exist
	Name        string `json:"name"` //check is nil (should fil)
	Description string `json:"description"`
}

type UpdateResponse struct { //PATCH
	ID          int       `json:"id"`
	ProjectID   int       `json:"projectId"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Priority    int       `json:"priority"`
	Removed     bool      `json:"removed"`   //false
	CreatedAt   time.Time `json:"createdAt"` //time.Now().Unix() - may be
}

//Payload(nil); URL: id=int & projectId=int //check is exist

type DeleteResponse struct { //DELETE
	ID         int  `json:"id"`
	CampaignID int  `json:"campaignId"`
	Removed    bool `json:"removed"` //true
}

//Payload(nil); URL: limit=int & offset=int

type GetListResponse struct {
	Meta  Meta    `json:"meta"`
	Goods []Goods `json:"goods"`
}

type Meta struct {
	Total   int `json:"total"`
	Removed int `json:"removed"`
	Limit   int `json:"limit"`
	Offset  int `json:"offset"`
}

type Goods struct {
	ID          int       `json:"id"`
	ProjectID   int       `json:"projectId"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Priority    int       `json:"priority"`
	Removed     bool      `json:"removed"`   //false
	CreatedAt   time.Time `json:"createdAt"` //time.Now().Unix() - may be
}

type ReoprioritizePayload struct { //PATCH, URL: id=int & projectId=int // check is exist
	NewPriority int `json:"newPriority"`
}

type ReoprioritizeResponse struct {
	Priorities Priorities `json:"priorities"`
}

type Priorities struct {
	ID       int `json:"id"`
	Priority int `json:"priority"`
}
