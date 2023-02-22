package auth

type RequestPayload struct {
	Username string `validate:"required"`
	Password string `validate:"required"`
}
