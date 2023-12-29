package controllers

import (
	"archive/zip"
	"encoding/base64"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"andalalin/initializers"
	"andalalin/models"
	"andalalin/repository"
	"andalalin/utils"

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
	var mutex sync.Mutex

	mutex.Lock()
	defer mutex.Unlock()

	bufferSize := 10
	resultChan := make(chan models.DataMaster, bufferSize)

	rows, err := dm.DB.Table("data_masters").Rows()
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Data error"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var result models.DataMaster
		if err := dm.DB.ScanRows(rows, &result); err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Data error"})
			return
		}

		resultChan <- result
	}

	close(resultChan)

	for result := range resultChan {
		respone := struct {
			IdDataMaster               uuid.UUID                        `json:"id_data_master,omitempty"`
			JenisProyek                []string                         `json:"jenis_proyek,omitempty"`
			Lokasi                     []string                         `json:"lokasi_pengambilan,omitempty"`
			KategoriRencanaPembangunan []string                         `json:"kategori_rencana,omitempty"`
			JenisRencanaPembangunan    []models.JenisRencanaPembangunan `json:"jenis_rencana,omitempty"`
			KategoriPerlengkapanUtama  []string                         `json:"kategori_utama,omitempty"`
			KategoriPerlengkapan       []models.KategoriPerlengkapan    `json:"kategori_perlengkapan,omitempty"`
			PerlengkapanLaluLintas     []models.JenisPerlengkapan       `json:"perlengkapan,omitempty"`
			Persyaratan                models.Persyaratan               `json:"persyaratan,omitempty"`
			Provinsi                   []models.Provinsi                `json:"provinsi,omitempty"`
			Kabupaten                  []models.Kabupaten               `json:"kabupaten,omitempty"`
			Kecamatan                  []models.Kecamatan               `json:"kecamatan,omitempty"`
			Kelurahan                  []models.Kelurahan               `json:"kelurahan,omitempty"`
			Jalan                      []models.Jalan                   `json:"jalan,omitempty"`
			UpdatedAt                  string                           `json:"update,omitempty"`
		}{
			IdDataMaster:               result.IdDataMaster,
			JenisProyek:                result.JenisProyek,
			Lokasi:                     result.LokasiPengambilan,
			KategoriRencanaPembangunan: result.KategoriRencanaPembangunan,
			JenisRencanaPembangunan:    result.JenisRencanaPembangunan,
			KategoriPerlengkapanUtama:  result.KategoriPerlengkapanUtama,
			KategoriPerlengkapan:       result.KategoriPerlengkapan,
			PerlengkapanLaluLintas:     result.PerlengkapanLaluLintas,
			Persyaratan:                result.Persyaratan,
			Provinsi:                   result.Provinsi,
			Kabupaten:                  result.Kabupaten,
			Kecamatan:                  result.Kecamatan,
			Kelurahan:                  result.Kelurahan,
			Jalan:                      result.Jalan,
			UpdatedAt:                  result.UpdatedAt,
		}
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": respone})
	}
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

func (dm *DataMasterControler) GetDataMasterByType(ctx *gin.Context) {
	tipe := ctx.Param("tipe")

	switch tipe {
	case "proyek":
		var mutex sync.Mutex

		mutex.Lock()
		defer mutex.Unlock()

		bufferSize := 10
		resultChan := make(chan models.DataMaster, bufferSize)

		rows, err := dm.DB.Table("data_masters").Select("id_data_master", "jenis_proyek").Rows()
		if err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Data error"})
			return
		}
		defer rows.Close()

		for rows.Next() {
			var result models.DataMaster
			if err := dm.DB.ScanRows(rows, &result); err != nil {
				ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Data error"})
				return
			}

			resultChan <- result
		}

		close(resultChan)

		for result := range resultChan {
			respone := struct {
				IdDataMaster uuid.UUID `json:"id_data_master,omitempty"`
				JenisProyek  []string  `json:"jenis_proyek,omitempty"`
			}{
				IdDataMaster: result.IdDataMaster,
				JenisProyek:  result.JenisProyek,
			}
			ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": respone})
		}
	case "wilayah":
		var mutex sync.Mutex

		mutex.Lock()
		defer mutex.Unlock()

		bufferSize := 10
		resultChan := make(chan models.DataMaster, bufferSize)

		rows, err := dm.DB.Table("data_masters").Select("id_data_master", "provinsi", "kabupaten", "kecamatan", "kelurahan").Rows()
		if err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Data error"})
			return
		}
		defer rows.Close()

		for rows.Next() {
			var result models.DataMaster
			if err := dm.DB.ScanRows(rows, &result); err != nil {
				ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Data error"})
				return
			}

			resultChan <- result
		}

		close(resultChan)

		for result := range resultChan {
			respone := struct {
				IdDataMaster uuid.UUID          `json:"id_data_master,omitempty"`
				Provinsi     []models.Provinsi  `json:"provinsi,omitempty"`
				Kabupaten    []models.Kabupaten `json:"kabupaten,omitempty"`
				Kecamatan    []models.Kecamatan `json:"kecamatan,omitempty"`
				Kelurahan    []models.Kelurahan `json:"kelurahan,omitempty"`
			}{
				IdDataMaster: result.IdDataMaster,
				Provinsi:     result.Provinsi,
				Kabupaten:    result.Kabupaten,
				Kecamatan:    result.Kecamatan,
				Kelurahan:    result.Kelurahan,
			}
			ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": respone})
		}
	case "jalan":
		var mutex sync.Mutex

		mutex.Lock()
		defer mutex.Unlock()

		bufferSize := 10
		resultChan := make(chan models.DataMaster, bufferSize)

		rows, err := dm.DB.Table("data_masters").Select("id_data_master", "jalan", "kabupaten", "kecamatan", "kelurahan").Rows()
		if err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Data error"})
			return
		}
		defer rows.Close()

		for rows.Next() {
			var result models.DataMaster
			if err := dm.DB.ScanRows(rows, &result); err != nil {
				ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Data error"})
				return
			}

			resultChan <- result
		}

		close(resultChan)

		for result := range resultChan {
			var id_kabupaten string

			for _, kabupaten := range result.Kabupaten {
				if kabupaten.Name == "Kota Banjarmasin" {
					id_kabupaten += kabupaten.Id
				}
			}

			kecamatan_filter := []models.Kecamatan{}

			for _, kecamatan := range result.Kecamatan {
				if kecamatan.IdKabupaten == id_kabupaten {
					kecamatan_filter = append(kecamatan_filter, models.Kecamatan{Id: kecamatan.Id, IdKabupaten: kecamatan.IdKabupaten, Name: kecamatan.Name})
				}
			}

			kelurahan_filter := []models.Kelurahan{}

			for _, kecamatan := range kecamatan_filter {
				for _, kelurahan := range result.Kelurahan {
					if kelurahan.IdKecamatan == kecamatan.Id {
						kelurahan_filter = append(kelurahan_filter, models.Kelurahan{Id: kelurahan.Id, IdKecamatan: kelurahan.IdKecamatan, Name: kelurahan.Name})
					}
				}
			}

			respone := struct {
				IdDataMaster uuid.UUID          `json:"id_data_master,omitempty"`
				Jalan        []models.Jalan     `json:"jalan,omitempty"`
				Kecamatan    []models.Kecamatan `json:"kecamatan,omitempty"`
				Kelurahan    []models.Kelurahan `json:"kelurahan,omitempty"`
			}{
				IdDataMaster: result.IdDataMaster,
				Jalan:        result.Jalan,
				Kecamatan:    kecamatan_filter,
				Kelurahan:    kelurahan_filter,
			}
			ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": respone})
		}
	case "pengambilan":
		var mutex sync.Mutex

		mutex.Lock()
		defer mutex.Unlock()

		bufferSize := 10
		resultChan := make(chan models.DataMaster, bufferSize)

		rows, err := dm.DB.Table("data_masters").Select("id_data_master", "lokasi_pengambilan").Rows()
		if err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Data error"})
			return
		}
		defer rows.Close()

		for rows.Next() {
			var result models.DataMaster
			if err := dm.DB.ScanRows(rows, &result); err != nil {
				ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Data error"})
				return
			}

			resultChan <- result
		}

		close(resultChan)

		for result := range resultChan {
			respone := struct {
				IdDataMaster uuid.UUID `json:"id_data_master,omitempty"`
				Lokasi       []string  `json:"lokasi_pengambilan,omitempty"`
			}{
				IdDataMaster: result.IdDataMaster,
				Lokasi:       result.LokasiPengambilan,
			}
			ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": respone})
		}
	case "kategori_pembangunan":
	case "jenis_pembangunan":
	case "kategori_utama":
	case "kategori_perlengkapan":
	case "jenis_perlengkapan":
	case "persyaratan_andalalin":
	case "persyaratam_perlalin":
	}
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
	id := ctx.Param("id")
	var payload *models.LokasiInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

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

	rows, err := dm.DB.Table("data_masters").Where("id_data_master", id).Select("lokasi_pengambilan", "updated_at").Rows()
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Data error"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		if err := dm.DB.ScanRows(rows, &master); err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Data error"})
			return
		}
	}

	exist := contains(master.LokasiPengambilan, payload.Lokasi)

	if exist {
		ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "Data sudah ada"})
		return
	}

	master.LokasiPengambilan = append(master.LokasiPengambilan, payload.Lokasi)

	loc, _ := time.LoadLocation("Asia/Singapore")
	now := time.Now().In(loc).Format("02-01-2006")

	master.UpdatedAt = now + " " + time.Now().In(loc).Format("15:04:05")

	resultsSave := dm.DB.Table("data_masters").Where("id_data_master", id).Select("lokasi_pengambilan", "updated_at").Updates(models.DataMaster{LokasiPengambilan: master.LokasiPengambilan, UpdatedAt: master.UpdatedAt})
	if resultsSave.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": resultsSave.Error})
		return
	}

	respone := struct {
		Lokasi []string `json:"lokasi_pengambilan,omitempty"`
	}{
		Lokasi: master.LokasiPengambilan,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": respone})
}

