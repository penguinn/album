package models

type ImageAuth struct {
	AccessTime int    `gorm:"column:access_time"`
	CreateTime int    `gorm:"column:create_time"`
	DelFlag    int    `gorm:"column:del_flag"`
	ID         int    `gorm:"column:id"`
	Token      string `gorm:"column:token"`
	UpdateTime int    `gorm:"column:update_time"`
	Value      string `gorm:"column:value"`
}

// TableName sets the insert table name for this struct type
func (i *ImageAuth) TableName() string {
	return "image_auth"
}
