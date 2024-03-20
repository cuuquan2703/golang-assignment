package service_test

import (
	"database/sql"
	"log"
	"reflect"
	"regexp"
	"server/repositories"
	"server/service"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return db, mock
}

var db, mock = NewMock()

var repo = &repositories.BookRepository{
	DB:    db,
	Table: "Book",
}

var bookService = service.BookService{
	Repo: repo,
}

func TestGetAllBooks(t *testing.T) {
	expected := []repositories.Book{
		{ISBN: "19123450", Name: "Atomic", Author: "Grahahm", PublishYear: 2022},
		{ISBN: "12235670", Name: "Skinner", Author: "Albert", PublishYear: 2001},
		{ISBN: "12223900", Name: "Short", Author: "Victor", PublishYear: 1998},
	}
	rows := sqlmock.NewRows([]string{"isbn", "name", "author", "publish_year"}).
		AddRow("19123450", "Atomic", "Grahahm", 2022).
		AddRow("12235670", "Skinner", "Albert", 2001).
		AddRow("12223900", "Short", "Victor", 1998)

	mock.ExpectQuery("SELECT isbn,name,author,publish_year from Book").WillReturnRows(rows)

	book, err := bookService.GetAllBooks()
	if err != nil {
		log.Fatal(err)
	}

	if !reflect.DeepEqual(book, expected) {
		t.Errorf("Returned books don't match expected books. Expected: %v, Actual: %v", expected, book)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestGetByISBN(t *testing.T) {
	expected := []repositories.Book{
		{ISBN: "12235670", Name: "Skinner", Author: "Albert", PublishYear: 2001},
	}

	expectedRows := sqlmock.NewRows([]string{"isbn", "nam", "author", "publish_year"}).
		AddRow("12235670", "Skinner", "Albert", 2001)
	mock.ExpectQuery(`SELECT (.*)`).WillReturnRows(expectedRows)

	book, err := bookService.GetByISBN("12235670")
	if err != nil {
		log.Fatal(err)
	}

	if !reflect.DeepEqual(book, expected[0]) {
		t.Errorf("Returned books don't match expected books. Expected: %v, Actual: %v", expected, book)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestGetByAuthor(t *testing.T) {
	expected := []repositories.Book{
		{ISBN: "12235670", Name: "Skinner", Author: "Albert", PublishYear: 2001},
		{ISBN: "12289970", Name: "Stlake", Author: "Albert", PublishYear: 1997},
	}
	expectedRows := sqlmock.NewRows([]string{"isbn", "nam", "author", "publish_year"}).
		AddRow("12235670", "Skinner", "Albert", 2001).
		AddRow("12289970", "Stlake", "Albert", 1997)
	mock.ExpectQuery(`SELECT (.*)`).WillReturnRows(expectedRows)

	book, err := bookService.GetByAuthor("Albert")
	if err != nil {
		log.Fatal(err)
	}

	if !reflect.DeepEqual(book, expected) {
		t.Errorf("Returned books don't match expected books. Expected: %v, Actual: %v", expected, book)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestGetInRange(t *testing.T) {
	expected := []repositories.Book{
		{ISBN: "12235670", Name: "Skinner", Author: "Albert", PublishYear: 2001},
		{ISBN: "19123450", Name: "Atomic", Author: "Grahahm", PublishYear: 2022},
	}
	expectedRows := sqlmock.NewRows([]string{"isbn", "nam", "author", "publish_year"}).
		AddRow("12235670", "Skinner", "Albert", 2001).
		AddRow("19123450", "Atomic", "Grahahm", 2022)
	mock.ExpectQuery(`SELECT (.*)`).WillReturnRows(expectedRows)

	book, err := bookService.GetInRange(1999, 2023)
	if err != nil {
		log.Fatal(err)
	}

	if !reflect.DeepEqual(book, expected) {
		t.Errorf("Returned books don't match expected books. Expected: %v, Actual: %v", expected, book)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestUpdate(t *testing.T) {

	var bookData = []repositories.Book{
		{ISBN: "19123450", Name: "Update 1", Author: "Author 1", PublishYear: 2022},
		{ISBN: "19126450", Name: "Update 2", Author: "Author 2", PublishYear: 2021},
	}

	for _, data := range bookData {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT isbn,name,author,publish_year from Book where "isbn"=$1`)).
			WithArgs(data.ISBN).
			WillReturnRows(sqlmock.NewRows([]string{"isbn", "name", "author", "publish_year"}).
				AddRow(data.ISBN, "Atomic", "Grahahm", 2022))
		mock.ExpectExec(regexp.QuoteMeta("UPDATE Book SET name = $1, publish_year = $2, author = $3 WHERE isbn = $4")).
			WithArgs(data.Name, data.PublishYear, data.Author, data.ISBN).
			WillReturnResult(sqlmock.NewResult(0, 1))
	}

	err := bookService.Update(bookData)
	if err != nil {
		t.Errorf("Error when updating db")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestDelete(t *testing.T) {

	var bookData = []repositories.Book{
		{ISBN: "19123450", Name: "Name 1", Author: "Author 1", PublishYear: 2022},
		{ISBN: "19126450", Name: "Name 1", Author: "Author 1", PublishYear: 2022},
	}

	for _, data := range bookData {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT isbn,name,author,publish_year from Book where "isbn"=$1`)).
			WithArgs(data.ISBN).
			WillReturnRows(sqlmock.NewRows([]string{"isbn", "name", "author", "publish_year"}).
				AddRow(data.ISBN, "Name 1", "Author 1", 2022))
		mock.ExpectExec(regexp.QuoteMeta("DELETE FROM Book WHERE isbn = $1")).
			WithArgs(data.ISBN).
			WillReturnResult(sqlmock.NewResult(0, 1))
	}

	err := bookService.Delete(bookData)
	if err != nil {
		t.Errorf("Error when delete db")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestInsert(t *testing.T) {

	var bookData = []repositories.Book{
		{ISBN: "19123450", Name: "Name 1", Author: "Author 1", PublishYear: 2022},
		{ISBN: "19126450", Name: "Name 2", Author: "Author 2", PublishYear: 2024},
	}

	for _, data := range bookData {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT isbn,name,author,publish_year from Book where "isbn"=$1`)).
			WithArgs(data.ISBN).
			WillReturnRows(sqlmock.NewRows([]string{"isbn", "name", "author", "publish_year"}).
				AddRow(nil, nil, nil, nil))
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO Book (isbn, name, publish_year, author) VALUES ($1, $2, $3, $4)")).
			WithArgs(data.ISBN, data.Name, data.PublishYear, data.Author).
			WillReturnResult(sqlmock.NewResult(0, 1))
	}

	err := bookService.Insert(bookData)
	if err != nil {
		t.Errorf("Error when inserting db")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}
