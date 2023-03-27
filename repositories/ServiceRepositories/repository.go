package ServiceRepositories

import (
	"gorm.io/gen"
	"rent-n-go-backend/query"
	"time"
)

var (
	Car    = car{}
	Driver = driver{}
	Tour   = tour{}
)

func activeOrder(db gen.Dao) gen.Dao {
	return db.Where(query.Orders.EndPeriod.Gt(time.Now())).Where(query.Orders.Status.Neq("completed"))
}
