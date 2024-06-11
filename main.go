package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	port := os.Getenv("PORT")

	if port == "" {
		port = "3000"
	}

	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	if err := handleDatabaseFiles(); err != nil {
		log.Fatalf("cannot open database: %s\n", err.Error())
	}

	initRoutes(r)

	r.Run(":" + port) // Server listening on port 3000 by default
}
