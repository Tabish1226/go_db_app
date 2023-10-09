package models

import "gorm.io/gorm"

type Book struct {
	ID        uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	Title     *string `json:"title"`
	Publisher *string `json:"publisher"`
	Author    Author
}

func MigrateBooks(db *gorm.DB) error {
	err := db.AutoMigrate(&Book{})
	return err
}
