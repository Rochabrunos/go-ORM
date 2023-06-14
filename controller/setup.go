package controller

import (
	"github.com/gin-gonic/gin"
)

func SendError(c *gin.Context, statuCode int, err error) {
	c.JSON(statuCode, gin.H{
		"error": err.Error(),
	})
}