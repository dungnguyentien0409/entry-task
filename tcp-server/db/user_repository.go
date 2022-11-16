package dao

import (
	"context"
	"database/sql"
	"entrytask/tcp-server/shared/dto"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(DB *sql.DB) *UserRepository {
	return &UserRepository{
		db: DB,
	}
}

func (u *UserRepository) Add(user dto.User) (*dto.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*3))
	defer cancel()
	res, err := u.db.ExecContext(ctx, "INSERT INTO user_tab (account, nickname, password, salt) values(?, ? ,?, ?)",
		user.Account, user.Nickname, user.Password, user.Salt)
	if err != nil {
		return nil, err
	}

	userId, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	user.Id = userId
	return &user, nil
}

func (u *UserRepository) GetUserByAccount(account string) (*dto.User, error) {
	statement, err := u.db.Prepare("select id, account, nickname, password, salt from user_tab where account = ?")
	if err != nil {
		return nil, err
	}
	defer statement.Close()

	row := statement.QueryRow(account)

	user := dto.User{}
	err = row.Scan(&user.Id, &user.Account, &user.Nickname, &user.Password, &user.Salt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
