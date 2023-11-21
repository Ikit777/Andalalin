package controllers

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Ikit777/E-Andalalin/initializers"
	"github.com/Ikit777/E-Andalalin/models"
	"github.com/Ikit777/E-Andalalin/repository"
	"github.com/Ikit777/E-Andalalin/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	_ "time/tzdata"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"

	"github.com/lukasjarosch/go-docx"
)

type AndalalinController struct {
	DB *gorm.DB
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

func NewAndalalinController(DB *gorm.DB) AndalalinController {
	return AndalalinController{DB}
}

func (ac *AndalalinController) Pengajuan(ctx *gin.Context) {
	var payload *models.DataAndalalin
	currentUser := ctx.MustGet("currentUser").(models.User)

	config, _ := initializers.LoadConfig()

	accessUser := ctx.MustGet("accessUser").(string)

	claim, error := utils.ValidateToken(accessUser, config.AccessTokenPublicKey)
	if error != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": error.Error()})
		return
	}

	credential := claim.Credentials[repository.AndalalinPengajuanCredential]

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

	kode := "andalalin/" + utils.Generate(6)
	tanggal := nowTime.Format("02") + " " + utils.Bulan(nowTime.Month()) + " " + nowTime.Format("2006")

	var path string

	if payload.Andalalin.Pemohon == "Perorangan" {
		path = "templates/tandaterimaTemplatePerorangan.html"
	} else {
		path = "templates/tandaterimaTemplate.html"
	}

	t, err := template.ParseFiles(path)
	if err != nil {
		log.Fatal("Error reading the email template:", err)
		return
	}

	bukti := struct {
		Tanggal    string
		Waktu      string
		Kode       string
		Nama       string
		Instansi   *string
		Nomor      string
		Pengembang string
	}{
		Tanggal:    tanggal,
		Waktu:      nowTime.Format("15:04:05"),
		Kode:       kode,
		Nama:       currentUser.Name,
		Instansi:   payload.Andalalin.NamaPerusahaan,
		Nomor:      payload.Andalalin.NomerPemohon,
		Pengembang: payload.Andalalin.NamaPengembang,
	}

	buffer := new(bytes.Buffer)
	if err = t.Execute(buffer, bukti); err != nil {
		log.Fatal("Eror saat membaca template:", err)
		return
	}

	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		log.Fatal("Eror generate pdf", err)
		return
	}

	// read the HTML page as a PDF page
	page := wkhtmltopdf.NewPageReader(bytes.NewReader(buffer.Bytes()))

	pdfg.AddPage(page)

	pdfg.Dpi.Set(300)
	pdfg.PageSize.Set(wkhtmltopdf.PageSizeA4)
	pdfg.Orientation.Set(wkhtmltopdf.OrientationPortrait)
	pdfg.MarginBottom.Set(20)
	pdfg.MarginLeft.Set(30)
	pdfg.MarginRight.Set(30)
	pdfg.MarginTop.Set(20)

	err = pdfg.Create()
	if err != nil {
		log.Fatal(err)
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	persyaratan := []models.PersyaratanPermohonan{}

	for key, files := range form.File {
		for _, file := range files {
			// Save the uploaded file with key as prefix
			filed, err := file.Open()

			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			defer filed.Close()

			data, err := io.ReadAll(filed)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			// Store the blob data in the map
			persyaratan = append(persyaratan, models.PersyaratanPermohonan{Persyaratan: key, Tipe: "Pdf", Berkas: data})

		}
	}

	dokumen := []models.DokumenPermohonan{}

	dokumen = append(dokumen, models.DokumenPermohonan{Role: "User", Dokumen: "Tanda terima pendaftaran", Tipe: "Pdf", Berkas: pdfg.Bytes()})

	permohonan := models.Andalalin{
		//Data Permohonan
		IdUser:            currentUser.ID,
		JenisAndalalin:    "Dokumen analisis dampak lalu lintas",
		Bangkitan:         payload.Andalalin.Bangkitan,
		Pemohon:           payload.Andalalin.Pemohon,
		Kategori:          payload.Andalalin.KategoriJenisRencanaPembangunan,
		Jenis:             payload.Andalalin.JenisRencanaPembangunan,
		Kode:              kode,
		LokasiPengambilan: payload.Andalalin.LokasiPengambilan,
		WaktuAndalalin:    nowTime.Format("15:04:05"),
		TanggalAndalalin:  tanggal,
		StatusAndalalin:   "Cek persyaratan",

		//Data Proyek
		NamaProyek:      payload.Andalalin.NamaProyek,
		JenisProyek:     payload.Andalalin.JenisProyek,
		NegaraProyek:    "Indonesia",
		ProvinsiProyek:  payload.Andalalin.ProvinsiProyek,
		KabupatenProyek: payload.Andalalin.KabupatenProyek,
		KecamatanProyek: payload.Andalalin.KecamatanProyek,
		KelurahanProyek: payload.Andalalin.KelurahanProyek,
		AlamatProyek:    payload.Andalalin.AlamatProyek,
		KodeJalan:       payload.Andalalin.KodeJalan,
		KodeJalanMerge:  payload.Andalalin.KodeJalanMerge,
		NamaJalan:       payload.Andalalin.NamaJalan,
		PangkalJalan:    payload.Andalalin.PangkalJalan,
		UjungJalan:      payload.Andalalin.UjungJalan,
		PanjangJalan:    payload.Andalalin.PanjangJalan,
		LebarJalan:      payload.Andalalin.LebarJalan,
		PermukaanJalan:  payload.Andalalin.PermukaanJalan,
		FungsiJalan:     payload.Andalalin.FungsiJalan,

		//Data Pemohon
		NikPemohon:             payload.Andalalin.NikPemohon,
		NamaPemohon:            currentUser.Name,
		EmailPemohon:           currentUser.Email,
		TempatLahirPemohon:     payload.Andalalin.TempatLahirPemohon,
		TanggalLahirPemohon:    payload.Andalalin.TanggalLahirPemohon,
		NegaraPemohon:          "Indonesia",
		ProvinsiPemohon:        payload.Andalalin.ProvinsiPemohon,
		KabupatenPemohon:       payload.Andalalin.KabupatenPemohon,
		KecamatanPemohon:       payload.Andalalin.KecamatanPemohon,
		KelurahanPemohon:       payload.Andalalin.KelurahanPemohon,
		AlamatPemohon:          payload.Andalalin.AlamatPemohon,
		JenisKelaminPemohon:    payload.Andalalin.JenisKelaminPemohon,
		NomerPemohon:           payload.Andalalin.NomerPemohon,
		JabatanPemohon:         payload.Andalalin.JabatanPemohon,
		NomerSertifikatPemohon: payload.Andalalin.NomerSertifikatPemohon,
		KlasifikasiPemohon:     payload.Andalalin.KlasifikasiPemohon,

		//Data Perusahaan
		NamaPerusahaan:              payload.Andalalin.NamaPerusahaan,
		NegaraPerusahaan:            "Indonesia",
		ProvinsiPerusahaan:          payload.Andalalin.ProvinsiPerusahaan,
		KabupatenPerusahaan:         payload.Andalalin.KabupatenPerusahaan,
		KecamatanPerusahaan:         payload.Andalalin.KecamatanPerusahaan,
		KelurahanPerusahaan:         payload.Andalalin.KelurahanPerusahaan,
		AlamatPerusahaan:            payload.Andalalin.AlamatPerusahaan,
		NomerPerusahaan:             payload.Andalalin.NomerPerusahaan,
		EmailPerusahaan:             payload.Andalalin.EmailPerusahaan,
		NamaPimpinan:                payload.Andalalin.NamaPimpinan,
		JabatanPimpinan:             payload.Andalalin.JabatanPimpinan,
		JenisKelaminPimpinan:        payload.Andalalin.JenisKelaminPimpinan,
		NegaraPimpinanPerusahaan:    "Indonesia",
		ProvinsiPimpinanPerusahaan:  payload.Andalalin.ProvinsiPimpinanPerusahaan,
		KabupatenPimpinanPerusahaan: payload.Andalalin.KabupatenPimpinanPerusahaan,
		KecamatanPimpinanPerusahaan: payload.Andalalin.KecamatanPimpinanPerusahaan,
		KelurahanPimpinanPerusahaan: payload.Andalalin.KelurahanPimpinanPerusahaan,
		AlamatPimpinan:              payload.Andalalin.AlamatPimpinan,

		//Data Pengembang
		NamaPengembang:                 payload.Andalalin.NamaPengembang,
		NegaraPengembang:               "Indonesia",
		ProvinsiPengembang:             payload.Andalalin.ProvinsiPengembang,
		KabupatenPengembang:            payload.Andalalin.KabupatenPengembang,
		KecamatanPengembang:            payload.Andalalin.KecamatanPengembang,
		KelurahanPengembang:            payload.Andalalin.KelurahanPengembang,
		AlamatPengembang:               payload.Andalalin.AlamatPengembang,
		NomerPengembang:                payload.Andalalin.NomerPengembang,
		EmailPengembang:                payload.Andalalin.EmailPengembang,
		NamaPimpinanPengembang:         payload.Andalalin.NamaPimpinanPengembang,
		JabatanPimpinanPengembang:      payload.Andalalin.JabatanPimpinanPengembang,
		JenisKelaminPimpinanPengembang: payload.Andalalin.JenisKelaminPimpinanPengembang,
		NegaraPimpinanPengembang:       "Indonesia",
		ProvinsiPimpinanPengembang:     payload.Andalalin.ProvinsiPimpinanPengembang,
		KabupatenPimpinanPengembang:    payload.Andalalin.KabupatenPimpinanPengembang,
		KecamatanPimpinanPengembang:    payload.Andalalin.KecamatanPimpinanPengembang,
		KelurahanPimpinanPengembang:    payload.Andalalin.KelurahanPimpinanPengembang,
		AlamatPimpinanPengembang:       payload.Andalalin.AlamatPimpinanPengembang,

		//Data Kegiatan
		Aktivitas:         payload.Andalalin.Aktivitas,
		Peruntukan:        payload.Andalalin.Peruntukan,
		TotalLuasLahan:    payload.Andalalin.TotalLuasLahan,
		KriteriaKhusus:    payload.Andalalin.KriteriaKhusus,
		NilaiKriteria:     payload.Andalalin.NilaiKriteria,
		Terbilang:         payload.Andalalin.Terbilang,
		LokasiBangunan:    payload.Andalalin.LokasiBangunan,
		LatitudeBangunan:  payload.Andalalin.LatitudeBangunan,
		LongitudeBangunan: payload.Andalalin.LongitudeBangunan,
		NomerSKRK:         payload.Andalalin.NomerSKRK,
		TanggalSKRK:       payload.Andalalin.TanggalSKRK,
		Catatan:           payload.Andalalin.Catatan,

		//Dokumen Permohonan
		Dokumen: dokumen,

		//Data Persyaratan
		Persyaratan: persyaratan,
	}

	result := ac.DB.Create(&permohonan)

	respone := &models.DaftarAndalalinResponse{
		IdAndalalin:      permohonan.IdAndalalin,
		Kode:             permohonan.Kode,
		TanggalAndalalin: permohonan.TanggalAndalalin,
		Nama:             permohonan.NamaPemohon,
		Pengembang:       permohonan.NamaPengembang,
		JenisAndalalin:   permohonan.JenisAndalalin,
		StatusAndalalin:  permohonan.StatusAndalalin,
	}

	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "eror saat mengirim data"})
		return
	} else {
		ac.ReleaseTicketLevel1(ctx, permohonan.IdAndalalin)

		var user []models.User

		ac.DB.Find(&user, "role = ?", "Operator")

		for _, users := range user {
			simpanNotif := models.Notifikasi{
				IdUser: users.ID,
				Title:  "Permohonan baru",
				Body:   "Permohonan baru dengan kode " + permohonan.Kode + " telah diajukan, silahkan menindaklanjuti permohonan",
			}

			ac.DB.Create(&simpanNotif)

			if users.PushToken != "" {
				notif := utils.Notification{
					IdUser: users.ID,
					Title:  "Permohonan baru",
					Body:   "Permohonan baru dengan kode " + permohonan.Kode + " telah diajukan, silahkan menindaklanjuti permohonan",
					Token:  users.PushToken,
				}

				utils.SendPushNotifications(&notif)
			}
		}

		ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": respone})
	}
}

