package repositories

import (
	"database/sql"
	"log"
	"os"
	er "server/error"

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

func NewBookRepository() (*BookRepository, error) {
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

	return &BookRepository{
		DB:    db,
		Table: "Book",
	}, nil
}

func (repo BookRepository) GetAllBooks() ([]Book, error) {
	books := []Book{}
	cmd := `SELECT isbn,name,author,publish_year from Book`
	row, err := repo.DB.Query(cmd)
	defer row.Close()
	for row.Next() {
		book := Book{}
		err := row.Scan(&book.ISBN, &book.Name, &book.Author, &book.PublishYear)
		if err != nil {
			return nil, er.New("Something went wrong")
		}
		books = append(books, book)
	}

	if len(books) == 0 {
		return nil, er.New("No books found")
	}

	return books, err
}

func (repo BookRepository) GetByISBN(isbn string) (Book, error) {
	book := Book{}
	cmd := `SELECT isbn,name,author,publish_year from Book where "isbn"=$1`
	row := repo.DB.QueryRow(cmd, isbn)
	err := row.Scan(&book.ISBN, &book.Name, &book.Author, &book.PublishYear)

	if err != nil {
		if err == sql.ErrNoRows {
			return book, er.New("No Book found")
		}

		return book, er.New("Something went wrong")
	}

	return book, err
}

func (repo BookRepository) GetByAuthor(author string) ([]Book, error) {
	books := []Book{}
	cmd := `SELECT isbn,name,author,publish_year from $1 where "author"=$2`
	row, err := repo.DB.Query(cmd, repo.Table, author)

	if err != nil {
		return nil, er.New("Something went wrong")
	}
	defer row.Close()
	for row.Next() {
		book := Book{}
		err := row.Scan(&book.ISBN, &book.Name, &book.Author, &book.PublishYear)
		if err != nil {
			return nil, er.New("Something went wrong")
		}
		books = append(books, book)

	}

	if len(books) == 0 {
		return nil, er.New("No books found")
	}

	return books, err
}

func (repo BookRepository) GetInRange(year1, year2 int) ([]Book, error) {
	books := []Book{}
	cmd := `SELECT isbn,name,author,publish_year from $1 where "publish_year"<=$3 and "publish_year">=$2`
	row, err := repo.DB.Query(cmd, repo.Table, year1, year2)

	if err != nil {
		return nil, er.New("Something went wrong")
	}
	defer row.Close()
	for row.Next() {
		book := Book{}
		err := row.Scan(&book.ISBN, &book.Name, &book.Author, &book.PublishYear)
		if err != nil {
			return nil, er.New("Something went wrong")
		}
		books = append(books, book)

	}

	if len(books) == 0 {
		return nil, er.New("No books found")
	}

	return books, err
}
