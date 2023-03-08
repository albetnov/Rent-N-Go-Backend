package user

type CreateUserPayload struct {
	Name        string `validate:"required"`
	Email       string `validate:"required,email"`
	PhoneNumber int    `validate:"required" form:"phone_number"`
	Role        string `validate:"required,oneof=admin user"`
	Nik         int
	Password    string `validate:"omitempty,eqfield=ConfirmPassword,min=8"`
}
