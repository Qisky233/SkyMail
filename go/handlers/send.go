package handlers

import (
	"encoding/base64"
	"fmt"
	"net"

	"github.com/gin-gonic/gin"
)

// 请求参数结构体
type MailRequest struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func SendMailHandler(c *gin.Context) {
	var req MailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	authResult, err := sendMail(Host, Hostname, Username, Password, Sender, SenderName, req.To, req.Subject, req.Body)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error sending mail"})
		return
	}

	if authResult == "Mail sent successfully" {
		c.JSON(200, gin.H{
			"message": "Mail sent successfully",
		})
	} else {
		c.JSON(500, gin.H{
			"message":  "Mail sending failed",
			"response": authResult,
		})
	}
}
func sendMail(host, hostname, username, password, sender, senderName, to, subject, body string) (string, error) {
	conn, err := net.Dial("tcp", host+":25") // 假设SMTP服务端口为25
	if err != nil {
		fmt.Println("Error connecting to SMTP server:", err)
		return "", err
	}
	defer conn.Close()

	// 读取初始连接响应
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading initial connection response:", err)
		return "", err
	}
	fmt.Println("Initial connection response:", string(buffer[:n]))

	// 发送HELO命令
	_, err = conn.Write([]byte("HELO " + hostname + "\r\n"))
	if err != nil {
		fmt.Println("Error sending HELO command:", err)
		return "", err
	}

	// 读取HELO响应
	n, err = conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading HELO response:", err)
		return "", err
	}
	fmt.Println("HELO response:", string(buffer[:n]))

	// 发送AUTH LOGIN命令
	_, err = conn.Write([]byte("AUTH LOGIN\r\n"))
	if err != nil {
		fmt.Println("Error sending AUTH LOGIN command:", err)
		return "", err
	}

	// 读取AUTH LOGIN响应
	n, err = conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading AUTH LOGIN response:", err)
		return "", err
	}
	fmt.Println("AUTH LOGIN response:", string(buffer[:n]))

	// 发送用户名（base64编码）
	encodedUsername := base64.StdEncoding.EncodeToString([]byte(username))
	_, err = conn.Write([]byte(encodedUsername + "\r\n"))
	if err != nil {
		fmt.Println("Error sending username:", err)
		return "", err
	}

	// 读取用户名响应
	n, err = conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading username response:", err)
		return "", err
	}
	fmt.Println("Username response:", string(buffer[:n]))

	// 发送密码（base64编码）
	encodedPassword := base64.StdEncoding.EncodeToString([]byte(password))
	_, err = conn.Write([]byte(encodedPassword + "\r\n"))
	if err != nil {
		fmt.Println("Error sending password:", err)
		return "", err
	}

	// 读取密码响应
	n, err = conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading password response:", err)
		return "", err
	}
	response := string(buffer[:n])
	fmt.Println("Password response:", response)

	// 检查响应代码是否为235（身份验证成功）
	if response[:3] != "235" {
		fmt.Println("Authentication failed:", response)
		return response, nil
	}

	// 发送 MAIL FROM 命令
	_, err = conn.Write([]byte("MAIL FROM:<" + sender + ">\r\n"))
	if err != nil {
		fmt.Println("Error sending MAIL FROM command:", err)
		return "", err
	}

	// 读取 MAIL FROM 响应
	n, err = conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading MAIL FROM response:", err)
		return "", err
	}
	fmt.Println("MAIL FROM response:", string(buffer[:n]))

	// 发送 RCPT TO 命令
	_, err = conn.Write([]byte("RCPT TO:<" + to + ">\r\n"))
	if err != nil {
		fmt.Println("Error sending RCPT TO command:", err)
		return "", err
	}

	// 读取 RCPT TO 响应
	n, err = conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading RCPT TO response:", err)
		return "", err
	}
	fmt.Println("RCPT TO response:", string(buffer[:n]))

	// 发送 DATA 命令
	_, err = conn.Write([]byte("DATA\r\n"))
	if err != nil {
		fmt.Println("Error sending DATA command:", err)
		return "", err
	}

	// 读取 DATA 响应
	n, err = conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading DATA response:", err)
		return "", err
	}
	fmt.Println("DATA response:", string(buffer[:n]))

	// 构建邮件内容
	mailContent := fmt.Sprintf(`Subject: %s
From: "%s" <%s>
To: <%s>

%s
.`, subject, senderName, sender, to, body)

	// 发送邮件内容
	_, err = conn.Write([]byte(mailContent + "\r\n"))
	if err != nil {
		fmt.Println("Error sending mail content:", err)
		return "", err
	}

	// 读取邮件内容响应
	n, err = conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading mail content response:", err)
		return "", err
	}
	fmt.Println("Mail content response:", string(buffer[:n]))

	return "Mail sent successfully", nil
}
