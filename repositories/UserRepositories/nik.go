package UserRepositories

import (
	"rent-n-go-backend/models/UserModels"
	"rent-n-go-backend/query"
	"strconv"
)

type nikRepository struct{}

func (n nikRepository) UpdateOrCreate(id uint, payload *UserModels.Nik) {
	nik := query.Nik
	if result, _ := nik.Where(nik.UserID.Eq(id)).Updates(payload); result.RowsAffected <= 0 {
		nik.Create(payload)
	}
}

func (n nikRepository) GetFromUserId(userId uint) (*UserModels.Nik, error) {
	return query.Nik.Where(query.Nik.UserID.Eq(userId)).First()
}

func (n nikRepository) OptionalCreate(nik int, userId uint) {
	if nik != 0 {
		Nik.UpdateOrCreate(userId, &UserModels.Nik{
			UserID:     userId,
			Nik:        strconv.FormatInt(int64(nik), 10),
			IsVerified: false,
		})
	}
}
