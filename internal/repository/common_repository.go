package repository

import (
	"encoding/json"
	"github.com/go-redsync/redsync/v4"
	"github.com/irvankadhafi/talent-hub-service/pkg/cacher"
	"github.com/sirupsen/logrus"
)

func findFromCacheByKey[T any](cacheKeeper cacher.CacheManager, key string) (item T, mutex *redsync.Mutex, err error) {
	var cachedData any

	cachedData, mutex, err = cacheKeeper.GetOrLock(key)
	if err != nil || cachedData == nil {
		return
	}

	cachedDataByte, _ := cachedData.([]byte)
	if cachedDataByte == nil {
		return
	}

	if err = json.Unmarshal(cachedDataByte, &item); err != nil {
		return
	}

	return
}

func storeNilCache(cache cacher.CacheManager, cacheKey string) {
	if err := cache.StoreNil(cacheKey); err != nil {
		logrus.Error(err)
	}
}
