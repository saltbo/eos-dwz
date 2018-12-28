package main

import (
	"eos-dwz/src/rest"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.New()
	rest.SetupRouter(router)
	router.Run(":8128")
}
