package main

import (
	"github.com/gin-gonic/gin"

	"happyladysauce/database/gorm"
	"happyladysauce/utils/conf"
	"happyladysauce/utils/logger"
)

func Init() {
	config.InitConfig("./conf", "db", config.YAML)
	database.InitGormDB()
	logger.InitLogger()
}

func main() {
	// 初始化路由
	router := gin.Default()
	// 注册路由

	

	// 启动服务
	router.Run(":8080")
}