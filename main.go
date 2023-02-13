package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"

	docs "github.com/aimotrens/html2pdf/docs"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var _wkhtmltopdfPath string

func main() {
	_wkhtmltopdfPath := os.Getenv("WKHTMLTOPDF_PATH")
	if _wkhtmltopdfPath == "" {
		_wkhtmltopdfPath = "./wkhtmltopdf"
	}

	if _, err := os.Stat(_wkhtmltopdfPath); err != nil {
		log.Fatal("wkhtmltopdf not found: '" + _wkhtmltopdfPath + "'")
	}

	docs.SwaggerInfo.Title = "HTML to PDF Konverter"
	docs.SwaggerInfo.BasePath = "/api"

	r := gin.Default()

	r.NoRoute(func(ctx *gin.Context) {
		ctx.Redirect(http.StatusTemporaryRedirect, "/swagger/index.html")
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	api := r.Group("/api")
	{
		api.GET("/healthcheck", HealthCheck)
		h2p := api.Group("/html2pdf")
		{
			h2p.POST("/convert", Convert)
		}
	}

	fmt.Println("Html2Pdf started.")
	r.Run(":8080")
}

// Convert HTML to PDF
// @Summary Konvertiert eine HTML Seite in ein A4 PDF-Dokument
// @Param request body string true "HTML-Seite"
// @Accept text/html
// @Produce  application/pdf
// @Success 200 {file} binary
// @Success 500 {string} string
// @Router /html2pdf/convert [post]
func Convert(g *gin.Context) {
	tmpID, err := uuid.NewV4()
	if err != nil {
		g.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	tmpFile := "/tmp/html2pdf-" + tmpID.String()
	defer os.Remove(tmpFile)

	data, err := io.ReadAll(g.Request.Body)
	if err != nil {
		g.AbortWithError(http.StatusBadRequest, err)
		return
	}

	p := exec.Command(_wkhtmltopdfPath, "--page-size", "A4", "-", tmpFile)
	w, err := p.StdinPipe()
	if err != nil {
		g.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	w.Write(data)
	w.Close()

	p.Stdout = os.Stdout
	p.Stderr = os.Stderr

	err = p.Start()
	if err != nil {
		g.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	p.Wait()

	if ec := p.ProcessState.ExitCode(); ec != 0 {
		g.AbortWithError(http.StatusInternalServerError, errors.New("wkhtmltopdf failed"))
		return
	}

	g.Header("content-type", "application/pdf")
	g.File(tmpFile)
}

// Healthcheck
// @Summary Gibt immer "OK" zurück
// @Produce  text/plain
// @Success 200 {string} OK
// @Router /healthcheck [get]
func HealthCheck(g *gin.Context) {
	g.String(http.StatusOK, "OK")
}
