package routes

import (
	"andalalin/controllers"
	"andalalin/middleware"

	"github.com/gin-gonic/gin"
)

type AndalalinRouteController struct {
	andalalinController controllers.AndalalinController
}

func NewRouteAndalalinController(andalalinController controllers.AndalalinController) AndalalinRouteController {
	return AndalalinRouteController{andalalinController}
}

func (uc *AndalalinRouteController) AndalalainRoute(rg *gin.RouterGroup) {

	router := rg.Group("permohonan")

	//Pengajuan
	router.POST("/pengajuan/andalalin", middleware.DeserializeUser(), uc.andalalinController.Pengajuan)
	router.POST("/pengajuan/perlalin", middleware.DeserializeUser(), uc.andalalinController.PengajuanPerlalin)

	//Get Permohonan
	router.GET("/all", middleware.DeserializeUser(), uc.andalalinController.GetPermohonan)
	router.GET("/user", middleware.DeserializeUser(), uc.andalalinController.GetPermohonanByIdUser)
	router.GET("/detail/:id_andalalin", middleware.DeserializeUser(), uc.andalalinController.GetPermohonanByIdAndalalin)
	router.GET("/status/:status_andalalin", middleware.DeserializeUser(), uc.andalalinController.GetPermohonanByStatus)
	router.GET("/petugas/:status", middleware.DeserializeUser(), uc.andalalinController.GetAndalalinTicketLevel2)

	//Get Berkas Permohonan
	router.GET("/berkas/:id_andalalin/:dokumen", middleware.DeserializeUser(), uc.andalalinController.GetDokumen)
	router.POST("/berkas/update/:id_andalalin", middleware.DeserializeUser(), uc.andalalinController.UpdateBerkas)

	//Menindaklanjuti Permohonan
	router.POST("/update/status/:id_andalalin/:status", middleware.DeserializeUser(), uc.andalalinController.UpdateStatusPermohonan)
	router.POST("/tolak/:id_andalalin/:pertimbangan", middleware.DeserializeUser(), uc.andalalinController.TolakPermohonan)
	router.POST("/tunda/:id_andalalin/:pertimbangan", middleware.DeserializeUser(), uc.andalalinController.TundaPermohonan)
	router.POST("/lanjutkan/:id_andalalin", middleware.DeserializeUser(), uc.andalalinController.LanjutkanPermohonan)

	//Menindaklanjuti Persyaratan Permohonan
	router.POST("/persyaratan/andalalin/pemeriksaan/:id_andalalin", middleware.DeserializeUser(), uc.andalalinController.CheckAdministrasi)
	router.POST("/persyaratan/perlalin/pemeriksaan/:id_andalalin", middleware.DeserializeUser(), uc.andalalinController.CheckAdministrasiPerlalin)

	//Upload Berkas
	router.POST("/upload/dokumen/:id_andalalin/:dokumen", middleware.DeserializeUser(), uc.andalalinController.UploadDokumen)

	//Pembuatan surat pernyataan dan keputusan
	router.POST("/pembuatan/pernyataan/:id_andalalin", middleware.DeserializeUser(), uc.andalalinController.PembuatanSuratPernyataan)
	router.POST("/pembuatan/keputusan/:id_andalalin", middleware.DeserializeUser(), uc.andalalinController.PembuatanSuratKeputusan)
	router.POST("/pembuatan/penyusun/:id_andalalin", middleware.DeserializeUser(), uc.andalalinController.PembuatanPenyusunDokumen)

	//Survei
	router.POST("/survey/:id_andalalin", middleware.DeserializeUser(), uc.andalalinController.IsiSurvey)
	router.GET("/survey/detail/:id_andalalin", middleware.DeserializeUser(), uc.andalalinController.GetSurvey)
	router.GET("/survey/all", middleware.DeserializeUser(), uc.andalalinController.GetAllSurvey)

	//Menindaklanjuti petugas permohonan
	router.POST("/survey/petugas/pilih/:id_andalalin", middleware.DeserializeUser(), uc.andalalinController.TambahPetugas)
	router.POST("/survey/petugas/ganti/:id_andalalin", middleware.DeserializeUser(), uc.andalalinController.GantiPetugas)

	//Pemeriksaan surat keputusan
	router.POST("/pemeriksaan/keputusan/:id_andalalin", middleware.DeserializeUser(), uc.andalalinController.PemeriksaanSuratKeputusan)

	//Pemeriksaan dokumen andalalin
	router.POST("/pemeriksaan/dokumen/:id_andalalin", middleware.DeserializeUser(), uc.andalalinController.PemeriksaanDokumenAndalalin)

	//Pemeriksanaan kelengkapan akhir
	router.POST("/pemeriksaan/kelengkapan/:id_andalalin", middleware.DeserializeUser(), uc.andalalinController.CheckKelengkapanAkhir)

	//Usulan tindakan
	router.POST("/usulan/:id_andalalin", middleware.DeserializeUser(), uc.andalalinController.UsulanTindakanPengelolaan)
	router.GET("/usulan/all", middleware.DeserializeUser(), uc.andalalinController.GetUsulan)
	router.GET("/usulan/detail/:id_andalalin", middleware.DeserializeUser(), uc.andalalinController.GetDetailUsulan)
	router.POST("/usulan/tindakan/:id_andalalin/:jenis_pelaksanaan", middleware.DeserializeUser(), uc.andalalinController.TindakanPengelolaan)
	router.DELETE("/usulan/hapus/:id_andalalin", middleware.DeserializeUser(), uc.andalalinController.HapusUsulan)

	//Laporan survei
	router.POST("/laporan/survei/:id_andalalin", middleware.DeserializeUser(), uc.andalalinController.LaporanSurvei)

	//Keputusan hasil
	router.POST("/keputusan/hasil/:id_andalalin", middleware.DeserializeUser(), uc.andalalinController.KeputusanHasil)

	//Pemasangan perlalin
	router.GET("/pemasangan", middleware.DeserializeUser(), uc.andalalinController.GetPermohonanPemasanganLalin)
	router.POST("/pemasangan/pasang/:id_andalalin", middleware.DeserializeUser(), uc.andalalinController.PemasanganPerlengkapanLaluLintas)
	router.GET("/pemasangan/detail/:id_andalalin", middleware.DeserializeUser(), uc.andalalinController.GetPemasangan)
	router.GET("/pemasangan/all", middleware.DeserializeUser(), uc.andalalinController.GetAllPemasangan)
}
