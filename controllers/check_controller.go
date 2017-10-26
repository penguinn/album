package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/penguinn/penguin/component/log"
	"net/http"
)

type CheckController struct {
	PostImageAuth func(*gin.Context) `path:"/check/image-auth"`
	PostUserName func(*gin.Context) `path:"/check/username"`
}

func(CheckController) Name() string {
	return "CheckController"
}

func NewCheckController() CheckController {
	return CheckController{
		PostImageAuth:PostImageAuth,
		PostUserName:PostUserName,
	}
}

type ImageAuthRequest struct {
	Type    int  `form:"type"` 		//1：代表主界面
}

func PostImageAuth(c *gin.Context) {
	var imageAuthRequest ImageAuthRequest
	err := c.Bind(&imageAuthRequest)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusOK, NewResponse(false, PARAM_ERROR_CODE))
		return
	}


}

func PostUserName(c *gin.Context) {

}
