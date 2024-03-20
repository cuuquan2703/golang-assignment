package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	_ "os"
	_ "server/logger"
	_ "strconv"

	_ "github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Author struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	BirthDate string `json:"birth_date"`
}

type AuthorRepository struct {
	DB    *sql.DB
	Table string
}

func NewAuthorRepository() (*AuthorRepository, error) {
	db, err := ConnectDB()
	if err != nil {
		return nil, err
	}
	return &AuthorRepository{
		DB:    db,
		Table: "Book",
	}, nil
}

func (repo AuthorRepository) GetAllAuthors() ([]Author, error) {
	authors := []Author{}
	cmd := `SELECT id,name,birth_date from Author`
	L.Info("Querying " + cmd)
	row, err := repo.DB.Query(cmd)
	if err != nil {
		L.Error("Error", err)
	} else {
		L.Info("Query successfully")
	}

	for row.Next() {
		author := Author{}
		err := row.Scan(&author.Id, &author.Name, &author.BirthDate)
		if err != nil {
			L.Error("Error", err)
			return nil, err
		}
		authors = append(authors, author)
	}

	if len(authors) == 0 {
		L.Error("Error ", errors.New("no author found"))
		return nil, errors.New("No author found")
	}
	defer row.Close()
	return authors, err
}

func (repo AuthorRepository) GetByID(id int) (Author, error) {
	author := Author{}
	cmd := `SELECT id,name,birth_date from Author WHERE id=$1`
	L.Info("Querying " + cmd)
	row := repo.DB.QueryRow(cmd, id)
	L.Info("Query successfully")
	err := row.Scan(&author.Id, &author.Name, &author.BirthDate)
	if err != nil {
		if err == sql.ErrNoRows {
			L.Error("Error ", errors.New("no author found"))
			return author, errors.New("No author found")
		}
		L.Error("Error ", err)
		return author, errors.New("Something went wrong")
	}
	return author, err
}

func (repo AuthorRepository) GetByName(name string) (Author, error) {
	author := Author{}
	cmd := `SELECT id,name,birth_date from Author WHERE "name"=$1`
	L.Info("Querying " + cmd)
	row := repo.DB.QueryRow(cmd, name)
	L.Info("Query successfully")
	err := row.Scan(&author.Id, &author.Name, &author.BirthDate)
	fmt.Println(author)
	if err != nil {
		if err == sql.ErrNoRows {
			L.Error("Error ", errors.New("no author found"))
			return author, errors.New("No author found")
		}
		L.Error("Error ", err)
		return author, errors.New("Something went wrong")
	}
	return author, err
}

func (repo AuthorRepository) Insert(data Author) (sql.Result, error) {
	res, err := repo.DB.Exec("INSERT INTO Author (name,birth_date) VALUES ($1,$2)", data.Name, data.BirthDate)
	if err != nil {
		L.Error("Error insert author ", err)
	} else {
		L.Info("Insert successfully author")
	}
	return res, err
}

func (repo AuthorRepository) Update(data Author) (sql.Result, error) {
	res, err := repo.DB.Exec("UPDATE Author SET name=$1, birth_date=$2 WHERE id=$3", data.Name, data.BirthDate, data.Id)
	if err != nil {
		L.Error("Error insert author ", err)
	} else {
		L.Info("Insert successfully author")
	}
	return res, err
}
