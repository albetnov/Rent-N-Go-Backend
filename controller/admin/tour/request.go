package tour

type TourPayload struct {
	Name     string `validate:"required"`
	Stock    int    `validate:"required,numeric"`
	Desc     string `validate:"required"`
	Price    int    `validate:"required,numeric"`
	CarID    int
	DriverID int
}
