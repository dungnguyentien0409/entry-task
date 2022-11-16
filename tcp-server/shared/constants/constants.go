package constants

const (
	REGISTER_CMD = 0
	LOGIN_CMD    = 1

	JWT_EXPIRE_TIME = 24

	PATH = "/Users/chris.nguyen/Desktop/entrytask/"

	DEFAULT_PAGE_INDEX = 0
	DEFAULT_PAGE_SIZE  = 10
)

const (
	REGISTER_SUCCESSED int = 1
	ACCOUNT_OR_PASSWORD_EMPTY
	ACCOUNT_EXISTED
	REGISTER_FAILED
	CANNOT_PARESE_REGISTER_REQUEST
	LOGIN_SUCCESSED
	ACCOUT_OR_PASSWORD_INCORRECT
	ACCOUNT_NOT_EXISTED
	WRONG_PASSWORD
	CANNOT_PARSE_LOGIN_REQUEST
)
