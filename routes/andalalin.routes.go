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

	router.GET("/permohonan", middleware.DeserializeUser(), uc.andalalinController.GetPermohonan)
	router.GET("/userpermohonan", middleware.DeserializeUser(), uc.andalalinController.GetPermohonanByIdUser)
	router.GET("/detailpermohonan/:id_andalalin", middleware.DeserializeUser(), uc.andalalinController.GetPermohonanByIdAndalalin)
	router.GET("/getdokumen/:id_andalalin/:dokumen", middleware.DeserializeUser(), uc.andalalinController.GetDokumen)
	router.GET("/statuspermohonan/:status_andalalin", middleware.DeserializeUser(), uc.andalalinController.GetPermohonanByStatus)

	router.GET("/petugaspermohonan/:status", middleware.DeserializeUser(), uc.andalalinController.GetAndalalinTicketLevel2)

	router.POST("/updatestatus/:id_andalalin/:status", middleware.DeserializeUser(), uc.andalalinController.UpdateStatusPermohonan)
	router.POST("/tolakpermohonan/:id_andalalin/:pertimbangan", middleware.DeserializeUser(), uc.andalalinController.TolakPermohonan)
	router.POST("/tundapermohonan/:id_andalalin/:pertimbangan", middleware.DeserializeUser(), uc.andalalinController.TundaPermohonan)
	router.POST("/lanjutkanpermohonan/:id_andalalin", middleware.DeserializeUser(), uc.andalalinController.LanjutkanPermohonan)

	router.POST("/persyaratantidaksesuai/:id_andalalin", middleware.DeserializeUser(), uc.andalalinController.PersyaratanTidakSesuai)
	router.POST("/persyaratanterpenuhi/:id_andalalin", middleware.DeserializeUser(), uc.andalalinController.PersyaratanTerpenuhi)
	router.POST("/checkadministrasi/:id_andalalin", middleware.DeserializeUser(), uc.andalalinController.CheckAdministrasi)
	router.POST("/uploaddokumen/:id_andalalin/:dokumen", middleware.DeserializeUser(), uc.andalalinController.UploadDokumen)

	router.POST("/pembuatansuratpernyataan/:id_andalalin", middleware.DeserializeUser(), uc.andalalinController.PembuatanSuratPernyataan)

	router.POST("/pembuatansuratkeputusan/:id_andalalin", middleware.DeserializeUser(), uc.andalalinController.PembuatanSuratKeputusan)

	router.POST("/updatepersyaratan/:id_andalalin", middleware.DeserializeUser(), uc.andalalinController.UpdatePersyaratan)

	router.POST("/pilihpetugas/:id_andalalin", middleware.DeserializeUser(), uc.andalalinController.TambahPetugas)
	router.POST("/gantipetugas/:id_andalalin", middleware.DeserializeUser(), uc.andalalinController.GantiPetugas)

	router.POST("/survey/:id_andalalin", middleware.DeserializeUser(), uc.andalalinController.IsiSurvey)
	router.GET("/getsurvey/:id_andalalin", middleware.DeserializeUser(), uc.andalalinController.GetSurvey)
	router.GET("/getallsurvey", middleware.DeserializeUser(), uc.andalalinController.GetAllSurvey)

	router.POST("/pemeriksaansuratkeputusan/:id_andalalin", middleware.DeserializeUser(), uc.andalalinController.PemeriksaanSuratKeputusan)

	router.POST("/checkkelengkapanakhir/:id_andalalin", middleware.DeserializeUser(), uc.andalalinController.CheckKelengkapanAkhir)

	router.POST("/usulantindakan/:id_andalalin", middleware.DeserializeUser(), uc.andalalinController.UsulanTindakanPengelolaan)
	router.GET("/getusulantindakan", middleware.DeserializeUser(), uc.andalalinController.GetUsulan)
	router.GET("/getdetailusulan/:id_andalalin", middleware.DeserializeUser(), uc.andalalinController.GetDetailUsulan)
	router.POST("/tindakanusulan/:id_andalalin/:jenis_pelaksanaan", middleware.DeserializeUser(), uc.andalalinController.TindakanPengelolaan)
	router.DELETE("/hapususulan/:id_andalalin", middleware.DeserializeUser(), uc.andalalinController.HapusUsulan)

	router.GET("/getallandalalintiket/:status", middleware.DeserializeUser(), uc.andalalinController.GetAllAndalalinByTiketLevel2)

	router.POST("/laporansurvei/:id_andalalin", middleware.DeserializeUser(), uc.andalalinController.LaporanSurvei)
	router.POST("/keputusanhasil/:id_andalalin", middleware.DeserializeUser(), uc.andalalinController.KeputusanHasil)

	router.GET("/getpermohonanpemasangan", middleware.DeserializeUser(), uc.andalalinController.GetPermohonanPemasanganLalin)
	router.POST("/pemasanganperlengkapan/:id_andalalin", middleware.DeserializeUser(), uc.andalalinController.PemasanganPerlengkapanLaluLintas)
	router.GET("/getpemasanganperlengkapan/:id_andalalin", middleware.DeserializeUser(), uc.andalalinController.GetPemasangan)
	router.GET("/getallpemasanganperlengkapan", middleware.DeserializeUser(), uc.andalalinController.GetAllPemasangan)

	router.POST("/surveimandiri", middleware.DeserializeUser(), uc.andalalinController.IsiSurveyMandiri)
	router.GET("/detailsurveimandiri/:id_survei", middleware.DeserializeUser(), uc.andalalinController.GetSurveiMandiri)
	router.GET("/daftarsurveimandiri", middleware.DeserializeUser(), uc.andalalinController.GetAllSurveiMandiri)
	router.GET("/daftarsurveimandiribypetugas", middleware.DeserializeUser(), uc.andalalinController.GetAllSurveiMandiriByPetugas)
	router.POST("/terimasurvei/:id_survei/:keterangan", middleware.DeserializeUser(), uc.andalalinController.TerimaSurvei)

	router.POST("/surveikepuasan/:id_andalalin", middleware.DeserializeUser(), uc.andalalinController.SurveiKepuasan)
	router.GET("/ceksurveikepuasan/:id_andalalin", middleware.DeserializeUser(), uc.andalalinController.CekSurveiKepuasan)
	router.GET("/hasilsurveikepuasan", middleware.DeserializeUser(), uc.andalalinController.HasilSurveiKepuasan)
}
