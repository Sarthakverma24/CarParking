package wrapper

import (
	"CarParking/db"
	"context"
	"database/sql"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type UserWrapper struct {
	Db db.Database
}

type SigninUser struct {
	Name     string `json:"name"`
	Phone    string `json:"number"`
	Email    string `json:"email"`
	Age      int    `json:"age"`
	UserName string `json:"username"`
	Password string `json:"password"`
}
type User struct {
	Username string `json:"user_name"`
	Password string `json:"password"`
}

type Users struct {
	UserName string `json:"user_name"`
	Name     string `json:"name"`
	PhoneNo  string `json:"phone_no"`
	Gmail    string `json:"gmail"`
	Age      int    `json:"age"`
}

func (d *UserWrapper) Login(ctx context.Context, user User) (string, error) {
	logger := db.Logger.Sugar()
	logger.Infow("Login user request", "username", user.Username)

	if d.Db == nil {
		db.Logger.Error("Database instance is nil")
		return "", errors.New("internal server error")
	}

	query := `SELECT user_name FROM user_data WHERE user_name = $1 AND password = $2`
	result, err := d.Db.GetData(ctx, query, user.Username, user.Password)
	if err != nil {
		db.Logger.Error("Failed to login user", zap.Error(err))
		return "", err
	}

	rows, ok := result.(*sql.Rows)
	if !ok {
		db.Logger.Error("Invalid result format: not *sql.Rows")
		return "", errors.New("internal error: invalid query result")
	}
	defer rows.Close()

	if !rows.Next() {
		// No rows returned
		return "", errors.New("invalid email or password")
	}

	var username string
	if err := rows.Scan(&username); err != nil {
		db.Logger.Error("Failed to scan row", zap.Error(err))
		return "", errors.New("internal error: scan failed")
	}

	return username, nil
}

func (d *UserWrapper) SignIn(ctx context.Context, user SigninUser) (string, error) {
	query := `INSERT INTO user_data (name, phone_no, gmail, age, user_name, password) VALUES ($1, $2, $3, $4, $5, $6)`
	err := d.Db.SetData(ctx, query, user.Name, user.Phone, user.Email, user.Age, user.UserName, user.Password)
	if err != nil {
		return "", err
	}
	return user.UserName, nil
}

func (d *UserWrapper) GetUsers(ctx context.Context) ([]Users, error) {
	if d.Db == nil {
		db.Logger.Error("Database instance is nil")
		return nil, errors.New("internal server error")
	}

	query := `SELECT user_name, name, phone_no, gmail, age FROM "user_data"`
	result, err := d.Db.GetData(ctx, query)
	if err != nil {
		db.Logger.Error("Failed to fetch lot data", zap.Error(err))
		return nil, err
	}

	rows, ok := result.(*sql.Rows)
	if !ok {
		db.Logger.Error("Invalid result format: not *sql.Rows")
		return nil, errors.New("internal error: invalid query result")
	}
	defer rows.Close()

	var users []Users
	for rows.Next() {
		var user Users
		err := rows.Scan(
			&user.UserName,
			&user.Name,
			&user.PhoneNo,
			&user.Gmail,
			&user.Age,
		)
		if err != nil {
			// handle error
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		db.Logger.Error("Row iteration error", zap.Error(err))
		return nil, errors.New("internal error: row iteration")
	}

	return users, nil
}
