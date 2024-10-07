package main

import (
	"log"

	"github.com/gin-contrib/cors"
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
	gin.SetMode(gin.ReleaseMode)
	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true // This allows all origins. For production, you might want to restrict this.
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}

	// Use CORS middleware
	router.Use(cors.New(config))

	// Serve static files from the "static/dist" directory
	router.Static("/static", "./static/dist")

	for _, route := range handlers.Routes {
		router.GET(route, handlers.SendClientHtml)
	}

	// Load HTML templates
	router.LoadHTMLGlob("static/template/*")

	// Use the separated handler function for the /files route
	router.GET("/files", handlers.GetFiles)
	// Use the separated handler function for the /files route
	router.GET("/default-files", handlers.GetDefaultFiles)

	// Add the new download handler
	router.GET("/download", handlers.DownloadFile)

	// Determine the port based on the environment
	port := "8080"
	if gin.Mode() == gin.ReleaseMode {
		port = "80"
	}

	// Log the server start message
	log.Printf("Server is starting on port %s", port)

	// Run the server on the determined port
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
