package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func DataHandler(c *gin.Context) {
	// Open the SQLite database
	db, err := sql.Open("sqlite3", "./handlers/data/data.db")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
		return
	}
	defer db.Close()

	// Get pagination parameters
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	// Query the total count of emails
	var total int
	err = db.QueryRow("SELECT COUNT(*) FROM emails").Scan(&total)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query total count"})
		return
	}

	// Query the email table with pagination
	rows, err := db.Query("SELECT id, file, content, label FROM emails LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query database"})
		return
	}
	defer rows.Close()

	// Parse the results
	var emails []map[string]interface{}
	for rows.Next() {
		var id int
		var file string
		var content string
		var label string
		if err := rows.Scan(&id, &file, &content, &label); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse database results"})
			return
		}
		emails = append(emails, gin.H{"id": id, "file": file, "content": content, "label": label})
	}

	// Return the results as JSON
	c.JSON(http.StatusOK, gin.H{
		"page":    page,
		"limit":   limit,
		"total":   total,
		"results": emails,
	})
}
