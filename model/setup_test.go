package model

import (
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

type Case[T any] struct {
	Title   string
	Input   []T
	Error   string
	Context []gin.Param
}

func MockContext() *gin.Context {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	return c
}
