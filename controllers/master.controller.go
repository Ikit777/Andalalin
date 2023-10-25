package controllers

import (
	"archive/zip"
	"encoding/base64"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

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

type file struct {
	Name string
	File []byte
}

func NewDataMasterControler(DB *gorm.DB) DataMasterControler {
	return DataMasterControler{DB}
}

func (dm *DataMasterControler) GetDataMaster(ctx *gin.Context) {
	var master models.DataMaster

	results := dm.DB.First(&master)

	if results.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": results.Error})
		return
	}

	respone := struct {
		IdDataMaster           uuid.UUID                  `json:"id_data_master,omitempty"`
		Lokasi                 []string                   `json:"lokasi_pengambilan,omitempty"`
		JenisRencana           []string                   `json:"jenis_rencana,omitempty"`
		RencanaPembangunan     []models.Rencana           `json:"rencana_pembangunan,omitempty"`
		KategoriPerlengkapan   []string                   `json:"kategori_perlengkapan,omitempty"`
		PerlengkapanLaluLintas []models.JenisPerlengkapan `json:"perlengkapan,omitempty"`
		Persyaratan            models.Persyaratan         `json:"persyaratan,omitempty"`
		Provinsi               []models.Provinsi          `json:"provinsi,omitempty"`
		Kabupaten              []models.Kabupaten         `json:"kabupaten,omitempty"`
		Kecamatan              []models.Kecamatan         `json:"kecamatan,omitempty"`
		Kelurahan              []models.Kelurahan         `json:"kelurahan,omitempty"`
		UpdatedAt              string                     `json:"update,omitempty"`
	}{
		IdDataMaster:           master.IdDataMaster,
		Lokasi:                 master.LokasiPengambilan,
		JenisRencana:           master.JenisRencanaPembangunan,
		RencanaPembangunan:     master.RencanaPembangunan,
		KategoriPerlengkapan:   master.KategoriPerlengkapan,
		PerlengkapanLaluLintas: master.PerlengkapanLaluLintas,
		Persyaratan:            master.Persyaratan,
		Provinsi:               master.Provinsi,
		Kabupaten:              master.Kabupaten,
		Kecamatan:              master.Kecamatan,
		Kelurahan:              master.Kelurahan,
		UpdatedAt:              master.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": respone})
}

func (dm *DataMasterControler) CheckDataMaster(ctx *gin.Context) {
	var master models.DataMaster

	results := dm.DB.First(&master)

	if results.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": results.Error})
		return
	}

	respone := struct {
		IdDataMaster uuid.UUID `json:"id_data_master,omitempty"`
		UpdatedAt    string    `json:"update,omitempty"`
	}{
		IdDataMaster: master.IdDataMaster,
		UpdatedAt:    master.UpdatedAt,
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

func compressFiles(zipFileName string, fileData []file) error {
	zipFile, err := os.Create(zipFileName)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	for _, data := range fileData {
		tmpFile, _ := os.CreateTemp("", "persyaratan.pdf")
		defer os.Remove(tmpFile.Name())

		_, _ = tmpFile.Write(data.File)
		tmpFile.Close()

		file, err := os.Open(tmpFile.Name())
		if err != nil {
			return err
		}
		defer file.Close()

		fileInfo, err := file.Stat()
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(fileInfo)
		if err != nil {
			return err
		}

		header.Name = data.Name

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		_, err = io.Copy(writer, file)
		if err != nil {
			return err
		}
	}

	return nil
}

func (dm *DataMasterControler) TambahLokasi(ctx *gin.Context) {
	lokasi := ctx.Param("lokasi")
	id := ctx.Param("id")

	config, _ := initializers.LoadConfig()

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

	loc, _ := time.LoadLocation("Asia/Singapore")
	now := time.Now().In(loc).Format("02-01-2006")

	master.UpdatedAt = now + " " + time.Now().In(loc).Format("15:04:05")

	resultsSave := dm.DB.Save(&master)
	if resultsSave.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": resultsSave.Error})
		return
	}

	respone := struct {
		Lokasi    []string `json:"lokasi_pengambilan,omitempty"`
		UpdatedAt string   `json:"update,omitempty"`
	}{
		UpdatedAt: master.UpdatedAt,
		Lokasi:    master.LokasiPengambilan,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": respone})
}

func (dm *DataMasterControler) HapusLokasi(ctx *gin.Context) {
	lokasi := ctx.Param("lokasi")
	id := ctx.Param("id")

	config, _ := initializers.LoadConfig()

	accessUser := ctx.MustGet("accessUser").(string)

	claim, error := utils.ValidateToken(accessUser, config.AccessTokenPublicKey)
	if error != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": error.Error()})
		return
	}

	credential := claim.Credentials[repository.ProductDeleteCredential]

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

	loc, _ := time.LoadLocation("Asia/Singapore")
	now := time.Now().In(loc).Format("02-01-2006")

	master.UpdatedAt = now + " " + time.Now().In(loc).Format("15:04:05")

	resultsSave := dm.DB.Save(&master)
	if resultsSave.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": resultsSave.Error})
		return
	}

	respone := struct {
		Lokasi    []string `json:"lokasi_pengambilan,omitempty"`
		UpdatedAt string   `json:"update,omitempty"`
	}{
		UpdatedAt: master.UpdatedAt,
		Lokasi:    master.LokasiPengambilan,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": respone})
}

func (dm *DataMasterControler) EditLokasi(ctx *gin.Context) {
	lokasi := ctx.Param("lokasi")
	newLokasi := ctx.Param("new_lokasi")
	id := ctx.Param("id")

	config, _ := initializers.LoadConfig()

	accessUser := ctx.MustGet("accessUser").(string)

	claim, error := utils.ValidateToken(accessUser, config.AccessTokenPublicKey)
	if error != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": error.Error()})
		return
	}

	credential := claim.Credentials[repository.ProductUpdateCredential]

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

	loc, _ := time.LoadLocation("Asia/Singapore")
	now := time.Now().In(loc).Format("02-01-2006")

	master.UpdatedAt = now + " " + time.Now().In(loc).Format("15:04:05")

	resultsSave := dm.DB.Save(&master)
	if resultsSave.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": resultsSave.Error})
		return
	}

	respone := struct {
		Lokasi    []string `json:"lokasi_pengambilan,omitempty"`
		UpdatedAt string   `json:"update,omitempty"`
	}{
		UpdatedAt: master.UpdatedAt,
		Lokasi:    master.LokasiPengambilan,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": respone})
}

