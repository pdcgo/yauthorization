package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.Static("/", "./assets")
	r.POST("/auth", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"asdasd": "asdasd"})
	})

	r.Run("localhost:8080")
}
