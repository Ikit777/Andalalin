package routes

import (
	"github.com/Ikit777/E-Andalalin/controllers"
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
}