func (ac *AndalalinController) PengajuanPerlalin(ctx *gin.Context) {
	var payload *models.DataPerlalin
	currentUser := ctx.MustGet("currentUser").(models.User)

	config, _ := initializers.LoadConfig()

	accessUser := ctx.MustGet("accessUser").(string)

	claim, error := utils.ValidateToken(accessUser, config.AccessTokenPublicKey)
	if error != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": error.Error()})
		return
	}

	credential := claim.Credentials[repository.AndalalinPengajuanCredential]

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

	kode := "perlalin/" + utils.Generate(6)
	tanggal := nowTime.Format("02") + " " + utils.Bulan(nowTime.Month()) + " " + nowTime.Format("2006")

	t, err := template.ParseFiles("templates/tandaterimaPerlalin.html")
	if err != nil {
		log.Fatal("Error reading the email template:", err)
		return
	}

	bukti := struct {
		Tanggal string
		Waktu   string
		Kode    string
		Nama    string
		Nomor   string
	}{
		Tanggal: tanggal,
		Waktu:   nowTime.Format("15:04:05"),
		Kode:    kode,
		Nama:    currentUser.Name,
		Nomor:   payload.Perlalin.NomerPemohon,
	}

	buffer := new(bytes.Buffer)
	if err = t.Execute(buffer, bukti); err != nil {
		log.Fatal("Eror saat membaca template:", err)
		return
	}

	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		log.Fatal("Eror generate pdf", err)
		return
	}

	// read the HTML page as a PDF page
	page := wkhtmltopdf.NewPageReader(bytes.NewReader(buffer.Bytes()))

	pdfg.AddPage(page)

	pdfg.Dpi.Set(300)
	pdfg.PageSize.Set(wkhtmltopdf.PageSizeA4)
	pdfg.Orientation.Set(wkhtmltopdf.OrientationPortrait)
	pdfg.MarginBottom.Set(20)
	pdfg.MarginLeft.Set(30)
	pdfg.MarginRight.Set(30)
	pdfg.MarginTop.Set(20)

	err = pdfg.Create()
	if err != nil {
		log.Fatal(err)
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dokumen := []models.PersyaratanPermohonan{}

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
			dokumen = append(dokumen, models.PersyaratanPermohonan{Persyaratan: key, Berkas: data})

		}
	}

	permohonan := models.Perlalin{
		IdUser:                      currentUser.ID,
		JenisAndalalin:              "Perlengkapan lalu lintas",
		Kategori:                    payload.Perlalin.Kategori,
		Jenis:                       payload.Perlalin.Jenis,
		Kode:                        kode,
		NikPemohon:                  payload.Perlalin.NikPemohon,
		NamaPemohon:                 currentUser.Name,
		EmailPemohon:                currentUser.Email,
		TempatLahirPemohon:          payload.Perlalin.TempatLahirPemohon,
		TanggalLahirPemohon:         payload.Perlalin.TanggalLahirPemohon,
		WilayahAdministratifPemohon: payload.Perlalin.WilayahAdministratifPemohon,
		AlamatPemohon:               payload.Perlalin.AlamatPemohon,
		JenisKelaminPemohon:         payload.Perlalin.JenisKelaminPemohon,
		NomerPemohon:                payload.Perlalin.NomerPemohon,
		LokasiPengambilan:           payload.Perlalin.LokasiPengambilan,
		WaktuAndalalin:              nowTime.Format("15:04:05"),
		TanggalAndalalin:            tanggal,
		Alasan:                      payload.Perlalin.Alasan,
		Peruntukan:                  payload.Perlalin.Peruntukan,
		LokasiPemasangan:            payload.Perlalin.LokasiPemasangan,
		LatitudePemasangan:          payload.Perlalin.LatitudePemasangan,
		LongitudePemasangan:         payload.Perlalin.LongitudePemasangan,
		Catatan:                     payload.Perlalin.Catatan,
		StatusAndalalin:             "Cek persyaratan",
		TandaTerimaPendaftaran:      pdfg.Bytes(),

		Persyaratan: dokumen,
	}

	result := ac.DB.Create(&permohonan)

	respone := &models.DaftarAndalalinResponse{
		IdAndalalin:      permohonan.IdAndalalin,
		Kode:             permohonan.Kode,
		TanggalAndalalin: permohonan.TanggalAndalalin,
		Nama:             permohonan.NamaPemohon,
		Email:            permohonan.EmailPemohon,
		Petugas:          permohonan.NamaPetugas,
		JenisAndalalin:   permohonan.JenisAndalalin,
		StatusAndalalin:  permohonan.StatusAndalalin,
	}

	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "eror saat mengirim data"})
		return
	} else {
		ac.ReleaseTicketLevel1(ctx, permohonan.IdAndalalin)
		var user []models.User

		ac.DB.Find(&user, "role = ?", "Operator")

		for _, users := range user {
			simpanNotif := models.Notifikasi{
				IdUser: users.ID,
				Title:  "Permohonan baru",
				Body:   "Permohonan baru dengan kode " + permohonan.Kode + " telah diajukan, silahkan menindaklanjuti permohonan",
			}

			ac.DB.Create(&simpanNotif)

			if users.PushToken != "" {
				notif := utils.Notification{
					IdUser: users.ID,
					Title:  "Permohonan baru",
					Body:   "Permohonan baru dengan kode " + permohonan.Kode + " telah diajukan, silahkan menindaklanjuti permohonan",
					Token:  users.PushToken,
				}

				utils.SendPushNotifications(&notif)
			}
		}
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": respone})
	}
}

func (ac *AndalalinController) ReleaseTicketLevel1(ctx *gin.Context, id uuid.UUID) {
	tiket := models.TiketLevel1{
		IdAndalalin: id,
		Status:      "Buka",
	}

	result := ac.DB.Create(&tiket)

	if result.Error != nil && strings.Contains(result.Error.Error(), "duplicate key value violates unique") {
		ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "Tiket level 1 sudah tersedia"})
		return
	} else if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Telah terjadi sesuatu"})
		return
	}
}

func (ac *AndalalinController) CloseTiketLevel1(ctx *gin.Context, id uuid.UUID) {
	var tiket models.TiketLevel1

	result := ac.DB.Model(&tiket).Where("id_andalalin = ? AND status = ?", id, "Buka").Update("status", "Tutup")
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Tiket level 1 tidak tersedia"})
		return
	}
}

func (ac *AndalalinController) ReleaseTicketLevel2(ctx *gin.Context, id uuid.UUID, petugas uuid.UUID) {
	var tiket1 models.TiketLevel1
	results := ac.DB.First(&tiket1, "id_andalalin = ?", id)

	if results.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": results.Error})
		return
	}

	tiket := models.TiketLevel2{
		IdTiketLevel1: tiket1.IdTiketLevel1,
		IdAndalalin:   id,
		IdPetugas:     petugas,
		Status:        "Buka",
	}

	result := ac.DB.Create(&tiket)

	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Telah terjadi sesuatu"})
		return
	}
}

func (ac *AndalalinController) CloseTiketLevel2(ctx *gin.Context, id uuid.UUID) {
	var tiket models.TiketLevel2

	result := ac.DB.Model(&tiket).Where("id_andalalin = ? AND status = ?", id, "Buka").Or("id_andalalin = ? AND status = ?", id, "Batal").Update("status", "Tutup")
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Telah terjadi sesuatu"})
		return
	}
}

func (ac *AndalalinController) TundaPermohonan(ctx *gin.Context) {
	id := ctx.Param("id_andalalin")
	pertimbangan := ctx.Param("pertimbangan")

	config, _ := initializers.LoadConfig()

	accessUser := ctx.MustGet("accessUser").(string)

	claim, error := utils.ValidateToken(accessUser, config.AccessTokenPublicKey)
	if error != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": error.Error()})
		return
	}

	credential := claim.Credentials[repository.AndalalinTindakLanjut]

	if !credential {
		// Return status 403 and permission denied error message.
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": true,
			"msg":   "Permission denied",
		})
		return
	}

	var perlalin models.Perlalin

	ac.DB.First(&perlalin, "id_andalalin = ?", id)

	if perlalin.IdAndalalin != uuid.Nil {
		ac.CloseTiketLevel1(ctx, perlalin.IdAndalalin)
		perlalin.StatusAndalalin = "Permohonan ditunda"
		perlalin.PertimbanganPenundaan = pertimbangan
		ac.DB.Save(&perlalin)

		data := utils.PermohonanDitolak{
			Kode:    perlalin.Kode,
			Nama:    perlalin.NamaPemohon,
			Tlp:     perlalin.NomerPemohon,
			Jenis:   perlalin.JenisAndalalin,
			Status:  perlalin.StatusAndalalin,
			Subject: "Permohonan ditunda",
		}

		utils.SendEmailPermohonanDitunda(perlalin.EmailPemohon, &data)

		var user models.User
		resultUser := ac.DB.First(&user, "id = ?", perlalin.IdUser)
		if resultUser.Error != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "User tidak ditemukan"})
			return
		}

		simpanNotif := models.Notifikasi{
			IdUser: user.ID,
			Title:  "Permohonan ditunda",
			Body:   "Permohonan anda dengan kode " + perlalin.Kode + " telah ditunda, silahkan cek permohonan pada aplikasi untuk lebih jelas",
		}

		ac.DB.Create(&simpanNotif)

		if user.PushToken != "" {
			notif := utils.Notification{
				IdUser: user.ID,
				Title:  "Permohonan ditunda",
				Body:   "Permohonan anda dengan kode " + perlalin.Kode + " telah ditunda, silahkan cek permohonan pada aplikasi untuk lebih jelas",
				Token:  user.PushToken,
			}

			utils.SendPushNotifications(&notif)
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (ac *AndalalinController) LanjutkanPermohonan(ctx *gin.Context) {
	id := ctx.Param("id_andalalin")

	config, _ := initializers.LoadConfig()

	accessUser := ctx.MustGet("accessUser").(string)

	claim, error := utils.ValidateToken(accessUser, config.AccessTokenPublicKey)
	if error != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": error.Error()})
		return
	}

	credential := claim.Credentials[repository.AndalalinTindakLanjut]

	if !credential {
		// Return status 403 and permission denied error message.
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": true,
			"msg":   "Permission denied",
		})
		return
	}

	var perlalin models.Perlalin

	ac.DB.First(&perlalin, "id_andalalin = ?", id)

	if perlalin.IdAndalalin != uuid.Nil {
		ac.CloseTiketLevel1(ctx, perlalin.IdAndalalin)
		perlalin.StatusAndalalin = "Cek persyaratan"
		perlalin.PertimbanganPenundaan = ""
		ac.DB.Save(&perlalin)

		var user models.User
		resultUser := ac.DB.First(&user, "id = ?", perlalin.IdUser)
		if resultUser.Error != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "User tidak ditemukan"})
			return
		}

		simpanNotif := models.Notifikasi{
			IdUser: user.ID,
			Title:  "Permohonan dilanjutkan",
			Body:   "Permohonan anda dengan kode " + perlalin.Kode + " telah dilanjutkan, silahkan cek permohonan pada aplikasi untuk lebih jelas",
		}

		ac.DB.Create(&simpanNotif)

		if user.PushToken != "" {
			notif := utils.Notification{
				IdUser: user.ID,
				Title:  "Permohonan dilanjutkan",
				Body:   "Permohonan anda dengan kode " + perlalin.Kode + " telah dilanjutkan, silahkan cek permohonan pada aplikasi untuk lebih jelas",
				Token:  user.PushToken,
			}

			utils.SendPushNotifications(&notif)
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (ac *AndalalinController) TolakPermohonan(ctx *gin.Context) {
	id := ctx.Param("id_andalalin")
	pertimbangan := ctx.Param("pertimbangan")

	config, _ := initializers.LoadConfig()

	accessUser := ctx.MustGet("accessUser").(string)

	claim, error := utils.ValidateToken(accessUser, config.AccessTokenPublicKey)
	if error != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": error.Error()})
		return
	}

	credential := claim.Credentials[repository.AndalalinTindakLanjut]

	if !credential {
		// Return status 403 and permission denied error message.
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": true,
			"msg":   "Permission denied",
		})
		return
	}

	var andalalin models.Andalalin
	var perlalin models.Perlalin

	ac.DB.First(&andalalin, "id_andalalin = ?", id)
	ac.DB.First(&perlalin, "id_andalalin = ?", id)

	if andalalin.IdAndalalin != uuid.Nil {
		ac.CloseTiketLevel1(ctx, andalalin.IdAndalalin)
		andalalin.StatusAndalalin = "Permohonan ditolak"
		andalalin.Pertimbangan = pertimbangan
		ac.DB.Save(&andalalin)

		data := utils.PermohonanDitolak{
			Kode:    andalalin.Kode,
			Nama:    andalalin.NamaPemohon,
			Tlp:     andalalin.NomerPemohon,
			Jenis:   andalalin.JenisAndalalin,
			Status:  andalalin.StatusAndalalin,
			Subject: "Permohonan ditolak",
		}

		utils.SendEmailPermohonanDitolak(andalalin.EmailPemohon, &data)

		var user models.User
		resultUser := ac.DB.First(&user, "id = ?", andalalin.IdUser)
		if resultUser.Error != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "User tidak ditemukan"})
			return
		}

		simpanNotif := models.Notifikasi{
			IdUser: user.ID,
			Title:  "Permohonan ditolak",
			Body:   "Permohonan anda dengan kode " + andalalin.Kode + " telah ditolak, silahkan cek permohonan pada aplikasi untuk lebih jelas",
		}

		ac.DB.Create(&simpanNotif)

		if user.PushToken != "" {
			notif := utils.Notification{
				IdUser: user.ID,
				Title:  "Permohonan ditolak",
				Body:   "Permohonan anda dengan kode " + andalalin.Kode + " telah ditolak, silahkan cek permohonan pada aplikasi untuk lebih jelas",
				Token:  user.PushToken,
			}

			utils.SendPushNotifications(&notif)
		}
	}

	if perlalin.IdAndalalin != uuid.Nil {
		ac.CloseTiketLevel1(ctx, perlalin.IdAndalalin)
		perlalin.StatusAndalalin = "Permohonan ditolak"
		perlalin.PertimbanganPenolakan = pertimbangan
		ac.DB.Save(&perlalin)

		data := utils.PermohonanDitolak{
			Kode:    perlalin.Kode,
			Nama:    perlalin.NamaPemohon,
			Tlp:     perlalin.NomerPemohon,
			Jenis:   perlalin.JenisAndalalin,
			Status:  perlalin.StatusAndalalin,
			Subject: "Permohonan ditolak",
		}

		utils.SendEmailPermohonanDitolak(perlalin.EmailPemohon, &data)

		var user models.User
		resultUser := ac.DB.First(&user, "id = ?", perlalin.IdUser)
		if resultUser.Error != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "User tidak ditemukan"})
			return
		}

		simpanNotif := models.Notifikasi{
			IdUser: user.ID,
			Title:  "Permohonan ditolak",
			Body:   "Permohonan anda dengan kode " + perlalin.Kode + " telah ditolak, silahkan cek permohonan pada aplikasi untuk lebih jelas",
		}

		ac.DB.Create(&simpanNotif)

		if user.PushToken != "" {
			notif := utils.Notification{
				IdUser: user.ID,
				Title:  "Permohonan ditolak",
				Body:   "Permohonan anda dengan kode " + perlalin.Kode + " telah ditolak, silahkan cek permohonan pada aplikasi untuk lebih jelas",
				Token:  user.PushToken,
			}

			utils.SendPushNotifications(&notif)
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (ac *AndalalinController) GetPermohonanByIdUser(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)

	var andalalin []models.Andalalin
	var perlalin []models.Perlalin

	resultsAndalalin := ac.DB.Order("tanggal_andalalin").Find(&andalalin, "id_user = ?", currentUser.ID)
	resultsPerlalin := ac.DB.Order("tanggal_andalalin").Find(&perlalin, "id_user = ?", currentUser.ID)

	if resultsAndalalin.Error != nil && resultsPerlalin != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Tidak ditemukan"})
		return
	} else {
		var respone []models.DaftarAndalalinResponse
		for _, s := range andalalin {
			respone = append(respone, models.DaftarAndalalinResponse{
				IdAndalalin:      s.IdAndalalin,
				Kode:             s.Kode,
				TanggalAndalalin: s.TanggalAndalalin,
				Nama:             s.NamaPemohon,
				Pengembang:       s.NamaPengembang,
				JenisAndalalin:   s.JenisAndalalin,
				StatusAndalalin:  s.StatusAndalalin,
			})
		}
		for _, s := range perlalin {
			respone = append(respone, models.DaftarAndalalinResponse{
				IdAndalalin:      s.IdAndalalin,
				Kode:             s.Kode,
				TanggalAndalalin: s.TanggalAndalalin,
				Nama:             s.NamaPemohon,
				Email:            s.EmailPemohon,
				Petugas:          s.NamaPetugas,
				JenisAndalalin:   s.JenisAndalalin,
				StatusAndalalin:  s.StatusAndalalin,
			})
		}
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "results": len(respone), "data": respone})
	}
}

