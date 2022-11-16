package cache

import (
	"entrytask/tcp-server/shared/dto"
	"time"
)

type ICache interface {
	Save(key string, user dto.User) error
	SaveWithExpire(key string, user dto.User, duration time.Duration) error
	Get(key string) (*dto.User, error)
	GetToken(key string) (string, error)
	SaveToken(key string, token string, duration time.Duration) error
}
