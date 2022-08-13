package models

import (
	"html"
	"mime/multipart"
	"strings"
	"studapp-blog/api/helpers"

	"github.com/go-playground/validator"
	"github.com/jinzhu/gorm"
)

type File struct {
	File multipart.File `json:"file,omitempty" validate:"required"`
}

type Image struct {
	gorm.Model
	Url string `json:"url,omitempty" validate:"required"`
}

func (img *Image) Prepare() {
	img.ID = 0
	img.Url = html.EscapeString(strings.TrimSpace(img.Url))
}

var (
	validate = validator.New()
)

type mediaUpload interface {
	FileUpload(file File) (string, error)
}

type media struct{}

func NewMediaUpload() mediaUpload {
	return &media{}
}

func (*media) FileUpload(file File) (string, error) {
	err := validate.Struct(file)
	if err != nil {
		return "", err
	}
	uploadUrl, err := helpers.ImageUploadHelper(file.File)
	if err != nil {
		return "", err
	}
	if err != nil {
		return " ", err
	}
	return uploadUrl, nil
}

func (img *Image) SaveImage() (*Image, error) {
	err := db.Debug().Create(img).Error
	if err != nil {
		return &Image{}, err
	}
	return img, nil
}

func (img *Image) DeleteByID(imgid uint) (int64, error) {
	err := db.Debug().Table("images").Where("id = ? ", imgid).Take(&img).Delete(Image{})
	if err.Error != nil {
		return 0, err.Error
	}
	return db.RowsAffected, nil
}

func (img *Image) UpdateImageByID(imgid uint) (*Image, error) {
	db := GetDB().Debug().Table("images").Where("id=?", imgid).UpdateColumns(
		map[string]interface{}{
			"url": img.Url,
		},
	)
	if db.Error != nil {
		return &Image{}, db.Error
	}
	err := GetDB().Debug().Table("images").Where("id=?", imgid).Take(&img).Error
	if err != nil {
		return &Image{}, err
	}
	return img, nil
}

func (img *Image) FindAll() ([]Image, error) {
	var images []Image
	err := db.Debug().Table("images").Find(&images).Error
	if err != nil {
		return nil, err
	}
	return images, nil
}