func (ac *AndalalinController) GetPermohonanByIdAndalalin(ctx *gin.Context) {
	id := ctx.Param("id_andalalin")

	currentUser := ctx.MustGet("currentUser").(models.User)

	var andalalin models.Andalalin
	var perlalin models.Perlalin

	ac.DB.First(&andalalin, "id_andalalin = ?", id)
	ac.DB.First(&perlalin, "id_andalalin = ?", id)

	var ticket2 models.TiketLevel2
	resultTiket2 := ac.DB.Not("status = ?", "Tutup").Where("id_andalalin = ?", id).First(&ticket2)
	var status string
	if resultTiket2.Error != nil {
		status = "Kosong"
	} else {
		status = ticket2.Status
	}

	var persyaratan_andalalin []string
	var persyaratan_perlalin []string

	for _, persyaratan := range andalalin.Persyaratan {
		persyaratan_andalalin = append(persyaratan_andalalin, persyaratan.Persyaratan)
	}

	for _, persyaratan := range perlalin.Persyaratan {
		persyaratan_perlalin = append(persyaratan_perlalin, persyaratan.Persyaratan)
	}

	var dokumen_andalalin_dinas []string

	for _, dokumen := range andalalin.Dokumen {
		dokumen_andalalin_dinas = append(dokumen_andalalin_dinas, dokumen.Dokumen)
	}

	var dokumen_andalalin_user []string

	for _, dokumen := range andalalin.Dokumen {
		if dokumen.Role == "User" {
			dokumen_andalalin_user = append(dokumen_andalalin_user, dokumen.Dokumen)
		}
	}

	if andalalin.IdAndalalin != uuid.Nil {
		if currentUser.Role == "User" {
			dataUser := models.AndalalinResponseUser{
				//Data Permohonan
				IdAndalalin:             andalalin.IdAndalalin,
				JenisAndalalin:          andalalin.JenisAndalalin,
				Bangkitan:               andalalin.Bangkitan,
				Pemohon:                 andalalin.Pemohon,
				JenisRencanaPembangunan: andalalin.Jenis,
				Kategori:                andalalin.Kategori,
				Kode:                    andalalin.Kode,
				LokasiPengambilan:       andalalin.LokasiPengambilan,
				WaktuAndalalin:          andalalin.WaktuAndalalin,
				TanggalAndalalin:        andalalin.TanggalAndalalin,
				StatusAndalalin:         andalalin.StatusAndalalin,

				//Data Proyek
				NamaProyek:      andalalin.NamaProyek,
				JenisProyek:     andalalin.JenisProyek,
				NamaJalan:       andalalin.NamaJalan,
				FungsiJalan:     andalalin.FungsiJalan,
				NegaraProyek:    andalalin.NegaraProyek,
				ProvinsiProyek:  andalalin.ProvinsiProyek,
				KabupatenProyek: andalalin.KabupatenProyek,
				KecamatanProyek: andalalin.KecamatanProyek,
				KelurahanProyek: andalalin.KelurahanProyek,
				AlamatProyek:    andalalin.AlamatProyek,

				//Data Pemohon
				NikPemohon:             andalalin.NikPemohon,
				EmailPemohon:           andalalin.EmailPemohon,
				NomerPemohon:           andalalin.NomerPemohon,
				NamaPemohon:            andalalin.NamaPemohon,
				JabatanPemohon:         andalalin.JabatanPemohon,
				NomerSertifikatPemohon: andalalin.NomerSertifikatPemohon,
				KlasifikasiPemohon:     andalalin.KlasifikasiPemohon,

				//Data perusahaan
				NamaPerusahaan: andalalin.NamaPerusahaan,

				//Data Pengembang
				NamaPengembang: andalalin.NamaPengembang,

				//Data Kegiatan
				Aktivitas:         andalalin.Aktivitas,
				Peruntukan:        andalalin.Peruntukan,
				TotalLuasLahan:    andalalin.TotalLuasLahan,
				LokasiBangunan:    andalalin.LokasiBangunan,
				LatitudeBangunan:  andalalin.LatitudeBangunan,
				LongitudeBangunan: andalalin.LongitudeBangunan,
				KriteriaKhusus:    andalalin.KriteriaKhusus,
				NilaiKriteria:     andalalin.NilaiKriteria,
				Catatan:           andalalin.Catatan,

				//Data Persyaratan dan Pertimbangan
				PersyaratanTidakSesuai: andalalin.PersyaratanTidakSesuai,
				Pertimbangan:           andalalin.Pertimbangan,

				Persyaratan: persyaratan_andalalin,

				//Dokumen Permohonan
				Dokumen: dokumen_andalalin_user,
			}

			ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": dataUser})
		} else {
			data := models.AndalalinResponse{
				//Data Permohonan
				IdAndalalin:       andalalin.IdAndalalin,
				JenisAndalalin:    andalalin.JenisAndalalin,
				Bangkitan:         andalalin.Bangkitan,
				Pemohon:           andalalin.Pemohon,
				Kategori:          andalalin.Kategori,
				Jenis:             andalalin.Jenis,
				Kode:              andalalin.Kode,
				LokasiPengambilan: andalalin.LokasiPengambilan,
				WaktuAndalalin:    andalalin.WaktuAndalalin,
				TanggalAndalalin:  andalalin.TanggalAndalalin,
				StatusAndalalin:   andalalin.StatusAndalalin,

				//Data Proyek
				NamaProyek:      andalalin.NamaProyek,
				JenisProyek:     andalalin.JenisProyek,
				NegaraProyek:    andalalin.NegaraProyek,
				ProvinsiProyek:  andalalin.ProvinsiProyek,
				KabupatenProyek: andalalin.KabupatenProyek,
				KecamatanProyek: andalalin.KecamatanProyek,
				KelurahanProyek: andalalin.KelurahanProyek,
				AlamatProyek:    andalalin.AlamatProyek,
				KodeJalan:       andalalin.KodeJalan,
				KodeJalanMerge:  andalalin.KodeJalanMerge,
				NamaJalan:       andalalin.NamaJalan,
				PangkalJalan:    andalalin.PangkalJalan,
				UjungJalan:      andalalin.UjungJalan,
				PanjangJalan:    andalalin.PanjangJalan,
				LebarJalan:      andalalin.LebarJalan,
				PermukaanJalan:  andalalin.PermukaanJalan,
				FungsiJalan:     andalalin.FungsiJalan,

				//Data Pemohon
				NikPemohon:             andalalin.NikPemohon,
				NamaPemohon:            andalalin.NamaPemohon,
				EmailPemohon:           andalalin.EmailPemohon,
				TempatLahirPemohon:     andalalin.TempatLahirPemohon,
				TanggalLahirPemohon:    andalalin.TanggalLahirPemohon,
				NegaraPemohon:          andalalin.NegaraPemohon,
				ProvinsiPemohon:        andalalin.ProvinsiPemohon,
				KabupatenPemohon:       andalalin.KabupatenPemohon,
				KecamatanPemohon:       andalalin.KecamatanPemohon,
				KelurahanPemohon:       andalalin.KelurahanPemohon,
				AlamatPemohon:          andalalin.AlamatPemohon,
				JenisKelaminPemohon:    andalalin.JenisKelaminPemohon,
				NomerPemohon:           andalalin.NomerPemohon,
				JabatanPemohon:         andalalin.JabatanPemohon,
				NomerSertifikatPemohon: andalalin.NomerSertifikatPemohon,
				KlasifikasiPemohon:     andalalin.KlasifikasiPemohon,

				//Data Perusahaan
				NamaPerusahaan:              andalalin.NamaPerusahaan,
				NegaraPerusahaan:            andalalin.NegaraPerusahaan,
				ProvinsiPerusahaan:          andalalin.ProvinsiPerusahaan,
				KabupatenPerusahaan:         andalalin.KabupatenPerusahaan,
				KecamatanPerusahaan:         andalalin.KecamatanPerusahaan,
				KelurahanPerusahaan:         andalalin.KelurahanPerusahaan,
				AlamatPerusahaan:            andalalin.AlamatPerusahaan,
				NomerPerusahaan:             andalalin.NomerPerusahaan,
				EmailPerusahaan:             andalalin.EmailPerusahaan,
				NamaPimpinan:                andalalin.NamaPimpinan,
				JabatanPimpinan:             andalalin.JabatanPimpinan,
				JenisKelaminPimpinan:        andalalin.JenisKelaminPimpinan,
				NegaraPimpinanPerusahaan:    andalalin.NegaraPimpinanPerusahaan,
				ProvinsiPimpinanPerusahaan:  andalalin.ProvinsiPimpinanPerusahaan,
				KabupatenPimpinanPerusahaan: andalalin.KabupatenPimpinanPerusahaan,
				KecamatanPimpinanPerusahaan: andalalin.KecamatanPimpinanPerusahaan,
				KelurahanPimpinanPerusahaan: andalalin.KelurahanPimpinanPerusahaan,
				AlamatPimpinan:              andalalin.AlamatPimpinan,

				//Data Pengembang
				NamaPengembang:                 andalalin.NamaPengembang,
				NegaraPengembang:               andalalin.NegaraPengembang,
				ProvinsiPengembang:             andalalin.ProvinsiPengembang,
				KabupatenPengembang:            andalalin.KabupatenPengembang,
				KecamatanPengembang:            andalalin.KecamatanPengembang,
				KelurahanPengembang:            andalalin.KelurahanPengembang,
				AlamatPengembang:               andalalin.AlamatPengembang,
				NomerPengembang:                andalalin.NomerPengembang,
				EmailPengembang:                andalalin.EmailPengembang,
				NamaPimpinanPengembang:         andalalin.NamaPimpinanPengembang,
				JabatanPimpinanPengembang:      andalalin.JabatanPimpinanPengembang,
				JenisKelaminPimpinanPengembang: andalalin.JenisKelaminPimpinanPengembang,
				NegaraPimpinanPengembang:       andalalin.NegaraPimpinanPengembang,
				ProvinsiPimpinanPengembang:     andalalin.ProvinsiPimpinanPengembang,
				KabupatenPimpinanPengembang:    andalalin.KabupatenPimpinanPengembang,
				KecamatanPimpinanPengembang:    andalalin.KecamatanPimpinanPengembang,
				KelurahanPimpinanPengembang:    andalalin.KelurahanPimpinanPengembang,
				AlamatPimpinanPengembang:       andalalin.AlamatPimpinanPengembang,

				//Data Kegiatan
				Aktivitas:         andalalin.Aktivitas,
				Peruntukan:        andalalin.Peruntukan,
				TotalLuasLahan:    andalalin.TotalLuasLahan,
				LokasiBangunan:    andalalin.LokasiBangunan,
				LatitudeBangunan:  andalalin.LatitudeBangunan,
				LongitudeBangunan: andalalin.LongitudeBangunan,
				KriteriaKhusus:    andalalin.KriteriaKhusus,
				NilaiKriteria:     andalalin.NilaiKriteria,
				NomerSKRK:         andalalin.NomerSKRK,
				TanggalSKRK:       andalalin.TanggalSKRK,
				Catatan:           andalalin.Catatan,

				//Data Persyaratan
				Persyaratan:            persyaratan_andalalin,
				PersyaratanTidakSesuai: andalalin.PersyaratanTidakSesuai,

				//Dokumen Permohonan
				Dokumen: dokumen_andalalin_dinas,

				//Data Persertujuan
				PersetujuanDokumen:           andalalin.PersetujuanDokumen,
				KeteranganPersetujuanDokumen: andalalin.KeteranganPersetujuanDokumen,

				//Data Pertimbangan
				Pertimbangan: andalalin.Pertimbangan,
			}
			ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": data})
		}
	}

	if perlalin.IdAndalalin != uuid.Nil {
		if currentUser.Role == "User" {
			dataUser := models.PerlalinResponseUser{
				//Data Permohonan
				IdAndalalin:      perlalin.IdAndalalin,
				JenisAndalalin:   perlalin.JenisAndalalin,
				Kategori:         perlalin.Kategori,
				Jenis:            perlalin.Jenis,
				Kode:             perlalin.Kode,
				WaktuAndalalin:   perlalin.WaktuAndalalin,
				TanggalAndalalin: perlalin.TanggalAndalalin,
				StatusAndalalin:  perlalin.StatusAndalalin,

				//Data Pemohon
				NamaPemohon:       perlalin.NamaPemohon,
				NikPemohon:        perlalin.NikPemohon,
				EmailPemohon:      perlalin.EmailPemohon,
				NomerPemohon:      perlalin.NomerPemohon,
				LokasiPengambilan: perlalin.LokasiPengambilan,

				//Data Kegiatan
				Alasan:                 perlalin.Alasan,
				Peruntukan:             perlalin.Peruntukan,
				LokasiPemasangan:       perlalin.LokasiPemasangan,
				LatitudePemasangan:     perlalin.LatitudePemasangan,
				LongitudePemasangan:    perlalin.LongitudePemasangan,
				PersyaratanTidakSesuai: perlalin.PersyaratanTidakSesuai,

				//Catatan
				Catatan: perlalin.Catatan,

				Tindakan:              perlalin.Tindakan,
				PertimbanganTindakan:  perlalin.PertimbanganTindakan,
				PertimbanganPenolakan: perlalin.PertimbanganPenolakan,
				PertimbanganPenundaan: perlalin.PertimbanganPenundaan,
			}

			ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": dataUser})
		} else {
			data := models.PerlalinResponse{
				//Data Permohonan
				IdAndalalin:      perlalin.IdAndalalin,
				JenisAndalalin:   perlalin.JenisAndalalin,
				Kategori:         perlalin.Kategori,
				Jenis:            perlalin.Jenis,
				Kode:             perlalin.Kode,
				WaktuAndalalin:   perlalin.WaktuAndalalin,
				TanggalAndalalin: perlalin.TanggalAndalalin,
				StatusAndalalin:  perlalin.StatusAndalalin,

				//Data Pemohon
				NikPemohon:                  perlalin.NikPemohon,
				NamaPemohon:                 perlalin.NamaPemohon,
				EmailPemohon:                perlalin.EmailPemohon,
				TempatLahirPemohon:          perlalin.TempatLahirPemohon,
				TanggalLahirPemohon:         perlalin.TanggalLahirPemohon,
				WilayahAdministratifPemohon: perlalin.WilayahAdministratifPemohon,
				AlamatPemohon:               perlalin.AlamatPemohon,
				JenisKelaminPemohon:         perlalin.JenisKelaminPemohon,
				NomerPemohon:                perlalin.NomerPemohon,
				LokasiPengambilan:           perlalin.LokasiPengambilan,

				//Data Kegiatan
				Alasan:                 perlalin.Alasan,
				Peruntukan:             perlalin.Peruntukan,
				LokasiPemasangan:       perlalin.LokasiPemasangan,
				LatitudePemasangan:     perlalin.LatitudePemasangan,
				LongitudePemasangan:    perlalin.LongitudePemasangan,
				PersyaratanTidakSesuai: perlalin.PersyaratanTidakSesuai,
				IdPetugas:              perlalin.IdPetugas,
				NamaPetugas:            perlalin.NamaPetugas,
				EmailPetugas:           perlalin.EmailPetugas,
				StatusTiketLevel2:      status,

				Persyaratan: persyaratan_perlalin,

				Catatan: perlalin.Catatan,

				Tindakan:              perlalin.Tindakan,
				PertimbanganTindakan:  perlalin.PertimbanganTindakan,
				PertimbanganPenolakan: perlalin.PertimbanganPenolakan,
				PertimbanganPenundaan: perlalin.PertimbanganPenundaan,
			}
			ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": data})
		}
	}
}

