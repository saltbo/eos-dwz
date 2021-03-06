package main

import (
	"eos-dwz/src/rest"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	rest.SetupRouter(router)
	router.Run(":8128")
}
