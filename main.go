package main

import (
	controllers "orm-golang/controller"

	"github.com/gin-gonic/gin"
	kgin "github.com/keploy/go-sdk/integrations/kgin/v1"
	"github.com/keploy/go-sdk/keploy"
)

func main() {
	r := gin.Default()
	port := "8080"
	k := keploy.New(keploy.Config{
		App: keploy.AppConfig{
			Name: "dvdrental",
			Port: port,
		},
		Server: keploy.ServerConfig{
			URL: "http://localhost:6789/api",
		},
	})
	kgin.GinV1(k, r)

	lang := r.Group("/languages")
	{
		lang.GET("/:id", controllers.GetByIdLanguageEndpoint)
		lang.GET("", controllers.GetAllLanguageEndpoint)
		lang.POST("", controllers.CreateLanguageEndpoint)
		lang.PUT("/:id", controllers.ModifyLanguageEndpoint)
		lang.DELETE("/:id", controllers.DeleteLanguageEndpoint)
	}

	category := r.Group("/categories")
	{
		category.GET("/:id", controllers.GetByIdCategoryEndpoint)
		category.GET("", controllers.GetAllCategoryEndpoint)
		category.POST("", controllers.CreateCategoryEndpoint)
		category.PUT("/:id", controllers.ModifyCategoryEndpoint)
		category.DELETE("/:id", controllers.DeleteCategoryEndpoint)
	}

	film := r.Group("/films")
	{
		film.GET("/:id", controllers.GetByIdFilmEndpoint)
		film.GET("", controllers.GetAllFilmsEndpoint)
		film.POST("", controllers.CreateFilmEndpoint)
		film.PUT("/:id", controllers.ModifyFilmEndpoint)
		film.DELETE("/:id", controllers.DeleteFilmEndpoint)
	}
	r.Run(":" + port) // listen and serve on 0.0.0.0:8080
}
