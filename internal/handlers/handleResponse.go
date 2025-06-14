package handlers

type DefaultResponse struct {
	Type    string `json:"type"`    // Error | Data | Message
	Message string `json:"message"` // Message
}
