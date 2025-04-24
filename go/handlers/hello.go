// handlers/hello.go
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// WelcomeHandler 处理 /hello 接口的请求
func WelcomeHandler(c *gin.Context) {
	// 判断请求方法
	if c.Request.Method == "GET" {
		// GET 请求返回 "运行成功" 字样
		c.String(http.StatusOK, "运行成功")
	} else if c.Request.Method == "POST" {
		// POST 请求返回 "运行成功" 的 JSON 格式响应数据
		c.JSON(http.StatusOK, gin.H{
			"message": "运行成功",
		})
	} else {
		// 其他请求方法返回 405 Method Not Allowed
		c.Status(http.StatusMethodNotAllowed)
	}
}
