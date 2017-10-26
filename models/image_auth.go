package models

import "github.com/penguinn/penguin/component/db"

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
func (p ImageAuth) Insert(token string, value string) error {
	imageAuth := ImageAuth{
		Token:token,
		Value:value,
	}
	conn, err := db.WriteModel(p)
	if err != nil {
		return nil
	}
	err = conn.Table(p.TableName()).Create(&imageAuth).Error
	return err
}

//R
func (p ImageAuth) ValidateAuthCode(token, value string) (bool, error) {
	var imageAuth ImageAuth
	conn, err := db.ReadModel(p)
	if err != nil {
		return false, err
	}
	err = conn.Table(p.TableName()).Where("token = ?", token).First(&imageAuth).Error
	if err != nil {
		return false, err
	}
	if imageAuth.Value != value{
		return false, nil
	}
	return true, nil
}

