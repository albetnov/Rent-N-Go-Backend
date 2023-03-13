package profile

type CompleteNikPayload struct {
	Nik int64 `validate:"required,numeric"`
}

type UpdateProfilePayload struct {
	Name        string `validate:"required"`
	Email       string `validate:"required,email"`
	PhoneNumber string `json:"phone_number" validate:"required,numeric"`
}

type UpdatePasswordPayload struct {
	OldPassword     string `json:"old_password" validate:"required"`
	Password        string `validate:"required,eqfield=ConfirmPassword,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required"`
}
