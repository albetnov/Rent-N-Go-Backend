package profile

import (
	"github.com/gofiber/fiber/v2"
	"os"
	"path"
	"rent-n-go-backend/models/user"
	"rent-n-go-backend/query"
	userRepository "rent-n-go-backend/repositories/user"
	"rent-n-go-backend/utils"
	"strconv"
)

func CurrentUser(c *fiber.Ctx) error {
	currentUser := utils.GetUser(c)

	return c.JSON(fiber.Map{
		"data": currentUser,
	})
}

func CompletionStatus(c *fiber.Ctx) error {
	userId := utils.GetUserId(c)

	status := 0

	if data, err := userRepository.Nik.GetFromUserId(userId); err == nil {
		if data.IsVerified {
			status += 10
		}
		status += 40
	}

	if data, err := userRepository.Sim.GetByUserId(userId); err == nil {
		if data.IsVerified {
			status += 10
		}

		status += 40
	}

	return c.JSON(fiber.Map{
		"message":    "Ok!",
		"percentage": status,
	})
}

func UpdateNik(c *fiber.Ctx) error {
	payload := utils.GetPayload[CompleteNikPayload](c)

	authId := utils.GetUserId(c)

	nikPayload := user.Nik{
		Nik:        strconv.FormatInt(payload.Nik, 10),
		UserID:     authId,
		IsVerified: false,
	}

	userRepository.Nik.UpdateOrCreate(authId, &nikPayload)

	return c.JSON(fiber.Map{
		"message": "NIK updated successfully",
		"data":    nikPayload,
	})
}

func UpdateSim(c *fiber.Ctx) error {
	fileName, err := utils.SaveFileFromPayload(c, "file_name")

	if err != nil {
		return utils.SafeThrow(c, err)
	}

	authId := utils.GetUserId(c)

	simPayload := user.Sim{
		UserID:     authId,
		IsVerified: false,
		FilePath:   fileName,
	}

	userRepository.Sim.UpdateOrCreate(authId, &simPayload)

	return c.JSON(fiber.Map{
		"message": "SIM updated successfully",
		"data":    simPayload,
	})
}

func UpdateProfile(c *fiber.Ctx) error {
	payload := utils.GetPayload[UpdateProfilePayload](c)

	authId := utils.GetUserId(c)

	updatePayload := user.User{
		Name:        payload.Name,
		Email:       payload.Email,
		PhoneNumber: payload.PhoneNumber,
	}

	if err := userRepository.User.UpdateById(c, authId, &updatePayload); err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message": "Successfully update user info.",
		"action":  "REFRESH_TOKEN",
	})
}

func UpdatePassword(c *fiber.Ctx) error {
	payload := utils.GetPayload[UpdatePasswordPayload](c)
	authId := utils.GetUserId(c)

	currentUser, _ := userRepository.User.GetById(authId)

	if !utils.ComparePassword(payload.OldPassword, currentUser.Password) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Nah, wrong password bro",
		})
	}

	password, err := utils.HashPassword(payload.Password)

	if err != nil {
		utils.SafeThrow(c, err)
	}

	passwordPayload := user.User{
		Password: password,
	}

	userRepository.User.UpdatePasswordById(authId, &passwordPayload)

	return c.JSON(fiber.Map{
		"message": "Password updated successfully",
		"action":  "LOGOUT",
	})
}

func DeleteAccount(c *fiber.Ctx) error {
	authId := utils.GetUserId(c)

	u := query.User

	user, _ := userRepository.User.GetById(authId)

	u.Select(u.Nik.Field()).Delete(user)

	if user.Sim.FilePath != "" {
		os.Remove(path.Join(utils.PublicPath(), user.Sim.FilePath))
	}

	u.Select(u.Sim.Field()).Delete(user)

	u.Where(u.ID.Eq(user.ID)).Delete()

	// Yes even though the account has been removed in both storage and database, their JWT is still active
	// out there, and the JWT itself is not associated with database, therefore we just said "scheduled" :v
	// since it will expire anyway.
	return c.JSON(fiber.Map{
		"message": "Your account has been scheduled for deletion.",
	})
}

func UpdatePhoto(c *fiber.Ctx) error {
	fileName, err := utils.SaveFileFromPayload(c, "file_name")

	if err != nil {
		return utils.SafeThrow(c, err)
	}

	authId := utils.GetUserId(c)

	userRepository.User.UpdateUserPhoto(authId, fileName)

	return c.JSON(fiber.Map{
		"message":   "Profile picture updated successfully",
		"file_name": fileName,
	})
}
