package UserRepositories

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gen"
	"rent-n-go-backend/models"
	"rent-n-go-backend/query"
	"rent-n-go-backend/repositories/BasicRepositories"
	"rent-n-go-backend/repositories/ServiceRepositories"
	"rent-n-go-backend/utils"
	"sync"
	"time"
)

type orderRepository struct {
	c             *fiber.Ctx
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

func (o *orderRepository) CreateOrder(ctx *fiber.Ctx, startPeriod, endPeriod, paymentMethod string, userId uint) orderRepository {
	o.startPeriod = utils.ParseISO8601Date(startPeriod)
	o.endPeriod = utils.ParseISO8601Date(endPeriod)
	o.paymentMethod = paymentMethod
	o.userId = userId
	o.c = ctx

	return *o
}

func withCar(db gen.Dao) gen.Dao {
	qo := query.Orders
	return db.Preload(qo.Car).
		Preload(qo.Car.Features.On(query.Features.Associate.Eq(BasicRepositories.Car))).
		Preload(qo.Car.Pictures.On(query.Pictures.Associate.Eq(BasicRepositories.Car)))
}

func withDriver(db gen.Dao) gen.Dao {
	qo := query.Orders
	return db.Preload(qo.Driver).
		Preload(qo.Driver.Features.On(query.Features.Associate.Eq(BasicRepositories.Driver))).
		Preload(qo.Driver.Pictures.On(query.Pictures.Associate.Eq(BasicRepositories.Driver)))
}

func withTour(db gen.Dao) gen.Dao {
	qo := query.Orders
	return db.Preload(qo.Tour).
		Preload(qo.Tour.Car).
		Preload(qo.Tour.Driver).
		Preload(qo.Tour.Features.On(query.Features.Associate.Eq(BasicRepositories.Tour))).
		Preload(qo.Tour.Pictures.On(query.Pictures.Associate.Eq(BasicRepositories.Tour)))
}

func orderPreload(db gen.Dao) gen.Dao {
	return db.
		Scopes(withCar).
		Scopes(withDriver).
		Scopes(withTour)

}

func (o orderRepository) GetUserOrder(userId uint, c *fiber.Ctx, filter string) ([]*models.Orders, int64, error) {
	qo := query.Orders

	builder := qo.Scopes(orderPreload)

	if userId != 0 {
		builder = builder.Where(qo.UserId.Eq(userId))
	}

	if filter != "" {
		builder = builder.Where(qo.Type.Eq(filter))
	}

	total, err := builder.Count()

	if err != nil {
		return nil, 0, err
	}

	orders, err := builder.
		Scopes(utils.Paginate(c)).
		Order(qo.Status, qo.ID).
		Find()

	if err != nil {
		return nil, 0, err
	}

	wg := new(sync.WaitGroup)
	qp := query.Pictures
	qf := query.Features

	if filter == "" || filter == "tour" {
		for _, v := range orders {
			wg.Add(4)
			go func(v *models.Orders) {
				tourCarPictures, _ := query.Cars.Pictures.Where(qp.Associate.Eq(BasicRepositories.Car)).Model(&v.Tour.Car).Find()
				for _, p := range tourCarPictures {
					v.Tour.Car.Pictures = append(v.Tour.Car.Pictures, *p)
				}
				wg.Done()
			}(v)

			go func(v *models.Orders) {
				tourCarFeatures, _ := query.Cars.Features.Where(qf.Associate.Eq(BasicRepositories.Car)).Model(&v.Tour.Car).Find()
				for _, f := range tourCarFeatures {
					v.Tour.Car.Features = append(v.Tour.Car.Features, *f)
				}
				wg.Done()
			}(v)

			go func(v *models.Orders) {
				tourDriverPictures, _ := query.Driver.Pictures.Where(qp.Associate.Eq(BasicRepositories.Driver)).Model(&v.Tour.Driver).Find()
				for _, p := range tourDriverPictures {
					v.Tour.Driver.Pictures = append(v.Tour.Driver.Pictures, *p)
				}
				wg.Done()
			}(v)

			go func(v *models.Orders) {
				tourDriverFeatures, _ := query.Driver.Features.Where(qf.Associate.Eq(BasicRepositories.Driver)).Model(&v.Tour.Driver).Find()
				for _, f := range tourDriverFeatures {
					v.Tour.Driver.Features = append(v.Tour.Driver.Features, *f)
				}
				wg.Done()
			}(v)
		}

		wg.Wait()
	}

	return orders, total, err
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
	currentStock, car, err := ServiceRepositories.Car.Ctx(o.c).CheckStock(carId)

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
	if ServiceRepositories.Driver.Ctx(o.c).CheckAvailability(driverId) {
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

	driver, _ := ServiceRepositories.Driver.Ctx(o.c).GetById(driverId)

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
	stock, tour, err := ServiceRepositories.Tour.Ctx(o.c).CheckStock(tourId)

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
