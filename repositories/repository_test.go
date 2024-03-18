package repositories_test

import (
	"database/sql"
	"log"
	"reflect"
	repositories "server/repositories"
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

func TestGetAllBooks(t *testing.T) {
	db, mock := NewMock()

	repo := repositories.BookRepository{
		DB:    db,
		Table: "Book",
	}
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

	book, err := repo.GetAllBooks()
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
	db, mock := NewMock()

	repo := repositories.BookRepository{
		DB:    db,
		Table: "Book",
	}
	expected := []repositories.Book{
		{ISBN: "12235670", Name: "Skinner", Author: "Albert", PublishYear: 2001},
	}
	sqlmock.NewRows([]string{"isbn", "name", "author", "publish_year"}).
		AddRow("19123450", "Atomic", "Grahahm", 2022).
		AddRow("12235670", "Skinner", "Albert", 2001).
		AddRow("12223900", "Short", "Victor", 1998)
	expectedRows := sqlmock.NewRows([]string{"isbn", "nam", "author", "publish_year"}).
		AddRow("12235670", "Skinner", "Albert", 2001)
	mock.ExpectQuery(`SELECT (.*)`).WillReturnRows(expectedRows)

	book, err := repo.GetByISBN("12235670")
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
	db, mock := NewMock()

	repo := repositories.BookRepository{
		DB:    db,
		Table: "Book",
	}
	expected := []repositories.Book{
		{ISBN: "12235670", Name: "Skinner", Author: "Albert", PublishYear: 2001},
		{ISBN: "12289970", Name: "Stlake", Author: "Albert", PublishYear: 1997},
	}
	sqlmock.NewRows([]string{"isbn", "name", "author", "publish_year"}).
		AddRow("19123450", "Atomic", "Grahahm", 2022).
		AddRow("12235670", "Skinner", "Albert", 2001).
		AddRow("12223900", "Short", "Victor", 1998).
		AddRow("12289970", "Stlake", "Albert", 1997)
	expectedRows := sqlmock.NewRows([]string{"isbn", "nam", "author", "publish_year"}).
		AddRow("12235670", "Skinner", "Albert", 2001).
		AddRow("12289970", "Stlake", "Albert", 1997)
	mock.ExpectQuery(`SELECT (.*)`).WillReturnRows(expectedRows)

	book, err := repo.GetByAuthor("Albert")
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
	db, mock := NewMock()

	repo := repositories.BookRepository{
		DB:    db,
		Table: "Book",
	}
	expected := []repositories.Book{
		{ISBN: "12235670", Name: "Skinner", Author: "Albert", PublishYear: 2001},
		{ISBN: "19123450", Name: "Atomic", Author: "Grahahm", PublishYear: 2022},
	}
	sqlmock.NewRows([]string{"isbn", "name", "author", "publish_year"}).
		AddRow("19123450", "Atomic", "Grahahm", 2022).
		AddRow("12235670", "Skinner", "Albert", 2001).
		AddRow("12223900", "Short", "Victor", 1998).
		AddRow("12289970", "Stlake", "Albert", 1997)
	expectedRows := sqlmock.NewRows([]string{"isbn", "nam", "author", "publish_year"}).
		AddRow("12235670", "Skinner", "Albert", 2001).
		AddRow("19123450", "Atomic", "Grahahm", 2022)
	mock.ExpectQuery(`SELECT (.*)`).WillReturnRows(expectedRows)

	book, err := repo.GetInRange(1999, 2023)
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
