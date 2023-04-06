package html2pdf

import (
	"errors"
	"io"
	"net/http"
	"os"
	"os/exec"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

// @Tags Html2Pdf
// @Summary Konvertiert eine HTML Seite in ein A4 PDF-Dokument
// @Param request body string true "HTML-Seite"
// @Accept text/html
// @Produce  application/pdf
// @Success 200 {file} binary
// @Success 500 {string} string
// @Router /html2pdf/convert [post]
func Convert(g *gin.Context, wkhtmltopdfPath string) {
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

	p := exec.Command(wkhtmltopdfPath, "--page-size", "A4", "-", tmpFile)
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
