package BasicRepositories

import "errors"

var (
	Features = featuresRepository{}
	Pictures = picturesRepository{}
)

const (
	Car    = "car"
	Driver = "driver"
)

var ErrInvalidArgument = errors.New("invalid argument being passed")
var ErrNotFound = errors.New("invalid id being passed. Target association not found")

func checkAssociation(associate string) error {
	if associate != "car" && associate != "driver" {
		return ErrInvalidArgument
	}

	return nil
}
