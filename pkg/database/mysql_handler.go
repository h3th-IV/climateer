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
	createUser      *sql.Stmt
	checkUser       *sql.Stmt
	getUserByEmail  *sql.Stmt
	getBySessionKey *sql.Stmt
}

func NewSQLDatabase(db *sql.DB) (*mysqlDatabase, error) {
	var (
		createUser      = "insert into Users(first_name, last_name, email, password_hash, phone, edu_institute, session_key) values(?,?,?,?,?,?);"
		checkUser       = "SELECT * FROM users WHERE email = ? AND password=?;"
		getUserByEmail  = "SELECT * FROM users WHERE email = ?;"
		getBySessionKey = "SELECT * FROM users WHERE session_key=?;"
		database        = &mysqlDatabase{}
		err             error
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
	return database, nil
}

func (db *mysqlDatabase) Close() error {
	db.createUser.Close()
	db.checkUser.Close()
	db.getUserByEmail.Close()
	db.getBySessionKey.Close()
	return nil
}
