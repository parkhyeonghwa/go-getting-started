package main

import (
        "fmt"
        "log"
        "net/http"
        "net/url"
        "os"

        "github.com/gin-gonic/gin"
        "github.com/line/line-bot-sdk-go/linebot"
)

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

        router.POST("/callback", func(c *gin.Context) {
                proxyURL, _ := url.Parse(os.Getenv("http://fixie:qkRJC1ZNBKlIZn9@velodrome.usefixie.com:80"))
                client := &http.Client{
                        Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)},
                }
                bot, err := linebot.NewClient(1473280524, "da986f754c64f8d54fb82074437b07a2", "u4b22e97d92c92e21ba3c4750cc5f95ca", linebot.WithHTTPClient(client))
                if err != nil {
                        fmt.Println(err)
                        return
                }

                received, err := bot.ParseRequest(c.Request)
                if err != nil {
                        if err == linebot.ErrInvalidSignature {
                                fmt.Println(err)
                        }
                        return
                }
                for _, result := range received.Results {
                        content := result.Content()
                        if content != nil && content.IsMessage && content.ContentType == linebot.ContentTypeText {
                                text, err := content.TextContent()
                                res, err := bot.SendText([]string{content.From}, "OK "+text.Text)
                                if err != nil {
                                        fmt.Println(res)
                                }
                        }
                }
        })

        router.Run(":" + port)
}