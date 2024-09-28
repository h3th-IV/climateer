package database

import (
	"context"
	"log"

	"github.com/h3th-IV/climateer/pkg/model"
)

func (db *mysqlDatabase) CreateUser(ctx context.Context, first_name string, last_name string, email string, password string, phone string, eduInstitute string, sessionkey string) (bool, error) {
	userQuery, err := db.createUser.ExecContext(ctx, first_name, last_name, email, password, phone, eduInstitute, sessionkey)
	if err != nil {
		return false, err
	}
	lastId, err := userQuery.LastInsertId()
	if err != nil {
		return false, err
	}
	if lastId == 0 || lastId < 1 {
		return false, err
	}
	return true, nil
}

func (db *mysqlDatabase) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	user := &model.User{}
	getUserByEmail := db.getUserByEmail.QueryRowContext(ctx, email)
	err := getUserByEmail.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.Phone, &user.EduInstitute, &user.SessionKey, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		log.Println("get user by email", err)
		return nil, err
	}
	return user, nil
}

func (db *mysqlDatabase) CheckUser(ctx context.Context, email string, password string) (*model.User, error) {
	user := &model.User{}
	getUserByEmail := db.checkUser.QueryRowContext(ctx, email, password)
	err := getUserByEmail.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.Phone, &user.EduInstitute, &user.SessionKey, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		log.Println("checkuser", err)
		return nil, err
	}
	return user, nil
}

func (db *mysqlDatabase) GetBySessionKey(ctx context.Context, sessionkey string) (*model.User, error) {
	user := &model.User{}
	getBySessionKey := db.getBySessionKey.QueryRowContext(ctx, sessionkey)
	err := getBySessionKey.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.Phone, &user.EduInstitute, &user.SessionKey, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}
