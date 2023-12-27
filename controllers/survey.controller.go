package controllers

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"andalalin/initializers"
	"andalalin/models"
	"andalalin/repository"
	"andalalin/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	_ "time/tzdata"
)

type SurveyController struct {
	DB *gorm.DB
}

func NewSurveyController(DB *gorm.DB) SurveyController {
	return SurveyController{DB}
}

type data struct {
	Jenis string
	Nilai int
	Hasil string
}

type komentar struct {
	Nama     string
	Komentar string
}

func interval(hasil float64) string {
	intervalNilai := ""
	if hasil < 26.00 {
		cekInterval := float64(hasil) / float64(25.00)
		intervalNilai = fmt.Sprintf("%.2f", cekInterval)
	} else if hasil >= 43.76 && hasil <= 62.50 {
		cekInterval := float64(hasil) / float64(25.00)
		intervalNilai = fmt.Sprintf("%.2f", cekInterval)
	} else if hasil >= 62.51 && hasil <= 81.25 {
		cekInterval := float64(hasil) / float64(25.00)
		intervalNilai = fmt.Sprintf("%.2f", cekInterval)
	} else if hasil >= 81.26 && hasil <= 100 {
		cekInterval := float64(hasil) / float64(25.00)
		intervalNilai = fmt.Sprintf("%.2f", cekInterval)
	}
	return intervalNilai
}

func mutu(hasil float64) string {
	mutuNilai := ""
	if hasil <= 43.75 {
		mutuNilai = "D"
	} else if hasil >= 43.76 && hasil <= 62.50 {
		mutuNilai = "C"
	} else if hasil >= 62.51 && hasil <= 81.25 {
		mutuNilai = "B"
	} else if hasil >= 81.26 && hasil <= 100 {
		mutuNilai = "A"
	}
	return mutuNilai
}

func kinerja(hasil float64) string {
	kinerjaNilai := ""
	if hasil <= 43.75 {
		kinerjaNilai = "Buruk"
	} else if hasil >= 43.76 && hasil <= 62.50 {
		kinerjaNilai = "Kurang baik"
	} else if hasil >= 62.51 && hasil <= 81.25 {
		kinerjaNilai = "Baik"
	} else if hasil >= 81.26 && hasil <= 100 {
		kinerjaNilai = "Sangat baik"
	}
	return kinerjaNilai
}

func getStartOfMonth(year int, month time.Month) time.Time {
	return time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
}

func getEndOfMonth(year int, month time.Month) time.Time {
	nextMonth := getStartOfMonth(year, month).AddDate(0, 1, 0)
	return nextMonth.Add(-time.Second)
}

