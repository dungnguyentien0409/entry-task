package services

import "entrytask/tcp-server/shared/model"

type ITCPService interface {
	Register(request model.Request) model.Response
	Login(request model.Request) model.Response
}
