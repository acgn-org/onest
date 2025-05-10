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

	databaseConfig := config.Database.Get()

	var err error
	switch databaseConfig.Type {
	case "sqlite":
		DB, err = gorm.Open(sqlite.Open(databaseConfig.DBFile), conf)
	case "mysql":
		DB, err = gorm.Open(mysql.Open(fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?parseTime=True&loc=Local&tls=%s",
			databaseConfig.User,
			databaseConfig.Password,
			databaseConfig.Host,
			databaseConfig.Port,
			databaseConfig.Database,
			databaseConfig.SSLMode,
		)), conf)
	default:
		logger.Fatalf("unsupported database type: %s", databaseConfig.Type)
	}
	if err != nil {
		logger.Fatalln("failed:", err)
	}

	if err := repository.AutoMigrate(DB); err != nil {
		logger.WithAction("migrate").Fatalln("failed:", err)
	}
}

func Begin() *gorm.DB {
	return DB.Begin()
}

func NewRepository[T any]() T {
	var repo T
	any(&repo).(repository.TypeRepository).SetDB(DB)
	return repo
}

func BeginRepository[T any]() T {
	var repo T
	any(&repo).(repository.TypeRepository).SetDB(Begin())
	return repo
}