func (sc *SurveyController) SurveiKepuasan(ctx *gin.Context) {
	var payload *models.SurveiKepuasanInput
	id := ctx.Param("id_andalalin")
	currentUser := ctx.MustGet("currentUser").(models.User)

	if err := ctx.ShouldBind(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	loc, _ := time.LoadLocation("Asia/Singapore")
	nowTime := time.Now().In(loc)

	tanggal := nowTime.Format("02") + " " + utils.Bulan(nowTime.Month()) + " " + nowTime.Format("2006")

	var andalalin models.Perlalin
	var perlalin models.Andalalin

	sc.DB.First(&andalalin, "id_andalalin = ?", id)
	sc.DB.First(&perlalin, "id_andalalin = ?", id)

	if andalalin.IdAndalalin != uuid.Nil {
		kepuasan := models.SurveiKepuasan{
			IdAndalalin:        andalalin.IdAndalalin,
			IdUser:             currentUser.ID,
			Nama:               currentUser.Name,
			Email:              currentUser.Email,
			KritikSaran:        payload.KritikSaran,
			TanggalPelaksanaan: tanggal,
			DataSurvei:         payload.DataSurvei,
		}

		result := sc.DB.Create(&kepuasan)

		if result.Error != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Telah terjadi sesuatu"})
			return
		}
	}

	if perlalin.IdAndalalin != uuid.Nil {
		kepuasan := models.SurveiKepuasan{
			IdAndalalin:        perlalin.IdAndalalin,
			IdUser:             currentUser.ID,
			Nama:               currentUser.Name,
			Email:              currentUser.Email,
			KritikSaran:        payload.KritikSaran,
			TanggalPelaksanaan: tanggal,
			DataSurvei:         payload.DataSurvei,
		}

		result := sc.DB.Create(&kepuasan)

		if result.Error != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Telah terjadi sesuatu"})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (sc *SurveyController) CekSurveiKepuasan(ctx *gin.Context) {
	id := ctx.Param("id_andalalin")

	var survei models.SurveiKepuasan

	result := sc.DB.First(&survei, "id_andalalin", id)

	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Telah terjadi sesuatu"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (sc *SurveyController) HasilSurveiKepuasan(ctx *gin.Context) {
	config, _ := initializers.LoadConfig()

	accessUser := ctx.MustGet("accessUser").(string)

	claim, error := utils.ValidateToken(accessUser, config.AccessTokenPublicKey)
	if error != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": error.Error()})
		return
	}

	credential := claim.Credentials[repository.AndalalinSurveiKepuasan]

	if !credential {
		// Return status 403 and permission denied error message.
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": true,
			"msg":   "Permission denied",
		})
		return
	}

	loc, _ := time.LoadLocation("Asia/Singapore")
	nowTime := time.Now().In(loc)

	startOfMonth := getStartOfMonth(nowTime.Year(), nowTime.Month())

	endOfMonth := getEndOfMonth(nowTime.Year(), nowTime.Month())

	periode := startOfMonth.Format("02") + " - " + endOfMonth.Format("02") + " " + utils.Bulan(nowTime.Month()) + " " + nowTime.Format("2006")

	var survei []models.SurveiKepuasan

	result := sc.DB.Where("tanggal_pelaksanaan LIKE ?", fmt.Sprintf("%%%s%%", utils.Bulan(nowTime.Month())+" "+nowTime.Format("2006"))).Find(&survei)

	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Telah terjadi sesuatu"})
		return
	}

	nilai := []data{}

	nilai = append(nilai, data{Jenis: "Persyaratan pelayanan", Nilai: 0, Hasil: "0"})
	nilai = append(nilai, data{Jenis: "Prosedur pelayanan", Nilai: 0, Hasil: "0"})
	nilai = append(nilai, data{Jenis: "Waktu pelayanan", Nilai: 0, Hasil: "0"})
	nilai = append(nilai, data{Jenis: "Biaya / tarif pelayanan", Nilai: 0, Hasil: "0"})
	nilai = append(nilai, data{Jenis: "Produk pelayanan", Nilai: 0, Hasil: "0"})
	nilai = append(nilai, data{Jenis: "Kompetensi pelaksana", Nilai: 0, Hasil: "0"})
	nilai = append(nilai, data{Jenis: "Perilaku / sikap petugas", Nilai: 0, Hasil: "0"})
	nilai = append(nilai, data{Jenis: "Maklumat pelayanan", Nilai: 0, Hasil: "0"})
	nilai = append(nilai, data{Jenis: "Ketersediaan sarana pengaduan", Nilai: 0, Hasil: "0"})

	komen := []komentar{}

	for _, data := range survei {
		komen = append(komen, komentar{Nama: data.Nama, Komentar: *data.KritikSaran})
		for _, isi := range data.DataSurvei {
			for i, item := range nilai {
				if item.Jenis == isi.Jenis {
					switch isi.Nilai {
					case "Sangat baik":
						nilai[i].Nilai = nilai[i].Nilai + 4
					case "Baik":
						nilai[i].Nilai = nilai[i].Nilai + 3
					case "Kurang baik":
						nilai[i].Nilai = nilai[i].Nilai + 2
					case "Buruk":
						nilai[i].Nilai = nilai[i].Nilai + 1
					}
					break
				}
			}
		}
	}

	total := 0

	for i, item := range nilai {
		hasil := float64(item.Nilai) * float64(100) / float64(len(survei)) / float64(4)
		nilai[i].Hasil = fmt.Sprintf("%.2f", hasil)
		total = total + item.Nilai
	}

	indeksHasil := float64(total) * float64(100) / float64(9) / float64(4) / float64(len(survei))
	indeks := fmt.Sprintf("%.2f", indeksHasil)

	hasil := struct {
		Periode        string     `json:"periode,omitempty"`
		Responden      string     `json:"responden,omitempty"`
		IndeksKepuasan string     `json:"indeks_kepuasan,omitempty"`
		NilaiInterval  string     `json:"nilai_interval,omitempty"`
		Mutu           string     `json:"mutu,omitempty"`
		Kinerja        string     `json:"kinerja,omitempty"`
		DataHasil      []data     `json:"hasil,omitempty"`
		Komentar       []komentar `json:"komentar,omitempty"`
	}{
		Periode:        periode,
		Responden:      strconv.Itoa(len(survei)),
		IndeksKepuasan: indeks,
		NilaiInterval:  interval(indeksHasil),
		Mutu:           mutu(indeksHasil),
		Kinerja:        kinerja(indeksHasil),
		DataHasil:      nilai,
		Komentar:       komen,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": hasil})
}

