package controllers

import "github.com/gin-gonic/gin"

type UserController struct {
	PostUserSignIn func(*gin.Context) `path:"/user/sign-in"`
	PostUserLogin func(*gin.Context) `path:"/user/login"`
	PostUserLogout func(*gin.Context) `path:"/user/logout"`
	PostCurrentUser func(*gin.Context) `path:"/user/current-user"`
	PostChangePassword func(*gin.Context) `path:"/user/change-password"`
}

func (UserController) Name() string {
	return "UserController"
}

func NewUserController() UserController{
	return UserController{
		PostUserSignIn: PostUserSignIn,
		PostUserLogin:PostUserLogin,
		PostUserLogout:PostUserLogout,
		PostCurrentUser:PostCurrentUser,
		PostChangePassword:PostChangePassword,
	}
}

type PostUserLoginRequest struct {

}

func PostUserSignIn(c *gin.Context) {

}

func PostUserLogin(c *gin.Context) {

}

func PostUserLogout(c *gin.Context) {

}

func PostCurrentUser(c *gin.Context) {

}

func PostChangePassword(c *gin.Context) {

}