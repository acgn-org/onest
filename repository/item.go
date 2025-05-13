package repository

import "gorm.io/gorm/clause"

type Item struct {
	ID        uint  `gorm:"primarykey"`
	ChannelID int64 `gorm:"not null"`

	Name   string `gorm:"not null"`
	Regexp string `gorm:"not null"`

	DateStart int32 `gorm:"not null"`
	DateEnd   int32 `gorm:"not null"`

	Process int64 `gorm:"not null"`

	Priority   int32  `gorm:"not null"`
	TargetPath string `gorm:"not null"`
}

type ItemRepository struct {
	Repository
}

func (repo ItemRepository) GetAllForUpdates() ([]Item, error) {
	var items []Item
	return items, repo.DB.Model(&Item{}).Clauses(clause.Locking{Strength: "UPDATE"}).Find(&items).Error
}

func (repo ItemRepository) UpdateProcess(id uint, process int64, dateEnd int32) error {
	model := &Item{
		ID:      id,
		Process: process,
		DateEnd: dateEnd,
	}
	return repo.DB.Model(&model).Select("process", "date_end").Updates(&model).Error
}
