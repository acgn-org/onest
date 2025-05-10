package queue

import (
	"github.com/acgn-org/onest/internal/database"
	"github.com/acgn-org/onest/repository"
)

func setDownloadError(id uint, isFatal bool, msg string, date int64) error {
	downloadRepo := database.BeginRepository[repository.DownloadRepository]()
	defer downloadRepo.Rollback()

	err := downloadRepo.UpdateDownloadError(id, isFatal, msg, date)
	if err != nil {
		return err
	}
	return tx.Commit().Error
	return downloadRepo.Commit().Error
}
