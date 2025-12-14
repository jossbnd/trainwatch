package main

import (
	"github.com/jossbnd/trainwatch/backend/internal/api"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	api.RegisterRoutes(r)
	r.Run() // listen and serve on
}