func (ac *AndalalinController) GetPermohonanByStatus(ctx *gin.Context) {
	status := ctx.Param("status_andalalin")

	config, _ := initializers.LoadConfig()

	accessUser := ctx.MustGet("accessUser").(string)

	claim, error := utils.ValidateToken(accessUser, config.AccessTokenPublicKey)
	if error != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": error.Error()})
		return
	}

	credential := claim.Credentials[repository.AndalalinGetCredential]

	if !credential {
		// Return status 403 and permission denied error message.
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": true,
			"msg":   "Permission denied",
		})
		return
	}

	var andalalin []models.Andalalin
	var perlalin []models.Perlalin

	resultsAndalalin := ac.DB.Order("tanggal_andalalin").Find(&andalalin, "status_andalalin = ?", status)
	resultsPerlalin := ac.DB.Order("tanggal_andalalin").Find(&perlalin, "status_andalalin = ?", status)

	if resultsAndalalin.Error != nil && resultsPerlalin != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Tidak ditemukan"})
		return
	} else {
		var respone []models.DaftarAndalalinResponse
		for _, s := range andalalin {
			respone = append(respone, models.DaftarAndalalinResponse{
				IdAndalalin:      s.IdAndalalin,
				Kode:             s.Kode,
				TanggalAndalalin: s.TanggalAndalalin,
				Nama:             s.NamaPemohon,
				Pengembang:       s.NamaPengembang,
				JenisAndalalin:   s.JenisAndalalin,
				StatusAndalalin:  s.StatusAndalalin,
			})
		}
		for _, s := range perlalin {
			respone = append(respone, models.DaftarAndalalinResponse{
				IdAndalalin:      s.IdAndalalin,
				Kode:             s.Kode,
				TanggalAndalalin: s.TanggalAndalalin,
				Nama:             s.NamaPemohon,
				Email:            s.EmailPemohon,
				Petugas:          s.NamaPetugas,
				JenisAndalalin:   s.JenisAndalalin,
				StatusAndalalin:  s.StatusAndalalin,
			})
		}
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "results": len(respone), "data": respone})
	}
}

func (ac *AndalalinController) GetPermohonan(ctx *gin.Context) {
	config, _ := initializers.LoadConfig()

	accessUser := ctx.MustGet("accessUser").(string)

	claim, error := utils.ValidateToken(accessUser, config.AccessTokenPublicKey)
	if error != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": error.Error()})
		return
	}

	credential := claim.Credentials[repository.AndalalinGetCredential]

	if !credential {
		// Return status 403 and permission denied error message.
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": true,
			"msg":   "Permission denied",
		})
		return
	}

	var andalalin []models.Andalalin
	var perlalin []models.Perlalin

	resultsAndalalin := ac.DB.Order("tanggal_andalalin").Find(&andalalin)
	resultsPerlalin := ac.DB.Order("tanggal_andalalin").Find(&perlalin)

	if resultsAndalalin.Error != nil && resultsPerlalin != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Tidak ditemukan"})
		return
	} else {
		var respone []models.DaftarAndalalinResponse
		for _, s := range andalalin {
			respone = append(respone, models.DaftarAndalalinResponse{
				IdAndalalin:      s.IdAndalalin,
				Kode:             s.Kode,
				TanggalAndalalin: s.TanggalAndalalin,
				Nama:             s.NamaPemohon,
				Pengembang:       s.NamaPengembang,
				JenisAndalalin:   s.JenisAndalalin,
				StatusAndalalin:  s.StatusAndalalin,
			})
		}
		for _, s := range perlalin {
			respone = append(respone, models.DaftarAndalalinResponse{
				IdAndalalin:      s.IdAndalalin,
				Kode:             s.Kode,
				TanggalAndalalin: s.TanggalAndalalin,
				Nama:             s.NamaPemohon,
				Email:            s.EmailPemohon,
				Petugas:          s.NamaPetugas,
				JenisAndalalin:   s.JenisAndalalin,
				StatusAndalalin:  s.StatusAndalalin,
			})
		}
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "results": len(respone), "data": respone})
	}
}

func (ac *AndalalinController) GetDokumen(ctx *gin.Context) {
	id := ctx.Param("id_andalalin")
	dokumen := ctx.Param("dokumen")

	var andalalin models.Andalalin
	var perlalin models.Perlalin

	ac.DB.First(&andalalin, "id_andalalin = ?", id)
	ac.DB.First(&perlalin, "id_andalalin = ?", id)

	var docs []byte
	var tipe string

	if andalalin.IdAndalalin != uuid.Nil {

		for _, item := range andalalin.Dokumen {
			if item.Dokumen == dokumen {
				docs = item.Berkas
				tipe = item.Tipe
				break
			}
		}

		for _, item := range andalalin.Persyaratan {
			if item.Persyaratan == dokumen {
				docs = item.Berkas
				tipe = item.Tipe
				break
			}
		}
	}

	if perlalin.IdAndalalin != uuid.Nil {
		if dokumen == "Tanda terima pendaftaran" {
			docs = perlalin.TandaTerimaPendaftaran
		}

		if dokumen == "Laporan survei" {
			docs = perlalin.LaporanSurvei
		}

		for _, item := range perlalin.Persyaratan {
			if item.Persyaratan == dokumen {
				docs = item.Berkas
				break
			}
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "tipe": tipe, "data": docs})
}

func (ac *AndalalinController) UpdatePersyaratan(ctx *gin.Context) {
	config, _ := initializers.LoadConfig()

	accessUser := ctx.MustGet("accessUser").(string)

	claim, error := utils.ValidateToken(accessUser, config.AccessTokenPublicKey)
	if error != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": error.Error()})
		return
	}

	credential := claim.Credentials[repository.AndalalinPersyaratanredential]

	if !credential {
		// Return status 403 and permission denied error message.
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": true,
			"msg":   "Permission denied",
		})
		return
	}

	id := ctx.Param("id_andalalin")
	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var andalalin *models.Andalalin
	var perlalin *models.Perlalin

	resultsAndalalin := ac.DB.Find(&andalalin, "id_andalalin = ?", id)
	resultsPerlalin := ac.DB.Find(&perlalin, "id_andalalin = ?", id)

	if resultsAndalalin.Error != nil && resultsPerlalin != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Tidak ditemukan"})
		return
	}

	if andalalin.IdAndalalin != uuid.Nil {
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

				for i := range andalalin.Persyaratan {
					if andalalin.Persyaratan[i].Persyaratan == key {
						andalalin.Persyaratan[i].Berkas = data
						break
					}
				}

			}
		}

		andalalin.StatusAndalalin = "Cek persyaratan"

		ac.DB.Save(&andalalin)
	}

	if perlalin.IdAndalalin != uuid.Nil {
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

				for i := range andalalin.Persyaratan {
					if andalalin.Persyaratan[i].Persyaratan == key {
						andalalin.Persyaratan[i].Berkas = data
						break
					}
				}

			}
		}

		perlalin.PersyaratanTidakSesuai = nil
		perlalin.StatusAndalalin = "Cek persyaratan"

		ac.DB.Save(&perlalin)
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "msg": "persyaratan berhasil diupdate"})
}

func (ac *AndalalinController) UploadDokumen(ctx *gin.Context) {
	id := ctx.Param("id_andalalin")
	dokumen := ctx.Param("dokumen")

	config, _ := initializers.LoadConfig()

	accessUser := ctx.MustGet("accessUser").(string)

	claim, error := utils.ValidateToken(accessUser, config.AccessTokenPublicKey)
	if error != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": error.Error()})
		return
	}

	credential := claim.Credentials[repository.AndalalinDokumenCredential]

	if !credential {
		// Return status 403 and permission denied error message.
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": true,
			"msg":   "Permission denied",
		})
		return
	}

	var andalalin models.Andalalin
	var perlalin models.Perlalin

	resultsAndalalin := ac.DB.First(&andalalin, "id_andalalin = ?", id)
	resultsPerlalin := ac.DB.First(&perlalin, "id_andalalin = ?", id)

	if resultsAndalalin.Error != nil && resultsPerlalin != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Tidak ditemukan"})
		return
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if andalalin.IdAndalalin != uuid.Nil {
		if dokumen == "Checklist administrasi" {
			itemIndex := -1

			for i, item := range andalalin.Dokumen {
				if item.Dokumen == "Checklist administrasi" {
					itemIndex = i
					break
				}
			}

			for _, files := range form.File {
				for _, file := range files {
					// Save the uploaded file with key as prefix
					filed, err := file.Open()

					if err != nil {
						ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
						return
					}
					defer filed.Close()

					data, err := io.ReadAll(filed)
					if err != nil {
						ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
						return
					}
					andalalin.Dokumen[itemIndex].Berkas = data
				}
			}

			andalalin.Dokumen[itemIndex].Role = "User"
			if andalalin.PersyaratanTidakSesuai != nil {
				andalalin.StatusAndalalin = "Persyaratan tidak terpenuhi"
				justString := strings.Join(andalalin.PersyaratanTidakSesuai, "\n")

				data := utils.PersyaratanTidakSesuai{
					Kode:        andalalin.Kode,
					Nama:        andalalin.NamaPemohon,
					Tlp:         andalalin.NomerPemohon,
					Jenis:       andalalin.JenisAndalalin,
					Status:      andalalin.StatusAndalalin,
					Persyaratan: justString,
					Subject:     "Persyaratan tidak terpenuhi",
				}

				utils.SendEmailPersyaratan(andalalin.EmailPemohon, &data)

				var user models.User
				resultUser := ac.DB.First(&user, "id = ?", andalalin.IdUser)
				if resultUser.Error != nil {
					ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "User tidak ditemukan"})
					return
				}

				simpanNotif := models.Notifikasi{
					IdUser: user.ID,
					Title:  "Persyaratan tidak terpenuhi",
					Body:   "Permohonan anda dengan kode " + andalalin.Kode + " terdapat persyaratan yang tidak terpenuhi, silahkan cek email atau permohonan untuk lebih lanjut",
				}

				ac.DB.Create(&simpanNotif)

				if user.PushToken != "" {
					notif := utils.Notification{
						IdUser: user.ID,
						Title:  "Persyaratan tidak terpenuhi",
						Body:   "Permohonan anda dengan kode " + andalalin.Kode + " terdapat persyaratan yang tidak terpenuhi, silahkan cek atau permohonan email untuk lebih lanjut",
						Token:  user.PushToken,
					}

					utils.SendPushNotifications(&notif)
				}
			} else {
				andalalin.StatusAndalalin = "Persyaratan terpenuhi"
			}
		}

		if dokumen == "Surat pernyataan kesanggupan" {
			itemPernyataan := -1

			for i, item := range andalalin.Dokumen {
				if item.Dokumen == "Surat pernyataan kesanggupan (word)" {
					itemPernyataan = i
					break
				}
			}

			for key, files := range form.File {
				for _, file := range files {
					// Save the uploaded file with key as prefix
					filed, err := file.Open()

					if err != nil {
						ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
						return
					}
					defer filed.Close()

					data, err := io.ReadAll(filed)
					if err != nil {
						ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
						return
					}
					if key == "Surat pernyataan kesanggupan (word)" {
						andalalin.Dokumen[itemPernyataan].Berkas = data
					} else {
						andalalin.Dokumen = append(andalalin.Dokumen, models.DokumenPermohonan{Role: "User", Dokumen: "Surat pernyataan kesanggupan (pdf)", Tipe: "Pdf", Berkas: data})
					}

				}
			}

			andalalin.StatusAndalalin = "Menunggu pembayaran"
		}

		if dokumen == "Bukti pembayaran" {
			for key, files := range form.File {
				for _, file := range files {
					// Save the uploaded file with key as prefix
					filed, err := file.Open()

					if err != nil {
						ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
						return
					}
					defer filed.Close()

					data, err := io.ReadAll(filed)
					if err != nil {
						ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
						return
					}
					andalalin.Dokumen = append(andalalin.Dokumen, models.DokumenPermohonan{Role: "User", Dokumen: key, Tipe: "Pdf", Berkas: data})

				}
			}
			andalalin.StatusAndalalin = "Pembuatan surat keputusan"
		}
		ac.DB.Save(&andalalin)
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "msg": "Dokumen berhasil diupload"})
}

