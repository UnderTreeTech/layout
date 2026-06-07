package model

// BaseRequest is the common base for all HTTP request structs.
// Every API request must embed *BaseRequest to carry the authentication token.
type BaseRequest struct {
	// Token is the user authentication token, passed via JSON body or form field.
	Token string `json:"token" form:"token" validate:"required"`
}
