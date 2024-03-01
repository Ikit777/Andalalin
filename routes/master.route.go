package routes

import (
	"andalalin/controllers"
	"andalalin/middleware"

	"github.com/gin-gonic/gin"
)

type DataMasterRouteController struct {
	dataMasterController controllers.DataMasterControler
}

func NewDataMasterRouteController(dataMasterController controllers.DataMasterControler) DataMasterRouteController {
	return DataMasterRouteController{dataMasterController}
}

func (dm *DataMasterRouteController) DataMasterRoute(rg *gin.RouterGroup) {
	router := rg.Group("/master")

	router.GET("/andalalin", dm.dataMasterController.GetDataMaster)
	router.GET("/tipe/:tipe", dm.dataMasterController.GetDataMasterByType)

	router.GET("/check", dm.dataMasterController.CheckDataMaster)

	router.POST("/tambah/lokasi/:id", middleware.DeserializeUser(), dm.dataMasterController.TambahLokasi)
	router.POST("/hapus/lokasi/:id", middleware.DeserializeUser(), dm.dataMasterController.HapusLokasi)
	router.POST("/edit/lokasi/:id", middleware.DeserializeUser(), dm.dataMasterController.EditLokasi)

	router.POST("/tambah/kategori/rencana/:id", middleware.DeserializeUser(), dm.dataMasterController.TambahKategori)
	router.POST("/hapus/kategori/rencana/:id", middleware.DeserializeUser(), dm.dataMasterController.HapusKategori)
	router.POST("/edit/kategori/rencana/:id", middleware.DeserializeUser(), dm.dataMasterController.EditKategori)

	router.POST("/tambah/jenis/rencana/:id", middleware.DeserializeUser(), dm.dataMasterController.TambahJenisRencanaPembangunan)
	router.POST("/hapus/jenis/rencana/:id", middleware.DeserializeUser(), dm.dataMasterController.HapusJenisRencanaPembangunan)
	router.POST("/edit/jenis/rencana/:id", middleware.DeserializeUser(), dm.dataMasterController.EditJenisRencanaPembangunan)

	router.POST("/tambah/kategori/utama/:id", middleware.DeserializeUser(), dm.dataMasterController.TambahKategoriUtamaPerlengkapan)
	router.POST("/hapus/kategori/utama/:id", middleware.DeserializeUser(), dm.dataMasterController.HapusKategoriUtamaPerlengkapan)
	router.POST("/edit/kategori/utama/:id", middleware.DeserializeUser(), dm.dataMasterController.EditKategoriUtamaPerlengkapan)

	router.POST("/tambah/kategori/perlengkapan/:id", middleware.DeserializeUser(), dm.dataMasterController.TambahKategoriPerlengkapan)
	router.POST("/hapus/kategori/perlengkapan/:id", middleware.DeserializeUser(), dm.dataMasterController.HapusKategoriPerlengkapan)
	router.POST("/edit/kategori/perlengkapan/:id", middleware.DeserializeUser(), dm.dataMasterController.EditKategoriPerlengkapan)

	router.POST("/tambah/jenis/perlengkapan/:id", middleware.DeserializeUser(), dm.dataMasterController.TambahPerlengkapan)
	router.POST("/hapus/jenis/perlengkapan/:id", middleware.DeserializeUser(), dm.dataMasterController.HapuspPerlengkapan)
	router.POST("/edit/jenis/perlengkapan/:id", middleware.DeserializeUser(), dm.dataMasterController.EditPerlengkapan)

	router.POST("/tambah/persyaratan/andalalin/:id", middleware.DeserializeUser(), dm.dataMasterController.TambahPersyaratanAndalalin)
	router.POST("/hapus/persyaratan/andalalin/:id", middleware.DeserializeUser(), dm.dataMasterController.HapusPersyaratanAndalalin)
	router.POST("/edit/persyaratan/andalalin/:id", middleware.DeserializeUser(), dm.dataMasterController.EditPersyaratanAndalalin)

	router.POST("/tambah/persyaratan/perlalin/:id", middleware.DeserializeUser(), dm.dataMasterController.TambahPersyaratanPerlalin)
	router.POST("/hapus/persyaratan/perlalin/:id", middleware.DeserializeUser(), dm.dataMasterController.HapusPersyaratanPerlalin)
	router.POST("/edit/persyaratan/perlalin/:id", middleware.DeserializeUser(), dm.dataMasterController.EditPersyaratanPerlalin)

	router.POST("/tambah/provinsi/:id", middleware.DeserializeUser(), dm.dataMasterController.TambahProvinsi)
	router.POST("/hapus/provinsi/:id", middleware.DeserializeUser(), dm.dataMasterController.HapusProvinsi)
	router.POST("/edit/provinsi/:id", middleware.DeserializeUser(), dm.dataMasterController.EditProvinsi)

	router.POST("/tambah/kabupaten/:id", middleware.DeserializeUser(), dm.dataMasterController.TambahKabupaten)
	router.POST("/hapus/kabupaten/:id", middleware.DeserializeUser(), dm.dataMasterController.HapusKabupaten)
	router.POST("/edit/kabupaten/:id", middleware.DeserializeUser(), dm.dataMasterController.EditKabupaten)

	router.POST("/tambah/kecamatan/:id", middleware.DeserializeUser(), dm.dataMasterController.TambahKecamatan)
	router.POST("/hapus/kecamatan/:id", middleware.DeserializeUser(), dm.dataMasterController.HapusKecamatan)
	router.POST("/edit/kecamatan/:id", middleware.DeserializeUser(), dm.dataMasterController.EditKecamatan)

	router.POST("/tambah/kelurahan/:id", middleware.DeserializeUser(), dm.dataMasterController.TambahKelurahan)
	router.POST("/hapus/kelurahan/:id", middleware.DeserializeUser(), dm.dataMasterController.HapusKelurahan)
	router.POST("/edit/kelurahan/:id", middleware.DeserializeUser(), dm.dataMasterController.EditKelurahan)

	router.POST("/tambah/jenis/proyek/:id", middleware.DeserializeUser(), dm.dataMasterController.TambahJenisProyek)
	router.POST("/hapus/jenis/proyek/:id", middleware.DeserializeUser(), dm.dataMasterController.HapusJenisProyek)
	router.POST("/edit/jenis/proyek/:id", middleware.DeserializeUser(), dm.dataMasterController.EditJenisProyek)

	router.POST("/tambah/jalan/:id", middleware.DeserializeUser(), dm.dataMasterController.TambahJalan)
	router.POST("/hapus/jalan/:id", middleware.DeserializeUser(), dm.dataMasterController.HapusJalan)
	router.POST("/edit/jalan/:id", middleware.DeserializeUser(), dm.dataMasterController.EditJalan)

	router.POST("/tambah/panduan/:id", middleware.DeserializeUser(), dm.dataMasterController.TambahPanduan)
	router.POST("/hapus/panduan/:id", middleware.DeserializeUser(), dm.dataMasterController.HapusPanduan)
	router.POST("/edit/panduan/:id", middleware.DeserializeUser(), dm.dataMasterController.EditPanduan)
	router.POST("/get/panduan/:id", middleware.DeserializeUser(), dm.dataMasterController.GetPanduan)
}