func (ac *AndalalinController) CheckAdministrasi(ctx *gin.Context) {
	var payload *models.Administrasi
	id := ctx.Param("id_andalalin")

	config, _ := initializers.LoadConfig()

	currentUser := ctx.MustGet("currentUser").(models.User)

	accessUser := ctx.MustGet("accessUser").(string)

	claim, error := utils.ValidateToken(accessUser, config.AccessTokenPublicKey)
	if error != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": error.Error()})
		return
	}

	credential := claim.Credentials[repository.AndalalinTindakLanjut]

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

	var andalalin models.Andalalin

	result := ac.DB.First(&andalalin, "id_andalalin = ?", id)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Permohonan tidak ditemukan"})
		return
	}

	loc, _ := time.LoadLocation("Asia/Singapore")
	nowTime := time.Now().In(loc)
	tanggal := nowTime.Format("02") + " " + utils.Bulan(nowTime.Month()) + " " + nowTime.Format("2006")

	t, err := template.ParseFiles("templates/checklistAdministrasi.html")
	if err != nil {
		log.Fatal("Error reading the email template:", err)
		return
	}

	var bangkitan string

	switch andalalin.Bangkitan {
	case "Bangkitan rendah":
		bangkitan = "RENDAH"
	case "Bangkitan sedang":
		bangkitan = "SEDANG"
	case "Bangkitan tinggi":
		bangkitan = "TINGGI"
	}

	itemIndex := -1

	for i, item := range andalalin.Dokumen {
		if item.Dokumen == "Checklist administrasi" {
			itemIndex = i
			break
		}
	}

	if itemIndex != -1 {
		administrasi := struct {
			Bangkitan   string
			Objek       string
			Lokasi      string
			Pengembang  string
			Sertifikat  string
			Klasifikasi string
			Nomor       string
			Diterima    string
			Pemeriksaan string
			Status      string
			Data        []models.DataAdministrasi
			Operator    string
			Nip         string
		}{
			Bangkitan:   bangkitan,
			Objek:       andalalin.Jenis,
			Lokasi:      andalalin.NamaJalan + ", " + andalalin.AlamatProyek + ", " + andalalin.KelurahanProyek + ", " + andalalin.KecamatanProyek + ", " + andalalin.KabupatenProyek + ", " + andalalin.ProvinsiProyek + ", " + andalalin.NegaraProyek,
			Pengembang:  andalalin.NamaPengembang,
			Sertifikat:  andalalin.NomerSertifikatPemohon,
			Klasifikasi: andalalin.KlasifikasiPemohon,
			Nomor:       payload.NomorSurat + ", " + payload.TanggalSurat,
			Diterima:    andalalin.TanggalAndalalin,
			Pemeriksaan: tanggal,
			Status:      "Revisi",
			Data:        payload.Data,
			Operator:    currentUser.Name,
			Nip:         *currentUser.NIP,
		}

		buffer := new(bytes.Buffer)
		if err = t.Execute(buffer, administrasi); err != nil {
			log.Fatal("Eror saat membaca template:", err)
			return
		}

		pdfg, err := wkhtmltopdf.NewPDFGenerator()
		if err != nil {
			log.Fatal("Eror generate pdf", err)
			return
		}

		// read the HTML page as a PDF page
		page := wkhtmltopdf.NewPageReader(bytes.NewReader(buffer.Bytes()))

		pdfg.AddPage(page)

		pdfg.Dpi.Set(300)
		pdfg.PageSize.Set(wkhtmltopdf.PageSizeA4)
		pdfg.Orientation.Set(wkhtmltopdf.OrientationPortrait)
		pdfg.MarginBottom.Set(20)
		pdfg.MarginLeft.Set(30)
		pdfg.MarginRight.Set(30)
		pdfg.MarginTop.Set(20)

		err = pdfg.Create()
		if err != nil {
			log.Fatal(err)
		}

		andalalin.Dokumen[itemIndex].Berkas = pdfg.Bytes()
		andalalin.Dokumen[itemIndex].Role = "Dishub"
	} else {
		administrasi := struct {
			Bangkitan   string
			Objek       string
			Lokasi      string
			Pengembang  string
			Sertifikat  string
			Klasifikasi string
			Nomor       string
			Diterima    string
			Pemeriksaan string
			Status      string
			Data        []models.DataAdministrasi
			Operator    string
			Nip         string
		}{
			Bangkitan:   bangkitan,
			Objek:       andalalin.Jenis,
			Lokasi:      andalalin.NamaJalan + ", " + andalalin.AlamatProyek + ", " + andalalin.KelurahanProyek + ", " + andalalin.KecamatanProyek + ", " + andalalin.KabupatenProyek + ", " + andalalin.ProvinsiProyek + ", " + andalalin.NegaraProyek,
			Pengembang:  andalalin.NamaPengembang,
			Sertifikat:  andalalin.NomerSertifikatPemohon,
			Klasifikasi: andalalin.KlasifikasiPemohon,
			Nomor:       payload.NomorSurat + ", " + payload.TanggalSurat,
			Diterima:    andalalin.TanggalAndalalin,
			Pemeriksaan: tanggal,
			Status:      "Baru",
			Data:        payload.Data,
			Operator:    currentUser.Name,
			Nip:         *currentUser.NIP,
		}

		buffer := new(bytes.Buffer)
		if err = t.Execute(buffer, administrasi); err != nil {
			log.Fatal("Eror saat membaca template:", err)
			return
		}

		pdfg, err := wkhtmltopdf.NewPDFGenerator()
		if err != nil {
			log.Fatal("Eror generate pdf", err)
			return
		}

		// read the HTML page as a PDF page
		page := wkhtmltopdf.NewPageReader(bytes.NewReader(buffer.Bytes()))

		pdfg.AddPage(page)

		pdfg.Dpi.Set(300)
		pdfg.PageSize.Set(wkhtmltopdf.PageSizeA4)
		pdfg.Orientation.Set(wkhtmltopdf.OrientationPortrait)
		pdfg.MarginBottom.Set(20)
		pdfg.MarginLeft.Set(30)
		pdfg.MarginRight.Set(30)
		pdfg.MarginTop.Set(20)

		err = pdfg.Create()
		if err != nil {
			log.Fatal(err)
		}

		andalalin.Dokumen = append(andalalin.Dokumen, models.DokumenPermohonan{Role: "Dishub", Dokumen: "Checklist administrasi", Tipe: "Pdf", Berkas: pdfg.Bytes()})
	}

	if andalalin.PersyaratanTidakSesuai != nil {
		andalalin.PersyaratanTidakSesuai = nil
	}

	for _, item := range payload.Data {
		if item.Tidak != "" && item.Kebutuhan == "Wajib" {
			andalalin.PersyaratanTidakSesuai = append(andalalin.PersyaratanTidakSesuai, item.Persyaratan)
		}
	}

	andalalin.Nomor = payload.NomorSurat
	andalalin.Tanggal = payload.TanggalSurat

	ac.DB.Save(&andalalin)

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (ac *AndalalinController) PersyaratanTerpenuhi(ctx *gin.Context) {
	id := ctx.Param("id_andalalin")

	config, _ := initializers.LoadConfig()

	accessUser := ctx.MustGet("accessUser").(string)

	claim, error := utils.ValidateToken(accessUser, config.AccessTokenPublicKey)
	if error != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": error.Error()})
		return
	}

	credential := claim.Credentials[repository.AndalalinTindakLanjut]

	if !credential {
		// Return status 403 and permission denied error message.
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": true,
			"msg":   "Permission denied",
		})
		return
	}

	var perlalin models.Perlalin

	ac.DB.First(&perlalin, "id_andalalin = ?", id)

	if perlalin.IdAndalalin != uuid.Nil {
		perlalin.StatusAndalalin = "Persyaratan terpenuhi"
		ac.DB.Save(&perlalin)
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (ac *AndalalinController) PersyaratanTidakSesuai(ctx *gin.Context) {
	id := ctx.Param("id_andalalin")
	var payload *models.PersayaratanTidakSesuaiInput

	config, _ := initializers.LoadConfig()

	accessUser := ctx.MustGet("accessUser").(string)

	claim, error := utils.ValidateToken(accessUser, config.AccessTokenPublicKey)
	if error != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": error.Error()})
		return
	}

	credential := claim.Credentials[repository.AndalalinTindakLanjut]

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

	var andalalin models.Andalalin
	var perlalin models.Perlalin

	ac.DB.First(&andalalin, "id_andalalin = ?", id)
	ac.DB.First(&perlalin, "id_andalalin = ?", id)

	if andalalin.IdAndalalin != uuid.Nil {
		andalalin.StatusAndalalin = "Persyaratan tidak terpenuhi"
		andalalin.PersyaratanTidakSesuai = payload.Persyaratan

		ac.DB.Save(&andalalin)

		justString := strings.Join(payload.Persyaratan, "\n")

		data := utils.PersyaratanTidakSesuai{
			Kode:        andalalin.Kode,
			Nama:        andalalin.NamaPemohon,
			Tlp:         andalalin.NomerPemohon,
			Jenis:       andalalin.JenisAndalalin,
			Status:      andalalin.StatusAndalalin,
			Persyaratan: justString,
			Subject:     "Persyaratan tidak terpenuhi",
		}

		utils.SendEmailPersyaratan(andalalin.EmailPemohon, &data)

		var user models.User
		resultUser := ac.DB.First(&user, "id = ?", andalalin.IdUser)
		if resultUser.Error != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "User tidak ditemukan"})
			return
		}

		simpanNotif := models.Notifikasi{
			IdUser: user.ID,
			Title:  "Persyaratan tidak terpenuhi",
			Body:   "Permohonan anda dengan kode " + andalalin.Kode + " terdapat persyaratan yang tidak terpenuhi, silahkan cek email atau permohonan untuk lebih lanjut",
		}

		ac.DB.Create(&simpanNotif)

		if user.PushToken != "" {
			notif := utils.Notification{
				IdUser: user.ID,
				Title:  "Persyaratan tidak terpenuhi",
				Body:   "Permohonan anda dengan kode " + andalalin.Kode + " terdapat persyaratan yang tidak terpenuhi, silahkan cek email atau permohonan untuk lebih lanjut",
				Token:  user.PushToken,
			}

			utils.SendPushNotifications(&notif)
		}

	}

	if perlalin.IdAndalalin != uuid.Nil {
		perlalin.StatusAndalalin = "Persyaratan tidak terpenuhi"
		perlalin.PersyaratanTidakSesuai = payload.Persyaratan

		ac.DB.Save(&perlalin)

		justString := strings.Join(payload.Persyaratan, "\n")

		data := utils.PersyaratanTidakSesuai{
			Kode:        perlalin.Kode,
			Nama:        perlalin.NamaPemohon,
			Tlp:         perlalin.NomerPemohon,
			Jenis:       perlalin.JenisAndalalin,
			Status:      perlalin.StatusAndalalin,
			Persyaratan: justString,
			Subject:     "Persyaratan tidak terpenuhi",
		}

		utils.SendEmailPersyaratan(perlalin.EmailPemohon, &data)

		var user models.User
		resultUser := ac.DB.First(&user, "id = ?", perlalin.IdUser)
		if resultUser.Error != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "User tidak ditemukan"})
			return
		}

		simpanNotif := models.Notifikasi{
			IdUser: user.ID,
			Title:  "Persyaratan tidak terpenuhi",
			Body:   "Permohonan anda dengan kode " + perlalin.Kode + " terdapat persyaratan yang tidak terpenuhi, silahkan cek email atau permohonan untuk lebih lanjut",
		}

		ac.DB.Create(&simpanNotif)

		if user.PushToken != "" {
			notif := utils.Notification{
				IdUser: user.ID,
				Title:  "Persyaratan tidak terpenuhi",
				Body:   "Permohonan anda dengan kode " + perlalin.Kode + " terdapat persyaratan yang tidak terpenuhi, silahkan cek email atau permohonan untuk lebih lanjut",
				Token:  user.PushToken,
			}

			utils.SendPushNotifications(&notif)
		}

	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (ac *AndalalinController) UpdateStatusPermohonan(ctx *gin.Context) {
	status := ctx.Param("status")
	id := ctx.Param("id_andalalin")

	config, _ := initializers.LoadConfig()

	accessUser := ctx.MustGet("accessUser").(string)

	claim, error := utils.ValidateToken(accessUser, config.AccessTokenPublicKey)
	if error != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": error.Error()})
		return
	}

	credential := claim.Credentials[repository.AndalalinUpdateCredential]

	if !credential {
		// Return status 403 and permission denied error message.
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": true,
			"msg":   "Permission denied",
		})
		return
	}

	var andalalin models.Andalalin
	var perlalin models.Perlalin

	ac.DB.First(&andalalin, "id_andalalin = ?", id)
	ac.DB.First(&perlalin, "id_andalalin = ?", id)

	if andalalin.IdAndalalin != uuid.Nil {
		andalalin.StatusAndalalin = status
		ac.DB.Save(&andalalin)

	}

	if perlalin.IdAndalalin != uuid.Nil {
		perlalin.StatusAndalalin = status
		ac.DB.Save(&perlalin)
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (ac *AndalalinController) PembuatanSuratPernyataan(ctx *gin.Context) {
	var payload *models.Kewajiban
	id := ctx.Param("id_andalalin")

	config, _ := initializers.LoadConfig()

	accessUser := ctx.MustGet("accessUser").(string)

	claim, error := utils.ValidateToken(accessUser, config.AccessTokenPublicKey)
	if error != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": error.Error()})
		return
	}

	credential := claim.Credentials[repository.AndalalinTindakLanjut]

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

	var andalalin models.Andalalin

	result := ac.DB.First(&andalalin, "id_andalalin = ?", id)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Permohonan tidak ditemukan"})
		return
	}

	var bangkitan string

	switch andalalin.Bangkitan {
	case "Bangkitan rendah":
		bangkitan = "Rendah"
	case "Bangkitan sedang":
		bangkitan = "Sedang"
	case "Bangkitan tinggi":
		bangkitan = "Tinggi"
	}

	listContent := ""
	for i, item := range payload.Kewajiban {
		if i == len(payload.Kewajiban)-1 {
			listContent += fmt.Sprint(i+1, ". ", item)
		} else {
			listContent += fmt.Sprint(i+1, ". ", item, "\n")
		}
	}

	replaceMap := docx.PlaceholderMap{
		"_nama_":       andalalin.NamaPimpinanPengembang,
		"_jabatan_":    andalalin.JabatanPimpinanPengembang,
		"_alamat_":     andalalin.AlamatPimpinanPengembang + ", " + andalalin.KelurahanPimpinanPengembang + ", " + andalalin.KecamatanPimpinanPengembang + ", " + andalalin.KabupatenPimpinanPengembang + ", " + andalalin.ProvinsiPimpinanPengembang + ", " + andalalin.NegaraPimpinanPengembang,
		"_pengembang_": andalalin.NamaPengembang,
		"_bangkitan_":  bangkitan,
		"_nomor_":      andalalin.Nomor,
		"_tanggal_":    andalalin.Tanggal[0:2],
		"_bulan_":      utils.Month(andalalin.Tanggal[3:5]),
		"_tahun_":      andalalin.Tanggal[6:10],
		"_kegiatan_":   andalalin.JenisProyek + " " + andalalin.Jenis,
		"_kewajiban_":  listContent,
	}

	doc, err := docx.Open("templates/suratPernyataanKesanggupan.docx")
	if err != nil {
		panic(err)
	}

	// replace the keys with values from replaceMap
	err = doc.ReplaceAll(replaceMap)
	if err != nil {
		panic(err)
	}

	tempFilePath := "temp.docx"
	err = doc.WriteToFile(tempFilePath)
	if err != nil {
		log.Fatal(err)
	}

	docBytes, err := os.ReadFile(tempFilePath)
	if err != nil {
		log.Fatal(err)
	}

	_ = os.Remove(tempFilePath)

	andalalin.StatusAndalalin = "Memberikan pernyataan"

	itemIndex := -1

	for i, item := range andalalin.Dokumen {
		if item.Dokumen == "Surat pernyataan kesanggupan (word)" {
			itemIndex = i
			break
		}
	}

	if itemIndex != -1 {
		andalalin.Dokumen[itemIndex].Berkas = docBytes
	} else {
		andalalin.Dokumen = append(andalalin.Dokumen, models.DokumenPermohonan{Role: "User", Dokumen: "Surat pernyataan kesanggupan (word)", Tipe: "Word", Berkas: docBytes})
	}

	ac.DB.Save(&andalalin)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Surat berhasil dibuat"})
}

