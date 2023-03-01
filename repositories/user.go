package repositories

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gen"
	"rent-n-go-backend/models"
	"rent-n-go-backend/query"
)

type userRepository struct {
}

func (ur userRepository) GetByEmail(email string) (*models.User, error) {
	return query.User.Where(query.User.Email.Eq(email)).First()
}

func (ur userRepository) GetById(id uint) (*models.User, error) {
	return query.User.Where(query.User.ID.Eq(id)).First()
}

func (ur userRepository) GetByEmailOrPhone(email string, phone string) (*models.User, error) {
	u := query.User
	return u.Where(u.Email.Eq(email)).Or(u.PhoneNumber.Eq(phone)).First()
}

func (ur userRepository) Create(user *models.User) error {
	return query.User.Create(user)
}

func (ur userRepository) UpdateById(c *fiber.Ctx, userId uint, user *models.User) error {
	current, _ := query.User.Where(query.User.ID.Eq(userId)).First()

	if current.Email != user.Email {
		if _, err := query.User.Where(query.User.Email.Eq(user.Email)).First(); err == nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Ups, email already exist.",
			})
		}
	}

	if current.PhoneNumber != user.PhoneNumber {
		if _, err := query.User.Where(query.User.PhoneNumber.Eq(user.PhoneNumber)).First(); err == nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Ups, phone number already exist.",
			})
		}
	}

	query.User.Where(query.User.ID.Eq(userId)).Updates(user)

	return nil
}

func (ur userRepository) UpdatePasswordById(userId uint, payload *models.User) (gen.ResultInfo, error) {
	RefreshToken.DeleteByUserId(userId)
	return query.User.Where(query.User.ID.Eq(userId)).Updates(payload)
}
