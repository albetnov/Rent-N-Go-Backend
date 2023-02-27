package profile

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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

	if data, err := repositories.GetNikFromUserId(userId); err == nil {
		if data.IsVerified {
			status += 10
		}
		status += 40
	}

	if data, err := repositories.GetSimByUserId(userId); err == nil {
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

	repositories.UpdateOrCreateNik(authId, &nikPayload)

	return c.JSON(fiber.Map{
		"message": "NIK updated successfully",
		"data":    nikPayload,
	})
}

func UpdateSim(c *fiber.Ctx) error {
	file, err := c.FormFile("file_name")

	if err != nil {
		return utils.SafeThrow(c, err)
	}

	reader, err := file.Open()
	if err != nil {
		utils.SafeThrow(c, err)
	}

	defer reader.Close()

	err = utils.CheckMimes(c, reader, []string{"image/jpg", "image/png", "image/jpeg"})

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	salt := uuid.New().String()

	c.SaveFile(file, path.Join(utils.PublicPath(), salt+file.Filename))

	authId := utils.GetUserId(c)

	simPayload := models.Sim{
		UserID:     authId,
		IsVerified: false,
		FilePath:   salt + file.Filename,
	}

	repositories.UpdateOrCrateSim(authId, &simPayload)

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

	if err := repositories.UpdateUserByUserId(c, authId, &updatePayload); err != nil {
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

	user, _ := repositories.GetUserById(authId)

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

	repositories.UpdateUserPasswordByUserId(authId, &passwordPayload)

	return c.JSON(fiber.Map{
		"message": "Password updated successfully",
		"action":  "LOGOUT",
	})
}

func DeleteAccount(c *fiber.Ctx) error {
	authId := utils.GetUserId(c)

	u := query.User

	user, _ := repositories.GetUserById(authId)

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
