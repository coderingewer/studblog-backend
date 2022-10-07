package models

import "github.com/jinzhu/gorm"

type FavoritesList struct {
	gorm.Model
	Name       string `gorm:"name" json:"name"`
	Hide       bool   `gorm:"hide" json:"hide"`
	CoverImgId uint   `gorm:"cover_img_id" json:"coverImgId"`
	CoverImage Image  `gorm:"cover_image" json:"coverImage"`
	Posts      []Post `gorm:"posts" json:"posts"`
	UserId     uint   `gorm:"user_id" json:"userId"`
}

type Favorite struct {
	gorm.Model
	ListId uint `gorm:"list_id" json:"listId"`
	PostId uint `gorm:"post_id" json:"postId"`
	UserId uint `gorm:"user_id" json:"userId"`
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
	err := GetDB().Table("favarite").Where("id=?", favid).Take(&fav).Delete(Favorite{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (favs FavoritesList) FavoritesList(favid uint) error {
	err := GetDB().Table("favarites").Where("id=?", favid).Take(&favs).Delete(FavoritesList{}).Error
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
			"hide": favs.Hide,
			"name": favs.Name,
		},
	)
	if err.Error != nil {
		return FavoritesList{}, err.Error
	}
	return favs, nil
}
