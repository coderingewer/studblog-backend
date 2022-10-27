package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"studapp-blog/api/auth"
	"studapp-blog/api/models"
	"studapp-blog/api/utils"

	"golang.org/x/crypto/bcrypt"
)

func Login(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	user.Prepare()
	err = user.Validate("login")
	if err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	email := strings.ToLower(user.Email)
	token, user, err := SignIn(email, user.Password)
	if err != nil {
		utils.ERROR(w, http.StatusBadRequest, err)
		fmt.Println(string(err.Error()))
		return
	}
	usr := models.ResponseUser{}
	usr.ID = user.ID
	usr.Token = token
	usr.Username = user.Username
	usr.Name = user.Name
	usr.UserRole = user.UserRole
	usr.Isvalid = user.Isvalid
	utils.JSON(w, http.StatusOK, usr)
}

func SignIn(email, password string) (string, models.User, error) {
	var err error
	user := models.User{}
	err = models.GetDB().Debug().Table("users").Where("email=?", email).First(&user).Error
	if err != nil {
		return "", user, errors.New("Kullanıcı bulunamadı")
	}
	err = models.VerifyPassword(user.Password, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", models.User{}, errors.New("Şifre Yanlış")
	}
	err = models.GetDB().Debug().Table("images").Where("id=?", user.ImageID).Take(&user.Image).Error
	token, err := auth.CreateToken(user.ID)
	if err != nil {
		return "", user, err
	}
	return token, user, nil
}
