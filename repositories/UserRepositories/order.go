package UserRepositories

import (
	"errors"
	"gorm.io/gen"
	"rent-n-go-backend/models"
	"rent-n-go-backend/query"
	"rent-n-go-backend/repositories/ServiceRepositories"
	"rent-n-go-backend/utils"
	"time"
)

type orderRepository struct {
	startPeriod   time.Time
	endPeriod     time.Time
	userId        uint
	paymentMethod string
}

const ORDER_COMPLETED = "completed"
const ORDER_ACTIVE = "active"

var CarIsOutOfStockErr = errors.New("The car is currently out of stock")

func (o *orderRepository) CreateOrder(startPeriod, endPeriod, paymentMethod string, userId uint) orderRepository {
	o.startPeriod = utils.ParseISO8601Date(startPeriod)
	o.endPeriod = utils.ParseISO8601Date(endPeriod)
	o.paymentMethod = paymentMethod
	o.userId = userId

	return *o
}

func orderPreload(o gen.Dao) gen.Dao {
	qo := query.Orders
	return o.Preload(qo.Car).
		Preload(qo.Driver).
		Preload(qo.Tour).
		Preload(qo.Tour.Car).
		Preload(qo.Tour.Driver)
}

func (o orderRepository) GetUserOrder(userId uint) ([]*models.Orders, error) {
	qo := query.Orders
	return qo.
		Scopes(orderPreload).
		Where(qo.UserId.Eq(userId)).
		Find()
}

func (o orderRepository) CreateCarOrder(carId uint) error {
	qo := query.Orders

	currentStock, car, err := ServiceRepositories.Car.CheckStock(carId)

	if err != nil {
		return err
	}

	if currentStock == 0 {
		return CarIsOutOfStockErr
	}

	if _, err := User.GetById(o.userId); err != nil {
		return err
	}

	dayDiff := o.endPeriod.YearDay() - o.startPeriod.YearDay()

	if err := qo.Create(&models.Orders{
		UserId:        o.userId,
		CarId:         &carId,
		Status:        ORDER_ACTIVE,
		StartPeriod:   o.startPeriod,
		EndPeriod:     o.endPeriod,
		TotalAmount:   car.Price * dayDiff,
		PaymentMethod: o.paymentMethod,
	}); err != nil {
		return err
	}

	return nil
}
