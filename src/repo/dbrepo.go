package repo

import (
	"database/sql"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"log"
	"time"
)

const maxOpenDBConn = 10
const maxIdleDBConn = 5
const maxDBLifeTime = 5 * time.Minute

var dbConn = &PostgresDBRepo{}

type PostgresDBRepo struct {
	DB  *sql.DB
}

func newDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Println("Error occurred in getting the database connection")
		return nil, err
	}
	
	err = db.Ping()
	if err != nil {
		log.Println("Error occurred in pinging database")
		return nil, err
	}
	
	return db, err
}

func testDB(db *sql.DB) error {
	err := db.Ping()
	return err
}

func ConnectSQL(dsn string) (*PostgresDBRepo, error) {
	db, err := newDB(dsn)
	if err != nil {
		panic("This error occurred in driver.go[ConnectSQL func] : " + err.Error())
	}

	db.SetMaxOpenConns(maxOpenDBConn)
	db.SetConnMaxIdleTime(maxIdleDBConn)
	db.SetConnMaxLifetime(maxDBLifeTime)

	dbConn.DB = db
	err = testDB(db)
	if err != nil {
		return nil, err
	}

	return dbConn, nil
}
