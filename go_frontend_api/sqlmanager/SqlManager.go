package sqlmanager

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

//Unused as of now, because we only implement SSO, and not database based authentication
func InitDB() (db *sql.DB) {
	// Connect to the postgres db
	//you might have to change the connection string to add your database credentials
	db, err := sql.Open("mysql", "XXXXXXXXXXXXX")
	if err != nil {
		panic(err)
	}
	return db
}
