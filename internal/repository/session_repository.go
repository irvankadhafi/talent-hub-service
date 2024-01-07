package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/irvankadhafi/talent-hub-service/internal/config"
	"github.com/irvankadhafi/talent-hub-service/internal/model"
	"github.com/irvankadhafi/talent-hub-service/pkg/cacher"
	"github.com/irvankadhafi/talent-hub-service/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

type sessionRepo struct {
	db           *gorm.DB
	cacheManager cacher.CacheManager
}

// NewSessionRepository sessionRepo constructor
func NewSessionRepository(
	db *gorm.DB,
	cacheManager cacher.CacheManager,
) model.SessionRepository {
	return &sessionRepo{
		db:           db,
		cacheManager: cacheManager,
	}
}

func (s *sessionRepo) Create(ctx context.Context, sess *model.Session) error {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":    utils.DumpIncomingContext(ctx),
		"userID": sess.CandidateID,
	})

	if err := s.db.WithContext(ctx).Create(sess).Error; err != nil {
		logger.Error(err)
		return err
	}

	if err := s.cacheToken(sess); err != nil {
		logger.Error(err)
	}

	return nil
}

// FindByToken find a session by it's token
func (s *sessionRepo) FindByToken(ctx context.Context, tokenType model.TokenType, token string) (*model.Session, error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":       utils.DumpIncomingContext(ctx),
		"tokenType": tokenType,
	})

	cacheKey := model.NewSessionTokenCacheKey(token)
	if !config.DisableCaching() {
		reply, mu, err := findFromCacheByKey[*model.Session](s.cacheManager, cacheKey)
		if err != nil {
			logger.Error(err)
			return nil, err
		}

		defer cacher.SafeUnlock(mu)

		if mu == nil {
			return reply, nil
		}
	}

	sess := &model.Session{}
	var err error
	switch tokenType {
	case model.AccessToken:
		err = s.db.Take(sess, "access_token = ?", token).Error
	case model.RefreshToken:
		err = s.db.Take(sess, "refresh_token = ?", token).Error
	}
	switch err {
	case nil:
	case gorm.ErrRecordNotFound:
		storeNilCache(s.cacheManager, cacheKey)
		return nil, nil
	default:
		logger.Error(err)
		return nil, err
	}

	if err = s.cacheToken(sess); err != nil {
		logger.Error(err)
	}

	return sess, nil
}

// FindByID find session by id
func (s *sessionRepo) FindByID(ctx context.Context, id int64) (*model.Session, error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx": utils.DumpIncomingContext(ctx),
		"id":  id,
	})

	cacheKey := s.newCacheKeyByID(id)
	if !config.DisableCaching() {
		sess, mu, err := findFromCacheByKey[*model.Session](s.cacheManager, cacheKey)
		if err != nil {
			return nil, err
		}
		defer cacher.SafeUnlock(mu)

		if mu == nil {
			return sess, nil
		}
	}

	sess := model.Session{}
	err := s.db.WithContext(ctx).Take(&sess, "id = ?", id).Error
	switch err {
	case nil:
	case gorm.ErrRecordNotFound:
		storeNilCache(s.cacheManager, cacheKey)
		return nil, nil
	default:
		logger.Error(err)
		return nil, err
	}

	if err = s.cacheToken(&sess); err != nil {
		logger.Error(err)
	}

	return &sess, nil
}

// CheckToken check whether the token exists or not in the cache
func (s *sessionRepo) CheckToken(ctx context.Context, token string) (exist bool, err error) {
	reply, err := s.cacheManager.Get(model.NewSessionTokenCacheKey(token))
	if err != nil {
		return false, err
	}

	bt, _ := reply.([]byte)
	return string(bt) != "", nil
}

// RefreshToken update access and refresh token string value and expired_at
func (s *sessionRepo) RefreshToken(ctx context.Context, oldSess, sess *model.Session) (*model.Session, error) {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":     utils.DumpIncomingContext(ctx),
		"session": utils.Dump(sess),
	})

	sess.UpdatedAt = time.Now()
	err := s.db.WithContext(ctx).Model(model.Session{}).Select(
		"access_token",
		"refresh_token",
		"access_token_expired_at",
		"refresh_token_expired_at",
		"user_agent",
		"ip_address",
		"updated_at",
	).Where("id = ?", sess.ID).Updates(sess).Error
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	if err = s.deleteCaches(oldSess); err != nil {
		logger.Error(err)
	}
	if err = s.deleteCaches(sess); err != nil {
		logger.Error(err)
	}

	return s.FindByID(ctx, sess.ID)
}

