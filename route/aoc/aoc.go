package aoc

import (
	"github.com/gin-gonic/gin"
	httpHandler "github.com/iamarpitzala/aca-reca-backend/internal/http"
	"github.com/iamarpitzala/aca-reca-backend/internal/middleware"
	"github.com/iamarpitzala/aca-reca-backend/internal/service"
)

func RegisterAOCRoutes(e *gin.RouterGroup, aocHandler *httpHandler.AOCHandler, tokenService *service.TokenService) {
	aoc := e.Group("/aoc")
	aoc.Use(middleware.AuthMiddleware(tokenService))

	aoc.GET("/type", aocHandler.GetAllAOCType)
	aoc.GET("/tax", aocHandler.GetAllAccountTax)
	aoc.GET("/account-types", aocHandler.GetAOCsByAccountType)
	aoc.GET("/account-type/:id", aocHandler.GetAOCByAccountTypeID)
	aoc.GET("/account-tax/:id", aocHandler.GetAOCByAccountTaxID)
	aoc.GET("/code/:code", aocHandler.GetAOCByCode)
	aoc.POST("", aocHandler.CreateAOC)
	aoc.POST("/", aocHandler.CreateAOC)
	aoc.GET("", aocHandler.GetAllAOCs)
	aoc.GET("/", aocHandler.GetAllAOCs)
	aoc.GET("/:id", aocHandler.GetAOCByID)
	aoc.PUT("/:id", aocHandler.UpdateAOC)
	aoc.PATCH("", aocHandler.DeleteAOC)
	aoc.PATCH("/bulk-tax", aocHandler.BulkUpdateTax)
	aoc.PATCH("/archive", aocHandler.ArchiveAOC)
}
