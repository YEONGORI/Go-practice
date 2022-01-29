package main

import (
	"go_src/test/controllers"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/", controllers.MainPage)

	r.Run()
}
