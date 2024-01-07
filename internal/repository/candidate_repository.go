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

type candidateRepository struct {
	db           *gorm.DB
	cacheManager cacher.CacheManager
}

func NewCandidateRepository(db *gorm.DB, cacheManager cacher.CacheManager) model.CandidateRepository {
	return &candidateRepository{
		db:           db,
		cacheManager: cacheManager,
	}
}

func (c *candidateRepository) FindByID(ctx context.Context, id int64) (*model.Candidate, error) {
	if id <= 0 {
		return nil, nil
	}
	logger := logrus.WithFields(logrus.Fields{
		"ctx": utils.DumpIncomingContext(ctx),
		"id":  id,
	})

	cacheKey := c.newCacheKeyByID(id)
	if !config.DisableCaching() {
		reply, mu, err := findFromCacheByKey[*model.Candidate](c.cacheManager, cacheKey)
		if err != nil {
			logger.Error(err)
			return nil, err
		}

		defer cacher.SafeUnlock(mu)

		if mu == nil {
			return reply, nil
		}
	}

	var candidate model.Candidate
	err := c.db.WithContext(ctx).Take(&candidate, "id = ?", id).Error
	switch err {
	case nil:
	case gorm.ErrRecordNotFound:
		storeNilCache(c.cacheManager, cacheKey)
		return nil, nil
	default:
		logger.Error(err)
		return nil, err
	}

	if err := c.cacheManager.StoreWithoutBlocking(cacher.NewItem(cacheKey, utils.Dump(candidate))); err != nil {
		logger.Error(err)
	}

	return &candidate, nil
}

func (c *candidateRepository) FindPasswordByID(ctx context.Context, id int64) ([]byte, error) {
	if id <= 0 {
		return nil, nil
	}
	logger := logrus.WithFields(logrus.Fields{
		"ctx": utils.DumpIncomingContext(ctx),
		"id":  id,
	})

	cacheKey := c.newPasswordCacheKeyByID(id)
	if !config.DisableCaching() {
		reply, mu, err := findFromCacheByKey[string](c.cacheManager, cacheKey)
		if err != nil {
			logger.Error(err)
			return nil, err
		}

		defer cacher.SafeUnlock(mu)

		if mu == nil {
			return []byte(reply), nil
		}
	}

	var pass string
	err := c.db.WithContext(ctx).Model(model.Candidate{}).Select("password").Take(&pass, "id = ?", id).Error
	switch err {
	case nil:
	case gorm.ErrRecordNotFound:
		storeNilCache(c.cacheManager, cacheKey)
		return nil, nil
	default:
		logger.Error(err)
		return nil, err
	}

	if err := c.cacheManager.StoreWithoutBlocking(cacher.NewItem(cacheKey, pass)); err != nil {
		logger.Error(err)
	}

	return []byte(pass), err
}

func (c *candidateRepository) FindByEmail(ctx context.Context, email string) (*model.Candidate, error) {
	if email == "" {
		return nil, nil
	}

	logger := logrus.WithFields(logrus.Fields{
		"ctx":   utils.DumpIncomingContext(ctx),
		"email": email,
	})

	cacheKey := c.newCacheKeyByEmail(email)
	if !config.DisableCaching() {
		id, mu, err := findFromCacheByKey[int64](c.cacheManager, cacheKey)
		if err != nil {
			logger.Error(err)
			return nil, err
		}

		defer cacher.SafeUnlock(mu)

		if mu == nil {
			return c.FindByID(ctx, id)
		}
	}

	var id int64
	err := c.db.WithContext(ctx).Model(model.Candidate{}).Select("id").Take(&id, "email = ?", email).Error
	switch err {
	case nil:
	case gorm.ErrRecordNotFound:
		storeNilCache(c.cacheManager, cacheKey)
		return nil, nil
	default:
		logger.Error(err)
		return nil, err
	}

	if err := c.cacheManager.StoreWithoutBlocking(cacher.NewItem(cacheKey, id)); err != nil {
		logger.Error(err)
	}

	return c.FindByID(ctx, id)
}

