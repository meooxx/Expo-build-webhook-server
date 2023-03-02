package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

const (
	GIN_MODE    = "GIN_MODE"
	GIN_DEBUG   = "debug"
	GIN_RELEASE = "release"
)

func init() {
	os.Setenv("TZ", "Asia/Shanghai")
}
func main() {
	if os.Getenv(GIN_MODE) != GIN_RELEASE {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	log.Println(os.Getenv(GIN_MODE))
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})
	r.POST("/webhook", handleHook)

	r.Run("0.0.0.0:10000")
}
