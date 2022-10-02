package mocks

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

func NewDatabase() (*gorm.DB, *sql.DB, sqlmock.Sqlmock) {
	// get db and mock
	sqlDB, mock, err := sqlmock.New(
	//sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp),
	)
	if err != nil {
		log.Fatalf("[sqlmock new] %s", err)
	}

	// create dialector
	dialector := mysql.New(mysql.Config{
		DSN:                       "sqlmock_db_0",
		Conn:                      sqlDB,
		DriverName:                "mysql",
		SkipInitializeWithVersion: true,
	})

	// open the database
	db, err := gorm.Open(dialector, &gorm.Config{PrepareStmt: false})
	if err != nil {
		log.Fatalf("[gorm open] %s", err)
	}

	return db, sqlDB, mock
}
