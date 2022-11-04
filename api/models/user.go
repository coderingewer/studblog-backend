package models

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model
	Username string `gorm:"size:255;not null; unique" json:"username"`
	Email    string `gorm:"size:255;not null;unique" json:"email"`
	Name     string `json:"name"`
	UserRole string `gorm:"size:20;not null;" json:"userRole"`
	Password string `gorm:"size:255;not null;" json:"password"`
	Isvalid  bool   `gorm:"not null;" json:"isvalid"`
	ImageID  uint   `json:"imageId"`
	Image    Image  `json:"image"`
	Posts    []Post `json:"posts"`
}

func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerifyPassword(hashesPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashesPassword), []byte(password))
}

func (u *User) BeforeSAve() error {
	hashedPassword, err := Hash(u.Password)
	if err != nil {
		fmt.Println(string(err.Error()))
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) Prepare() {
	u.ID = 0
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	u.DeletedAt = nil
	u.Isvalid = false
	u.UserRole = "KULLANICI"
	u.Image.ID = u.ImageID
	u.Email = strings.ToLower(u.Email)
	u.Image.Url = "https://res.cloudinary.com/ddeatrwxs/image/upload/v1661888774/studapp/spwdg9relpxlmozdck6b.png"
}

func (u *User) Validate(action string) error {
	switch strings.ToLower(action) {
	case "login":
		if u.Password == "" {
			return errors.New("Şifre Zorunlu")
		}
		if u.Email == "" {
			return errors.New("Email Zorunlu")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Doğrulanmamış E Mail Adresi")
		}
		return nil

	default:
		if u.Username == "" {
			return errors.New("Kullanıcı Adı Zorulu")
		}
		if u.Password == "" {
			return errors.New("Şifre Zorulu")
		}
		if u.Username == "" {
			return errors.New("E Posta Adresi Zorulu")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Onaylanmamış E Mail Adresi")
		}
		return nil
	}
}

func (u *User) SaveUser() (*User, error) {
	db := GetDB()
	err := db.Debug().Create(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) FindAllUsers() ([]User, error) {
	var err error
	db := GetDB()
	users := []User{}
	err = db.Debug().Table("users").Limit(100).Find(&users).Error
	if err != nil {
		return []User{}, errors.New("Kullanıclar yüklenemedi")
	}
	if len(users) > 0 {
		for i, _ := range users {
			err = GetDB().Debug().Table("images").Where("id=?", &users[i].ImageID).Take(&users[i].Image).Error
			if err != nil {
				return []User{}, errors.New("Kullanıcıya ait profil fotoğrafı yok")
			}
		}
	}
	return users, nil
}

func (u *User) FindByID(uid uint) (*User, error) {
	var err error
	db = GetDB()
	err = db.Debug().Table("users").Where("id=?", uid).Take(&u).Error
	if err != nil {
		return &User{}, errors.New("Kullanıcı bulunamadı")
	}
	err = db.Debug().Table("posts").Where("user_id=?", uid).Find(&u.Posts).Error
	if err != nil {
		return &User{}, errors.New("Kullanıcıya ait gönderi verisi alınamadı")
	}
	err = db.Debug().Table("images").Where("id=?", u.ImageID).Find(&u.Image).Error
	if err != nil {
		return &User{}, errors.New("Kullanıcının profil fotoğrafı alınamadı")
	}
	return u, err
}

func (u *User) FindByUserName(username string) (*User, error) {
	var err error
	db = GetDB()
	err = db.Debug().Table("users").Where("username=?", username).Take(&u).Error
	if err != nil {
		return &User{}, errors.New("Kullanıcı bulunamadı")
	}
	err = db.Debug().Table("posts").Where("user_id=?", u.ID).Find(&u.Posts).Error
	if err != nil {
		return &User{}, errors.New("Kullanıcıya ait gönderi verisi alınamadı")
	}
	err = db.Debug().Table("images").Where("id=?", u.ImageID).Find(&u.Image).Error
	if err != nil {
		return &User{}, errors.New("Kullanıcının profil fotoğrafı alınamadı")
	}
	return u, err
}

func (u *User) UpdateAUser(uid uint) (*User, error) {
	err := u.BeforeSAve()
	if err != nil {
		log.Fatal(err)
	}
	db := GetDB().Table("users").Where("id=?", uid).UpdateColumn(
		map[string]interface{}{
			"username":   u.Username,
			"email":      u.Email,
			"name":       u.Name,
			"updated_at": time.Now(),
		},
	)
	if db.Error != nil {
		return &User{}, errors.New("Kullanıcı Güncellenemedi")
	}
	err = GetDB().Table("users").Where("id=?", uid).Take(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) DeleteByID(uid uint) (int64, error) {
	db := GetDB().Debug().Table("users").Where("id=?", uid).Take(&u).Delete(&User{})
	if db.Error != nil {
		return 0, errors.New("Kullanıcı silinemedi")
	}
	return db.RowsAffected, nil

}

func (u *User) UpdatePassword(uid uint, password string) error {
	db := GetDB().Table("users").Where("id=?", uid).UpdateColumn(
		map[string]interface{}{
			"password":   password,
			"updated_at": time.Now(),
		},
	)
	if db.Error != nil {
		return errors.New("Şifre güncelleme hatası")
	}
	err := GetDB().Table("users").Where("id=?", uid).Take(&u).Error
	if err != nil {
		return err
	}
	return nil
}

func (u *User) UpdateAUserByAdmin(email string) (*User, error) {
	err := u.BeforeSAve()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("user: ", u)
	db := GetDB().Table("users").Where("email=?", email).UpdateColumn(
		map[string]interface{}{
			"user_role":  u.UserRole,
			"isvalid":    u.Isvalid,
			"updated_at": time.Now(),
		},
	)
	if db.Error != nil {
		return &User{}, errors.New("Kullanıcı Güncellenemedi")
	}
	err = GetDB().Table("users").Where("email=?", email).Take(&u).Error
	if err != nil {
		return &User{}, err
	}
	fmt.Println(u)
	return u, nil
}