func (dm *DataMasterControler) TambahKategori(ctx *gin.Context) {
	kategori := ctx.Param("kategori")
	id := ctx.Param("id")

	config, _ := initializers.LoadConfig()

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

	loc, _ := time.LoadLocation("Asia/Singapore")
	now := time.Now().In(loc).Format("02-01-2006")

	master.UpdatedAt = now + " " + time.Now().In(loc).Format("15:04:05")

	resultsSave := dm.DB.Save(&master)
	if resultsSave.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": resultsSave.Error})
		return
	}

	respone := struct {
		JenisRencana []string `json:"jenis_rencana,omitempty"`
		UpdatedAt    string   `json:"update,omitempty"`
	}{
		JenisRencana: master.JenisRencanaPembangunan,
		UpdatedAt:    master.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": respone})
}

func (dm *DataMasterControler) HapusKategori(ctx *gin.Context) {
	kategori := ctx.Param("kategori")
	id := ctx.Param("id")

	config, _ := initializers.LoadConfig()

	accessUser := ctx.MustGet("accessUser").(string)

	claim, error := utils.ValidateToken(accessUser, config.AccessTokenPublicKey)
	if error != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": error.Error()})
		return
	}

	credential := claim.Credentials[repository.ProductDeleteCredential]

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

	loc, _ := time.LoadLocation("Asia/Singapore")
	now := time.Now().In(loc).Format("02-01-2006")

	master.UpdatedAt = now + " " + time.Now().In(loc).Format("15:04:05")

	resultsSave := dm.DB.Save(&master)
	if resultsSave.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": resultsSave.Error})
		return
	}

	respone := struct {
		JenisRencana []string `json:"jenis_rencana,omitempty"`
		UpdatedAt    string   `json:"update,omitempty"`
	}{
		JenisRencana: master.JenisRencanaPembangunan,
		UpdatedAt:    master.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": respone})
}

func (dm *DataMasterControler) EditKategori(ctx *gin.Context) {
	kategori := ctx.Param("kategori")
	newKategori := ctx.Param("new_kategori")
	id := ctx.Param("id")

	config, _ := initializers.LoadConfig()

	accessUser := ctx.MustGet("accessUser").(string)

	claim, error := utils.ValidateToken(accessUser, config.AccessTokenPublicKey)
	if error != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": error.Error()})
		return
	}

	credential := claim.Credentials[repository.ProductUpdateCredential]

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
	itemIndexRencana := -1

	for i, item := range master.JenisRencanaPembangunan {
		if item == kategori {
			itemIndex = i
			break
		}
	}

	if itemIndex != -1 {
		master.JenisRencanaPembangunan[itemIndex] = newKategori
	}

	for i, item := range master.RencanaPembangunan {
		if item.Kategori == kategori {
			itemIndexRencana = i
			break
		}
	}

	if itemIndexRencana != -1 {
		master.RencanaPembangunan[itemIndexRencana].Kategori = newKategori
	}

	loc, _ := time.LoadLocation("Asia/Singapore")
	now := time.Now().In(loc).Format("02-01-2006")

	master.UpdatedAt = now + " " + time.Now().In(loc).Format("15:04:05")

	resultsSave := dm.DB.Save(&master)
	if resultsSave.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": resultsSave.Error})
		return
	}

	respone := struct {
		JenisRencana []string `json:"jenis_rencana,omitempty"`
		UpdatedAt    string   `json:"update,omitempty"`
	}{
		JenisRencana: master.JenisRencanaPembangunan,
		UpdatedAt:    master.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": respone})
}

func (dm *DataMasterControler) TambahJenisRencanaPembangunan(ctx *gin.Context) {
	kategori := ctx.Param("kategori")
	rencana := ctx.Param("rencana")
	kriteria := ctx.Param("kriteria")
	satuan := ctx.Param("satuan")
	id := ctx.Param("id")

	config, _ := initializers.LoadConfig()

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
				if item.Jenis == rencana {
					jenisExists = true
					ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "Data sudah ada"})
					return
				}
			}
		}
	}

	if !kategoriExists {
		jenis_rencana := []models.JenisRencana{}
		jenis_rencana = append(jenis_rencana, models.JenisRencana{Jenis: rencana,
			Kriteria: kriteria,
			Satuan:   satuan})

		master.RencanaPembangunan = append(master.RencanaPembangunan, models.Rencana{Kategori: kategori, JenisRencana: jenis_rencana})
	}

	if !jenisExists && kategoriExists {
		master.RencanaPembangunan[itemIndex].JenisRencana = append(master.RencanaPembangunan[itemIndex].JenisRencana, models.JenisRencana{Jenis: rencana,
			Kriteria: kriteria,
			Satuan:   satuan})
	}

	loc, _ := time.LoadLocation("Asia/Singapore")
	now := time.Now().In(loc).Format("02-01-2006")

	master.UpdatedAt = now + " " + time.Now().In(loc).Format("15:04:05")

	resultsSave := dm.DB.Save(&master)
	if resultsSave.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": resultsSave.Error})
		return
	}

	respone := struct {
		JenisRencana       []string         `json:"jenis_rencana,omitempty"`
		RencanaPembangunan []models.Rencana `json:"rencana_pembangunan,omitempty"`
		UpdatedAt          string           `json:"update,omitempty"`
	}{
		JenisRencana:       master.JenisRencanaPembangunan,
		RencanaPembangunan: master.RencanaPembangunan,
		UpdatedAt:          master.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": respone})
}

