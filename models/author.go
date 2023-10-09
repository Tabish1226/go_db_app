package models

import "gorm.io/gorm"

type Author struct {
	ID    uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name  *string `json:"username"`
	Books []Book
}

func MigrateAuthors(db *gorm.DB) error {
	err := db.AutoMigrate(&Author{})
	return err
}
