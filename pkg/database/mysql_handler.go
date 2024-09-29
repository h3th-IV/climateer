package database

import (
	"database/sql"
	"io"
)

var (
	_ Database  = &mysqlDatabase{}
	_ io.Closer = &mysqlDatabase{}
)

type mysqlDatabase struct {
	createUser            *sql.Stmt
	checkUser             *sql.Stmt
	getUserByEmail        *sql.Stmt
	getBySessionKey       *sql.Stmt
	insertUserMeasurement *sql.Stmt
	getCountryByName      *sql.Stmt
	getIndicatorByCode    *sql.Stmt
}

func NewSQLDatabase(db *sql.DB) (*mysqlDatabase, error) {
	var (
		createUser            = "insert into Users(first_name, last_name, email, password_hash, phone, edu_institute, session_key) values(?,?,?,?,?,?,?);"
		checkUser             = "SELECT * FROM Users WHERE email = ? AND password_hash=?;"
		getUserByEmail        = "SELECT * FROM Users WHERE email = ?;"
		getBySessionKey       = "SELECT * FROM Users WHERE session_key=?;"
		insertUserMeasurement = "INSERT INTO UserMeasurements(user_id, country_id, indicator_id, year, value) VALUES(?,?,?,?,?);"
		getCountryByName      = "SELECT id FROM Countries WHERE country_name = ?;"
		getIndicatorByCode    = "SELECT id FROM Indicators WHERE IndicatorCode = ?;"
		database              = &mysqlDatabase{}
		err                   error
	)
	if database.createUser, err = db.Prepare(createUser); err != nil {
		return nil, err
	}
	if database.checkUser, err = db.Prepare(checkUser); err != nil {
		return nil, err
	}
	if database.getUserByEmail, err = db.Prepare(getUserByEmail); err != nil {
		return nil, err
	}
	if database.getBySessionKey, err = db.Prepare(getBySessionKey); err != nil {
		return nil, err
	}
	if database.insertUserMeasurement, err = db.Prepare(insertUserMeasurement); err != nil {
		return nil, err
	}
	if database.getCountryByName, err = db.Prepare(getCountryByName); err != nil {
		return nil, err
	}
	if database.getIndicatorByCode, err = db.Prepare(getIndicatorByCode); err != nil {
		return nil, err
	}
	return database, nil
}
