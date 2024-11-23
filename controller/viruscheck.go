package controller

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"log"
	"hash-api/utils"
	"github.com/joho/godotenv"
	"os"
)

func Viruscheck(c *gin.Context) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	var virustotalAPIKey = os.Getenv("API_KEY")
	hash := c.DefaultPostForm("hash", "")
	if hash == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No hash provided",
		})
		return
	}

	status, err := utils.CheckVirusTotalHash(hash, virustotalAPIKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to check VirusTotal",
		})
		return
	}

	c.HTML(http.StatusOK, "hash_result.html", gin.H{
		"hash":   hash,
		"status": status,
	})
}