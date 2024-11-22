package controller

import (
	"github.com/gin-gonic/gin"
	"html/template"
	"log"
)

func Hashpage(c *gin.Context) {
	tmpl, err := template.ParseFiles("templates/check.html")
	if err != nil {
		log.Fatal(err)
	}
	c.Header("Content-Type", "text/html")
	err = tmpl.Execute(c.Writer, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func Index(c *gin.Context) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Fatal(err)
	}
	c.Header("Content-Type", "text/html")
	err = tmpl.Execute(c.Writer, nil)
	if err != nil {
		log.Fatal(err)
	}
}