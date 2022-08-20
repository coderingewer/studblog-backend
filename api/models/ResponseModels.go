package models

import "time"

type ResponseUser struct {
	ID        uint   `json:"ID"`
	Name      string `json:"name"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	UserImage Image  `json:"user_image"`
	Token     string `json:"token"`
	UserRole  string `json:"user_role"`
}
type UserResponse struct {
	ID        uint   `json:"ID"`
	Name      string `json:"name"`
	Username  string `json:"username"`
	UserImage Image  `json:"user_image"`
	Email     string `json:"email"`
	UserRole  string `json:"user_role"`
}
type UpdatePost struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type PostResponse struct {
	ID        uint         `json:"ID"`
	Title     string       `json:"title"`
	Content   string       `json:"content"`
	Views     []View       `json:"views"`
	Image     Image        `json:"image"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	IsValid   bool         `json:"isValid"`
	Sender    UserResponse `json:"sender"`
}

func (usrR UserResponse) UserToResponse(usr User) UserResponse {
	usrR.ID = usr.ID
	usrR.Name = usr.Name
	usrR.UserImage = usr.Image
	usrR.Username = usr.Username
	usrR.Email = usr.Email
	usrR.UserRole = usr.UserRole
	return usrR
}

func (pstR PostResponse) PostToResponse(post Post) PostResponse {
	pstR.ID = post.ID
	pstR.Title = post.Title
	pstR.Content = post.Content
	pstR.Views = post.Views
	pstR.IsValid = post.IsValid
	pstR.CreatedAt = post.CreatedAt
	pstR.UpdatedAt = post.UpdatedAt
	pstR.Image = post.Image
	usrR := UserResponse{}
	postUsr := usrR.UserToResponse(post.Sender)
	pstR.Sender = postUsr
	return pstR
}
