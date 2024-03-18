package repositories_test

import (
	"database/sql"
	"log"
	"reflect"
	"regexp"
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
		AuthorRepo: repositories.AuthorRepository{
			DB:    db,
			Table: "author",
		},
		BookAuthorRepo: repositories.BookAuthorRepository{
			DB:    db,
			Table: "book_author",
		},
	}
	expected := []repositories.Book{
		{ISBN: "19123450", Name: "Atomic", Author: repositories.Author{Id: 1, Name: "Thmoas", BirthDate: "17-04-2002"}, PublishYear: 2022},
		{ISBN: "12235670", Name: "Skinner", Author: repositories.Author{Id: 2, Name: "Albert", BirthDate: "17-04-2002"}, PublishYear: 2001},
		{ISBN: "12223900", Name: "Short", Author: repositories.Author{Id: 3, Name: "Vicotr", BirthDate: "17-04-2002"}, PublishYear: 1998},
	}
	Bookrows := sqlmock.NewRows([]string{"isbn", "name", "author", "publish_year"}).
		AddRow("19123450", "Atomic", 1, 2022).
		AddRow("12235670", "Skinner", 2, 2001).
		AddRow("12223900", "Short", 3, 1998)
	sqlmock.NewRows([]string{"id", "name", "birth_date"}).
		AddRow(1, "Thmoas", "17-04-2002").
		AddRow(2, "Albert", "17-04-2002").
		AddRow(3, "Vicotr", "17-04-2002")
	sqlmock.NewRows([]string{"id_book", "id_author"}).
		AddRow("19123450", 1).AddRow("12235670", 2).AddRow("12223900", 3)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT b.isbn,b.name,ba.id_author,b.publish_year from Book b JOIN book_author ba ON b.isbn = ba.id_book`)).WillReturnRows(Bookrows)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id,name,birth_date from Author WHERE id=$1`)).WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "birth_date"}).
		AddRow(1, "Thmoas", "17-04-2002"))
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id,name,birth_date from Author WHERE id=$1`)).WithArgs(2).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "birth_date"}).
		AddRow(2, "Albert", "17-04-2002"))
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id,name,birth_date from Author WHERE id=$1`)).WithArgs(3).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "birth_date"}).
		AddRow(3, "Vicotr", "17-04-2002"))

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
		AuthorRepo: repositories.AuthorRepository{
			DB:    db,
			Table: "author",
		},
		BookAuthorRepo: repositories.BookAuthorRepository{
			DB:    db,
			Table: "book_author",
		},
	}
	expected := []repositories.Book{
		{ISBN: "12223900", Name: "Short", Author: repositories.Author{Id: 3, Name: "Vicotr", BirthDate: "17-04-2002"}, PublishYear: 1998},
	}
	sqlmock.NewRows([]string{"isbn", "name", "author", "publish_year"}).
		AddRow("19123450", "Atomic", 1, 2022).
		AddRow("12235670", "Skinner", 2, 2001).
		AddRow("12223900", "Short", 3, 1998)
	sqlmock.NewRows([]string{"id", "name", "birth_date"}).
		AddRow(1, "Thmoas", "17-04-2002").
		AddRow(2, "Albert", "17-04-2002").
		AddRow(3, "Vicotr", "17-04-2002")
	sqlmock.NewRows([]string{"id_book", "id_author"}).
		AddRow("19123450", 1).AddRow("12235670", 2).AddRow("12223900", 3)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT b.isbn,b.name,ba.id_author,b.publish_year from Book b 
	JOIN book_author ba ON b.isbn = ba.id_book 
	WHERE b.isbn = $1`)).WithArgs("12223900").WillReturnRows(sqlmock.NewRows([]string{"isbn", "name", "author", "publish_year"}).
		AddRow("12223900", "Short", 3, 1998))
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id,name,birth_date from Author WHERE id=$1")).WithArgs(3).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "birth_date"}).
		AddRow(3, "Vicotr", "17-04-2002"))

	book, err := repo.GetByISBN("12223900")
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
		AuthorRepo: repositories.AuthorRepository{
			DB:    db,
			Table: "author",
		},
		BookAuthorRepo: repositories.BookAuthorRepository{
			DB:    db,
			Table: "book_author",
		},
	}
	expected := []repositories.Book{
		{ISBN: "12223900", Name: "Short", Author: repositories.Author{Id: 3, Name: "Vicotr", BirthDate: "17-04-2002"}, PublishYear: 1998},
	}
	sqlmock.NewRows([]string{"isbn", "name", "author", "publish_year"}).
		AddRow("19123450", "Atomic", 1, 2022).
		AddRow("12235670", "Skinner", 2, 2001).
		AddRow("12223900", "Short", 3, 1998)
	sqlmock.NewRows([]string{"id", "name", "birth_date"}).
		AddRow(1, "Thmoas", "17-04-2002").
		AddRow(2, "Albert", "17-04-2002").
		AddRow(3, "Vicotr", "17-04-2002")
	sqlmock.NewRows([]string{"id_book", "id_author"}).
		AddRow("19123450", 1).AddRow("12235670", 2).AddRow("12223900", 3)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT b.isbn,b.name,ba.id_author,b.publish_year from Book b 
	JOIN book_author ba ON b.isbn= ba.id_book 
	JOIN author a ON ba.id_author = a.id 
	WHERE a.name = $1`)).WithArgs("Vicotr").WillReturnRows(sqlmock.NewRows([]string{"isbn", "name", "author", "publish_year"}).
		AddRow("12223900", "Short", 3, 1998))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id,name,birth_date from Author WHERE id=$1")).WithArgs(3).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "birth_date"}).
		AddRow(3, "Vicotr", "17-04-2002"))
	book, err := repo.GetByAuthor("Vicotr")
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
		AuthorRepo: repositories.AuthorRepository{
			DB:    db,
			Table: "author",
		},
		BookAuthorRepo: repositories.BookAuthorRepository{
			DB:    db,
			Table: "book_author",
		},
	}
	expected := []repositories.Book{
		{ISBN: "19123450", Name: "Atomic", Author: repositories.Author{Id: 1, Name: "Thmoas", BirthDate: "17-04-2002"}, PublishYear: 2022},
		{ISBN: "12235670", Name: "Skinner", Author: repositories.Author{Id: 2, Name: "Albert", BirthDate: "17-04-2002"}, PublishYear: 2001},
	}
	sqlmock.NewRows([]string{"isbn", "name", "author", "publish_year"}).
		AddRow("19123450", "Atomic", 1, 2022).
		AddRow("12235670", "Skinner", 2, 2001).
		AddRow("12223900", "Short", 3, 1998)
	sqlmock.NewRows([]string{"id", "name", "birth_date"}).
		AddRow(1, "Thmoas", "17-04-2002").
		AddRow(2, "Albert", "17-04-2002").
		AddRow(3, "Vicotr", "17-04-2002")
	sqlmock.NewRows([]string{"id_book", "id_author"}).
		AddRow("19123450", 1).AddRow("12235670", 2).AddRow("12223900", 3)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT b.isbn,b.name,ba.id_author,b.publish_year from Book b 
	JOIN book_author ba ON b.isbn = ba.id_book  
	WHERE b.publish_year<=$2 and b.publish_year>=$1`)).WithArgs(1999, 2023).WillReturnRows(sqlmock.NewRows([]string{"isbn", "name", "author", "publish_year"}).
		AddRow("19123450", "Atomic", 1, 2022).
		AddRow("12235670", "Skinner", 2, 2001))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id,name,birth_date from Author WHERE id=$1")).WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "birth_date"}).
		AddRow(1, "Thmoas", "17-04-2002"))
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id,name,birth_date from Author WHERE id=$1")).WithArgs(2).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "birth_date"}).
		AddRow(2, "Albert", "17-04-2002"))

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

// func TestFailGetByISBN(t *testing.T) {
// 	db, mock := NewMock()

// 	repo := repositories.BookRepository{
// 		DB:    db,
// 		Table: "Book",
// 	}
// 	expected := "No books found"
// 	sqlmock.NewRows([]string{"isbn", "name", "author", "publish_year"}).
// 		AddRow("19123450", "Atomic", "Grahahm", 2022).
// 		AddRow("12235670", "Skinner", "Albert", 2001).
// 		AddRow("12223900", "Short", "Victor", 1998)
// 	expectedRows := sqlmock.NewRows([]string{"isbn", "nam", "author", "publish_year"})
// 	mock.ExpectQuery(`SELECT (.*)`).WillReturnRows(expectedRows)

// 	book, err := repo.GetByISBN("12235671")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	if !reflect.DeepEqual(err, expected) {
// 		t.Errorf("Returned books don't match expected books. Expected: %v, Actual: %v", expected, book)
// 	}
// 	// if err := mock.ExpectationsWereMet(); err != nil {
// 	// 	t.Errorf("Unfulfilled expectations: %s", err)
// 	// }
// }

// func TestGetByAuthor(t *testing.T) {
// 	db, mock := NewMock()

// 	repo := repositories.BookRepository{
// 		DB:    db,
// 		Table: "Book",
// 	}
// 	expected := []repositories.Book{
// 		{ISBN: "12235670", Name: "Skinner", Author: "Albert", PublishYear: 2001},
// 		{ISBN: "12289970", Name: "Stlake", Author: "Albert", PublishYear: 1997},
// 	}
// 	sqlmock.NewRows([]string{"isbn", "name", "author", "publish_year"}).
// 		AddRow("19123450", "Atomic", "Grahahm", 2022).
// 		AddRow("12235670", "Skinner", "Albert", 2001).
// 		AddRow("12223900", "Short", "Victor", 1998).
// 		AddRow("12289970", "Stlake", "Albert", 1997)
// 	expectedRows := sqlmock.NewRows([]string{"isbn", "nam", "author", "publish_year"}).
// 		AddRow("12235670", "Skinner", "Albert", 2001).
// 		AddRow("12289970", "Stlake", "Albert", 1997)
// 	mock.ExpectQuery(`SELECT (.*)`).WillReturnRows(expectedRows)

// 	book, err := repo.GetByAuthor("Albert")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	if !reflect.DeepEqual(book, expected) {
// 		t.Errorf("Returned books don't match expected books. Expected: %v, Actual: %v", expected, book)
// 	}
// 	if err := mock.ExpectationsWereMet(); err != nil {
// 		t.Errorf("Unfulfilled expectations: %s", err)
// 	}
// }

// func TestGetInRange(t *testing.T) {
// 	db, mock := NewMock()

// 	repo := repositories.BookRepository{
// 		DB:    db,
// 		Table: "Book",
// 	}
// 	expected := []repositories.Book{
// 		{ISBN: "12235670", Name: "Skinner", Author: "Albert", PublishYear: 2001},
// 		{ISBN: "19123450", Name: "Atomic", Author: "Grahahm", PublishYear: 2022},
// 	}
// 	sqlmock.NewRows([]string{"isbn", "name", "author", "publish_year"}).
// 		AddRow("19123450", "Atomic", "Grahahm", 2022).
// 		AddRow("12235670", "Skinner", "Albert", 2001).
// 		AddRow("12223900", "Short", "Victor", 1998).
// 		AddRow("12289970", "Stlake", "Albert", 1997)
// 	expectedRows := sqlmock.NewRows([]string{"isbn", "nam", "author", "publish_year"}).
// 		AddRow("12235670", "Skinner", "Albert", 2001).
// 		AddRow("19123450", "Atomic", "Grahahm", 2022)
// 	mock.ExpectQuery(`SELECT (.*)`).WillReturnRows(expectedRows)

// 	book, err := repo.GetInRange(1999, 2023)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	if !reflect.DeepEqual(book, expected) {
// 		t.Errorf("Returned books don't match expected books. Expected: %v, Actual: %v", expected, book)
// 	}
// 	if err := mock.ExpectationsWereMet(); err != nil {
// 		t.Errorf("Unfulfilled expectations: %s", err)
// 	}
// }
