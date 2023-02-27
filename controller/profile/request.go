package profile

type CompleteNikPayload struct {
	Nik int64 `validate:"required,numeric"`
}
