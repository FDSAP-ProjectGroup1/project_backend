package models

import "gorm.io/gorm"

type Scheds struct {
	ID     uint    `gorm:"primary key;autoIncrement" json:"id"`
	Title  *string `json:"title"`
	Reason *string `json:"reason"`
}

func MigrateSched(db *gorm.DB) error {
	err := db.AutoMigrate(&Scheds{})
	return err
}
