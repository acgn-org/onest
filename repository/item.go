package repository

import "gorm.io/gorm/clause"

type Item struct {
	ID        uint  `gorm:"primarykey"`
	ChannelID int64 `gorm:"index:idx_channel_date;not null"`

	Name    string `gorm:"not null"`
	Regexp  string `gorm:"not null"`
	Pattern string `gorm:"not null"`

	DateStart int32 `gorm:"not null"`
	DateEnd   int32 `gorm:"index:idx_date;index:idx_channel_date;not null"`

	Process int64 `gorm:"not null"`

	Priority   int32  `gorm:"not null"`
	TargetPath string `gorm:"not null"`
}

type ItemRepository struct {
	Repository
}

func (repo ItemRepository) FirstItemByID(id uint) (*Item, error) {
	var item Item
	return &item, repo.DB.Model(&item).Where("id = ?", id).First(&item).Error
}

func (repo ItemRepository) GetAllForUpdates() ([]Item, error) {
	var items []Item
	return items, repo.DB.Model(&Item{}).Clauses(clause.Locking{Strength: "UPDATE"}).Find(&items).Error
}

func (repo ItemRepository) GetWithDateEnd(dateEnd int32) ([]Item, error) {
	var items []Item
	return items, repo.DB.Model(&Item{}).Where("date_end > ?", dateEnd).Find(&items).Error
}

func (repo ItemRepository) UpdateProcess(id uint, process int64, dateEnd int32) error {
	model := &Item{
		ID:      id,
		Process: process,
		DateEnd: dateEnd,
	}
	return repo.DB.Model(&model).Select("process", "date_end").Updates(&model).Error
}
