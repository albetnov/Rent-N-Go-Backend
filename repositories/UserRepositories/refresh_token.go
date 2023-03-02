package UserRepositories

import (
	"gorm.io/gen"
	"rent-n-go-backend/models/UserModels"
	"rent-n-go-backend/query"
)

type refreshTokenRepository struct {
}

func (rtr refreshTokenRepository) GetByUserId(userId uint) (*UserModels.RefreshToken, error) {
	return query.RefreshToken.Where(query.RefreshToken.UserID.Eq(userId)).First()
}

func (rtr refreshTokenRepository) DeleteByTokenId(id uint) (gen.ResultInfo, error) {
	return query.RefreshToken.Where(query.RefreshToken.ID.Eq(id)).Delete()
}

func (rtr refreshTokenRepository) UpdateOrCreateByUserId(id uint, payload *UserModels.RefreshToken) {
	rt := query.RefreshToken
	if result, _ := rt.Where(rt.UserID.Eq(id)).Updates(payload); result.RowsAffected <= 0 {
		rt.Create(payload)
	}
}

func (rtr refreshTokenRepository) DeleteByUserId(userId uint) (gen.ResultInfo, error) {
	return query.RefreshToken.Where(query.RefreshToken.UserID.Eq(userId)).Delete()
}
