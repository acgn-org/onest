package database

import (
	"context"
	"errors"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"time"
)

type Logger struct {
	Entry *log.Entry
}

func (l *Logger) LogMode(gormlogger.LogLevel) gormlogger.Interface {
	return l
}

func (l *Logger) Info(ctx context.Context, s string, args ...interface{}) {
	l.Entry.WithContext(ctx).Infof(s, args)
}

func (l *Logger) Warn(ctx context.Context, s string, args ...interface{}) {
	l.Entry.WithContext(ctx).Warnf(s, args)
}

func (l *Logger) Error(ctx context.Context, s string, args ...interface{}) {
	l.Entry.WithContext(ctx).Errorf(s, args)
}

func (l *Logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, _ := fc()
	logger := l.Entry.WithContext(ctx).WithFields(log.Fields{
		"elapsed": elapsed,
	})

	if err != nil {
		logger := logger.WithError(err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Traceln(sql)
			return
		}
		logger.Errorln(sql)
		return
	}

	logger.Tracef(sql)
}
