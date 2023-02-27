package repositories

import (
	"rent-n-go-backend/models"
	"rent-n-go-backend/query"
)

func UpdateOrCreateNik(id uint, payload *models.Nik) {
	nik := query.Nik
	if result, _ := nik.Where(nik.UserID.Eq(id)).Updates(payload); result.RowsAffected <= 0 {
		nik.Create(payload)
	}
}

func GetNikFromUserId(userId uint) (*models.Nik, error) {
	return query.Nik.Where(query.Nik.UserID.Eq(userId)).First()
}
