package auth

type LoginPayload struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required"`
}

type RefreshPayload struct {
	RefreshToken string `json:"refresh_token" validate:"required,min=100"`
}
