package user

type baseUserPayload struct {
	Name        string `validate:"required"`
	Email       string `validate:"required,email"`
	PhoneNumber int    `validate:"required" form:"phone_number"`
	Role        string `validate:"required,oneof=admin user"`
	Nik         int
}

type CreateUserPayload struct {
	baseUserPayload
	Password string `validate:"required,min=8"`
}

type UpdateUserPayload struct {
	baseUserPayload
	Password string `validate:"passwordable"`
}
