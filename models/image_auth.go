package models

import (
	"fmt"
	"github.com/penguinn/penguin/component/db"
	"strings"
)

type ImageAuth struct {
	ID         int    `gorm:"column:id"`
	Token      string `gorm:"column:token"`
	Value      string `gorm:"column:value"`
	DelFlag    int    `gorm:"column:del_flag"`
	CreateTime int    `gorm:"column:create_time"`
	UpdateTime int    `gorm:"column:update_time"`
	AccessTime int    `gorm:"column:access_time"`
}

// TableName sets the insert table name for this struct type
func (i *ImageAuth) TableName() string {
	return "image_auth"
}

func (ImageAuth) ConnectionName() string {
	return "default"
}

//C
func (p ImageAuth) Insert(token string, value string) (err error) {
	imageAuth := ImageAuth{
		Token: token,
		Value: value,
	}
	if p.ValidateTokenIsExist(token) {
		updateMap := map[string]interface{}{
			"value": value,
		}
		err = p.update(token, updateMap)
	} else {
		conn, err := db.WriteModel(p)
		if err != nil {
			return nil
		}
		err = conn.Create(&imageAuth).Error
	}
	return
}

//R
func (p ImageAuth) ValidateAuthCode(token, value string) (bool, error) {
	var imageAuth ImageAuth
	conn, err := db.ReadModel(p)
	if err != nil {
		return false, err
	}
	fmt.Println(token)
	err = conn.Scopes(db.NotDeletedScope).Where("token = ?", token).First(&imageAuth).Error
	if err != nil {
		return false, err
	}
	if strings.ToLower(imageAuth.Value) != strings.ToLower(value) {
		return false, nil
	}
	return true, nil
}

func (p ImageAuth) ValidateTokenIsExist(token string) bool {
	var count int
	conn, err := db.ReadModel(p)
	if err != nil {
		return false
	}
	err = conn.Scopes(db.NotDeletedScope).Where("token = ?", token).Count(&count).Error
	if err != nil {
		return false
	}
	if count < 1 {
		return false
	}
	return true
}

func (p ImageAuth) update(token string, updateMap map[string]interface{}) error {
	conn, err := db.WriteModel(p)
	if err != nil {
		return err
	}
	return conn.Scopes(db.NotDeletedScope).Where("token = ?", token).Update(updateMap).Error
}
