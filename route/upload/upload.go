package upload

import (
	"github.com/gin-gonic/gin"
	httpHandler "github.com/iamarpitzala/aca-reca-backend/internal/http"
	"github.com/iamarpitzala/aca-reca-backend/internal/middleware"
	"github.com/iamarpitzala/aca-reca-backend/internal/service"
)

func RegisterUploadRoutes(e *gin.RouterGroup, uploadHandler *httpHandler.UploadHandler, tokenService *service.TokenService) {
	upload := e.Group("/upload")
	upload.Use(middleware.AuthMiddleware(tokenService))

	upload.POST("/image", uploadHandler.UploadImage)
	upload.POST("/document", uploadHandler.UploadDocument)
}
