package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/aimotrens/html2pdf/app/healthcheck"
	"github.com/aimotrens/html2pdf/app/html2pdf"
	_ "github.com/aimotrens/html2pdf/docs"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title HTML zu PDF Konverter API
// @description API für den HTML zu PDF Konverter
// @version 1.0.0
// @BasePath /api

// @contact.name   aimotrens
// @contact.url    https://github.com/aimotrens/html2pdf

// @license.name  MIT Lizenz
// @license.url   https://github.com/aimotrens/html2pdf/blob/main/LICENSE

func main() {
	wkhtmltopdfPath := os.Getenv("WKHTMLTOPDF_PATH")
	if wkhtmltopdfPath == "" {
		wkhtmltopdfPath = "./wkhtmltopdf"
	}

	if _, err := os.Stat(wkhtmltopdfPath); err != nil {
		log.Fatal("wkhtmltopdf not found: '" + wkhtmltopdfPath + "'")
	}

	r := gin.Default()

	r.NoRoute(func(ctx *gin.Context) {
		ctx.Redirect(http.StatusTemporaryRedirect, "/swagger/index.html")
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	api := r.Group("/api")
	{
		hc := api.Group("/healthcheck")
		{
			hc.GET("/", healthcheck.HealthCheck)
		}

		h2p := api.Group("/html2pdf")
		{
			h2p.POST("/convert", func(ctx *gin.Context) { html2pdf.Convert(ctx, wkhtmltopdfPath) })
		}
	}

	fmt.Println("Html2Pdf started.")
	r.Run(":8080")
}
