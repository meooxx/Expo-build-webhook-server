package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	// "net/http"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	log.Println("1\n23")
	r := gin.Default()
	r.POST("/webhook", handleHook)

	r.Run("0.0.0.0:10000")
}
