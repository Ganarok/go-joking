package main

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Joke struct {
	Setup    string
	Delivery string
}

func initRoutes(router *gin.Engine) {
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello World!",
		})
	})

	router.GET("/joke", func(c *gin.Context) {
		res, getErr := http.Get("https://v2.jokeapi.dev/joke/Any?lang=fr")

		if getErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": getErr.Error(),
			})
			return
		}

		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		var joke Joke

		jsonErr := json.Unmarshal(body, &joke)

		if jsonErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": jsonErr.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"setup":    joke.Setup,
			"delivery": joke.Delivery,
		})
	})
}
