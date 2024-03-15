package main

import (
	"database/sql"
	"errors"
	"os"
	"server/logger"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const DB_URL = "DB_URL"

var L = logger.CreateLog()

func main() {
	var url string
	err := godotenv.Load()
	if err != nil {
		L.Error("Error loading .env file:", err)
	}
	if _, exist := os.LookupEnv(DB_URL); exist {
		url = os.Getenv(DB_URL)
	} else {
		L.Error("Error ", errors.New("Env variable not found"))
	}

	db, err := sql.Open("postgres", url)
	if err != nil {
		L.Error("Error in open DB", err)
	}

	L.Info("Insert mock data")
	_, e := db.Exec(`Call insert_book($1, $2 , $3, $4)`, "133", "Abcyx", 2021, "a1")
	if e != nil {
		L.Error("Error while inserting mock data", e)
	}
	_, e1 := db.Exec(`Call insert_book($1, $2 , $3, $4)`, "129", "Bo", 2009, "a10")
	if e1 != nil {
		L.Error("Error while inserting mock data", e1)
	}
	_, e2 := db.Exec(`Call insert_book($1, $2 , $3, $4)`, "273", "Atomic ", 2010, "a2")
	if e2 != nil {
		L.Error("Error while inserting mock data", e2)
	}
	_, e3 := db.Exec(`Call insert_book($1, $2 , $3, $4)`, "103", "Atomic 3", 2023, "a4")
	if e3 != nil {
		L.Error("Error while inserting mock data", e3)
	}

	defer db.Close()
}
