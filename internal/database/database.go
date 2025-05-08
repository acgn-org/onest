package database

import (
	"fmt"
	"github.com/acgn-org/onest/internal/config"
	"github.com/acgn-org/onest/internal/logfield"
	"github.com/acgn-org/onest/repository"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	logger := logfield.New(logfield.ComDatabase)

	conf := &gorm.Config{
		SkipDefaultTransaction: true,
		Logger: &Logger{
			Entry: logger.Entry,
		},
	}

	logger = logger.WithAction("connect")

	var err error
	switch config.Database.Type {
	case "sqlite":
		DB, err = gorm.Open(sqlite.Open(config.Database.DBFile), conf)
	case "mysql":
		DB, err = gorm.Open(mysql.Open(fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?parseTime=True&loc=Local&tls=%s",
			config.Database.User,
			config.Database.Password,
			config.Database.Host,
			config.Database.Port,
			config.Database.Database,
			config.Database.SSLMode,
		)), conf)
	default:
		logger.Fatalf("unsupported database type: %s", config.Database.Type)
	}
	if err != nil {
		logger.Fatalln("failed:", err)
	}

	if err := repository.AutoMigrate(DB); err != nil {
		logger.WithAction("migrate").Fatalln("failed:", err)
	}
}
