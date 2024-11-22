package main

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"html/template"
	"encoding/json"
)

type VirusTotalResponse struct {
	Data struct {
		Attributes struct {
			TotalVotes struct {
				Harmless int `json:"harmless"`
				Malicious int `json:"malicious"`
			} `json:"total_votes"`
		} `json:"attributes"`
	} `json:"data"`
}

func checkVirusTotalHash(fileHash,virustotalAPIKey string) (string, error) {

	url := fmt.Sprintf("https://www.virustotal.com/api/v3/files/%s", fileHash)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("x-apikey", virustotalAPIKey)
	req.Header.Add("accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var virustotalResponse VirusTotalResponse
	if err := json.Unmarshal(body, &virustotalResponse); err != nil {
		return "", err
	}

	if virustotalResponse.Data.Attributes.TotalVotes.Malicious > 0 {
		return "Malicious", nil
	}
	return "Harmless", nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	var virustotalAPIKey = os.Getenv("API_KEY")
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.LoadHTMLGlob("templates/*")
	r.Static("/public", "./public")
	r.SetTrustedProxies([]string{"localhost"})

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

	r.GET("/checkhash", func(c *gin.Context) {
		tmpl, err := template.ParseFiles("templates/check.html")
		if err != nil {
			log.Fatal(err)
		}
		c.Header("Content-Type", "text/html")
		err = tmpl.Execute(c.Writer, nil)
		if err != nil {
			log.Fatal(err)
		}
	})

	r.POST("/checkhash", func(c *gin.Context) {
		hash := c.DefaultPostForm("hash", "")
		if hash == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "No hash provided",
			})
			return
		}

		status, err := checkVirusTotalHash(hash, virustotalAPIKey)
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
	})

	r.Run(":8080")
}
