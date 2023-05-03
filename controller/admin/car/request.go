package car

type CarPayload struct {
	Name          string   `validate:"required"`
	Stock         int      `validate:"required,numeric"`
	Desc          string   `validate:"required"`
	Price         int      `validate:"required,numeric"`
	FeaturesIcon  []string `validate:"required" form:"features-icon"`
	FeaturesLabel []string `validate:"required" form:"features-label"`
}

type EditCarPayload struct {
	Name  string `validate:"required"`
	Stock int    `validate:"required,numeric"`
	Desc  string `validate:"required"`
	Price int    `validate:"required,numeric"`
}
