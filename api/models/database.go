package models

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

const (
	host     = "dpg-cedm5e02i3mr7lhduv20-a"
	port     = 5432
	user     = "blogapp_user"
	password = "Ts9mtd9oFR9F8dWVEsEyT7fByFYngFLT"
	dbname   = "blogapp"
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
		Image{}, Tag{}, PostTag{}, Like{}, View{}, Favorite{}, FavoritesList{}, Comment{})
	fmt.Println("Veri Tabanı bağlantısı başarılı!")
}
func GetDB() *gorm.DB {
	return db
}
