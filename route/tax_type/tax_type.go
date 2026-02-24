package tax_type

import (
	"github.com/gin-gonic/gin"
	httpHandler "github.com/iamarpitzala/aca-reca-backend/internal/http"
	"github.com/iamarpitzala/aca-reca-backend/internal/middleware"
	"github.com/iamarpitzala/aca-reca-backend/internal/service"
)

func RegisterTaxTypeRoutes(e *gin.RouterGroup, handler *httpHandler.TaxTypeHandler, tokenService *service.TokenService) {
	g := e.Group("/tax-types")
	g.Use(middleware.AuthMiddleware(tokenService))
	g.GET("", handler.GetTaxTypes)
}
