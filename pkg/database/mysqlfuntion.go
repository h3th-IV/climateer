package database

import (
	"context"
	"database/sql"
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

func (db *mysqlDatabase) AddUserMeasurement(ctx context.Context, userID, countryID, indicatorID, year int, value float64) (bool, error) {
	result, err := db.insertUserMeasurement.ExecContext(ctx, userID, countryID, indicatorID, year, value)
	if err != nil {
		log.Println("error adding user measurement:", err)
		return false, err
	}
	lid_result, err := result.LastInsertId()
	if err != nil || lid_result == 0 {
		log.Println("no new rows inserted or error:", err)
		return false, err
	}
	rwa_result, err := result.RowsAffected()
	if err != nil || rwa_result == 0 {
		log.Println("no rows affected or error:", err)
		return false, err
	}

	return true, nil
}

func (db *mysqlDatabase) GetCountryIDByName(ctx context.Context, countryName string) (int, error) {
	var countryID int
	err := db.getCountryByName.QueryRowContext(ctx, countryName).Scan(&countryID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("no country found with that name:", countryName)
			return 0, nil
		}
		log.Println("error fetching country by name:", err)
		return 0, err
	}
	return countryID, nil
}

func (db *mysqlDatabase) GetIndicatorIDByCode(ctx context.Context, indicatorCode string) (int, error) {
	var indicatorID int
	err := db.getIndicatorByCode.QueryRowContext(ctx, indicatorCode).Scan(&indicatorID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("no indicator found with that code:", indicatorCode)
			return 0, nil
		}
		log.Println("error fetching indicator by code:", err)
		return 0, err
	}
	return indicatorID, nil
}
