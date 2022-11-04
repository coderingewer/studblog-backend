package models

import (
	"fmt"
	"time"
)

type ResponseUser struct {
	ID        uint   `json:"ID"`
	Name      string `json:"name"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	UserImage Image  `json:"user_image"`
	Token     string `json:"token"`
	UserRole  string `json:"user_role"`
	Isvalid   bool   `json:"isValid"`
}

type UserResponse struct {
	ID        uint   `json:"ID"`
	Name      string `json:"name"`
	Username  string `json:"username"`
	UserImage Image  `json:"user_image"`
	Email     string `json:"email"`
	UserRole  string `json:"user_role"`
	Isvalid   bool   `json:"isValid"`
}

type UserDetailResponse struct {
	ID           uint   `json:"ID"`
	Name         string `json:"name"`
	Username     string `json:"username"`
	UserImageUrl string `json:"userImageUrl"`
	UserImageID  uint   `json:"userImageId"`
	Email        string `json:"email"`
	UserRole     string `json:"user_role"`
	Isvalid      bool   `json:"isValid"`
}

type UpdatePost struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type UpdatePassUser struct {
	Password    string `json:"password"`
	NewPassword string `json:"newPassword"`
}

type PostResponse struct {
	ID        uint         `json:"ID"`
	Title     string       `json:"title"`
	Content   string       `json:"content"`
	Category  string       ` json:"category"`
	Views     []View       `json:"views"`
	Image     Image        `json:"image"`
	CreatedAt string       `json:"created_at"`
	UpdatedAt string       `json:"updated_at"`
	IsValid   bool         `json:"isValid"`
	Sender    UserResponse `json:"sender"`
}

type PostdetailReponse struct {
	ID             uint      `json:"ID"`
	Title          string    `json:"title"`
	Content        string    `json:"content"`
	Category       string    ` json:"category"`
	Views          int       `json:"views"`
	ImageUrl       string    `json:"image"`
	ImageId        uint      `json:"imageId"`
	CreatedAt      string    `json:"created_at"`
	UpdatedAt      string    `json:"updated_at"`
	IsValid        bool      `json:"isValid"`
	SenderId       uint      `json:"senderId"`
	SenderImageUrl string    `json:"senderImageUrl"`
	SenderUserName string    `json:"senderUserName"`
	SenderName     string    `json:"senderName"`
	Likes          int       `json:"likes"`
	Comments       []Comment `json:"comments"`
}

type CommentdetailReponse struct {
	ID             uint   `json:"ID"`
	Content        string `json:"content"`
	ImageId        uint   `json:"imageId"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
	IsValid        bool   `json:"isValid"`
	SenderId       uint   `json:"senderId"`
	SenderImageUrl string `json:"senderImageUrl"`
	SenderUserName string `json:"senderUserName"`
	SenderName     string `json:"senderName"`
}

type ResponseFavoritesList struct {
	ID          uint       `json:"ID"`
	Name        string     `json:"name"`
	Hide        bool       ` json:"hide"`
	CoverImgId  uint       `json:"coverImgId"`
	CoverImgUrl string     `json:"coverImageUrl"`
	UserId      uint       `json:"userId"`
	User        User       `json:"user"`
	Items       []Favorite `json:"items"`
}

type ResponseFavorite struct {
	ID     uint `json:"ID"`
	ListId uint `json:"listId"`
	PostId uint `json:"postId"`
	UserId uint `json:"userId"`
	Post   Post `json:"post"`
}

func (u *UpdatePassUser) BeforeSAve() error {
	hashedPassword, err := Hash(u.NewPassword)
	if err != nil {
		fmt.Println(string(err.Error()))
		return err
	}
	u.NewPassword = string(hashedPassword)
	return nil
}

func (usrR UserResponse) UserToResponse(usr User) UserResponse {
	usrR.ID = usr.ID
	usrR.Name = usr.Name
	usrR.UserImage = usr.Image
	usrR.Username = usr.Username
	usrR.Email = usr.Email
	usrR.UserRole = usr.UserRole
	usrR.Isvalid = usr.Isvalid

	return usrR
}

func (usrR UserDetailResponse) UserToUserDetailResponse(usr User) UserDetailResponse {
	usrR.ID = usr.ID
	usrR.Name = usr.Name
	usrR.UserImageID = usr.Image.ID
	usrR.UserImageUrl = usr.Image.Url
	usrR.Username = usr.Username
	usrR.Email = usr.Email
	usrR.UserRole = usr.UserRole
	usrR.Isvalid = usr.Isvalid
	return usrR
}

func CustomDate(Date time.Time) string {
	return fmt.Sprintf("%v/%v/%v", Date.Day(), Date.Local().Month(), Date.Year())
}

func (pstR PostResponse) PostToResponse(post Post) PostResponse {
	pstR.ID = post.ID
	pstR.Title = post.Title
	pstR.Content = post.Content
	pstR.Views = post.Views
	pstR.IsValid = post.IsValid
	pstR.CreatedAt = CustomDate(post.CreatedAt)
	pstR.UpdatedAt = CustomDate(post.UpdatedAt)
	pstR.Image = post.Image
	pstR.Category = post.Category
	usrR := UserResponse{}
	postUsr := usrR.UserToResponse(post.Sender)
	pstR.Sender = postUsr
	return pstR
}

func (pstR PostdetailReponse) PostToPostDetailResponse(post Post) PostdetailReponse {
	pstR.ID = post.ID
	pstR.Title = post.Title
	pstR.Content = post.Content
	pstR.Views = len(post.Views)
	pstR.Likes = len(post.Likes)
	pstR.IsValid = post.IsValid
	pstR.CreatedAt = CustomDate(post.CreatedAt)
	pstR.UpdatedAt = CustomDate(post.UpdatedAt)
	pstR.ImageUrl = post.Image.Url
	pstR.ImageId = post.Image.ID
	pstR.Category = post.Category
	pstR.SenderId = post.Sender.ID
	pstR.SenderName = post.Sender.Name
	pstR.SenderImageUrl = post.Sender.Image.Url
	pstR.SenderUserName = post.Sender.Username
	return pstR
}

func (cmtR CommentdetailReponse) CommentToCommentDetailResponse(comment Comment) CommentdetailReponse {
	cmtR.ID = comment.ID
	cmtR.Content = comment.Content
	cmtR.IsValid = comment.Validated
	cmtR.CreatedAt = CustomDate(comment.CreatedAt)
	cmtR.UpdatedAt = CustomDate(comment.UpdatedAt)
	cmtR.SenderId = comment.UserID
	cmtR.SenderName = comment.User.Name
	cmtR.SenderImageUrl = comment.User.Image.Url
	cmtR.SenderUserName = comment.User.Username
	return cmtR
}

func (favRes ResponseFavoritesList) FavoriteToResponse(favs FavoritesList) ResponseFavoritesList {
	favRes.ID = favs.ID
	favRes.Name = favs.Name
	favRes.Hide = favs.Hide
	favRes.Items = favs.Items
	favRes.User = favs.User
	favRes.CoverImgId = favs.CoverImgId
	favRes.CoverImgUrl = favs.CoverImage.Url
	return favRes
}
