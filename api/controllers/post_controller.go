package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"studapp-blog/api/auth"
	"studapp-blog/api/models"
	"studapp-blog/api/utils"

	"github.com/gorilla/mux"
)

func CreatePost(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	post := models.Post{}
	err = json.Unmarshal(body, &post)
	if err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		utils.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	user := models.User{}
	author, err := user.FindByID(uid)
	if !author.Isvalid {
		utils.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	photo := models.Image{}
	photo.Prepare()
	img, err := photo.SaveImage()
	if err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	post.PhotoID = uint(img.ID)
	post.UserID = uint(uid)
	post.Category = strings.ToLower(post.Category)
	post.Prepare()
	postCreated, err := post.Save()
	if err != nil {
		formattedError := utils.FormatError(err.Error())
		utils.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	respPost := models.PostdetailReponse{}
	restPst := respPost.PostToPostDetailResponse(*postCreated)

	/*
		tag := models.Tag{}
		err = tag.CreatTag(post.ID)
		if err != nil {
			utils.ERROR(w, http.StatusUnprocessableEntity, err)
			return
		}*/

	w.Header().Set("Location", fmt.Sprintf("%s%s%d", r.Host, r.URL, post.ID))
	utils.JSON(w, http.StatusCreated, restPst)
}

func GetPosts(w http.ResponseWriter, r *http.Request) {
	post := models.Post{}
	posts, err := post.FindAllPosts()
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	respPosts := []models.PostResponse{}
	respPost := models.PostResponse{}
	if len(posts) > 0 {
		for i, _ := range posts {
			restPst := respPost.PostToResponse(posts[i])
			respPosts = append(respPosts, restPst)
		}
	}
	utils.JSON(w, http.StatusOK, respPosts)
}

func GetPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		utils.ERROR(w, http.StatusBadRequest, err)
	}
	post := models.Post{}
	postReceived, err := post.FindByID(uint(pid))
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	view := models.View{}
	err = view.ViewPost(postReceived.ID)
	if err != nil {
		utils.ERROR(w, http.StatusBadRequest, err)
	}
	respPost := models.PostdetailReponse{}
	restPst := respPost.PostToPostDetailResponse(*postReceived)
	utils.JSON(w, http.StatusOK, restPst)
}

func UpdatePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		utils.ERROR(w, http.StatusBadRequest, err)
		return
	}

	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		utils.ERROR(w, http.StatusUnauthorized, errors.New("Yetkisi yok"))
		return
	}
	post := models.Post{}
	err = models.GetDB().Debug().Table("posts").Where("id = ?", pid).Take(&post).Error
	if err != nil {
		utils.ERROR(w, http.StatusNotFound, err)
		return
	}
	post.UserID = uint(uid)
	if uid != post.UserID {
		utils.ERROR(w, http.StatusUnauthorized, errors.New("Yetkisi yok"))
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("pufw")
		utils.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	postUpdate := models.Post{}
	json.Unmarshal(body, &postUpdate)

	postUpdate.UserID = post.UserID
	postUpdate.Category = strings.ToLower(postUpdate.Category)
	if uid != postUpdate.UserID {
		utils.ERROR(w, http.StatusUnauthorized, errors.New("Yetkisi yok"))
		return
	}
	postUpdate.Prepare()
	postUpdate.ID = post.ID

	postUpdated, err := postUpdate.UpdatePost(uint(pid))
	if err != nil {
		formattedError := utils.FormatError(err.Error())
		utils.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	respPost := models.PostResponse{}
	restPst := respPost.PostToResponse(*postUpdated)

	utils.JSON(w, http.StatusOK, restPst)

}

func DeletePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		utils.ERROR(w, http.StatusBadRequest, err)
		return
	}
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		utils.ERROR(w, http.StatusUnauthorized, errors.New("Yetkisi yok"))
		return
	}
	post := models.Post{}
	err = models.GetDB().Debug().Table("posts").Where("id = ?", pid).Take(&post).Error
	if err != nil {
		utils.ERROR(w, http.StatusNotFound, err)
		return
	}

	view := models.View{}
	err = view.DeleteViews(uint(pid))
	if err != nil {
		utils.ERROR(w, http.StatusNotFound, err)
		fmt.Println("hua")
		return
	}
	if uid != post.UserID {
		utils.ERROR(w, http.StatusUnauthorized, errors.New("Yetkisi yok"))
		return
	}

	_, err = post.DeleteByID(uint(pid))
	if err != nil {
		utils.ERROR(w, http.StatusBadRequest, err)
		return
	}

	w.Header().Set("Entity", fmt.Sprintf("%d", pid))
	utils.JSON(w, http.StatusNoContent, "")
}

func GetPostsByUserID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	post := models.Post{}
	uid, err := strconv.ParseUint(vars["userId"], 10, 64)
	if err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	posts, err := post.FinBYUserID(uint(uid))
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	respPosts := []models.PostResponse{}
	respPost := models.PostResponse{}
	if len(posts) > 0 {
		for i, _ := range posts {
			restPst := respPost.PostToResponse(posts[i])
			respPosts = append(respPosts, restPst)
		}
	}
	utils.JSON(w, http.StatusOK, respPosts)
}

func GetPostsByCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	post := models.Post{}
	category := vars["category"]
	posts, err := post.FindByCategory(strings.ToLower(category))
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	utils.JSON(w, http.StatusOK, posts)
}

func UpdatePostImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["postId"], 10, 64)
	if err != nil {
		utils.ERROR(w, http.StatusBadRequest, err)
		return
	}
	post := models.Post{}
	updatePost, err := post.FindByID(uint(pid))
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		utils.ERROR(w, http.StatusUnauthorized, errors.New("Yetkisi yok"))
		return
	}
	if updatePost.UserID != uid {
		utils.ERROR(w, http.StatusUnauthorized, errors.New("Yetkisi yok"))

		return
	}

	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	formFile, _, err := r.FormFile("file")
	uploadUrl, err := models.NewMediaUpload().FileUpload(models.File{File: formFile})
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	imgid := updatePost.Image.ID

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

func LikePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		utils.ERROR(w, http.StatusBadRequest, err)
		return
	}
	like := models.Like{}
	like.PostID = uint(pid)
	err = like.LikePOst(uint(pid))
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	utils.JSON(w, http.StatusOK, "")
}

func ViewPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		utils.ERROR(w, http.StatusBadRequest, err)
		return
	}
	view := models.View{}
	view.PostID = uint(pid)
	err = view.ViewPost(uint(pid))
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	utils.JSON(w, http.StatusOK, "")
}

func UnLikePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		utils.ERROR(w, http.StatusBadRequest, err)
		return
	}
	like := models.Like{}
	err = like.UnlikePost(uint(pid))
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	utils.JSON(w, http.StatusOK, "")
}

func GetPopularPosts(w http.ResponseWriter, r *http.Request) {
	post := models.Post{}
	posts, err := post.FindAllPosts()
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	respPosts := []models.PostdetailReponse{}
	respPost := models.PostdetailReponse{}
	if len(posts) > 0 {
		for i, _ := range posts {
			if utils.CaluclateLastDay(posts[i].CreatedAt) {
				restPst := respPost.PostToPostDetailResponse(posts[i])
				respPosts = append(respPosts, restPst)
			}
		}
		sort.Slice(respPosts, func(i, j int) bool {
			return respPosts[i].Views > respPosts[j].Views
		})
	}
	if len(respPosts) >= 5 {
		respPosts = respPosts[:5]
	}
	utils.JSON(w, http.StatusOK, respPosts)
}
