package main

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
	"html/template"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.LoadHTMLGlob("templates/*")

	r.Static("/public", "./public")
	r.SetTrustedProxies([]string{"localhost"})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/", func(c *gin.Context) {
		tmpl, err := template.ParseFiles("templates/index.html")
		if err != nil {
			log.Fatal(err)
		}
		c.Header("Content-Type", "text/html")
		err = tmpl.Execute(c.Writer, nil)
		if err != nil {
			log.Fatal(err)
		}
	})

	r.POST("/upload", func(c *gin.Context) {
		file, _, err := c.Request.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to read file",
			})
			return
		}
		defer file.Close()

		md5Hash := md5.New()
		sha256Hash := sha256.New()

		_, err = io.Copy(md5Hash, file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to calculate MD5 hash",
			})
			return
		}

		file.Seek(0, 0)
		_, err = io.Copy(sha256Hash, file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to calculate SHA256 hash",
			})
			return
		}

		md5HashString := hex.EncodeToString(md5Hash.Sum(nil))
		sha256HashString := hex.EncodeToString(sha256Hash.Sum(nil))

		c.HTML(http.StatusOK, "results.html", gin.H{
			"md5":    md5HashString,
			"sha256": sha256HashString,
		})
	})
	r.Run(":8080")
}
