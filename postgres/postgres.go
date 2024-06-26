package postgres

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)


func NewDb() (*sql.DB, error) {
	databaseSource := os.Getenv("DATABASE_URL")
	//log.Println(databaseSource)
	db, err := sql.Open("postgres", databaseSource)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}	
	return db, nil
}