func (dm *DataMasterControler) HapusLokasi(ctx *gin.Context) {
	id := ctx.Param("id")

	var payload *models.LokasiInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

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

	rows, err := dm.DB.Table("data_masters").Where("id_data_master", id).Select("lokasi_pengambilan", "updated_at").Rows()
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Data error"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		if err := dm.DB.ScanRows(rows, &master); err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Data error"})
			return
		}
	}

	for i, item := range master.LokasiPengambilan {
		if item == payload.Lokasi {
			master.LokasiPengambilan = append(master.LokasiPengambilan[:i], master.LokasiPengambilan[i+1:]...)
			break
		}
	}

	loc, _ := time.LoadLocation("Asia/Singapore")
	now := time.Now().In(loc).Format("02-01-2006")

	master.UpdatedAt = now + " " + time.Now().In(loc).Format("15:04:05")

	resultsSave := dm.DB.Table("data_masters").Where("id_data_master", id).Select("lokasi_pengambilan", "updated_at").Updates(models.DataMaster{LokasiPengambilan: master.LokasiPengambilan, UpdatedAt: master.UpdatedAt})
	if resultsSave.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": resultsSave.Error})
		return
	}

	respone := struct {
		Lokasi []string `json:"lokasi_pengambilan,omitempty"`
	}{
		Lokasi: master.LokasiPengambilan,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": respone})
}

func (dm *DataMasterControler) EditLokasi(ctx *gin.Context) {
	id := ctx.Param("id")

	var payload *models.LokasiEdit

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

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

	rows, err := dm.DB.Table("data_masters").Where("id_data_master", id).Select("lokasi_pengambilan", "updated_at").Rows()
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Data error"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		if err := dm.DB.ScanRows(rows, &master); err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Data error"})
			return
		}
	}

	itemIndex := -1

	for i, item := range master.LokasiPengambilan {
		if item == payload.Lokasi {
			itemIndex = i
			break
		}
	}

	if itemIndex != -1 {
		master.LokasiPengambilan[itemIndex] = payload.LokasiEdit
	}

	loc, _ := time.LoadLocation("Asia/Singapore")
	now := time.Now().In(loc).Format("02-01-2006")

	master.UpdatedAt = now + " " + time.Now().In(loc).Format("15:04:05")

	resultsSave := dm.DB.Table("data_masters").Where("id_data_master", id).Select("lokasi_pengambilan", "updated_at").Updates(models.DataMaster{LokasiPengambilan: master.LokasiPengambilan, UpdatedAt: master.UpdatedAt})
	if resultsSave.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": resultsSave.Error})
		return
	}

	respone := struct {
		Lokasi []string `json:"lokasi_pengambilan,omitempty"`
	}{
		Lokasi: master.LokasiPengambilan,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": respone})
}

func (dm *DataMasterControler) TambahKategori(ctx *gin.Context) {
	id := ctx.Param("id")

	var payload *models.KategoriRencanaInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

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

	exist := contains(master.KategoriRencanaPembangunan, payload.Kategori)

	if exist {
		ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "Data sudah ada"})
		return
	}

	master.KategoriRencanaPembangunan = append(master.KategoriRencanaPembangunan, payload.Kategori)

	loc, _ := time.LoadLocation("Asia/Singapore")
	now := time.Now().In(loc).Format("02-01-2006")

	master.UpdatedAt = now + " " + time.Now().In(loc).Format("15:04:05")

	resultsSave := dm.DB.Save(&master)
	if resultsSave.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": resultsSave.Error})
		return
	}

	respone := struct {
		KategoriRencanaPembangunan []string `json:"kategori_rencana,omitempty"`
		UpdatedAt                  string   `json:"update,omitempty"`
	}{
		KategoriRencanaPembangunan: master.KategoriRencanaPembangunan,
		UpdatedAt:                  master.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": respone})
}

