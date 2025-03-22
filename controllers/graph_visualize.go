package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func VisualizeStreamGraph(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
	})
}

func VisualizeQueryGraph(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
	})
}