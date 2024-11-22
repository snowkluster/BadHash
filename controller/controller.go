package controller

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"github.com/gin-gonic/gin"
	"log"
	"hash-api/utils"
	"github.com/joho/godotenv"
	"os"
)


func UploadFile(c *gin.Context) {
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
}

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