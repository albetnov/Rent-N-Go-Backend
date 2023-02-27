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
	user := utils.GetUser(c)

	userId := uint(user["id"].(float64))

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

	auth := utils.GetUser(c)
	authId := uint(auth["id"].(float64))

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

	salt := uuid.New().String()

	c.SaveFile(file, path.Join(utils.PublicPath(), salt+file.Filename))

	auth := utils.GetUser(c)
	authId := uint(auth["id"].(float64))

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

func DeleteAccount(c *fiber.Ctx) error {
	auth := utils.GetUser(c)
	authId := uint(auth["id"].(float64))

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
