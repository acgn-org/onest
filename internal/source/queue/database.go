package queue

import (
	"github.com/acgn-org/onest/internal/database"
	"github.com/acgn-org/onest/repository"
)

func setDownloadError(id uint, isFatal bool, msg string, date int64) error {
	tx := database.Begin()
	defer tx.Rollback()

	downloadRepo := repository.DownloadRepository{DB: tx}
	err := downloadRepo.UpdateDownloadError(id, isFatal, msg, date)
	if err != nil {
		return err
	}
	return tx.Commit().Error
}
