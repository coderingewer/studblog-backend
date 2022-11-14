package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"studapp-blog/api/auth"
	"studapp-blog/api/models"
	"studapp-blog/api/utils"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
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
	user.BeforeSAve()
	err = user.Validate("register")
	if err != nil {
		formatedError := utils.FormatError(err.Error())
		utils.ERROR(w, http.StatusBadRequest, formatedError)
		return
	}
	usrImg := models.Image{}
	img, err := usrImg.SaveImage()
	if err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	user.ImageID = img.ID

	/*
		token, err := auth.CreateToken(user.ID)
		if err != nil {
			utils.ERROR(w, http.StatusUnprocessableEntity, err)
			return
		}

			data := utils.EmailData{}
			data.FirstName = user.Name
			data.URL = string(os.Getenv("DOMAIN")) + "/api/users/confirm/" + string(token)
			data.Subject = "Hesabınızı onaylayın"
			err = utils.SendMail(user.Email, data, "confirm")
			if err != nil {
				utils.ERROR(w, http.StatusUnprocessableEntity, err)
				return
			}*/
	userCreated, err := user.SaveUser()
	usrRes := models.UserResponse{}
	usrR := usrRes.UserToResponse(*userCreated)
	if err != nil {
		formatedError := utils.FormatError(err.Error())
		utils.ERROR(w, http.StatusUnprocessableEntity, formatedError)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("%s%s%d", r.Host, r.RequestURI, userCreated.ID))
	utils.JSON(w, http.StatusCreated, usrR)

}
func CreateUserByAdmin(w http.ResponseWriter, r *http.Request) {
	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		utils.ERROR(w, http.StatusUnauthorized, errors.New("Yetkilendirilmemiş"))
		return
	}
	user := models.User{}
	usr, err := user.FindByID(uint(tokenID))
	if usr.UserRole != "SUPER-USER" {
		utils.ERROR(w, http.StatusUnauthorized, errors.New("Yetkilendirilmemiş"))
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	err = json.Unmarshal(body, &user)
	if err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	err = user.Validate("register")
	if err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	usrImg := models.Image{}
	img, err := usrImg.SaveImage()
	if err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	user.ImageID = img.ID
	fmt.Println("id: ", img.ID)

	/*
		token, err := auth.CreateToken(user.ID)
		if err != nil {
			utils.ERROR(w, http.StatusUnprocessableEntity, err)
			return
		}

		data := utils.EmailData{}
		data.FirstName = user.Name
		data.URL = string(os.Getenv("DOMAIN")) + "/api/users/confirm/" + string(token)
		data.Subject = "Hesabınızı onaylayın"
		err = utils.SendMail(user.Email, data, "confirm")
		if err != nil {
			utils.ERROR(w, http.StatusUnprocessableEntity, err)
			return
			}*/
	user.BeforeSAve()
	user.Prepare()
	userCreated, err := user.SaveUser()
	usrRes := models.UserResponse{}
	usrR := usrRes.UserToResponse(*userCreated)
	if err != nil {
		formatedError := utils.FormatError(err.Error())
		utils.ERROR(w, http.StatusUnprocessableEntity, formatedError)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("%s%s%d", r.Host, r.RequestURI, userCreated.ID))
	utils.JSON(w, http.StatusCreated, usrR)

}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	users, err := user.FindAllUsers()
	if err != nil {
		formattedError := utils.FormatError(err.Error())
		utils.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	usersResp := []models.UserDetailResponse{}
	userResp := models.UserDetailResponse{}
	if len(users) > 0 {
		for i, _ := range users {
			respUsr := userResp.UserToUserDetailResponse(users[i])
			usersResp = append(usersResp, respUsr)
		}
	}
	utils.JSON(w, http.StatusOK, usersResp)
}

