package main

import (
	"database/sql"
	"os"
	er "server/error"
	"github.com/joho/godotenv"
	"log"
	_ "github.com/lib/pq" 
)

const DB_URL = "DB_URL"

func main() {
	var url string
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}
	if _, exist := os.LookupEnv(DB_URL); exist {
		url = os.Getenv(DB_URL)
	} else {
		er.Check(er.New("Env variable not found"))
	}

	db, err := sql.Open("postgres", url)
	er.Check(err)

	_,e := db.Exec(`Call insert_book($1, $2 , $3, $4)`,"133", "Abcyx" , 2021, "a1")
	er.Check(e)
	_,e1 := db.Exec(`Call insert_book($1, $2 , $3, $4)`,"129", "Bo" , 2009, "a10")
	er.Check(e1)
	_,e2 := db.Exec(`Call insert_book($1, $2 , $3, $4)`,"273", "Atomic " , 2010, "a2")
	er.Check(e2)
	_,e3 := db.Exec(`Call insert_book($1, $2 , $3, $4)`,"103", "Atomic 3" , 2023, "a4")
	er.Check(e3)

	defer db.Close()
}