package main

import (
	controllers "orm-golang/controller"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	lang := r.Group("/languages")
	{
		lang.GET("/:id", controllers.GetLangByIdEndpoint)
		lang.GET("", controllers.GetAllLangsEndpoint)
		lang.POST("", controllers.CreateLangEndpoint)
		lang.PUT("/:id", controllers.ModifyLangEndpoint)
		lang.DELETE("/:id", controllers.DeleteLangEndpoint)
	}

	category := r.Group("/categories")
	{
		category.GET("/:id", controllers.GetByIdEndpoint)
		category.GET("", controllers.GetAllEndpoint)
		category.POST("", controllers.CreateEndpoint)
		category.PUT("/:id", controllers.ModifyEndpoint)
		category.DELETE("/:id", controllers.DeleteEndpoint)
	}

	film := r.Group("/films")
	{
		film.GET("/:id", controllers.GetFilmByIdEndpoint)
		film.GET("", controllers.GetAllFilmsEndpoint)
		film.POST("", controllers.CreateFilmEndpoint)
		film.PUT("/:id", controllers.ModifyFilmEndpoint)
		film.DELETE("/:id", controllers.DeleteFilmEndpoint)
	}
	r.Run(":8080") // listen and serve on 0.0.0.0:8080
}
