package controllers

import (
	"net/http"
	"strconv"
	"studapp-blog/api/models"
	"studapp-blog/api/utils"

	"github.com/gorilla/mux"
)

func ImgUpload(w http.ResponseWriter, r *http.Request) {
	img := models.Image{}
	formFile, _, err := r.FormFile("file")
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	uploadUrl, err := models.NewMediaUpload().FileUpload(models.File{File: formFile})
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	img.Url = uploadUrl
	image, err := img.SaveImage()
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	utils.JSON(w, http.StatusOK, image)
}

func UpdateImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	imgid, err := strconv.ParseUint(vars["imageId"], 10, 64)
	if err != nil {
		utils.ERROR(w, http.StatusBadRequest, err)
		return
	}
	formFile, _, err := r.FormFile("file")
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	uploadUrl, err := models.NewMediaUpload().FileUpload(models.File{File: formFile})
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}

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

func GetAllImages(w http.ResponseWriter, r *http.Request) {
	media := models.Image{}
	medias, err := media.FindAll()
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	utils.JSON(w, http.StatusOK, medias)
}
