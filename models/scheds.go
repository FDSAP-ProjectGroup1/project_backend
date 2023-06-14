package models

import "gorm.io/gorm"

type Scheds struct {
	ID     uint    `gorm:"primary key;autoIncrement" json:"id"`
	Date   *string `json:"date"`
	Time   *string `json:"time"`
	Title  *string `json:"title"`
	Reason *string `json:"reason"`
}

func MigrateScheds(db *gorm.DB) error {
	err := db.AutoMigrate(&Scheds{})
	return err
}
