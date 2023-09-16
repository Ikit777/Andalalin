package routes

import (
	"github.com/Ikit777/E-Andalalin/controllers"
	"github.com/Ikit777/E-Andalalin/middleware"
	"github.com/gin-gonic/gin"
)

type AndalalinRouteController struct {
	andalalinController controllers.AndalalinController
}

func NewRouteAndalalinController(andalalinController controllers.AndalalinController) AndalalinRouteController {
	return AndalalinRouteController{andalalinController}
}

func (uc *AndalalinRouteController) AndalalainRoute(rg *gin.RouterGroup) {

	router := rg.Group("andalalin")

	router.POST("/pengajuan", middleware.DeserializeUser(), uc.andalalinController.Pengajuan)
	router.POST("/pengajuanperlalin", middleware.DeserializeUser(), uc.andalalinController.PengajuanPerlalin)

	router.GET("/persyaratan/:id_andalalin", middleware.DeserializeUser(), uc.andalalinController.GetPersyaratan)
	router.GET("/perusahaan/:id_andalalin", middleware.DeserializeUser(), uc.andalalinController.GetPerusahaan)
	router.GET("/permohonan", middleware.DeserializeUser(), uc.andalalinController.GetPermohonan)
	router.GET("/userpermohonan", middleware.DeserializeUser(), uc.andalalinController.GetPermohonanByIdUser)
	router.GET("/detailpermohonan/:id_andalalin", middleware.DeserializeUser(), uc.andalalinController.GetPermohonanByIdAndalalin)
	router.GET("/detailperlalin/:id_andalalin", middleware.DeserializeUser(), uc.andalalinController.GetPermohonanByIdPerlalin)
	router.GET("/statuspermohonan/:status_andalalin", middleware.DeserializeUser(), uc.andalalinController.GetPermohonanByStatus)

	router.GET("/tiketpermohonan/:status", middleware.DeserializeUser(), uc.andalalinController.GetAndalalinTicketLevel1)
	router.GET("/petugaspermohonan/:status", middleware.DeserializeUser(), uc.andalalinController.GetAndalalinTicketLevel2)

	router.POST("/updatestatus/:id_andalalin/:status", middleware.DeserializeUser(), uc.andalalinController.UpdateStatusPermohonan)

	router.POST("/persyaratantidaksesuai/:id_andalalin", middleware.DeserializeUser(), uc.andalalinController.PersyaratanTidakSesuai)
	router.POST("/persyaratanterpenuhi/:id_andalalin", middleware.DeserializeUser(), uc.andalalinController.PersyaratanTerpenuhi)

	router.POST("/updatepersyaratan/:id_andalalin", middleware.DeserializeUser(), uc.andalalinController.UpdatePersyaratan)

	router.POST("/pilihpetugas/:id_andalalin", middleware.DeserializeUser(), uc.andalalinController.TambahPetugas)
	router.POST("/gantipetugas/:id_andalalin", middleware.DeserializeUser(), uc.andalalinController.GantiPetugas)

	router.POST("/survey/:id_andalalin", middleware.DeserializeUser(), uc.andalalinController.IsiSurvey)
	router.GET("/getsurvey/:id_andalalin", middleware.DeserializeUser(), uc.andalalinController.GetSurvey)
	router.GET("/getallsurvey", middleware.DeserializeUser(), uc.andalalinController.GetAllSurvey)

	router.GET("/getbap/:id_andalalin", middleware.DeserializeUser(), uc.andalalinController.GetBAP)
	router.POST("/bap/:id_andalalin", middleware.DeserializeUser(), uc.andalalinController.LaporanBAP)

	router.POST("/persetujuan/:id_andalalin", middleware.DeserializeUser(), uc.andalalinController.PersetujuanDokumen)
	router.GET("/getpersetujuan/:id_andalalin", middleware.DeserializeUser(), uc.andalalinController.GetPersetujuanDokumen)

	router.POST("/sk/:id_andalalin", middleware.DeserializeUser(), uc.andalalinController.LaporanSK)
	router.GET("/getsk/:id_andalalin", middleware.DeserializeUser(), uc.andalalinController.GetSK)

	router.POST("/usulantindakan/:id_andalalin", middleware.DeserializeUser(), uc.andalalinController.UsulanTindakanPengelolaan)
	router.GET("/getusulantindakan", middleware.DeserializeUser(), uc.andalalinController.GetUsulan)
	router.GET("/getdetailusulan/:id_andalalin", middleware.DeserializeUser(), uc.andalalinController.GetDetailUsulan)
	router.POST("/tindakanusulan/:id_andalalin/:jenis_pelaksanaan", middleware.DeserializeUser(), uc.andalalinController.TindakanPengelolaan)
	router.DELETE("/hapususulan/:id_andalalin", middleware.DeserializeUser(), uc.andalalinController.HapusUsulan)

	router.GET("/getallandalalintiket/:status", middleware.DeserializeUser(), uc.andalalinController.GetAllAndalalinByTiketLevel2)
}
