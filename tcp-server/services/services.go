package services

import (
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	dao "entrytask/tcp-server/db"
	"entrytask/tcp-server/shared/constants"
	"entrytask/tcp-server/shared/dto"
	"entrytask/tcp-server/shared/jwt"
	"entrytask/tcp-server/shared/model"
	"entrytask/tcp-server/shared/redis"
	"github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql"
	"math/rand"
	"time"
)

type TCPService struct {
	userRepository dao.IUserRepository
	userCache      cache.ICache
	privateKey     *rsa.PrivateKey
}

func NewTCPService(userRepository dao.IUserRepository, userCache cache.ICache, keyBytes []byte) *TCPService {
	rand.Seed(time.Now().UnixNano())
	privateKey, _ := jwt.ParseRSAPrivateKeyFromPEM(keyBytes)

	return &TCPService{
		userRepository: userRepository,
		userCache:      userCache,
		privateKey:     privateKey,
	}
}

func (p *TCPService) Register(request model.Request) model.Response {
	param := request.Data.(map[string]interface{})
	registerRequest := model.RegisterRequest{}
	response := model.Response{}

	jsonbody, err := json.Marshal(param)
	if err != nil {
		response.Status = constants.CANNOT_PARESE_REGISTER_REQUEST
		response.Message = "Register request is incorrect"
		return response
	}

	if err = json.Unmarshal(jsonbody, &registerRequest); err != nil {
		response.Status = constants.CANNOT_PARESE_REGISTER_REQUEST
		response.Message = "Register request is incorrect"
		return response
	}

	if registerRequest.Account == "" || registerRequest.PassWord == "" {
		response.Status = constants.ACCOUNT_OR_PASSWORD_EMPTY
		response.Message = "Account or Password is empty"
		return response
	}
	if registerRequest.Nickname == "" {
		registerRequest.Nickname = registerRequest.Account
	}
	_, e := p.userRepository.GetUserByAccount(registerRequest.Account)
	if e == nil {
		response.Status = constants.ACCOUNT_EXISTED
		response.Message = "Account existed"
		return response
	}

	salt := GenerateSalt()
	hashPassword := GenerateHashPassword(registerRequest.PassWord, salt)

	user := dto.User{Account: registerRequest.Account, Nickname: registerRequest.Nickname, Password: hashPassword, Salt: salt}
	_, err = p.userRepository.Add(user)
	if err != nil {
		response.Status = constants.REGISTER_FAILED
		response.Message = "register failed"
		return response
	}

	response.Status = constants.REGISTER_SUCCESSED
	response.Message = "register success"
	response.Data = nil
	return response
}

func (p *TCPService) Login(request model.Request) model.Response {
	param := request.Data.(map[string]interface{})
	loginRequest := model.LoginRequest{}
	response := model.Response{}

	jsonbody, err := json.Marshal(param)
	if err != nil {
		response.Status = constants.CANNOT_PARSE_LOGIN_REQUEST
		response.Message = "Login request is incorrect"
		return response
	}

	if err = json.Unmarshal(jsonbody, &loginRequest); err != nil {
		response.Status = constants.CANNOT_PARSE_LOGIN_REQUEST
		response.Message = "Login request is incorrect"
		return response
	}

	if loginRequest.Account == "" || loginRequest.PassWord == "" {
		response.Status = constants.ACCOUT_OR_PASSWORD_INCORRECT
		response.Message = "Account or password is incorrect"
		return response
	}

	user, err := p.userCache.Get(loginRequest.Account)

	if err != nil {
		user, err = p.userRepository.GetUserByAccount(loginRequest.Account)
		if err != nil {
			response.Status = constants.ACCOUNT_NOT_EXISTED
			response.Message = "Account not existed"
			return response
		}

		if !CheckPassword(loginRequest.PassWord, *user) {
			response.Status = constants.WRONG_PASSWORD
			response.Message = "Wrong password"
			return response
		}

		p.userCache.SaveWithExpire(user.Account, *user, 2*time.Minute)
	}

	token, _ := jwtHelper.GenerateToken(user.Account, constants.JWT_EXPIRE_TIME, p.privateKey)

	response.Status = constants.LOGIN_SUCCESSED
	response.Message = "login success"
	response.Data = model.LoginResponse{token}

	return response
}

func CheckPassword(password string, user dto.User) bool {
	hashPassword := GenerateHashPassword(password, user.Salt)

	return hashPassword == user.Password
}

func GenerateSalt() string {
	salt := ""

	for i := 0; i < 64; i++ {
		random_int := rand.Intn(74) + 48
		character := rune(random_int)
		salt += string(character)
	}

	return salt
}

func GenerateHashPassword(password string, salt string) string {
	hash := sha256.New()
	hash.Write([]byte(salt))
	hash.Write([]byte(password))

	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}