func (ac *AndalalinController) PembuatanSuratKeputusan(ctx *gin.Context) {}

func (ac *AndalalinController) TambahPetugas(ctx *gin.Context) {
	var payload *models.TambahPetugas
	id := ctx.Param("id_andalalin")

	config, _ := initializers.LoadConfig()

	accessUser := ctx.MustGet("accessUser").(string)

	claim, error := utils.ValidateToken(accessUser, config.AccessTokenPublicKey)
	if error != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": error.Error()})
		return
	}

	credential := claim.Credentials[repository.AndalalinAddOfficerCredential]

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

	var perlalin models.Perlalin

	resultsPerlalin := ac.DB.First(&perlalin, "id_andalalin = ?", id)

	if resultsPerlalin.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Tidak ditemukan"})
		return
	}

	if perlalin.IdAndalalin != uuid.Nil {
		perlalin.IdPetugas = payload.IdPetugas
		perlalin.NamaPetugas = payload.NamaPetugas
		perlalin.EmailPetugas = payload.EmailPetugas
		perlalin.StatusAndalalin = "Survei lapangan"

		ac.DB.Save(&perlalin)

		ac.ReleaseTicketLevel2(ctx, perlalin.IdAndalalin, payload.IdPetugas)
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Tambah petugas berhasil"})
}

func (ac *AndalalinController) GantiPetugas(ctx *gin.Context) {
	var payload *models.TambahPetugas
	id := ctx.Param("id_andalalin")

	config, _ := initializers.LoadConfig()

	accessUser := ctx.MustGet("accessUser").(string)

	claim, error := utils.ValidateToken(accessUser, config.AccessTokenPublicKey)
	if error != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": error.Error()})
		return
	}

	credential := claim.Credentials[repository.AndalalinAddOfficerCredential]

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

	var perlalin models.Perlalin

	resultsPerlalin := ac.DB.First(&perlalin, "id_andalalin = ?", id)

	if resultsPerlalin.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Tidak ditemukan"})
		return
	}

	if perlalin.IdAndalalin != uuid.Nil {
		perlalin.IdPetugas = payload.IdPetugas
		perlalin.NamaPetugas = payload.NamaPetugas
		perlalin.EmailPetugas = payload.EmailPetugas
		if perlalin.StatusAndalalin == "Survei lapangan" {
			ac.CloseTiketLevel2(ctx, perlalin.IdAndalalin)

			ac.ReleaseTicketLevel2(ctx, perlalin.IdAndalalin, payload.IdPetugas)
		}

		ac.DB.Save(&perlalin)

	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Ubah petugas berhasil"})
}

func (ac *AndalalinController) GetAndalalinTicketLevel2(ctx *gin.Context) {
	status := ctx.Param("status")
	currentUser := ctx.MustGet("currentUser").(models.User)

	config, _ := initializers.LoadConfig()

	accessUser := ctx.MustGet("accessUser").(string)

	claim, error := utils.ValidateToken(accessUser, config.AccessTokenPublicKey)
	if error != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": error.Error()})
		return
	}

	credential := claim.Credentials[repository.AndalalinTicket2Credential]

	if !credential {
		// Return status 403 and permission denied error message.
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": true,
			"msg":   "Permission denied",
		})
		return
	}

	var ticket []models.TiketLevel2

	results := ac.DB.Find(&ticket, "status = ? AND id_petugas = ?", status, currentUser.ID)

	if results.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": results.Error})
		return
	} else {
		var respone []models.DaftarAndalalinResponse
		for _, s := range ticket {
			var perlalin models.Perlalin
			var usulan models.UsulanPengelolaan

			ac.DB.First(&perlalin, "id_andalalin = ? AND id_petugas = ?", s.IdAndalalin, currentUser.ID)

			ac.DB.First(&usulan, "id_andalalin = ?", s.IdAndalalin)

			if perlalin.IdAndalalin != uuid.Nil && usulan.IdUsulan == uuid.Nil {
				respone = append(respone, models.DaftarAndalalinResponse{
					IdAndalalin:      perlalin.IdAndalalin,
					Kode:             perlalin.Kode,
					TanggalAndalalin: perlalin.TanggalAndalalin,
					Nama:             perlalin.NamaPemohon,
					Email:            perlalin.EmailPemohon,
					Petugas:          perlalin.NamaPetugas,
					JenisAndalalin:   perlalin.JenisAndalalin,
					StatusAndalalin:  perlalin.StatusAndalalin,
				})
			}

		}
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "results": len(respone), "data": respone})
	}
}

func (ac *AndalalinController) IsiSurvey(ctx *gin.Context) {
	var payload *models.DataSurvey
	currentUser := ctx.MustGet("currentUser").(models.User)
	id := ctx.Param("id_andalalin")

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

	var ticket1 models.TiketLevel1
	var ticket2 models.TiketLevel2

	resultTiket1 := ac.DB.Find(&ticket1, "id_andalalin = ?", id)
	resultTiket2 := ac.DB.Find(&ticket2, "id_andalalin = ?", id)
	if resultTiket1.Error != nil && resultTiket2.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Tiket tidak ditemukan"})
		return
	}

	var perlalin models.Perlalin
	resultsPerlalin := ac.DB.First(&perlalin, "id_andalalin = ?", id)

	if resultsPerlalin.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Tidak ditemukan"})
		return
	}

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

	if perlalin.IdAndalalin != uuid.Nil {
		survey := models.Survei{
			IdAndalalin:   perlalin.IdAndalalin,
			IdTiketLevel1: ticket1.IdTiketLevel1,
			IdTiketLevel2: ticket2.IdTiketLevel2,
			IdPetugas:     currentUser.ID,
			Petugas:       currentUser.Name,
			EmailPetugas:  currentUser.Email,
			Lokasi:        payload.Data.Lokasi,
			Keterangan:    payload.Data.Keterangan,
			Foto1:         blobs["foto1"],
			Foto2:         blobs["foto2"],
			Foto3:         blobs["foto3"],
			Latitude:      payload.Data.Latitude,
			Longitude:     payload.Data.Longitude,
			TanggalSurvei: tanggal,
			WaktuSurvei:   nowTime.Format("15:04:05"),
		}

		result := ac.DB.Create(&survey)

		if result.Error != nil && strings.Contains(result.Error.Error(), "duplicate key value violates unique") {
			ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "Data survey sudah tersedia"})
			return
		} else if result.Error != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Telah terjadi sesuatu"})
			return
		}

		perlalin.StatusAndalalin = "Laporan survei"

		ac.DB.Save(&perlalin)

		ac.CloseTiketLevel2(ctx, perlalin.IdAndalalin)
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success"})
}

func (ac *AndalalinController) GetAllSurvey(ctx *gin.Context) {
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

	var survey []models.Survei

	results := ac.DB.Find(&survey, "id_petugas = ?", currentUser.ID)

	if results.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": results.Error})
		return
	} else {
		var respone []models.DaftarAndalalinResponse
		for _, s := range survey {
			var perlalin models.Perlalin

			ac.DB.First(&perlalin, "id_andalalin = ?", s.IdAndalalin)

			if perlalin.IdAndalalin != uuid.Nil {
				respone = append(respone, models.DaftarAndalalinResponse{
					IdAndalalin:      perlalin.IdAndalalin,
					Kode:             perlalin.Kode,
					TanggalAndalalin: perlalin.TanggalAndalalin,
					Nama:             perlalin.NamaPemohon,
					Email:            perlalin.EmailPemohon,
					Petugas:          perlalin.NamaPetugas,
					JenisAndalalin:   perlalin.JenisAndalalin,
					StatusAndalalin:  perlalin.StatusAndalalin,
				})
			}

		}
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "results": len(respone), "data": respone})
	}
}

func (ac *AndalalinController) GetSurvey(ctx *gin.Context) {
	id := ctx.Param("id_andalalin")

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

	var survey *models.Survei

	result := ac.DB.First(&survey, "id_andalalin = ?", id)
	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": survey})
}

func (ac *AndalalinController) PersetujuanDokumen(ctx *gin.Context) {
	var payload *models.Persetujuan
	id := ctx.Param("id_andalalin")

	config, _ := initializers.LoadConfig()

	accessUser := ctx.MustGet("accessUser").(string)

	claim, error := utils.ValidateToken(accessUser, config.AccessTokenPublicKey)
	if error != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": error.Error()})
		return
	}

	credential := claim.Credentials[repository.AndalalinPersetujuanCredential]

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

	var andalalin models.Andalalin

	result := ac.DB.First(&andalalin, "id_andalalin = ?", id)
	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error})
		return
	}

	andalalin.PersetujuanDokumen = payload.Persetujuan
	andalalin.KeteranganPersetujuanDokumen = payload.Keterangan
	if payload.Persetujuan == "Dokumen disetujui" {
		andalalin.StatusAndalalin = "Pembuatan surat keputusan"
		ac.CloseTiketLevel2(ctx, andalalin.IdAndalalin)
	} else {
		andalalin.StatusAndalalin = "Berita acara pemeriksaan"
	}

	ac.DB.Save(&andalalin)

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (ac *AndalalinController) PermohonanSelesai(ctx *gin.Context, id uuid.UUID) {
	var andalalin models.Andalalin

	result := ac.DB.First(&andalalin, "id_andalalin = ?", id)
	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error})
		return
	}

	andalalin.StatusAndalalin = "Permohonan selesai"

	ac.DB.Save(&andalalin)

	data := utils.PermohonanSelesai{
		Kode:    andalalin.Kode,
		Nama:    andalalin.NamaPemohon,
		Tlp:     andalalin.NomerPemohon,
		Jenis:   andalalin.JenisAndalalin,
		Status:  andalalin.StatusAndalalin,
		Subject: "Permohonan telah selesai",
	}

	utils.SendEmailPermohonanSelesai(andalalin.EmailPemohon, &data)

	var user models.User
	resultUser := ac.DB.First(&user, "id = ?", andalalin.IdUser)
	if resultUser.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "User tidak ditemukan"})
		return
	}

	simpanNotif := models.Notifikasi{
		IdUser: user.ID,
		Title:  "Permohonan selesai",
		Body:   "Permohonan anda dengan kode " + andalalin.Kode + " telah selesai, silahkan cek permohonan pada aplikasi untuk lebih jelas",
	}

	ac.DB.Create(&simpanNotif)

	if user.PushToken != "" {
		notif := utils.Notification{
			IdUser: user.ID,
			Title:  "Permohonan selesai",
			Body:   "Permohonan anda dengan kode " + andalalin.Kode + " telah selesai, silahkan cek permohonan pada aplikasi untuk lebih jelas",
			Token:  user.PushToken,
		}

		utils.SendPushNotifications(&notif)
	}
}

func (ac *AndalalinController) UsulanTindakanPengelolaan(ctx *gin.Context) {
	var payload *models.InputUsulanPengelolaan
	id := ctx.Param("id_andalalin")
	currentUser := ctx.MustGet("currentUser").(models.User)

	config, _ := initializers.LoadConfig()

	accessUser := ctx.MustGet("accessUser").(string)

	claim, error := utils.ValidateToken(accessUser, config.AccessTokenPublicKey)
	if error != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": error.Error()})
		return
	}

	credential := claim.Credentials[repository.AndalalinKelolaTiket]

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

	var ticket1 models.TiketLevel1
	var ticket2 models.TiketLevel2

	resultTiket1 := ac.DB.Find(&ticket1, "id_andalalin = ?", id)
	resultTiket2 := ac.DB.Find(&ticket2, "id_andalalin = ?", id)
	if resultTiket1.Error != nil && resultTiket2.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Tiket tidak ditemukan"})
		return
	}

	var perlalin models.Perlalin

	resultsPerlalin := ac.DB.First(&perlalin, "id_andalalin = ?", id)

	if resultsPerlalin.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Tidak ditemukan"})
		return
	}

	if perlalin.IdAndalalin != uuid.Nil {
		usul := models.UsulanPengelolaan{
			IdAndalalin:                perlalin.IdAndalalin,
			IdTiketLevel1:              ticket1.IdTiketLevel1,
			IdTiketLevel2:              ticket2.IdTiketLevel2,
			IdPengusulTindakan:         currentUser.ID,
			NamaPengusulTindakan:       currentUser.Name,
			PertimbanganUsulanTindakan: payload.PertimbanganUsulanTindakan,
			KeteranganUsulanTindakan:   payload.KeteranganUsulanTindakan,
		}

		result := ac.DB.Create(&usul)

		if result.Error != nil && strings.Contains(result.Error.Error(), "duplicate key value violates unique") {
			ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "Usulan sudah ada"})
			return
		} else if result.Error != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Telah terjadi sesuatu"})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (ac *AndalalinController) GetUsulan(ctx *gin.Context) {
	config, _ := initializers.LoadConfig()

	accessUser := ctx.MustGet("accessUser").(string)

	claim, error := utils.ValidateToken(accessUser, config.AccessTokenPublicKey)
	if error != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": error.Error()})
		return
	}

	credential := claim.Credentials[repository.AndalalinKelolaTiket]

	if !credential {
		// Return status 403 and permission denied error message.
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": true,
			"msg":   "Permission denied",
		})
		return
	}

	var usulan []models.UsulanPengelolaan

	results := ac.DB.Find(&usulan)

	if results.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": results.Error})
		return
	} else {
		var respone []models.DaftarAndalalinResponse
		for _, s := range usulan {
			var ticket2 models.TiketLevel2
			resultTiket2 := ac.DB.Not("status = ?", "Tunda").Where("id_andalalin = ? AND status = ?", s.IdAndalalin, "Buka").First(&ticket2)
			if resultTiket2.Error == nil {
				var perlalin models.Perlalin

				resultsPerlalin := ac.DB.First(&perlalin, "id_andalalin = ?", ticket2.IdAndalalin)

				if resultsPerlalin.Error != nil {
					ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Tidak ditemukan"})
					return
				}

				if perlalin.IdAndalalin != uuid.Nil {
					respone = append(respone, models.DaftarAndalalinResponse{
						IdAndalalin:      perlalin.IdAndalalin,
						Kode:             perlalin.Kode,
						TanggalAndalalin: perlalin.TanggalAndalalin,
						Nama:             perlalin.NamaPemohon,
						Email:            perlalin.EmailPemohon,
						Petugas:          perlalin.NamaPetugas,
						JenisAndalalin:   perlalin.JenisAndalalin,
						StatusAndalalin:  perlalin.StatusAndalalin,
					})
				}

			}
		}
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "results": len(respone), "data": respone})
	}
}