// DeleteByCandidateIDAndMaxRemainderSession delete session by candidate id
func (s *sessionRepo) DeleteByCandidateIDAndMaxRemainderSession(ctx context.Context, candidateID int64, maxRemainderSess int) error {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":              utils.DumpIncomingContext(ctx),
		"candidateID":      candidateID,
		"maxRemainderSess": maxRemainderSess,
	})

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		deleteIDs, cacheKeys, err := s.getOffsetIDsAndCacheKeysByCandidateIDAndMaxActiveSess(ctx, tx, candidateID, maxRemainderSess)
		if err != nil {
			logger.Error(err)
			return err
		}

		if len(deleteIDs) == 0 {
			return nil
		}

		if err := tx.Delete(&model.Session{}, deleteIDs).Error; err != nil {
			logger.Error(err)
			return err
		}

		if err := s.cacheManager.DeleteByKeys(cacheKeys); err != nil {
			logger.Error(err)
			return err
		}

		return nil
	})
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

// Delete deletes existing session by id.
func (s *sessionRepo) Delete(ctx context.Context, session *model.Session) error {
	logger := logrus.WithFields(logrus.Fields{
		"ctx":     utils.DumpIncomingContext(ctx),
		"session": utils.Dump(session),
	})

	if err := s.db.WithContext(ctx).Delete(session).Error; err != nil {
		logger.Error(err)
		return err
	}

	if err := s.deleteCaches(session); err != nil {
		logger.Error(err)
	}

	return nil
}

func (s *sessionRepo) getOffsetIDsAndCacheKeysByCandidateIDAndMaxActiveSess(ctx context.Context, tx *gorm.DB, candidateID int64, maxRemainderSess int) ([]int64, []string, error) {
	var (
		deleteIDs []int64
		cacheKeys []string
		sessions  = []model.Session{}
		limit     = config.SessionDeleteBatchSize()
		offset    = maxRemainderSess
	)

	logger := logrus.WithFields(logrus.Fields{
		"ctx": utils.DumpIncomingContext(ctx),
		"DeleteByCandidateIDAndMaxRemainderSession": candidateID,
		"maxRemainderSess":                          maxRemainderSess,
	})

	for {
		err := tx.WithContext(ctx).
			Where(&model.Session{CandidateID: candidateID}).
			Order("refresh_token_expired_at desc").
			Offset(offset).Limit(limit).
			Find(&sessions).Error
		if err != nil {
			logger.Error(err)
			return deleteIDs, cacheKeys, err
		}

		if len(sessions) == 0 {
			break
		}

		for _, session := range sessions {
			deleteIDs = append(deleteIDs, session.ID)
			cacheKeys = append(cacheKeys,
				model.NewSessionTokenCacheKey(session.AccessToken),
				model.NewSessionTokenCacheKey(session.RefreshToken),
				s.newCacheKeyByID(session.ID),
			)
		}

		offset += limit
		sessions = nil
	}

	return deleteIDs, cacheKeys, nil
}

func (s *sessionRepo) cacheToken(session *model.Session) error {
	sess, err := json.Marshal(session)
	if err != nil {
		return err
	}

	now := time.Now()
	return s.cacheManager.StoreMultiWithoutBlocking([]cacher.Item{
		cacher.NewItemWithCustomTTL(model.NewSessionTokenCacheKey(session.AccessToken), sess, session.AccessTokenExpiredAt.Sub(now)),
		cacher.NewItemWithCustomTTL(s.newCacheKeyByID(session.ID), sess, session.AccessTokenExpiredAt.Sub(now)),
		cacher.NewItemWithCustomTTL(model.NewSessionTokenCacheKey(session.RefreshToken), sess, session.RefreshTokenExpiredAt.Sub(now)),
	})
}

func (s *sessionRepo) deleteCaches(session *model.Session) error {
	return s.cacheManager.DeleteByKeys([]string{
		model.NewSessionTokenCacheKey(session.AccessToken),
		model.NewSessionTokenCacheKey(session.RefreshToken),
		s.newCacheKeyByID(session.ID),
	})
}

func (s *sessionRepo) newCacheKeyByID(id int64) string {
	return fmt.Sprintf("cache:object:session:id:%d", id)
}
