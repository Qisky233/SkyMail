package handlers

import (
	"app/handlers/model"
	"app/handlers/utils"
	"bytes"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/ssh"
)

// GetMailHandler 处理 GET /mail 请求，通过SSH解析邮件
func DetailHandler(c *gin.Context) {
	// 获取查询参数 fileName
	fileName := c.Query("fileName")
	if fileName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "FileName parameter is required"})
		return
	}

	// 调试信息1：打印接收到的文件名
	fmt.Printf("[DEBUG] 请求文件名: %s\n", fileName)

	// 构造远程邮件文件路径（使用正斜杠）
	remoteMailDir := "/www/vmail/aluo18.top/1212121/cur"
	remoteFilePath := filepath.ToSlash(filepath.Join(remoteMailDir, fileName))

	// 调试信息2：打印远程文件路径
	fmt.Printf("[DEBUG] 远程文件路径: %s\n", remoteFilePath)

	// 创建SSH配置
	config := &ssh.ClientConfig{
		User:            User,
		Auth:            []ssh.AuthMethod{ssh.Password(Pwd)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         15 * time.Second, // 延长超时时间
	}

	// 调试信息3：打印SSH连接信息
	fmt.Printf("[DEBUG] 尝试连接到 %s 用户 %s\n", Host, User)

	// 建立SSH连接
	client, err := ssh.Dial("tcp", Host+":22", config)
	if err != nil {
		errorMsg := fmt.Sprintf("SSH连接失败: %v", err)
		fmt.Printf("[ERROR] %s\n", errorMsg)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   errorMsg,
			"details": "请检查SSH服务是否正常运行",
		})
		return
	}
	defer client.Close()

	// 调试信息4：连接成功
	fmt.Println("[DEBUG] SSH连接成功")

	// 创建会话读取远程文件内容
	session, err := client.NewSession()
	if err != nil {
		errorMsg := fmt.Sprintf("创建SSH会话失败: %v", err)
		fmt.Printf("[ERROR] %s\n", errorMsg)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errorMsg,
		})
		return
	}
	defer session.Close()

	// 捕获输出
	var stdout, stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr

	// 执行命令读取远程文件内容
	readCmd := fmt.Sprintf("cat %s", remoteFilePath)
	if err := session.Run(readCmd); err != nil {
		// 调试信息5：命令执行失败详情
		fmt.Printf("[ERROR] 命令执行失败: %v\n", err)
		fmt.Printf("[ERROR] 标准错误输出: %s\n", stderr.String())
		fmt.Printf("[ERROR] 标准输出: %s\n", stdout.String())

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":    "读取远程文件失败",
			"details":  stderr.String(),
			"output":   stdout.String(),
			"command":  readCmd,
			"filePath": remoteFilePath,
		})
		return
	}

	// 调试信息6：成功读取远程文件内容
	fmt.Printf("[DEBUG] 成功读取远程文件内容，大小: %d bytes\n", len(stdout.String()))

	// 调用 ParseMailHandler 方法处理邮件内容
	mailData, err := utils.ParseMailHandler(stdout.String())
	if err != nil {
		errorMsg := fmt.Sprintf("邮件解析失败: %v", err)
		fmt.Printf("[ERROR] %s\n", errorMsg)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   errorMsg,
			"details": "无法解析邮件内容",
		})
		return
	}

	// 调试信息7：成功解析邮件内容
	fmt.Printf("[DEBUG] 邮件解析成功，返回数据: %v\n", mailData)

	// 合并所有 parts 的 body 内容
	var combinedBody strings.Builder
	for _, part := range mailData.Parts {
		if part.Body != "" {
			combinedBody.WriteString(part.Body)
			// combinedBody.WriteString("\n") // 添加换行符分隔不同部分
		}
	}

	// 调用 TestSpamFilter 方法检测垃圾邮件
	isSpam := model.TestSpamFilter(combinedBody.String())

	// 返回成功结果
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   mailData,
		"meta": gin.H{
			"file":   fileName,
			"path":   remoteFilePath,
			"size":   len(stdout.String()),
			"isSpam": isSpam,
		},
	})
}
