package api

import (
	"github.com/febster16/go-auth/common"
	"github.com/febster16/go-auth/internal/api/handlers"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(r *gin.Engine) {
	r.POST("/signup", handlers.Signup)
	r.POST("/login", handlers.Login)
	r.GET("/validate", common.RequireAuthMiddleware, handlers.Validate)
	r.PATCH("/change-password", handlers.ChangePassword)
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
