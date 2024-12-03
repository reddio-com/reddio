package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/log"
	bridge_utils "github.com/reddio-com/reddio/bridge/utils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

type gormLogger struct {
	gethLogger log.Logger
}

func (g *gormLogger) LogMode(level logger.LogLevel) logger.Interface {
	return g
}

func (g *gormLogger) Info(_ context.Context, msg string, data ...interface{}) {
	infoMsg := fmt.Sprintf(msg, data...)
	g.gethLogger.Info("gorm", "info message", infoMsg)
}

func (g *gormLogger) Warn(_ context.Context, msg string, data ...interface{}) {
	warnMsg := fmt.Sprintf(msg, data...)
	g.gethLogger.Warn("gorm", "warn message", warnMsg)
}

func (g *gormLogger) Error(_ context.Context, msg string, data ...interface{}) {
	errMsg := fmt.Sprintf(msg, data...)
	g.gethLogger.Error("gorm", "err message", errMsg)
}

func (g *gormLogger) Trace(_ context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rowsAffected := fc()
	g.gethLogger.Debug("gorm", "line", utils.FileWithLineNum(), "cost", elapsed, "sql", sql, "rowsAffected", rowsAffected, "err", err)
}

// InitDB init the db handler
func InitDB(config *Config) (*gorm.DB, error) {
	tmpGormLogger := gormLogger{
		gethLogger: log.Root(),
	}
	var dialector gorm.Dialector
	switch config.DriverName {
	case "mysql":
		dialector = mysql.Open(config.DSN)
	default:
		return nil, fmt.Errorf("unsupported driver: %s", config.DriverName)
	}
	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: &tmpGormLogger,
		NowFunc: func() time.Time {
			return bridge_utils.NowUTC()
		},
	})
	if err != nil {
		return nil, err
	}

	sqlDB, pingErr := Ping(db)
	if pingErr != nil {
		return nil, pingErr
	}

	sqlDB.SetConnMaxLifetime(time.Minute * 10)
	sqlDB.SetConnMaxIdleTime(time.Minute * 5)

	sqlDB.SetMaxOpenConns(config.MaxOpenNum)
	sqlDB.SetMaxIdleConns(config.MaxIdleNum)

	return db, nil
}

// CloseDB close the db handler. notice the db handler only can close when then program exit.
func CloseDB(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	if err := sqlDB.Close(); err != nil {
		return err
	}
	return nil
}

// Ping check db status
func Ping(db *gorm.DB) (*sql.DB, error) {
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	if err = sqlDB.Ping(); err != nil {
		return nil, err
	}
	return sqlDB, nil
}
