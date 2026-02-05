package aoc

import (
	"github.com/gin-gonic/gin"
	"github.com/iamarpitzala/aca-reca-backend/config"
	httpHandler "github.com/iamarpitzala/aca-reca-backend/internal/http"
	"github.com/iamarpitzala/aca-reca-backend/internal/middleware"
	"github.com/iamarpitzala/aca-reca-backend/internal/service"
)

func RegisterAOCRoutes(e *gin.RouterGroup, aocHandler *httpHandler.AOCHandler) {
	aoc := e.Group("/aoc")
	cfg := config.Load()

	tokenService := service.NewTokenService(cfg.JWT)
	aoc.Use(middleware.AuthMiddleware(tokenService))

	aoc.GET("/type", aocHandler.GetAllAOCType)
	aoc.GET("/tax", aocHandler.GetAllAccountTax)
	aoc.POST("/", aocHandler.CreateAOC)
	aoc.GET("/:id", aocHandler.GetAOCByID)
	aoc.GET("/code/:code", aocHandler.GetAOCByCode)
	aoc.GET("/account-type/:accountTypeId", aocHandler.GetAOCByAccountTypeID)
	aoc.GET("/account-tax/:accountTaxId", aocHandler.GetAOCByAccountTaxID)
	aoc.PUT("/:id", aocHandler.UpdateAOC)
	aoc.DELETE("/:id", aocHandler.DeleteAOC)
}
