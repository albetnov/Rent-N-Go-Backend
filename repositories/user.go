package repositories

import (
	"rent-n-go-backend/models"
	"rent-n-go-backend/query"
)

func GetUserByEmail(email string) (*models.User, error) {
	return query.User.Where(query.User.Email.Eq(email)).First()
}

func GetUserById(id uint) (*models.User, error) {
	return query.User.Where(query.User.ID.Eq(id)).First()
}

func GetUserByEmailOrPhone(email string, phone string) (*models.User, error) {
	u := query.User
	return u.Where(u.Email.Eq(email)).Or(u.PhoneNumber.Eq(phone)).First()
}

func CreateUser(user *models.User) error {
	return query.User.Create(user)
}
