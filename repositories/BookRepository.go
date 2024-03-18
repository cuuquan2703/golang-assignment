package repositories

import (
	"database/sql"
	"errors"
	"fmt"
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
	Author      Author `json:"author"`
}

type BookRepository struct {
	DB             *sql.DB
	Table          string
	AuthorRepo     AuthorRepository
	BookAuthorRepo BookAuthorRepository
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
		AuthorRepo: AuthorRepository{
			DB:    db,
			Table: "author",
		},
		BookAuthorRepo: BookAuthorRepository{
			DB:    db,
			Table: "book_author",
		},
	}, nil
}

// var AuthorRepo, _ = NewAuthorRepository()
// var repo.BookAuthorRepo, _ = Newrepo.BookAuthorRepository()

func (repo BookRepository) GetAllBooks() ([]Book, error) {
	books := []Book{}
	cmd := `SELECT b.isbn,b.name,ba.id_author,b.publish_year from Book b JOIN book_author ba ON b.isbn = ba.id_book`
	L.Info("Querying " + cmd)
	row, err := repo.DB.Query(cmd)
	if err != nil {
		L.Error("Error", err)
	} else {
		L.Info("Query successfully")
	}

	for row.Next() {

		book := Book{}
		err := row.Scan(&book.ISBN, &book.Name, &book.Author.Id, &book.PublishYear)

		if err != nil {
			L.Error("Error", err)
			return nil, err
		}
		author, _ := repo.AuthorRepo.GetByID(book.Author.Id)
		book.Author = author
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
	cmd := `SELECT b.isbn,b.name,ba.id_author,b.publish_year from Book b 
			JOIN book_author ba ON b.isbn = ba.id_book 
			WHERE b.isbn = $1`
	L.Info("Querying " + cmd)
	row := repo.DB.QueryRow(cmd, isbn)
	L.Info("Query successfully")
	err := row.Scan(&book.ISBN, &book.Name, &book.Author.Id, &book.PublishYear)

	if err != nil {
		if err == sql.ErrNoRows {
			L.Error("Error ", errors.New("no books found"))
			return book, errors.New("No Book found")
		}
		L.Error("Error ", err)
		return book, errors.New("Something went wrong")
	}
	author, _ := repo.AuthorRepo.GetByID(book.Author.Id)
	book.Author = author
	return book, err
}

func (repo BookRepository) GetByAuthor(authorName string) ([]Book, error) {
	books := []Book{}
	cmd := `SELECT b.isbn,b.name,ba.id_author,b.publish_year from Book b 
			JOIN book_author ba ON b.isbn= ba.id_book 
			JOIN author a ON ba.id_author = a.id 
			WHERE a.name = $1`
	L.Info("Querying " + cmd)
	row, err := repo.DB.Query(cmd, authorName)
	if err != nil {
		L.Error("Error", err)
	} else {
		L.Info("Query successfully")
	}

	if err != nil {
		L.Error("Error ", err)
		return nil, err
	}
	defer row.Close()
	for row.Next() {
		book := Book{}
		err := row.Scan(&book.ISBN, &book.Name, &book.Author.Id, &book.PublishYear)
		if err != nil {
			L.Error("Error ", err)
			return nil, err
		}
		author, _ := repo.AuthorRepo.GetByID(book.Author.Id)
		book.Author = author
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
	cmd := `SELECT b.isbn,b.name,ba.id_author,b.publish_year from Book b 
			JOIN book_author ba ON b.isbn = ba.id_book  
			WHERE b.publish_year<=$2 and b.publish_year>=$1`
	L.Info("Querying " + cmd)
	row, err := repo.DB.Query(cmd, year1, year2)
	if err != nil {
		L.Error("Error", err)
	} else {
		L.Info("Query successfully")
	}

	if err != nil {
		L.Error("Error ", err)
		return nil, err
	}
	defer row.Close()
	for row.Next() {
		book := Book{}
		err := row.Scan(&book.ISBN, &book.Name, &book.Author.Id, &book.PublishYear)
		if err != nil {
			L.Error("Error ", err)
			return nil, err
		}
		author, _ := repo.AuthorRepo.GetByID(book.Author.Id)
		book.Author = author
		books = append(books, book)

	}

	if len(books) == 0 {
		L.Error("Error ", errors.New("no books found"))
		return nil, errors.New("No books found")
	}

	return books, err
}
func (repo BookRepository) Insert(data Book) (sql.Result, error) {
	author, err := repo.AuthorRepo.GetByName(data.Author.Name)
	if author == (Author{}) {
		_, err = repo.AuthorRepo.Insert(data.Author.Name)
		if err != nil {
			L.Error("Error insert author ", err)
		} else {
			L.Info("Insert successfully author")
		}
	}
	data.Author, _ = repo.AuthorRepo.GetByName(data.Author.Name)

	res, err2 := repo.DB.Exec(" INSERT INTO Book (isbn, name, publish_year) VALUES ($1, $2, $3);", data.ISBN, data.Name, data.PublishYear)
	_, err3 := repo.BookAuthorRepo.Insert(data.ISBN, data.Author.Id)
	if err2 != nil {
		L.Error("Error insert books ", err2)
	} else {
		L.Info("Insert successfully books")
	}
	if err3 != nil {
		L.Error("Error insert book_author ", err3)
	} else {
		L.Info("Insert successfully book_author")
	}
	return res, err
}

func (repo BookRepository) Delete(data Book) (sql.Result, error) {
	_, err1 := repo.BookAuthorRepo.Delete(data.ISBN)
	if err1 != nil {
		L.Error("Error Delete author ", err1)
	} else {
		L.Info("Delete successfully author")
	}
	res, err := repo.DB.Exec("DELETE FROM Book WHERE isbn = $1", data.ISBN)
	if err != nil {
		L.Error("Error Delete author ", err)
	} else {
		L.Info("Delete successfully author")
	}
	return res, err
}

func (repo BookRepository) Update(data Book) (sql.Result, error) {
	fmt.Println(data)
	res, err2 := repo.DB.Exec("UPDATE Book SET name = $1, publish_year = $2 WHERE isbn = $3", data.Name, data.PublishYear, data.ISBN)
	if err2 != nil {
		L.Error("Error update book ", err2)
	} else {
		L.Info("update successfully book")
	}
	return res, err2
}
