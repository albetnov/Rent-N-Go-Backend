package UserRepositories

import (
	"github.com/gofiber/fiber/v2"
	"os"
	"rent-n-go-backend/models/UserModels"
	"rent-n-go-backend/query"
	"rent-n-go-backend/utils"
	"strings"
)

type simRepository struct {
}

func (s simRepository) GetByUserId(userId uint) (*UserModels.Sim, error) {
	return query.Sim.Where(query.Sim.UserID.Eq(userId)).First()
}

func (sr simRepository) UpdateOrCreate(userId uint, payload *UserModels.Sim) {
	s := query.Sim

	preCond := s.Where(s.UserID.Eq(userId))

	if result, err := preCond.First(); err == nil {
		os.Remove(utils.AssetPath("sim", result.FilePath))
		preCond.Updates(payload)
	} else {
		s.Create(payload)
	}
}

func (s simRepository) OptionalCreate(c *fiber.Ctx, payload string, sess utils.SessionStore, userId uint, fallback string) error {
	simFile, err := utils.SaveFileFromPayload(c, payload, utils.AssetPath("sim"))

	if err != nil {
		if strings.Contains(err.Error(), utils.NO_UPLOADED_FILE) {
			return nil
		}

		sess.SetSession("error", utils.GetErrorMessage(err))
		return c.RedirectBack(fallback)
	}

	Sim.UpdateOrCreate(userId, &UserModels.Sim{
		UserID:     userId,
		IsVerified: false,
		FilePath:   simFile,
	})

	return nil
}
