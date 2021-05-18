package db

import (
	"database/sql"
	"time"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	_ "github.com/go-sql-driver/mysql"
)

// DB *grom.DB
type DB struct {
	*goqu.Database
}

// ConnectDB : connecting DB
func ConnectDB() (*DB, error) {
	// ?parseTime=true
	// https://stackoverflow.com/questions/45040319/unsupported-scan-storing-driver-value-type-uint8-into-type-time-time
	db, err := sql.Open("mysql", "api:keepgrowth$@/book_report?parseTime=true")

	if err != nil {
		panic(err)
	}

	// https://github.com/go-sql-driver/mysql/#important-settings
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetConnMaxIdleTime(10)

	errPing := db.Ping()
	if errPing != nil {
		panic(err.Error())
	}

	qb := goqu.New("mysql", db)

	return &DB{qb}, nil
}
