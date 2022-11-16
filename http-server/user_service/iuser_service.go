package user_service

import (
	"entrytask/http-server/shared/model"
)

type IUserService interface {
	Register(request model.RegisterRequest) (response model.Response)
	Login(request model.LoginRequest) (response model.Response)
}
