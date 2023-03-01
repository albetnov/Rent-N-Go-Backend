package user

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gen"
	"gorm.io/gen/field"
	"os"
	"path"
	"rent-n-go-backend/models/user"
	"rent-n-go-backend/query"
	"rent-n-go-backend/utils"
)

type userRepository struct {
}

func (ur userRepository) GetByEmail(email string) (*user.User, error) {
	return query.User.Where(query.User.Email.Eq(email)).First()
}

func (ur userRepository) GetById(id uint) (*user.User, error) {
	return query.User.Preload(field.Associations).Where(query.User.ID.Eq(id)).First()
}

func (ur userRepository) GetByEmailOrPhone(email string, phone string) (*user.User, error) {
	u := query.User
	return u.Where(u.Email.Eq(email)).Or(u.PhoneNumber.Eq(phone)).First()
}

func (ur userRepository) Create(user *user.User) error {
	return query.User.Create(user)
}

func (ur userRepository) UpdateById(c *fiber.Ctx, userId uint, user *user.User) error {
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

func (ur userRepository) UpdatePasswordById(userId uint, payload *user.User) (gen.ResultInfo, error) {
	RefreshToken.DeleteByUserId(userId)
	return query.User.Where(query.User.ID.Eq(userId)).Updates(payload)
}

func (ur userRepository) UpdateUserPhoto(userId uint, fileName string) {
	qup := query.UserPhoto

	preCond := qup.Where(qup.UserID.Eq(userId))

	if result, err := preCond.First(); err == nil {
		os.Remove(path.Join(utils.PublicPath(), result.PhotoPath))
		preCond.Update(qup.PhotoPath, fileName)
		return
	}

	qup.Create(&user.Photo{PhotoPath: fileName, UserID: userId})
}
