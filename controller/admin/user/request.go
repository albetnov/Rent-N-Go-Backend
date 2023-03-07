package user

type CreateUserPayload struct {
	Name            string `validate:"required"`
	Email           string `validate:"required,email"`
	PhoneNumber     int    `validate:"required"`
	Role            string `validate:"required,oneof=admin user"`
	Nik             int
	Password        string `validate:"eqfield=ConfirmPassword,min=8"`
	ConfirmPassword string `validate:"required_with=Password,min=8"`
}
