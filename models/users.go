package models

import "gorm.io/gorm"

type Users struct {
	ID       uint    `gorm:"primary key;autoIncrement" json:"id"`
	Fullname *string `json:"fullname"`
	Name     *string `json:"email"`
	Password *string `json:"password"`
	Address  *string `json:"address"`
}

func MigrateUsers(db *gorm.DB) error {
	err := db.AutoMigrate(&Users{})
	return err
}
