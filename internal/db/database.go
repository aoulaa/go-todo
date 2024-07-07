package db

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"os"
	"strconv"
)

var Db *sql.DB

func ConnectDatabase() {
	err := godotenv.Load("/Users/aoulaa/Documents/projects/go-todo/.env")
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	host := os.Getenv("POSTGRES_HOST")
	port, err := strconv.Atoi(os.Getenv("POSTGRES_PORT"))
	if err != nil {
		fmt.Println("Error while converting port to integer")
		panic(err)
	}
	user := os.Getenv("POSTGRES_USER")
	dbname := os.Getenv("POSTGRES_DB")
	password := os.Getenv("POSTGRES_PASSWORD")

	psqlSetup := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		host, port, user, dbname, password)
	db, errSql := sql.Open("postgres", psqlSetup)

	if errSql != nil {
		fmt.Println("There is an error while connecting to the database ", err)
		panic(err)
	} else {
		Db = db
		fmt.Println("Successfully connected to database!")
	}
}

func CloseDatabase() {
	Db.Close()
	fmt.Println("Database connection closed!")
}
