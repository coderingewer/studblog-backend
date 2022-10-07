package models

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

/*
const (
	host     = "ec2-99-81-16-126.eu-west-1.compute.amazonaws.com"
	port     = 5432
	user     = "ywmpggsikbsbqm"
	password = "dc9395dcf455a164c32173b40c82ca99835ae28d371b7977b9b122c272be51bc"
	dbname   = "d6anbpmcouj6s2"
)*/

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "uyumak"
	dbname   = "blogapp"
)

var db *gorm.DB

func init() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	fmt.Println(psqlInfo)
	conn, err := gorm.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	db = conn
	db.Debug().AutoMigrate(User{}, Post{},
		Image{}, Tag{}, PostTag{}, Like{}, View{}, Favorite{}, FavoritesList{})
	fmt.Println("Veri Tabanı bağlantısı başarılı!")
}
func GetDB() *gorm.DB {
	return db
}
