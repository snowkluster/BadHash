package main

import (
	"github.com/gin-gonic/gin"
	"hash-api/controller"
	"github.com/joho/godotenv"
	"os"
	"log"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	var port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.LoadHTMLGlob("templates/*")
	r.Static("/public", "./public")
	r.SetTrustedProxies([]string{"localhost"})

	r.GET("/", controller.Index)
	r.POST("/upload", controller.UploadFile)
	r.GET("/checkhash", controller.Hashpage)
	r.POST("/checkhash", controller.Viruscheck)
	r.Run(":" + port)
}
