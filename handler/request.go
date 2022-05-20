package handler

// loginRequest is the request body for login endpoint
type loginRequest struct {
	Email    string `json:"email" validate:"required,exists=users"`
	Password string `json:"password" validate:"required"`
}

// registerRequest is the request body for register endpoint
type registerRequest struct {
	Email    string `json:"email" validate:"required,email,unique=users"`
	Password string `json:"password" validate:"required,min=6"`
}

type googleAuthRequest struct {
	IdToken string `json:"idToken"`
}
