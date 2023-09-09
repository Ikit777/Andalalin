package controllers

import (
	"net/http"

	"github.com/Ikit777/E-Andalalin/initializers"
	"github.com/Ikit777/E-Andalalin/models"
	"github.com/Ikit777/E-Andalalin/repository"
	"github.com/Ikit777/E-Andalalin/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
		IdDataMaster        uuid.UUID                  `json:"id_data_master,omitempty"`
		Lokasi              []string                   `json:"lokasi_pengambilan,omitempty"`
		JenisRencana        []string                   `json:"jenis_rencana,omitempty"`
		RencanaPembangunan  []models.Rencana           `json:"rencana_pembangunan,omitempty"`
		PersyaratanTambahan models.PersyaratanTambahan `json:"persyaratan_tambahan,omitempty"`
	}{
		IdDataMaster:        data.IdDataMaster,
		Lokasi:              data.LokasiPengambilan,
		JenisRencana:        data.JenisRencanaPembangunan,
		RencanaPembangunan:  data.RencanaPembangunan,
		PersyaratanTambahan: data.PersyaratanTambahan,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": respone})
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func (dm *DataMasterControler) TambahLokasi(ctx *gin.Context) {
	lokasi := ctx.Param("lokasi")
	id := ctx.Param("id")

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

	var master models.DataMaster

	resultsData := dm.DB.Where("id_data_master", id).First(&master)

	if resultsData.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": resultsData.Error})
		return
	}

	exist := contains(master.LokasiPengambilan, lokasi)

	if exist {
		ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "Data sudah ada"})
		return
	}

	master.LokasiPengambilan = append(master.LokasiPengambilan, lokasi)

	resultsSave := dm.DB.Save(&master)
	if resultsSave.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": resultsSave.Error})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (dm *DataMasterControler) HapusLokasi(ctx *gin.Context) {
	lokasi := ctx.Param("lokasi")
	id := ctx.Param("id")

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

	var master models.DataMaster

	resultsData := dm.DB.Where("id_data_master", id).First(&master)

	if resultsData.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": resultsData.Error})
		return
	}

	for i, item := range master.LokasiPengambilan {
		if item == lokasi {
			master.LokasiPengambilan = append(master.LokasiPengambilan[:i], master.LokasiPengambilan[i+1:]...)
			break
		}
	}

	resultsSave := dm.DB.Save(&master)
	if resultsSave.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": resultsSave.Error})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (dm *DataMasterControler) EditLokasi(ctx *gin.Context) {
	lokasi := ctx.Param("lokasi")
	newLokasi := ctx.Param("new_lokasi")
	id := ctx.Param("id")

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

	var master models.DataMaster

	resultsData := dm.DB.Where("id_data_master", id).First(&master)

	if resultsData.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": resultsData.Error})
		return
	}

	itemIndex := -1

	for i, item := range master.LokasiPengambilan {
		if item == lokasi {
			itemIndex = i
			break
		}
	}

	if itemIndex != -1 {
		master.LokasiPengambilan[itemIndex] = newLokasi
	}

	resultsSave := dm.DB.Save(&master)
	if resultsSave.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": resultsSave.Error})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (dm *DataMasterControler) TambahKategori(ctx *gin.Context) {
	kategori := ctx.Param("kategori")
	id := ctx.Param("id")

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

	var master models.DataMaster

	resultsData := dm.DB.Where("id_data_master", id).First(&master)

	if resultsData.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": resultsData.Error})
		return
	}

	exist := contains(master.JenisRencanaPembangunan, kategori)

	if exist {
		ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "Data sudah ada"})
		return
	}

	master.JenisRencanaPembangunan = append(master.JenisRencanaPembangunan, kategori)

	resultsSave := dm.DB.Save(&master)
	if resultsSave.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": resultsSave.Error})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (dm *DataMasterControler) HapusKategori(ctx *gin.Context) {
	kategori := ctx.Param("kategori")
	id := ctx.Param("id")

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

	var master models.DataMaster

	resultsData := dm.DB.Where("id_data_master", id).First(&master)

	if resultsData.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": resultsData.Error})
		return
	}

	for i, item := range master.JenisRencanaPembangunan {
		if item == kategori {
			master.JenisRencanaPembangunan = append(master.JenisRencanaPembangunan[:i], master.JenisRencanaPembangunan[i+1:]...)
			break
		}
	}

	for i, item := range master.RencanaPembangunan {
		if item.Kategori == kategori {
			master.RencanaPembangunan = append(master.RencanaPembangunan[:i], master.RencanaPembangunan[i+1:]...)
			break
		}
	}

	resultsSave := dm.DB.Save(&master)
	if resultsSave.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": resultsSave.Error})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (dm *DataMasterControler) EditKategori(ctx *gin.Context) {
	kategori := ctx.Param("kategori")
	newKategori := ctx.Param("new_kategori")
	id := ctx.Param("id")

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

	var master models.DataMaster

	resultsData := dm.DB.Where("id_data_master", id).First(&master)

	if resultsData.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": resultsData.Error})
		return
	}

	itemIndex := -1

	for i, item := range master.JenisRencanaPembangunan {
		if item == kategori {
			itemIndex = i
			break
		}
	}

	if itemIndex != -1 {
		master.JenisRencanaPembangunan[itemIndex] = newKategori
	}

	resultsSave := dm.DB.Save(&master)
	if resultsSave.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": resultsSave.Error})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (dm *DataMasterControler) TambahJenisRencanaPembangunan(ctx *gin.Context) {
	kategori := ctx.Param("kategori")
	rencana := ctx.Param("rencana")
	id := ctx.Param("id")

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

	var master models.DataMaster

	resultsData := dm.DB.Where("id_data_master", id).First(&master)

	if resultsData.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": resultsData.Error})
		return
	}

	kategoriExists := false
	jenisExists := false
	itemIndex := 0

	for i := range master.RencanaPembangunan {
		if master.RencanaPembangunan[i].Kategori == kategori {
			kategoriExists = true
			itemIndex = i
			for _, item := range master.RencanaPembangunan[i].JenisRencana {
				if item == rencana {
					jenisExists = true
					ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "Data sudah ada"})
					return
				}
			}
		}
	}

	if !kategoriExists {
		jenis := models.Rencana{
			Kategori:     kategori,
			JenisRencana: []string{rencana},
		}
		master.RencanaPembangunan = append(master.RencanaPembangunan, jenis)
	}

	if !jenisExists && kategoriExists {
		master.RencanaPembangunan[itemIndex].JenisRencana = append(master.RencanaPembangunan[itemIndex].JenisRencana, rencana)
	}

	resultsSave := dm.DB.Save(&master)
	if resultsSave.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": resultsSave.Error})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}
