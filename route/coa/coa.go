package coa

import (
	"github.com/gin-gonic/gin"
	httpHandler "github.com/iamarpitzala/aca-reca-backend/internal/http"
	"github.com/iamarpitzala/aca-reca-backend/internal/middleware"
	"github.com/iamarpitzala/aca-reca-backend/internal/service"
)

func RegisterCOARoutes(e *gin.RouterGroup, coaHandler *httpHandler.COAHandler, tokenService *service.TokenService) {
	coa := e.Group("/coa")
	coa.Use(middleware.AuthMiddleware(tokenService))

	coa.GET("/type", coaHandler.GetAllCOAType)
	coa.GET("/tax", coaHandler.GetAllAccountTax)
	coa.GET("/account-types", coaHandler.GetCOAsByAccountType)
	coa.GET("/account-type/:id", coaHandler.GetCOAByAccountTypeID)
	coa.GET("/code/:code", coaHandler.GetCOAByCode)
	coa.POST("", coaHandler.CreateCOA)
	coa.POST("/", coaHandler.CreateCOA)
	coa.GET("", coaHandler.GetAllCOAs)
	coa.GET("/", coaHandler.GetAllCOAs)
	coa.GET("/:id", coaHandler.GetCOAByID)
	coa.PUT("/:id", coaHandler.UpdateCOA)
	coa.PATCH("", coaHandler.DeleteCOA)
	coa.PATCH("/bulk-tax", coaHandler.BulkUpdateCOATax)
	coa.PATCH("/archive", coaHandler.ArchiveCOA)
}
