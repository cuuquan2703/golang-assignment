package repositories

import (
	"database/sql"
	"errors"
	_ "os"
	_ "server/logger"
	_ "strconv"
	_ "time"

	_ "github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type BookAuthor struct {
	IdBook   string `json:"id_book"`
	IdAuthor int    `json:"id_author"`
}

type BookAuthorRepository struct {
	DB    *sql.DB
	Table string
}

func NewBookAuthorRepository() (*BookAuthorRepository, error) {
	db, err := ConnectDB()
	if err != nil {
		return nil, err
	}
	return &BookAuthorRepository{
		DB:    db,
		Table: "Book",
	}, nil
}

func (repo BookAuthorRepository) GetAll() ([]BookAuthor, error) {
	Bookauthors := []BookAuthor{}
	cmd := `SELECT id,name,birth_date from Book_Author`
	L.Info("Querying " + cmd)
	row, err := repo.DB.Query(cmd)
	if err != nil {
		L.Error("Error", err)
	} else {
		L.Info("Query successfully")
	}

	for row.Next() {
		Bookauthor := BookAuthor{}
		err := row.Scan(&Bookauthor.IdBook, &Bookauthor.IdAuthor)
		if err != nil {
			L.Error("Error", err)
			return nil, err
		}
		Bookauthors = append(Bookauthors, Bookauthor)
	}

	if len(Bookauthors) == 0 {
		L.Error("Error ", errors.New("no Bookauthor found"))
		return nil, errors.New("No Bookauthor found")
	}
	defer row.Close()
	return Bookauthors, err
}

func (repo BookAuthorRepository) Insert(IdBook string, IdAuthor int) (sql.Result, error) {
	L.Info("Insert in to book_author table")
	res, err := repo.DB.Exec(" INSERT INTO Book_Author (id_book,id_author) VALUES ($1, $2);", IdBook, IdAuthor)
	if err != nil {
		L.Error("Error ", err)
	} else {
		L.Info("Insert successfully")
	}
	return res, err
}

func (repo BookAuthorRepository) Delete(IdBook string) (sql.Result, error) {
	L.Info("Delete from book_author table")
	res, err := repo.DB.Exec(" DELETE FROM Book_Author WHERE id_book = $1;", IdBook)
	if err != nil {
		L.Error("Error ", err)
	} else {
		L.Info("Delete successfully")
	}
	return res, err
}
