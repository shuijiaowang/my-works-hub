package main

import (
	"log"
	"workhub/config"
	"workhub/routes"
)

func main() {
	config.InitConfig() //读取配置文件

	// 创建路由
	r := routes.SetupRouter()

	// 启动服务
	if err := r.Run(":7789"); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}
