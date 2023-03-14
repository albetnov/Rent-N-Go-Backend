package BasicRepositories

import (
	"gorm.io/gen"
	"rent-n-go-backend/models"
	"rent-n-go-backend/query"
)

type featuresRepository struct {
}

func (f featuresRepository) Insert(associate string, associateId uint, iconKey string, value string) error {
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

	err := query.Features.Create(&models.Features{
		Associate:   associate,
		AssociateId: int(associateId),
		IconKey:     iconKey,
		Value:       value,
	})

	if err != nil {
		return err
	}

	return nil
}

func (f featuresRepository) DeleteById(id uint) (gen.ResultInfo, error) {
	qf := query.Features
	return qf.Where(qf.ID.Eq(id)).Delete()
}

func (f featuresRepository) DeleteByModuleId(associate string, associateId uint) (gen.ResultInfo, error) {
	if err := checkAssociation(associate); err != nil {
		return gen.ResultInfo{}, err
	}

	qf := query.Features
	return qf.Where(qf.Associate.Eq(associate)).Where(qf.AssociateId.Eq(int(associateId))).Delete()
}

func (f featuresRepository) GetByModule(associate string, associateId uint) ([]*models.Features, error) {
	if err := checkAssociation(associate); err != nil {
		return nil, err
	}

	qf := query.Features

	return qf.Where(qf.Associate.Eq(associate)).Where(qf.AssociateId.Eq(int(associateId))).Find()
}