func (dm *DataMasterControler) HapusJenisRencanaPembangunan(ctx *gin.Context) {
	kategori := ctx.Param("kategori")
	rencana := ctx.Param("rencana")
	id := ctx.Param("id")

	config, _ := initializers.LoadConfig()

	accessUser := ctx.MustGet("accessUser").(string)

	claim, error := utils.ValidateToken(accessUser, config.AccessTokenPublicKey)
	if error != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": error.Error()})
		return
	}

	credential := claim.Credentials[repository.ProductDeleteCredential]

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

	for i := range master.RencanaPembangunan {
		if master.RencanaPembangunan[i].Kategori == kategori {
			for j, item := range master.RencanaPembangunan[i].JenisRencana {
				if item.Jenis == rencana {
					master.RencanaPembangunan[i].JenisRencana = append(master.RencanaPembangunan[i].JenisRencana[:j], master.RencanaPembangunan[i].JenisRencana[j+1:]...)
				}
			}
		}
	}

	loc, _ := time.LoadLocation("Asia/Singapore")
	now := time.Now().In(loc).Format("02-01-2006")

	master.UpdatedAt = now + " " + time.Now().In(loc).Format("15:04:05")

	resultsSave := dm.DB.Save(&master)
	if resultsSave.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": resultsSave.Error})
		return
	}

	respone := struct {
		JenisRencana       []string         `json:"jenis_rencana,omitempty"`
		RencanaPembangunan []models.Rencana `json:"rencana_pembangunan,omitempty"`
		UpdatedAt          string           `json:"update,omitempty"`
	}{
		JenisRencana:       master.JenisRencanaPembangunan,
		RencanaPembangunan: master.RencanaPembangunan,
		UpdatedAt:          master.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": respone})
}

func (dm *DataMasterControler) EditJenisRencanaPembangunan(ctx *gin.Context) {
	kategori := ctx.Param("kategori")
	rencana := ctx.Param("rencana")
	newRencana := ctx.Param("rencana_new")
	kriteria := ctx.Param("kriteria")
	satuan := ctx.Param("satuan")
	id := ctx.Param("id")

	config, _ := initializers.LoadConfig()

	accessUser := ctx.MustGet("accessUser").(string)

	claim, error := utils.ValidateToken(accessUser, config.AccessTokenPublicKey)
	if error != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": error.Error()})
		return
	}

	credential := claim.Credentials[repository.ProductUpdateCredential]

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

	itemIndexKategori := -1
	itemIndexRencana := -1

	for i := range master.RencanaPembangunan {
		if master.RencanaPembangunan[i].Kategori == kategori {
			itemIndexKategori = i
			for j, item := range master.RencanaPembangunan[i].JenisRencana {
				if item.Jenis == rencana {
					itemIndexRencana = j
					break
				}
			}
		}
	}

	if itemIndexKategori != -1 && itemIndexRencana != -1 {
		master.RencanaPembangunan[itemIndexKategori].JenisRencana[itemIndexRencana] = models.JenisRencana{Jenis: newRencana,
			Kriteria: kriteria,
			Satuan:   satuan}
	}

	loc, _ := time.LoadLocation("Asia/Singapore")
	now := time.Now().In(loc).Format("02-01-2006")

	master.UpdatedAt = now + " " + time.Now().In(loc).Format("15:04:05")

	resultsSave := dm.DB.Save(&master)
	if resultsSave.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": resultsSave.Error})
		return
	}

	respone := struct {
		RencanaPembangunan []models.Rencana `json:"rencana_pembangunan,omitempty"`
		UpdatedAt          string           `json:"update,omitempty"`
	}{
		RencanaPembangunan: master.RencanaPembangunan,
		UpdatedAt:          master.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": respone})
}

func (dm *DataMasterControler) TambahKategoriPerlengkapan(ctx *gin.Context) {
	kategori := ctx.Param("kategori")
	id := ctx.Param("id")

	config, _ := initializers.LoadConfig()

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

	exist := contains(master.KategoriPerlengkapan, kategori)

	if exist {
		ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "Data sudah ada"})
		return
	}

	master.KategoriPerlengkapan = append(master.KategoriPerlengkapan, kategori)

	loc, _ := time.LoadLocation("Asia/Singapore")
	now := time.Now().In(loc).Format("02-01-2006")

	master.UpdatedAt = now + " " + time.Now().In(loc).Format("15:04:05")

	resultsSave := dm.DB.Save(&master)
	if resultsSave.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": resultsSave.Error})
		return
	}

	respone := struct {
		KategoriPerlengkapan []string `json:"kategori_perlengkapan,omitempty"`
		UpdatedAt            string   `json:"update,omitempty"`
	}{
		KategoriPerlengkapan: master.KategoriPerlengkapan,
		UpdatedAt:            master.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": respone})
}

func (dm *DataMasterControler) HapusKategoriPerlengkapan(ctx *gin.Context) {
	kategori := ctx.Param("kategori")
	id := ctx.Param("id")

	config, _ := initializers.LoadConfig()

	accessUser := ctx.MustGet("accessUser").(string)

	claim, error := utils.ValidateToken(accessUser, config.AccessTokenPublicKey)
	if error != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": error.Error()})
		return
	}

	credential := claim.Credentials[repository.ProductDeleteCredential]

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

	for i, item := range master.KategoriPerlengkapan {
		if item == kategori {
			master.KategoriPerlengkapan = append(master.KategoriPerlengkapan[:i], master.KategoriPerlengkapan[i+1:]...)
			break
		}
	}

	for i, item := range master.PerlengkapanLaluLintas {
		if item.Kategori == kategori {
			master.PerlengkapanLaluLintas = append(master.PerlengkapanLaluLintas[:i], master.PerlengkapanLaluLintas[i+1:]...)
			break
		}
	}

	loc, _ := time.LoadLocation("Asia/Singapore")
	now := time.Now().In(loc).Format("02-01-2006")

	master.UpdatedAt = now + " " + time.Now().In(loc).Format("15:04:05")

	resultsSave := dm.DB.Save(&master)
	if resultsSave.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": resultsSave.Error})
		return
	}

	respone := struct {
		KategoriPerlengkapan []string `json:"kategori_perlengkapan,omitempty"`
		UpdatedAt            string   `json:"update,omitempty"`
	}{
		KategoriPerlengkapan: master.KategoriPerlengkapan,
		UpdatedAt:            master.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": respone})
}