func ConfirmAcoount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token := vars["token"]
	uid, err := auth.ExtractTokenFromStringID(token)
	if err != nil {
		utils.ERROR(w, http.StatusUnauthorized, errors.New("onay başarısız"))
		return
	}
	user := models.User{}
	userGotten, err := user.FindByID(uint(uid))
	userGotten.Isvalid = true
	if err != nil {
		utils.ERROR(w, http.StatusUnauthorized, err)
		return
	}
	utils.JSON(w, http.StatusOK, nil)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		utils.ERROR(w, http.StatusBadRequest, err)
		return
	}
	user := models.User{}
	userGotten, err := user.FindByID(uint(uid))
	if err != nil {
		utils.ERROR(w, http.StatusBadRequest, err)
		return
	}
	usrRes := models.UserDetailResponse{}
	usrR := usrRes.UserToUserDetailResponse(*userGotten)
	utils.JSON(w, http.StatusOK, usrR)
}

func GetUserByUserName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]
	user := models.User{}
	userGotten, err := user.FindByUserName(username)
	if err != nil {
		utils.ERROR(w, http.StatusBadRequest, err)
		return
	}
	usrRes := models.UserDetailResponse{}
	usrR := usrRes.UserToUserDetailResponse(*userGotten)
	utils.JSON(w, http.StatusOK, usrR)
}

func GetUserByToken(w http.ResponseWriter, r *http.Request) {
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		utils.ERROR(w, http.StatusUnauthorized, errors.New("Yetkilendirilmemiş"))
		fmt.Println("hao2")
	}
	user := models.User{}
	userGotten, err := user.FindByID(uint(uid))
	if err != nil {
		utils.ERROR(w, http.StatusBadRequest, err)
		return
	}
	usrRes := models.UserResponse{}
	usrR := usrRes.UserToResponse(*userGotten)
	utils.JSON(w, http.StatusOK, usrR)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["id"], 10, 64)
	user := models.User{}
	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		utils.ERROR(w, http.StatusUnauthorized, errors.New("Yetkilendirilmemiş"))
		fmt.Println("hao2")
	}

	admin, err := user.FindByID(uint(tokenID))
	if err != nil {
		formattedError := utils.FormatError(err.Error())
		utils.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}

	if tokenID != uint(uid) && admin.UserRole != "SUPER-USER" {
		utils.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		fmt.Println("hao1")
	}
	if err != nil {
		utils.ERROR(w, http.StatusBadRequest, err)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, err)
		fmt.Println("s")
		return
	}
	err = json.Unmarshal(body, &user)
	if err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, err)
		fmt.Println("d")
		return
	}

	user.Prepare()
	updatedUser, err := user.UpdateAUser(uint(uid))
	if err != nil {
		formattedError := utils.FormatError(err.Error())
		utils.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	usrRes := models.UserResponse{}
	usrR := usrRes.UserToResponse(*updatedUser)
	utils.JSON(w, http.StatusOK, usrR)
}

func UpdatePassword(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["id"], 10, 64)
	user := models.User{}
	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		utils.ERROR(w, http.StatusUnauthorized, errors.New("Yetkilendirilmemiş"))
		fmt.Println("hao2")
	}
	usr, err := user.FindByID(uint(uid))
	if tokenID != uint(uid) {
		utils.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		fmt.Println("hao1")
	}
	if err != nil {
		utils.ERROR(w, http.StatusBadRequest, err)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	UpdatePassUsr := models.UpdatePassUser{}
	err = json.Unmarshal(body, &UpdatePassUsr)
	if err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	_ = UpdatePassUsr.BeforeSAve()
	if err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	err = models.VerifyPassword(usr.Password, UpdatePassUsr.Password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		fmt.Println(user.Password)
		utils.ERROR(w, http.StatusUnauthorized, errors.New("Şifre yanlış"))
		return
	}

	err = user.UpdatePassword(uint(uid), UpdatePassUsr.NewPassword)
	if err != nil {
		formattedError := utils.FormatError(err.Error())
		utils.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	fmt.Println(usr.Password)
	utils.JSON(w, http.StatusOK, "")
}

func UpdatePasswordByAdmin(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["id"], 10, 64)
	user := models.User{}
	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		utils.ERROR(w, http.StatusUnauthorized, errors.New("Yetkilendirilmemiş"))
		fmt.Println("hao2")
	}
	admin, err := user.FindByID(uint(tokenID))
	if err != nil {
		formattedError := utils.FormatError(err.Error())
		utils.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}

	if admin.UserRole != "SUPER-USER" {
		utils.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		fmt.Println("hao1")
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	UpdatePassUsr := models.UpdatePassUser{}
	err = json.Unmarshal(body, &UpdatePassUsr)
	if err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	user.Prepare()
	if err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	err = user.UpdatePassword(uint(uid), UpdatePassUsr.NewPassword)
	if err != nil {
		formattedError := utils.FormatError(err.Error())
		utils.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	utils.JSON(w, http.StatusOK, "")
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := models.User{}

	uid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		utils.ERROR(w, http.StatusUnauthorized, errors.New("Yetkilendirilmemiş"))
	}
	if tokenID != uint(uid) {
		utils.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
	}
	_, err = user.DeleteByID(uint(uid))
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", uid))
	utils.JSON(w, http.StatusNoContent, "")

}

