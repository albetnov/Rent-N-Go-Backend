package order

type PlaceOrderPayload struct {
	DriverId      uint   `json:"driver_id" validate:"omitempty,numeric"`
	CarId         uint   `json:"car_id" validate:"numeric,required_with=DriverId"`
	TourId        uint   `json:"tour_id" validate:"omitempty,numeric"`
	PaymentMethod string `json:"payment_method" validate:"required"`
	StartPeriod   string `json:"start_period" validate:"required,ISO8601date"`
	EndPeriod     string `json:"end_period" validate:"required,ISO8601date,afteriso=StartPeriod"`
}