func (dm *DataMasterControler) EditKategoriPerlengkapan(ctx *gin.Context) {
	kategori := ctx.Param("kategori")
	newKategori := ctx.Param("new_kategori")
	id := ctx.Param("id")

	config, _ := initializers.LoadConfig()

	accessUser := ctx.MustGet("accessUser").(string)

	claim, error := utils.ValidateToken(accessUser, config.AccessTokenPublicKey)
	if error != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": error.Error()})
		return
	}

	credential := claim.Credentials[repository.ProductUpdateCredential]

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
	itemIndexKategori := -1

	for i, item := range master.KategoriPerlengkapan {
		if item == kategori {
			itemIndex = i
			break
		}
	}

	if itemIndex != -1 {
		master.KategoriPerlengkapan[itemIndex] = newKategori
	}

	for i, item := range master.PerlengkapanLaluLintas {
		if item.Kategori == kategori {
			itemIndexKategori = i
			break
		}
	}

	if itemIndexKategori != -1 {
		master.PerlengkapanLaluLintas[itemIndexKategori].Kategori = newKategori
	}

	loc, _ := time.LoadLocation("Asia/Singapore")
	now := time.Now().In(loc).Format("02-01-2006")

	master.UpdatedAt = now + " " + time.Now().In(loc).Format("15:04:05")

	resultsSave := dm.DB.Save(&master)
	if resultsSave.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": resultsSave.Error})
		return
	}

	respone := struct {
		KategoriPerlengkapan []string `json:"kategori_perlengkapan,omitempty"`
		UpdatedAt            string   `json:"update,omitempty"`
	}{
		KategoriPerlengkapan: master.KategoriPerlengkapan,
		UpdatedAt:            master.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": respone})
}

func (dm *DataMasterControler) TambahPerlengkapan(ctx *gin.Context) {
	kategori := ctx.Param("kategori")
	perlengkapan := ctx.Param("perlengkapan")
	id := ctx.Param("id")

	config, _ := initializers.LoadConfig()

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
	perlengkapanExist := false
	itemIndex := 0

	for i := range master.PerlengkapanLaluLintas {
		if master.PerlengkapanLaluLintas[i].Kategori == kategori {
			kategoriExists = true
			itemIndex = i
			for _, item := range master.PerlengkapanLaluLintas[i].Perlengkapan {
				if item.JenisPerlengkapan == perlengkapan {
					perlengkapanExist = true
					ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "Data sudah ada"})
					return
				}
			}
		}
	}

	file, err := ctx.FormFile("perlengkapan")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	uploadedFile, err := file.Open()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer uploadedFile.Close()

	data, err := io.ReadAll(uploadedFile)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !kategoriExists {
		perlengkapan := models.PerlengkapanItem{
			JenisPerlengkapan:  perlengkapan,
			GambarPerlengkapan: data,
		}
		jenis := models.JenisPerlengkapan{
			Kategori:     kategori,
			Perlengkapan: []models.PerlengkapanItem{perlengkapan},
		}
		master.PerlengkapanLaluLintas = append(master.PerlengkapanLaluLintas, jenis)
	}

	if !perlengkapanExist && kategoriExists {
		perlengkapan := models.PerlengkapanItem{
			JenisPerlengkapan:  perlengkapan,
			GambarPerlengkapan: data,
		}
		master.PerlengkapanLaluLintas[itemIndex].Perlengkapan = append(master.PerlengkapanLaluLintas[itemIndex].Perlengkapan, perlengkapan)
	}

	loc, _ := time.LoadLocation("Asia/Singapore")
	now := time.Now().In(loc).Format("02-01-2006")

	master.UpdatedAt = now + " " + time.Now().In(loc).Format("15:04:05")

	resultsSave := dm.DB.Save(&master)
	if resultsSave.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": resultsSave.Error})
		return
	}

	respone := struct {
		PerlengkapanLaluLintas []models.JenisPerlengkapan `json:"perlengkapan,omitempty"`
		UpdatedAt              string                     `json:"update,omitempty"`
	}{
		PerlengkapanLaluLintas: master.PerlengkapanLaluLintas,
		UpdatedAt:              master.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": respone})
}

func (dm *DataMasterControler) HapuspPerlengkapan(ctx *gin.Context) {
	kategori := ctx.Param("kategori")
	perlengkapan := ctx.Param("perlengkapan")
	id := ctx.Param("id")

	config, _ := initializers.LoadConfig()

	accessUser := ctx.MustGet("accessUser").(string)

	claim, error := utils.ValidateToken(accessUser, config.AccessTokenPublicKey)
	if error != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": error.Error()})
		return
	}

	credential := claim.Credentials[repository.ProductDeleteCredential]

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

	for i := range master.PerlengkapanLaluLintas {
		if master.PerlengkapanLaluLintas[i].Kategori == kategori {
			for j, item := range master.PerlengkapanLaluLintas[i].Perlengkapan {
				if item.JenisPerlengkapan == perlengkapan {
					master.PerlengkapanLaluLintas[i].Perlengkapan = append(master.PerlengkapanLaluLintas[i].Perlengkapan[:j], master.PerlengkapanLaluLintas[i].Perlengkapan[j+1:]...)
				}
			}
		}
	}

	loc, _ := time.LoadLocation("Asia/Singapore")
	now := time.Now().In(loc).Format("02-01-2006")

	master.UpdatedAt = now + " " + time.Now().In(loc).Format("15:04:05")

	resultsSave := dm.DB.Save(&master)
	if resultsSave.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": resultsSave.Error})
		return
	}

	respone := struct {
		PerlengkapanLaluLintas []models.JenisPerlengkapan `json:"perlengkapan,omitempty"`
		UpdatedAt              string                     `json:"update,omitempty"`
	}{
		PerlengkapanLaluLintas: master.PerlengkapanLaluLintas,
		UpdatedAt:              master.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": respone})
}

