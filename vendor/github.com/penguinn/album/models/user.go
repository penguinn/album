package models

import "github.com/penguinn/penguin/component/db"

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

//U
func (p User) Update(id int, updateMap map[string]interface{}) error {
	conn, err := db.ReadModel(p)
	if err != nil {
		return err
	}

	err = conn.Table(p.TableName()).Where("id = ?", id).Update(updateMap).Error
	return err
}