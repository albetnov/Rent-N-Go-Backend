package profile

import (
	"github.com/gofiber/fiber/v2"
	"rent-n-go-backend/models/UserModels"
	"rent-n-go-backend/repositories/UserRepositories"
	"rent-n-go-backend/utils"
	"strconv"
)

func getSim(sim UserModels.Sim) string {
	hasSim := "Not Uploaded"

	if sim.ID > 0 {
		if sim.IsVerified {
			hasSim = "Done!"
		} else {
			hasSim = "Not Verified"
		}
	}

	return hasSim
}

func getNik(nik UserModels.Nik) string {
	hasNik := "Not Filled"

	if nik.ID > 0 {
		if nik.IsVerified {
			hasNik = "Done!"
		} else {
			hasNik = "Not Verified"
		}
	}

	return hasNik
}

func CurrentUser(c *fiber.Ctx) error {
	currentUser := utils.GetUser(c)

	user, err := UserRepositories.User.GetAllById(uint(currentUser["id"].(float64)))

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Something went wrong, or user not found",
		})
	}

	if user.Photo.PhotoPath == "" {
		// set default photo
		user.Photo.PhotoPath = "https://source.unsplash.com/500x500?potrait"
	} else {
		user.Photo.PhotoPath = utils.FormatUrl(c, user.Photo.PhotoPath, "user")
	}

	currentUserStats := fiber.Map{
		"name":  user.Name,
		"phone": user.PhoneNumber,
		"photo": user.Photo,
		"email": user.Email,
		"role":  user.Role,
		"nik":   getNik(user.Nik),
		"sim":   getSim(user.Sim),
	}

	return c.JSON(fiber.Map{
		"data": currentUserStats,
	})
}

func CompletionStatus(c *fiber.Ctx) error {
	userId := utils.GetUserId(c)

	status := 0

	if data, err := UserRepositories.Nik.GetFromUserId(userId); err == nil {
		if data.IsVerified {
			status += 10
		}
		status += 40
	}

	if data, err := UserRepositories.Sim.GetByUserId(userId); err == nil {
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

	nikPayload := UserModels.Nik{
		Nik:        strconv.FormatInt(payload.Nik, 10),
		UserID:     authId,
		IsVerified: false,
	}

	UserRepositories.Nik.UpdateOrCreate(authId, &nikPayload)

	return c.JSON(fiber.Map{
		"message": "NIK updated successfully",
		"data":    nikPayload,
	})
}

func UpdateSim(c *fiber.Ctx) error {
	fileName, err := utils.SaveFileFromPayload(c, "file_name", "sim")

	if err != nil {
		return utils.SafeThrow(c, err)
	}

	authId := utils.GetUserId(c)

	simPayload := UserModels.Sim{
		UserID:     authId,
		IsVerified: false,
		FilePath:   utils.AssetPath("sim", fileName),
	}

	UserRepositories.Sim.UpdateOrCreate(authId, &simPayload)

	simPayload.FilePath = utils.FormatUrl(c, fileName, "sim")

	return c.JSON(fiber.Map{
		"message": "SIM updated successfully",
		"data":    simPayload,
	})
}

func UpdateProfile(c *fiber.Ctx) error {
	payload := utils.GetPayload[UpdateProfilePayload](c)

	authId := utils.GetUserId(c)

	updatePayload := UserModels.User{
		Name:        payload.Name,
		Email:       payload.Email,
		PhoneNumber: payload.PhoneNumber,
	}

	if err := UserRepositories.User.UpdateById(c, authId, &updatePayload); err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message": "Successfully update User Info.",
	})
}

func UpdatePassword(c *fiber.Ctx) error {
	payload := utils.GetPayload[UpdatePasswordPayload](c)
	authId := utils.GetUserId(c)

	currentUser, _ := UserRepositories.User.GetById(authId)

	if !utils.ComparePassword(payload.OldPassword, currentUser.Password) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Nah, wrong password bro",
		})
	}

	password, err := utils.HashPassword(payload.Password)

	if err != nil {
		utils.SafeThrow(c, err)
	}

	passwordPayload := UserModels.User{
		Password: password,
	}

	UserRepositories.User.UpdatePasswordById(authId, &passwordPayload)

	return c.JSON(fiber.Map{
		"message": "Password updated successfully",
		"action":  "LOGOUT",
	})
}

func DeleteAccount(c *fiber.Ctx) error {
	authId := utils.GetUserId(c)

	if err := UserRepositories.User.DeleteById(authId); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User not found, something is wrong...",
		})
	}

	// Yes even though the account has been removed in both storage and database, their JWT is still active
	// out there, and the JWT itself is not associated with database, therefore we just said "scheduled" :v
	// since it will expire anyway.
	return c.JSON(fiber.Map{
		"message": "Your account has been scheduled for deletion.",
	})
}

func UpdatePhoto(c *fiber.Ctx) error {
	fileName, err := utils.SaveFileFromPayload(c, "file_name", "user")

	if err != nil {
		return utils.SafeThrow(c, err)
	}

	authId := utils.GetUserId(c)

	UserRepositories.User.UpdateUserPhoto(authId, fileName)

	return c.JSON(fiber.Map{
		"message":   "Profile picture updated successfully",
		"file_name": fileName,
	})
}
