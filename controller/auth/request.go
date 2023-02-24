package auth

type LoginPayload struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required"`
}

type RefreshPayload struct {
	RefreshToken string `json:"refresh_token" validate:"required,min=100"`
}

type RegisterPayload struct {
	Name            string `validate:"required"`
	Email           string `validate:"required,email"`
	PhoneNumber     string `json:"phone_number" validate:"required"`
	Password        string `validate:"required,eqfield=ConfirmPassword"`
	ConfirmPassword string `json:"confirm_password" validate:"required"`
}
