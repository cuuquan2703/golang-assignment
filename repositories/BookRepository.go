package repositories

import (
	"database/sql"
	"errors"
	"os"
	"server/logger"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const DB_URL = "DB_URL"

type Book struct {
	ISBN        string `json:"isbn"`
	Name        string `json:"name"`
	PublishYear int    `json:"publish_year"`
	Author      string `json:"author"`
}

type BookRepository struct {
	DB    *sql.DB
	Table string
}

var L = logger.CreateLog()

func ConnectDB() (*sql.DB, error) {
	var url string
	err := godotenv.Load()
	if err != nil {
		L.Error("Error loading .env file:", err)
	}
	if _, exist := os.LookupEnv(DB_URL); exist {
		url = os.Getenv(DB_URL)
	} else {
		L.Error("Error loading .env file:", err)
	}

	db, err := sql.Open("postgres", url)
	if err != nil {
		L.Error("Error open db:", err)
	}
	return db, err
}

func NewBookRepository() (*BookRepository, error) {
	db, err := ConnectDB()
	if err != nil {
		return nil, err
	}
	return &BookRepository{
		DB:    db,
		Table: "Book",
	}, nil
}

func (repo BookRepository) GetAllBooks() ([]Book, error) {
	books := []Book{}
	cmd := `SELECT isbn,name,author,publish_year from Book`
	L.Info("Querying " + cmd)
	row, err := repo.DB.Query(cmd)
	L.Info("Query successfully")

	for row.Next() {
		book := Book{}
		err := row.Scan(&book.ISBN, &book.Name, &book.Author, &book.PublishYear)
		if err != nil {
			L.Error("Error", err)
			return nil, err
		}
		books = append(books, book)
	}

	if len(books) == 0 {
		L.Error("Error ", errors.New("no books found"))
		return nil, errors.New("No books found")
	}
	defer row.Close()
	return books, err
}

func (repo BookRepository) GetByISBN(isbn string) (Book, error) {
	book := Book{}
	cmd := `SELECT isbn,name,author,publish_year from Book where "isbn"=$1`
	L.Info("Querying " + cmd)
	row := repo.DB.QueryRow(cmd, isbn)
	L.Info("Query successfully")
	err := row.Scan(&book.ISBN, &book.Name, &book.Author, &book.PublishYear)

	if err != nil {
		if err == sql.ErrNoRows {
			L.Error("Error ", errors.New("no books found"))
			return book, errors.New("No Book found")
		}
		L.Error("Error ", err)
		return book, errors.New("Something went wrong")
	}

	return book, err
}

func (repo BookRepository) GetByAuthor(author string) ([]Book, error) {
	books := []Book{}
	cmd := `SELECT isbn,name,author,publish_year from Book where "author"=$1`
	L.Info("Querying " + cmd)
	row, err := repo.DB.Query(cmd, author)
	L.Info("Query successfully")

	if err != nil {
		L.Error("Error ", err)
		return nil, err
	}
	defer row.Close()
	for row.Next() {
		book := Book{}
		err := row.Scan(&book.ISBN, &book.Name, &book.Author, &book.PublishYear)
		if err != nil {
			L.Error("Error ", err)
			return nil, err
		}
		books = append(books, book)

	}

	if len(books) == 0 {
		L.Error("Error ", errors.New("no books found"))
		return nil, errors.New("No books found")
	}

	return books, err
}

func (repo BookRepository) GetInRange(year1, year2 int) ([]Book, error) {
	books := []Book{}
	cmd := `SELECT isbn,name,author,publish_year from Book where "publish_year"<=$2 and "publish_year">=$1`
	L.Info("Querying " + cmd)
	row, err := repo.DB.Query(cmd, year1, year2)
	L.Info("Query successfully")

	if err != nil {
		L.Error("Error ", err)
		return nil, err
	}
	defer row.Close()
	for row.Next() {
		book := Book{}
		err := row.Scan(&book.ISBN, &book.Name, &book.Author, &book.PublishYear)
		if err != nil {
			L.Error("Error ", err)
			return nil, err
		}
		books = append(books, book)

	}

	if len(books) == 0 {
		L.Error("Error ", errors.New("no books found"))
		return nil, errors.New("No books found")
	}

	return books, err
}
