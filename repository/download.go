package repository

import (
	"github.com/zelenin/go-tdlib/client"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Download struct {
	ID uint `gorm:"primarykey"`

	ItemID uint  `gorm:"index:idx_item_status;uniqueIndex:idx_item_unique;not null"`
	MsgID  int64 `gorm:"uniqueIndex:idx_item_unique;not null"`
	Text   string
	Size   int64 `gorm:"not null"`
	Date   int32 `gorm:"index:idx_global_queue,priority:4,sort:asc;not null"`

	Priority    int32 `gorm:"index:idx_global_queue,priority:3,sort:desc;not null"`
	Downloading bool  `gorm:"index:idx_global_queue,priority:2;default:false"`
	Downloaded  bool  `gorm:"index:idx_global_queue,priority:1;index:idx_item_status;default:false;not null"`

	FatalError bool `gorm:"index:idx_item_status;default:0;not null"`
	Error      string
	ErrorAt    int64 `gorm:"index:idx_item_status;default:0;not null"`
}

type DownloadWithChannelID struct {
	Download
	ChannelID int64
}

type DownloadTask struct {
	ID          uint         `json:"id"`
	ItemID      uint         `json:"item_id"`
	MsgID       int64        `json:"msg_id"`
	Text        string       `json:"text"`
	Size        int64        `json:"size"`
	Date        int32        `json:"date"`
	Priority    int32        `json:"priority"`
	Downloading bool         `json:"downloading"`
	Downloaded  bool         `json:"downloaded"`
	FatalError  bool         `json:"fatal_error"`
	Error       string       `json:"error"`
	ErrorAt     int64        `json:"error_at"`
	File        *client.File `json:"file,omitempty" gorm:"-"`
}

type DownloadForm struct {
	MsgID    int64 `json:"msg_id" form:"msg_id" binding:"required"`
	Priority int32 `json:"priority" form:"priority" binding:"min=1,max=32"`
}

type DownloadRepository struct {
	Repository
}

func (repo DownloadRepository) joinInnerItems(tx *gorm.DB) *gorm.DB {
	return tx.Joins("INNER JOIN items ON items.id = downloads.item_id")
}

func (repo DownloadRepository) modelWithChannelID(tx *gorm.DB) *gorm.DB {
	tx = tx.Model(&Download{}).Select("downloads.*", "items.channel_id")
	return repo.joinInnerItems(tx)
}

func (repo DownloadRepository) CreateAll(models []Download) error {
	return repo.DB.Create(&models).Error
}

func (repo DownloadRepository) CreateWithMessages(item uint, priority int32, messages []*client.Message) ([]Download, error) {
	var models = make([]Download, 0, len(messages))
	for _, message := range messages {
		videoContent, ok := message.Content.(*client.MessageVideo)
		if !ok {
			continue
		}
		models = append(models, Download{
			ItemID:   item,
			MsgID:    message.Id,
			Text:     videoContent.Caption.Text,
			Size:     videoContent.Video.Video.Size,
			Date:     message.Date,
			Priority: priority,
		})
	}
	if len(models) == 0 {
		return models, nil
	}

	return models, repo.DB.Model(&Download{}).Create(&models).Error
}

func (repo DownloadRepository) FirstByID(id uint) (*Download, error) {
	var download Download
	return &download, repo.DB.Model(&Download{}).Where("id = ?", id).First(&download).Error
}

func (repo DownloadRepository) FirstByIDWithChannelID(id uint) (*DownloadWithChannelID, error) {
	var download DownloadWithChannelID
	tx := repo.modelWithChannelID(repo.DB)
	return &download, tx.Where("downloads.id = ?", id).First(&download).Error
}

func (repo DownloadRepository) FirstByIDForUpdate(id uint) (*Download, error) {
	var download Download
	return &download, repo.DB.Model(&Download{}).Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", id).First(&download).Error
}

func (repo DownloadRepository) GetForDownload(limit *int) ([]DownloadWithChannelID, error) {
	var models []DownloadWithChannelID
	tx := repo.modelWithChannelID(repo.DB)
	tx = tx.Where("downloading=? AND downloaded=?", false, false).Order("downloads.priority DESC,date ASC,downloads.id ASC")
	if limit != nil {
		tx = tx.Limit(*limit)
	}
	return models, tx.Find(&models).Error
}

func (repo DownloadRepository) GetIDByItemForUpdates(itemID uint) ([]uint, error) {
	var ids []uint
	return ids, repo.DB.Model(&Download{}).Clauses(clause.Locking{Strength: "UPDATE"}).Select("id").Where("item_id = ?", itemID).Find(&ids).Error
}

func (repo DownloadRepository) GetDownloading() ([]DownloadWithChannelID, error) {
	var downloads []DownloadWithChannelID
	tx := repo.modelWithChannelID(repo.DB)
	return downloads, tx.Where("downloads.downloaded=? AND downloads.downloading=?", false, true).Find(&downloads).Error
}

func (repo DownloadRepository) GetByID(ids ...uint) ([]DownloadTask, error) {
	var tasks []DownloadTask
	return tasks, repo.DB.Model(&Download{}).Where("id IN ?", ids).Find(&tasks).Error
}

func (repo DownloadRepository) GetByItemID(id uint) ([]DownloadTask, error) {
	var tasks []DownloadTask
	return tasks, repo.DB.Model(&Download{}).Where("item_id = ?", id).Find(&tasks).Error
}

func (repo DownloadRepository) SetDownloading(id uint) error {
	return repo.DB.Model(&Download{}).Where("id=?", id).Update("downloading", true).Error
}

func (repo DownloadRepository) UpdatePriority(id uint, priority int32) (bool, error) {
	result := repo.DB.Model(&Download{}).Where("id=?", id).Update("priority", priority)
	return result.RowsAffected > 0, result.Error
}

func (repo DownloadRepository) UpdateDownloadError(id uint, err string, date int64) error {
	model := Download{
		ID:      id,
		Error:   err,
		ErrorAt: date,
	}
	return repo.DB.Model(&model).Select("error", "error_at").Updates(&model).Error
}

func (repo DownloadRepository) UpdateDownloadFatal(id uint, error string, errorAt int64) error {
	model := Download{
		ID:          id,
		Downloading: false,
		Downloaded:  true,
		FatalError:  true,
		Error:       error,
		ErrorAt:     errorAt,
	}
	return repo.DB.Model(&model).Select(
		"downloading", "downloaded", "fatal_error", "error", "error_at",
	).Updates(&model).Error
}

func (repo DownloadRepository) UpdateResetDownloadState(id uint) (bool, error) {
	model := Download{
		ID:          id,
		Downloading: false,
		FatalError:  false,
		Error:       "",
		ErrorAt:     0,
	}
	result := repo.DB.Model(&model).Select(
		"downloading", "downloaded", "fatal_error", "error", "error_at",
	).Updates(&model)
	return result.RowsAffected > 0, result.Error
}

func (repo DownloadRepository) UpdateDownloadComplete(id uint) error {
	model := Download{
		ID:          id,
		Downloading: false,
		Downloaded:  true,
		FatalError:  false,
	}
	return repo.DB.Model(&model).Select("downloading", "downloaded", "fatal_error").Updates(&model).Error
}

func (repo DownloadRepository) DeleteByID(id uint) (bool, error) {
	result := repo.DB.Model(&Download{}).Where("id=?", id).Delete(nil)
	return result.RowsAffected > 0, result.Error
}

func (repo DownloadRepository) DeleteByItemID(id uint) error {
	return repo.DB.Model(&Download{}).Where("item_id=?", id).Delete(nil).Error
}
