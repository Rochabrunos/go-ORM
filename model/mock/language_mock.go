package mock

import (
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MockedLanguageModel struct {
	ID   int
	Name string
}

func (l MockedLanguageModel) GetById(*gin.Context, *gorm.DB) error {
	return errors.New(MockedErrorMessage)
}
func (l MockedLanguageModel) GetAll(*gin.Context, *gorm.DB) error {
	return errors.New(MockedErrorMessage)
}
func (l MockedLanguageModel) CreateNew(*gin.Context, *gorm.DB) error {
	return errors.New(MockedErrorMessage)
}
func (l MockedLanguageModel) UpdateById(*gin.Context, *gorm.DB) error {
	return errors.New(MockedErrorMessage)
}
func (l MockedLanguageModel) DeleteById(*gin.Context, *gorm.DB) error {
	return errors.New(MockedErrorMessage)
}
