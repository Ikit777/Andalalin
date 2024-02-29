package routes

import (
	"andalalin/controllers"
	"andalalin/middleware"

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
	router.GET("/kepuasan/hasil/periode/:waktu", middleware.DeserializeUser(), sc.surveyController.HasilSurveiKepuasanTertentu)

	//Pengaduan
	router.POST("/pengaduan", middleware.DeserializeUser(), sc.surveyController.IsiPengaduan)
	router.GET("/pengaduan/detail/:id_pengaduan", middleware.DeserializeUser(), sc.surveyController.GetPengaduan)
	router.GET("/pengaduan/daftar", middleware.DeserializeUser(), sc.surveyController.GetAllPengaduan)
	router.GET("/pengaduan/daftar/petugas", middleware.DeserializeUser(), sc.surveyController.GetAllPengaduanByPetugas)
	router.POST("/pengaduan/terima/:id_pengaduan", middleware.DeserializeUser(), sc.surveyController.TerimPengaduan)
}
