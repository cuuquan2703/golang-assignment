package service

import (
	"database/sql"
	"errors"
	"server/logger"
	"server/repositories"
)

type BookService struct {
	Repo *repositories.BookRepository
}

var L = logger.CreateLog()

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

func (service BookService) Update(bookData []repositories.Book) (sql.Result, error) {
	var err, err2 error
	var res sql.Result
	for _, data := range bookData {
		existingBook, errGet := service.Repo.GetByISBN(data.ISBN)
		if err != nil {
			L.Error("Error GetByISBN: ", err)
			err = errGet
		}
		if existingBook == (repositories.Book{}) {
			err = errors.New("book not found")
		}

		//
		res, err2 = service.Repo.Update(data)
		if err2 != nil {
			L.Error("Error Update: ", err2)
			err = err2
		}
	}
	return res, err
}

func (service BookService) Delete(bookData []repositories.Book) (sql.Result, error) {
	var err, err2 error
	var res sql.Result
	for _, data := range bookData {
		existingBook, errGet := service.Repo.GetByISBN(data.ISBN)
		if err != nil {
			L.Error("Error GetByISBN: ", err)
			err = errGet
		}
		if existingBook == (repositories.Book{}) {
			err = errors.New("book not found")
		}

		//
		res, err2 = service.Repo.Delete(data)
		if err2 != nil {
			L.Error("Error Update: ", err2)
			err = err2
		}
	}
	return res, err
}

func (service BookService) Insert(bookData []repositories.Book) (sql.Result, error) {
	var err, err2 error
	var res sql.Result
	for _, data := range bookData {

		//
		res, err2 = service.Repo.Insert(data)
		if err2 != nil {
			L.Error("Error Update: ", err2)
			err = err2
		}
	}
	return res, err
}
