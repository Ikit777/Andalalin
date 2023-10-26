package routes

import (
	"github.com/Ikit777/E-Andalalin/controllers"
	"github.com/Ikit777/E-Andalalin/middleware"
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

	router.GET("/check", dm.dataMasterController.CheckDataMaster)

	router.POST("/tambahlokasi/:id/:lokasi", middleware.DeserializeUser(), dm.dataMasterController.TambahLokasi)
	router.POST("/hapuslokasi/:id/:lokasi", middleware.DeserializeUser(), dm.dataMasterController.HapusLokasi)
	router.POST("/editlokasi/:id/:lokasi/:new_lokasi", middleware.DeserializeUser(), dm.dataMasterController.EditLokasi)

	router.POST("/tambahkategori/:id/:kategori", middleware.DeserializeUser(), dm.dataMasterController.TambahKategori)
	router.POST("/hapuskategori/:id/:kategori", middleware.DeserializeUser(), dm.dataMasterController.HapusKategori)
	router.POST("/editkategori/:id/:kategori/:new_kategori", middleware.DeserializeUser(), dm.dataMasterController.EditKategori)

	router.POST("/tambahpembangunan/:id/:kategori/:rencana/:kriteria/:satuan", middleware.DeserializeUser(), dm.dataMasterController.TambahJenisRencanaPembangunan)
	router.POST("/hapuspembangunan/:id/:kategori/:rencana", middleware.DeserializeUser(), dm.dataMasterController.HapusJenisRencanaPembangunan)
	router.POST("/editpembangunan/:id/:kategori/:rencana/:rencana_new/:kriteria/:satuan", middleware.DeserializeUser(), dm.dataMasterController.EditJenisRencanaPembangunan)

	router.POST("/tambahkategoriperlengkapan/:id/:kategori", middleware.DeserializeUser(), dm.dataMasterController.TambahKategoriPerlengkapan)
	router.POST("/hapuskategoriperlengkapan/:id/:kategori", middleware.DeserializeUser(), dm.dataMasterController.HapusKategoriPerlengkapan)
	router.POST("/editkategoriperlengkapan/:id/:kategori/:new_kategori", middleware.DeserializeUser(), dm.dataMasterController.EditKategoriPerlengkapan)

	router.POST("/tambahperlengkapan/:id/:kategori/:perlengkapan", middleware.DeserializeUser(), dm.dataMasterController.TambahPerlengkapan)
	router.POST("/hapusperlengkapan/:id/:kategori/:perlengkapan", middleware.DeserializeUser(), dm.dataMasterController.HapuspPerlengkapan)
	router.POST("/editperlengkapan/:id/:kategori/:perlengkapan/:perlengkapan_new", middleware.DeserializeUser(), dm.dataMasterController.EditPerlengkapan)

	router.POST("/tambahpersyaratanandalalin/:id", middleware.DeserializeUser(), dm.dataMasterController.TambahPersyaratanAndalalin)
	router.POST("/hapuspersyaratanandalalin/:id/:persyaratan", middleware.DeserializeUser(), dm.dataMasterController.HapusPersyaratanAndalalin)
	router.POST("/editpersyaratanandalalin/:id/:persyaratan", middleware.DeserializeUser(), dm.dataMasterController.EditPersyaratanAndalalin)

	router.POST("/tambahpersyaratanperlalin/:id", middleware.DeserializeUser(), dm.dataMasterController.TambahPersyaratanPerlalin)
	router.POST("/hapuspersyaratanperlalin/:id/:persyaratan", middleware.DeserializeUser(), dm.dataMasterController.HapusPersyaratanPerlalin)
	router.POST("/editpersyaratanperlalin/:id/:persyaratan", middleware.DeserializeUser(), dm.dataMasterController.EditPersyaratanPerlalin)

	router.POST("/tambahprovinsi/:id/:provinsi", middleware.DeserializeUser(), dm.dataMasterController.TambahProvinsi)
	router.POST("/hapusprovinsi/:id/:provinsi", middleware.DeserializeUser(), dm.dataMasterController.HapusProvinsi)
	router.POST("/editprovinsi/:id/:provinsi/:new_provinsi", middleware.DeserializeUser(), dm.dataMasterController.EditProvinsi)

	router.POST("/tambahkabupaten/:id/:provinsi/:kabupaten", middleware.DeserializeUser(), dm.dataMasterController.TambahKabupaten)
	router.POST("/hapuskabupaten/:id/:kabupaten", middleware.DeserializeUser(), dm.dataMasterController.HapusKabupaten)
	router.POST("/editkabupaten/:id/:kabupaten/:new_kabupaten", middleware.DeserializeUser(), dm.dataMasterController.EditKabupaten)

	router.POST("/tambahkecamatan/:id/:kabupaten/:kecamatan", middleware.DeserializeUser(), dm.dataMasterController.TambahKecamatan)
	router.POST("/hapuskecamatan/:id/:kecamatan", middleware.DeserializeUser(), dm.dataMasterController.HapusKecamatan)
	router.POST("/editkecamatan/:id/:kecamatan/:new_kecamatan", middleware.DeserializeUser(), dm.dataMasterController.EditKecamatan)

	router.POST("/tambahkelurahan/:id/:kecamatan/:kelurahan", middleware.DeserializeUser(), dm.dataMasterController.TambahKelurahan)
	router.POST("/hapuskelurahan/:id/:kelurahan", middleware.DeserializeUser(), dm.dataMasterController.HapusKelurahan)
	router.POST("/editkelurahan/:id/:kelurahan/:new_kelurahan", middleware.DeserializeUser(), dm.dataMasterController.EditKelurahan)
}
