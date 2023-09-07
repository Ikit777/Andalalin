package controllers

import (
	"net/http"

	"github.com/Ikit777/E-Andalalin/initializers"
	"github.com/Ikit777/E-Andalalin/models"
	"github.com/Ikit777/E-Andalalin/repository"
	"github.com/Ikit777/E-Andalalin/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DataMasterControler struct {
	DB *gorm.DB
}

func NewDataMasterControler(DB *gorm.DB) DataMasterControler {
	return DataMasterControler{DB}
}

func (dm *DataMasterControler) GetDataMaster(ctx *gin.Context) {
	var data models.DataMaster

	results := dm.DB.First(&data)

	if results.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": results.Error})
		return
	}

	respone := struct {
		Lokasi              []string                   `json:"lokasi_pengambilan,omitempty"`
		JenisRencana        []string                   `json:"jenis_rencana,omitempty"`
		RencanaPembangunan  []models.Rencana           `json:"rencana_pembangunan,omitempty"`
		PersyaratanTambahan models.PersyaratanTambahan `json:"persyaratan_tambahan,omitempty"`
	}{
		Lokasi:              data.LokasiPengambilan,
		JenisRencana:        data.JenisRencanaPembangunan,
		RencanaPembangunan:  data.RencanaPembangunan,
		PersyaratanTambahan: data.PersyaratanTambahan,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": respone})
}

func (dm *DataMasterControler) TambahLokasi(ctx *gin.Context) {
	var payload *models.LokasiInput

	config, _ := initializers.LoadConfig(".")

	accessUser := ctx.MustGet("accessUser").(string)

	claim, error := utils.ValidateToken(accessUser, config.AccessTokenPublicKey)
	if error != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": error.Error()})
		return
	}

	credential := claim.Credentials[repository.ProductAddCredential]

	if !credential {
		// Return status 403 and permission denied error message.
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": true,
			"msg":   "Permission denied",
		})
		return
	}

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	lokasi := models.LokasiInput{
		Lokasi: payload.Lokasi,
	}

	var master models.DataMaster

	results := dm.DB.First(&master.LokasiPengambilan, models.Lokasi{payload.Lokasi})
	if results.Error != nil {
		ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "Data sudah ada"})
		return
	}

	dm.DB.Model(&master).Association("LokasiPengambilan").Append(&lokasi)

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})

}