func (dm *DataMasterControler) EditPerlengkapan(ctx *gin.Context) {
	kategori := ctx.Param("kategori")
	perlengkapan := ctx.Param("perlengkapan")
	newPerlengkapan := ctx.Param("perlengkapan_new")
	id := ctx.Param("id")

	config, _ := initializers.LoadConfig()

	accessUser := ctx.MustGet("accessUser").(string)

	claim, error := utils.ValidateToken(accessUser, config.AccessTokenPublicKey)
	if error != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": error.Error()})
		return
	}

	credential := claim.Credentials[repository.ProductUpdateCredential]

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

	itemIndexKategori := -1
	itemIndexPerlengkapan := -1

	for i := range master.PerlengkapanLaluLintas {
		if master.PerlengkapanLaluLintas[i].Kategori == kategori {
			itemIndexKategori = i
			for j, item := range master.PerlengkapanLaluLintas[i].Perlengkapan {
				if item.JenisPerlengkapan == perlengkapan {
					itemIndexPerlengkapan = j
					break
				}
			}
		}
	}

	file, _ := ctx.FormFile("perlengkapan")

	if itemIndexKategori != -1 && itemIndexPerlengkapan != -1 {
		master.PerlengkapanLaluLintas[itemIndexKategori].Perlengkapan[itemIndexPerlengkapan].JenisPerlengkapan = newPerlengkapan
		if file != nil {
			uploadedFile, err := file.Open()
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			defer uploadedFile.Close()

			data, err := io.ReadAll(uploadedFile)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			master.PerlengkapanLaluLintas[itemIndexKategori].Perlengkapan[itemIndexPerlengkapan].GambarPerlengkapan = data
		}

	}

	loc, _ := time.LoadLocation("Asia/Singapore")
	now := time.Now().In(loc).Format("02-01-2006")

	master.UpdatedAt = now + " " + time.Now().In(loc).Format("15:04:05")

	resultsSave := dm.DB.Save(&master)
	if resultsSave.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": resultsSave.Error})
		return
	}

	respone := struct {
		PerlengkapanLaluLintas []models.JenisPerlengkapan `json:"perlengkapan,omitempty"`
		UpdatedAt              string                     `json:"update,omitempty"`
	}{
		PerlengkapanLaluLintas: master.PerlengkapanLaluLintas,
		UpdatedAt:              master.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": respone})
}

func (dm *DataMasterControler) TambahPersyaratanAndalalin(ctx *gin.Context) {
	var payload *models.PersyaratanAndalalinInput
	id := ctx.Param("id")

	config, _ := initializers.LoadConfig()

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

	var master models.DataMaster

	resultsData := dm.DB.Where("id_data_master", id).First(&master)

	if resultsData.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": resultsData.Error})
		return
	}

	persyaratanExist := false

	for i := range master.Persyaratan.PersyaratanAndalalin {
		if master.Persyaratan.PersyaratanAndalalin[i].Persyaratan == payload.Persyaratan {
			persyaratanExist = true
			ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "Data sudah ada"})
			return
		}
	}

	if !persyaratanExist {
		persyaratan := models.PersyaratanAndalalinInput{
			Bangkitan:             payload.Bangkitan,
			Persyaratan:           payload.Persyaratan,
			KeteranganPersyaratan: payload.KeteranganPersyaratan,
		}
		master.Persyaratan.PersyaratanAndalalin = append(master.Persyaratan.PersyaratanAndalalin, persyaratan)
	}

	loc, _ := time.LoadLocation("Asia/Singapore")
	now := time.Now().In(loc).Format("02-01-2006")

	master.UpdatedAt = now + " " + time.Now().In(loc).Format("15:04:05")

	resultsSave := dm.DB.Save(&master)
	if resultsSave.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": resultsSave.Error})
		return
	}

	respone := struct {
		Persyaratan models.Persyaratan `json:"persyaratan,omitempty"`
		UpdatedAt   string             `json:"update,omitempty"`
	}{
		Persyaratan: master.Persyaratan,
		UpdatedAt:   master.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": respone})
}

func (dm *DataMasterControler) HapusPersyaratanAndalalin(ctx *gin.Context) {
	persyaratan := ctx.Param("persyaratan")
	id := ctx.Param("id")

	config, _ := initializers.LoadConfig()

	accessUser := ctx.MustGet("accessUser").(string)

	claim, error := utils.ValidateToken(accessUser, config.AccessTokenPublicKey)
	if error != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": error.Error()})
		return
	}

	credential := claim.Credentials[repository.ProductDeleteCredential]

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

	for i := range master.Persyaratan.PersyaratanAndalalin {
		if master.Persyaratan.PersyaratanAndalalin[i].Persyaratan == persyaratan {
			master.Persyaratan.PersyaratanAndalalin = append(master.Persyaratan.PersyaratanAndalalin[:i], master.Persyaratan.PersyaratanAndalalin[i+1:]...)
			break
		}
	}

	var andalalin []models.Andalalin

	results := dm.DB.Find(&andalalin)

	if results.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": results.Error})
		return
	} else {
		dataFile := []file{}
		for i, permohonan := range andalalin {
			for j, tambahan := range permohonan.Persyaratan {
				if tambahan.Persyaratan == persyaratan {
					oldSubstr := "/"
					newSubstr := "-"

					result := strings.Replace(permohonan.Kode, oldSubstr, newSubstr, -1)
					fileName := result + ".pdf"

					dataFile = append(dataFile, file{Name: fileName, File: tambahan.Berkas})
					andalalin[i].Persyaratan = append(andalalin[i].Persyaratan[:j], andalalin[i].Persyaratan[j+1:]...)
					break
				}
			}
		}

		zipFile := persyaratan + ".zip"
		error = compressFiles(zipFile, dataFile)
		if error != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": error})
			return
		}

		zipData, errorZip := os.ReadFile(zipFile)
		if errorZip != nil {
			ctx.JSON(http.StatusNoContent, gin.H{"status": "error", "message": errorZip})
			return
		}

		base64ZipData := base64.StdEncoding.EncodeToString(zipData)

		dm.DB.Save(&andalalin)

		loc, _ := time.LoadLocation("Asia/Singapore")
		now := time.Now().In(loc).Format("02-01-2006")

		master.UpdatedAt = now + " " + time.Now().In(loc).Format("15:04:05")

		resultsSave := dm.DB.Save(&master)
		if resultsSave.Error != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": resultsSave.Error})
			return
		}

		respone := struct {
			Persyaratan models.Persyaratan `json:"persyaratan,omitempty"`
			UpdatedAt   string             `json:"update,omitempty"`
		}{
			Persyaratan: master.Persyaratan,
			UpdatedAt:   master.UpdatedAt,
		}

		ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": respone, "file": base64ZipData})
	}
}

