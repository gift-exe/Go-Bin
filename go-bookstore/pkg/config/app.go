package config

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	db *gorm.DB //creeating a db variable of type gorm.DB
)

func Connect() *gorm.DB {
	db, err := gorm.Open("mysql", "root:AwesomeGod003#@!@/go-bookstore-db?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}

	return db
}
