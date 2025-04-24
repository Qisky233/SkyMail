package main

import (
	"app/handlers"
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// gin.SetMode(gin.ReleaseMode) // 设置为生产环境模式
	r := gin.Default()

	// 配置 CORS
	r.Use(cors.Default())

	// 提供静态文件服务
	r.Static("/static", "./")

	// 定义路由
	r.GET("/", handlers.WelcomeHandler)
	// 获取邮件列表
	r.GET("/v1/inbox", handlers.GetInboxListHandler)
	// 发送邮件
	r.POST("/v1/send", handlers.SendMailHandler)
	// 解析邮件
	// r.GET("/v1/parse", handlers.ParseMailHandler)
	// 获取邮件详情
	r.GET("/v1/detail", handlers.DetailHandler)
	// 获取数据集分页内容
	r.GET("/v1/data", handlers.DataHandler)
	// 搜索数据集
	r.GET("/v1/data/search", handlers.SearchHandler)
	// 获取数据集详情
	r.POST("/v1/data/detail", handlers.DataDetailHandler)
	// 垃圾邮件检测
	r.POST("/v1/spam", handlers.SpamHandler)
	// 获取检测历史
	r.GET("/v1/spam/history", handlers.HistoryHandler)

	port := ":7010"
	fmt.Printf("Server running at http://localhost%s\n", port)
	r.Run(port)
}
