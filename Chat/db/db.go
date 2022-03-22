package db

import (
	"fmt"
	"goChallenge/chat/config"
	"goChallenge/chat/models"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Database struct {
	*gorm.DB
}

var DB *gorm.DB

// Opening a database and save the reference to `Database` struct.
func Init() *gorm.DB {
	dbConfig := config.DbConfig

	dsn := fmt.Sprintf("%s:%s@%s(%s:%s)/%s?parseTime=True",
		dbConfig.User,
		dbConfig.Pwd,
		"tcp",
		dbConfig.Host,
		"3306",
		dbConfig.Db)

	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		log.Fatal("db err: (Init) ", err)
	}
	db.AutoMigrate(&models.User{}, &models.Message{})

	DB = db

	return db
}

func Close() {
	DB.Close()
}
