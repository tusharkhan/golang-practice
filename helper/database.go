package helper

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

func ConnectDatabase() (*sql.DB, error) {
	envLoadingError := godotenv.Load()

	if envLoadingError != nil {
		panic(envLoadingError)
	}

	var DBNAME string = os.Getenv("DB_NAME")
	var DBHOST string = os.Getenv("DB_HOST")
	var DBUSER string = os.Getenv("DB_USER")
	var DBPASS string = os.Getenv("DB_PASSWPRD")
	var DBPORT string = os.Getenv("DB_PORT")

	// databse connection string
	var connectionString string = "host=" + DBHOST + " port=" + DBPORT + " user=" + DBUSER + " password=" + DBPASS + " dbname=" + DBNAME

	db, databaseConnectionError := sql.Open("pgx", connectionString)

	if databaseConnectionError != nil {
		panic(databaseConnectionError)
	}

	dbPingError := db.Ping()

	if dbPingError != nil {
		return nil, dbPingError
	}

	fmt.Println("Connected to database...")

	return db, nil
}
