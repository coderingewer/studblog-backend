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

func CreateComment(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	comment := models.Comment{}
	err = json.Unmarshal(body, &comment)
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
	comment.UserID = uid
	comment.Prepare()
	commetcommented, err := comment.CraateComment()
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, errors.New(http.StatusText(http.StatusInternalServerError)))
	}
	respComment := models.CommentdetailReponse{}
	respC := respComment.CommentToCommentDetailResponse(*commetcommented)
	w.Header().Set("Location", fmt.Sprintf("%s%s%d", r.Host, r.URL, comment.ID))
	utils.JSON(w, http.StatusCreated, respC)
}

func DeleteComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		utils.ERROR(w, http.StatusBadRequest, err)
		return
	}
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		utils.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	comment := models.Comment{}
	err = models.GetDB().Debug().Table("comments").Where("id=?", cid).Take(comment).Error
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, errors.New(http.StatusText(http.StatusInternalServerError)))
		return
	}

	if uid != comment.UserID {
		utils.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	err = comment.DeleteComment(uint(cid))
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, errors.New(http.StatusText(http.StatusInternalServerError)))
		return
	}
	utils.JSON(w, http.StatusOK, "")
}

func GetComments(w http.ResponseWriter, r *http.Request) {
	comment := models.Comment{}
	comments, err := comment.FindAll()
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	respcomments := []models.CommentdetailReponse{}
	respcomment := models.CommentdetailReponse{}
	if len(comments) > 0 {
		for i, _ := range comments {
			restPst := respcomment.CommentToCommentDetailResponse(comments[i])
			respcomments = append(respcomments, restPst)
		}
	}
	utils.JSON(w, http.StatusOK, respcomments)
}

func GetCommentsByPostID(w http.ResponseWriter, r *http.Request) {
	comment := models.Comment{}
	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["postId"], 10, 64)
	if err != nil {
		utils.ERROR(w, http.StatusBadRequest, err)
	}
	comments, err := comment.FindByPostID(uint(pid))
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	respcomments := []models.CommentdetailReponse{}
	respcomment := models.CommentdetailReponse{}
	if len(comments) > 0 {
		for i, _ := range comments {
			restPst := respcomment.CommentToCommentDetailResponse(comments[i])
			respcomments = append(respcomments, restPst)
		}
	}
	utils.JSON(w, http.StatusOK, respcomments)
}

func GetCommentsByUserID(w http.ResponseWriter, r *http.Request) {
	comment := models.Comment{}
	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["userId"], 10, 64)
	if err != nil {
		utils.ERROR(w, http.StatusBadRequest, err)
	}
	comments, err := comment.FindByUserID(uint(uid))
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	respcomments := []models.CommentdetailReponse{}
	respcomment := models.CommentdetailReponse{}
	if len(comments) > 0 {
		for i, _ := range comments {
			restPst := respcomment.CommentToCommentDetailResponse(comments[i])
			respcomments = append(respcomments, restPst)
		}
	}
	utils.JSON(w, http.StatusOK, respcomments)
}

func GetComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		utils.ERROR(w, http.StatusBadRequest, err)
	}
	cmt := models.Comment{}
	commentReceived, err := cmt.FindByID(uint(cid))
	if err != nil {
		utils.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	if err != nil {
		utils.ERROR(w, http.StatusBadRequest, err)
	}
	respC := models.CommentdetailReponse{}
	restCmt := respC.CommentToCommentDetailResponse(commentReceived)
	utils.JSON(w, http.StatusOK, restCmt)
}

func UpdateComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		utils.ERROR(w, http.StatusBadRequest, err)
		return
	}

	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		utils.ERROR(w, http.StatusUnauthorized, errors.New("Yetkisi yok"))
		return
	}
	usr := models.User{}
	user, err := usr.FindByID(uid)
	cmt := models.Comment{}
	err = models.GetDB().Debug().Table("comments").Where("id = ?", cid).Take(&cmt).Error
	if err != nil {
		utils.ERROR(w, http.StatusNotFound, err)
		return
	}
	if uid != cmt.UserID || user.UserRole != "SUPER-USER" {
		utils.ERROR(w, http.StatusUnauthorized, errors.New("Yetkisi yok"))
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	cmtUpdate := models.Comment{}
	json.Unmarshal(body, &cmtUpdate)

	cmtUpdate.UserID = cmt.UserID
	if uid != cmtUpdate.UserID {
		utils.ERROR(w, http.StatusUnauthorized, errors.New("Yetkisi yok"))
		return
	}
	cmtUpdate.Prepare()
	cmtUpdate.ID = cmt.ID

	cmtUpdated, err := cmtUpdate.UpdateComment(uint(cid))
	if err != nil {
		formattedError := utils.FormatError(err.Error())
		utils.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	respComm := models.CommentdetailReponse{}
	restCmt := respComm.CommentToCommentDetailResponse(*cmtUpdated)

	utils.JSON(w, http.StatusOK, restCmt)

}
