package postgres

import (
	"database/sql"
	"log"
	"os"

)


func NewDb() (*sql.DB, error) {
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
	return db, nil
}



func NewDbTest() (*sql.DB, error) {
	return nil, nil
}
