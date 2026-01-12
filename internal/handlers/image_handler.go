package handlers

import (
	"fmt"
	"log"
	"net/http"
	"path"

	"github.com/gin-gonic/gin"
)

func UploadHandler(c *gin.Context) {
    file, err := c.FormFile("file")

	if err != nil {
		c.String(http.StatusRequestEntityTooLarge, fmt.Sprintf("file is too big"))
	}

    log.Println(file.Filename)

	dst := path.Join("./temp", file.Filename)

	c.SaveUploadedFile(file, dst)

    c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
}