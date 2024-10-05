package handlers

import (
	"fileshare/main/helpers"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func GetFiles(c *gin.Context) {
	// Get query parameters
	limit := c.DefaultQuery("limit", "50")
	offset := c.DefaultQuery("offset", "0")
	path := c.DefaultQuery("path", "D://")
	sortOrder := c.DefaultQuery("sort", "asc")
	sortBy := c.DefaultQuery("sortBy", "name") // new query parameter to sort by name or modification date

	// Convert limit and offset to integer
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		limitInt = 10
	}

	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		offsetInt = 0
	}

	fileList, err := helpers.ReadDirectory(path)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Sort files
	if sortBy == "modification" {
		if sortOrder == "desc" {
			sort.Slice(fileList, func(i, j int) bool {
				return fileList[i]["modification"].(time.Time).After(fileList[j]["modification"].(time.Time))
			})
		} else {
			sort.Slice(fileList, func(i, j int) bool {
				return fileList[i]["modification"].(time.Time).Before(fileList[j]["modification"].(time.Time))
			})
		}
	} else {
		if sortOrder == "desc" {
			sort.Slice(fileList, func(i, j int) bool {
				return fileList[i]["name"].(string) > fileList[j]["name"].(string)
			})
		} else {
			sort.Slice(fileList, func(i, j int) bool {
				return fileList[i]["name"].(string) < fileList[j]["name"].(string)
			})
		}
	}

	// Apply offset and limit for lazy loading
	hasMore := false
	if offsetInt < len(fileList) {
		fileList = fileList[offsetInt:]
	} else {
		fileList = []map[string]interface{}{}
	}

	if len(fileList) > limitInt {
		fileList = fileList[:limitInt]
		hasMore = true
	}

	c.JSON(http.StatusOK, gin.H{
		"path":    path,
		"hasMore": hasMore,
		"files":   fileList,
	})
}
