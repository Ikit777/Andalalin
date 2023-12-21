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

	//Get detail user
	router.GET("/me", middleware.DeserializeUser(), uc.userController.GetMe)

	//Get semua user
	router.GET("/all", middleware.DeserializeUser(), uc.userController.GetUsers)

	//Get user berdasarkan email
	router.GET("/email/:emailUser", middleware.DeserializeUser(), uc.userController.GetUserByEmail)

	//GEt user berdasarkan role
	router.GET("role/:role", middleware.DeserializeUser(), uc.userController.GetUsersSortRole)

	//Menindaklanjuti user
	router.POST("/add", middleware.DeserializeUser(), uc.userController.Add)
	router.DELETE("/delete", middleware.DeserializeUser(), uc.userController.Delete)

	//Merubah akun
	router.POST("/edit/account", middleware.DeserializeUser(), uc.userController.EditAkun)

	//Merubah foto profil
	router.POST("/edit/photo", middleware.DeserializeUser(), uc.userController.UpdatePhoto)

	//Reset password
	router.POST("/password/forgot", uc.userController.ForgotPassword)
	router.PATCH("/password/reset/:resetToken", uc.userController.ResetPassword)

	//Get notifikasi
	router.GET("/notification", middleware.DeserializeUser(), uc.userController.GetNotifikasi)

	//Bersihkan notifikasi
	router.DELETE("/notification/delete", middleware.DeserializeUser(), uc.userController.ClearNotifikasi)

	//Get semua petugas untuk pilih petugas
	router.GET("/petugas", middleware.DeserializeUser(), uc.userController.GetPetugas)
}
