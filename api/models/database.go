package models

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

const (
	host     = "ec2-52-212-228-71.eu-west-1.compute.amazonaws.com"
	port     = 5432
	user     = "rfraosvrrrjwhz"
	password = "3909ea924b9f57d3225430e04b6bba39027db1bb3d88e4c9e6f8b880572dd984"
	dbname   = "d8jmg90h8phnrj"
)

var db *gorm.DB

func init() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=require",
		host, port, user, password, dbname)

	fmt.Println(psqlInfo)
	conn, err := gorm.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	db = conn
	db.Debug().AutoMigrate(User{}, Post{},
		Image{}, Tag{}, PostTag{}, Like{}, View{})
	fmt.Println("DB Successfully connected!")
}
func GetDB() *gorm.DB {
	return db
}
