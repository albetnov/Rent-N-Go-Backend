package auth

type LoginRequest struct {
	Email    string `validate:"email,required,min=1"`
	Password string `validate:"required,min=1"`
}
