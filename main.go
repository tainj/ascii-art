package main

import (
	"asscll_art/internal/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
    // Создаем роутер
    router := gin.Default()
    router.MaxMultipartMemory = 8 << 2

    // Регистрируем хэндлер /upload
    router.POST("/upload", handlers.UploadHandler)

    // Запускаем роутер
    router.Run(":8080")
}