func DeleteUserByAdmin(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := models.User{}
	uid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		utils.ERROR(w, http.StatusUnauthorized, errors.New("Yetkilendirilmemiş"))
	}
	admin, err := user.FindByID(uint(tokenID))
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, errors.New("Yetkilendirilmemiş"))
	}
	if admin.UserRole != "SUPER-USER" {
		utils.ERROR(w, http.StatusUnauthorized, errors.New("Yetkilendirilmemiş"))
	}

	_, err = user.DeleteByID(uint(uid))
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", uid))
	utils.JSON(w, http.StatusNoContent, "")
}

func UpdateUsermage(w http.ResponseWriter, r *http.Request) {
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		utils.ERROR(w, http.StatusUnauthorized, errors.New("Yetkisi yok"))
		return
	}
	user := models.User{}
	userr, err := user.FindByID(uint(uid))
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	if userr.ID != uid {
		utils.ERROR(w, http.StatusUnauthorized, errors.New("Yetkisi yok"))

		return
	}

	formFile, _, err := r.FormFile("file")
	uploadUrl, err := models.NewMediaUpload().FileUpload(models.File{File: formFile})
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	imgid := userr.ImageID

	img := models.Image{}

	err = models.GetDB().Debug().Table("images").Where("id = ?", imgid).Take(&img).Error
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	imageUpdate := models.Image{}
	imageUpdate.Prepare()
	imageUpdate.ID = img.ID
	imageUpdate.Url = uploadUrl
	imgUpdated, err := imageUpdate.UpdateImageByID(uint(imgid))
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)

		return
	}
	utils.JSON(w, http.StatusOK, imgUpdated)
}

func UpdateUserByAdmin(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	email := vars["email"]

	user := models.User{}
	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		utils.ERROR(w, http.StatusUnauthorized, errors.New("Yetkilendirilmemiş"))
		fmt.Println("hao2")
		return
	}

	admin, err := user.FindByID(uint(tokenID))
	if err != nil {
		formattedError := utils.FormatError(err.Error())
		utils.ERROR(w, http.StatusInternalServerError, formattedError)
		fmt.Println("c")
		return
	}

	if admin.UserRole != "SUPER-USER" {
		utils.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		fmt.Println("hao1")
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, err)
		fmt.Println("s")
		return
	}
	updateUser := models.User{}
	err = json.Unmarshal(body, &updateUser)
	if err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, err)
		fmt.Println("d")
		return
	}
	updatedUser, err := updateUser.UpdateAUserByAdmin(email)
	if err != nil {
		formattedError := utils.FormatError(err.Error())
		utils.ERROR(w, http.StatusInternalServerError, formattedError)
		fmt.Println("g")
		return
	}
	usrRes := models.UserResponse{}
	usrR := usrRes.UserToResponse(*updatedUser)
	utils.JSON(w, http.StatusOK, usrR)
}
