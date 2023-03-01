package profile

import (
	"github.com/gofiber/fiber/v2"
	"os"
	"path"
	"rent-n-go-backend/models"
	"rent-n-go-backend/query"
	"rent-n-go-backend/repositories"
	"rent-n-go-backend/utils"
	"strconv"
)

func CurrentUser(c *fiber.Ctx) error {
	user := utils.GetUser(c)

	return c.JSON(fiber.Map{
		"data": user,
	})
}

func CompletionStatus(c *fiber.Ctx) error {
	userId := utils.GetUserId(c)

	status := 0

	if data, err := repositories.Nik.GetFromUserId(userId); err == nil {
		if data.IsVerified {
			status += 10
		}
		status += 40
	}

	if data, err := repositories.Sim.GetByUserId(userId); err == nil {
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

	nikPayload := models.Nik{
		Nik:        strconv.FormatInt(payload.Nik, 10),
		UserID:     authId,
		IsVerified: false,
	}

	repositories.Nik.UpdateOrCreate(authId, &nikPayload)

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

	simPayload := models.Sim{
		UserID:     authId,
		IsVerified: false,
		FilePath:   fileName,
	}

	repositories.Sim.UpdateOrCreate(authId, &simPayload)

	return c.JSON(fiber.Map{
		"message": "SIM updated successfully",
		"data":    simPayload,
	})
}

func UpdateProfile(c *fiber.Ctx) error {
	payload := utils.GetPayload[UpdateProfilePayload](c)

	authId := utils.GetUserId(c)

	updatePayload := models.User{
		Name:        payload.Name,
		Email:       payload.Email,
		PhoneNumber: payload.PhoneNumber,
	}

	if err := repositories.User.UpdateById(c, authId, &updatePayload); err != nil {
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

	user, _ := repositories.User.GetById(authId)

	if !utils.ComparePassword(payload.OldPassword, user.Password) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Nah, wrong password bro",
		})
	}

	password, err := utils.HashPassword(payload.Password)

	if err != nil {
		utils.SafeThrow(c, err)
	}

	passwordPayload := models.User{
		Password: password,
	}

	repositories.User.UpdatePasswordById(authId, &passwordPayload)

	return c.JSON(fiber.Map{
		"message": "Password updated successfully",
		"action":  "LOGOUT",
	})
}

func DeleteAccount(c *fiber.Ctx) error {
	authId := utils.GetUserId(c)

	u := query.User

	user, _ := repositories.User.GetById(authId)

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

	repositories.User.UpdateUserPhoto(authId, fileName)

	return c.JSON(fiber.Map{
		"message":   "Profile picture updated successfully",
		"file_name": fileName,
	})
}
