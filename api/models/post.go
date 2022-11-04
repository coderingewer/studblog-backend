package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

type Post struct {
	gorm.Model
	Title    string    `json:"title"`
	Sender   User      `json:"sender"`
	Content  string    ` gorm:"not null" json:"content"`
	Category string    `json:"category"`
	UserID   uint      `gorm:"not null" json:"userId"`
	PhotoID  uint      `json:"photoId"`
	Image    Image     `json:"image"`
	Likes    []Like    `json:"likes"`
	Views    []View    `json:"views"`
	Comments []Comment `json:"comments"`
	IsValid  bool      `json:"isValid"`
	PosTags  []PostTag `gorm:"many2many:post_tags" json:"post_tags"`
}

func (p *Post) Prepare() {
	p.ID = 0
	p.Sender = User{}
	p.IsValid = true
}

func (p *Post) Save() (*Post, error) {
	err := db.Debug().Create(&p).Error
	if err != nil {
		return &Post{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?",
			p.UserID).Take(&p.Sender).Error
		if err != nil {
			fmt.Println("hoh")
			return &Post{}, err
		}
		err = GetDB().Debug().Table("images").Where("id=?", &p.PhotoID).Take(&p.Image).Error
		if err != nil {
			return &Post{}, err
		}
	}
	return p, nil
}

func (p *Post) FindAllPosts() ([]Post, error) {
	posts := []Post{}
	err := GetDB().Debug().Table("posts").Order("created_at desc").Find(&posts).Error
	if err != nil {
		return []Post{}, err
	}
	if len(posts) > 0 {
		for i, _ := range posts {
			err := GetDB().Debug().Table("users").Where("id=?", &posts[i].UserID).Take(&posts[i].Sender).Error
			err = GetDB().Debug().Table("images").Where("id=?", &posts[i].Sender.ImageID).Take(&posts[i].Sender.Image).Error
			if err != nil {
				return []Post{}, err
			}
			err = GetDB().Debug().Table("images").Where("id=?", &posts[i].PhotoID).Take(&posts[i].Image).Error
			if err != nil {
				return []Post{}, err
			}
			err = GetDB().Debug().Table("views").Where("post_id=?", &posts[i].ID).Find(&posts[i].Views).Error
			if err != nil {
				return []Post{}, err
			}
		}
	}
	return posts, nil
}

func (post *Post) FindByID(pid uint) (*Post, error) {

	p := Post{}

	err := GetDB().Debug().Table("posts").Where("id = ?", pid).Take(&p).Error
	if err != nil {
		return &Post{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &Post{}, errors.New("Gönderi bulunamadı")
	}
	err = GetDB().Debug().Table("users").Where("id=?", &p.UserID).Take(&p.Sender).Error
	if err != nil {
		return &Post{}, err
	}
	err = GetDB().Debug().Table("images").Where("id=?", &p.Sender.ImageID).Take(&p.Sender.Image).Error
	if err != nil {
		return &Post{}, err
	}
	err = GetDB().Debug().Table("images").Where("id=?", &p.PhotoID).Take(&p.Image).Error
	if err != nil {
		return &Post{}, err
	}
	err = GetDB().Debug().Table("views").Where("post_id=?", &p.ID).Find(&p.Views).Error
	if err != nil {
		return &Post{}, err
	}
	err = GetDB().Debug().Table("comments").Where("post_id=?", &p.ID).Find(&p.Comments).Error
	if err != nil {
		return &Post{}, err
	}
	return &p, nil
}

func (post *Post) FindByCategory(category string) ([]Post, error) {
	posts := []Post{}
	err := GetDB().Debug().Table("posts").Where("category= ?", category).Order("created_at desc").Find(&posts).Error
	if err != nil {
		return []Post{}, err
	}
	if len(posts) > 0 {
		for i, _ := range posts {
			err := GetDB().Debug().Table("users").Where("id=?", &posts[i].UserID).Take(&posts[i].Sender).Error
			err = GetDB().Debug().Table("images").Where("id=?", &posts[i].Sender.ImageID).Take(&posts[i].Sender.Image).Error
			if err != nil {
				return []Post{}, err
			}
			err = GetDB().Debug().Table("images").Where("id=?", &posts[i].PhotoID).Take(&posts[i].Image).Error
			if err != nil {
				return []Post{}, err
			}
			err = GetDB().Debug().Table("views").Where("post_id=?", &posts[i].ID).Find(&posts[i].Views).Error
			if err != nil {
				return []Post{}, err
			}
		}
	}
	return posts, nil
}

func (p *Post) UpdatePost(pid uint) (*Post, error) {
	db := GetDB().Debug().Table("posts").Where("id=?", pid).UpdateColumns(
		map[string]interface{}{
			"title":      p.Title,
			"content":    p.Content,
			"category":   p.Category,
			"updated_at": time.Now(),
		},
	)
	if db.Error != nil {
		return &Post{}, db.Error
	}
	err := GetDB().Debug().Table("posts").Where("id=?", pid).Take(&p).Error
	if err != nil {
		return &Post{}, err
	}
	return p, nil
}

func (p *Post) DeleteByID(pid uint) (int64, error) {
	db := GetDB().Debug().Table("posts").Where("id=? ", pid).Take(&p).Delete(&Post{})
	if db.Error != nil {
		return 0, nil
	}
	return db.RowsAffected, nil
}

func (p *Post) FinBYUserID(uid uint) ([]Post, error) {
	posts := []Post{}
	err := GetDB().Debug().Table("posts").Where("user_id = ?", uid).Order("created_at desc").Find(&posts).Error
	if err != nil {
		return []Post{}, err
	}
	if len(posts) > 0 {
		for i, _ := range posts {
			err := GetDB().Debug().Table("users").Where("id=?", &posts[i].UserID).Take(&posts[i].Sender).Error
			if err != nil {
				return []Post{}, err
			}
			err = GetDB().Debug().Table("images").Where("id=?", &posts[i].PhotoID).Take(&posts[i].Image).Error
			if err != nil {
				return []Post{}, err
			}
		}
	}
	return posts, nil
}

func (p *Post) UpdatePostImage(pid uint) (*Post, error) {
	db := GetDB().Debug().Table("posts").Where("id=?", pid).UpdateColumns(
		map[string]interface{}{
			"photo_id": p.PhotoID,
		},
	)
	if db.Error != nil {
		return &Post{}, db.Error
	}
	err := GetDB().Debug().Table("posts").Where("id=?", pid).Take(&p).Error
	if err != nil {
		return &Post{}, err
	}
	return p, nil
}

func (p *Post) FindPopularPosts() ([]Post, error) {
	posts := []Post{}
	err := GetDB().Debug().Table("posts").Order("created_at desc").Limit(5).Find(&posts).Limit(10).Error
	if err != nil {
		return []Post{}, err
	}
	if len(posts) > 0 {
		for i, _ := range posts {
			err := GetDB().Debug().Table("users").Where("id=?", &posts[i].UserID).Take(&posts[i].Sender).Error
			err = GetDB().Debug().Table("images").Where("id=?", &posts[i].Sender.ImageID).Take(&posts[i].Sender.Image).Error
			if err != nil {
				return []Post{}, err
			}
			err = GetDB().Debug().Table("images").Where("id=?", &posts[i].PhotoID).Take(&posts[i].Image).Error
			if err != nil {
				return []Post{}, err
			}
			err = GetDB().Debug().Table("views").Where("post_id=?", &p.ID).Find(&p.Views).Error
			if err != nil {
				return []Post{}, err
			}

		}
	}
	return posts, nil
}
