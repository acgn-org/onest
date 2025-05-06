package repository

import "gorm.io/gorm"

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&Item{},
		&Download{},
	)
}

type Repository interface {
	DB() *gorm.DB
}
