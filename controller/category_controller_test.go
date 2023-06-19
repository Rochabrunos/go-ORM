package controller

import (
	"fmt"
	"net/http/httptest"
	mock "orm-golang/model/mock"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGetCategByIdEndpoint(t *testing.T) {
	t.Run("Test fail to get a category", func(t *testing.T) {
		handler = Handler{}.New(&mock.MockedCategory{})
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.AddParam("id", "1")
		GetCategByIdEndpoint(c)
		fmt.Println("result", w)

	})
	t.Run("Test return a category as expect", func(t *testing.T) {

	})
}
