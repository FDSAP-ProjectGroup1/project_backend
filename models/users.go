package models

import "gorm.io/gorm"

type Users struct {
	ID       uint    `gorm:"primary key;autoIncrement" json:"id"`
	Name     *string `json:"name"`
	UserName *string `json:"username"`
	Password *string `json:"password"`
	Address  *string `json:"address"`
}

func MigrateUsers(db *gorm.DB) error {
	err := db.AutoMigrate(&Users{})
	return err
}
