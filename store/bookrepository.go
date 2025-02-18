package store

import (
	"fmt"
	"log"

	"github.com/polyakovaa/standartserver3/internal/app/models"
)

type BookRepository struct {
	store *Store
}

var (
	tableBook string = "books"
)

func (br *BookRepository) Create(b *models.Book) (*models.Book, error) {
	query := fmt.Sprintf("INSERT INTO %s (title, author, content) VALUES ($1, $2, $3) RETURNING id", tableBook)
	if err := br.store.db.QueryRow(query, b.Title, b.Author, b.Content).Scan(&b.ID); err != nil {
		return nil, err
	}
	return b, nil
}

func (br *BookRepository) DeleteById(id int) (*models.Book, error) {
	book, ok, err := br.FindBookById(id)
	if err != nil {
		return nil, err
	}
	if ok {
		query := fmt.Sprintf("delete from %s where id=$1", tableBook)
		_, err = br.store.db.Exec(query, id)
		if err != nil {
			return nil, err
		}
	}

	return book, nil
}

func (br *BookRepository) FindBookById(id int) (*models.Book, bool, error) {
	books, err := br.SelectAll()
	found := false
	if err != nil {
		return nil, found, err
	}
	var bookFound *models.Book
	for _, b := range books {
		if b.ID == id {
			bookFound = b
			found = true
		}
	}

	return bookFound, found, nil

}

func (br *BookRepository) SelectAll() ([]*models.Book, error) {
	query := fmt.Sprintf("SELECT * FROM %s", tableBook)
	rows, err := br.store.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	books := make([]*models.Book, 0)
	for rows.Next() {
		b := models.Book{}
		err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.Content)
		if err != nil {
			log.Println(err)
			continue
		}
		books = append(books, &b)
	}
	return books, nil
}
