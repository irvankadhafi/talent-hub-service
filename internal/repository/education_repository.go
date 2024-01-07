package repository

import (
	"context"
	"fmt"
	"github.com/irvankadhafi/talent-hub-service/internal/config"
	"github.com/irvankadhafi/talent-hub-service/internal/model"
	"github.com/irvankadhafi/talent-hub-service/pkg/cacher"
	"github.com/irvankadhafi/talent-hub-service/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type educationRepository struct {
	db           *gorm.DB
	cacheManager cacher.CacheManager
}

func NewEducationRepository(
	db *gorm.DB,
	cacheManager cacher.CacheManager,
) model.EducationRepository {
	return &educationRepository{
		db:           db,
		cacheManager: cacheManager,
	}
}

func (e *educationRepository) FindByID(ctx context.Context, id int64) (*model.Education, error) {
	if id <= 0 {
		return nil, nil
	}
	logger := logrus.WithFields(logrus.Fields{
		"ctx": utils.DumpIncomingContext(ctx),
		"id":  id,
	})

	cacheKey := e.newCacheKeyByID(id)
	if !config.DisableCaching() {
		reply, mu, err := findFromCacheByKey[*model.Education](e.cacheManager, cacheKey)
		if err != nil {
			logger.Error(err)
			return nil, err
		}

		defer cacher.SafeUnlock(mu)

		if mu == nil {
			return reply, nil
		}
	}

	var education model.Education
	err := e.db.WithContext(ctx).Take(&education, "id = ?", id).Error
	switch err {
	case nil:
	case gorm.ErrRecordNotFound:
		storeNilCache(e.cacheManager, cacheKey)
		return nil, nil
	default:
		logger.Error(err)
		return nil, err
	}

	if err := e.cacheManager.StoreWithoutBlocking(cacher.NewItem(cacheKey, utils.Dump(education))); err != nil {
		logger.Error(err)
	}

	return &education, nil
}

func (e *educationRepository) FindByCandidateID(ctx context.Context, candidateID int64) (*model.Education, error) {
	if candidateID <= 0 {
		return nil, nil
	}

	logger := logrus.WithFields(logrus.Fields{
		"ctx":         utils.DumpIncomingContext(ctx),
		"candidateID": candidateID,
	})

	cacheKey := e.newCacheKeyByCandidateID(candidateID)
	if !config.DisableCaching() {
		id, mu, err := findFromCacheByKey[int64](e.cacheManager, cacheKey)
		if err != nil {
			logger.Error(err)
			return nil, err
		}

		defer cacher.SafeUnlock(mu)

		if mu == nil {
			return e.FindByID(ctx, id)
		}
	}

	var id int64
	err := e.db.WithContext(ctx).Model(model.Education{}).Select("id").Take(&id, "candidate_id = ?", candidateID).Error
	switch err {
	case nil:
	case gorm.ErrRecordNotFound:
		storeNilCache(e.cacheManager, cacheKey)
		return nil, nil
	default:
		logger.Error(err)
		return nil, err
	}

	if err := e.cacheManager.StoreWithoutBlocking(cacher.NewItem(cacheKey, id)); err != nil {
		logger.Error(err)
	}

	return e.FindByID(ctx, id)
}

func (e *educationRepository) Create(ctx context.Context, education *model.Education) error {
	logger := logrus.WithContext(ctx).WithFields(logrus.Fields{
		"ctx":       utils.DumpIncomingContext(ctx),
		"education": utils.Dump(education),
	})

	if err := e.db.WithContext(ctx).Create(education).Error; err != nil {
		logger.Error(err)
		return err
	}

	if err := e.deleteCommonCache(education); err != nil {
		logger.Error(err)
	}

	return nil
}

func (e *educationRepository) Update(ctx context.Context, education *model.Education) error {
	logger := logrus.WithContext(ctx).WithFields(logrus.Fields{
		"ctx":       utils.DumpIncomingContext(ctx),
		"education": utils.Dump(education),
	})

	if err := e.db.WithContext(ctx).Model(model.Education{}).
		Where("id = ?", education.ID).Updates(education).Error; err != nil {
		logger.Error(err)
		return err
	}

	if err := e.deleteCommonCache(education); err != nil {
		logger.Error(err)
	}

	return nil
}

func (e *educationRepository) newCacheKeyByID(id int64) string {
	return fmt.Sprintf("cache:object:education:id:%d", id)
}

func (e *educationRepository) newCacheKeyByCandidateID(candidateID int64) string {
	return fmt.Sprintf("cache:id:education:candidate_id:%d", candidateID)
}

func (e *educationRepository) deleteCommonCache(education *model.Education) error {
	cacheKeys := []string{
		e.newCacheKeyByID(education.ID),
		e.newCacheKeyByCandidateID(education.CandidateID),
	}

	return e.cacheManager.DeleteByKeys(cacheKeys)
}