func (ac *AndalalinController) GetDetailUsulan(ctx *gin.Context) {
	id := ctx.Param("id_andalalin")

	config, _ := initializers.LoadConfig()

	accessUser := ctx.MustGet("accessUser").(string)

	claim, error := utils.ValidateToken(accessUser, config.AccessTokenPublicKey)
	if error != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": error.Error()})
		return
	}

	credential := claim.Credentials[repository.AndalalinKelolaTiket]

	if !credential {
		// Return status 403 and permission denied error message.
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": true,
			"msg":   "Permission denied",
		})
		return
	}

	var usulan models.UsulanPengelolaan

	resultUsulan := ac.DB.First(&usulan, "id_andalalin = ?", id)
	if resultUsulan.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Telah terjadi sesuatu"})
		return
	}

	data := struct {
		NamaPengusulTindakan       string  `json:"nama,omitempty"`
		PertimbanganUsulanTindakan string  `json:"pertimbangan,omitempty"`
		KeteranganUsulanTindakan   *string `json:"keterangan,omitempty"`
	}{
		NamaPengusulTindakan:       usulan.NamaPengusulTindakan,
		PertimbanganUsulanTindakan: usulan.PertimbanganUsulanTindakan,
		KeteranganUsulanTindakan:   usulan.KeteranganUsulanTindakan,
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": data})
}

func (ac *AndalalinController) TindakanPengelolaan(ctx *gin.Context) {
	id := ctx.Param("id_andalalin")
	jenis := ctx.Param("jenis_pelaksanaan")

	config, _ := initializers.LoadConfig()

	accessUser := ctx.MustGet("accessUser").(string)

	claim, error := utils.ValidateToken(accessUser, config.AccessTokenPublicKey)
	if error != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": error.Error()})
		return
	}

	credential := claim.Credentials[repository.AndalalinKelolaTiket]

	if !credential {
		// Return status 403 and permission denied error message.
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": true,
			"msg":   "Permission denied",
		})
		return
	}

	var tiket models.TiketLevel2

	result := ac.DB.Model(&tiket).Where("id_andalalin = ? AND status = ?", id, "Buka").Or("id_andalalin = ? AND status = ?", id, "Tunda").Update("status", jenis)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Telah terjadi sesuatu"})
		return
	}

	var usulan models.UsulanPengelolaan

	ac.DB.First(&usulan, "id_andalalin = ?", id)
	var perlalin models.Perlalin

	resultsPerlalin := ac.DB.First(&perlalin, "id_andalalin = ?", id)

	if resultsPerlalin.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Tidak ditemukan"})
		return
	}

	var userPengusul models.User
	resulPengusul := ac.DB.First(&userPengusul, "id = ?", usulan.IdPengusulTindakan)
	if resulPengusul.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "User tidak ditemukan"})
		return
	}

	if perlalin.IdAndalalin != uuid.Nil {
		var userPetugas models.User
		resultPetugas := ac.DB.First(&userPetugas, "id = ?", perlalin.IdPetugas)
		if resultPetugas.Error != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "User tidak ditemukan"})
			return
		}

		switch jenis {
		case "Tunda":
			simpanNotifPengusul := models.Notifikasi{
				IdUser: userPengusul.ID,
				Title:  "Pelaksanaan survei ditunda",
				Body:   "Usulan tindakan anda pada permohonan dengan kode " + perlalin.Kode + " telah diputuskan bahwa pelaksanaan survei ditunda",
			}

			ac.DB.Create(&simpanNotifPengusul)

			simpanNotifPetugas := models.Notifikasi{
				IdUser: userPetugas.ID,
				Title:  "Pelaksanaan survei ditunda",
				Body:   "Pelakasnaan survei pada permohonan dengan kode " + perlalin.Kode + " dibatalkan",
			}

			ac.DB.Create(&simpanNotifPetugas)

			if userPengusul.PushToken != "" {
				notifPengusul := utils.Notification{
					IdUser: userPengusul.ID,
					Title:  "Pelaksanaan survei ditunda",
					Body:   "Usulan tindakan anda pada permohonan dengan kode " + perlalin.Kode + " telah diputuskan bahwa pelaksanaan survei ditunda",
					Token:  userPengusul.PushToken,
				}

				utils.SendPushNotifications(&notifPengusul)
			}

			if userPetugas.PushToken != "" {
				notifPetugas := utils.Notification{
					IdUser: userPetugas.ID,
					Title:  "Pelaksanaan survei ditunda",
					Body:   "Pelakasnaan survei pada permohonan dengan kode " + perlalin.Kode + " ditunda",
					Token:  userPetugas.PushToken,
				}

				utils.SendPushNotifications(&notifPetugas)
			}
		case "Batal":
			simpanNotifPengusul := models.Notifikasi{
				IdUser: userPengusul.ID,
				Title:  "Pelaksanaan survei dibatalkan",
				Body:   "Usulan tindakan anda pada permohonan dengan kode " + perlalin.Kode + " telah diputuskan bahwa pelaksanaan survei dibatalkan",
			}

			ac.DB.Create(&simpanNotifPengusul)

			simpanNotifPetugas := models.Notifikasi{
				IdUser: userPetugas.ID,
				Title:  "Pelaksanaan survei dibatalkan",
				Body:   "Pelakasnaan survei pada permohonan dengan kode " + perlalin.Kode + " dibatalkan",
			}

			ac.DB.Create(&simpanNotifPetugas)

			if userPengusul.PushToken != "" {
				notifPengusul := utils.Notification{
					IdUser: userPengusul.ID,
					Title:  "Pelaksanaan survei dibatalkan",
					Body:   "Usulan tindakan anda pada permohonan dengan kode " + perlalin.Kode + " telah diputuskan bahwa pelaksanaan survei dibatalkan",
					Token:  userPengusul.PushToken,
				}

				utils.SendPushNotifications(&notifPengusul)
			}

			if userPetugas.PushToken != "" {
				notifPetugas := utils.Notification{
					IdUser: userPetugas.ID,
					Title:  "Pelaksanaan survei dibatalkan",
					Body:   "Pelakasnaan survei pada permohonan dengan kode " + perlalin.Kode + " dibatalkan",
					Token:  userPetugas.PushToken,
				}

				utils.SendPushNotifications(&notifPetugas)
			}
		case "Buka":
			simpanNotifPetugas := models.Notifikasi{
				IdUser: userPetugas.ID,
				Title:  "Pelaksanaan survei dilanjutkan",
				Body:   "Pelaksanaan survei pada permohonan dengan kode " + perlalin.Kode + " telah dilanjutkan kembali",
			}

			ac.DB.Create(&simpanNotifPetugas)

			if userPetugas.PushToken != "" {
				notifPetugas := utils.Notification{
					IdUser: userPetugas.ID,
					Title:  "Pelaksanaan survei dilanjutkan",
					Body:   "Pelaksanaan survei pada permohonan dengan kode " + perlalin.Kode + " telah dilanjutkan kembali",
					Token:  userPetugas.PushToken,
				}

				utils.SendPushNotifications(&notifPetugas)
			}
		}
	}

	if jenis == "Batal" || jenis == "Buka" {
		ac.DB.Delete(&models.UsulanPengelolaan{}, "id_andalalin = ?", id)
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success"})
}

func (ac *AndalalinController) HapusUsulan(ctx *gin.Context) {
	id := ctx.Param("id_andalalin")

	config, _ := initializers.LoadConfig()

	accessUser := ctx.MustGet("accessUser").(string)

	claim, error := utils.ValidateToken(accessUser, config.AccessTokenPublicKey)
	if error != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": error.Error()})
		return
	}

	credential := claim.Credentials[repository.AndalalinKelolaTiket]

	if !credential {
		// Return status 403 and permission denied error message.
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": true,
			"msg":   "Permission denied",
		})
		return
	}

	var usulan models.UsulanPengelolaan

	resultUsulan := ac.DB.First(&usulan, "id_andalalin = ?", id)
	if resultUsulan.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "User tidak ditemukan"})
		return
	}

	var perlalin models.Perlalin

	resultsPerlalin := ac.DB.First(&perlalin, "id_andalalin = ?", id)

	if resultsPerlalin.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Tidak ditemukan"})
		return
	}

	var userPengusul models.User
	resulPengusul := ac.DB.First(&userPengusul, "id = ?", usulan.IdPengusulTindakan)
	if resulPengusul.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "User tidak ditemukan"})
		return
	}

	if perlalin.IdAndalalin != uuid.Nil {
		simpanNotifPengusul := models.Notifikasi{
			IdUser: userPengusul.ID,
			Title:  "Usulan tindakan dihapus",
			Body:   "Usulan tindakan anda pada permohonan dengan kode " + perlalin.Kode + " telah dihapus",
		}

		ac.DB.Create(&simpanNotifPengusul)

		if userPengusul.PushToken != "" {
			notifPengusul := utils.Notification{
				IdUser: userPengusul.ID,
				Title:  "Usulan tindakan dihapus",
				Body:   "Usulan tindakan anda pada permohonan dengan kode " + perlalin.Kode + " telah dihapus",
				Token:  userPengusul.PushToken,
			}

			utils.SendPushNotifications(&notifPengusul)
		}
	}

	results := ac.DB.Delete(&models.UsulanPengelolaan{}, "id_andalalin = ?", id)

	if results.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": results.Error})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success"})
}

func (ac *AndalalinController) GetAllAndalalinByTiketLevel2(ctx *gin.Context) {
	status := ctx.Param("status")

	config, _ := initializers.LoadConfig()

	accessUser := ctx.MustGet("accessUser").(string)

	claim, error := utils.ValidateToken(accessUser, config.AccessTokenPublicKey)
	if error != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": error.Error()})
		return
	}

	credential := claim.Credentials[repository.AndalalinTicket2Credential]

	if !credential {
		// Return status 403 and permission denied error message.
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": true,
			"msg":   "Permission denied",
		})
		return
	}

	var ticket []models.TiketLevel2

	results := ac.DB.Find(&ticket, "status = ?", status)

	if results.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": results.Error})
		return
	} else {
		var respone []models.DaftarAndalalinResponse
		for _, s := range ticket {
			var perlalin models.Perlalin

			resultsPerlalin := ac.DB.First(&perlalin, "id_andalalin = ?", s.IdAndalalin)

			if resultsPerlalin.Error != nil {
				ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Tidak ditemukan"})
				return
			}

			if perlalin.IdAndalalin != uuid.Nil {
				respone = append(respone, models.DaftarAndalalinResponse{
					IdAndalalin:      perlalin.IdAndalalin,
					Kode:             perlalin.Kode,
					TanggalAndalalin: perlalin.TanggalAndalalin,
					Nama:             perlalin.NamaPemohon,
					Email:            perlalin.EmailPemohon,
					Petugas:          perlalin.NamaPetugas,
					JenisAndalalin:   perlalin.JenisAndalalin,
					StatusAndalalin:  perlalin.StatusAndalalin,
				})
			}

		}
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "results": len(respone), "data": respone})
	}
}

func (ac *AndalalinController) LaporanSurvei(ctx *gin.Context) {
	id := ctx.Param("id_andalalin")

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

	var perlalin models.Perlalin

	result := ac.DB.First(&perlalin, "id_andalalin = ?", id)
	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error})
		return
	}

	file, err := ctx.FormFile("ls")
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

	perlalin.LaporanSurvei = data
	perlalin.StatusAndalalin = "Menunggu hasil keputusan"

	resultLaporan := ac.DB.Save(&perlalin)

	if resultLaporan.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Telah terjadi sesuatu"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success"})
}

func (ac *AndalalinController) IsiSurveyMandiri(ctx *gin.Context) {
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
		Keterangan:    payload.Data.Keterangan,
		Foto1:         blobs["foto1"],
		Foto2:         blobs["foto2"],
		Foto3:         blobs["foto3"],
		Latitude:      payload.Data.Latitude,
		Longitude:     payload.Data.Longitude,
		StatusSurvei:  "Perlu tindakan",
		TanggalSurvei: tanggal,
		WaktuSurvei:   nowTime.Format("15:04:05"),
	}

	result := ac.DB.Create(&survey)

	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Telah terjadi sesuatu"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": survey})
}

func (ac *AndalalinController) GetAllSurveiMandiri(ctx *gin.Context) {
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

	results := ac.DB.Find(&survey)

	if results.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": results.Error})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "results": len(survey), "data": survey})
}

func (ac *AndalalinController) GetAllSurveiMandiriByPetugas(ctx *gin.Context) {
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

	results := ac.DB.Find(&survey, "id_petugas = ?", currentUser.ID)

	if results.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": results.Error})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "results": len(survey), "data": survey})
}

