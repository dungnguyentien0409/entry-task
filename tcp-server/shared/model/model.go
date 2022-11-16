package model

type Request struct {
	Cmd  int         `json:"cmd_type"`
	Data interface{} `json:"data"`
}

type Response struct {
	Status  int         `json:"error_code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type RegisterRequest struct {
	Account  string `json:"account"`
	PassWord string `json:"password"`
	Nickname string `json:"nickname"`
}

type LoginRequest struct {
	Account  string `json:"account"`
	PassWord string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}
