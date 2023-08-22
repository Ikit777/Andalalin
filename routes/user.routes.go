package routes

import (
	"github.com/Ikit777/E-Andalalin/controllers"
	"github.com/Ikit777/E-Andalalin/middleware"
	"github.com/gin-gonic/gin"
)

type UserRouteController struct {
	userController controllers.UserController
}

func NewRouteUserController(userController controllers.UserController) UserRouteController {
	return UserRouteController{userController}
}

func (uc *UserRouteController) UserRoute(rg *gin.RouterGroup) {

	router := rg.Group("users")

	router.GET("/me", middleware.DeserializeUser(), uc.userController.GetMe)

	router.GET("/", middleware.DeserializeUser(), uc.userController.GetUsers)
	router.GET("/email/:emailUser", middleware.DeserializeUser(), uc.userController.GetUserByEmail)
	router.GET("role/:role", middleware.DeserializeUser(), uc.userController.GetUsersSortRole)

	router.POST("/add", middleware.DeserializeUser(), uc.userController.Add)
	router.DELETE("/delete", middleware.DeserializeUser(), uc.userController.Delete)

	router.POST("/forgotpassword", uc.userController.ForgotPassword)
	router.PATCH("/resetpassword/:resetToken", uc.userController.ResetPassword)

	router.POST("/updatephoto", middleware.DeserializeUser(), uc.userController.UpdatePhoto)

	router.GET("/notifikasi", middleware.DeserializeUser(), uc.userController.GetNotifikasi)
	router.DELETE("/clearnotifikasi", middleware.DeserializeUser(), uc.userController.ClearNotifikasi)
}