func (dm *DataMasterControler) HapusKategori(ctx *gin.Context) {
	id := ctx.Param("id")

	var payload *models.KategoriRencanaInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

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

	for i, item := range master.KategoriRencanaPembangunan {
		if item == payload.Kategori {
			master.KategoriRencanaPembangunan = append(master.KategoriRencanaPembangunan[:i], master.KategoriRencanaPembangunan[i+1:]...)
			break
		}
	}

	for i, item := range master.JenisRencanaPembangunan {
		if item.Kategori == payload.Kategori {
			master.JenisRencanaPembangunan = append(master.JenisRencanaPembangunan[:i], master.JenisRencanaPembangunan[i+1:]...)
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
		KategoriRencanaPembangunan []string `json:"kategori_rencana,omitempty"`
		UpdatedAt                  string   `json:"update,omitempty"`
	}{
		KategoriRencanaPembangunan: master.KategoriRencanaPembangunan,
		UpdatedAt:                  master.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": respone})
}

func (dm *DataMasterControler) EditKategori(ctx *gin.Context) {
	id := ctx.Param("id")

	var payload *models.KategoriRencanaEdit

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

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

	for i, item := range master.KategoriRencanaPembangunan {
		if item == payload.Kategori {
			itemIndex = i
			break
		}
	}

	if itemIndex != -1 {
		master.KategoriRencanaPembangunan[itemIndex] = payload.KategoriEdit
	}

	for i, item := range master.JenisRencanaPembangunan {
		if item.Kategori == payload.Kategori {
			itemIndexRencana = i
			break
		}
	}

	if itemIndexRencana != -1 {
		master.JenisRencanaPembangunan[itemIndexRencana].Kategori = payload.KategoriEdit
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
		KategoriRencanaPembangunan []string `json:"kategori_rencana,omitempty"`
		UpdatedAt                  string   `json:"update,omitempty"`
	}{
		KategoriRencanaPembangunan: master.KategoriRencanaPembangunan,
		UpdatedAt:                  master.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": respone})
}

func (dm *DataMasterControler) TambahJenisRencanaPembangunan(ctx *gin.Context) {
	id := ctx.Param("id")

	var payload *models.JenisRencanaPembangunanInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

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

	for i := range master.JenisRencanaPembangunan {
		if master.JenisRencanaPembangunan[i].Kategori == payload.Kategori {
			kategoriExists = true
			itemIndex = i
			for _, item := range master.JenisRencanaPembangunan[i].JenisRencana {
				if item.Jenis == payload.Jenis {
					jenisExists = true
					ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "Data sudah ada"})
					return
				}
			}
		}
	}

	if !kategoriExists {
		jenis_rencana := []models.JenisRencana{}
		jenis_rencana = append(jenis_rencana, models.JenisRencana{Jenis: payload.Jenis,
			Kriteria: payload.Kriteria,
			Satuan:   payload.Satuan, Terbilang: payload.Terbilang})

		master.JenisRencanaPembangunan = append(master.JenisRencanaPembangunan, models.JenisRencanaPembangunan{Kategori: payload.Kategori, JenisRencana: jenis_rencana})
	}

	if !jenisExists && kategoriExists {
		master.JenisRencanaPembangunan[itemIndex].JenisRencana = append(master.JenisRencanaPembangunan[itemIndex].JenisRencana, models.JenisRencana{Jenis: payload.Jenis,
			Kriteria: payload.Kriteria,
			Satuan:   payload.Satuan, Terbilang: payload.Terbilang})
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
		KategoriRencanaPembangunan []string                         `json:"kategori_rencana,omitempty"`
		JenisRencanaPembangunan    []models.JenisRencanaPembangunan `json:"jenis_rencana,omitempty"`
		UpdatedAt                  string                           `json:"update,omitempty"`
	}{
		KategoriRencanaPembangunan: master.KategoriRencanaPembangunan,
		JenisRencanaPembangunan:    master.JenisRencanaPembangunan,
		UpdatedAt:                  master.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": respone})
}

func (dm *DataMasterControler) HapusJenisRencanaPembangunan(ctx *gin.Context) {
	id := ctx.Param("id")

	var payload *models.JenisRencanaPembangunanHapus

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

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

	for i := range master.JenisRencanaPembangunan {
		if master.JenisRencanaPembangunan[i].Kategori == payload.Kategori {
			for j, item := range master.JenisRencanaPembangunan[i].JenisRencana {
				if item.Jenis == payload.Jenis {
					master.JenisRencanaPembangunan[i].JenisRencana = append(master.JenisRencanaPembangunan[i].JenisRencana[:j], master.JenisRencanaPembangunan[i].JenisRencana[j+1:]...)
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
		KategoriRencanaPembangunan []string                         `json:"kategori_rencana,omitempty"`
		JenisRencanaPembangunan    []models.JenisRencanaPembangunan `json:"jenis_rencana,omitempty"`
		UpdatedAt                  string                           `json:"update,omitempty"`
	}{
		KategoriRencanaPembangunan: master.KategoriRencanaPembangunan,
		JenisRencanaPembangunan:    master.JenisRencanaPembangunan,
		UpdatedAt:                  master.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": respone})
}

func (dm *DataMasterControler) EditJenisRencanaPembangunan(ctx *gin.Context) {
	id := ctx.Param("id")

	var payload *models.JenisRencanaPembangunanEdit

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

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

	for i := range master.JenisRencanaPembangunan {
		if master.JenisRencanaPembangunan[i].Kategori == payload.Kategori {
			itemIndexKategori = i
			for j, item := range master.JenisRencanaPembangunan[i].JenisRencana {
				if item.Jenis == payload.Jenis {
					itemIndexRencana = j
					break
				}
			}
		}
	}

	if itemIndexKategori != -1 && itemIndexRencana != -1 {
		master.JenisRencanaPembangunan[itemIndexKategori].JenisRencana[itemIndexRencana] = models.JenisRencana{Jenis: payload.JenisEdit,
			Kriteria: payload.Kriteria,
			Satuan:   payload.Satuan, Terbilang: payload.Terbilang}
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
		JenisRencanaPembangunan []models.JenisRencanaPembangunan `json:"jenis_rencana,omitempty"`
		UpdatedAt               string                           `json:"update,omitempty"`
	}{
		JenisRencanaPembangunan: master.JenisRencanaPembangunan,
		UpdatedAt:               master.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": respone})
}

func (dm *DataMasterControler) TambahKategoriUtamaPerlengkapan(ctx *gin.Context) {
	var payload *models.KategoriPerlengkapanUtamaInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

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

	exist := contains(master.KategoriPerlengkapanUtama, payload.Kategori)

	if exist {
		ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "Data sudah ada"})
		return
	}

	master.KategoriPerlengkapanUtama = append(master.KategoriPerlengkapanUtama, payload.Kategori)

	loc, _ := time.LoadLocation("Asia/Singapore")
	now := time.Now().In(loc).Format("02-01-2006")

	master.UpdatedAt = now + " " + time.Now().In(loc).Format("15:04:05")

	resultsSave := dm.DB.Save(&master)
	if resultsSave.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": resultsSave.Error})
		return
	}

	respone := struct {
		KategoriPerlengkapanUtama []string `json:"kategori_utama,omitempty"`
		UpdatedAt                 string   `json:"update,omitempty"`
	}{
		KategoriPerlengkapanUtama: master.KategoriPerlengkapanUtama,
		UpdatedAt:                 master.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": respone})
}

func (dm *DataMasterControler) HapusKategoriUtamaPerlengkapan(ctx *gin.Context) {
	var payload *models.KategoriPerlengkapanUtamaInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

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

	for i, item := range master.KategoriPerlengkapanUtama {
		if item == payload.Kategori {
			master.KategoriPerlengkapanUtama = append(master.KategoriPerlengkapanUtama[:i], master.KategoriPerlengkapanUtama[i+1:]...)
			break
		}
	}

	for i, item := range master.KategoriPerlengkapan {
		if item.KategoriUtama == payload.Kategori {
			master.KategoriPerlengkapan = append(master.KategoriPerlengkapan[:i], master.KategoriPerlengkapan[i+1:]...)
			break
		}
	}

	for i, item := range master.PerlengkapanLaluLintas {
		if item.KategoriUtama == payload.Kategori {
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
		KategoriPerlengkapanUtama []string `json:"kategori_utama,omitempty"`
		UpdatedAt                 string   `json:"update,omitempty"`
	}{
		KategoriPerlengkapanUtama: master.KategoriPerlengkapanUtama,
		UpdatedAt:                 master.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": respone})
}

func (dm *DataMasterControler) EditKategoriUtamaPerlengkapan(ctx *gin.Context) {
	var payload *models.KategoriPerlengkapanUtamaEdit

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

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

	itemIndex := -1
	itemIndexKategori := -1
	itemIndexPerlengkapan := -1

	for i, item := range master.KategoriPerlengkapanUtama {
		if item == payload.Kategori {
			itemIndex = i
			break
		}
	}

	if itemIndex != -1 {
		master.KategoriPerlengkapanUtama[itemIndex] = payload.KategoriEdit
	}

	for i, item := range master.KategoriPerlengkapan {
		if item.KategoriUtama == payload.Kategori {
			itemIndexKategori = i
			break
		}
	}

	if itemIndexKategori != -1 {
		master.KategoriPerlengkapan[itemIndexKategori].KategoriUtama = payload.KategoriEdit
	}

	for i, item := range master.PerlengkapanLaluLintas {
		if item.KategoriUtama == payload.Kategori {
			itemIndexPerlengkapan = i
			break
		}
	}

	if itemIndexPerlengkapan != -1 {
		master.PerlengkapanLaluLintas[itemIndexPerlengkapan].KategoriUtama = payload.KategoriEdit
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
		KategoriPerlengkapanUtama []string `json:"kategori_utama,omitempty"`
		UpdatedAt                 string   `json:"update,omitempty"`
	}{
		KategoriPerlengkapanUtama: master.KategoriPerlengkapanUtama,
		UpdatedAt:                 master.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": respone})
}

func (dm *DataMasterControler) TambahKategoriPerlengkapan(ctx *gin.Context) {
	id := ctx.Param("id")

	var payload *models.KategoriPerlengkapanInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

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

	for _, kategori := range master.KategoriPerlengkapan {
		if kategori.KategoriUtama == payload.KategoriUtama && kategori.Kategori == payload.Kategori {
			kategoriExists = true
			break
		}
	}

	if !kategoriExists {
		ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "Data sudah ada"})
		return
	}

	if kategoriExists {
		master.KategoriPerlengkapan = append(master.KategoriPerlengkapan, models.KategoriPerlengkapan{KategoriUtama: payload.KategoriUtama, Kategori: payload.Kategori})
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
		KategoriPerlengkapan []models.KategoriPerlengkapan `json:"kategori_perlengkapan,omitempty"`
		UpdatedAt            string                        `json:"update,omitempty"`
	}{
		KategoriPerlengkapan: master.KategoriPerlengkapan,
		UpdatedAt:            master.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": respone})
}

func (dm *DataMasterControler) HapusKategoriPerlengkapan(ctx *gin.Context) {
	id := ctx.Param("id")

	var payload *models.KategoriPerlengkapanInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

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
		if item.KategoriUtama == payload.KategoriUtama && item.Kategori == payload.Kategori {
			master.KategoriPerlengkapan = append(master.KategoriPerlengkapan[:i], master.KategoriPerlengkapan[i+1:]...)
			break
		}
	}

	for i, item := range master.PerlengkapanLaluLintas {
		if item.KategoriUtama == payload.KategoriUtama && item.Kategori == payload.Kategori {
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
		KategoriPerlengkapan []models.KategoriPerlengkapan `json:"kategori_perlengkapan,omitempty"`
		UpdatedAt            string                        `json:"update,omitempty"`
	}{
		KategoriPerlengkapan: master.KategoriPerlengkapan,
		UpdatedAt:            master.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": respone})
}

func (dm *DataMasterControler) EditKategoriPerlengkapan(ctx *gin.Context) {
	id := ctx.Param("id")

	var payload *models.KategoriPerlengkapanEdit

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

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
		if item.KategoriUtama == payload.KategoriUtama && item.Kategori == payload.Kategori {
			itemIndex = i
			break
		}
	}

	if itemIndex != -1 {
		master.KategoriPerlengkapan[itemIndex].Kategori = payload.KategoriEdit
	}

	for i, item := range master.PerlengkapanLaluLintas {
		if item.KategoriUtama == payload.KategoriUtama && item.Kategori == payload.Kategori {
			itemIndexKategori = i
			break
		}
	}

	if itemIndexKategori != -1 {
		master.PerlengkapanLaluLintas[itemIndexKategori].Kategori = payload.KategoriEdit
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
		KategoriPerlengkapan []models.KategoriPerlengkapan `json:"kategori_perlengkapan,omitempty"`
		UpdatedAt            string                        `json:"update,omitempty"`
	}{
		KategoriPerlengkapan: master.KategoriPerlengkapan,
		UpdatedAt:            master.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": respone})
}

func (dm *DataMasterControler) TambahPerlengkapan(ctx *gin.Context) {
	id := ctx.Param("id")

	var payload *models.DataPerlengkapan

	if err := ctx.ShouldBind(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

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

	for i, item := range master.PerlengkapanLaluLintas {
		if item.KategoriUtama == payload.Perlengkapan.KategoriUtama && item.Kategori == payload.Perlengkapan.Kategori {
			kategoriExists = true
			itemIndex = i
			for _, item := range master.PerlengkapanLaluLintas[i].Perlengkapan {
				if item.JenisPerlengkapan == payload.Perlengkapan.Jenis {
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
			JenisPerlengkapan:  payload.Perlengkapan.Jenis,
			GambarPerlengkapan: data,
		}
		jenis := models.JenisPerlengkapan{
			KategoriUtama: payload.Perlengkapan.KategoriUtama,
			Kategori:      payload.Perlengkapan.Kategori,
			Perlengkapan:  []models.PerlengkapanItem{perlengkapan},
		}
		master.PerlengkapanLaluLintas = append(master.PerlengkapanLaluLintas, jenis)
	}

	if !perlengkapanExist && kategoriExists {
		perlengkapan := models.PerlengkapanItem{
			JenisPerlengkapan:  payload.Perlengkapan.Jenis,
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
	id := ctx.Param("id")

	var payload *models.JenisPerlengkapanInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

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

	for i, perlengkapan := range master.PerlengkapanLaluLintas {
		if perlengkapan.KategoriUtama == payload.KategoriUtama && perlengkapan.Kategori == payload.Kategori {
			for j, item := range master.PerlengkapanLaluLintas[i].Perlengkapan {
				if item.JenisPerlengkapan == payload.Jenis {
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
	id := ctx.Param("id")

	var payload *models.DataPerlengkapanEdit

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

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

	for i, perlengkapan := range master.PerlengkapanLaluLintas {
		if perlengkapan.KategoriUtama == payload.Perlengkapan.KategoriUtama && perlengkapan.Kategori == payload.Perlengkapan.Kategori {
			itemIndexKategori = i
			for j, item := range master.PerlengkapanLaluLintas[i].Perlengkapan {
				if item.JenisPerlengkapan == payload.Perlengkapan.Jenis {
					itemIndexPerlengkapan = j
					break
				}
			}
		}
	}

	file, _ := ctx.FormFile("perlengkapan")

	if itemIndexKategori != -1 && itemIndexPerlengkapan != -1 {
		master.PerlengkapanLaluLintas[itemIndexKategori].Perlengkapan[itemIndexPerlengkapan].JenisPerlengkapan = payload.Perlengkapan.JenisEdit
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
			Kebutuhan:             payload.Kebutuhan,
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
	id := ctx.Param("id")

	var payload *models.PersyaratanAndalalinHapus

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

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
		if master.Persyaratan.PersyaratanAndalalin[i].Persyaratan == payload.Persyaratan {
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
			for j, tambahan := range permohonan.BerkasPermohonan {
				if tambahan.Nama == payload.Persyaratan {
					oldSubstr := "/"
					newSubstr := "-"

					result := strings.Replace(permohonan.Kode, oldSubstr, newSubstr, -1)
					fileName := result + ".pdf"

					dataFile = append(dataFile, file{Name: fileName, File: tambahan.Berkas})
					andalalin[i].BerkasPermohonan = append(andalalin[i].BerkasPermohonan[:j], andalalin[i].BerkasPermohonan[j+1:]...)
					break
				}
			}
		}

		zipFile := payload.Persyaratan + ".zip"
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
	var payload *models.PersyaratanAndalalinEdit
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
		if master.Persyaratan.PersyaratanAndalalin[i].Persyaratan == payload.Persyaratan {
			itemIndex = i
			break
		}
	}

	if itemIndex != -1 {
		if master.Persyaratan.PersyaratanAndalalin[itemIndex].Kebutuhan != payload.Kebutuhan {
			master.Persyaratan.PersyaratanAndalalin[itemIndex].Kebutuhan = payload.Kebutuhan
		}

		if master.Persyaratan.PersyaratanAndalalin[itemIndex].Persyaratan != payload.PersyaratanEdit {
			master.Persyaratan.PersyaratanAndalalin[itemIndex].Persyaratan = payload.PersyaratanEdit
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
			Kebutuhan:             payload.Kebutuhan,
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
	id := ctx.Param("id")

	var payload *models.PersyaratanPerlalinHapus

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

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
		if master.Persyaratan.PersyaratanPerlalin[i].Persyaratan == payload.Persyaratan {
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
			for j, tambahan := range permohonan.BerkasPermohonan {
				if tambahan.Nama == payload.Persyaratan {
					oldSubstr := "/"
					newSubstr := "-"

					result := strings.Replace(permohonan.Kode, oldSubstr, newSubstr, -1)
					fileName := result + ".pdf"

					dataFile = append(dataFile, file{Name: fileName, File: tambahan.Berkas})
					perlalin[i].BerkasPermohonan = append(perlalin[i].BerkasPermohonan[:j], perlalin[i].BerkasPermohonan[j+1:]...)
					break
				}
			}
		}

		zipFile := payload.Persyaratan + ".zip"
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
	var payload *models.PersyaratanPerlalinEdit
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
		if master.Persyaratan.PersyaratanPerlalin[i].Persyaratan == payload.Persyaratan {
			itemIndex = i
			break
		}
	}

	if itemIndex != -1 {
		if master.Persyaratan.PersyaratanPerlalin[itemIndex].Kebutuhan != payload.Kebutuhan {
			master.Persyaratan.PersyaratanPerlalin[itemIndex].Kebutuhan = payload.Kebutuhan
		}

		if master.Persyaratan.PersyaratanPerlalin[itemIndex].Persyaratan != payload.PersyaratanEdit {
			master.Persyaratan.PersyaratanPerlalin[itemIndex].Persyaratan = payload.PersyaratanEdit
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
	id := ctx.Param("id")

	var payload *models.ProvinsiInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

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

	rows, err := dm.DB.Table("data_masters").Where("id_data_master", id).Select("provinsi", "updated_at").Rows()
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Data error"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		if err := dm.DB.ScanRows(rows, &master); err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Data error"})
			return
		}
	}

	exist := false
	notExistId := false
	var dataId string

	for _, item := range master.Provinsi {
		if item.Name == payload.Provinsi {
			exist = true
			break
		}
	}

	for _, item := range master.Provinsi {
		generate := strconv.Itoa(rand.Intn(10)*10 + rand.Intn(10))
		if item.Id != generate {
			dataId += generate
			notExistId = true
			break
		}
	}

	if !exist && notExistId {
		data := models.Provinsi{
			Id:   dataId,
			Name: payload.Provinsi,
		}
		master.Provinsi = append(master.Provinsi, data)
	}

	loc, _ := time.LoadLocation("Asia/Singapore")
	now := time.Now().In(loc).Format("02-01-2006")

	master.UpdatedAt = now + " " + time.Now().In(loc).Format("15:04:05")

	resultsSave := dm.DB.Table("data_masters").Where("id_data_master", id).Select("provinsi", "updated_at").Updates(models.DataMaster{Provinsi: master.Provinsi, UpdatedAt: master.UpdatedAt})
	if resultsSave.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": resultsSave.Error})
		return
	}

	respone := struct {
		Provinsi []models.Provinsi `json:"provinsi,omitempty"`
	}{
		Provinsi: master.Provinsi,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": respone})
}

func (dm *DataMasterControler) HapusProvinsi(ctx *gin.Context) {
	id := ctx.Param("id")

	var payload *models.ProvinsiInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

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

	rows, err := dm.DB.Table("data_masters").Where("id_data_master", id).Select("provinsi", "kabupaten", "kecamatan", "kelurahan", "updated_at").Rows()
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Data error"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		if err := dm.DB.ScanRows(rows, &master); err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Data error"})
			return
		}
	}

	var id_provinsi string
	var id_kabupaten string
	var id_kecamatan string

	for i, item := range master.Provinsi {
		if item.Name == payload.Provinsi {
			id_provinsi = item.Id
			master.Provinsi = append(master.Provinsi[:i], master.Provinsi[i+1:]...)
			break
		}
	}

	for i, item := range master.Kabupaten {
		if item.IdProvinsi == id_provinsi {
			id_kabupaten = item.Id
			master.Kabupaten = append(master.Kabupaten[:i], master.Kabupaten[i+1:]...)
			break
		}
	}

	for i, item := range master.Kecamatan {
		if item.IdKabupaten == id_kabupaten {
			id_kecamatan = item.Id
			master.Kecamatan = append(master.Kecamatan[:i], master.Kecamatan[i+1:]...)
			break
		}
	}

	for i, item := range master.Kelurahan {
		if item.IdKecamatan == id_kecamatan {
			master.Kelurahan = append(master.Kelurahan[:i], master.Kelurahan[i+1:]...)
			break
		}
	}

	loc, _ := time.LoadLocation("Asia/Singapore")
	now := time.Now().In(loc).Format("02-01-2006")

	master.UpdatedAt = now + " " + time.Now().In(loc).Format("15:04:05")

	resultsSave := dm.DB.Table("data_masters").Where("id_data_master", id).Select("provinsi", "kabupaten", "kecamatan", "kelurahan", "updated_at").Updates(models.DataMaster{Provinsi: master.Provinsi, Kabupaten: master.Kabupaten, Kecamatan: master.Kecamatan, Kelurahan: master.Kelurahan, UpdatedAt: master.UpdatedAt})
	if resultsSave.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": resultsSave.Error})
		return
	}

	respone := struct {
		Provinsi []models.Provinsi `json:"provinsi,omitempty"`
	}{
		Provinsi: master.Provinsi,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": respone})
}

func (dm *DataMasterControler) EditProvinsi(ctx *gin.Context) {
	id := ctx.Param("id")

	var payload *models.ProvinsiEdit

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

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

	rows, err := dm.DB.Table("data_masters").Where("id_data_master", id).Select("provinsi", "updated_at").Rows()
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Data error"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		if err := dm.DB.ScanRows(rows, &master); err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Data error"})
			return
		}
	}

	itemIndex := -1

	for i, item := range master.Provinsi {
		if item.Name == payload.Provinsi {
			itemIndex = i
			break
		}
	}

	if itemIndex != -1 {
		master.Provinsi[itemIndex].Name = payload.ProvinsiEdit
	}

	loc, _ := time.LoadLocation("Asia/Singapore")
	now := time.Now().In(loc).Format("02-01-2006")

	master.UpdatedAt = now + " " + time.Now().In(loc).Format("15:04:05")

	resultsSave := dm.DB.Table("data_masters").Where("id_data_master", id).Select("provinsi", "updated_at").Updates(models.DataMaster{Provinsi: master.Provinsi, UpdatedAt: master.UpdatedAt})
	if resultsSave.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": resultsSave.Error})
		return
	}

	respone := struct {
		Provinsi []models.Provinsi `json:"provinsi,omitempty"`
	}{
		Provinsi: master.Provinsi,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": respone})
}

func (dm *DataMasterControler) TambahKabupaten(ctx *gin.Context) {
	id := ctx.Param("id")

	var payload *models.KabupatenInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

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

	rows, err := dm.DB.Table("data_masters").Where("id_data_master", id).Select("provinsi", "kabupaten", "updated_at").Rows()
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Data error"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		if err := dm.DB.ScanRows(rows, &master); err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Data error"})
			return
		}
	}

	exist := false
	var id_provinsi string

	for _, item := range master.Provinsi {
		if item.Name == payload.Provinsi {
			id_provinsi += item.Id
			break
		}
	}

	for _, item := range master.Kabupaten {
		if item.Name == payload.Kabupaten {
			exist = true
			break
		}
	}

	notExistId := false
	var dataId string

	for _, item := range master.Kabupaten {
		generate := id_provinsi + strconv.Itoa(rand.Intn(10)*10+rand.Intn(10))
		if item.Id != generate {
			dataId += generate
			notExistId = true
			break
		}
	}

	if !exist && notExistId {
		data := models.Kabupaten{
			Id:         dataId,
			IdProvinsi: id_provinsi,
			Name:       payload.Kabupaten,
		}
		master.Kabupaten = append(master.Kabupaten, data)
	}

	loc, _ := time.LoadLocation("Asia/Singapore")
	now := time.Now().In(loc).Format("02-01-2006")

	master.UpdatedAt = now + " " + time.Now().In(loc).Format("15:04:05")

	resultsSave := dm.DB.Table("data_masters").Where("id_data_master", id).Select("kabupaten", "updated_at").Updates(models.DataMaster{Kabupaten: master.Kabupaten, UpdatedAt: master.UpdatedAt})
	if resultsSave.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": resultsSave.Error})
		return
	}

	respone := struct {
		Kabupaten []models.Kabupaten `json:"kabupaten,omitempty"`
	}{
		Kabupaten: master.Kabupaten,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": respone})
}

func (dm *DataMasterControler) HapusKabupaten(ctx *gin.Context) {
	id := ctx.Param("id")

	var payload *models.KabupatenHapus

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

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

	rows, err := dm.DB.Table("data_masters").Where("id_data_master", id).Select("kabupaten", "kecamatan", "kelurahan", "updated_at").Rows()
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Data error"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		if err := dm.DB.ScanRows(rows, &master); err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Data error"})
			return
		}
	}

	var id_kabupaten string
	var id_kecamatan string

	for i, item := range master.Kabupaten {
		if item.Name == payload.Kabupaten {
			id_kabupaten = item.Id
			master.Kabupaten = append(master.Kabupaten[:i], master.Kabupaten[i+1:]...)
			break
		}
	}

	for i, item := range master.Kecamatan {
		if item.IdKabupaten == id_kabupaten {
			id_kecamatan = item.Id
			master.Kecamatan = append(master.Kecamatan[:i], master.Kecamatan[i+1:]...)
			break
		}
	}

	for i, item := range master.Kelurahan {
		if item.IdKecamatan == id_kecamatan {
			master.Kelurahan = append(master.Kelurahan[:i], master.Kelurahan[i+1:]...)
			break
		}
	}

	loc, _ := time.LoadLocation("Asia/Singapore")
	now := time.Now().In(loc).Format("02-01-2006")

	master.UpdatedAt = now + " " + time.Now().In(loc).Format("15:04:05")

	resultsSave := dm.DB.Table("data_masters").Where("id_data_master", id).Select("kabupaten", "kecamatan", "kelurahan", "updated_at").Updates(models.DataMaster{Kabupaten: master.Kabupaten, Kecamatan: master.Kecamatan, Kelurahan: master.Kelurahan, UpdatedAt: master.UpdatedAt})
	if resultsSave.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": resultsSave.Error})
		return
	}

	respone := struct {
		Kabupaten []models.Kabupaten `json:"kabupaten,omitempty"`
	}{
		Kabupaten: master.Kabupaten,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": respone})
}

