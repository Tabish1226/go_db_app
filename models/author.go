package models

import "gorm.io/gorm"

type Author struct {
	ID    uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name  *string `json:"name"`
	Books []Book  `gorm:"foreignKey:AuthorID"`
}

func MigrateAuthors(db *gorm.DB) error {
	err := db.AutoMigrate(&Author{})
	return err
}
