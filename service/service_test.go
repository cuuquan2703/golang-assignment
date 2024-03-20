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
	AuthorRepo: repositories.AuthorRepository{
		DB:    db,
		Table: "author",
	},
}

var bookService = service.BookService{
	Repo: repo,
}

func TestGetAllBooks(t *testing.T) {
	expected := []repositories.Book{
		{ISBN: "19123450", Name: "Atomic", Author: repositories.Author{Id: 1, Name: "Thmoas", BirthDate: "17-04-2002"}, PublishYear: 2022},
		{ISBN: "12235670", Name: "Skinner", Author: repositories.Author{Id: 2, Name: "Albert", BirthDate: "17-04-2002"}, PublishYear: 2001},
		{ISBN: "12223900", Name: "Short", Author: repositories.Author{Id: 3, Name: "Vicotr", BirthDate: "17-04-2002"}, PublishYear: 1998},
	}
	Bookrows := sqlmock.NewRows([]string{"isbn", "name", "author", "publish_year"}).
		AddRow("19123450", "Atomic", 1, 2022).
		AddRow("12235670", "Skinner", 2, 2001).
		AddRow("12223900", "Short", 3, 1998)
	Authorrows := sqlmock.NewRows([]string{"id", "name", "birth_date"}).
		AddRow(1, "Thmoas", "17-04-2002").
		AddRow(2, "Albert", "17-04-2002").
		AddRow(3, "Vicotr", "17-04-2002")
	mock.ExpectQuery("SELECT isbn,name,author,publish_year from Book").WillReturnRows(Bookrows)
	mock.ExpectQuery("SELECT id,name,birth_date from Author").WillReturnRows(Authorrows)
	mock.ExpectQuery("SELECT (.*)").WillReturnRows(Authorrows)
	mock.ExpectQuery("SELECT (.*)").WillReturnRows(Authorrows)

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
		{ISBN: "12223900", Name: "Short", Author: repositories.Author{Id: 3, Name: "Vicotr", BirthDate: "17-04-2002"}, PublishYear: 1998},
	}
	mock.ExpectQuery(regexp.QuoteMeta("SELECT isbn,name,id_author,publish_year from Book where isbn=$1")).WithArgs("12223900").WillReturnRows(sqlmock.NewRows([]string{"isbn", "name", "author", "publish_year"}).
		AddRow("12223900", "Short", 3, 1998))
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id,name,birth_date from Author WHERE id=$1")).WithArgs(3).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "birth_date"}).
		AddRow(3, "Vicotr", "17-04-2002"))

	book, err := bookService.GetByISBN("12223900")
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
		{ISBN: "12235670", Name: "Skinner", Author: repositories.Author{Id: 2, Name: "Albert", BirthDate: "17-04-2002"}, PublishYear: 2001},
	}
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id,name,birth_date from Author WHERE "name"=$1`)).WithArgs("Albert").WillReturnRows(sqlmock.NewRows([]string{"id", "name", "birth_date"}).
		AddRow(2, "Albert", "17-04-2002"))
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT isbn,name,id_author,publish_year from Book where "id_author"=$1`)).WithArgs(2).WillReturnRows(sqlmock.NewRows([]string{"isbn", "name", "author", "publish_year"}).
		AddRow("12235670", "Skinner", 2, 2001))

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
		{ISBN: "19123450", Name: "Atomic", Author: repositories.Author{Id: 1, Name: "Thmoas", BirthDate: "17-04-2002"}, PublishYear: 2022},
		{ISBN: "12235670", Name: "Skinner", Author: repositories.Author{Id: 2, Name: "Albert", BirthDate: "17-04-2002"}, PublishYear: 2001},
	}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT isbn,name,id_author,publish_year from Book where "publish_year"<=$2 and "publish_year">=$1`)).WithArgs(1999, 2023).WillReturnRows(sqlmock.NewRows([]string{"isbn", "name", "author", "publish_year"}).
		AddRow("19123450", "Atomic", 1, 2022).
		AddRow("12235670", "Skinner", 2, 2001))
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id,name,birth_date from Author WHERE id=$1`)).WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "birth_date"}).
		AddRow(1, "Thmoas", "17-04-2002"))
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id,name,birth_date from Author WHERE id=$1`)).WithArgs(2).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "birth_date"}).
		AddRow(2, "Albert", "17-04-2002"))

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

	newBook := repositories.Book{
		ISBN:        "123456789",
		Name:        "Updated Book Name",
		PublishYear: 2024,
		Author:      repositories.Author{},
	}

	mock.ExpectExec(regexp.QuoteMeta("UPDATE Book SET name = $1, publish_year = $2 WHERE isbn = $3")).
		WithArgs(newBook.Name, newBook.PublishYear, newBook.ISBN).
		WillReturnResult((sqlmock.NewResult(0, 1)))

	_, err := bookService.Update(newBook)

	if err != nil {
		t.Errorf("Error when updating data")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestDelete(t *testing.T) {

	newBook := repositories.Book{
		ISBN:        "123456789",
		Name:        "Name",
		PublishYear: 2024,
		Author:      repositories.Author{},
	}

	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM Book WHERE isbn = $1")).
		WithArgs(newBook.ISBN).
		WillReturnResult((sqlmock.NewResult(0, 1)))

	_, err := bookService.Delete(newBook)

	if err != nil {
		t.Errorf("Error when deleting data")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestInsertCaseExistAuthor(t *testing.T) {

	existAuthor := repositories.Author{
		Id:        1,
		Name:      "Author",
		BirthDate: "17-04-2002",
	}

	newBook := repositories.Book{
		ISBN:        "123456789",
		Name:        "Name",
		PublishYear: 2024,
		Author:      existAuthor,
	}
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id,name,birth_date from Author WHERE "name"=$1`)).
		WithArgs(existAuthor.Name).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "birth_date"}).
			AddRow(1, "Author", "17-04-2002"))

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO Book (isbn, name, publish_year, id_author) VALUES ($1, $2, $3, $4);")).
		WithArgs(newBook.ISBN, newBook.Name, newBook.PublishYear, existAuthor.Id).
		WillReturnResult((sqlmock.NewResult(0, 1)))
	_, err := bookService.Insert(newBook)

	if err != nil {
		t.Errorf("Error when Insert data")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestInsertCaseNotExistAuthor(t *testing.T) {

	notExistAuthor := repositories.Author{
		Id:        1,
		Name:      "Author",
		BirthDate: "17-04-2002",
	}

	newBook := repositories.Book{
		ISBN:        "123456789",
		Name:        "Name",
		PublishYear: 2024,
		Author:      notExistAuthor,
	}
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id,name,birth_date from Author WHERE "name"=$1`)).
		WithArgs(notExistAuthor.Name).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "birth_date"}).
			AddRow(nil, nil, nil))

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO Author (name,birth_date) VALUES ($1,$2)")).
		WithArgs(notExistAuthor.Name, notExistAuthor.BirthDate).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id,name,birth_date from Author WHERE "name"=$1`)).
		WithArgs(notExistAuthor.Name).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "birth_date"}).
			AddRow(1, "Author", "17-04-2002"))

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO Book (isbn, name, publish_year, id_author) VALUES ($1, $2, $3, $4);")).
		WithArgs(newBook.ISBN, newBook.Name, newBook.PublishYear, notExistAuthor.Id).
		WillReturnResult((sqlmock.NewResult(0, 1)))
	_, err := bookService.Insert(newBook)

	if err != nil {
		t.Errorf("Error when Insert data")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}
