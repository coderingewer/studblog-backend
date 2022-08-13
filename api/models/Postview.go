package models

import "github.com/jinzhu/gorm"

type View struct {
	gorm.Model
	PostID uint `json:"postId"`
}

func (view *View) ViewPost(pid uint) error {
	err := GetDB().Create(&view).Error
	if err != nil {
		return err
	}
	return nil
}
