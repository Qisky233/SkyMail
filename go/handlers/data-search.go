package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func SearchHandler(c *gin.Context) {
	// Open the SQLite database
	db, err := sql.Open("sqlite3", "./handlers/data/data.db")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
		return
	}
	defer db.Close()

	// Get query parameters
	queryType := c.Query("type")
	keyword := c.Query("keyword")
	limitStr := c.DefaultQuery("limit", "10")
	pageStr := c.DefaultQuery("page", "1")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
		return
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page parameter"})
		return
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Build the query
	var rows *sql.Rows
	var totalQuery string
	var total int

	if queryType == "file" {
		// Use LIKE for file search
		rows, err = db.Query("SELECT id, file, label, content FROM emails WHERE file LIKE ? LIMIT ? OFFSET ?", "%"+keyword+"%", limit, offset)
		totalQuery = "SELECT COUNT(*) FROM emails WHERE file LIKE ?"
	} else if queryType == "label" {
		// Use exact match for label search
		rows, err = db.Query("SELECT id, file, label, content FROM emails WHERE label = ? LIMIT ? OFFSET ?", keyword, limit, offset)
		totalQuery = "SELECT COUNT(*) FROM emails WHERE label = ?"
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid type parameter"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute query"})
		return
	}
	defer rows.Close()

	// Query the total count of matching emails
	var totalParam interface{}
	if queryType == "file" {
		totalParam = "%" + keyword + "%"
	} else if queryType == "label" {
		totalParam = keyword
	}

	err = db.QueryRow(totalQuery, totalParam).Scan(&total)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query total count"})
		return
	}

	// Parse the results
	var results []map[string]interface{}
	for rows.Next() {
		var id int
		var file, label, content string
		if err := rows.Scan(&id, &file, &label, &content); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse query results"})
			return
		}
		results = append(results, gin.H{
			"id":      id,
			"file":    file,
			"label":   label,
			"content": content,
		})
	}

	// Return the results
	c.JSON(http.StatusOK, gin.H{
		"limit":   limit,
		"page":    page,
		"total":   total,
		"results": results,
	})
}
