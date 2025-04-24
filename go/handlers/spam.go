package handlers

import (
	"app/handlers/model"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func SpamHandler(c *gin.Context) {
	var requestBody struct {
		Content string `json:"content"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	content := requestBody.Content

	// Open the SQLite database
	db, err := sql.Open("sqlite3", "./handlers/data/data.db")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
		return
	}
	defer db.Close()

	// Create the spam table if it doesn't exist
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS spam (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		content TEXT NOT NULL,
		time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		result TEXT NOT NULL
	);
	`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create table"})
		return
	}

	// Use model.TestSpamFilter to process the content
	isSpam := model.TestSpamFilter(content)
	result := "normal"
	if isSpam {
		result = "spam"
	}

	// Get the current time in Shanghai timezone
	shanghaiTime := time.Now().In(time.FixedZone("CST", 8*3600)).Format("2006-01-02 15:04:05")

	// Insert the content and result into the spam table
	insertSQL := "INSERT INTO spam (content, time, result) VALUES (?, ?, ?)"
	_, err = db.Exec(insertSQL, content, shanghaiTime, result)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert data"})
		return
	}

	fmt.Println("Received content:", content)
	c.JSON(http.StatusOK, gin.H{"status": "success", "content": content, "result": result})
}
