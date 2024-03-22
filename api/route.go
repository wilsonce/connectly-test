package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/wilsonce/connectly-test/api/dto"
	"github.com/wilsonce/connectly-test/initialize"
	"github.com/wilsonce/connectly-test/model"
	"github.com/wilsonce/connectly-test/service"
	"github.com/wilsonce/connectly-test/tg"
	"net/http"
)

const JwtSecret = "Zh1w*Vu3xeR@4JBZC_d#"

func InitRoute(e *gin.Engine) {
	e.POST("/login", func(c *gin.Context) {
		var user dto.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if user.Username == "admin" && user.Password == "123456" {
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"username": user.Username,
			})
			s, err := token.SignedString([]byte(JwtSecret))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": "Login ok", "token": s})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "Login error"})
	})

	botGroup := e.Group("/bot", Auth())
	botGroup.POST("/add", func(c *gin.Context) {
		var bot dto.Bot
		if err := c.ShouldBindJSON(&bot); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		botBuilder := tg.NewBotBuilder()
		tmpBot := botBuilder.Build(bot.Token)
		if tmpBot == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bot add error"})
			return
		}
		botService := service.NewBotService()
		botM := &model.WBot{
			BotToken: bot.Token,
			BotName:  bot.BotName,
		}
		err := botService.Add(context.Background(), botM)
		if err != nil {
			initialize.Logger.Error(err.Error())
		}
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	botGroup.POST("/stop", func(c *gin.Context) {
		var bot dto.Bot
		if err := c.ShouldBindJSON(&bot); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		botBuilder := tg.NewBotBuilder()
		botBuilder.Remove(bot.Token)
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	e.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
}
