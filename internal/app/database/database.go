package database

import (
	"database/sql"
	"fmt"
	errorlog "jwt/internal/app/errorLog"
	"jwt/internal/structs"
	"os"

	_ "github.com/lib/pq"
)

type Database struct {
	db *sql.DB
}

func New() *Database {
	db, err := initDB()
	if err != nil {
		errorlog.ErrorPrintFatal("can not connect to db", err)
	}

	err = createTables(db)
	if err != nil {
		errorlog.ErrorPrintFatal("error creating tables", err)
	}

	return &Database{db: db}
}

func createTables(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS refreshtoken (
    hash TEXT NOT NULL,
    ip TEXT NOT NULL UNIQUE);`
	_, err := db.Exec(query)
	return err
}

func (DB *Database) PutInDB(hash []byte, ip string) error {
	query := `INSERT INTO refreshtoken (hash, ip)
	VALUES ($1, $2)
	ON CONFLICT (ip)
	DO UPDATE SET hash = EXCLUDED.hash;`

	_, err := DB.db.Query(query, hash, ip)
	if err != nil {
		return err
	}
	return nil

}

func (DB *Database) GetDataFromDB(ip string) (structs.RefreshToken, error) {
	query := "SELECT * FROM refreshtoken WHERE ip = $1;"

	var hash, dbIP string
	var data structs.RefreshToken
	err := DB.db.QueryRow(query, ip).Scan(&hash, &dbIP)
	if err != nil {
		if err == sql.ErrNoRows {
			return data, nil
		}
		return data, err
	}
	data.Hash = []byte(hash)
	data.IP = dbIP
	return data, nil
}

func initDB() (*sql.DB, error) {

	psqlconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", os.Getenv("host"), os.Getenv("port"), os.Getenv("user"), os.Getenv("password"), os.Getenv("dbname"))
	db, err := sql.Open("postgres", psqlconn)
	return db, err
}
