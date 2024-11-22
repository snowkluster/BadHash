package main

import (
	"github.com/gin-gonic/gin"
	"hash-api/controller"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.LoadHTMLGlob("templates/*")
	r.Static("/public", "./public")
	r.SetTrustedProxies([]string{"localhost"})

	r.GET("/", controller.Index)
	r.POST("/upload", controller.UploadFile)
	r.GET("/checkhash", controller.Hashpage)
	r.POST("/checkhash", controller.Viruscheck)
	r.Run(":8080")
}
