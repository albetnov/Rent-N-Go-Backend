package auth

type LoginPayload struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required"`
}
