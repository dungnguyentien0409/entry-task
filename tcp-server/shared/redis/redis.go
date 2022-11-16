package cache

import (
	"context"
	"entrytask/tcp-server/shared/dto"
	"github.com/bsm/redislock"
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"time"
)

type Cache struct {
	redisCache *cache.Cache
	ctx        context.Context
}

func NewCache(addr string) *Cache {
	redisClient := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	cache := cache.New(&cache.Options{
		Redis:      redisClient,
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})

	locker := redislock.New(redisClient)

	locker.Obtain(context.Background(), "cart-key", time.Second, nil)

	return &Cache{
		redisCache: cache,
		ctx:        context.TODO(),
	}
}

func (userCache *Cache) Save(key string, user dto.User) error {
	err := userCache.redisCache.Set(&cache.Item{
		Ctx:   userCache.ctx,
		Key:   key,
		Value: user,
		TTL:   time.Hour * 24,
	})

	return err
}

func (userCache *Cache) Get(key string) (*dto.User, error) {
	var user dto.User
	err := userCache.redisCache.Get(userCache.ctx, key, &user)
	if err != nil {
		return &user, err
	}
	return &user, nil
}

func (userCache *Cache) SaveWithExpire(key string, user dto.User, duration time.Duration) error {
	err := userCache.redisCache.Set(&cache.Item{
		Ctx:   userCache.ctx,
		Key:   key,
		Value: user,
		TTL:   duration,
	})

	return err
}

func (userCache *Cache) GetToken(key string) (string, error) {
	var token string
	err := userCache.redisCache.Get(userCache.ctx, key, &token)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (userCache *Cache) SaveToken(key string, token string, duration time.Duration) error {
	err := userCache.redisCache.Set(&cache.Item{
		Ctx:   userCache.ctx,
		Key:   key,
		Value: token,
		TTL:   duration,
	})

	return err
}
