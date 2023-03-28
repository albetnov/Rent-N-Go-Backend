package BasicRepositories

import (
	"gorm.io/gen"
	"os"
	"rent-n-go-backend/models"
	"rent-n-go-backend/query"
	"rent-n-go-backend/utils"
)

type picturesRepository struct {
}

func (f picturesRepository) Insert(associate string, associateId uint, fileName string) error {
	if err := checkAssociation(associate); err != nil {
		return err
	}

	if associate == "car" {
		if _, err := query.Cars.Where(query.Cars.ID.Eq(associateId)).First(); err != nil {
			return ErrNotFound
		}
	}

	if associate == "driver" {
		if _, err := query.Driver.Where(query.Driver.ID.Eq(associateId)).First(); err != nil {
			return ErrNotFound
		}
	}

	err := query.Pictures.Create(&models.Pictures{
		Associate:   associate,
		AssociateId: int(associateId),
		FileName:    utils.AssetPath(associate, fileName),
	})

	if err != nil {
		return err
	}

	return nil
}

func (f picturesRepository) DeleteById(id uint) (gen.ResultInfo, error) {
	qf := query.Pictures
	builder := qf.Where(qf.ID.Eq(id))

	if result, err := builder.First(); err != nil {
		os.Remove(result.FileName)
	}

	return builder.Delete()
}

func (f picturesRepository) DeleteByModuleId(associate string, associateId uint) (gen.ResultInfo, error) {
	if err := checkAssociation(associate); err != nil {
		return gen.ResultInfo{}, err
	}

	qf := query.Pictures
	builder := qf.Where(qf.Associate.Eq(associate)).Where(qf.AssociateId.Eq(int(associateId)))

	if result, err := builder.Find(); err != nil {
		for _, v := range result {
			os.Remove(v.FileName)
		}
	}

	return builder.Delete()
}

func (f picturesRepository) GetByModule(associate string, associateId uint) ([]*models.Pictures, error) {
	if err := checkAssociation(associate); err != nil {
		return nil, err
	}

	qf := query.Pictures

	return qf.Where(qf.Associate.Eq(associate)).Where(qf.AssociateId.Eq(int(associateId))).Find()
}
