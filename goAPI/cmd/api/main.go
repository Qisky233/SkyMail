package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-ole/go-ole" // 添加此行
	"skyAPI/internal/app/hmail/handler"
	"skyAPI/internal/config"
)

func main() {
	// 初始化 COM 库（必须第一行！）
	if err := ole.CoInitialize(0); err != nil {
		log.Fatal("COM 初始化失败:", err)
	}
	defer ole.CoUninitialize() // 确保卸载

	// 加载配置
	cfg, err := config.LoadConfig("configs/config.yaml")
	if err != nil {
		log.Fatal("加载配置失败:", err)
	}

	// 初始化 Gin
	router := gin.Default()

	// 注册路由
	mailHandler := handler.NewMailHandler(cfg)
	router.GET("/api/version", mailHandler.GetVersion)
	router.POST("/api/domains", mailHandler.CreateDomain)

	// 启动服务
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("服务启动在 %s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal("服务启动失败:", err)
	}
}