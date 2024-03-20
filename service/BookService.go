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
	for _, data := range bookData {
		existingBook, err := service.Repo.GetByISBN(data.ISBN)
		if err != nil {
			L.Error("Error: ", err)
			return err
		}
		if existingBook == (repositories.Book{}) {
			return errors.New("Book not found")
		}

		_, err2 := service.Repo.Update(data.ISBN, data.Name, data.Author, data.PublishYear)
		if err2 != nil {
			L.Error("Error: ", err2)
			return err2
		}
	}
	return nil
}

func (service BookService) Delete(bookData []repositories.Book) error {
	for _, data := range bookData {
		existingBook, err := service.Repo.GetByISBN(data.ISBN)
		if err != nil {
			L.Error("Error: ", err)
			return err
		}
		if existingBook == (repositories.Book{}) {
			return errors.New("Book not found")
		}

		_, err2 := service.Repo.Delete(data.ISBN)
		if err2 != nil {
			return err2
		}
	}
	return nil
}

func (service BookService) Insert(bookData []repositories.Book) error {
	for _, data := range bookData {
		_, err := service.Repo.GetByISBN(data.ISBN)
		if err == nil {
			L.Error("Error: ", err)
			return err
		}

		_, err2 := service.Repo.Insert(data.ISBN, data.Name, data.Author, data.PublishYear)
		if err2 != nil {
			L.Error("Error: ", err2)
			return err2
		}
	}
	return nil
}
