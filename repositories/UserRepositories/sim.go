package UserRepositories

import (
	"os"
	"path"
	"rent-n-go-backend/models/UserModels"
	"rent-n-go-backend/query"
	"rent-n-go-backend/utils"
)

type simRepository struct {
}

func (s simRepository) GetByUserId(userId uint) (*UserModels.Sim, error) {
	return query.Sim.Where(query.Sim.UserID.Eq(userId)).First()
}

func (sr simRepository) UpdateOrCreate(userId uint, payload *UserModels.Sim) {
	s := query.Sim

	preCond := s.Where(s.UserID.Eq(userId))

	if result, err := preCond.First(); err == nil {
		os.Remove(path.Join(utils.PublicPath(), result.FilePath))
		preCond.Updates(payload)
	} else {
		s.Create(payload)
	}
}
