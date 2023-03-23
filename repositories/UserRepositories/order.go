package UserRepositories

import (
	"gorm.io/gen"
	"rent-n-go-backend/models"
	"rent-n-go-backend/query"
)

type orderRepository struct {
}

func orderPreload(o gen.Dao) gen.Dao {
	qo := query.Orders
	return o.Preload(qo.Car).
		Preload(qo.Driver).
		Preload(qo.Driver.Car).
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