func (sc *SurveyController) HasilSurveiKepuasanTertentu(ctx *gin.Context) {
	config, _ := initializers.LoadConfig()
	waktu := ctx.Param("waktu")

	accessUser := ctx.MustGet("accessUser").(string)

	claim, error := utils.ValidateToken(accessUser, config.AccessTokenPublicKey)
	if error != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": error.Error()})
		return
	}

	credential := claim.Credentials[repository.AndalalinSurveiKepuasan]

	if !credential {
		// Return status 403 and permission denied error message.
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": true,
			"msg":   "Permission denied",
		})
		return
	}

	tahun, err := strconv.Atoi(waktu[3:7])
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Telah terjadi sesuatu"})
		return
	}

	convertBulan, err := strconv.Atoi(waktu[0:2])
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Telah terjadi sesuatu"})
		return
	}

	bulan := time.Month(convertBulan)

	startOfMonth := getStartOfMonth(tahun, bulan)

	endOfMonth := getEndOfMonth(tahun, bulan)

	periode := startOfMonth.Format("02") + " - " + endOfMonth.Format("02") + " " + utils.Bulan(bulan) + " " + waktu[3:7]

	var survei []models.SurveiKepuasan

	result := sc.DB.Where("tanggal_pelaksanaan LIKE ?", fmt.Sprintf("%%%s%%", utils.Bulan(bulan)+" "+waktu[3:7])).Find(&survei)

	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Telah terjadi sesuatu"})
		return
	}

	nilai := []data{}

	nilai = append(nilai, data{Jenis: "Persyaratan pelayanan", Nilai: 0, Hasil: "0"})
	nilai = append(nilai, data{Jenis: "Prosedur pelayanan", Nilai: 0, Hasil: "0"})
	nilai = append(nilai, data{Jenis: "Waktu pelayanan", Nilai: 0, Hasil: "0"})
	nilai = append(nilai, data{Jenis: "Biaya / tarif pelayanan", Nilai: 0, Hasil: "0"})
	nilai = append(nilai, data{Jenis: "Produk pelayanan", Nilai: 0, Hasil: "0"})
	nilai = append(nilai, data{Jenis: "Kompetensi pelaksana", Nilai: 0, Hasil: "0"})
	nilai = append(nilai, data{Jenis: "Perilaku / sikap petugas", Nilai: 0, Hasil: "0"})
	nilai = append(nilai, data{Jenis: "Maklumat pelayanan", Nilai: 0, Hasil: "0"})
	nilai = append(nilai, data{Jenis: "Ketersediaan sarana pengaduan", Nilai: 0, Hasil: "0"})

	komen := []komentar{}

	for _, data := range survei {
		komen = append(komen, komentar{Nama: data.Nama, Komentar: *data.KritikSaran})
		for _, isi := range data.DataSurvei {
			for i, item := range nilai {
				if item.Jenis == isi.Jenis {
					switch isi.Nilai {
					case "Sangat baik":
						nilai[i].Nilai = nilai[i].Nilai + 4
					case "Baik":
						nilai[i].Nilai = nilai[i].Nilai + 3
					case "Kurang baik":
						nilai[i].Nilai = nilai[i].Nilai + 2
					case "Buruk":
						nilai[i].Nilai = nilai[i].Nilai + 1
					}
					break
				}
			}
		}
	}

	total := 0

	for i, item := range nilai {
		hasil := float64(item.Nilai) * float64(100) / float64(len(survei)) / float64(4)
		nilai[i].Hasil = fmt.Sprintf("%.2f", hasil)
		total = total + item.Nilai
	}

	indeksHasil := float64(total) * float64(100) / float64(9) / float64(4) / float64(len(survei))
	indeks := fmt.Sprintf("%.2f", indeksHasil)

	hasil := struct {
		Periode        string     `json:"periode,omitempty"`
		Responden      string     `json:"responden,omitempty"`
		IndeksKepuasan string     `json:"indeks_kepuasan,omitempty"`
		NilaiInterval  string     `json:"nilai_interval,omitempty"`
		Mutu           string     `json:"mutu,omitempty"`
		Kinerja        string     `json:"kinerja,omitempty"`
		DataHasil      []data     `json:"hasil,omitempty"`
		Komentar       []komentar `json:"komentar,omitempty"`
	}{
		Periode:        periode,
		Responden:      strconv.Itoa(len(survei)),
		IndeksKepuasan: indeks,
		NilaiInterval:  interval(indeksHasil),
		Mutu:           mutu(indeksHasil),
		Kinerja:        kinerja(indeksHasil),
		DataHasil:      nilai,
		Komentar:       komen,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": hasil})
}

