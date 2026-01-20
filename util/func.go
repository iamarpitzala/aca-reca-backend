package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func BindAndValidate(c *gin.Context, rq interface{}) error {
	if err := c.ShouldBind(rq); err != nil {
		return err
	}

	if err := validator.New().Struct(rq); err != nil {
		return err
	}

	return nil
}
