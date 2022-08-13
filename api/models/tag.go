package models

import (
	"strings"

	"github.com/jinzhu/gorm"
)

type Tag struct {
	gorm.Model
	Tagname string `gorm:"not null; unique " json:"tagname"`
	Flowers []User `gorm:"foreignkey:TagID" json:"flowers"`
}

type PostTag struct {
	gorm.Model
	PostID uint `json:"post_id"`
	TagID  uint `json:"tag_id"`
}

func (tag *Tag) Prepare() {
}

func (tag *Tag) CreatTag(pid uint) error {
	post := Post{}
	err := GetDB().Where("id = ?", pid).First(&post).Error
	if err != nil {
		return err
	}
	var text string
	text = post.Content
	text = strings.ToLower(text)
	stext := strings.Split(text, " ")

	if len(stext) > 0 {
		for i, _ := range stext {
			if stext[i][0:1] == "#" {
				if tag.GetTag(stext[i][1:]) == false {
					tag.Tagname = stext[i][1:]
					tag.SaveTag()
					tag.SavePostTag(pid)
				} else {
					tag.SavePostTag(pid)
				}
			}
		}
		return nil
	}
	return nil
}

func (tag *Tag) GetTag(tagname string) bool {
	err := GetDB().Table("tags").Where("tagname = ?", tagname).First(&tag).Error
	if err != nil {
		return false
	}
	return true
}

func (tag *Tag) SaveTag() error {
	err := GetDB().Create(&tag).Error
	if err != nil {
		return err
	}
	return nil
}

func (tag *Tag) SavePostTag(pid uint) error {
	posttag := PostTag{}
	posttag.PostID = pid
	posttag.TagID = tag.ID
	err := GetDB().Create(&posttag).Error
	if err != nil {
		return err
	}
	return nil
}
