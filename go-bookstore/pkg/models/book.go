package models

import (
	"github.com/Gift-py/go-bookstore/pkg/config"
	"github.com/jinzhu/gorm"
)

var db *gorm.DB //kinda like a gloal variable

type Book struct {
	gorm.Model
	Name        string `json:"name"`
	Author      string `json:"author"`
	Publication string `json:"publication"`
}

func init() {
	db = config.Connect()
	db.AutoMigrate(&Book{})
}

func (b *Book) CreateBook() *Book {
	//similar to an object method.
	//in this case this function is a method of the struct book
	db.NewRecord(b)
	db.Create(&b)
	return b
}

func GetAllBooks() []Book {
	var Books []Book
	db.Find(&Books)
	return Books
}

func GetBookById(Id int64) (*Book, *gorm.DB) {
	var getBook Book
	db := db.Where("ID=?", Id).Find(&getBook) //find the book with "Where()", store value in getBook with "Find()"
	return &getBook, db
}

func DeleteBook(Id int64) Book {
	var book Book
	db.Where("ID=?", Id).Delete(book)
	return book
}
