package UserRepositories

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gen"
	"gorm.io/gen/field"
	"os"
	"path"
	"rent-n-go-backend/models/UserModels"
	"rent-n-go-backend/query"
	"rent-n-go-backend/utils"
)

type userRepository struct {
}

func (ur userRepository) GetByEmail(email string) (*UserModels.User, error) {
	return query.User.Where(query.User.Email.Eq(email)).First()
}

func (ur userRepository) GetById(id uint) (*UserModels.User, error) {
	return query.User.Preload(field.Associations).Where(query.User.ID.Eq(id)).First()
}

func (ur userRepository) GetByEmailOrPhone(email string, phone string) (*UserModels.User, error) {
	u := query.User
	return u.Where(u.Email.Eq(email)).Or(u.PhoneNumber.Eq(phone)).First()
}

func (ur userRepository) Create(user *UserModels.User) error {
	return query.User.Create(user)
}

func (ur userRepository) UpdateById(c *fiber.Ctx, userId uint, user *UserModels.User) error {
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

func (ur userRepository) UpdatePasswordById(userId uint, payload *UserModels.User) (gen.ResultInfo, error) {
	RefreshToken.DeleteByUserId(userId)
	return query.User.Where(query.User.ID.Eq(userId)).Updates(payload)
}

func (ur userRepository) UpdateUserPhoto(userId uint, fileName string) {
	qup := query.UserPhoto

	preCond := qup.Where(qup.UserID.Eq(userId))

	if result, err := preCond.First(); err == nil {
		os.Remove(path.Join(utils.AssetPath("user"), result.PhotoPath))
		preCond.Update(qup.PhotoPath, fileName)
		return
	}

	qup.Create(&UserModels.UserPhoto{PhotoPath: fileName, UserID: userId})
}

func (ur userRepository) GetAllById(userId uint) (*UserModels.User, error) {
	qu := query.User
	return qu.Where(qu.ID.Eq(userId)).Preload(field.Associations).First()
}
