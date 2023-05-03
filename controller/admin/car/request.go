package car

type CarPayload struct {
	Name    string `validate:"required"`
	Stock   int    `validate:"required,numeric"`
	Desc    string `validate:"required"`
	Price   int    `validate:"required,numeric"`
	Seats   int    `validate:"required,numeric"`
	Baggage int    `validate:"required,numeric"`
}
