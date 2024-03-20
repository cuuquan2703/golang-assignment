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
	BookAuthorRepo: repositories.BookAuthorRepository{
		DB:    db,
		Table: "book_author",
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

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT b.isbn,b.name,ba.id_author,b.publish_year from Book b JOIN book_author ba ON b.isbn = ba.id_book`)).WillReturnRows(Bookrows)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id,name,birth_date from Author WHERE id=$1`)).WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "birth_date"}).
		AddRow(1, "Thmoas", "17-04-2002"))
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id,name,birth_date from Author WHERE id=$1`)).WithArgs(2).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "birth_date"}).
		AddRow(2, "Albert", "17-04-2002"))
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id,name,birth_date from Author WHERE id=$1`)).WithArgs(3).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "birth_date"}).
		AddRow(3, "Vicotr", "17-04-2002"))

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

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT b.isbn,b.name,ba.id_author,b.publish_year from Book b 
	JOIN book_author ba ON b.isbn = ba.id_book 
	WHERE b.isbn = $1`)).WithArgs("12223900").WillReturnRows(sqlmock.NewRows([]string{"isbn", "name", "author", "publish_year"}).
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
		{ISBN: "12223900", Name: "Short", Author: repositories.Author{Id: 3, Name: "Vicotr", BirthDate: "17-04-2002"}, PublishYear: 1998},
	}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT b.isbn,b.name,ba.id_author,b.publish_year from Book b 
	JOIN book_author ba ON b.isbn= ba.id_book 
	JOIN author a ON ba.id_author = a.id 
	WHERE a.name = $1`)).WithArgs("Vicotr").WillReturnRows(sqlmock.NewRows([]string{"isbn", "name", "author", "publish_year"}).
		AddRow("12223900", "Short", 3, 1998))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id,name,birth_date from Author WHERE id=$1")).WithArgs(3).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "birth_date"}).
		AddRow(3, "Vicotr", "17-04-2002"))
	book, err := bookService.GetByAuthor("Vicotr")
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

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT b.isbn,b.name,ba.id_author,b.publish_year from Book b 
	JOIN book_author ba ON b.isbn = ba.id_book  
	WHERE b.publish_year<=$2 and b.publish_year>=$1`)).WithArgs(1999, 2023).WillReturnRows(sqlmock.NewRows([]string{"isbn", "name", "author", "publish_year"}).
		AddRow("19123450", "Atomic", 1, 2022).
		AddRow("12235670", "Skinner", 2, 2001))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id,name,birth_date from Author WHERE id=$1")).WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "birth_date"}).
		AddRow(1, "Thmoas", "17-04-2002"))
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id,name,birth_date from Author WHERE id=$1")).WithArgs(2).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "birth_date"}).
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

	oldData := []repositories.Book{
		{ISBN: "123456789", Name: "Old Book Name", PublishYear: 2022, Author: repositories.Author{
			Id:        1,
			Name:      "Albert",
			BirthDate: "12-02-1998",
		}},
		{ISBN: "113456789", Name: "Old Book Name", PublishYear: 2024, Author: repositories.Author{
			Id:        2,
			Name:      "Kale",
			BirthDate: "12-02-1988",
		}},
	}

	newData := []repositories.Book{
		{ISBN: "123456789", Name: "Updated Book Name", Author: repositories.Author{
			Id:        1,
			Name:      "Albert",
			BirthDate: "12-02-1998"}, PublishYear: 2024},
		{ISBN: "113456789", Name: "Updated Book Name", Author: repositories.Author{
			Id:        2,
			Name:      "Kale",
			BirthDate: "12-02-1988",
		}, PublishYear: 2024},
	}

	for index, _ := range oldData {

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT b.isbn,b.name,ba.id_author,b.publish_year from Book b 
		JOIN book_author ba ON b.isbn = ba.id_book 
		WHERE b.isbn = $1`)).
			WithArgs(oldData[index].ISBN).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "publish_year", "author"}).
				AddRow(oldData[index].ISBN, oldData[index].Name, oldData[index].Author.Id, oldData[index].PublishYear))

		mock.ExpectExec(regexp.QuoteMeta("UPDATE Book SET name = $1, publish_year = $2 WHERE isbn = $3")).
			WithArgs(newData[index].Name, newData[index].PublishYear, newData[index].ISBN).
			WillReturnResult((sqlmock.NewResult(0, 1)))
	}

	_, err := bookService.Update(newData)

	if err != nil {
		t.Errorf("Error when updating data")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestDelete(t *testing.T) {

	data := []repositories.Book{
		{ISBN: "123456789", Name: "Old Book Name", Author: repositories.Author{
			Id:        1,
			Name:      "Albert",
			BirthDate: "12-02-1998",
		}, PublishYear: 2022},
		{ISBN: "113456789", Name: "Old Book Name", Author: repositories.Author{
			Id:        2,
			Name:      "Kale",
			BirthDate: "12-02-1988",
		}, PublishYear: 2024},
	}
	for index, _ := range data {

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT b.isbn,b.name,ba.id_author,b.publish_year from Book b 
		JOIN book_author ba ON b.isbn = ba.id_book 
		WHERE b.isbn = $1`)).
			WithArgs(data[index].ISBN).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "publish_year", "author"}).
				AddRow(data[index].ISBN, data[index].Name, data[index].Author.Id, data[index].PublishYear))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT id,name,birth_date from Author WHERE id=$1`)).
			WithArgs(data[index].Author.Id).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "birth_dat"}).
				AddRow(data[index].Author.Id, data[index].Author.Name, data[index].Author.BirthDate))

		mock.ExpectExec(regexp.QuoteMeta(" DELETE FROM Book_Author WHERE id_book = $1;")).
			WithArgs(data[index].ISBN).
			WillReturnResult((sqlmock.NewResult(0, 1)))

		mock.ExpectExec(regexp.QuoteMeta("DELETE FROM Book WHERE isbn = $1")).
			WithArgs(data[index].ISBN).
			WillReturnResult((sqlmock.NewResult(0, 1)))
	}

	_, err := bookService.Delete(data)

	if err != nil {
		t.Errorf("Error when deleting data")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestInsertCaseExistAuthor(t *testing.T) {

	existAuthor := []repositories.Author{
		{Id: 1, Name: "Albert", BirthDate: "12-02-1998"},
		{Id: 2, Name: "Kale", BirthDate: "12-02-1988"},
	}

	data := []repositories.Book{
		{ISBN: "123456789", Name: "New Book Name", PublishYear: 2022, Author: repositories.Author{
			Name:      "Albert",
			BirthDate: "12-02-1998",
		}},
		{ISBN: "113456789", Name: "New Book Name", PublishYear: 2024, Author: repositories.Author{
			Name:      "Kale",
			BirthDate: "12-02-1988",
		}},
	}

	for index, _ := range data {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT id,name,birth_date from Author WHERE "name"=$1`)).
			WithArgs(data[index].Author.Name).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "birth_date"}).
				AddRow(existAuthor[index].Id, existAuthor[index].Name, existAuthor[index].BirthDate))

		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO Book (isbn, name, publish_year) VALUES ($1, $2, $3);")).
			WithArgs(data[index].ISBN, data[index].Name, data[index].PublishYear).
			WillReturnResult((sqlmock.NewResult(0, 1)))

		mock.ExpectExec(regexp.QuoteMeta(" INSERT INTO Book_Author (id_book,id_author) VALUES ($1, $2);")).
			WithArgs(data[index].ISBN, existAuthor[index].Id).
			WillReturnResult((sqlmock.NewResult(0, 1)))
	}
	_, err := bookService.Insert(data)

	if err != nil {
		t.Errorf("Error when Insert data")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestInsertCaseNotExistAuthor(t *testing.T) {

	genAuthor := []repositories.Author{
		{Id: 1, Name: "Albert", BirthDate: "12-02-1998"},
		{Id: 2, Name: "Kale", BirthDate: "12-02-1988"},
	}

	data := []repositories.Book{
		{ISBN: "123456789", Name: "New Book Name", PublishYear: 2022, Author: repositories.Author{
			Name:      "Albert",
			BirthDate: "12-02-1998",
		}},
		{ISBN: "113456789", Name: "New Book Name", PublishYear: 2024, Author: repositories.Author{
			Name:      "Kale",
			BirthDate: "12-02-1988",
		}},
	}

	for index, _ := range data {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT id,name,birth_date from Author WHERE "name"=$1`)).
			WithArgs(data[index].Author.Name).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "birth_date"}).
				AddRow(nil, nil, nil))

		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO Author (name,birth_date) VALUES ($1,$2)")).
			WithArgs(data[index].Author.Name, data[index].Author.BirthDate).
			WillReturnResult((sqlmock.NewResult(0, 1)))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT id,name,birth_date from Author WHERE "name"=$1`)).
			WithArgs(data[index].Author.Name).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "birth_date"}).
				AddRow(genAuthor[index].Id, genAuthor[index].Name, genAuthor[index].BirthDate))

		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO Book (isbn, name, publish_year) VALUES ($1, $2, $3);")).
			WithArgs(data[index].ISBN, data[index].Name, data[index].PublishYear).
			WillReturnResult((sqlmock.NewResult(0, 1)))

		mock.ExpectExec(regexp.QuoteMeta(" INSERT INTO Book_Author (id_book,id_author) VALUES ($1, $2);")).
			WithArgs(data[index].ISBN, genAuthor[index].Id).
			WillReturnResult((sqlmock.NewResult(0, 1)))
	}
	_, err := bookService.Insert(data)

	if err != nil {
		t.Errorf("Error when Insert data")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}