func (dm *DataMasterControler) EditPersyaratanAndalalin(ctx *gin.Context) {
	var payload *models.PersyaratanAndalalinInput
	id := ctx.Param("id")
	syarat := ctx.Param("persyaratan")

	config, _ := initializers.LoadConfig()

	accessUser := ctx.MustGet("accessUser").(string)

	claim, error := utils.ValidateToken(accessUser, config.AccessTokenPublicKey)
	if error != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": error.Error()})
		return
	}

	credential := claim.Credentials[repository.ProductUpdateCredential]

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

	var master models.DataMaster

	resultsData := dm.DB.Where("id_data_master", id).First(&master)

	if resultsData.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": resultsData.Error})
		return
	}

	itemIndex := -1

	for i := range master.Persyaratan.PersyaratanAndalalin {
		if master.Persyaratan.PersyaratanAndalalin[i].Persyaratan == syarat {
			itemIndex = i
			break
		}
	}

	if itemIndex != -1 {
		if master.Persyaratan.PersyaratanAndalalin[itemIndex].Persyaratan != payload.Persyaratan {
			master.Persyaratan.PersyaratanAndalalin[itemIndex].Persyaratan = payload.Persyaratan
		}

		if master.Persyaratan.PersyaratanAndalalin[itemIndex].KeteranganPersyaratan != payload.KeteranganPersyaratan {
			master.Persyaratan.PersyaratanAndalalin[itemIndex].KeteranganPersyaratan = payload.KeteranganPersyaratan
		}

	}

	loc, _ := time.LoadLocation("Asia/Singapore")
	now := time.Now().In(loc).Format("02-01-2006")

	master.UpdatedAt = now + " " + time.Now().In(loc).Format("15:04:05")

	resultsSave := dm.DB.Save(&master)
	if resultsSave.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": resultsSave.Error})
		return
	}

	respone := struct {
		Persyaratan models.Persyaratan `json:"persyaratan,omitempty"`
		UpdatedAt   string             `json:"update,omitempty"`
	}{
		Persyaratan: master.Persyaratan,
		UpdatedAt:   master.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": respone})
}

func (dm *DataMasterControler) TambahPersyaratanPerlalin(ctx *gin.Context) {
	var payload *models.PersyaratanPerlalinInput
	id := ctx.Param("id")

	config, _ := initializers.LoadConfig()

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

	var master models.DataMaster

	resultsData := dm.DB.Where("id_data_master", id).First(&master)

	if resultsData.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": resultsData.Error})
		return
	}

	persyaratanExist := false

	for i := range master.Persyaratan.PersyaratanPerlalin {
		if master.Persyaratan.PersyaratanPerlalin[i].Persyaratan == payload.Persyaratan {
			persyaratanExist = true
			ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "Data sudah ada"})
			return
		}
	}

	if !persyaratanExist {
		persyaratan := models.PersyaratanPerlalinInput{
			Persyaratan:           payload.Persyaratan,
			KeteranganPersyaratan: payload.KeteranganPersyaratan,
		}
		master.Persyaratan.PersyaratanPerlalin = append(master.Persyaratan.PersyaratanPerlalin, persyaratan)
	}

	loc, _ := time.LoadLocation("Asia/Singapore")
	now := time.Now().In(loc).Format("02-01-2006")

	master.UpdatedAt = now + " " + time.Now().In(loc).Format("15:04:05")

	resultsSave := dm.DB.Save(&master)
	if resultsSave.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": resultsSave.Error})
		return
	}

	respone := struct {
		Persyaratan models.Persyaratan `json:"persyaratan,omitempty"`
		UpdatedAt   string             `json:"update,omitempty"`
	}{
		Persyaratan: master.Persyaratan,
		UpdatedAt:   master.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": respone})
}

func (dm *DataMasterControler) HapusPersyaratanPerlalin(ctx *gin.Context) {
	persyaratan := ctx.Param("persyaratan")
	id := ctx.Param("id")

	config, _ := initializers.LoadConfig()

	accessUser := ctx.MustGet("accessUser").(string)

	claim, error := utils.ValidateToken(accessUser, config.AccessTokenPublicKey)
	if error != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": error.Error()})
		return
	}

	credential := claim.Credentials[repository.ProductDeleteCredential]

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

	for i := range master.Persyaratan.PersyaratanPerlalin {
		if master.Persyaratan.PersyaratanPerlalin[i].Persyaratan == persyaratan {
			master.Persyaratan.PersyaratanPerlalin = append(master.Persyaratan.PersyaratanPerlalin[:i], master.Persyaratan.PersyaratanPerlalin[i+1:]...)
			break
		}
	}

	var perlalin []models.Perlalin

	results := dm.DB.Find(&perlalin)

	if results.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": results.Error})
		return
	} else {
		dataFile := []file{}
		for i, permohonan := range perlalin {
			for j, tambahan := range permohonan.Persyaratan {
				if tambahan.Persyaratan == persyaratan {
					oldSubstr := "/"
					newSubstr := "-"

					result := strings.Replace(permohonan.Kode, oldSubstr, newSubstr, -1)
					fileName := result + ".pdf"

					dataFile = append(dataFile, file{Name: fileName, File: tambahan.Berkas})
					perlalin[i].Persyaratan = append(perlalin[i].Persyaratan[:j], perlalin[i].Persyaratan[j+1:]...)
					break
				}
			}
		}

		zipFile := persyaratan + ".zip"
		error = compressFiles(zipFile, dataFile)
		if error != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": error})
			return
		}

		zipData, errorZip := os.ReadFile(zipFile)
		if errorZip != nil {
			ctx.JSON(http.StatusNoContent, gin.H{"status": "error", "message": errorZip})
			return
		}

		base64ZipData := base64.StdEncoding.EncodeToString(zipData)

		dm.DB.Save(&perlalin)

		loc, _ := time.LoadLocation("Asia/Singapore")
		now := time.Now().In(loc).Format("02-01-2006")

		master.UpdatedAt = now + " " + time.Now().In(loc).Format("15:04:05")

		resultsSave := dm.DB.Save(&master)
		if resultsSave.Error != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": resultsSave.Error})
			return
		}

		respone := struct {
			Persyaratan models.Persyaratan `json:"persyaratan,omitempty"`
			UpdatedAt   string             `json:"update,omitempty"`
		}{
			Persyaratan: master.Persyaratan,
			UpdatedAt:   master.UpdatedAt,
		}

		ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": respone, "file": base64ZipData})
	}
}

