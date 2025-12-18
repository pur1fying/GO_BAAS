package database

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pur1fying/GO_BAAS/internal/config"
	"github.com/pur1fying/GO_BAAS/internal/logger"
)

var DB *sqlx.DB

const SqlDriver = "mysql"

var SqlConfig mysql.Config

func Init() error {
	logger.HighLight("Database Init")
	SqlConfig.Net = "tcp"
	SqlConfig.Addr = fmt.Sprintf("%s:%d", config.Config.Database.Host, config.Config.Database.Port)
	SqlConfig.User = config.Config.Database.Username
	SqlConfig.Passwd = config.Config.Database.Password
	SqlConfig.DBName = config.Config.Database.Name
	logger.BAASInfo("Addr   :", SqlConfig.Addr)
	logger.BAASInfo("DBName :", SqlConfig.DBName)
	logger.BAASInfo("User   :", SqlConfig.User)

	var err error
	DB, err = sqlx.Open(SqlDriver, SqlConfig.FormatDSN())
	if err != nil {
		return errors.New("Database Open Error : " + err.Error())
	}

	DB.SetMaxOpenConns(config.Config.Database.MaxOpenConn)
	DB.SetMaxIdleConns(config.Config.Database.MaxIdleConn)
	DB.SetConnMaxLifetime(time.Duration(config.Config.Database.ConnMaxLifetime) * time.Minute)

	if err := DB.Ping(); err != nil {
		return errors.New("Database Ping Error : " + err.Error())
	}

	return nil
}

func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
