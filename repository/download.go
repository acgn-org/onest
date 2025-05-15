package repository

import "github.com/zelenin/go-tdlib/client"

type Download struct {
	ID uint `gorm:"primarykey"`

	ItemID uint `gorm:"index;not null"`
	Item   Item `gorm:"foreignKey:ItemID;constraint:OnDelete:CASCADE"`

	MsgID int64 `gorm:"not null"`
	Text  string
	Size  int64 `gorm:"not null"`
	Date  int32 `gorm:"index:idx_global_queue;not null;priority:3"`

	Priority    int32 `gorm:"not null"`
	Downloading bool  `gorm:"index:idx_global_queue;default:false;priority:2"`
	Downloaded  bool  `gorm:"index:idx_global_queue;default:false;priority:1"`

	FatalError bool
	Error      string
	ErrorAt    int64
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
	File        *client.File `json:"file,omitempty"`
}

type DownloadRepository struct {
	Repository
}

func (repo DownloadRepository) CountQueued() (int64, error) {
	var count int64
	return count, repo.DB.Model(&Download{}).Where("downloaded=? AND downloading=?", false, false).Count(&count).Error
}

func (repo DownloadRepository) EarliestToDownload(limit int) ([]Download, error) {
	var models []Download
	return models, repo.DB.Model(&Download{}).Where("downloading=? AND downloaded=?", false, false).Order("date ASC").Limit(limit).Find(&models).Error
}

func (repo DownloadRepository) GetDownloading() ([]Download, error) {
	var downloads []Download
	return downloads, repo.DB.Model(&Download{}).Preload("Item").Where("downloaded=? AND downloading=?", false, true).Find(&downloads).Error
}

func (repo DownloadRepository) GetDownloadTaskInfo(tasks []DownloadTask) error {
	return repo.DB.Model(&Download{}).Omit("msg_id", "priority", "fatal_error").Where("id").Find(&tasks).Error
}

func (repo DownloadRepository) CreateWithMessages(item uint, priority int32, messages []*client.Message) ([]Download, error) {
	var models = make([]Download, len(messages))
	var i int
	for _, message := range messages {
		videoContent, ok := message.Content.(*client.MessageVideo)
		if !ok {
			continue
		}

		models[i] = Download{
			ItemID:   item,
			MsgID:    message.Id,
			Text:     videoContent.Caption.Text,
			Size:     videoContent.Video.Video.Size,
			Date:     message.Date,
			Priority: priority,
		}
		i++
	}
	models = models[:i]
	if len(models) == 0 {
		return models, nil
	}

	return models, repo.DB.Model(&Download{}).Create(&models).Error
}

func (repo DownloadRepository) SetDownloading(id uint) error {
	return repo.DB.Model(&Download{}).Where("id=?", id).Update("downloading", true).Error
}

func (repo DownloadRepository) UpdateDownloadError(id uint, err string, date int64) error {
	model := Download{
		ID:      id,
		Error:   err,
		ErrorAt: date,
	}
	return repo.DB.Model(&model).Select("error", "error_at").Updates(&model).Error
}

func (repo DownloadRepository) UpdateDownloadFatal(id uint) error {
	model := Download{
		ID:          id,
		Downloading: false,
		Downloaded:  true,
		FatalError:  true,
	}
	return repo.DB.Model(&model).Select("downloading", "downloaded", "fatal_error").Updates(&model).Error
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
