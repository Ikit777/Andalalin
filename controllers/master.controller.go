package controllers

import (
	"net/http"

	"github.com/Ikit777/E-Andalalin/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DataMasterControll struct {
	DB *gorm.DB
}

func NewDataMasterControll(DB *gorm.DB) DataMasterControll {
	return DataMasterControll{DB}
}

func (dm *DataMasterControll) GetDataMaster(ctx *gin.Context) {
	var data models.DataMaster

	results := dm.DB.First(&data)

	if results.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": results.Error})
		return
	}

	respone := struct {
		Lokasi             []string       `json:"lokasi_pengambilan,omitempty"`
		JenisRencana       []string       `json:"jenis_rencana,omitempty"`
		RencanaPembangunan models.Rencana `json:"rencana_pembangunan,omitempty"`
	}{
		Lokasi:             data.LokasiPengambilan,
		JenisRencana:       data.JenisRencanaPembangunan,
		RencanaPembangunan: data.RencanaPembangunan,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": respone})
}
