package main

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"openchat/storage"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func getUserFprt(r http.Request) string {
	userStr := r.Host + r.RemoteAddr
	hash := sha256.Sum256([]byte(userStr))
	userId := base64.StdEncoding.EncodeToString(hash[:])
	//Add User if not exists
	storage.UserExists(userId)
	return userId
}

func getChats() ([]storage.Chat, error) {
	results, err := storage.QueryChats()
	if err != nil {
		return nil, err
	}
	return results, nil
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.LoadHTMLGlob("templates/**/*")
	r.StaticFile("/htmx.min.js", "./lib/htmx.min.js")
	r.StaticFile("/main.css", "./lib/main.css")
	r.StaticFile("/logo.png", "./images/logo.png")
	r.StaticFile("/favicon.ico", "./images/favicon.ico")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{})
	})

	r.GET("/chats", func(c *gin.Context) {
		results, err := getChats()
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
			return
		}

		c.HTML(http.StatusOK, "chats.tmpl", gin.H{
			"chats": results,
		})
	})

	r.POST("/chats", func(c *gin.Context) {
		var chat storage.Chat
		if err := c.ShouldBind(&chat); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		chat.User = getUserFprt(*c.Request)
		storage.AddChat(chat)
		c.Header("HX-Trigger", "newChat")
		c.JSON(http.StatusCreated, gin.H{"message": "Success"})
	})

	r.GET("/:uuid", func(c *gin.Context) {
		chat, err := storage.QueryChat(c.Params.ByName("uuid"))
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
			return
		}

		c.HTML(http.StatusOK, "chat.tmpl", gin.H{
			"data": chat,
		})
	})

	r.GET("/replies/chats/:uuid", func(c *gin.Context) {
		results, err := storage.QueryReplies(c.Params.ByName("uuid"))
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
			return
		}

		chat, err := storage.QueryChat(c.Params.ByName("uuid"))
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
			return
		}

		c.HTML(http.StatusOK, "replies.tmpl", gin.H{
			"replies": results,
			"chat":    chat,
		})
	})

	r.POST("/replies/chats/:uuid", func(c *gin.Context) {
		var chat storage.Chat
		if err := c.ShouldBind(&chat); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		chat.UUID = c.Params.ByName("uuid")
		chat.User = getUserFprt(*c.Request)
		storage.AddReply(chat)
		c.Header("HX-Trigger", "newReplie_"+c.Params.ByName("uuid"))
		c.JSON(http.StatusCreated, gin.H{"message": "Success"})
	})

	r.GET("/upvotes/chats/:uuid", func(c *gin.Context) {
		chat, err := storage.QueryUpvote(c.Params.ByName("uuid"))
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
			return
		}
		c.HTML(http.StatusOK, "upvotes.tmpl", gin.H{"upvotes": chat.Upvotes})
	})

	r.POST("/upvotes/chats/:uuid", func(c *gin.Context) {
		storage.AddUpvote(c.Params.ByName("uuid"), getUserFprt(*c.Request))
		c.Header("HX-Trigger", "newUpvote_"+c.Params.ByName("uuid"))
		c.JSON(http.StatusCreated, gin.H{"message": "Success"})
	})

	return r
}

func main() {
	fmt.Println("PLEASE REMEMBER TO RUN /storage/migrations.sql BEFORE FRESH SETUP")

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	if os.Getenv("APP_ENV") == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := setupRouter()
	if os.Getenv("APP_URL") != "127.0.0.1" {
		r.SetTrustedProxies([]string{os.Getenv("APP_URL")})
	}
	r.Run(":80")
}
