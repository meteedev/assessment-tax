package postgres

import (
	"database/sql"
	"log"
	"os"

)

type Postgres struct {
	Db *sql.DB
}

type DbConfiguration struct {
	Host     string
	Port     string
	User     string
	Password string
	DbName   string
	SslMode  string
}

func New() (*Postgres, error) {
	databaseSource := os.Getenv("DATABASE_URL")
	log.Println(databaseSource)
	db, err := sql.Open("postgres", databaseSource)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return &Postgres{Db: db}, nil
}