func (dm *DataMasterControler) EditKabupaten(ctx *gin.Context) {
	id := ctx.Param("id")

	var payload *models.KabupatenEdit

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

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

	rows, err := dm.DB.Table("data_masters").Where("id_data_master", id).Select("kabupaten", "updated_at").Rows()
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Data error"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		if err := dm.DB.ScanRows(rows, &master); err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Data error"})
			return
		}
	}
	itemIndex := -1

	for i, item := range master.Kabupaten {
		if item.Name == payload.Kabupaten {
			itemIndex = i
			break
		}
	}

	if itemIndex != -1 {
		master.Kabupaten[itemIndex].Name = payload.KabupatenEdit
	}

	loc, _ := time.LoadLocation("Asia/Singapore")
	now := time.Now().In(loc).Format("02-01-2006")

	master.UpdatedAt = now + " " + time.Now().In(loc).Format("15:04:05")

	resultsSave := dm.DB.Table("data_masters").Where("id_data_master", id).Select("kabupaten", "updated_at").Updates(models.DataMaster{Kabupaten: master.Kabupaten, UpdatedAt: master.UpdatedAt})
	if resultsSave.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": resultsSave.Error})
		return
	}

	respone := struct {
		Kabupaten []models.Kabupaten `json:"kabupaten,omitempty"`
	}{
		Kabupaten: master.Kabupaten,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": respone})
}

