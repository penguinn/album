package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/penguinn/album/components/aliyun-oss"
	"github.com/penguinn/penguin/component/log"
	"net/http"
)

type UploadController struct {
	PostGetSignURL  func(*gin.Context) `path:"/upload/sign-url"`
	PostUploadImage func(*gin.Context) `path:"/upload/image"`
}

func (UploadController) Name() string {
	return "UploadController"
}

func NewUploadController() UploadController {
	return UploadController{
		PostGetSignURL:  PostGetSignURL,
		PostUploadImage: PostUploadImage,
	}
}

type PostGetSignURLRequest struct {
	ObjectName string `form:"objectName"`
}

type PostGetSignURLResponse struct {
	SignURL string `json:"signURL"`
}

func PostGetSignURL(c *gin.Context) {
	var postGetSignURLRequest PostGetSignURLRequest
	err := c.Bind(&postGetSignURLRequest)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusOK, NewResponse(false, PARAM_ERROR_CODE))
		return
	}
	signURL, err := aliyun_oss.GetSignURL(postGetSignURLRequest.ObjectName)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusOK, NewResponse(false, SERVER_ERROR_CODE))
		return
	}
	postGetSignURLResponse := PostGetSignURLResponse{
		SignURL: signURL,
	}
	c.JSON(http.StatusOK, NewResponse(postGetSignURLResponse))
}

func PostUploadImage(c *gin.Context) {
	multipartForm, err := c.MultipartForm()
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusOK, NewResponse(false, PARAM_ERROR_CODE))
		return
	}
	objectName := multipartForm.Value["objectName"][0]

	file := multipartForm.File["uploadFile"][0]
	err = c.SaveUploadedFile(file, objectName)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusOK, NewResponse(false, PARAM_ERROR_CODE))
		return
	}

	err = aliyun_oss.FileUpload(objectName, objectName)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusOK, NewResponse(false, SERVER_ERROR_CODE))
		return
	}
	c.JSON(http.StatusOK, NewResponse(true))
}
