package user_service

import (
	"encoding/json"
	"entrytask/http-server/shared/model"
	"entrytask/protocol"
	"entrytask/tcp-server/shared/constants"
)

type UserService struct {
	tcpClient *protocol.Client
}

func NewUserService(minCap int, maxCap int, addr string) (*UserService, error) {
	tcpClient, er := protocol.NewClient(minCap, maxCap, addr)
	return &UserService{
		tcpClient: tcpClient}, er
}

func (userService *UserService) Register(request model.RegisterRequest) (response model.Response) {
	response = Call(constants.REGISTER_CMD, request, userService.tcpClient.Call)
	return response
}

func (userService *UserService) Login(request model.LoginRequest) (response model.Response) {
	response = Call(constants.LOGIN_CMD, request, userService.tcpClient.Call)
	return response
}

func Call(cmd int, request interface{}, sendRequest func(data []byte) ([]byte, error)) model.Response {
	tReq := model.Request{
		Cmd:  cmd,
		Data: request,
	}

	data, _ := json.Marshal(tReq)
	rawResp, _ := sendRequest(data)

	resp := model.Response{}
	_ = json.Unmarshal(rawResp, &resp)
	return resp
}
