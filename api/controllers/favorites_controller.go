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
)

func NewFavsList(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	list := models.FavoritesList{}
	err = json.Unmarshal(body, &list)
	if err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		utils.ERROR(w, http.StatusUnauthorized, errors.New("Yetkilendirilmemiş"))
		return
	}
	photo := models.Image{}
	img, err := photo.SaveImage()
	if err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	list.CoverImgId = img.ID
	list.UserId = uint(uid)
	listcreated, err := list.CreateAList()
	if err != nil {
		formattedError := utils.FormatError(err.Error())
		utils.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("%s%s%d", r.Host, r.URL, list.ID))
	utils.JSON(w, http.StatusCreated, listcreated)
}

func AddItemToList(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	item := models.Favorite{}
	err = json.Unmarshal(body, &item)
	if err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	favs := models.FavoritesList{}
	favsList, err := favs.FindByID(item.ListId)

	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		utils.ERROR(w, http.StatusUnauthorized, errors.New("Yetkilendirilmemiş"))
		return
	}
	if favsList.UserId != uid {
		utils.ERROR(w, http.StatusUnauthorized, errors.New("Liste size ait değil"))
		return
	}
	itemAdded, err := item.AddPostToList()
	if err != nil {
		formattedError := utils.FormatError(err.Error())
		utils.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("%s%s%d", r.Host, r.URL, item.ID))
	utils.JSON(w, http.StatusCreated, itemAdded)
}

func DeleteFavsList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	listid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		utils.ERROR(w, http.StatusBadRequest, err)
		return
	}
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		utils.ERROR(w, http.StatusUnauthorized, errors.New("Yetkisi yok"))
		return
	}
	list := models.FavoritesList{}
	listGeted, err := list.FindByID(uint(listid))

	if uid != listGeted.UserId {
		utils.ERROR(w, http.StatusUnauthorized, errors.New("Yetkisi yok"))
		return
	}
	err = listGeted.DeleteList(uint(listid))
	if err != nil {
		utils.ERROR(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", listid))
	utils.JSON(w, http.StatusNoContent, "")
}

func DeleteFavFromList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	listitemId, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		utils.ERROR(w, http.StatusBadRequest, err)
		return
	}
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		utils.ERROR(w, http.StatusUnauthorized, errors.New("Yetkisi yok"))
		return
	}
	listitem := models.Favorite{}
	listitemGeted, err := listitem.FindByID(uint(listitemId))

	if uid != listitemGeted.UserId {
		utils.ERROR(w, http.StatusUnauthorized, errors.New("Yetkisi yok"))
		return
	}
	err = listitemGeted.RemoveFromList(uint(listitemId))
	if err != nil {
		utils.ERROR(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", listitemId))
	utils.JSON(w, http.StatusNoContent, "")
}

func UpdateFavsList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	listId, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		utils.ERROR(w, http.StatusBadRequest, err)
		return
	}
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		utils.ERROR(w, http.StatusUnauthorized, errors.New("Yetkisi yok"))
		return
	}
	list := models.FavoritesList{}
	listGeted, err := list.FindByID(uint(listId))
	if uid != listGeted.UserId {
		utils.ERROR(w, http.StatusUnauthorized, errors.New("Yetkisi yok"))
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	listeUpdate := models.FavoritesList{}
	json.Unmarshal(body, &listeUpdate)
	listeUpdate.UserId = listGeted.UserId
	listeUpdate.ID = listGeted.ID
	listUpdated, err := listeUpdate.UpdateFavList(uint(listId))
	if err != nil {
		formattedError := utils.FormatError(err.Error())
		utils.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	utils.JSON(w, http.StatusOK, listUpdated)

}
