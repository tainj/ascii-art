package handlers

import (
	"asscll_art/internal/services"
	"fmt"
	"log"
	"net/http"
	"path"

	"github.com/gin-gonic/gin"
)

func UploadHandler(c *gin.Context) {
    file, err := c.FormFile("file")

	if err != nil {
		c.String(http.StatusRequestEntityTooLarge, "file is too big")
	}

    log.Println(file.Filename)

	dst := path.Join("./temp", file.Filename)

	c.SaveUploadedFile(file, dst)

	_, err = services.ConvertToASCII(dst)
	if err != nil {
		c.String(http.StatusRequestEntityTooLarge, "file is too big")
	}

	_, err = services.ScaleGrayImageJpg(dst, 120) 
	if err != nil {
		c.String(http.StatusRequestEntityTooLarge, "file is too big")
	}

    c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
}