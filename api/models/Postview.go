package models

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

type View struct {
	gorm.Model
	PostID uint `json:"postId"`
}

func (view *View) ViewPost(pid uint) error {
	view.PostID = pid
	err := GetDB().Create(&view).Error
	if err != nil {
		return err
	}
	return nil
}

func (view *View) DeleteViews(pid uint) error {
	views := []View{}
	err := GetDB().Debug().Table("views").Where("post_id=?", pid).Find(&views).Error
	if err != nil {
		fmt.Println("hua-1")
		return err
	}
	if len(views) > 0 {
		for i, _ := range views {
			db := GetDB().Debug().Table("views").Where("id=? ", views[i].ID).Take(&views[i]).Delete(&View{})
			if db.Error != nil {
				fmt.Println("hua-2")
				return db.Error
			}
		}
	}
	return nil
}
