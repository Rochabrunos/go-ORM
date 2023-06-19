package mock

import (
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MockedCategory struct {
	ID   int
	Name string
}

func (c MockedCategory) GetCategoryById(*gin.Context, *gorm.DB) error {
	return errors.New("Mocked error")
}
func (c MockedCategory) GetAllCategories(*gin.Context, *gorm.DB) error {
	return errors.New("Mocked error")
}
func (c MockedCategory) CreateNewCategory(*gin.Context, *gorm.DB) error {
	return errors.New("Mocked error")
}
func (c MockedCategory) UpdateCategoryById(*gin.Context, *gorm.DB) error {
	return errors.New("Mocked error")
}
func (c MockedCategory) DeleteCategoryById(*gin.Context, *gorm.DB) error {
	return errors.New("Mocked error")
}
