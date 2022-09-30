package models

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

const (
	host     = "ec2-34-247-72-29.eu-west-1.compute.amazonaws.com"
	port     = 5432
	user     = "gttcdewdwapakv"
	password = "3e932956a5c1089b1723046c7b1e05e101f8d7a5db2a9273bf40e4d762931d6d"
	dbname   = "d2lvv2j2o59m80"
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
