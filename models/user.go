package models

import (
	"github.com/penguinn/penguin/component/db"
)

type User struct {
	ID         int    `gorm:"column:id"`
	Username   string `gorm:"column:username"`
	Password   string `gorm:"column:password"`
	DelFlag    int    `gorm:"column:del_flag"`
	AccessTime int64  `gorm:"column:access_time"`
	CreateTime int64  `gorm:"column:create_time"`
	UpdateTime int    `gorm:"column:update_time"`
}

// TableName sets the insert table name for this struct type
func (u *User) TableName() string {
	return "user"
}

func (User) ConnectionName() string {
	return "default"
}

func (p User) Insert(username, password string) (int, error) {
	user := User{
		Username: username,
		Password: password,
	}
	conn, err := db.ReadModel(p)
	if err != nil {
		return 0, err
	}
	err = conn.Table(p.TableName()).Create(&user).Error
	return user.ID, err
}

//R
func (p User) ValidateUsername(username string) (bool, error) {
	var user User
	conn, err := db.ReadModel(p)
	if err != nil {
		return false, err
	}
	err = conn.Table(p.TableName()).Where("username = ?", username).First(&user).Error
	if err != nil {
		return false, err
	}
	if user.ID == 0 {
		return false, nil
	}
	return true, nil
}

func (p User) SelectByUsername(username string) (*User, error) {
	var user User
	conn, err := db.ReadModel(p)
	if err != nil {
		return nil, err
	}
	err = conn.Table(p.TableName()).Scopes(db.NotDeletedScope).Where("username = ?", username).Find(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (p User) SelectByUserID(userID int) (*User, error) {
	var user User
	conn, err := db.ReadModel(p)
	if err != nil {
		return nil, err
	}
	err = conn.Table(p.TableName()).Scopes(db.NotDeletedScope).Where("id = ?", userID).Find(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

//U
func (p User) Update(id int, updateMap map[string]interface{}) error {
	conn, err := db.WriteModel(p)
	if err != nil {
		return err
	}

	err = conn.Table(p.TableName()).Where("id = ?", id).Update(updateMap).Error
	return err
}