func (ac *AndalalinController) GetSurveiMandiri(ctx *gin.Context) {
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

	result := ac.DB.First(&survey, "id_survey = ?", id)
	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": survey})
}

func (ac *AndalalinController) TerimaSurvei(ctx *gin.Context) {
	id := ctx.Param("id_survei")
	keterangan := ctx.Param("keterangan")

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

	result := ac.DB.First(&survey, "id_survey = ?", id)
	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error})
		return
	}
	survey.StatusSurvei = "Survei diterima"
	survey.KeteranganTindakan = keterangan

	ac.DB.Save(&survey)

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": survey})
}

func (ac *AndalalinController) KeputusanHasil(ctx *gin.Context) {
	id := ctx.Param("id_andalalin")
	var payload *models.KeputusanHasil

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

	credential := claim.Credentials[repository.AndalalinKeputusanHasil]

	if !credential {
		// Return status 403 and permission denied error message.
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": true,
			"msg":   "Permission denied",
		})
		return
	}

	var perlalin models.Perlalin

	result := ac.DB.First(&perlalin, "id_andalalin = ?", id)
	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error})
		return
	}

	if payload.Keputusan == "Pemasangan ditunda" {
		perlalin.Tindakan = payload.Keputusan
		perlalin.PertimbanganTindakan = payload.Pertimbangan
		perlalin.StatusAndalalin = "Tunda pemasangan"
	} else if payload.Keputusan == "Pemasangan disegerakan" {
		perlalin.Tindakan = payload.Keputusan
		perlalin.StatusAndalalin = "Pemasangan sedang dilakukan"
	}

	resultKeputusan := ac.DB.Save(&perlalin)

	if resultKeputusan.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Telah terjadi sesuatu"})
		return
	}

	var mutex sync.Mutex

	updateChannelTunda := make(chan struct{})
	updateChannelDisegerakan := make(chan struct{})

	if payload.Keputusan == "Pemasangan ditunda" {
		go func() {
			duration := 3 * 24 * time.Hour
			timer := time.NewTimer(duration)

			select {
			case <-timer.C:
				mutex.Lock()
				defer mutex.Unlock()

				var data models.Perlalin

				result := ac.DB.First(&data, "id_andalalin = ?", id)
				if result.Error != nil {
					ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error})
					return
				}

				if data.StatusAndalalin == "Tunda pemasangan" {
					ac.CloseTiketLevel1(ctx, data.IdAndalalin)
					ac.BatalkanPermohonan(ctx, id)
					data.Tindakan = "Permohonan dibatalkan"
					data.PertimbanganTindakan = "Permohonan dibatalkan"
					data.StatusAndalalin = "Permohonan dibatalkan"
					ac.DB.Save(&data)
					updateChannelTunda <- struct{}{}
				}
			case <-updateChannelTunda:
				// The update was canceled, do nothing
			}
		}()
	} else if payload.Keputusan == "Pemasangan disegerakan" {
		if perlalin.StatusAndalalin == "Tunda pemasangan" {
			close(updateChannelTunda)
		}

		go func() {
			duration := 3 * 24 * time.Hour
			timer := time.NewTimer(duration)

			select {
			case <-timer.C:
				mutex.Lock()
				defer mutex.Unlock()

				var data models.Perlalin

				result := ac.DB.First(&data, "id_andalalin = ?", id)
				if result.Error != nil {
					ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error})
					return
				}

				if data.StatusAndalalin == "Pemasangan sedang dilakukan" {
					data.Tindakan = "Pemasangan ditunda"
					data.PertimbanganTindakan = "Pemasangan ditunda"
					data.StatusAndalalin = "Tunda pemasangan"
					ac.DB.Save(&data)

					updateChannelTunda = make(chan struct{})

					go func() {
						duration := 3 * 24 * time.Hour
						timer := time.NewTimer(duration)

						select {
						case <-timer.C:
							mutex.Lock()
							defer mutex.Unlock()

							var data models.Perlalin

							result := ac.DB.First(&data, "id_andalalin = ?", id)
							if result.Error != nil {
								ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error})
								return
							}

							if data.StatusAndalalin == "Tunda pemasangan" {
								ac.CloseTiketLevel1(ctx, data.IdAndalalin)
								ac.BatalkanPermohonan(ctx, id)
								updateChannelTunda <- struct{}{}
								updateChannelDisegerakan <- struct{}{}
							}

						case <-updateChannelTunda:
							// The update was canceled, do nothing
						}
					}()
				}
			case <-updateChannelDisegerakan:
				// The update was canceled, do nothing
			}
		}()
	} else if payload.Keputusan == "Batalkan permohonan" {
		if perlalin.StatusAndalalin == "Tunda pemasangan" {
			close(updateChannelTunda)
		}

		ac.CloseTiketLevel1(ctx, perlalin.IdAndalalin)
		ac.BatalkanPermohonan(ctx, id)
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (ac *AndalalinController) BatalkanPermohonan(ctx *gin.Context, id string) {
	var permohonan models.Perlalin

	result := ac.DB.First(&permohonan, "id_andalalin = ?", id)
	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error})
		return
	}

	permohonan.Tindakan = "Permohonan dibatalkan"
	permohonan.PertimbanganTindakan = "Permohonan dibatalkan"
	permohonan.StatusAndalalin = "Permohonan dibatalkan"
	ac.DB.Save(&permohonan)

	data := utils.PermohonanDibatalkan{
		Kode:    permohonan.Kode,
		Nama:    permohonan.NamaPemohon,
		Tlp:     permohonan.NomerPemohon,
		Jenis:   permohonan.JenisAndalalin,
		Status:  permohonan.StatusAndalalin,
		Subject: "Permohonan dibatalkan",
	}

	utils.SendEmailPermohonanDibatalkan(permohonan.EmailPemohon, &data)

	var user models.User
	resultUser := ac.DB.First(&user, "id = ?", permohonan.IdUser)
	if resultUser.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "User tidak ditemukan"})
		return
	}

	simpanNotif := models.Notifikasi{
		IdUser: user.ID,
		Title:  "Permohonan dibatalkan",
		Body:   "Permohonan anda dengan kode " + permohonan.Kode + " telah dibatalkan, silahkan cek permohonan pada aplikasi untuk lebih jelas",
	}

	ac.DB.Create(&simpanNotif)

	if user.PushToken != "" {
		notif := utils.Notification{
			IdUser: user.ID,
			Title:  "Permohonan dibatalkan",
			Body:   "Permohonan anda dengan kode " + permohonan.Kode + " telah dibatalkan, silahkan cek permohona pada aplikasi untuk lebih jelas",
			Token:  user.PushToken,
		}

		utils.SendPushNotifications(&notif)
	}
}

func (ac *AndalalinController) SurveiKepuasan(ctx *gin.Context) {
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

	ac.DB.First(&andalalin, "id_andalalin = ?", id)
	ac.DB.First(&perlalin, "id_andalalin = ?", id)

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

		result := ac.DB.Create(&kepuasan)

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

		result := ac.DB.Create(&kepuasan)

		if result.Error != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Telah terjadi sesuatu"})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (ac *AndalalinController) CekSurveiKepuasan(ctx *gin.Context) {
	id := ctx.Param("id_andalalin")

	var survei models.SurveiKepuasan

	result := ac.DB.First(&survei, "id_andalalin", id)

	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Telah terjadi sesuatu"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (ac *AndalalinController) HasilSurveiKepuasan(ctx *gin.Context) {
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

	result := ac.DB.Where("tanggal_pelaksanaan LIKE ?", fmt.Sprintf("%%%s%%", utils.Bulan(nowTime.Month()))).Find(&survei)

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

func (ac *AndalalinController) GetPermohonanPemasanganLalin(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)

	config, _ := initializers.LoadConfig()

	accessUser := ctx.MustGet("accessUser").(string)

	claim, error := utils.ValidateToken(accessUser, config.AccessTokenPublicKey)
	if error != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": error.Error()})
		return
	}

	credential := claim.Credentials[repository.AndalalinGetCredential]

	if !credential {
		// Return status 403 and permission denied error message.
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": true,
			"msg":   "Permission denied",
		})
		return
	}

	var perlalin []models.Perlalin

	ac.DB.Order("tanggal_andalalin").Find(&perlalin)

	var respone []models.DaftarAndalalinResponse
	for _, s := range perlalin {
		if s.StatusAndalalin == "Pemasangan sedang dilakukan" && s.IdPetugas == currentUser.ID {
			respone = append(respone, models.DaftarAndalalinResponse{
				IdAndalalin:      s.IdAndalalin,
				Kode:             s.Kode,
				TanggalAndalalin: s.TanggalAndalalin,
				Nama:             s.NamaPemohon,
				Email:            s.EmailPemohon,
				Petugas:          s.NamaPetugas,
				JenisAndalalin:   s.JenisAndalalin,
				StatusAndalalin:  s.StatusAndalalin,
			})
		}
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "results": len(respone), "data": respone})
}

func (ac *AndalalinController) PemasanganPerlengkapanLaluLintas(ctx *gin.Context) {
	var payload *models.DataSurvey
	currentUser := ctx.MustGet("currentUser").(models.User)
	id := ctx.Param("id_andalalin")

	config, _ := initializers.LoadConfig()

	accessUser := ctx.MustGet("accessUser").(string)

	claim, error := utils.ValidateToken(accessUser, config.AccessTokenPublicKey)
	if error != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": error.Error()})
		return
	}

	credential := claim.Credentials[repository.AndalalinPemasanganCredential]

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

	var ticket1 models.TiketLevel1

	resultTiket1 := ac.DB.Find(&ticket1, "id_andalalin = ?", id)
	if resultTiket1.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Tiket tidak ditemukan"})
		return
	}

	var perlalin models.Perlalin
	resultsPerlalin := ac.DB.First(&perlalin, "id_andalalin = ?", id)

	if resultsPerlalin.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Tidak ditemukan"})
		return
	}

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

	survey := models.Pemasangan{
		IdAndalalin:       perlalin.IdAndalalin,
		IdTiketLevel1:     ticket1.IdTiketLevel1,
		IdPetugas:         currentUser.ID,
		Petugas:           currentUser.Name,
		EmailPetugas:      currentUser.Email,
		Lokasi:            payload.Data.Lokasi,
		Keterangan:        payload.Data.Keterangan,
		Foto1:             blobs["foto1"],
		Foto2:             blobs["foto2"],
		Foto3:             blobs["foto3"],
		Latitude:          payload.Data.Latitude,
		Longitude:         payload.Data.Longitude,
		WaktuPemasangan:   tanggal,
		TanggalPemasangan: nowTime.Format("15:04:05"),
	}

	result := ac.DB.Create(&survey)

	if result.Error != nil && strings.Contains(result.Error.Error(), "duplicate key value violates unique") {
		ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "Data survey sudah tersedia"})
		return
	} else if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Telah terjadi sesuatu"})
		return
	}

	perlalin.StatusAndalalin = "Pemasangan selesai"

	ac.DB.Save(&perlalin)

	ac.PemasanganSelesai(ctx, perlalin)
	ac.CloseTiketLevel1(ctx, perlalin.IdAndalalin)

	ctx.JSON(http.StatusCreated, gin.H{"status": "success"})
}

func (ac *AndalalinController) PemasanganSelesai(ctx *gin.Context, permohonan models.Perlalin) {
	data := utils.Pemasangan{
		Kode:    permohonan.Kode,
		Nama:    permohonan.NamaPemohon,
		Tlp:     permohonan.NomerPemohon,
		Jenis:   permohonan.JenisAndalalin,
		Status:  permohonan.StatusAndalalin,
		Subject: "Pemasangan selesai",
	}

	utils.SendEmailPemasangan(permohonan.EmailPemohon, &data)

	var user models.User
	resultUser := ac.DB.First(&user, "id = ?", permohonan.IdUser)
	if resultUser.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "User tidak ditemukan"})
		return
	}

	simpanNotif := models.Notifikasi{
		IdUser: user.ID,
		Title:  "Pemasangan selesai",
		Body:   "Permohonan anda dengan kode " + permohonan.Kode + " telah selesai pemasangan perlengkapan lalu lintas, silahkan cek permohonan pada aplikasi untuk lebih jelas",
	}

	ac.DB.Create(&simpanNotif)

	if user.PushToken != "" {
		notif := utils.Notification{
			IdUser: user.ID,
			Title:  "Pemasangan selesai",
			Body:   "Permohonan anda dengan kode " + permohonan.Kode + " telah selesai pemasangan perlengkapan lalu lintas, silahkan cek permohonan pada aplikasi untuk lebih jelas",
			Token:  user.PushToken,
		}

		utils.SendPushNotifications(&notif)
	}
}

func (ac *AndalalinController) GetAllPemasangan(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)
	config, _ := initializers.LoadConfig()

	accessUser := ctx.MustGet("accessUser").(string)

	claim, error := utils.ValidateToken(accessUser, config.AccessTokenPublicKey)
	if error != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": error.Error()})
		return
	}

	credential := claim.Credentials[repository.AndalalinPemasanganCredential]

	if !credential {
		// Return status 403 and permission denied error message.
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": true,
			"msg":   "Permission denied",
		})
		return
	}

	var pemasangan []models.Pemasangan

	results := ac.DB.Find(&pemasangan, "id_petugas = ?", currentUser.ID)

	if results.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": results.Error})
		return
	} else {
		var respone []models.DaftarAndalalinResponse
		for _, s := range pemasangan {
			var perlalin models.Perlalin

			ac.DB.First(&perlalin, "id_andalalin = ?", s.IdAndalalin)

			if perlalin.IdAndalalin != uuid.Nil {
				respone = append(respone, models.DaftarAndalalinResponse{
					IdAndalalin:      perlalin.IdAndalalin,
					Kode:             perlalin.Kode,
					TanggalAndalalin: perlalin.TanggalAndalalin,
					Nama:             perlalin.NamaPemohon,
					Email:            perlalin.EmailPemohon,
					Petugas:          perlalin.NamaPetugas,
					JenisAndalalin:   perlalin.JenisAndalalin,
					StatusAndalalin:  perlalin.StatusAndalalin,
				})
			}

		}
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "results": len(respone), "data": respone})
	}
}

func (ac *AndalalinController) GetPemasangan(ctx *gin.Context) {
	id := ctx.Param("id_andalalin")

	var pemasangan *models.Pemasangan

	result := ac.DB.First(&pemasangan, "id_andalalin = ?", id)
	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": pemasangan})
}
