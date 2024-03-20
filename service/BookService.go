package service

import (
	"database/sql"
	"server/repositories"
)

type BookService struct {
	Repo *repositories.BookRepository
}

func (service BookService) GetAllBooks() ([]repositories.Book, error) {
	return service.Repo.GetAllBooks()
}

func (service BookService) GetByISBN(isbn string) (repositories.Book, error) {
	return service.Repo.GetByISBN(isbn)
}

func (service BookService) GetByAuthor(author string) ([]repositories.Book, error) {
	return service.Repo.GetByAuthor(author)
}

func (service BookService) GetInRange(year1, year2 int) ([]repositories.Book, error) {
	return service.Repo.GetInRange(year1, year2)
}

func (service BookService) Update(isbn, name, author string, publish_year int) (sql.Result, error) {
	return service.Repo.Update(isbn, name, author, publish_year)
}

func (service BookService) Delete(isbn string) (sql.Result, error) {
	return service.Repo.Delete(isbn)
}

func (service BookService) Insert(isbn, name, author string, publish_year int) (sql.Result, error) {
	return service.Repo.Insert(isbn, name, author, publish_year)
}
