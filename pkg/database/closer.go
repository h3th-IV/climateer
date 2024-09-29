package database

func (db *mysqlDatabase) Close() error {
	db.createUser.Close()
	db.checkUser.Close()
	db.getUserByEmail.Close()
	db.getBySessionKey.Close()
	db.insertUserMeasurement.Close()
	db.getIndicatorByCode.Close()
	return nil
}
