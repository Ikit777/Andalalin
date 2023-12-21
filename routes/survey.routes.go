package routes

import (
	"github.com/Ikit777/E-Andalalin/controllers"
	"github.com/Ikit777/E-Andalalin/middleware"
	"github.com/gin-gonic/gin"
)

type SurveyRouteController struct {
	surveyController controllers.SurveyController
}

func NewSurveyRouteController(surveyController controllers.SurveyController) SurveyRouteController {
	return SurveyRouteController{surveyController}
}

func (sc *SurveyRouteController) SurveyRoute(rg *gin.RouterGroup) {
	router := rg.Group("/survey")

	//Survei kepuasan
	router.POST("/kepuasan/:id_andalalin", middleware.DeserializeUser(), sc.surveyController.SurveiKepuasan)
	router.GET("/kepuasan/cek/:id_andalalin", middleware.DeserializeUser(), sc.surveyController.CekSurveiKepuasan)
	router.GET("/kepuasan/hasil", middleware.DeserializeUser(), sc.surveyController.HasilSurveiKepuasan)

	//Survei mandiri
	router.POST("/mandiri", middleware.DeserializeUser(), sc.surveyController.IsiSurveyMandiri)
	router.GET("/mandiri/detail/:id_survei", middleware.DeserializeUser(), sc.surveyController.GetSurveiMandiri)
	router.GET("/mandiri/daftar", middleware.DeserializeUser(), sc.surveyController.GetAllSurveiMandiri)
	router.GET("/mandiri/daftar/petugas", middleware.DeserializeUser(), sc.surveyController.GetAllSurveiMandiriByPetugas)
	router.POST("/mendiri/terima/:id_survei", middleware.DeserializeUser(), sc.surveyController.TerimaSurvei)
}
