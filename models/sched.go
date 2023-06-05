package models

import (
	"gorm.io/gorm"
)

type Scheds struct {
	ID     uint    `gorm:"SERIAL PRIMARY KEY" json:"id"`
	Date   *string `json:"date"`
	Time   *string `json:"time"`
	Title  *string `json:"title"`
	Reason *string `json:"reason"`
}

func MigrateSched(db *gorm.DB) error {
	err := db.AutoMigrate(&Scheds{})
	return err
}
