package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Comment struct {
	gorm.Model
	PostID    uint   `json:"postId"`
	UserID    uint   `json:"userId"`
	User      User   `json:"user"`
	Content   string `json:"content"`
	Reviewed  bool   `json:"reviewed"`
	Validated bool   `json:"validated"`
}

func (c *Comment) Prepare() {
	c.Validated = true
	c.Reviewed = false
}

func (cmt *Comment) CraateComment() (*Comment, error) {
	err := GetDB().Debug().Create(cmt).Error
	if err != nil {
		return &Comment{}, err
	}
	err = GetDB().Debug().Table("users").Where("id=?", cmt.UserID).Take(cmt.User).Error
	if err != nil {
		return &Comment{}, err
	}
	return cmt, nil
}

func (cmt *Comment) DeleteComment(cid uint) error {
	db := GetDB().Debug().Table("comments").Where("id=?", cid).Delete(Comment{})
	if db.Error != nil {
		return db.Error
	}
	return nil
}

func (cmt Comment) FindByID(cid uint) (Comment, error) {
	err := GetDB().Debug().Table("comments").Where("id=?", cid).Take(cmt).Error
	if err != nil {
		return Comment{}, err
	}

	err = GetDB().Debug().Table("users").Where("id=?", cmt.UserID).Take(cmt.User).Error
	if err != nil {
		return Comment{}, err
	}
	return cmt, nil
}

func (cmt Comment) FindByUserID(uid uint) ([]Comment, error) {
	comments := []Comment{}

	err := GetDB().Debug().Table("comments").Where("user_id=?", uid).Find(comments).Error
	if err != nil {
		return []Comment{}, err
	}
	if len(comments) > 0 {
		for i, _ := range comments {
			err = GetDB().Debug().Table("users").Where("id=?", comments[i].UserID).Take(comments[i].User).Error
			if err != nil {
				return []Comment{}, err
			}
		}
	}
	return comments, nil
}

func (cmt Comment) FindByPostID(pid uint) ([]Comment, error) {
	comments := []Comment{}

	err := GetDB().Debug().Table("comments").Where("post_id=?", pid).Find(comments).Error
	if err != nil {
		return []Comment{}, err
	}
	if len(comments) > 0 {
		for i, _ := range comments {
			err = GetDB().Debug().Table("users").Where("id=?", comments[i].UserID).Take(comments[i].User).Error
			if err != nil {
				return []Comment{}, err
			}
		}
	}
	return comments, nil
}

func (cmt *Comment) FindAll() ([]Comment, error) {
	comments := []Comment{}

	err := GetDB().Debug().Table("comments").Find(comments).Error
	if err != nil {
		return []Comment{}, err
	}
	if len(comments) > 0 {
		for i, _ := range comments {
			err = GetDB().Debug().Table("users").Where("id=?", comments[i].UserID).Take(comments[i].User).Error
			if err != nil {
				return []Comment{}, err
			}
		}
	}
	return comments, nil
}

func (cmt *Comment) UpdateComment(cid uint) (*Comment, error) {
	db := GetDB().Debug().Table("comments").Where("id=?", cid).UpdateColumns(
		map[string]interface{}{
			"content":    cmt.Content,
			"validated":  cmt.Validated,
			"reviewed":   cmt.Reviewed,
			"updated_at": time.Now(),
		},
	)
	if db.Error != nil {
		return &Comment{}, db.Error
	}
	err := GetDB().Debug().Table("posts").Where("id=?", cid).Take(&cmt).Error
	if err != nil {
		return &Comment{}, err
	}
	return cmt, nil
}
