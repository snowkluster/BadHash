package controller

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"github.com/gin-gonic/gin"
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