func (dm *DataMasterControler) TambahKecamatan(ctx *gin.Context) {
	id := ctx.Param("id")

	var payload *models.KecamatanInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

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

	rows, err := dm.DB.Table("data_masters").Where("id_data_master", id).Select("kabupaten", "kecamatan", "updated_at").Rows()
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Data error"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		if err := dm.DB.ScanRows(rows, &master); err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Data error"})
			return
		}
	}

	exist := false
	var id_kabupaten string

	for _, item := range master.Kabupaten {
		if item.Name == payload.Kabupaten {
			id_kabupaten += item.Id
			break
		}
	}

	for _, item := range master.Kecamatan {
		if item.Name == payload.Kecamatan {
			exist = true
			break
		}
	}

	notExistId := false
	var dataId string

	for _, item := range master.Kecamatan {
		generate := id_kabupaten + strconv.Itoa(rand.Intn(10)*10+rand.Intn(10))
		if item.Id != generate {
			dataId += generate
			notExistId = true
			break
		}
	}

	if !exist && notExistId {
		data := models.Kecamatan{
			Id:          dataId,
			IdKabupaten: id_kabupaten,
			Name:        payload.Kecamatan,
		}
		master.Kecamatan = append(master.Kecamatan, data)
	}

	loc, _ := time.LoadLocation("Asia/Singapore")
	now := time.Now().In(loc).Format("02-01-2006")

	master.UpdatedAt = now + " " + time.Now().In(loc).Format("15:04:05")

	resultsSave := dm.DB.Table("data_masters").Where("id_data_master", id).Select("kecamatan", "updated_at").Updates(models.DataMaster{Kecamatan: master.Kecamatan, UpdatedAt: master.UpdatedAt})
	if resultsSave.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": resultsSave.Error})
		return
	}

	respone := struct {
		Kecamatan []models.Kecamatan `json:"kecamatan,omitempty"`
	}{
		Kecamatan: master.Kecamatan,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": respone})
}

