package service

import (
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

func (service BookService) Update(bookData []repositories.Book) error {
	var err error
	for _, data := range bookData {
		existingBook, errGet := service.Repo.GetByISBN(data.ISBN)
		if errGet != nil {
			L.Error("Error: ", errGet)
			err = errGet
		}
		if existingBook == (repositories.Book{}) {
			err = errors.New("Book not found")
		}

		_, errUpdate := service.Repo.Update(data.ISBN, data.Name, data.Author, data.PublishYear)
		if errUpdate != nil {
			L.Error("Error: ", errUpdate)
			err = errUpdate
		}
	}
	return err
}

func (service BookService) Delete(bookData []repositories.Book) error {
	var err error
	for _, data := range bookData {
		existingBook, errGet := service.Repo.GetByISBN(data.ISBN)
		if errGet != nil {
			L.Error("Error: ", errGet)
			err = errGet
		}
		if existingBook == (repositories.Book{}) {
			err = errors.New("Book not found")
		}

		_, err2 := service.Repo.Delete(data.ISBN)
		if err2 != nil {
			err = err2
		}
	}
	return err
}

func (service BookService) Insert(bookData []repositories.Book) error {
	var err error
	for _, data := range bookData {
		_, errGet := service.Repo.GetByISBN(data.ISBN)
		if errGet == nil {
			L.Error("Error: ", errGet)
			err = errGet
		}

		_, err2 := service.Repo.Insert(data.ISBN, data.Name, data.Author, data.PublishYear)
		if err2 != nil {
			L.Error("Error: ", err2)
			err = err2
		}
	}
	return err
}
