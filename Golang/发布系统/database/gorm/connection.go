package database

import (
	"fmt"
	"log"
	"os"
	"path"

	"happyladysauce/utils/conf"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var PostDB *gorm.DB

func InitGormDB() {
    host := config.Config.GetString("mysql.host")
    port := config.Config.GetString("mysql.port")
    user := config.Config.GetString("mysql.user")
    password := config.Config.GetString("mysql.password")
    dbname := config.Config.GetString("mysql.dbname")
    // 读取日志文件名（与 conf/db.yaml 的 mysql.log 对应）
    logFileName := config.Config.GetString("mysql.log")
    if logFileName == "" {
        // 兜底，避免仅传入目录导致 openFile 试图打开目录
        logFileName = "post.db.log"
    }

    // 日志控制
    // 确保日志目录存在
    if err := os.MkdirAll("./log", 0755); err != nil {
        panic(err)
    }
    logFile, err := os.OpenFile(path.Join("./log", logFileName), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        panic(err)
    }
	newLogger := logger.New(
		log.New(logFile, "gorm:", log.LstdFlags),
		logger.Config{
			SlowThreshold:             200,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	DSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, dbname)

	db, err := gorm.Open(mysql.Open(DSN), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}

	PostDB = db
}