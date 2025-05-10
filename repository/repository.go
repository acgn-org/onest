package repository

import "gorm.io/gorm"

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&Item{},
		&Download{},
	)
}

type TypeRepository interface {
	SetDB(*gorm.DB)
}

type Repository struct {
	DB *gorm.DB
}

func (repo *Repository) SetDB(db *gorm.DB) {
	repo.DB = db
}

func (repo *Repository) Rollback() *gorm.DB {
	return repo.DB.Rollback()
}

func (repo *Repository) Commit() *gorm.DB {
	return repo.DB.Commit()
}