func (dm *DataMasterControler) HapusKecamatan(ctx *gin.Context) {
	id := ctx.Param("id")

	var payload *models.KecamatanHapus

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

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

	rows, err := dm.DB.Table("data_masters").Where("id_data_master", id).Select("kecamatan", "kelurahan", "updated_at").Rows()
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Data error"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		if err := dm.DB.ScanRows(rows, &master); err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Data error"})
			return
		}
	}

	var id_kecamatan string

	for i, item := range master.Kecamatan {
		if item.Name == payload.Kecamatan {
			id_kecamatan = item.Id
			master.Kecamatan = append(master.Kecamatan[:i], master.Kecamatan[i+1:]...)
			break
		}
	}

	for i, item := range master.Kelurahan {
		if item.IdKecamatan == id_kecamatan {
			master.Kelurahan = append(master.Kelurahan[:i], master.Kelurahan[i+1:]...)
			break
		}
	}

	loc, _ := time.LoadLocation("Asia/Singapore")
	now := time.Now().In(loc).Format("02-01-2006")

	master.UpdatedAt = now + " " + time.Now().In(loc).Format("15:04:05")

	resultsSave := dm.DB.Table("data_masters").Where("id_data_master", id).Select("kecamatan", "kelurahan", "updated_at").Updates(models.DataMaster{Kecamatan: master.Kecamatan, Kelurahan: master.Kelurahan, UpdatedAt: master.UpdatedAt})
	if resultsSave.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": resultsSave.Error})
		return
	}

	respone := struct {
		Kecamatan []models.Kecamatan `json:"kecamatan,omitempty"`
	}{
		Kecamatan: master.Kecamatan,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": respone})
}

