package response

import (
	"encoding/json"
	"net/http"
	"time"
) //time.Now().Unix() - may be => can be replaced by INTEGER

func (r DefaultResponse) Response(w http.ResponseWriter, header int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(header)
	json.NewEncoder(w).Encode(DefaultResponse{Type: r.Type, Message: r.Message})
}

type DefaultResponse struct {
	Type    string `json:"type"`    // Error | Data | Message
	Message string `json:"message"` // Message
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

type CreatePayload struct { //POST, URL: projectId=int;
	Name string `json:"name"`
}

type UpdatePayload struct { //PATCH, URL: id=int & projectId=int //check is exist
	Name        string `json:"name"` //check is nil (should fil)
	Description string `json:"description"`
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

type ReoprioritizePayload struct { //PATCH, URL: id=int & projectId=int // check is exist
	NewPriority *int `json:"newPriority"`
}

type ReoprioritizeResponse struct {
	Priorities []Priorities `json:"priorities"`
}

type Priorities struct {
	ID       int `json:"id"`
	Priority int `json:"priority"`
}

type NatsForClick struct {
	Id          uint32    `json:"Id"`
	ProjectId   uint32    `json:"ProjectId"`
	Name        string    `json:"Name"`
	Description string    `json:"Description"`
	Priority    uint32    `json:"Priority"`
	Removed     uint8     `json:"Removed"`
	EventTime   time.Time `json:"EventTime"`
}
