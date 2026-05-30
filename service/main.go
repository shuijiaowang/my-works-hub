package main

import (
	"log"
	"workhub/config"
	"workhub/db"
	"workhub/routes"
)

func main() {
	config.InitConfig() //读取配置文件
	if err := db.InitDB(); err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}

	// 创建路由
	r := routes.SetupRouter()

	// 启动服务
	if err := r.Run(":9002"); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}