func (dm *DataMasterControler) EditKecamatan(ctx *gin.Context) {
	id := ctx.Param("id")

	var payload *models.KecamatanEdit

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

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

	rows, err := dm.DB.Table("data_masters").Where("id_data_master", id).Select("kecamatan", "updated_at").Rows()
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Data error"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		if err := dm.DB.ScanRows(rows, &master); err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Data error"})
			return
		}
	}

	itemIndex := -1

	for i, item := range master.Kecamatan {
		if item.Name == payload.Kecamatan {
			itemIndex = i
			break
		}
	}

	if itemIndex != -1 {
		master.Kecamatan[itemIndex].Name = payload.KecamatanEdit
	}

	loc, _ := time.LoadLocation("Asia/Singapore")
	now := time.Now().In(loc).Format("02-01-2006")

	master.UpdatedAt = now + " " + time.Now().In(loc).Format("15:04:05")

	resultsSave := dm.DB.Table("data_masters").Where("id_data_master", id).Select("kecamatan", "updated_at").Updates(models.DataMaster{Kecamatan: master.Kecamatan, UpdatedAt: master.UpdatedAt})
	if resultsSave.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": resultsSave.Error})
		return
	}

	respone := struct {
		Kecamatan []models.Kecamatan `json:"kecamatan,omitempty"`
	}{
		Kecamatan: master.Kecamatan,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": respone})
}

func (dm *DataMasterControler) TambahKelurahan(ctx *gin.Context) {
	id := ctx.Param("id")

	var payload *models.KelurahanInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

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

	rows, err := dm.DB.Table("data_masters").Where("id_data_master", id).Select("kecamatan", "kelurahan", "updated_at").Rows()
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Data error"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		if err := dm.DB.ScanRows(rows, &master); err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Data error"})
			return
		}
	}

	exist := false
	var id_kecamatan string

	for _, item := range master.Kecamatan {
		if item.Name == payload.Kecamatan {
			id_kecamatan += item.Id
			break
		}
	}

	for _, item := range master.Kelurahan {
		if item.Name == payload.Kelurahan {
			exist = true
			break
		}
	}

	notExistId := false
	var dataId string

	for _, item := range master.Kelurahan {
		generate := id_kecamatan + strconv.Itoa(rand.Intn(10)*10+rand.Intn(10))
		if item.Id != generate {
			dataId += generate
			notExistId = true
			break
		}
	}

	if !exist && notExistId {
		data := models.Kelurahan{
			Id:          dataId,
			IdKecamatan: id_kecamatan,
			Name:        payload.Kelurahan,
		}
		master.Kelurahan = append(master.Kelurahan, data)
	}

	loc, _ := time.LoadLocation("Asia/Singapore")
	now := time.Now().In(loc).Format("02-01-2006")

	master.UpdatedAt = now + " " + time.Now().In(loc).Format("15:04:05")

	resultsSave := dm.DB.Table("data_masters").Where("id_data_master", id).Select("kelurahan", "updated_at").Updates(models.DataMaster{Kelurahan: master.Kelurahan, UpdatedAt: master.UpdatedAt})
	if resultsSave.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": resultsSave.Error})
		return
	}

	respone := struct {
		Kelurahan []models.Kelurahan `json:"kelurahan,omitempty"`
	}{
		Kelurahan: master.Kelurahan,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": respone})
}

func (dm *DataMasterControler) HapusKelurahan(ctx *gin.Context) {
	id := ctx.Param("id")

	var payload *models.KelurahanHapus

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

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

	rows, err := dm.DB.Table("data_masters").Where("id_data_master", id).Select("kelurahan", "updated_at").Rows()
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Data error"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		if err := dm.DB.ScanRows(rows, &master); err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Data error"})
			return
		}
	}

	for i, item := range master.Kelurahan {
		if item.Name == payload.Kelurahan {
			master.Kelurahan = append(master.Kelurahan[:i], master.Kelurahan[i+1:]...)
			break
		}
	}

	loc, _ := time.LoadLocation("Asia/Singapore")
	now := time.Now().In(loc).Format("02-01-2006")

	master.UpdatedAt = now + " " + time.Now().In(loc).Format("15:04:05")

	resultsSave := dm.DB.Table("data_masters").Where("id_data_master", id).Select("kelurahan", "updated_at").Updates(models.DataMaster{Kelurahan: master.Kelurahan, UpdatedAt: master.UpdatedAt})
	if resultsSave.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": resultsSave.Error})
		return
	}

	respone := struct {
		Kelurahan []models.Kelurahan `json:"kelurahan,omitempty"`
	}{
		Kelurahan: master.Kelurahan,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": respone})
}

func (dm *DataMasterControler) EditKelurahan(ctx *gin.Context) {
	id := ctx.Param("id")

	var payload *models.KelurahanEdit

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

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

	rows, err := dm.DB.Table("data_masters").Where("id_data_master", id).Select("kelurahan", "updated_at").Rows()
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Data error"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		if err := dm.DB.ScanRows(rows, &master); err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Data error"})
			return
		}
	}

	itemIndex := -1

	for i, item := range master.Kelurahan {
		if item.Name == payload.Kelurahan {
			itemIndex = i
			break
		}
	}

	if itemIndex != -1 {
		master.Kelurahan[itemIndex].Name = payload.KelurahanEdit
	}

	loc, _ := time.LoadLocation("Asia/Singapore")
	now := time.Now().In(loc).Format("02-01-2006")

	master.UpdatedAt = now + " " + time.Now().In(loc).Format("15:04:05")

	resultsSave := dm.DB.Table("data_masters").Where("id_data_master", id).Select("kelurahan", "updated_at").Updates(models.DataMaster{Kelurahan: master.Kelurahan, UpdatedAt: master.UpdatedAt})
	if resultsSave.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": resultsSave.Error})
		return
	}

	respone := struct {
		Kelurahan []models.Kelurahan `json:"kelurahan,omitempty"`
	}{
		Kelurahan: master.Kelurahan,
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": respone})
}

func (dm *DataMasterControler) TambahJenisProyek(ctx *gin.Context) {
	id := ctx.Param("id")

	var payload *models.JenisProyekInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

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

	rows, err := dm.DB.Table("data_masters").Where("id_data_master", id).Select("jenis_proyek", "updated_at").Rows()
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Data error"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		if err := dm.DB.ScanRows(rows, &master); err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Data error"})
			return
		}
	}

	exist := contains(master.JenisProyek, payload.Jenis)

	if exist {
		ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "Data sudah ada"})
		return
	}

	master.JenisProyek = append(master.JenisProyek, payload.Jenis)

	loc, _ := time.LoadLocation("Asia/Singapore")
	now := time.Now().In(loc).Format("02-01-2006")

	master.UpdatedAt = now + " " + time.Now().In(loc).Format("15:04:05")
	resultsSave := dm.DB.Table("data_masters").Where("id_data_master", id).Select("jenis_proyek", "updated_at").Updates(models.DataMaster{JenisProyek: master.JenisProyek, UpdatedAt: master.UpdatedAt})
	if resultsSave.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": resultsSave.Error})
		return
	}

	respone := struct {
		JenisProyek []string `json:"jenis_proyek,omitempty"`
	}{
		JenisProyek: master.JenisProyek,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": respone})
}

