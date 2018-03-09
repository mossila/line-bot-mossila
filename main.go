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

type replyObject struct {
	replyToken string
	messages   []ResponseMessage
}
type Message struct {
	Id string
	ResponseMessage
}

type ResponseMessage struct {
	Type string
	Text string
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

// BotPostMessage Line with call with this format
type BotPostMessage struct {
	Events []PostEvent
}

// PostEvent Event inside events array
type PostEvent struct {
	ReplyToken string
	Type       string
	timestamp  int64
	source     Source
	message    Message
}

// Source source of message contain user id
type Source struct {
	Type   string
	UserID string
}

func bot(c *gin.Context) {
	var pm BotPostMessage
	if err := c.BindJSON(&pm); err != nil {
		if pm.Events[0].ReplyToken != "" {
			m := ResponseMessage{Type: "text", Text: "Hello"}
			ms := []ResponseMessage{m}
			c.JSON(200, replyObject{
				replyToken: pm.Events[0].ReplyToken,
				messages:   ms,
			})
		}
	}
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
	router.POST("/bot", bot)

	router.GET("/verify", func(c *gin.Context) {
		c.JSON(200, verify())
	})

	router.Run(":" + port)
}
