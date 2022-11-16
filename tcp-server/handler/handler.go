package handler

import (
	"crypto/rsa"
	"encoding/json"
	"entrytask/tcp-server/services"
	"entrytask/tcp-server/shared/constants"
	"entrytask/tcp-server/shared/model"
)

var PrivateKey *rsa.PrivateKey

type Handler struct {
	handlerRegistry map[int]func(request model.Request) (response model.Response)
	tcpService      services.ITCPService
}

func NewHandler(tcpService services.ITCPService) Handler {
	handlerRegistry := map[int]func(request model.Request) (response model.Response){
		constants.REGISTER_CMD: tcpService.Register,
		constants.LOGIN_CMD:    tcpService.Login,
	}

	return Handler{
		handlerRegistry: handlerRegistry,
		tcpService:      tcpService,
	}
}

func (handler *Handler) OnHandle(data []byte) ([]byte, error) {
	request := model.Request{}
	_ = json.Unmarshal(data, &request)

	fHandle := handler.handlerRegistry[request.Cmd]
	result := fHandle(request)

	resp, _ := json.Marshal(result)
	return resp, nil
}
