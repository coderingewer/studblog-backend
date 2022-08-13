package models

import "time"

type ResponseUser struct {
	ID        uint      `json:"ID"`
	Name      string    `json:"name"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	UserImage Image     `json:"user_image"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserRole  string    `json:"user_role"`
}
type UpdatePost struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}
