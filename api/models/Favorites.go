package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type FavoritesList struct {
	gorm.Model
	Name       string     `gorm:"name" json:"name"`
	Hide       bool       `gorm:"hide" json:"hide"`
	CoverImgId uint       `gorm:"cover_img_id" json:"coverImgId"`
	CoverImage Image      `gorm:"cover_image" json:"coverImage"`
	UserId     uint       `gorm:"user_id" json:"userId"`
	User       User       `gorm:"user" json:"user"`
	Items      []Favorite `gorm:"items" json:"items"`
}

type Favorite struct {
	gorm.Model
	ListId uint `gorm:"list_id" json:"listId"`
	PostId uint `gorm:"post_id" json:"postId"`
	UserId uint `gorm:"user_id" json:"userId"`
	Post   Post `gorm:"post" json:"post"`
}

func (favs FavoritesList) CreateAList() (FavoritesList, error) {
	err := GetDB().Create(&favs).Error
	if err != nil {
		return FavoritesList{}, err
	}
	return favs, nil
}

func (fav Favorite) AddPostToList() (Favorite, error) {
	err := GetDB().Create(&fav).Error
	if err != nil {
		return Favorite{}, err
	}
	return fav, nil
}

func (fav Favorite) RemoveFromList(favid uint) error {
	err := GetDB().Table("favorites").Where("id=?", favid).Take(&fav).Delete(Favorite{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (favs *FavoritesList) DeleteList(favsid uint) error {
	err := GetDB().Table("favorites_lists").Where("id=?", favsid).Take(&favs).Delete(&FavoritesList{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (favs FavoritesList) FindFavsListsByUserID(userId uint) ([]FavoritesList, error) {
	lists := []FavoritesList{}
	err := GetDB().Table("favorites_lists").Where("user_id=?", userId).Take(&lists).Error
	if err != nil {
		return []FavoritesList{}, err
	}
	return lists, nil

}

func (favs FavoritesList) UpdateFavList(favsId uint) (FavoritesList, error) {
	err := GetDB().Table("favorites_lists").Where("id=?", favsId).UpdateColumns(
		map[string]interface{}{
			"hide":       favs.Hide,
			"name":       favs.Name,
			"updated_at": time.Now(),
		},
	)
	if err.Error != nil {
		return FavoritesList{}, err.Error
	}
	return favs, nil
}

func (favs *FavoritesList) FindByID(id uint) (*FavoritesList, error) {
	err := GetDB().Table("favorites_lists").Where("id = ?", id).Take(&favs).Error
	if err != nil {
		return &FavoritesList{}, err
	}
	err = GetDB().Table("favorites").Where("list_id = ?", id).Find(&favs.Items).Error
	if err != nil {
		return &FavoritesList{}, err
	}
	if len(favs.Items) > 0 {
		for i, _ := range favs.Items {
			err = GetDB().Table("posts").Where("id = ?", favs.Items[i].PostId).Find(&favs.Items[i].Post).Error
			err = GetDB().Table("images").Where("id = ?", favs.Items[i].Post.PhotoID).Find(&favs.Items[i].Post.Image).Error
			err = GetDB().Table("users").Where("id = ?", favs.Items[i].Post.UserID).Find(&favs.Items[i].Post.Sender).Error
			err = GetDB().Table("images").Where("id = ?", favs.Items[i].Post.Sender.ImageID).Find(&favs.Items[i].Post.Sender.Image).Error
		}
	}
	return favs, nil
}

func (fav *Favorite) FindByID(id uint) (*Favorite, error) {
	err := GetDB().Table("favorites").Where("id = ?", id).Take(&fav).Error
	if err != nil {
		return &Favorite{}, err
	}
	return fav, nil
}

func (favs *Favorite) FinByListID(listId uint) (*[]Favorite, error) {
	items := []Favorite{}
	err := GetDB().Table("favorites").Where("list_id = ?", listId).Find(&items).Error
	if err != nil {
		return &[]Favorite{}, err
	}
	return &items, nil
}
