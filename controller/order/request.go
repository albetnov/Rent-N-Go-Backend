package order

type PlaceOrderPayload struct {
	CarId         uint   `json:"car_id" validate:"required"`
	DriverId      uint   `json:"driver_id" validate:"omitempty,numeric,required_with=CarId"`
	TourId        uint   `json:"tour_id" validate:"omitempty,numeric,required_with_all=DriverId CarId"`
	PaymentMethod string `json:"payment_method" validate:"required"`
	StartPeriod   string `json:"start_period" validate:"required,ISO8601date"`
	EndPeriod     string `json:"end_period" validate:"required,ISO8601date,afteriso=StartPeriod"`
}
