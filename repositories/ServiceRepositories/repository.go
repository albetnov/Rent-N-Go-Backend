package ServiceRepositories

import (
	"gorm.io/gen"
	"rent-n-go-backend/query"
)

var (
	Car    = car{}
	Driver = driver{}
	Tour   = tour{}
)

func activeOrder(db gen.Dao) gen.Dao {
	return db.Where(query.Orders.Status.Eq("active"))
}
