package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type BaseController struct {
	GetBase   func(*gin.Context) `path:""`
}

func (BaseController) Name() string {
	return "BaseController"
}

func NewBaseController() BaseController {
	return BaseController{
		GetBase:GetBase,
	}
}

func GetBase(c *gin.Context) {
	c.JSON(http.StatusOK, NewResponse(true))
}
