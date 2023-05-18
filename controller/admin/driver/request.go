package driver

type DriverPayload struct {
	Name  string `validate:"required"`
	Desc  string `validate:"required"`
	Price int    `validate:"required,numeric"`
}