func (dm *DataMasterControler) EditPersyaratanPerlalin(ctx *gin.Context) {
	var payload *models.PersyaratanPerlalinInput
	id := ctx.Param("id")
	syarat := ctx.Param("persyaratan")

	config, _ := initializers.LoadConfig()

	accessUser := ctx.MustGet("accessUser").(string)

	claim, error := utils.ValidateToken(accessUser, config.AccessTokenPublicKey)
	if error != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": error.Error()})
		return
	}

	credential := claim.Credentials[repository.ProductUpdateCredential]

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

	var master models.DataMaster

	resultsData := dm.DB.Where("id_data_master", id).First(&master)

	if resultsData.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": resultsData.Error})
		return
	}

	itemIndex := -1

	for i := range master.Persyaratan.PersyaratanPerlalin {
		if master.Persyaratan.PersyaratanPerlalin[i].Persyaratan == syarat {
			itemIndex = i
			break
		}
	}

	if itemIndex != -1 {
		if master.Persyaratan.PersyaratanPerlalin[itemIndex].Persyaratan != payload.Persyaratan {
			master.Persyaratan.PersyaratanPerlalin[itemIndex].Persyaratan = payload.Persyaratan
		}

		if master.Persyaratan.PersyaratanPerlalin[itemIndex].KeteranganPersyaratan != payload.KeteranganPersyaratan {
			master.Persyaratan.PersyaratanPerlalin[itemIndex].KeteranganPersyaratan = payload.KeteranganPersyaratan
		}
	}

	loc, _ := time.LoadLocation("Asia/Singapore")
	now := time.Now().In(loc).Format("02-01-2006")

	master.UpdatedAt = now + " " + time.Now().In(loc).Format("15:04:05")

	resultsSave := dm.DB.Save(&master)
	if resultsSave.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": resultsSave.Error})
		return
	}

	respone := struct {
		Persyaratan models.Persyaratan `json:"persyaratan,omitempty"`
		UpdatedAt   string             `json:"update,omitempty"`
	}{
		Persyaratan: master.Persyaratan,
		UpdatedAt:   master.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": respone})
}

func (dm *DataMasterControler) TambahProvinsi(ctx *gin.Context) {
	provinsi := ctx.Param("provinsi")
	id := ctx.Param("id")

	config, _ := initializers.LoadConfig()

	accessUser := ctx.MustGet("accessUser").(string)

	claim, error := utils.ValidateToken(accessUser, config.AccessTokenPublicKey)
	if error != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": error.Error()})
		return
	}

	credential := claim.Credentials[repository.ProductDeleteCredential]

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

	exist := false
	dataId := utils.Encode(2)

	for _, item := range master.Provinsi {
		if item.Name == provinsi {
			exist = true
			break
		}
	}

	for _, item := range master.Provinsi {
		if item.Id == dataId {
			dataId = utils.Encode(2)
		} else {
			break
		}
	}

	if !exist {
		data := models.Provinsi{
			Id:   dataId,
			Name: provinsi,
		}
		master.Provinsi = append(master.Provinsi, data)
	}

	loc, _ := time.LoadLocation("Asia/Singapore")
	now := time.Now().In(loc).Format("02-01-2006")

	master.UpdatedAt = now + " " + time.Now().In(loc).Format("15:04:05")

	resultsSave := dm.DB.Save(&master)
	if resultsSave.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": resultsSave.Error})
		return
	}

	respone := struct {
		Provinsi  []models.Provinsi `json:"provinsi,omitempty"`
		UpdatedAt string            `json:"update,omitempty"`
	}{
		UpdatedAt: master.UpdatedAt,
		Provinsi:  master.Provinsi,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": respone})
}

func (dm *DataMasterControler) HapusProvinsi(ctx *gin.Context) {
	provinsi := ctx.Param("provinsi")
	id := ctx.Param("id")

	config, _ := initializers.LoadConfig()

	accessUser := ctx.MustGet("accessUser").(string)

	claim, error := utils.ValidateToken(accessUser, config.AccessTokenPublicKey)
	if error != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": error.Error()})
		return
	}

	credential := claim.Credentials[repository.ProductDeleteCredential]

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

	for i, item := range master.Provinsi {
		if item.Name == provinsi {
			master.Provinsi = append(master.Provinsi[:i], master.Provinsi[i+1:]...)
			break
		}
	}

	loc, _ := time.LoadLocation("Asia/Singapore")
	now := time.Now().In(loc).Format("02-01-2006")

	master.UpdatedAt = now + " " + time.Now().In(loc).Format("15:04:05")

	resultsSave := dm.DB.Save(&master)
	if resultsSave.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": resultsSave.Error})
		return
	}

	respone := struct {
		Provinsi  []models.Provinsi `json:"provinsi,omitempty"`
		UpdatedAt string            `json:"update,omitempty"`
	}{
		UpdatedAt: master.UpdatedAt,
		Provinsi:  master.Provinsi,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": respone})
}

