package handlers

import (
	"app/handlers/utils"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DataDetailHandler(c *gin.Context) {
	// 读取请求体内容
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read request body"})
		return
	}

	// 将请求体内容转换为字符串
	mailContent := string(body)

	fmt.Println("Received mail content:" + mailContent)

	// 解析邮件内容
	emailResponse, err := utils.ParseMailHandler(mailContent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回解析结果
	c.JSON(http.StatusOK, emailResponse)
}
