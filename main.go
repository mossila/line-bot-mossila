package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
)

// VerifyResult From line
type VerifyResult struct {
	Mid       string `json:"mid"`
	ChannelID int64  `json:"channelId"`
}

func verify() *VerifyResult {
	url := "https://api.line.me/v1/oauth/verify"

	req, _ := http.NewRequest("GET", url, nil)
	token := "WTiCfuhDuI3IWEeYdyQsnEiTwG2Hvf/vdQoNd4bE47ZdGLk66tp2WriXq2vxi2VH3/PxUToTFympfewhp5dADEOWDS6GeFHlmw4dFi393T8onJXmDpACsqZSLM49aUkBgl+a9JmevOAkasIWMNE12QdB04t89/1O/w1cDnyilFU="
	req.Header.Add("Authorization", "Bearer "+token)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(string(body))
	result := &VerifyResult{}
	json.Unmarshal(body, result)
	return result
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})
	router.GET("/bot", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello"})
	})
	router.POST("/bot", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello"})
	})

	router.GET("/verify", func(c *gin.Context) {
		c.JSON(200, verify())
	})

	router.Run(":" + port)
}
