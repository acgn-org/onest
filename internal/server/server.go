package server

import (
	"github.com/gin-gonic/gin"
)

func New() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	Engine := gin.Default()

	Api(Engine.Group("api"))

	return Engine
}
