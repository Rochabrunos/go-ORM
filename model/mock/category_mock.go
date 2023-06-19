package mock

import (
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MockedCategoryModel struct {
	ID   int
	Name string
}

var MockedErrorMessage = "mocked error"

func (c MockedCategoryModel) GetById(*gin.Context, *gorm.DB) error {
	return errors.New(MockedErrorMessage)
}
func (c MockedCategoryModel) GetAll(*gin.Context, *gorm.DB) error {
	return errors.New(MockedErrorMessage)
}
func (c MockedCategoryModel) CreateNew(*gin.Context, *gorm.DB) error {
	return errors.New(MockedErrorMessage)
}
func (c MockedCategoryModel) UpdateById(*gin.Context, *gorm.DB) error {
	return errors.New(MockedErrorMessage)
}
func (c MockedCategoryModel) DeleteById(*gin.Context, *gorm.DB) error {
	return errors.New(MockedErrorMessage)
}
