package main

import (
	"github.com/gin-gonic/gin"
	controllers"orm-golang/controller"
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
		category.GET("/:id", controllers.GetCategByIdEndpoint)
		category.GET("", controllers.GetAllCategsEndpoint)
		category.POST("", controllers.CreateCategEndpoint)
		category.PUT("/:id", controllers.ModifyCategEndpoint)
		category.DELETE("/:id", controllers.DeleteCategEndpoint)
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

