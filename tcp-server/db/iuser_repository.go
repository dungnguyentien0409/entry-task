package dao

import "entrytask/tcp-server/shared/dto"

type IUserRepository interface {
	Add(user dto.User) (*dto.User, error)
	GetUserByAccount(username string) (*dto.User, error)
}
