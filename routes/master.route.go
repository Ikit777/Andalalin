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

	router.POST("/tambahlokasi/:id/:lokasi", middleware.DeserializeUser(), dm.dataMasterController.TambahLokasi)
	router.POST("/hapuslokasi/:id/:lokasi", middleware.DeserializeUser(), dm.dataMasterController.HapusLokasi)
	router.POST("/editlokasi/:id/:lokasi/:new_lokasi", middleware.DeserializeUser(), dm.dataMasterController.EditLokasi)

	router.POST("/tambahkategori/:id/:kategori", middleware.DeserializeUser(), dm.dataMasterController.TambahKategori)
	router.POST("/hapuskategori/:id/:kategori", middleware.DeserializeUser(), dm.dataMasterController.HapusKategori)
	router.POST("/editkategori/:id/:kategori/:new_kategori", middleware.DeserializeUser(), dm.dataMasterController.EditKategori)

	router.POST("/tambahpembangunan/:id/:kategori/:rencana", middleware.DeserializeUser(), dm.dataMasterController.TambahJenisRencanaPembangunan)
	router.POST("/hapuspembangunan/:id/:kategori/:rencana", middleware.DeserializeUser(), dm.dataMasterController.HapusJenisRencanaPembangunan)
}
