package repository

type Download struct {
	ID uint `gorm:"primarykey"`

	ItemID uint `gorm:"index;not null"`
	Item   Item `gorm:"foreignKey:ItemID;constraint:OnDelete:CASCADE"`

	MsgID int64 `gorm:"not null"`
	Text  string
	Size  int32 `gorm:"not null"`
	Date  int32 `gorm:"index:idx_global_queue;not null;priority:3"`

	Priority    int32 `gorm:"not null"`
	Downloading bool  `gorm:"index:idx_global_queue;default:false;priority:2"`
	Downloaded  bool  `gorm:"index:idx_global_queue;default:false;priority:1"`

	FatalError bool
	Error      string
	ErrorAt    int64
}

type DownloadRepository struct {
	Repository
}

func (repo DownloadRepository) CountQueued() (int64, error) {
	var count int64
	return count, repo.DB.Model(&Download{}).Where("downloaded=? AND downloading=?", false, false).Count(&count).Error
}

func (repo DownloadRepository) GetDownloading() ([]Download, error) {
	var downloads []Download
	return downloads, repo.DB.Model(&Download{}).Preload("Item").Where("downloaded=? AND downloading=?", false, true).Find(&downloads).Error
}

func (repo DownloadRepository) UpdateDownloadError(id uint, err string, date int64) error {
	model := Download{
		ID:      id,
		Error:   err,
		ErrorAt: date,
	}
	return repo.DB.Model(&model).Select("error", "error_at").Updates(&model).Error
}

func (repo DownloadRepository) SetDownloading(id uint) error {
	return repo.DB.Model(&Download{}).Where("id=?", id).Update("downloading", true).Error
}