func (dm *DataMasterControler) HapusJenisProyek(ctx *gin.Context) {
	id := ctx.Param("id")

	var payload *models.JenisProyekInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

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

	rows, err := dm.DB.Table("data_masters").Where("id_data_master", id).Select("jenis_proyek", "updated_at").Rows()
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Data error"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		if err := dm.DB.ScanRows(rows, &master); err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Data error"})
			return
		}
	}

	for i, item := range master.JenisProyek {
		if item == payload.Jenis {
			master.JenisProyek = append(master.JenisProyek[:i], master.JenisProyek[i+1:]...)
			break
		}
	}

	loc, _ := time.LoadLocation("Asia/Singapore")
	now := time.Now().In(loc).Format("02-01-2006")

	master.UpdatedAt = now + " " + time.Now().In(loc).Format("15:04:05")

	resultsSave := dm.DB.Table("data_masters").Where("id_data_master", id).Select("jenis_proyek", "updated_at").Updates(models.DataMaster{JenisProyek: master.JenisProyek, UpdatedAt: master.UpdatedAt})
	if resultsSave.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": resultsSave.Error})
		return
	}

	respone := struct {
		JenisProyek []string `json:"jenis_proyek,omitempty"`
	}{
		JenisProyek: master.JenisProyek,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": respone})
}

func (dm *DataMasterControler) EditJenisProyek(ctx *gin.Context) {
	id := ctx.Param("id")

	var payload *models.JenisProyekEdit

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

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

	rows, err := dm.DB.Table("data_masters").Where("id_data_master", id).Select("jenis_proyek", "updated_at").Rows()
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Data error"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		if err := dm.DB.ScanRows(rows, &master); err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Data error"})
			return
		}
	}

	itemIndex := -1

	for i, item := range master.JenisProyek {
		if item == payload.Jenis {
			itemIndex = i
			break
		}
	}

	if itemIndex != -1 {
		master.JenisProyek[itemIndex] = payload.JenisEdit
	}

	loc, _ := time.LoadLocation("Asia/Singapore")
	now := time.Now().In(loc).Format("02-01-2006")

	master.UpdatedAt = now + " " + time.Now().In(loc).Format("15:04:05")

	resultsSave := dm.DB.Table("data_masters").Where("id_data_master", id).Select("jenis_proyek", "updated_at").Updates(models.DataMaster{JenisProyek: master.JenisProyek, UpdatedAt: master.UpdatedAt})
	if resultsSave.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": resultsSave.Error})
		return
	}

	respone := struct {
		JenisProyek []string `json:"jenis_proyek,omitempty"`
	}{
		JenisProyek: master.JenisProyek,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": respone})
}

func (dm *DataMasterControler) TambahJalan(ctx *gin.Context) {
	var payload *models.JalanInput
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

	rows, err := dm.DB.Table("data_masters").Where("id_data_master", id).Select("jalan", "updated_at").Rows()
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Data error"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		if err := dm.DB.ScanRows(rows, &master); err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Data error"})
			return
		}
	}

	jalanExist := false

	for _, item := range master.Jalan {
		if item.KodeJalan == payload.KodeJalan && item.Nama == payload.Nama {
			jalanExist = true
			ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "Data sudah ada"})
			return
		}
	}

	if !jalanExist {
		jalan := models.Jalan{
			KodeProvinsi:  "63",
			KodeKabupaten: "71",
			KodeKecamatan: payload.KodeKecamatan,
			KodeKelurahan: payload.KodeKelurahan,
			KodeJalan:     payload.KodeJalan,
			Nama:          payload.Nama,
			Pangkal:       payload.Pangkal,
			Ujung:         payload.Ujung,
			Kelurahan:     payload.Kelurahan,
			Kecamatan:     payload.Kecamatan,
			Panjang:       payload.Panjang,
			Lebar:         payload.Lebar,
			Permukaan:     payload.Permukaan,
			Fungsi:        payload.Fungsi,
		}
		master.Jalan = append(master.Jalan, jalan)
	}

	loc, _ := time.LoadLocation("Asia/Singapore")
	now := time.Now().In(loc).Format("02-01-2006")

	master.UpdatedAt = now + " " + time.Now().In(loc).Format("15:04:05")

	resultsSave := dm.DB.Table("data_masters").Where("id_data_master", id).Select("jalan", "updated_at").Updates(models.DataMaster{Jalan: master.Jalan, UpdatedAt: master.UpdatedAt})
	if resultsSave.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": resultsSave.Error})
		return
	}

	respone := struct {
		Jalan []models.Jalan `json:"jalan,omitempty"`
	}{
		Jalan: master.Jalan,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": respone})
}

func (dm *DataMasterControler) HapusJalan(ctx *gin.Context) {
	id := ctx.Param("id")

	var payload *models.JalanHapus

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

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

	rows, err := dm.DB.Table("data_masters").Where("id_data_master", id).Select("jalan", "updated_at").Rows()
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Data error"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		if err := dm.DB.ScanRows(rows, &master); err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Data error"})
			return
		}
	}

	for i, item := range master.Jalan {
		if item.KodeJalan == payload.Kode {
			master.Jalan = append(master.Jalan[:i], master.Jalan[i+1:]...)
			break
		}
	}

	loc, _ := time.LoadLocation("Asia/Singapore")
	now := time.Now().In(loc).Format("02-01-2006")

	master.UpdatedAt = now + " " + time.Now().In(loc).Format("15:04:05")

	resultsSave := dm.DB.Table("data_masters").Where("id_data_master", id).Select("jalan", "updated_at").Updates(models.DataMaster{Jalan: master.Jalan, UpdatedAt: master.UpdatedAt})
	if resultsSave.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": resultsSave.Error})
		return
	}

	respone := struct {
		Jalan []models.Jalan `json:"jalan,omitempty"`
	}{
		Jalan: master.Jalan,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": respone})
}

func (dm *DataMasterControler) EditJalan(ctx *gin.Context) {
	var payload *models.JalanEdit
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

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	var master models.DataMaster

	rows, err := dm.DB.Table("data_masters").Where("id_data_master", id).Select("jalan", "updated_at").Rows()
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Data error"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		if err := dm.DB.ScanRows(rows, &master); err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Data error"})
			return
		}
	}

	itemIndex := -1

	for i, item := range master.Jalan {
		if item.Nama == payload.Jalan && item.KodeJalan == payload.Kode {
			itemIndex = i
			break
		}
	}

	if itemIndex != -1 {
		master.Jalan[itemIndex].KodeJalan = payload.KodeJalan
		master.Jalan[itemIndex].KodeJalan = payload.KodeJalan
		master.Jalan[itemIndex].Nama = payload.Nama
		master.Jalan[itemIndex].Pangkal = payload.Pangkal
		master.Jalan[itemIndex].Ujung = payload.Ujung
		master.Jalan[itemIndex].Kelurahan = payload.Kelurahan
		master.Jalan[itemIndex].Kecamatan = payload.Kecamatan
		master.Jalan[itemIndex].Panjang = payload.Panjang
		master.Jalan[itemIndex].Lebar = payload.Lebar
		master.Jalan[itemIndex].Permukaan = payload.Permukaan
		master.Jalan[itemIndex].Fungsi = payload.Fungsi
	}

	loc, _ := time.LoadLocation("Asia/Singapore")
	now := time.Now().In(loc).Format("02-01-2006")

	master.UpdatedAt = now + " " + time.Now().In(loc).Format("15:04:05")

	resultsSave := dm.DB.Table("data_masters").Where("id_data_master", id).Select("jalan", "updated_at").Updates(models.DataMaster{Jalan: master.Jalan, UpdatedAt: master.UpdatedAt})
	if resultsSave.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": resultsSave.Error})
		return
	}

	respone := struct {
		Jalan []models.Jalan `json:"jalan,omitempty"`
	}{
		Jalan: master.Jalan,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": respone})
}
