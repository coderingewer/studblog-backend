package models

import "github.com/jinzhu/gorm"

type Like struct {
	gorm.Model
	PostID uint `json:"postId"`
}

func (like *Like) LikePOst(pid uint) error {
	err := GetDB().Create(&like).Error
	if err != nil {
		return err
	}
	return nil
}

func (like *Like) UnlikePost(pid uint) error {
	err := GetDB().Delete(&like).Error
	if err != nil {
		return err
	}
	return nil
}
