package mock

import (
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MockedFilmModel struct {
	ID   int
	Name string
}

func (f MockedFilmModel) GetById(*gin.Context, *gorm.DB) error {
	return errors.New(MockedErrorMessage)
}
func (f MockedFilmModel) GetAll(*gin.Context, *gorm.DB) error {
	return errors.New(MockedErrorMessage)
}
func (f MockedFilmModel) CreateNew(*gin.Context, *gorm.DB) error {
	return errors.New(MockedErrorMessage)
}
func (f MockedFilmModel) UpdateById(*gin.Context, *gorm.DB) error {
	return errors.New(MockedErrorMessage)
}
func (f MockedFilmModel) DeleteById(*gin.Context, *gorm.DB) error {
	return errors.New(MockedErrorMessage)
}