func (c *candidateRepository) FindByPhone(ctx context.Context, phone string) (*model.Candidate, error) {
	if phone == "" {
		return nil, nil
	}

	logger := logrus.WithFields(logrus.Fields{
		"ctx":   utils.DumpIncomingContext(ctx),
		"phone": phone,
	})

	cacheKey := c.newCacheKeyByPhone(phone)
	if !config.DisableCaching() {
		id, mu, err := findFromCacheByKey[int64](c.cacheManager, cacheKey)
		if err != nil {
			logger.Error(err)
			return nil, err
		}

		defer cacher.SafeUnlock(mu)

		if mu == nil {
			return c.FindByID(ctx, id)
		}
	}

	var id int64
	err := c.db.WithContext(ctx).Model(model.Candidate{}).Select("id").Take(&id, "phone = ?", phone).Error
	switch err {
	case nil:
	case gorm.ErrRecordNotFound:
		storeNilCache(c.cacheManager, cacheKey)
		return nil, nil
	default:
		logger.Error(err)
		return nil, err
	}

	if err := c.cacheManager.StoreWithoutBlocking(cacher.NewItem(cacheKey, id)); err != nil {
		logger.Error(err)
	}

	return c.FindByID(ctx, id)
}

func (c *candidateRepository) FindUnscopedByEmail(ctx context.Context, email string) (*model.Candidate, error) {
	if email == "" {
		return nil, nil
	}

	logger := logrus.WithFields(logrus.Fields{
		"ctx":   utils.DumpIncomingContext(ctx),
		"email": email,
	})

	var id int64
	err := c.db.WithContext(ctx).Model(model.Candidate{}).Unscoped().Select("id").Take(&id, "email = ?", email).Error
	switch err {
	case nil:
	case gorm.ErrRecordNotFound:
		return nil, nil
	default:
		logger.Error(err)
		return nil, err
	}

	return c.FindByID(ctx, id)
}

func (c *candidateRepository) FindUnscopedByPhone(ctx context.Context, phone string) (*model.Candidate, error) {
	if phone == "" {
		return nil, nil
	}

	logger := logrus.WithFields(logrus.Fields{
		"ctx":   utils.DumpIncomingContext(ctx),
		"phone": phone,
	})

	var id int64
	err := c.db.WithContext(ctx).Model(model.Candidate{}).Unscoped().Select("id").Take(&id, "phone = ?", phone).Error
	switch err {
	case nil:
	case gorm.ErrRecordNotFound:
		return nil, nil
	default:
		logger.Error(err)
		return nil, err
	}

	return c.FindByID(ctx, id)
}

func (c *candidateRepository) Create(ctx context.Context, candidate *model.Candidate) error {
	logger := logrus.WithContext(ctx).WithFields(logrus.Fields{
		"ctx":       utils.DumpIncomingContext(ctx),
		"candidate": utils.Dump(candidate),
	})

	if err := c.db.WithContext(ctx).Create(candidate).Error; err != nil {
		logger.Error(err)
		return err
	}

	if err := c.deleteCommonCache(candidate); err != nil {
		logger.Error(err)
	}

	return nil
}

func (c *candidateRepository) Update(ctx context.Context, candidate *model.Candidate) error {
	logger := logrus.WithContext(ctx).WithFields(logrus.Fields{
		"ctx":       utils.DumpIncomingContext(ctx),
		"candidate": utils.Dump(candidate),
	})

	if err := c.db.WithContext(ctx).Model(model.Candidate{}).
		Where("id = ?", candidate.ID).Updates(candidate).Error; err != nil {
		logger.Error(err)
		return err
	}

	if err := c.deleteCommonCache(candidate); err != nil {
		logger.Error(err)
	}

	return nil
}

func (c *candidateRepository) deleteCommonCache(candidate *model.Candidate) error {
	cacheKeys := []string{
		c.newCacheKeyByID(candidate.ID),
		c.newCacheKeyByEmail(candidate.Email.String),
		c.newCacheKeyByPhone(candidate.Phone.String),
	}

	return c.cacheManager.DeleteByKeys(cacheKeys)
}

func (c *candidateRepository) newCacheKeyByID(id int64) string {
	return fmt.Sprintf("cache:object:candidate:id:%d", id)
}

func (c *candidateRepository) newCacheKeyByEmail(email string) string {
	return fmt.Sprintf("cache:id:candidate:email:%s", email)
}

func (c *candidateRepository) newCacheKeyByPhone(phone string) string {
	return fmt.Sprintf("cache:id:candidate:phone:%s", phone)
}

func (c *candidateRepository) newPasswordCacheKeyByID(id int64) string {
	return fmt.Sprintf("cache:password:id:%d", id)
}
