package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	config.SetupDatabase()
	routes.SetupRoutes(r)
	r.Run()
}
