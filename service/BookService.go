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

func (service BookService) Update(data repositories.Book) (sql.Result, error) {
	return service.Repo.Update(data)
}

func (service BookService) Delete(data repositories.Book) (sql.Result, error) {
	return service.Repo.Delete(data)
}

func (service BookService) Insert(data repositories.Book) (sql.Result, error) {
	return service.Repo.Insert(data)
}
