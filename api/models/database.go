package models

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "uyumak"
	dbname   = "gulikk"
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
		Image{}, Tag{}, PostTag{}, Like{}, View{})

	fmt.Println("DB Successfully connected!")
	img := Image{}
	usrimg, _ := img.SaveImage()
	user := User{Name: "blogAdmin", Username: "superruser", Email: "usersuperr@super.com", UserRole: "SUPER-USER", ImageID: usrimg.ID, Password: "superUser.2"}
	user.BeforeSAve()
	_, _ = user.SaveUser()
	fmt.Println("Seed user sent")
}

func GetDB() *gorm.DB {
	return db
}

func seed() {

}
