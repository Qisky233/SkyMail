package model

import (
	"log"
)

// TestSpamFilter 测试垃圾邮件过滤器
// 参数: content string - 邮件内容
// 返回: bool - true=垃圾邮件, false=正常邮件
func TestSpamFilter(content string) bool {
	// 初始化垃圾邮件过滤器
	filter, err := NewSpamFilter("./handlers/model/spam_filter.gonb")
	if err != nil {
		log.Fatalf("Failed to initialize spam filter: %v", err)
	}

	// 检测邮件内容是否为垃圾邮件
	isSpam := filter.Predict(content, 1.5) // 使用默认阈值
	return isSpam
}