func (dm *DataMasterControler) EditProvinsi(ctx *gin.Context) {
	provinsi := ctx.Param("provinsi")
	newProvinsi := ctx.Param("new_provinsi")
	id := ctx.Param("id")

	config, _ := initializers.LoadConfig()

	accessUser := ctx.MustGet("accessUser").(string)

	claim, error := utils.ValidateToken(accessUser, config.AccessTokenPublicKey)
	if error != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": error.Error()})
		return
	}

	credential := claim.Credentials[repository.ProductUpdateCredential]

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

	for i, item := range master.Provinsi {
		if item.Name == provinsi {
			itemIndex = i
			break
		}
	}

	if itemIndex != -1 {
		master.Provinsi[itemIndex].Name = newProvinsi
	}

	loc, _ := time.LoadLocation("Asia/Singapore")
	now := time.Now().In(loc).Format("02-01-2006")

	master.UpdatedAt = now + " " + time.Now().In(loc).Format("15:04:05")

	resultsSave := dm.DB.Save(&master)
	if resultsSave.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": resultsSave.Error})
		return
	}

	respone := struct {
		Provinsi  []models.Provinsi `json:"provinsi,omitempty"`
		UpdatedAt string            `json:"update,omitempty"`
	}{
		UpdatedAt: master.UpdatedAt,
		Provinsi:  master.Provinsi,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": respone})
}

func (dm *DataMasterControler) TambahKabupaten(ctx *gin.Context) {
	kabupaten := ctx.Param("kabupaten")
	provinsi := ctx.Param("provinsi")
	id := ctx.Param("id")

	config, _ := initializers.LoadConfig()

	accessUser := ctx.MustGet("accessUser").(string)

	claim, error := utils.ValidateToken(accessUser, config.AccessTokenPublicKey)
	if error != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": error.Error()})
		return
	}

	credential := claim.Credentials[repository.ProductDeleteCredential]

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

	exist := false
	dataId := utils.Encode(2)
	var id_provinsi string

	for _, item := range master.Provinsi {
		if item.Name == provinsi {
			id_provinsi = item.Id
			break
		}
	}

	for _, item := range master.Kabupaten {
		if item.Name == kabupaten {
			exist = true
			break
		}
	}

	for _, item := range master.Provinsi {
		if item.Id == dataId {
			dataId = utils.Encode(2)
		} else {
			break
		}
	}

	if !exist {
		data := models.Kabupaten{
			Id:         id_provinsi + dataId,
			IdProvinsi: id_provinsi,
			Name:       kabupaten,
		}
		master.Kabupaten = append(master.Kabupaten, data)
	}

	loc, _ := time.LoadLocation("Asia/Singapore")
	now := time.Now().In(loc).Format("02-01-2006")

	master.UpdatedAt = now + " " + time.Now().In(loc).Format("15:04:05")

	resultsSave := dm.DB.Save(&master)
	if resultsSave.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": resultsSave.Error})
		return
	}

	respone := struct {
		Kabupaten []models.Kabupaten `json:"kabupaten,omitempty"`
		UpdatedAt string             `json:"update,omitempty"`
	}{
		UpdatedAt: master.UpdatedAt,
		Kabupaten: master.Kabupaten,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": respone})
}

func (dm *DataMasterControler) HapusKabupaten(ctx *gin.Context) {
	kabupaten := ctx.Param("kabupaten")
	id := ctx.Param("id")

	config, _ := initializers.LoadConfig()

	accessUser := ctx.MustGet("accessUser").(string)

	claim, error := utils.ValidateToken(accessUser, config.AccessTokenPublicKey)
	if error != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": error.Error()})
		return
	}

	credential := claim.Credentials[repository.ProductDeleteCredential]

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

	for i, item := range master.Kabupaten {
		if item.Name == kabupaten {
			master.Provinsi = append(master.Provinsi[:i], master.Provinsi[i+1:]...)
			break
		}
	}

	loc, _ := time.LoadLocation("Asia/Singapore")
	now := time.Now().In(loc).Format("02-01-2006")

	master.UpdatedAt = now + " " + time.Now().In(loc).Format("15:04:05")

	resultsSave := dm.DB.Save(&master)
	if resultsSave.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": resultsSave.Error})
		return
	}

	respone := struct {
		Kabupaten []models.Kabupaten `json:"kabupaten,omitempty"`
		UpdatedAt string             `json:"update,omitempty"`
	}{
		UpdatedAt: master.UpdatedAt,
		Kabupaten: master.Kabupaten,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": respone})
}

func (dm *DataMasterControler) EditKabupaten(ctx *gin.Context) {
	kabupaten := ctx.Param("kabupaten")
	newKabupaten := ctx.Param("new_kabupaten")
	id := ctx.Param("id")

	config, _ := initializers.LoadConfig()

	accessUser := ctx.MustGet("accessUser").(string)

	claim, error := utils.ValidateToken(accessUser, config.AccessTokenPublicKey)
	if error != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": error.Error()})
		return
	}

	credential := claim.Credentials[repository.ProductUpdateCredential]

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

	for i, item := range master.Kabupaten {
		if item.Name == kabupaten {
			itemIndex = i
			break
		}
	}

	if itemIndex != -1 {
		master.Kabupaten[itemIndex].Name = newKabupaten
	}

	loc, _ := time.LoadLocation("Asia/Singapore")
	now := time.Now().In(loc).Format("02-01-2006")

	master.UpdatedAt = now + " " + time.Now().In(loc).Format("15:04:05")

	resultsSave := dm.DB.Save(&master)
	if resultsSave.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": resultsSave.Error})
		return
	}

	respone := struct {
		Kabupaten []models.Kabupaten `json:"kabupaten,omitempty"`
		UpdatedAt string             `json:"update,omitempty"`
	}{
		UpdatedAt: master.UpdatedAt,
		Kabupaten: master.Kabupaten,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": respone})
}
