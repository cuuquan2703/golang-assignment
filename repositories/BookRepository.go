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
	Author      Author `json:"id_author"`
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

var AuthorRepo, _ = NewAuthorRepository()

func (repo BookRepository) GetAllBooks() ([]Book, error) {
	books := []Book{}
	cmd := `SELECT isbn,name,id_author,publish_year from Book`
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
		author, _ := AuthorRepo.GetByID(book.Author.Id)
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
	cmd := `SELECT isbn,name,id_author,publish_year from Book where "isbn"=$1`
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
	author, _ := AuthorRepo.GetByID(book.Author.Id)
	book.Author = author
	return book, err
}

func (repo BookRepository) GetByAuthor(author string) ([]Book, error) {
	books := []Book{}
	cmd := `SELECT isbn,name,id_author,publish_year from Book where "author"=$1`
	L.Info("Querying " + cmd)
	row, err := repo.DB.Query(cmd, author)
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
		err := row.Scan(&book.ISBN, &book.Name, &book.Author, &book.PublishYear)
		if err != nil {
			L.Error("Error ", err)
			return nil, err
		}
		author, _ := AuthorRepo.GetByID(book.Author.Id)
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
	cmd := `SELECT isbn,name,id_author,publish_year from Book where "publish_year"<=$2 and "publish_year">=$1`
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
		err := row.Scan(&book.ISBN, &book.Name, &book.Author, &book.PublishYear)
		if err != nil {
			L.Error("Error ", err)
			return nil, err
		}
		author, _ := AuthorRepo.GetByID(book.Author.Id)
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
	author, err := AuthorRepo.GetByName(data.Author.Name)
	if author == (Author{}) {
		_, err = AuthorRepo.Insert(data.Author.Name)
		if err != nil {
			L.Error("Error insert author ", err)
		} else {
			L.Info("Insert successfully author")
		}
	}
	data.Author, _ = AuthorRepo.GetByName(data.Author.Name)
	res, err2 := repo.DB.Exec(" INSERT INTO Book (isbn, name, publish_year, id_author) VALUES ($1, $2, $3, $4);", data.ISBN, data.Name, data.PublishYear, data.Author.Id)
	if err2 != nil {
		L.Error("Error insert books ", err2)
	} else {
		L.Info("Insert successfully books")
	}
	return res, err
}

func (repo BookRepository) Delete(data Book) (sql.Result, error) {
	res, err := repo.DB.Exec("DELETE FROM Book WHERE isbn = $1", data.ISBN)
	if err != nil {
		L.Error("Error insert author ", err)
	} else {
		L.Info("Insert successfully author")
	}
	return res, err
}

func (repo BookRepository) Update(data Book) (sql.Result, error) {
	_, err := AuthorRepo.Update(data.Author)
	if err != nil {
		L.Error("Error update author ", err)
	}
	data.Author, _ = AuthorRepo.GetByID(data.Author.Id)
	res, err2 := repo.DB.Exec("UPDATE Book SET name = $1, publish_year = $2, id_author = $3 WHERE isbn = $4", data.Name, data.PublishYear, data.Author.Id, data.ISBN)
	if err2 != nil {
		L.Error("Error insert author ", err2)
	} else {
		L.Info("Insert successfully author")
	}
	return res, err2
}
