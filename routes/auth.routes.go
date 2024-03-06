package routes

import (
	"andalalin/controllers"
	"andalalin/middleware"

	"github.com/gin-gonic/gin"
)

type AuthRouteController struct {
	authController controllers.AuthController
}

func NewAuthRouteController(authController controllers.AuthController) AuthRouteController {
	return AuthRouteController{authController}
}

func (rc *AuthRouteController) AuthRoute(rg *gin.RouterGroup) {
	router := rg.Group("/auth")

	//Pendaftaran akun
	router.POST("/register", rc.authController.SignUp)

	//Login aplikasi
	router.POST("/login", rc.authController.SignIn)

	//Refresh token
	router.POST("/refresh", rc.authController.RefreshAccessToken)

	//Verifikasi akun
	router.GET("/verification/:verificationCode", rc.authController.VerifyEmail)
	router.POST("/verification/resend", rc.authController.ResendVerification)

	//Logout
	router.POST("/logout", middleware.DeserializeUser(), rc.authController.LogoutUser)
}
