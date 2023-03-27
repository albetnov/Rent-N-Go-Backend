package UserRepositories

import (
	"errors"
	"gorm.io/gen"
	"rent-n-go-backend/models"
	"rent-n-go-backend/query"
	"rent-n-go-backend/repositories/ServiceRepositories"
	"rent-n-go-backend/utils"
	"sync"
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

var CarIsOutOfStockErr = errors.New("the car is currently out of stock")
var DriverIsNotAvailableErr = errors.New("driver is currently not available at the moment")
var TourIsNotAvailableErr = errors.New("seems like tour has disappeared")
var CarNotFound = errors.New("car not found")

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
		Order(qo.Status).
		Find()
}

func activeOrder(db gen.Dao) gen.Dao {
	return db.Where(query.Orders.Status.Neq(ORDER_COMPLETED))
}

func (o orderRepository) HasOrder(userId uint) bool {
	qo := query.Orders

	total, _ := qo.Scopes(activeOrder).Where(qo.UserId.Eq(userId)).Count()

	return total > 0
}

func (o orderRepository) checkCar(carId uint) (*models.Cars, error) {
	currentStock, car, err := ServiceRepositories.Car.CheckStock(carId)

	if err != nil {
		return nil, CarNotFound
	}

	if currentStock == 0 {
		return nil, CarIsOutOfStockErr
	}

	return car, nil
}

func (o orderRepository) CreateCarOrder(carId uint) error {
	qo := query.Orders

	car, err := o.checkCar(carId)

	if err != nil {
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
		Type:          "car",
	}); err != nil {
		return err
	}

	return nil
}

func (o orderRepository) checkDriver(driverId uint) error {
	if ServiceRepositories.Driver.CheckAvailability(driverId) {
		return DriverIsNotAvailableErr
	}

	return nil
}

func (o orderRepository) CreateDriverOrder(carId, driverId uint) error {
	errCh := make(chan error, 2)
	carCh := make(chan *models.Cars)
	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()
		result, err := o.checkCar(carId)
		if err != nil {
			errCh <- err
			return
		}
		carCh <- result
	}()

	go func() {
		defer wg.Done()
		err := o.checkDriver(driverId)
		if err != nil {
			errCh <- err
		}
	}()

	wg.Wait()
	close(errCh)

	if err := <-errCh; err != nil {
		return err
	}

	car := <-carCh

	driver, _ := ServiceRepositories.Driver.GetById(driverId)

	dayDiff := o.endPeriod.YearDay() - o.startPeriod.YearDay()

	price := (car.Price + driver["price"].(int)) * dayDiff

	qo := query.Orders

	if err := qo.Create(&models.Orders{
		UserId:        o.userId,
		CarId:         &carId,
		Status:        ORDER_ACTIVE,
		StartPeriod:   o.startPeriod,
		EndPeriod:     o.endPeriod,
		TotalAmount:   price,
		PaymentMethod: o.paymentMethod,
		DriverId:      &driverId,
		Type:          "driver",
	}); err != nil {
		return err
	}

	return nil
}

func (o orderRepository) CreateTourOrder(tourId uint) error {
	stock, tour, err := ServiceRepositories.Tour.CheckStock(tourId)

	if err != nil || stock <= 0 {
		return TourIsNotAvailableErr
	}

	errCh := make(chan error)
	wg := new(sync.WaitGroup)

	wg.Add(2)

	go func() {
		defer wg.Done()
		if _, err := o.checkCar(tour.CarId); err != nil {
			errCh <- err
		}
	}()

	go func() {
		defer wg.Done()

		if err := o.checkDriver(tour.DriverId); err != nil {
			errCh <- err
		}
	}()

	wg.Wait()
	close(errCh)

	if err := <-errCh; err != nil {
		return err
	}

	qo := query.Orders
	if err := qo.Create(&models.Orders{
		TourId:        &tourId,
		Type:          "tour",
		TotalAmount:   tour.Price,
		StartPeriod:   o.startPeriod,
		UserId:        o.userId,
		EndPeriod:     o.endPeriod,
		PaymentMethod: o.paymentMethod,
		Status:        ORDER_ACTIVE,
	}); err != nil {
		return err
	}

	return nil
}
