package repository

import (
	"github.com/jinzhu/copier"
	"gorm.io/gorm/clause"
)

type Item struct {
	ID        uint  `gorm:"primarykey" json:"id"`
	ChannelID int64 `gorm:"index:idx_scan;not null" json:"channel_id"`

	Name    string `gorm:"not null" json:"name"`
	Regexp  string `gorm:"not null" json:"regexp"`
	Pattern string `gorm:"not null" json:"pattern"`

	MatchPattern string `gorm:"not null" json:"match_pattern"`
	MatchContent string `gorm:"not null" json:"match_content"`

	DateStart int32 `gorm:"not null" json:"date_start"`
	DateEnd   int32 `gorm:"index:idx_date;index:idx_scan;not null" json:"date_end"`

	Process int64 `gorm:"not null" json:"process"`

	Priority   int32  `gorm:"not null" json:"priority"`
	TargetPath string `gorm:"not null" json:"target_path"`
}

type NewItemForm struct {
	ChannelID    int64  `json:"channel_id" form:"channel_id" binding:"required"`
	Name         string `json:"name" form:"name" binding:"required"`
	Regexp       string `json:"regexp" form:"regexp" binding:"required"`
	Pattern      string `json:"pattern" form:"pattern" binding:"required"`
	MatchPattern string `json:"match_pattern" form:"match_pattern" binding:"required"`
	MatchContent string `json:"match_content" form:"match_content" binding:"required"`
	DateStart    int32  `json:"date_start" from:"date_start" binding:"required"`
	DateEnd      int32  `json:"date_end" from:"date_end" binding:"required"`
	Process      int64  `json:"process" form:"process" binding:"required"`
	Priority     int32  `json:"priority" form:"priority" binding:"min=1,max=32"`
	TargetPath   string `json:"target_path" form:"target_path" binding:"required"`
}

type UpdateItemForm struct {
	Name         string `json:"name" form:"name"`
	Regexp       string `json:"regexp" form:"regexp"`
	Pattern      string `json:"pattern" form:"pattern"`
	MatchPattern string `json:"match_pattern" form:"match_pattern" binding:"required"`
	MatchContent string `json:"match_content" form:"match_content" binding:"required"`
	Priority     int32  `json:"priority" form:"priority" binding:"min=1,max=32"`
	TargetPath   string `json:"target_path" form:"target_path"`
}

type ItemRepository struct {
	Repository
}

func (repo ItemRepository) CreateWithForm(form *NewItemForm) (*Item, error) {
	var item Item
	if err := copier.Copy(&item, form); err != nil {
		panic(err)
	}
	return &item, repo.DB.Model(&Item{}).Create(&item).Error
}

func (repo ItemRepository) FirstItemByID(id uint) (*Item, error) {
	var item Item
	return &item, repo.DB.Model(&item).Where("id = ?", id).First(&item).Error
}

func (repo ItemRepository) FirstItemByIDForUpdates(id uint) (*Item, error) {
	var item Item
	return &item, repo.DB.Model(&item).Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", id).First(&item).Error
}

func (repo ItemRepository) GetForUpdates(dateEndAfter int32, channelIds ...int64) ([]Item, error) {
	var items []Item
	tx := repo.DB.Model(&Item{})
	if len(channelIds) != 0 {
		tx = tx.Where("channel_id IN ?", channelIds)
	}
	return items, tx.Where("date_end >= ?", channelIds, dateEndAfter).Clauses(clause.Locking{Strength: "UPDATE"}).Find(&items).Error
}

func (repo ItemRepository) GetActive(dateEnd int32) ([]Item, error) {
	var itemsToDownload []Item
	if err := repo.DB.Model(&Item{}).Where(
		"EXISTS (?)", repo.DB.Model(&Download{}).Where("downloads.item_id = items.id AND downloads.downloaded = ?", false).Limit(1),
	).Find(&itemsToDownload).Error; err != nil {
		return nil, err
	}

	var itemsRecentlyActive []Item
	tx := repo.DB.Model(&Item{})
	if len(itemsToDownload) > 0 {
		var itemsToDownloadIds = make([]uint, len(itemsToDownload))
		for i, item := range itemsToDownload {
			itemsToDownloadIds[i] = item.ID
		}
		tx = tx.Where("id NOT IN (?)", itemsToDownloadIds)
	}
	if err := tx.Where("date_end > ?", dateEnd).Find(&itemsRecentlyActive).Error; err != nil {
		return nil, err
	}

	var result = make([]Item, 0, len(itemsToDownload)+len(itemsRecentlyActive))
	result = append(result, itemsToDownload...)
	result = append(result, itemsRecentlyActive...)
	return result, nil
}

func (repo ItemRepository) GetError() ([]Item, error) {
	var items []Item
	return items, repo.DB.Model(&Item{}).Where(
		"EXISTS (?)", repo.DB.Model(&Download{}).Where(
			"downloads.item_id = items.id",
		).Where(
			"(downloaded = FALSE AND fatal_error = FALSE AND error_at > 0) OR (downloaded = TRUE AND fatal_error = TRUE)",
		),
	).Find(&items).Error
}

func (repo ItemRepository) UpdateProcess(id uint, process int64, dateEnd int32) error {
	model := &Item{
		ID:      id,
		Process: process,
		DateEnd: dateEnd,
	}
	return repo.DB.Model(&model).Select("process", "date_end").Updates(&model).Error
}

func (repo ItemRepository) UpdatesItemWithForm(id uint, form *UpdateItemForm) (bool, error) {
	var item Item
	if err := copier.Copy(&item, form); err != nil {
		panic(err)
	}
	item.ID = id
	result := repo.DB.Model(&item).Updates(&item)
	return result.RowsAffected > 0, result.Error
}

func (repo ItemRepository) DeleteByID(id uint) error {
	return repo.DB.Model(&Item{}).Where("id = ?", id).Delete(nil).Error
}
