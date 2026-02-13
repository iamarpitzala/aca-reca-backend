package entry

import (
	httpHandler "github.com/iamarpitzala/aca-reca-backend/internal/http"
	"github.com/iamarpitzala/aca-reca-backend/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/iamarpitzala/aca-reca-backend/internal/service"
)

func RegisterEntryRoute(e *gin.RouterGroup, entry *httpHandler.EntryHandler, tokenService *service.TokenService) {
	entr := e.Group("/entry")
	entr.Use(middleware.AuthMiddleware(tokenService))
	entr.POST("/", entry.AddEntry)
}
