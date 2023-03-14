package profile

type UpdateProfilePayload struct {
	Name     string `validate:"required"`
	Password string `validate:"passwordable"`
}