func (sc *SurveyController) IsiSurveyMandiri(ctx *gin.Context) {
	var payload *models.DataSurvey
	currentUser := ctx.MustGet("currentUser").(models.User)

	config, _ := initializers.LoadConfig()

	accessUser := ctx.MustGet("accessUser").(string)

	claim, error := utils.ValidateToken(accessUser, config.AccessTokenPublicKey)
	if error != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": error.Error()})
		return
	}

	credential := claim.Credentials[repository.AndalalinSurveyCredential]

	if !credential {
		// Return status 403 and permission denied error message.
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": true,
			"msg":   "Permission denied",
		})
		return
	}

	if err := ctx.ShouldBind(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	loc, _ := time.LoadLocation("Asia/Singapore")
	nowTime := time.Now().In(loc)

	tanggal := nowTime.Format("02") + " " + utils.Bulan(nowTime.Month()) + " " + nowTime.Format("2006")

	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	blobs := make(map[string][]byte)

	for key, files := range form.File {
		for _, file := range files {
			// Save the uploaded file with key as prefix
			file, err := file.Open()

			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			defer file.Close()

			data, err := io.ReadAll(file)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			// Store the blob data in the map
			blobs[key] = data
		}
	}

	survey := models.SurveiMandiri{
		IdPetugas:     currentUser.ID,
		Petugas:       currentUser.Name,
		EmailPetugas:  currentUser.Email,
		Lokasi:        payload.Data.Lokasi,
		Catatan:       payload.Data.Catatan,
		Foto1:         blobs["foto1"],
		Foto2:         blobs["foto2"],
		Foto3:         blobs["foto3"],
		Latitude:      payload.Data.Latitude,
		Longitude:     payload.Data.Longitude,
		StatusSurvei:  "Perlu tindakan",
		TanggalSurvei: tanggal,
		WaktuSurvei:   nowTime.Format("15:04:05"),
	}

	result := sc.DB.Create(&survey)

	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Telah terjadi sesuatu"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": survey})
}

func (sc *SurveyController) GetAllSurveiMandiri(ctx *gin.Context) {
	config, _ := initializers.LoadConfig()

	accessUser := ctx.MustGet("accessUser").(string)

	claim, error := utils.ValidateToken(accessUser, config.AccessTokenPublicKey)
	if error != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": error.Error()})
		return
	}

	credential := claim.Credentials[repository.AndalalinSurveyCredential]

	if !credential {
		// Return status 403 and permission denied error message.
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": true,
			"msg":   "Permission denied",
		})
		return
	}

	var survey []models.SurveiMandiri

	results := sc.DB.Find(&survey)

	if results.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": results.Error})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "results": len(survey), "data": survey})
}

func (sc *SurveyController) GetAllSurveiMandiriByPetugas(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)
	config, _ := initializers.LoadConfig()

	accessUser := ctx.MustGet("accessUser").(string)

	claim, error := utils.ValidateToken(accessUser, config.AccessTokenPublicKey)
	if error != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": error.Error()})
		return
	}

	credential := claim.Credentials[repository.AndalalinSurveyCredential]

	if !credential {
		// Return status 403 and permission denied error message.
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": true,
			"msg":   "Permission denied",
		})
		return
	}

	var survey []models.SurveiMandiri

	results := sc.DB.Find(&survey, "id_petugas = ?", currentUser.ID)

	if results.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": results.Error})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "results": len(survey), "data": survey})
}

func (sc *SurveyController) GetSurveiMandiri(ctx *gin.Context) {
	id := ctx.Param("id_survei")

	config, _ := initializers.LoadConfig()

	accessUser := ctx.MustGet("accessUser").(string)

	claim, error := utils.ValidateToken(accessUser, config.AccessTokenPublicKey)
	if error != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": error.Error()})
		return
	}

	credential := claim.Credentials[repository.AndalalinSurveyCredential]

	if !credential {
		// Return status 403 and permission denied error message.
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": true,
			"msg":   "Permission denied",
		})
		return
	}

	var survey *models.SurveiMandiri

	result := sc.DB.First(&survey, "id_survey = ?", id)
	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": survey})
}

func (sc *SurveyController) TerimaSurvei(ctx *gin.Context) {
	id := ctx.Param("id_survei")
	var payload *models.TerimaSurveiMandiri

	config, _ := initializers.LoadConfig()

	accessUser := ctx.MustGet("accessUser").(string)

	claim, error := utils.ValidateToken(accessUser, config.AccessTokenPublicKey)
	if error != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": error.Error()})
		return
	}

	credential := claim.Credentials[repository.AndalalinSurveyCredential]

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

	var survey *models.SurveiMandiri

	result := sc.DB.First(&survey, "id_survey = ?", id)
	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error})
		return
	}
	survey.StatusSurvei = "Survei diterima"
	survey.CatatanTindakan = payload.Catatan

	sc.DB.Save(&survey)

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": survey})
}
