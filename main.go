package main

import (
	"asscll_art/internal/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
    // Создаем роутер
    router := gin.Default()
    router.MaxMultipartMemory = 8 << 2
    router.POST("/upload", handlers.UploadHandler)

    router.Run(":8080")
}