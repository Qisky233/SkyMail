package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

// HistoryHandler 处理分页查询请求
func HistoryHandler(c *gin.Context) {
	// 获取查询参数
	pageStr := c.Query("page")
	offsetStr := c.Query("limit")

	// 设置默认值
	page := 1
	offset := 10

	// 解析分页参数
	if pageStr != "" {
		page, _ = strconv.Atoi(pageStr)
	}
	if offsetStr != "" {
		offset, _ = strconv.Atoi(offsetStr)
	}

	// 计算分页参数
	limit := offset
	skip := (page - 1) * offset

	// 查询数据库获取数据
	db, err := sql.Open("sqlite3", "./handlers/data/data.db")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "msg": "Failed to connect to database"})
		return
	}
	defer db.Close()

	// 查询总记录数
	var total int
	err = db.QueryRow("SELECT COUNT(*) FROM spam").Scan(&total)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "msg": "Failed to query total records"})
		return
	}

	// 查询分页数据
	rows, err := db.Query("SELECT id, content, time, result FROM spam ORDER BY id DESC LIMIT ? OFFSET ?", limit, skip)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "msg": "Failed to query data"})
		return
	}
	defer rows.Close()

	// 解析查询结果
	var data []map[string]interface{}
	for rows.Next() {
		var record struct {
			ID      int
			Content string
			Time    time.Time
			Result  string
		}
		if err := rows.Scan(&record.ID, &record.Content, &record.Time, &record.Result); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "msg": "Failed to scan data"})
			return
		}
		data = append(data, map[string]interface{}{
			"id":      record.ID,
			"content": record.Content,
			"time":    record.Time.Format("2006-01-02 15:04:05"),
			"result":  record.Result,
		})
	}

	// 返回数据
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": gin.H{
			"total": total,
			"items": data,
		},
	})
}
