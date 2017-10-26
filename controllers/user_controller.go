package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/penguinn/album/components/user"
	"github.com/penguinn/album/models"
	"github.com/penguinn/album/utils"
	"github.com/penguinn/penguin/component/log"
	"net/http"
)

type UserController struct {
	PostUserSignIn     func(*gin.Context) `path:"/user/sign-in"`
	PostUserLogin      func(*gin.Context) `path:"/user/login"`
	PostUserLogout     func(*gin.Context) `path:"/user/logout"`
	PostCurrentUser    func(*gin.Context) `path:"/user/current-user"`
	PostChangePassword func(*gin.Context) `path:"/user/change-password"`
}

func (UserController) Name() string {
	return "UserController"
}

func NewUserController() UserController {
	return UserController{
		PostUserSignIn:     PostUserSignIn,
		PostUserLogin:      PostUserLogin,
		PostUserLogout:     PostUserLogout,
		PostCurrentUser:    PostCurrentUser,
		PostChangePassword: PostChangePassword,
	}
}

type PostUserSignInRequest struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
	AuthCode string `form:"authCode" binding:"required"`
}

type PostUserSignInResponse struct {
	UserID int `json:"userID"`
}

type PostUserLoginRequest struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

type PostChangePasswordRequest struct {
	Password    string `form:"password" binding:"required"`
	NewPassword string `form:"newPassword" binding:"required"`
}

func PostUserSignIn(c *gin.Context) {
	var postUserSignInRequest PostUserSignInRequest
	err := c.Bind(&postUserSignInRequest)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusOK, NewResponse(false, PARAM_ERROR_CODE))
		return
	}
	//validate username
	isExist, err := models.User{}.ValidateUsername(postUserSignInRequest.Username)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusOK, NewResponse(false, SERVER_ERROR_CODE))
		return
	}
	if isExist {
		c.JSON(http.StatusOK, NewResponse(false, USERNAME_EXIST_CODE))
		return
	}
	//validate authCode
	cookie, err := c.Request.Cookie("sess-token")
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusOK, NewResponse(false, PARAM_ERROR_CODE, "获取不到token"))
		return
	}
	token := cookie.Value
	ok, err := models.ImageAuth{}.ValidateAuthCode(token, postUserSignInRequest.AuthCode)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusOK, NewResponse(false, SERVER_ERROR_CODE))
		return
	}
	if !ok {
		errStr := "图形验证码错误"
		log.Error(errStr)
		c.JSON(http.StatusOK, NewResponse(false, IMAGE_AUTH_CODE, errStr))
		return
	}

	userID, err := models.User{}.Insert(postUserSignInRequest.Username, utils.HashPassword(postUserSignInRequest.Password))
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusOK, NewResponse(false, SERVER_ERROR_CODE))
		return
	}
	postUserSignInResponse := PostUserSignInResponse{
		UserID: userID,
	}
	c.JSON(http.StatusOK, NewResponse(postUserSignInResponse, SUCCESS_CODE))
}

func PostUserLogin(c *gin.Context) {
	var postUserLoginRequest PostUserLoginRequest
	err := c.Bind(&postUserLoginRequest)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusOK, NewResponse(false, PARAM_ERROR_CODE))
		return
	}
	userInst, err := models.User{}.SelectByUsername(postUserLoginRequest.Username)
	if err == gorm.ErrRecordNotFound {
		errStr := "用户名不存在"
		log.Error(errStr + " " + postUserLoginRequest.Username)
		c.JSON(http.StatusOK, NewResponse(false, AUTH_FAILED_CODE, errStr))
		return
	}
	if err != nil {
		errStr := "系统异常"
		c.JSON(http.StatusOK, NewResponse(false, AUTH_FAILED_CODE, errStr))
		return
	}

	if !utils.ComparePassword(userInst.Password, postUserLoginRequest.Password) {
		errStr := "密码错误"
		log.Error(errStr)
		c.JSON(http.StatusOK, NewResponse(false, AUTH_FAILED_CODE, errStr))
		return
	}

	userManager := user.FromContext(c)
	err = userManager.Login(userInst.ID, userInst.Username)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusOK, NewResponse(false, AUTH_FAILED_CODE, err.Error()))
		return
	}
	c.JSON(http.StatusOK, NewResponse(true))
}

func PostUserLogout(c *gin.Context) {
	userManager := user.FromContext(c)
	userManager.Logout()
	c.JSON(http.StatusOK, NewResponse(true))
}

func PostCurrentUser(c *gin.Context) {

}

func PostChangePassword(c *gin.Context) {
	var postChangePasswordRequest PostChangePasswordRequest
	err := c.Bind(&postChangePasswordRequest)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusOK, NewResponse(false, PARAM_ERROR_CODE))
		return
	}
	userManager := user.FromContext(c)
	userID := userManager.GetID()

	userInst, err := models.User{}.SelectByUserID(userID)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusOK, NewResponse(false, SERVER_ERROR_CODE))
		return
	}
	if !utils.ComparePassword(userInst.Password, postChangePasswordRequest.Password) {
		errStr := "原密码不正确"
		log.Error(errStr)
		c.JSON(http.StatusOK, NewResponse(false, SERVER_ERROR_CODE, errStr))
		return
	}

	updateMap := map[string]interface{}{
		"password": utils.HashPassword(postChangePasswordRequest.NewPassword),
	}

	err = models.User{}.Update(userID, updateMap)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusOK, NewResponse(false, SERVER_ERROR_CODE, err.Error()))
		return
	}
	c.JSON(http.StatusOK, NewResponse(true))
}
