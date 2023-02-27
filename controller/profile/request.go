package profile

type CompleteNikPayload struct {
	Nik int64 `validate:"required,numeric"`
}

type UpdateProfilePayload struct {
	Name        string `validate:"required"`
	Email       string `validate:"required,email"`
	PhoneNumber string `json:"phone_number" validate:"required"`
}
