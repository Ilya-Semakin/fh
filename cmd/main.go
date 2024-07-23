package main

import (
	"github.com/Ilya-Semakin/fh/config"
	"github.com/Ilya-Semakin/fh/routres"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	config.SetupDatabase()
	routes.SetupRoutes(r)
	r.Run()
}
