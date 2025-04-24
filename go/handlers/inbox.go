package handlers

import (
	"bytes"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/ssh"
)

func GetInboxListHandler(c *gin.Context) {
	files, err := getInboxFiles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"files": files})
}

// getInboxFiles 获取收件箱文件列表及详细信息
func getInboxFiles() ([]map[string]interface{}, error) {
	config := &ssh.ClientConfig{
		User: User,
		Auth: []ssh.AuthMethod{
			ssh.Password(Pwd),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", Host+":22", config)
	if err != nil {
		return nil, fmt.Errorf("SSH 连接失败：%v", err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return nil, fmt.Errorf("创建会话失败：%v", err)
	}
	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b
	cmd := fmt.Sprintf("ls -l %s", MailDir)
	if err := session.Run(cmd); err != nil {
		return nil, fmt.Errorf("执行命令失败：%v", err)
	}

	lines := bytes.Split(b.Bytes(), []byte("\n"))
	files := make([]map[string]interface{}, 0)
	for _, line := range lines {
		if bytes.HasPrefix(line, []byte("总用量")) || len(line) == 0 {
			continue
		}
		fileInfo := parseFileInfo(line)
		if fileInfo != nil {
			files = append(files, fileInfo)
		}
	}
	return files, nil
}

// parseFileInfo 解析文件详细信息
func parseFileInfo(line []byte) map[string]interface{} {
	re := regexp.MustCompile(`^([drwx-]{10})\s+(\d+)\s+(\S+)\s+(\S+)\s+(\d+)\s+(\S+\s+\S+\s+\S+)\s+(\S+)$`)
	match := re.FindSubmatch(line)
	if match == nil {
		return nil
	}

	fileType := string(match[1][0])
	size, _ := strconv.Atoi(string(match[5]))

	// 尝试解析时间
	modTime := parseModTime(string(match[6]))
	if modTime.IsZero() {
		return nil // 如果时间解析失败，跳过该文件
	}

	return map[string]interface{}{
		"name":    string(match[7]),
		"type":    fileType,
		"size":    size,
		"modTime": modTime.Format("2006-01-02 15:04:05"),
	}
}

// parseModTime 解析文件修改时间
func parseModTime(modTimeStr string) time.Time {
	// 中文月份映射
	monthMap := map[string]time.Month{
		"1月":  time.January,
		"2月":  time.February,
		"3月":  time.March,
		"4月":  time.April,
		"5月":  time.May,
		"6月":  time.June,
		"7月":  time.July,
		"8月":  time.August,
		"9月":  time.September,
		"10月": time.October,
		"11月": time.November,
		"12月": time.December,
	}

	// 替换中文月份为英文月份
	for k, v := range monthMap {
		modTimeStr = regexp.MustCompile(k).ReplaceAllString(modTimeStr, v.String())
	}

	// 尝试解析格式 "Jan 2 15:04"
	format := "Jan 2 15:04"
	t, err := time.Parse(format, modTimeStr)
	if err == nil {
		// 如果时间是今年的，直接返回
		if t.Year() == time.Now().Year() {
			return t
		}
		// 如果时间不是今年的，需要调整年份
		t = t.AddDate(time.Now().Year()-t.Year(), 0, 0)
		return t
	}

	// 如果无法解析，返回零值
	return time.Time{}
}
