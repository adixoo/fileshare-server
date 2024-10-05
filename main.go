package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"fileshare/main/handlers"
)

type FileInfo struct {
	Name string
	Type string
	Size int64
}

func main() {
	// Create a default Gin router
	router := gin.Default()

	// Define a route for the root path
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to the Gin development server!",
		})
	})

	// Use the separated handler function for the /files route
	router.GET("/files", handlers.GetFiles)

	// Add the new download handler
	router.GET("/download", handlers.DownloadFile)

	// Run the server on port 8080
	router.Run(":8080")
}
