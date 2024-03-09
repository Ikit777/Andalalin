package controllers

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"andalalin/initializers"
	"andalalin/models"
	"andalalin/repository"
	"andalalin/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	_ "time/tzdata"

	"github.com/lukasjarosch/go-docx"
)

type AndalalinController struct {
	DB *gorm.DB
}

func NewAndalalinController(DB *gorm.DB) AndalalinController {
	return AndalalinController{DB}
}

func findItem(array []string, target string) int {
	for i, value := range array {
		if value == target {
			return i
		}
	}
	return -1
}

// func generatePDF(htmlContent string) ([]byte, error) {
// 	ctx, cancel := chromedp.NewContext(context.Background())
// 	defer cancel()

// 	var pdfContent []byte
// 	err := chromedp.Run(ctx,
// 		chromedp.Navigate("about:blank"),
// 		chromedp.ActionFunc(func(ctx context.Context) error {
// 			frameTree, err := page.GetFrameTree().Do(ctx)
// 			if err != nil {
// 				return err
// 			}

// 			return page.SetDocumentContent(frameTree.Frame.ID, htmlContent).Do(ctx)
// 		}),
// 		chromedp.ActionFunc(func(ctx context.Context) error {
// 			err := chromedp.ActionFunc(func(ctx context.Context) error {
// 				buf, _, err := page.PrintToPDF().WithPaperHeight(11.7).WithPaperWidth(8.3).WithMarginBottom(1).WithMarginLeft(1).WithMarginRight(1).WithMarginTop(1).WithDisplayHeaderFooter(false).WithPrintBackground(false).Do(ctx)
// 				if err != nil {
// 					return err
// 				}
// 				pdfContent = buf
// 				return nil
// 			}).Do(ctx)
// 			return err
// 		}),
// 	)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return pdfContent, nil
// }

func generatePDF(htmlContent []byte) ([]byte, error) {
	cmd := exec.Command("xvfb-run", "--", "wkhtmltopdf", "--margin-top", "1in", "--margin-bottom", "1in", "--margin-left", "1in", "--margin-right", "1in", "-", "-") // "-" indicates stdin and stdout respectively
	cmd.Stdin = bytes.NewReader(htmlContent)
	var pdfBuffer bytes.Buffer
	cmd.Stdout = &pdfBuffer

	// Run the command
	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	pdfBytes := pdfBuffer.Bytes()
	return pdfBytes, nil
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

	kode := "andalalin/" + utils.Generate(6) + "/" + nowTime.Format("2006")
	tanggal := nowTime.Format("02") + " " + utils.Bulan(nowTime.Month()) + " " + nowTime.Format("2006")

	var path string

	if payload.Andalalin.Pemohon == "Perorangan" {
		path = "templates/templateTandaTerimaAndalalinPerorangan.html"
	} else {
		path = "templates/templateTandaTerimaAndalalinNonPerorangan.html"
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
		Tanggal:  tanggal,
		Waktu:    nowTime.Format("15:04:05"),
		Kode:     kode,
		Nama:     currentUser.Name,
		Instansi: payload.Andalalin.NamaPerusahaan,
		Nomor:    payload.Andalalin.NomerPemohon,
	}

	buffer := new(bytes.Buffer)
	if err = t.Execute(buffer, bukti); err != nil {
		log.Fatal("Eror saat membaca template:", err)
		return
	}

	pdfContent, err := generatePDF(buffer.Bytes())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	berkas := []models.BerkasPermohonan{}

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

			if http.DetectContentType(data) == "application/pdf" {
				berkas = append(berkas, models.BerkasPermohonan{Nama: key, Tipe: "Pdf", Status: "Selesai", Berkas: data})
			} else {
				berkas = append(berkas, models.BerkasPermohonan{Nama: key, Tipe: "Word", Status: "Selesai", Berkas: data})
			}
		}
	}

	berkas = append(berkas, models.BerkasPermohonan{Nama: "Tanda terima pendaftaran", Tipe: "Pdf", Status: "Selesai", Berkas: pdfContent})

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
		StatusJalan:     payload.Andalalin.StatusJalan,

		//Data Pemohon
		NikPemohon:          payload.Andalalin.NikPemohon,
		NamaPemohon:         currentUser.Name,
		EmailPemohon:        currentUser.Email,
		TempatLahirPemohon:  payload.Andalalin.TempatLahirPemohon,
		TanggalLahirPemohon: payload.Andalalin.TanggalLahirPemohon,
		NegaraPemohon:       "Indonesia",
		ProvinsiPemohon:     payload.Andalalin.ProvinsiPemohon,
		KabupatenPemohon:    payload.Andalalin.KabupatenPemohon,
		KecamatanPemohon:    payload.Andalalin.KecamatanPemohon,
		KelurahanPemohon:    payload.Andalalin.KelurahanPemohon,
		AlamatPemohon:       payload.Andalalin.AlamatPemohon,
		JenisKelaminPemohon: payload.Andalalin.JenisKelaminPemohon,
		NomerPemohon:        payload.Andalalin.NomerPemohon,
		JabatanPemohon:      payload.Andalalin.JabatanPemohon,

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

		//Data Konsultan
		NamaKonsultan:                  payload.Andalalin.NamaKonsultan,
		NegaraKonsultan:                "Indonesia",
		ProvinsiKonsultan:              payload.Andalalin.ProvinsiKonsultan,
		KabupatenKonsultan:             payload.Andalalin.KabupatenKonsultan,
		KecamatanKonsultan:             payload.Andalalin.KecamatanKonsultan,
		KelurahanKonsultan:             payload.Andalalin.KelurahanKonsultan,
		AlamatKonsultan:                payload.Andalalin.AlamatKonsultan,
		NomerKonsultan:                 payload.Andalalin.NomerKonsultan,
		EmailKonsultan:                 payload.Andalalin.EmailKonsultan,
		NamaPenyusunDokumen:            payload.Andalalin.NamaPenyusunDokumen,
		JenisKelaminPenyusunDokumen:    payload.Andalalin.JenisKelaminPenyusunDokumen,
		NegaraPenyusunDokumen:          "Indonesia",
		ProvinsiPenyusunDokumen:        payload.Andalalin.ProvinsiPenyusunDokumen,
		KabupatenPenyusunDokumen:       payload.Andalalin.KabupatenPenyusunDokumen,
		KecamatanPenyusunDokumen:       payload.Andalalin.KecamatanPenyusunDokumen,
		KelurahanPenyusunDokumen:       payload.Andalalin.KelurahanPenyusunDokumen,
		AlamatPenyusunDokumen:          payload.Andalalin.AlamatPenyusunDokumen,
		NomerSertifikatPenyusunDokumen: payload.Andalalin.NomerSertifikatPenyusunDokumen,
		KlasifikasiPenyusunDokumen:     payload.Andalalin.KlasifikasiPenyusunDokumen,

		//Data Kegiatan
		Aktivitas:         payload.Andalalin.Aktivitas,
		Peruntukan:        payload.Andalalin.Peruntukan,
		TotalLuasLahan:    payload.Andalalin.TotalLuasLahan + " m²",
		KriteriaKhusus:    payload.Andalalin.KriteriaKhusus,
		NilaiKriteria:     payload.Andalalin.NilaiKriteria,
		Terbilang:         payload.Andalalin.Terbilang,
		LokasiBangunan:    payload.Andalalin.LokasiBangunan,
		LatitudeBangunan:  payload.Andalalin.LatitudeBangunan,
		LongitudeBangunan: payload.Andalalin.LongitudeBangunan,
		NomerSKRK:         payload.Andalalin.NomerSKRK,
		TanggalSKRK:       payload.Andalalin.TanggalSKRK,
		Catatan:           payload.Andalalin.Catatan,

		BerkasPermohonan:       berkas,
		StatusBerkasPermohonan: "Baru",
	}

	result := ac.DB.Create(&permohonan)

	respone := &models.DaftarAndalalinResponse{
		IdAndalalin:      permohonan.IdAndalalin,
		Kode:             permohonan.Kode,
		TanggalAndalalin: permohonan.TanggalAndalalin,
		Nama:             permohonan.NamaPemohon,
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
				Body:   "Permohonan baru dengan kode " + permohonan.Kode + " telah tersedia",
			}

			ac.DB.Create(&simpanNotif)

			if users.PushToken != "" {
				notif := utils.Notification{
					IdUser: users.ID,
					Title:  "Permohonan baru",
					Body:   "Permohonan baru dengan kode " + permohonan.Kode + " telah tersedia",
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

	kode := "perlalin/" + utils.Generate(6) + "/" + nowTime.Format("2006")
	tanggal := nowTime.Format("02") + " " + utils.Bulan(nowTime.Month()) + " " + nowTime.Format("2006")

	t, err := template.ParseFiles("templates/templateTandaTerimaPerlalin.html")
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

	pdfContent, err := generatePDF(buffer.Bytes())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	berkas := []models.BerkasPermohonan{}

	foto := []models.DataFoto{}

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

			switch http.DetectContentType(data) {
			case "application/pdf":
				berkas = append(berkas, models.BerkasPermohonan{Nama: key, Tipe: "Pdf", Status: "Selesai", Berkas: data})
			case "application/vnd.openxmlformats-officedocument.wordprocessingml.document":
				berkas = append(berkas, models.BerkasPermohonan{Nama: key, Tipe: "Word", Status: "Selesai", Berkas: data})
			case "image/jpeg":
				foto = append(foto, models.DataFoto{Id: key, Foto: data})
			}
		}
	}

	berkas = append(berkas, models.BerkasPermohonan{Nama: "Tanda terima pendaftaran", Tipe: "Pdf", Status: "Selesai", Berkas: pdfContent})

	perlengkapan := []models.Perlengkapan{}

	for _, data := range payload.Perlalin.Perlengkapan {
		foto_perlengkapan := []models.Foto{}
		for _, foto := range foto {
			if foto.Id == data.IdPerlengkapan {
				foto_perlengkapan = append(foto_perlengkapan, foto.Foto)
			}
		}

		input := models.Perlengkapan{
			IdPerlengkapan:       data.IdPerlengkapan,
			StatusPerlengkapan:   "Pemeriksaan",
			KategoriUtama:        data.KategoriUtama,
			KategoriPerlengkapan: data.KategoriPerlengkapan,
			JenisPerlengkapan:    data.JenisPerlengkapan,
			GambarPerlengkapan:   data.GambarPerlengkapan,
			LokasiPemasangan:     data.LokasiPemasangan,
			FotoLokasi:           foto_perlengkapan,
			LatitudePemasangan:   data.LatitudePemasangan,
			LongitudePemasangan:  data.LongitudePemasangan,
			Detail:               data.Detail,
			Alasan:               data.Alasan,
		}
		perlengkapan = append(perlengkapan, input)
	}

	permohonan := models.Perlalin{
		IdUser:              currentUser.ID,
		JenisAndalalin:      "Perlengkapan lalu lintas",
		Perlengkapan:        perlengkapan,
		Kode:                kode,
		NikPemohon:          payload.Perlalin.NikPemohon,
		NamaPemohon:         currentUser.Name,
		EmailPemohon:        currentUser.Email,
		TempatLahirPemohon:  payload.Perlalin.TempatLahirPemohon,
		TanggalLahirPemohon: payload.Perlalin.TanggalLahirPemohon,
		NegaraPemohon:       "Indonesia",
		ProvinsiPemohon:     payload.Perlalin.ProvinsiPemohon,
		KabupatenPemohon:    payload.Perlalin.KabupatenPemohon,
		KecamatanPemohon:    payload.Perlalin.KecamatanPemohon,
		KelurahanPemohon:    payload.Perlalin.KelurahanPemohon,
		AlamatPemohon:       payload.Perlalin.AlamatPemohon,
		JenisKelaminPemohon: payload.Perlalin.JenisKelaminPemohon,
		NomerPemohon:        payload.Perlalin.NomerPemohon,
		WaktuAndalalin:      nowTime.Format("15:04:05"),
		TanggalAndalalin:    tanggal,
		Catatan:             payload.Perlalin.Catatan,
		StatusAndalalin:     "Cek persyaratan",
		BerkasPermohonan:    berkas,
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
				Body:   "Permohonan baru dengan kode " + permohonan.Kode + " telah tersedia",
			}

			ac.DB.Create(&simpanNotif)

			if users.PushToken != "" {
				notif := utils.Notification{
					IdUser: users.ID,
					Title:  "Permohonan baru",
					Body:   "Permohonan baru dengan kode " + permohonan.Kode + " telah tersedia",
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

	if result.Error != nil {
		fmt.Println(result.Error)

		if strings.Contains(strings.ToLower(result.Error.Error()), "unique constraint") {
			ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "Data is exist"})
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "An error occurred on the server. Please try again later"})
			return
		}
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
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": results.Error})
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
	var payload *models.Pertimbangan
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
		andalalin.StatusAndalalin = "Permohonan ditunda"
		andalalin.PertimbanganPenundaan = payload.Pertimbangan
		ac.DB.Save(&andalalin)

		var user models.User
		resultUser := ac.DB.First(&user, "id = ?", andalalin.IdUser)
		if resultUser.Error != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "User tidak ditemukan"})
			return
		}

		simpanNotif := models.Notifikasi{
			IdUser: user.ID,
			Title:  "Permohonan ditunda",
			Body:   "Permohonan anda dengan kode " + andalalin.Kode + " telah ditunda",
		}

		ac.DB.Create(&simpanNotif)

		if user.PushToken != "" {
			notif := utils.Notification{
				IdUser: user.ID,
				Title:  "Permohonan ditunda",
				Body:   "Permohonan anda dengan kode " + andalalin.Kode + " telah ditunda",
				Token:  user.PushToken,
			}

			utils.SendPushNotifications(&notif)

		}
	}

	if perlalin.IdAndalalin != uuid.Nil {
		ac.CloseTiketLevel1(ctx, perlalin.IdAndalalin)
		perlalin.StatusAndalalin = "Permohonan ditunda"
		perlalin.PertimbanganPenundaan = payload.Pertimbangan
		ac.DB.Save(&perlalin)

		var user models.User
		resultUser := ac.DB.First(&user, "id = ?", perlalin.IdUser)
		if resultUser.Error != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "User tidak ditemukan"})
			return
		}

		simpanNotif := models.Notifikasi{
			IdUser: user.ID,
			Title:  "Permohonan ditunda",
			Body:   "Permohonan anda dengan kode " + perlalin.Kode + " telah ditunda",
		}

		ac.DB.Create(&simpanNotif)

		if user.PushToken != "" {
			notif := utils.Notification{
				IdUser: user.ID,
				Title:  "Permohonan ditunda",
				Body:   "Permohonan anda dengan kode " + perlalin.Kode + " telah ditunda",
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

	var andalalin models.Perlalin
	var perlalin models.Perlalin

	ac.DB.First(&andalalin, "id_andalalin = ?", id)
	ac.DB.First(&perlalin, "id_andalalin = ?", id)

	if andalalin.IdAndalalin != uuid.Nil {
		ac.CloseTiketLevel1(ctx, andalalin.IdAndalalin)
		andalalin.StatusAndalalin = "Cek persyaratan"
		andalalin.PertimbanganPenundaan = ""
		ac.DB.Save(&andalalin)

		var user models.User
		resultUser := ac.DB.First(&user, "id = ?", andalalin.IdUser)
		if resultUser.Error != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "User tidak ditemukan"})
			return
		}

		simpanNotif := models.Notifikasi{
			IdUser: user.ID,
			Title:  "Permohonan dilanjutkan",
			Body:   "Permohonan anda dengan kode " + andalalin.Kode + " telah dilanjutkan",
		}

		ac.DB.Create(&simpanNotif)

		if user.PushToken != "" {
			notif := utils.Notification{
				IdUser: user.ID,
				Title:  "Permohonan dilanjutkan",
				Body:   "Permohonan anda dengan kode " + andalalin.Kode + " telah dilanjutkan",
				Token:  user.PushToken,
			}

			utils.SendPushNotifications(&notif)

		}
	}

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
			Body:   "Permohonan anda dengan kode " + perlalin.Kode + " telah dilanjutkan",
		}

		ac.DB.Create(&simpanNotif)

		if user.PushToken != "" {
			notif := utils.Notification{
				IdUser: user.ID,
				Title:  "Permohonan dilanjutkan",
				Body:   "Permohonan anda dengan kode " + perlalin.Kode + " telah dilanjutkan",
				Token:  user.PushToken,
			}

			utils.SendPushNotifications(&notif)

		}
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (ac *AndalalinController) TolakPermohonan(ctx *gin.Context) {
	id := ctx.Param("id_andalalin")

	var payload *models.Pertimbangan
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
		andalalin.PertimbanganPenolakan = payload.Pertimbangan
		ac.DB.Save(&andalalin)

		var user models.User
		resultUser := ac.DB.First(&user, "id = ?", andalalin.IdUser)
		if resultUser.Error != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "User tidak ditemukan"})
			return
		}

		simpanNotif := models.Notifikasi{
			IdUser: user.ID,
			Title:  "Permohonan ditolak",
			Body:   "Permohonan anda dengan kode " + andalalin.Kode + " telah ditolak",
		}

		ac.DB.Create(&simpanNotif)

		if user.PushToken != "" {
			notif := utils.Notification{
				IdUser: user.ID,
				Title:  "Permohonan ditolak",
				Body:   "Permohonan anda dengan kode " + andalalin.Kode + " telah ditolak",
				Token:  user.PushToken,
			}

			utils.SendPushNotifications(&notif)

		}
	}

	if perlalin.IdAndalalin != uuid.Nil {
		ac.CloseTiketLevel1(ctx, perlalin.IdAndalalin)
		perlalin.StatusAndalalin = "Permohonan ditolak"
		perlalin.PertimbanganPenolakan = payload.Pertimbangan
		ac.DB.Save(&perlalin)

		var user models.User
		resultUser := ac.DB.First(&user, "id = ?", perlalin.IdUser)
		if resultUser.Error != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "User tidak ditemukan"})
			return
		}

		simpanNotif := models.Notifikasi{
			IdUser: user.ID,
			Title:  "Permohonan ditolak",
			Body:   "Permohonan anda dengan kode " + perlalin.Kode + " telah ditolak",
		}

		ac.DB.Create(&simpanNotif)

		if user.PushToken != "" {
			notif := utils.Notification{
				IdUser: user.ID,
				Title:  "Permohonan ditolak",
				Body:   "Permohonan anda dengan kode " + perlalin.Kode + " telah ditolak",
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
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Tidak ditemukan"})
		return
	} else {
		var respone []models.DaftarAndalalinResponse
		for _, s := range andalalin {
			respone = append(respone, models.DaftarAndalalinResponse{
				IdAndalalin:      s.IdAndalalin,
				Kode:             s.Kode,
				TanggalAndalalin: s.TanggalAndalalin,
				Nama:             s.NamaPemohon,
				Email:            s.EmailPemohon,
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

	var master models.DataMaster

	ac.DB.Select("persyaratan").First(&master)

	var persyaratan_andalalin []string
	for _, persyaratan := range master.Persyaratan.PersyaratanAndalalin {
		if persyaratan.Bangkitan == andalalin.Bangkitan {
			persyaratan_andalalin = append(persyaratan_andalalin, persyaratan.Persyaratan)
		}
	}

	var persyaratan_perlalin []string
	for _, persyaratan := range master.Persyaratan.PersyaratanPerlalin {
		persyaratan_perlalin = append(persyaratan_perlalin, persyaratan.Persyaratan)

	}

	var persyaratan_dishub []string
	var berkas_dishub []string
	for _, dokumen := range andalalin.BerkasPermohonan {
		index := findItem(persyaratan_andalalin, dokumen.Nama)

		if index != -1 {
			persyaratan_dishub = append(persyaratan_dishub, dokumen.Nama)
		} else {
			berkas_dishub = append(berkas_dishub, dokumen.Nama)
		}
	}

	var persyaratan_user []string
	var berkas_user []string
	for _, dokumen := range andalalin.BerkasPermohonan {
		if dokumen.Status == "Selesai" {
			index := findItem(persyaratan_andalalin, dokumen.Nama)

			if index != -1 {
				persyaratan_user = append(persyaratan_user, dokumen.Nama)
			} else {
				berkas_user = append(berkas_user, dokumen.Nama)
			}
		}
	}

	var kelengkapan_user []models.KelengkapanTidakSesuaiResponse
	for _, dokumen := range andalalin.KelengkapanTidakSesuai {
		if dokumen.Role == "User" {
			kelengkapan_user = append(kelengkapan_user, models.KelengkapanTidakSesuaiResponse{Dokumen: dokumen.Dokumen, Tipe: dokumen.Tipe})
		}
	}

	var kelengkapan_dushub []models.KelengkapanTidakSesuaiResponse
	for _, dokumen := range andalalin.KelengkapanTidakSesuai {
		if dokumen.Role == "Dishub" {
			kelengkapan_dushub = append(kelengkapan_dushub, models.KelengkapanTidakSesuaiResponse{Dokumen: dokumen.Dokumen, Tipe: dokumen.Tipe})
		}
	}

	var berkas_persyaratan_perlalin []string
	var berkas_permohonan_perlalin []string
	for _, dokumen := range perlalin.BerkasPermohonan {
		index := findItem(persyaratan_perlalin, dokumen.Nama)

		if index != -1 {
			berkas_persyaratan_perlalin = append(berkas_persyaratan_perlalin, dokumen.Nama)
		} else {
			berkas_permohonan_perlalin = append(berkas_permohonan_perlalin, dokumen.Nama)
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
				NegaraProyek:    andalalin.NegaraProyek,
				ProvinsiProyek:  andalalin.ProvinsiProyek,
				KabupatenProyek: andalalin.KabupatenProyek,
				KecamatanProyek: andalalin.KecamatanProyek,
				KelurahanProyek: andalalin.KelurahanProyek,
				AlamatProyek:    andalalin.AlamatProyek,

				//Data Pemohon
				NikPemohon:     andalalin.NikPemohon,
				EmailPemohon:   andalalin.EmailPemohon,
				NomerPemohon:   andalalin.NomerPemohon,
				NamaPemohon:    andalalin.NamaPemohon,
				JabatanPemohon: andalalin.JabatanPemohon,

				//Data perusahaan
				NamaPerusahaan: andalalin.NamaPerusahaan,

				//Data Konsultan
				NamaKonsultan:                  andalalin.NamaKonsultan,
				NomerKonsultan:                 andalalin.NomerKonsultan,
				EmailKonsultan:                 andalalin.EmailKonsultan,
				NamaPenyusunDokumen:            andalalin.NamaPenyusunDokumen,
				NomerSertifikatPenyusunDokumen: andalalin.NomerSertifikatPenyusunDokumen,
				KlasifikasiPenyusunDokumen:     andalalin.KlasifikasiPenyusunDokumen,

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
				PertimbanganPenundaan:  andalalin.PertimbanganPenundaan,
				PertimbanganPenolakan:  andalalin.PertimbanganPenolakan,

				BerkasPermohonan:      berkas_user,
				PersyaratanPermohonan: persyaratan_user,

				KelengkapanTidakSesuai: kelengkapan_user,
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
				StatusJalan:     andalalin.StatusJalan,

				//Data Pemohon
				NikPemohon:          andalalin.NikPemohon,
				NamaPemohon:         andalalin.NamaPemohon,
				EmailPemohon:        andalalin.EmailPemohon,
				TempatLahirPemohon:  andalalin.TempatLahirPemohon,
				TanggalLahirPemohon: andalalin.TanggalLahirPemohon,
				NegaraPemohon:       andalalin.NegaraPemohon,
				ProvinsiPemohon:     andalalin.ProvinsiPemohon,
				KabupatenPemohon:    andalalin.KabupatenPemohon,
				KecamatanPemohon:    andalalin.KecamatanPemohon,
				KelurahanPemohon:    andalalin.KelurahanPemohon,
				AlamatPemohon:       andalalin.AlamatPemohon,
				JenisKelaminPemohon: andalalin.JenisKelaminPemohon,
				NomerPemohon:        andalalin.NomerPemohon,
				JabatanPemohon:      andalalin.JabatanPemohon,

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

				//Data Konsultan
				NamaKonsultan:                  andalalin.NamaKonsultan,
				NegaraKonsultan:                andalalin.NegaraKonsultan,
				ProvinsiKonsultan:              andalalin.ProvinsiKonsultan,
				KabupatenKonsultan:             andalalin.KabupatenKonsultan,
				KecamatanKonsultan:             andalalin.KecamatanKonsultan,
				KelurahanKonsultan:             andalalin.KelurahanKonsultan,
				AlamatKonsultan:                andalalin.AlamatKonsultan,
				NomerKonsultan:                 andalalin.NomerKonsultan,
				EmailKonsultan:                 andalalin.EmailKonsultan,
				NamaPenyusunDokumen:            andalalin.NamaPenyusunDokumen,
				JenisKelaminPenyusunDokumen:    andalalin.JenisKelaminPenyusunDokumen,
				NegaraPenyusunDokumen:          andalalin.NegaraPenyusunDokumen,
				ProvinsiPenyusunDokumen:        andalalin.ProvinsiPenyusunDokumen,
				KabupatenPenyusunDokumen:       andalalin.KabupatenPenyusunDokumen,
				KecamatanPenyusunDokumen:       andalalin.KecamatanPenyusunDokumen,
				KelurahanPenyusunDokumen:       andalalin.KelurahanPenyusunDokumen,
				AlamatPenyusunDokumen:          andalalin.AlamatPenyusunDokumen,
				NomerSertifikatPenyusunDokumen: andalalin.NomerSertifikatPenyusunDokumen,
				KlasifikasiPenyusunDokumen:     andalalin.KlasifikasiPenyusunDokumen,

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

				BerkasPermohonan:      berkas_dishub,
				PersyaratanPermohonan: persyaratan_dishub,

				PersyaratanTidakSesuai: andalalin.PersyaratanTidakSesuai,

				HasilAsistensiDokumen:   andalalin.HasilAsistensiDokumen,
				CatatanAsistensiDokumen: andalalin.CatatanAsistensiDokumen,

				//Data Pemeriksaan Surat Persetujuan
				HasilPemeriksaan:   andalalin.HasilPemeriksaan,
				CatatanPemeriksaan: andalalin.CatatanPemeriksaan,

				//Data Pertimbangan
				PertimbanganPenundaan: andalalin.PertimbanganPenundaan,
				PertimbanganPenolakan: andalalin.PertimbanganPenolakan,

				KelengkapanTidakSesuai: kelengkapan_dushub,
			}
			ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": data})
		}
	}

	var perlengkapan []models.PerlengkapanResponse

	for _, data := range perlalin.Perlengkapan {
		perlengkapan = append(perlengkapan, models.PerlengkapanResponse{IdPerlengkapan: data.IdPerlengkapan, StatusPerlengkapan: data.StatusPerlengkapan, JenisPerlengkapan: data.JenisPerlengkapan, GambarPerlengkapan: data.GambarPerlengkapan, LokasiPemasangan: data.LokasiPemasangan})
	}

	if perlalin.IdAndalalin != uuid.Nil {
		if currentUser.Role == "User" {
			dataUser := models.PerlalinResponseUser{
				//Data Permohonan
				IdAndalalin:      perlalin.IdAndalalin,
				JenisAndalalin:   perlalin.JenisAndalalin,
				Perlengkapan:     perlengkapan,
				Kode:             perlalin.Kode,
				WaktuAndalalin:   perlalin.WaktuAndalalin,
				TanggalAndalalin: perlalin.TanggalAndalalin,
				StatusAndalalin:  perlalin.StatusAndalalin,

				//Data Pemohon
				NamaPemohon:  perlalin.NamaPemohon,
				NikPemohon:   perlalin.NikPemohon,
				EmailPemohon: perlalin.EmailPemohon,
				NomerPemohon: perlalin.NomerPemohon,

				PersyaratanTidakSesuai: perlalin.PersyaratanTidakSesuai,

				PersyaratanPermohonan: berkas_persyaratan_perlalin,
				BerkasPermohonan:      berkas_permohonan_perlalin,

				//Catatan
				Catatan: perlalin.Catatan,

				PertimbanganPembatalan: perlalin.PertimbanganPembatalan,
				PertimbanganPenolakan:  perlalin.PertimbanganPenolakan,
				PertimbanganPenundaan:  perlalin.PertimbanganPenundaan,
			}

			ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": dataUser})
		} else {
			data := models.PerlalinResponse{
				//Data Permohonan
				IdAndalalin:      perlalin.IdAndalalin,
				JenisAndalalin:   perlalin.JenisAndalalin,
				Perlengkapan:     perlengkapan,
				Kode:             perlalin.Kode,
				WaktuAndalalin:   perlalin.WaktuAndalalin,
				TanggalAndalalin: perlalin.TanggalAndalalin,
				StatusAndalalin:  perlalin.StatusAndalalin,

				//Data Pemohon
				NikPemohon:          perlalin.NikPemohon,
				NamaPemohon:         perlalin.NamaPemohon,
				EmailPemohon:        perlalin.EmailPemohon,
				TempatLahirPemohon:  perlalin.TempatLahirPemohon,
				TanggalLahirPemohon: perlalin.TanggalLahirPemohon,
				NegaraPemohon:       perlalin.NegaraPemohon,
				ProvinsiPemohon:     perlalin.ProvinsiPemohon,
				KabupatenPemohon:    perlalin.KabupatenPemohon,
				KecamatanPemohon:    perlalin.KecamatanPemohon,
				KelurahanPemohon:    perlalin.KelurahanPemohon,
				AlamatPemohon:       perlalin.AlamatPemohon,
				JenisKelaminPemohon: perlalin.JenisKelaminPemohon,
				NomerPemohon:        perlalin.NomerPemohon,

				PersyaratanTidakSesuai: perlalin.PersyaratanTidakSesuai,
				IdPetugas:              perlalin.IdPetugas,
				NamaPetugas:            perlalin.NamaPetugas,
				EmailPetugas:           perlalin.EmailPetugas,
				StatusTiketLevel2:      status,

				PersyaratanPermohonan: berkas_persyaratan_perlalin,
				BerkasPermohonan:      berkas_permohonan_perlalin,

				Catatan: perlalin.Catatan,

				PertimbanganPembatalan: perlalin.PertimbanganPembatalan,
				PertimbanganPenolakan:  perlalin.PertimbanganPenolakan,
				PertimbanganPenundaan:  perlalin.PertimbanganPenundaan,
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
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Tidak ditemukan"})
		return
	} else {
		var respone []models.DaftarAndalalinResponse
		for _, s := range andalalin {
			respone = append(respone, models.DaftarAndalalinResponse{
				IdAndalalin:      s.IdAndalalin,
				Kode:             s.Kode,
				TanggalAndalalin: s.TanggalAndalalin,
				Nama:             s.NamaPemohon,
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
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Tidak ditemukan"})
		return
	} else {
		var respone []models.DaftarAndalalinResponse
		for _, s := range andalalin {
			respone = append(respone, models.DaftarAndalalinResponse{
				IdAndalalin:      s.IdAndalalin,
				Kode:             s.Kode,
				TanggalAndalalin: s.TanggalAndalalin,
				Nama:             s.NamaPemohon,
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
		for _, item := range andalalin.BerkasPermohonan {
			if item.Nama == dokumen {
				docs = item.Berkas
				tipe = item.Tipe
				break
			}
		}
	}

	if perlalin.IdAndalalin != uuid.Nil {
		for _, item := range perlalin.BerkasPermohonan {
			if item.Nama == dokumen {
				docs = item.Berkas
				tipe = item.Tipe
				break
			}
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "tipe": tipe, "data": docs})
}

func (ac *AndalalinController) UpdateBerkas(ctx *gin.Context) {
	config, _ := initializers.LoadConfig()

	accessUser := ctx.MustGet("accessUser").(string)

	claim, error := utils.ValidateToken(accessUser, config.AccessTokenPublicKey)
	if error != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": error.Error()})
		return
	}

	credential := claim.Credentials[repository.AndalalinBerkasCredential]

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
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Tidak ditemukan"})
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

				if andalalin.StatusAndalalin == "Kelengkapan tidak terpenuhi" {
					if andalalin.KelengkapanTidakSesuai != nil {
						for i, kelengkapan := range andalalin.KelengkapanTidakSesuai {
							if kelengkapan.Dokumen == key {
								andalalin.KelengkapanTidakSesuai = append(andalalin.KelengkapanTidakSesuai[:i], andalalin.KelengkapanTidakSesuai[i+1:]...)
								break
							}
						}
					}
				}

				var berkas_permohonan []string

				for _, berkas := range andalalin.BerkasPermohonan {
					berkas_permohonan = append(berkas_permohonan, berkas.Nama)
				}

				index := findItem(berkas_permohonan, key)

				if index != -1 {
					for i := range andalalin.BerkasPermohonan {
						if andalalin.BerkasPermohonan[i].Nama == key {
							andalalin.BerkasPermohonan[i].Berkas = data
							break
						}
					}
				} else {
					if http.DetectContentType(data) == "application/pdf" {
						andalalin.BerkasPermohonan = append(andalalin.BerkasPermohonan, models.BerkasPermohonan{Nama: key, Tipe: "Pdf", Status: "Selesai", Berkas: data})
					} else {
						andalalin.BerkasPermohonan = append(andalalin.BerkasPermohonan, models.BerkasPermohonan{Nama: key, Tipe: "Word", Status: "Selesai", Berkas: data})
					}
				}
			}
		}

		if andalalin.StatusAndalalin == "Kelengkapan tidak terpenuhi" {
			if len(andalalin.KelengkapanTidakSesuai) == 0 {
				andalalin.StatusAndalalin = "Cek kelengkapan akhir"
			}
		} else {
			andalalin.StatusAndalalin = "Cek persyaratan"
		}

		andalalin.StatusBerkasPermohonan = "Revisi"

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

				for i, berkas := range perlalin.BerkasPermohonan {
					if berkas.Nama == key {
						perlalin.BerkasPermohonan[i].Berkas = data
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

func (ac *AndalalinController) UpdateLokasiProyek(ctx *gin.Context) {
	var payload *models.PerbaruiLokasi
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

	credential := claim.Credentials[repository.AndalalinTindakLanjut]

	if !credential {
		// Return status 403 and permission denied error message.
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": true,
			"msg":   "Permission denied",
		})
		return
	}

	id := ctx.Param("id_andalalin")

	var andalalin *models.Andalalin

	resultsAndalalin := ac.DB.First(&andalalin, "id_andalalin = ?", id)

	if resultsAndalalin.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Tidak ditemukan"})
		return
	}

	if andalalin.IdAndalalin != uuid.Nil {
		andalalin.LokasiBangunan = payload.Lokasi
		andalalin.LatitudeBangunan = payload.Latitude
		andalalin.LongitudeBangunan = payload.Longitude

		ac.DB.Save(&andalalin)
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

	credential := claim.Credentials[repository.AndalalinBerkasCredential]

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
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Tidak ditemukan"})
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

			for i, item := range andalalin.BerkasPermohonan {
				if item.Nama == "Checklist administrasi" {
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
					andalalin.BerkasPermohonan[itemIndex].Berkas = data
				}
			}

			andalalin.BerkasPermohonan[itemIndex].Status = "Selesai"
			if andalalin.PersyaratanTidakSesuai != nil {
				andalalin.StatusAndalalin = "Persyaratan tidak terpenuhi"

				var user models.User
				resultUser := ac.DB.First(&user, "id = ?", andalalin.IdUser)
				if resultUser.Error != nil {
					ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "User tidak ditemukan"})
					return
				}

				simpanNotif := models.Notifikasi{
					IdUser: user.ID,
					Title:  "Persyaratan tidak terpenuhi",
					Body:   "Permohonan anda dengan kode " + andalalin.Kode + " terdapat persyaratan yang tidak terpenuhi",
				}

				ac.DB.Create(&simpanNotif)

				if user.PushToken != "" {
					notif := utils.Notification{
						IdUser: user.ID,
						Title:  "Persyaratan tidak terpenuhi",
						Body:   "Permohonan anda dengan kode " + andalalin.Kode + " terdapat persyaratan yang tidak terpenuhi",
						Token:  user.PushToken,
					}

					utils.SendPushNotifications(&notif)

				}
			} else {
				andalalin.StatusAndalalin = "Persyaratan terpenuhi"
				ac.ReleaseTicketLevel2(ctx, andalalin.IdAndalalin, andalalin.IdAndalalin)
			}
		}

		if dokumen == "Surat pernyataan kesanggupan" {
			itemPernyataan := -1

			for i, item := range andalalin.BerkasPermohonan {
				if item.Nama == "Surat pernyataan kesanggupan (word)" {
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
						andalalin.BerkasPermohonan[itemPernyataan].Berkas = data
					} else {
						andalalin.BerkasPermohonan = append(andalalin.BerkasPermohonan, models.BerkasPermohonan{Status: "Selesai", Nama: "Surat pernyataan kesanggupan (pdf)", Tipe: "Pdf", Berkas: data})
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
					andalalin.BerkasPermohonan = append(andalalin.BerkasPermohonan, models.BerkasPermohonan{Status: "Selesai", Nama: key, Tipe: "Pdf", Berkas: data})

				}
			}

			switch andalalin.Bangkitan {
			case "Bangkitan rendah":
				andalalin.StatusAndalalin = "Pembuatan surat keputusan"
			case "Bangkitan sedang":
				andalalin.StatusAndalalin = "Pembuatan surat keputusan"
			case "Bangkitan tinggi":
			}
		}

		if dokumen == "Penyusun dokumen" {
			itemIndex := -1

			for i, item := range andalalin.BerkasPermohonan {
				if item.Nama == "Penyusun dokumen analsis dampak lalu lintas" {
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
					andalalin.BerkasPermohonan[itemIndex].Berkas = data
				}
			}

			andalalin.BerkasPermohonan[itemIndex].Status = "Selesai"
			andalalin.StatusAndalalin = "Pemeriksaan dokumen andalalin"
		}

		if dokumen == "Catatan asistensi dokumen" {
			itemIndex := -1

			for i, item := range andalalin.BerkasPermohonan {
				if item.Nama == "Catatan asistensi dokumen analisis dampak lalu lintas" {
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
					andalalin.BerkasPermohonan[itemIndex].Berkas = data
				}
			}

			andalalin.BerkasPermohonan[itemIndex].Status = "Selesai"
			andalalin.StatusAndalalin = andalalin.HasilAsistensiDokumen
		}

		if dokumen == "Surat keputusan persetujuan teknis andalalin" {
			itenKeputusan := -1

			for i, item := range andalalin.BerkasPermohonan {
				if item.Nama == "Surat keputusan persetujuan teknis andalalin" {
					itenKeputusan = i
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
					andalalin.BerkasPermohonan[itenKeputusan].Berkas = data
				}
			}

			andalalin.StatusAndalalin = "Cek kelengkapan akhir"
			ac.CloseTiketLevel2(ctx, andalalin.IdAndalalin)
		}

		if dokumen == "Checklist kelengkapan akhir" {
			itemKelengkapan := -1

			for i, item := range andalalin.BerkasPermohonan {
				if item.Nama == "Checklist kelengkapan akhir" {
					itemKelengkapan = i
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
					andalalin.BerkasPermohonan[itemKelengkapan].Berkas = data
				}
			}

			andalalin.BerkasPermohonan[itemKelengkapan].Status = "Selesai"

			if andalalin.KelengkapanTidakSesuai != nil {
				andalalin.StatusAndalalin = "Kelengkapan tidak terpenuhi"
			} else {
				andalalin.StatusAndalalin = "Permohonan selesai"

				itenKeputusan := -1

				for i, item := range andalalin.BerkasPermohonan {
					if item.Nama == "Surat keputusan persetujuan teknis andalalin" {
						itenKeputusan = i
						break
					}
				}

				andalalin.BerkasPermohonan[itenKeputusan].Status = "Selesai"

				var user models.User
				resultUser := ac.DB.First(&user, "id = ?", andalalin.IdUser)
				if resultUser.Error != nil {
					ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "User tidak ditemukan"})
					return
				}

				simpanNotif := models.Notifikasi{
					IdUser: user.ID,
					Title:  "Permohonan selesai",
					Body:   "Permohonan anda dengan kode " + andalalin.Kode + " telah selesai",
				}

				ac.DB.Create(&simpanNotif)

				if user.PushToken != "" {
					notif := utils.Notification{
						IdUser: user.ID,
						Title:  "Permohonan selesai",
						Body:   "Permohonan anda dengan kode " + andalalin.Kode + " telah selesai",
						Token:  user.PushToken,
					}

					utils.SendPushNotifications(&notif)

				}
			}
		}

		if dokumen == "Dokumen andalalin" {
			itemWord := -1
			itemPdf := -1

			for i, item := range andalalin.BerkasPermohonan {
				if item.Nama == "Dokumen hasil analisis dampak lalu lintas (word)" {
					itemWord = i
					break
				}
			}

			for i, item := range andalalin.BerkasPermohonan {
				if item.Nama == "Dokumen hasil analisis dampak lalu lintas (pdf)" {
					itemPdf = i
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
					if key == "Dokumen hasil analisis dampak lalu lintas (word)" {
						andalalin.BerkasPermohonan[itemWord].Berkas = data
					} else {
						andalalin.BerkasPermohonan[itemPdf].Berkas = data
					}

				}
			}

			andalalin.StatusAndalalin = "Pemeriksaan dokumen andalalin"
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

	switch andalalin.Bangkitan {
	case "Bangkitan rendah":
		t, err := template.ParseFiles("templates/checklistAdministrasiBangkitanRendah.html")
		if err != nil {
			log.Fatal("Error reading the email template:", err)
			return
		}

		itemIndex := -1

		for i, item := range andalalin.BerkasPermohonan {
			if item.Nama == "Checklist administrasi" {
				itemIndex = i
				break
			}
		}

		administrasi := struct {
			Bangkitan   string
			Objek       string
			Lokasi      string
			Pengembang  string
			Nomor       string
			Diterima    string
			Pemeriksaan string
			Status      string
			Data        []models.DataAdministrasi
			Operator    string
			Nip         string
		}{
			Bangkitan:   "RENDAH",
			Objek:       andalalin.Jenis,
			Lokasi:      "Jalan " + andalalin.NamaJalan + ", Kelurahan " + andalalin.KelurahanProyek + ", Kecamatan " + andalalin.KecamatanProyek + ", Kabupaten " + andalalin.KabupatenProyek + ", Provinsi " + andalalin.ProvinsiProyek,
			Pengembang:  *andalalin.NamaPerusahaan,
			Nomor:       payload.NomorSurat + ", " + payload.TanggalSurat[0:2] + " " + utils.Month(payload.TanggalSurat[3:5]) + " " + payload.TanggalSurat[6:10],
			Diterima:    andalalin.TanggalAndalalin,
			Pemeriksaan: tanggal,
			Status:      andalalin.StatusBerkasPermohonan,
			Data:        payload.Data,
			Operator:    currentUser.Name,
			Nip:         *currentUser.NIP,
		}

		buffer := new(bytes.Buffer)
		if err = t.Execute(buffer, administrasi); err != nil {
			log.Fatal("Eror saat membaca template:", err)
			return
		}

		pdfContent, err := generatePDF(buffer.Bytes())
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if itemIndex != -1 {
			andalalin.BerkasPermohonan[itemIndex].Berkas = pdfContent
			andalalin.BerkasPermohonan[itemIndex].Status = "Menunggu"
		} else {
			andalalin.BerkasPermohonan = append(andalalin.BerkasPermohonan, models.BerkasPermohonan{Status: "Menunggu", Nama: "Checklist administrasi", Tipe: "Pdf", Berkas: pdfContent})
		}
	case "Bangkitan sedang":
		t, err := template.ParseFiles("templates/checklistAdministrasiBangkitanSedang.html")
		if err != nil {
			log.Fatal("Error reading the email template:", err)
			return
		}

		itemIndex := -1

		for i, item := range andalalin.BerkasPermohonan {
			if item.Nama == "Checklist administrasi" {
				itemIndex = i
				break
			}
		}

		administrasi := struct {
			Bangkitan   string
			Objek       string
			Lokasi      string
			Pengembang  string
			Penyusun    string
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
			Bangkitan:   "SEDANG",
			Objek:       andalalin.Jenis,
			Lokasi:      "Jalan " + andalalin.NamaJalan + ", Kelurahan " + andalalin.KelurahanProyek + ", Kecamatan " + andalalin.KecamatanProyek + ", Kabupaten " + andalalin.KabupatenProyek + ", Provinsi " + andalalin.ProvinsiProyek,
			Pengembang:  *andalalin.NamaPerusahaan,
			Penyusun:    *andalalin.NamaPenyusunDokumen,
			Sertifikat:  *andalalin.NomerSertifikatPenyusunDokumen,
			Klasifikasi: *andalalin.KlasifikasiPenyusunDokumen,
			Nomor:       payload.NomorSurat + ", " + payload.TanggalSurat[0:2] + " " + utils.Month(payload.TanggalSurat[3:5]) + " " + payload.TanggalSurat[6:10],
			Diterima:    andalalin.TanggalAndalalin,
			Pemeriksaan: tanggal,
			Status:      andalalin.StatusBerkasPermohonan,
			Data:        payload.Data,
			Operator:    currentUser.Name,
			Nip:         *currentUser.NIP,
		}

		buffer := new(bytes.Buffer)
		if err = t.Execute(buffer, administrasi); err != nil {
			log.Fatal("Eror saat membaca template:", err)
			return
		}

		pdfContent, err := generatePDF(buffer.Bytes())
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if itemIndex != -1 {
			andalalin.BerkasPermohonan[itemIndex].Berkas = pdfContent
			andalalin.BerkasPermohonan[itemIndex].Status = "Menunggu"
		} else {
			andalalin.BerkasPermohonan = append(andalalin.BerkasPermohonan, models.BerkasPermohonan{Status: "Menunggu", Nama: "Checklist administrasi", Tipe: "Pdf", Berkas: pdfContent})
		}
	case "Bangkitan tinggi":
		t, err := template.ParseFiles("templates/checklistAdministrasiBangkitanTinggi.html")
		if err != nil {
			log.Fatal("Error reading the email template:", err)
			return
		}

		itemIndex := -1

		for i, item := range andalalin.BerkasPermohonan {
			if item.Nama == "Checklist administrasi" {
				itemIndex = i
				break
			}
		}

		administrasi := struct {
			Bangkitan   string
			Objek       string
			Lokasi      string
			Pengembang  string
			Penyusun    string
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
			Bangkitan:   "TINGGI",
			Objek:       andalalin.Jenis,
			Lokasi:      "Jalan " + andalalin.NamaJalan + ", Kelurahan " + andalalin.KelurahanProyek + ", Kecamatan " + andalalin.KecamatanProyek + ", Kabupaten " + andalalin.KabupatenProyek + ", Provinsi " + andalalin.ProvinsiProyek,
			Pengembang:  *andalalin.NamaPerusahaan,
			Penyusun:    *andalalin.NamaPenyusunDokumen,
			Sertifikat:  *andalalin.NomerSertifikatPenyusunDokumen,
			Klasifikasi: *andalalin.KlasifikasiPenyusunDokumen,
			Nomor:       payload.NomorSurat + ", " + payload.TanggalSurat[0:2] + " " + utils.Month(payload.TanggalSurat[3:5]) + " " + payload.TanggalSurat[6:10],
			Diterima:    andalalin.TanggalAndalalin,
			Pemeriksaan: tanggal,
			Status:      andalalin.StatusBerkasPermohonan,
			Data:        payload.Data,
			Operator:    currentUser.Name,
			Nip:         *currentUser.NIP,
		}

		buffer := new(bytes.Buffer)
		if err = t.Execute(buffer, administrasi); err != nil {
			log.Fatal("Eror saat membaca template:", err)
			return
		}

		pdfContent, err := generatePDF(buffer.Bytes())
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if itemIndex != -1 {
			andalalin.BerkasPermohonan[itemIndex].Berkas = pdfContent
			andalalin.BerkasPermohonan[itemIndex].Status = "Menunggu"
		} else {
			andalalin.BerkasPermohonan = append(andalalin.BerkasPermohonan, models.BerkasPermohonan{Status: "Menunggu", Nama: "Checklist administrasi", Tipe: "Pdf", Berkas: pdfContent})
		}
	}

	if andalalin.PersyaratanTidakSesuai != nil {
		andalalin.PersyaratanTidakSesuai = nil
	}

	for _, item := range payload.Data {
		if item.Tidak != "" && item.Kebutuhan == "Wajib" {
			andalalin.PersyaratanTidakSesuai = append(andalalin.PersyaratanTidakSesuai, models.PersayaratanTidakSesuai{Persyaratan: item.Persyaratan, Tipe: item.Tipe})
		}
	}

	andalalin.Nomor = payload.NomorSurat
	andalalin.Tanggal = payload.TanggalSurat
	andalalin.StatusAndalalin = "Persetujuan administrasi"

	ac.DB.Save(&andalalin)

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (ac *AndalalinController) CheckAdministrasiPerlalin(ctx *gin.Context) {
	var payload *models.AdministrasiPerlalin
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

	var perlalin models.Perlalin

	result := ac.DB.First(&perlalin, "id_andalalin = ?", id)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Permohonan tidak ditemukan"})
		return
	}

	if perlalin.PersyaratanTidakSesuai != nil {
		perlalin.PersyaratanTidakSesuai = nil
	}

	for _, item := range payload.Data {
		if item.Tidak != "" && item.Kebutuhan == "Wajib" {
			perlalin.PersyaratanTidakSesuai = append(perlalin.PersyaratanTidakSesuai, models.PersayaratanTidakSesuai{Persyaratan: item.Persyaratan, Tipe: item.Tipe})
		}
	}

	if perlalin.PersyaratanTidakSesuai == nil {
		perlalin.StatusAndalalin = "Persyaratan terpenuhi"
	} else {
		perlalin.StatusAndalalin = "Persyaratan tidak terpenuhi"
	}

	ac.DB.Save(&perlalin)

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

func (ac *AndalalinController) PembuatanSuratPermohonan(ctx *gin.Context) {
	var payload *models.DataSuratPermohonan

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

	credential := claim.Credentials[repository.AndalalinPengajuanCredential]

	if !credential {
		// Return status 403 and permission denied error message.
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": true,
			"msg":   "Permission denied",
		})
		return
	}

	switch payload.Bangkitan {
	case "Bangkitan rendah":

		replaceMap := docx.PlaceholderMap{
			"_nama_":        payload.Nama,
			"_pemohon_":     *payload.Jabatan + " " + payload.Pengembang,
			"_pengembang_":  payload.Pengembang,
			"_klasifikasi_": "Rendah",
			"_jenis_":       payload.Jenis,
			"_proyek_":      payload.Proyek,
			"_jalan_":       payload.Jalan,
			"_kelurahan_":   payload.Kelurahan,
			"_kecamatan_":   payload.Kecamatan,
			"_kabupaten_":   payload.Kabupaten,
			"_provinsi_":    payload.Provinsi,
			"_status_":      payload.StatusJalan,
		}

		var docRendah string

		if payload.Pemohon == "Perorangan" {
			docRendah = "templates/suratPermohonanBangkitanRendahPerorangan.docx"
		} else {
			docRendah = "templates/suratPermohonanBangkitanRendahNonPerorangan.docx"
		}

		doc, err := docx.Open(docRendah)
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

		ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": docBytes})
	case "Bangkitan sedang":
		replaceMap := docx.PlaceholderMap{
			"_nama_":        payload.Nama,
			"_pemohon_":     *payload.Jabatan + " " + payload.Pengembang,
			"_pengembang_":  payload.Pengembang,
			"_klasifikasi_": "Sedang",
			"_jenis_":       payload.Jenis,
			"_proyek_":      payload.Proyek,
			"_jalan_":       payload.Jalan,
			"_kelurahan_":   payload.Kelurahan,
			"_kecamatan_":   payload.Kecamatan,
			"_kabupaten_":   payload.Kabupaten,
			"_provinsi_":    payload.Provinsi,
			"_status_":      payload.StatusJalan,
			"_konsultan_":   *payload.Konsultan,
		}

		var docRendah string

		if payload.Pemohon == "Perorangan" {
			docRendah = "templates/suratPermohonanBangkitanSedangPerorangan.docx"
		} else {
			docRendah = "templates/suratPermohonanBangkitanSedangNonPerorangan.docx"
		}

		doc, err := docx.Open(docRendah)
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

		ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": docBytes})
	case "Bangkitan tinggi":
		replaceMap := docx.PlaceholderMap{
			"_nama_":        payload.Nama,
			"_pemohon_":     *payload.Jabatan + " " + payload.Pengembang,
			"_pengembang_":  payload.Pengembang,
			"_klasifikasi_": "Tinggi",
			"_jenis_":       payload.Jenis,
			"_proyek_":      payload.Proyek,
			"_jalan_":       payload.Jalan,
			"_kelurahan_":   payload.Kelurahan,
			"_kecamatan_":   payload.Kecamatan,
			"_kabupaten_":   payload.Kabupaten,
			"_provinsi_":    payload.Provinsi,
			"_status_":      payload.StatusJalan,
			"_konsultan_":   *payload.Konsultan,
		}

		var docRendah string

		if payload.Pemohon == "Perorangan" {
			docRendah = "templates/suratPermohonanBangkitanTinggiPerorangan.docx"
		} else {
			docRendah = "templates/suratPermohonanBangkitanTinggiNonPerorangan.docx"
		}

		doc, err := docx.Open(docRendah)
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

		ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": docBytes})
	}
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

	switch andalalin.Bangkitan {
	case "Bangkitan rendah":
		listContent := ""
		for i, item := range payload.Kewajiban {
			if i == len(payload.Kewajiban)-1 {
				listContent += fmt.Sprint(i+1, ". ", item)
			} else {
				listContent += fmt.Sprint(i+1, ". ", item, "\n")
			}
		}

		replaceMap := docx.PlaceholderMap{
			"_nama_":       andalalin.NamaPemohon,
			"_jabatan_":    *andalalin.JabatanPemohon,
			"_alamat_":     andalalin.AlamatPemohon + ", Kelurahan " + andalalin.KelurahanPemohon + ", Kecamatan " + andalalin.KecamatanPemohon + ", Kabupaten " + andalalin.KabupatenPemohon + ", Provinsi " + andalalin.ProvinsiPemohon,
			"_pengembang_": *andalalin.NamaPerusahaan,
			"_nomor_":      andalalin.Nomor,
			"_tanggal_":    andalalin.Tanggal[0:2],
			"_bulan_":      utils.Month(andalalin.Tanggal[3:5]),
			"_tahun_":      andalalin.Tanggal[6:10],
			"_kegiatan_":   andalalin.JenisProyek + " " + andalalin.NamaProyek + ", di Jalan " + andalalin.NamaJalan + ", Kelurahan " + andalalin.KelurahanProyek + ", Kecamatan " + andalalin.KecamatanProyek + ", Kabupaten " + andalalin.KabupatenProyek + ", Provinsi " + andalalin.ProvinsiProyek,
			"_kewajiban_":  listContent,
		}

		var docRendah string

		if andalalin.Pemohon == "Perorangan" {
			docRendah = "templates/suratPernyataanKesanggupanBangkitanRendahPerorangan.docx"
		} else {
			docRendah = "templates/suratPernyataanKesanggupanBangkitanRendahNonPerorangan.docx"
		}

		doc, err := docx.Open(docRendah)
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

		andalalin.StatusAndalalin = "Menunggu surat pernyataan"

		itemIndex := -1

		for i, item := range andalalin.BerkasPermohonan {
			if item.Nama == "Surat pernyataan kesanggupan (word)" {
				itemIndex = i
				break
			}
		}

		if itemIndex != -1 {
			andalalin.BerkasPermohonan[itemIndex].Berkas = docBytes
		} else {
			andalalin.BerkasPermohonan = append(andalalin.BerkasPermohonan, models.BerkasPermohonan{Status: "Selesai", Nama: "Surat pernyataan kesanggupan (word)", Tipe: "Word", Berkas: docBytes})
		}
	case "Bangkitan sedang":
		listContent := ""
		for i, item := range payload.Kewajiban {
			if i == len(payload.Kewajiban)-1 {
				listContent += fmt.Sprint(i+1, ". ", item)
			} else {
				listContent += fmt.Sprint(i+1, ". ", item, "\n")
			}
		}

		replaceMap := docx.PlaceholderMap{
			"_nama_":        andalalin.NamaPemohon,
			"_jabatan_":     *andalalin.JabatanPemohon,
			"_alamat_":      andalalin.AlamatPemohon + ", Kelurahan " + andalalin.KelurahanPemohon + ", Kecamatan " + andalalin.KecamatanPemohon + ", Kabupaten " + andalalin.KabupatenPemohon + ", Provinsi " + andalalin.ProvinsiPemohon,
			"_pengembang_":  *andalalin.NamaPerusahaan,
			"_nomor_":       andalalin.Nomor,
			"_tanggal_":     andalalin.Tanggal[0:2],
			"_bulan_":       utils.Month(andalalin.Tanggal[3:5]),
			"_tahun_":       andalalin.Tanggal[6:10],
			"_kegiatan_":    andalalin.JenisProyek + " " + andalalin.NamaProyek + ", di Jalan " + andalalin.NamaJalan + ", Kelurahan " + andalalin.KelurahanProyek + ", Kecamatan " + andalalin.KecamatanProyek + ", Kabupaten " + andalalin.KabupatenProyek + ", Provinsi " + andalalin.ProvinsiProyek,
			"_rekomendasi_": andalalin.JenisProyek + " " + andalalin.NamaProyek + ", di Jalan " + andalalin.NamaJalan + ", Kelurahan " + andalalin.KelurahanProyek + ", Kecamatan " + andalalin.KecamatanProyek + ", Kabupaten " + andalalin.KabupatenProyek + ", Provinsi " + andalalin.ProvinsiProyek,
			"_kewajiban_":   listContent,
		}

		var docRendah string

		if andalalin.Pemohon == "Perorangan" {
			docRendah = "templates/suratPernyataanKesanggupanBangkitanSedangPerorangan.docx"
		} else {
			docRendah = "templates/suratPernyataanKesanggupanBangkitanSedangNonPerorangan.docx"
		}

		doc, err := docx.Open(docRendah)
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

		andalalin.StatusAndalalin = "Menunggu surat pernyataan"

		itemIndex := -1

		for i, item := range andalalin.BerkasPermohonan {
			if item.Nama == "Surat pernyataan kesanggupan (word)" {
				itemIndex = i
				break
			}
		}

		if itemIndex != -1 {
			andalalin.BerkasPermohonan[itemIndex].Berkas = docBytes
		} else {
			andalalin.BerkasPermohonan = append(andalalin.BerkasPermohonan, models.BerkasPermohonan{Status: "Selesai", Nama: "Surat pernyataan kesanggupan (word)", Tipe: "Word", Berkas: docBytes})
		}
	case "Bangkitan tinggi":
		listContent := ""
		for i, item := range payload.Kewajiban {
			if i == len(payload.Kewajiban)-1 {
				listContent += fmt.Sprint(i+1, ". ", item)
			} else {
				listContent += fmt.Sprint(i+1, ". ", item, "\n")
			}
		}

		replaceMap := docx.PlaceholderMap{
			"_nama_":       andalalin.NamaPemohon,
			"_jabatan_":    *andalalin.JabatanPemohon,
			"_alamat_":     andalalin.AlamatPemohon + ", Kelurahan " + andalalin.KelurahanPemohon + ", Kecamatan " + andalalin.KecamatanPemohon + ", Kabupaten " + andalalin.KabupatenPemohon + ", Provinsi " + andalalin.ProvinsiPemohon,
			"_pengembang_": *andalalin.NamaPerusahaan,
			"_nomor_":      andalalin.Nomor,
			"_tanggal_":    andalalin.Tanggal[0:2],
			"_bulan_":      utils.Month(andalalin.Tanggal[3:5]),
			"_tahun_":      andalalin.Tanggal[6:10],
			"_kegiatan_":   andalalin.JenisProyek + " " + andalalin.NamaProyek + ", di Jalan " + andalalin.NamaJalan + ", Kelurahan " + andalalin.KelurahanProyek + ", Kecamatan " + andalalin.KecamatanProyek + ", Kabupaten " + andalalin.KabupatenProyek + ", Provinsi " + andalalin.ProvinsiProyek,
			"_kewajiban_":  listContent,
		}

		var docRendah string

		if andalalin.Pemohon == "Perorangan" {
			docRendah = "templates/suratPernyataanKesanggupanBangkitanTinggiPerorangan.docx"
		} else {
			docRendah = "templates/suratPernyataanKesanggupanBangkitanTinggiNonPerorangan.docx"
		}

		doc, err := docx.Open(docRendah)
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

		andalalin.StatusAndalalin = "Menunggu surat pernyataan"

		itemIndex := -1

		for i, item := range andalalin.BerkasPermohonan {
			if item.Nama == "Surat pernyataan kesanggupan (word)" {
				itemIndex = i
				break
			}
		}

		if itemIndex != -1 {
			andalalin.BerkasPermohonan[itemIndex].Berkas = docBytes
		} else {
			andalalin.BerkasPermohonan = append(andalalin.BerkasPermohonan, models.BerkasPermohonan{Status: "Selesai", Nama: "Surat pernyataan kesanggupan (word)", Tipe: "Word", Berkas: docBytes})
		}
	}

	ac.DB.Save(&andalalin)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Surat berhasil dibuat"})
}

func (ac *AndalalinController) PembuatanSuratKeputusan(ctx *gin.Context) {
	var payload *models.Keputusan
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

	loc, _ := time.LoadLocation("Asia/Singapore")
	nowTime := time.Now().In(loc)

	switch andalalin.Bangkitan {
	case "Bangkitan rendah":
		t, err := template.ParseFiles("templates/suratKeputusanBangkitanRendah.html")
		if err != nil {
			log.Fatal("Error reading the email template:", err)
			return
		}

		var kegiatan string

		if *andalalin.NilaiKriteria == "" || andalalin.NilaiKriteria == nil {
			kegiatan = "Dengan luas lahan total sebesar ± " + andalalin.TotalLuasLahan + " <i>(terbilang meter persegi)</i>"
		} else {
			kegiatan = "Dengan luas lahan total sebesar ± " + andalalin.TotalLuasLahan + " <i>(terbilang meter persegi)</i> dan " + strings.ToLower(*andalalin.KriteriaKhusus) + " sebesar ± " + *andalalin.NilaiKriteria + " <i>(terbilang " + *andalalin.Terbilang + ")</i>"
		}

		keputusan := struct {
			NomorKeputusan     string
			JenisProyek        string
			JenisProyekJudul   string
			NamaProyek         string
			NamaProyekJudul    string
			Pengembang         string
			AlamatPengembang   string
			NomorPengembang    string
			NamaPimpinan       string
			JabatanPimpinan    string
			JalanJudul         string
			KelurahanJudul     string
			KabupatenJudul     string
			StatusJudul        string
			ProvinsiJudul      string
			NomorSurat         string
			TanggalSurat       string
			NomorKesanggupan   string
			TanggalKesanggupan string
			Jalan              string
			Kelurahan          string
			Kabupaten          string
			Status             string
			Provinsi           string
			Kegiatan           template.HTML
			NamaKadis          string
			NipKadis           string
			NomorLampiran      string
			TahunTerbit        string
			Data               []models.DataKeputusan
		}{
			NomorKeputusan:     payload.NomorKeputusan,
			JenisProyek:        andalalin.JenisProyek,
			JenisProyekJudul:   strings.ToUpper(andalalin.JenisProyek),
			NamaProyek:         andalalin.NamaProyek,
			NamaProyekJudul:    strings.ToUpper(andalalin.NamaProyek),
			Pengembang:         *andalalin.NamaPerusahaan,
			AlamatPengembang:   *andalalin.AlamatPerusahaan,
			NomorPengembang:    *andalalin.NomerPerusahaan,
			NamaPimpinan:       *andalalin.NamaPimpinan,
			JabatanPimpinan:    *andalalin.JabatanPimpinan,
			JalanJudul:         strings.ToUpper("JALAN " + andalalin.NamaJalan + " " + "DENGAN NOMOR RUAS JALAN " + andalalin.KodeJalan),
			KelurahanJudul:     strings.ToUpper("KELURAHAN " + andalalin.KelurahanProyek),
			KabupatenJudul:     strings.ToUpper("KABUPATEN " + andalalin.KabupatenProyek),
			StatusJudul:        strings.ToUpper(andalalin.StatusJalan),
			ProvinsiJudul:      strings.ToUpper("PROVINSI " + andalalin.ProvinsiProyek),
			NomorSurat:         andalalin.Nomor,
			TanggalSurat:       andalalin.Tanggal[0:2] + " " + utils.Month(andalalin.Tanggal[3:5]) + " " + andalalin.Tanggal[6:10],
			NomorKesanggupan:   payload.NomorKesanggupan,
			TanggalKesanggupan: payload.TanggalKesanggupan[0:2] + " " + utils.Month(payload.TanggalKesanggupan[3:5]) + " " + payload.TanggalKesanggupan[6:10],
			Jalan:              "Jalan " + andalalin.NamaJalan + " " + "dengan Nomor Ruas Jalan " + andalalin.KodeJalan,
			Kelurahan:          "Kelurahan " + andalalin.KelurahanProyek,
			Kabupaten:          "Kabupaten " + andalalin.KabupatenProyek,
			Status:             andalalin.StatusJalan,
			Provinsi:           "Provinsi " + andalalin.ProvinsiProyek,
			Kegiatan:           template.HTML(kegiatan),
			NamaKadis:          payload.NamaKadis,
			NipKadis:           payload.NipKadis,
			NomorLampiran:      payload.NomorLampiran,
			TahunTerbit:        nowTime.Format("2006"),
			Data:               payload.Data,
		}

		buffer := new(bytes.Buffer)
		if err = t.Execute(buffer, keputusan); err != nil {
			log.Fatal("Eror saat membaca template:", err)
			return
		}

		pdfContent, err := generatePDF(buffer.Bytes())
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		itemIndex := -1

		for i, item := range andalalin.BerkasPermohonan {
			if item.Nama == "Surat keputusan persetujuan teknis andalalin" {
				itemIndex = i
				break
			}
		}

		if itemIndex != -1 {
			andalalin.BerkasPermohonan[itemIndex].Berkas = pdfContent
			andalalin.BerkasPermohonan[itemIndex].Status = "Menunggu"
		} else {
			andalalin.BerkasPermohonan = append(andalalin.BerkasPermohonan, models.BerkasPermohonan{Status: "Menunggu", Nama: "Surat keputusan persetujuan teknis andalalin", Tipe: "Pdf", Berkas: pdfContent})
		}
	case "Bangkitan sedang":
	case "Bangkitan tinggi":
	}

	andalalin.StatusAndalalin = "Pemeriksaan surat keputusan"

	var user []models.User

	ac.DB.Find(&user, "role = ?", "Admin")

	for _, users := range user {
		simpanNotif := models.Notifikasi{
			IdUser: users.ID,
			Title:  "Surat keputusan baru",
			Body:   "Surat keputusan baru untuk permohonan dengan kode " + andalalin.Kode + " telah tersedia",
		}

		ac.DB.Create(&simpanNotif)

		if users.PushToken != "" {
			notif := utils.Notification{
				IdUser: users.ID,
				Title:  "Surat keputusan baru",
				Body:   "Surat keputusan baru untuk permohonan dengan kode " + andalalin.Kode + " telah tersedia",
				Token:  users.PushToken,
			}

			utils.SendPushNotifications(&notif)

		}
	}

	ac.DB.Save(&andalalin)

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (ac *AndalalinController) CheckKelengkapanAkhir(ctx *gin.Context) {
	var payload *models.KelengkapanAkhir
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

	switch andalalin.Bangkitan {
	case "Bangkitan rendah":
		t, err := template.ParseFiles("templates/checklistKelengkapanAkhirBangkitanRendah.html")
		if err != nil {
			log.Fatal("Error reading the email template:", err)
			return
		}

		kelengkapan := struct {
			Bangkitan   string
			Objek       string
			Lokasi      string
			Pengembang  string
			Diterima    string
			Pemeriksaan string
			Data        []models.DataKelengkapanAkhir
			Operator    string
			Nip         string
		}{
			Bangkitan:   "RENDAH",
			Objek:       andalalin.Jenis,
			Lokasi:      "Jalan " + andalalin.NamaJalan + ", Kelurahan " + andalalin.KelurahanProyek + ", Kecamatan " + andalalin.KecamatanProyek + ", Kabupaten " + andalalin.KabupatenProyek + ", Provinsi " + andalalin.ProvinsiProyek,
			Pengembang:  *andalalin.NamaPerusahaan,
			Diterima:    andalalin.TanggalAndalalin,
			Pemeriksaan: tanggal,
			Data:        payload.Kelengkapan,
			Operator:    currentUser.Name,
			Nip:         *currentUser.NIP,
		}

		buffer := new(bytes.Buffer)
		if err = t.Execute(buffer, kelengkapan); err != nil {
			log.Fatal("Eror saat membaca template:", err)
			return
		}

		pdfContent, err := generatePDF(buffer.Bytes())
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		itemIndex := -1

		for i, item := range andalalin.BerkasPermohonan {
			if item.Nama == "Checklist kelengkapan akhir" {
				itemIndex = i
				break
			}
		}

		if itemIndex != -1 {
			andalalin.BerkasPermohonan[itemIndex].Berkas = pdfContent
			andalalin.BerkasPermohonan[itemIndex].Status = "Menunggu"
		} else {
			andalalin.BerkasPermohonan = append(andalalin.BerkasPermohonan, models.BerkasPermohonan{Status: "Menunggu", Nama: "Checklist kelengkapan akhir", Tipe: "Pdf", Berkas: pdfContent})
		}

	case "Bangkitan sedang":
		t, err := template.ParseFiles("templates/checklistKelengkapanAkhirBangkitanSedang.html")
		if err != nil {
			log.Fatal("Error reading the email template:", err)
			return
		}

		kelengkapan := struct {
			Bangkitan   string
			Objek       string
			Lokasi      string
			Pengembang  string
			Pemohon     string
			Sertifikat  string
			Klasifikasi string
			Diterima    string
			Pemeriksaan string
			Data        []models.DataKelengkapanAkhir
			Operator    string
			Nip         string
		}{
			Bangkitan:   "SEDANG",
			Objek:       andalalin.Jenis,
			Lokasi:      "Jalan " + andalalin.NamaJalan + ", Kelurahan " + andalalin.KelurahanProyek + ", Kecamatan " + andalalin.KecamatanProyek + ", Kabupaten " + andalalin.KabupatenProyek + ", Provinsi " + andalalin.ProvinsiProyek,
			Pengembang:  *andalalin.NamaPerusahaan,
			Pemohon:     *andalalin.NamaPenyusunDokumen,
			Sertifikat:  *andalalin.NomerSertifikatPenyusunDokumen,
			Klasifikasi: *andalalin.KlasifikasiPenyusunDokumen,
			Diterima:    andalalin.TanggalAndalalin,
			Pemeriksaan: tanggal,
			Data:        payload.Kelengkapan,
			Operator:    currentUser.Name,
			Nip:         *currentUser.NIP,
		}

		buffer := new(bytes.Buffer)
		if err = t.Execute(buffer, kelengkapan); err != nil {
			log.Fatal("Eror saat membaca template:", err)
			return
		}

		pdfContent, err := generatePDF(buffer.Bytes())
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		itemIndex := -1

		for i, item := range andalalin.BerkasPermohonan {
			if item.Nama == "Checklist kelengkapan akhir" {
				itemIndex = i
				break
			}
		}

		if itemIndex != -1 {
			andalalin.BerkasPermohonan[itemIndex].Berkas = pdfContent
			andalalin.BerkasPermohonan[itemIndex].Status = "Menunggu"
		} else {
			andalalin.BerkasPermohonan = append(andalalin.BerkasPermohonan, models.BerkasPermohonan{Status: "Menunggu", Nama: "Checklist kelengkapan akhir", Tipe: "Pdf", Berkas: pdfContent})
		}

	case "Bangkitan tinggi":
		t, err := template.ParseFiles("templates/checklistKelengkapanAkhirBangkitanTinggi.html")
		if err != nil {
			log.Fatal("Error reading the email template:", err)
			return
		}

		kelengkapan := struct {
			Bangkitan   string
			Objek       string
			Lokasi      string
			Pengembang  string
			Pemohon     string
			Sertifikat  string
			Klasifikasi string
			Diterima    string
			Pemeriksaan string
			Data        []models.DataKelengkapanAkhir
			Operator    string
			Nip         string
		}{
			Bangkitan:   "TINGGI",
			Objek:       andalalin.Jenis,
			Lokasi:      "Jalan " + andalalin.NamaJalan + ", Kelurahan " + andalalin.KelurahanProyek + ", Kecamatan " + andalalin.KecamatanProyek + ", Kabupaten " + andalalin.KabupatenProyek + ", Provinsi " + andalalin.ProvinsiProyek,
			Pengembang:  *andalalin.NamaPerusahaan,
			Pemohon:     *andalalin.NamaPenyusunDokumen,
			Sertifikat:  *andalalin.NomerSertifikatPenyusunDokumen,
			Klasifikasi: *andalalin.KlasifikasiPenyusunDokumen,
			Diterima:    andalalin.TanggalAndalalin,
			Pemeriksaan: tanggal,
			Data:        payload.Kelengkapan,
			Operator:    currentUser.Name,
			Nip:         *currentUser.NIP,
		}

		buffer := new(bytes.Buffer)
		if err = t.Execute(buffer, kelengkapan); err != nil {
			log.Fatal("Eror saat membaca template:", err)
			return
		}

		pdfContent, err := generatePDF(buffer.Bytes())
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		itemIndex := -1

		for i, item := range andalalin.BerkasPermohonan {
			if item.Nama == "Checklist kelengkapan akhir" {
				itemIndex = i
				break
			}
		}

		if itemIndex != -1 {
			andalalin.BerkasPermohonan[itemIndex].Berkas = pdfContent
			andalalin.BerkasPermohonan[itemIndex].Status = "Menunggu"
		} else {
			andalalin.BerkasPermohonan = append(andalalin.BerkasPermohonan, models.BerkasPermohonan{Status: "Menunggu", Nama: "Checklist kelengkapan akhir", Tipe: "Pdf", Berkas: pdfContent})
		}
	}

	if andalalin.KelengkapanTidakSesuai != nil {
		andalalin.KelengkapanTidakSesuai = nil
	}

	for _, data := range payload.Kelengkapan {
		if data.Tidak != "" {
			for _, kelengkapan := range data.Dokumen {
				andalalin.KelengkapanTidakSesuai = append(andalalin.KelengkapanTidakSesuai, models.KelengkapanTidakSesuai{Dokumen: kelengkapan.Dokumen, Tipe: kelengkapan.Tipe, Role: data.Role})
			}
		}
	}

	andalalin.StatusAndalalin = "Persetujuan kelengkapan akhir"

	ac.DB.Save(&andalalin)

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (ac *AndalalinController) PembuatanPenyusunDokumen(ctx *gin.Context) {
	var payload *models.PenyusunDokumen
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

	switch andalalin.Bangkitan {
	case "Bangkitan sedang":
		t, err := template.ParseFiles("templates/penyusunanDokumenAndalalin.html")
		if err != nil {
			log.Fatal("Error reading the email template:", err)
			return
		}

		penyusun := struct {
			Bangkitan   string
			Objek       string
			Lokasi      string
			Pengembang  string
			Pemohon     string
			Sertifikat  string
			Klasifikasi string
			Diterima    string
			Pemeriksaan string
			Data        []models.DataPenyusunDokumen
			Operator    string
			Nip         string
		}{
			Bangkitan:   "SEDANG",
			Objek:       andalalin.Jenis,
			Lokasi:      "Jalan " + andalalin.NamaJalan + ", Kelurahan " + andalalin.KelurahanProyek + ", Kecamatan " + andalalin.KecamatanProyek + ", Kabupaten " + andalalin.KabupatenProyek + ", Provinsi " + andalalin.ProvinsiProyek,
			Pengembang:  *andalalin.NamaPerusahaan,
			Pemohon:     *andalalin.NamaPenyusunDokumen,
			Sertifikat:  *andalalin.NomerSertifikatPenyusunDokumen,
			Klasifikasi: *andalalin.KlasifikasiPenyusunDokumen,
			Diterima:    andalalin.TanggalAndalalin,
			Pemeriksaan: tanggal,
			Data:        payload.Penyusun,
			Operator:    currentUser.Name,
			Nip:         *currentUser.NIP,
		}

		buffer := new(bytes.Buffer)
		if err = t.Execute(buffer, penyusun); err != nil {
			log.Fatal("Eror saat membaca template:", err)
			return
		}

		pdfContent, err := generatePDF(buffer.Bytes())
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		itemIndex := -1

		for i, item := range andalalin.BerkasPermohonan {
			if item.Nama == "Penyusun dokumen analsis dampak lalu lintas" {
				itemIndex = i
				break
			}
		}

		if itemIndex != -1 {
			andalalin.BerkasPermohonan[itemIndex].Berkas = pdfContent
			andalalin.BerkasPermohonan[itemIndex].Status = "Menunggu"
		} else {
			andalalin.BerkasPermohonan = append(andalalin.BerkasPermohonan, models.BerkasPermohonan{Status: "Menunggu", Nama: "Penyusun dokumen analsis dampak lalu lintas", Tipe: "Pdf", Berkas: pdfContent})
		}
	case "Bangkitan tinggi":
		t, err := template.ParseFiles("templates/penyusunanDokumenAndalalin.html")
		if err != nil {
			log.Fatal("Error reading the email template:", err)
			return
		}

		penyusun := struct {
			Bangkitan   string
			Objek       string
			Lokasi      string
			Pengembang  string
			Pemohon     string
			Sertifikat  string
			Klasifikasi string
			Diterima    string
			Pemeriksaan string
			Data        []models.DataPenyusunDokumen
			Operator    string
			Nip         string
		}{
			Bangkitan:   "TINGGI",
			Objek:       andalalin.Jenis,
			Lokasi:      "Jalan " + andalalin.NamaJalan + ", Kelurahan " + andalalin.KelurahanProyek + ", Kecamatan " + andalalin.KecamatanProyek + ", Kabupaten " + andalalin.KabupatenProyek + ", Provinsi " + andalalin.ProvinsiProyek,
			Pengembang:  *andalalin.NamaPerusahaan,
			Pemohon:     *andalalin.NamaPenyusunDokumen,
			Sertifikat:  *andalalin.NomerSertifikatPenyusunDokumen,
			Klasifikasi: *andalalin.KlasifikasiPenyusunDokumen,
			Diterima:    andalalin.TanggalAndalalin,
			Pemeriksaan: tanggal,
			Data:        payload.Penyusun,
			Operator:    currentUser.Name,
			Nip:         *currentUser.NIP,
		}

		buffer := new(bytes.Buffer)
		if err = t.Execute(buffer, penyusun); err != nil {
			log.Fatal("Eror saat membaca template:", err)
			return
		}

		pdfContent, err := generatePDF(buffer.Bytes())
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		itemIndex := -1

		for i, item := range andalalin.BerkasPermohonan {
			if item.Nama == "Penyusun dokumen analsis dampak lalu lintas" {
				itemIndex = i
				break
			}
		}

		if itemIndex != -1 {
			andalalin.BerkasPermohonan[itemIndex].Berkas = pdfContent
			andalalin.BerkasPermohonan[itemIndex].Status = "Menunggu"
		} else {
			andalalin.BerkasPermohonan = append(andalalin.BerkasPermohonan, models.BerkasPermohonan{Status: "Menunggu", Nama: "Penyusun dokumen analsis dampak lalu lintas", Tipe: "Pdf", Berkas: pdfContent})
		}
	}

	andalalin.StatusAndalalin = "Persetujuan penyusun dokumen"

	ac.DB.Save(&andalalin)

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (ac *AndalalinController) PemeriksaanDokumenAndalalin(ctx *gin.Context) {
	var payload *models.PemeriksaanDokumen
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

	switch andalalin.Bangkitan {
	case "Bangkitan sedang":
		t, err := template.ParseFiles("templates/catatanAsistensiDokumen.html")
		if err != nil {
			log.Fatal("Error reading the email template:", err)
			return
		}

		pemeriksaan := struct {
			Bangkitan   string
			Objek       string
			Lokasi      string
			Pengembang  string
			Pemohon     string
			Sertifikat  string
			Klasifikasi string
			Diterima    string
			Pemeriksaan string
			Data        []models.CatatanAsistensi
			Operator    string
			Nip         string
		}{
			Bangkitan:   "SEDANG",
			Objek:       andalalin.Jenis,
			Lokasi:      "Jalan " + andalalin.NamaJalan + ", Kelurahan " + andalalin.KelurahanProyek + ", Kecamatan " + andalalin.KecamatanProyek + ", Kabupaten " + andalalin.KabupatenProyek + ", Provinsi " + andalalin.ProvinsiProyek,
			Pengembang:  *andalalin.NamaPerusahaan,
			Pemohon:     *andalalin.NamaPenyusunDokumen,
			Sertifikat:  *andalalin.NomerSertifikatPenyusunDokumen,
			Klasifikasi: *andalalin.KlasifikasiPenyusunDokumen,
			Diterima:    andalalin.TanggalAndalalin,
			Pemeriksaan: tanggal,
			Data:        payload.Pemeriksaan,
			Operator:    currentUser.Name,
			Nip:         *currentUser.NIP,
		}

		buffer := new(bytes.Buffer)
		if err = t.Execute(buffer, pemeriksaan); err != nil {
			log.Fatal("Eror saat membaca template:", err)
			return
		}

		pdfContent, err := generatePDF(buffer.Bytes())
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		itemIndex := -1

		for i, item := range andalalin.BerkasPermohonan {
			if item.Nama == "Catatan asistensi dokumen analisis dampak lalu lintas" {
				itemIndex = i
				break
			}
		}

		if itemIndex != -1 {
			andalalin.BerkasPermohonan[itemIndex].Berkas = pdfContent
			andalalin.BerkasPermohonan[itemIndex].Status = "Menunggu"
		} else {
			andalalin.BerkasPermohonan = append(andalalin.BerkasPermohonan, models.BerkasPermohonan{Status: "Menunggu", Nama: "Catatan asistensi dokumen analisis dampak lalu lintas", Tipe: "Pdf", Berkas: pdfContent})
		}
	case "Bangkitan tinggi":
		t, err := template.ParseFiles("templates/catatanAsistensiDokumen.html")
		if err != nil {
			log.Fatal("Error reading the email template:", err)
			return
		}

		pemeriksaan := struct {
			Bangkitan   string
			Objek       string
			Lokasi      string
			Pengembang  string
			Pemohon     string
			Sertifikat  string
			Klasifikasi string
			Diterima    string
			Pemeriksaan string
			Data        []models.CatatanAsistensi
			Operator    string
			Nip         string
		}{
			Bangkitan:   "SEDANG",
			Objek:       andalalin.Jenis,
			Lokasi:      "Jalan " + andalalin.NamaJalan + ", Kelurahan " + andalalin.KelurahanProyek + ", Kecamatan " + andalalin.KecamatanProyek + ", Kabupaten " + andalalin.KabupatenProyek + ", Provinsi " + andalalin.ProvinsiProyek,
			Pengembang:  *andalalin.NamaPerusahaan,
			Pemohon:     *andalalin.NamaPenyusunDokumen,
			Sertifikat:  *andalalin.NomerSertifikatPenyusunDokumen,
			Klasifikasi: *andalalin.KlasifikasiPenyusunDokumen,
			Diterima:    andalalin.TanggalAndalalin,
			Pemeriksaan: tanggal,
			Data:        payload.Pemeriksaan,
			Operator:    currentUser.Name,
			Nip:         *currentUser.NIP,
		}

		buffer := new(bytes.Buffer)
		if err = t.Execute(buffer, pemeriksaan); err != nil {
			log.Fatal("Eror saat membaca template:", err)
			return
		}

		pdfContent, err := generatePDF(buffer.Bytes())
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		itemIndex := -1

		for i, item := range andalalin.BerkasPermohonan {
			if item.Nama == "Catatan asistensi dokumen analisis dampak lalu lintas" {
				itemIndex = i
				break
			}
		}

		if itemIndex != -1 {
			andalalin.BerkasPermohonan[itemIndex].Berkas = pdfContent
			andalalin.BerkasPermohonan[itemIndex].Status = "Menunggu"
		} else {
			andalalin.BerkasPermohonan = append(andalalin.BerkasPermohonan, models.BerkasPermohonan{Status: "Menunggu", Nama: "Catatan asistensi dokumen analisis dampak lalu lintas", Tipe: "Pdf", Berkas: pdfContent})
		}
	}

	andalalin.HasilAsistensiDokumen = payload.Status
	andalalin.CatatanAsistensiDokumen = payload.Pemeriksaan
	andalalin.StatusAndalalin = "Persetujuan asistensi dokumen"

	ac.DB.Save(&andalalin)

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (ac *AndalalinController) PemeriksaanSuratKeputusan(ctx *gin.Context) {
	var payload *models.Pemeriksaan
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
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": result.Error})
		return
	}

	andalalin.HasilPemeriksaan = payload.Hasil
	andalalin.CatatanPemeriksaan = payload.Catatan
	if payload.Hasil == "Surat keputusan terpenuhi" {
		andalalin.StatusAndalalin = "Persetujuan surat keputusan"
	} else {
		andalalin.StatusAndalalin = "Pembuatan surat keputusan"
	}

	ac.DB.Save(&andalalin)

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

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
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Tidak ditemukan"})
		return
	}

	if perlalin.IdAndalalin != uuid.Nil {
		perlalin.IdPetugas = payload.IdPetugas
		perlalin.NamaPetugas = payload.NamaPetugas
		perlalin.EmailPetugas = payload.EmailPetugas
		perlalin.StatusAndalalin = "Survei lapangan"

		for i := range perlalin.Perlengkapan {
			perlalin.Perlengkapan[i].StatusPerlengkapan = "Survei"
		}

		ac.DB.Save(&perlalin)

		ac.ReleaseTicketLevel2(ctx, perlalin.IdAndalalin, payload.IdPetugas)
	}

	var user models.User
	resultUser := ac.DB.First(&user, "id = ?", perlalin.IdPetugas)
	if resultUser.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "User tidak ditemukan"})
		return
	}

	simpanNotif := models.Notifikasi{
		IdUser: user.ID,
		Title:  "Tugas baru",
		Body:   "Survei lapangan untuk permohonan dengan kode " + perlalin.Kode + " telah tersedia",
	}

	ac.DB.Create(&simpanNotif)

	if user.PushToken != "" {
		notif := utils.Notification{
			IdUser: user.ID,
			Title:  "Tugas baru",
			Body:   "Survei lapangan untuk permohonan dengan kode " + perlalin.Kode + " telah tersedia",
			Token:  user.PushToken,
		}

		utils.SendPushNotifications(&notif)

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
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Tidak ditemukan"})
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

	var user models.User
	resultUser := ac.DB.First(&user, "id = ?", perlalin.IdPetugas)
	if resultUser.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "User tidak ditemukan"})
		return
	}

	if perlalin.StatusAndalalin == "Survei lapangan" {
		simpanNotif := models.Notifikasi{
			IdUser: user.ID,
			Title:  "Tugas baru",
			Body:   "Survei lapangan untuk permohonan dengan kode " + perlalin.Kode + " telah tersedia",
		}

		ac.DB.Create(&simpanNotif)

		if user.PushToken != "" {
			notif := utils.Notification{
				IdUser: user.ID,
				Title:  "Tugas baru",
				Body:   "Survei lapangan untuk permohonan dengan kode " + perlalin.Kode + " telah tersedia",
				Token:  user.PushToken,
			}

			utils.SendPushNotifications(&notif)

		}
	} else {
		simpanNotif := models.Notifikasi{
			IdUser: user.ID,
			Title:  "Tugas baru",
			Body:   "Pemasangan perlengkapan untuk permohonan dengan kode " + perlalin.Kode + " telah tersedia",
		}

		ac.DB.Create(&simpanNotif)

		if user.PushToken != "" {
			notif := utils.Notification{
				IdUser: user.ID,
				Title:  "Tugas baru",
				Body:   "Pemasangan perlengkapan untuk permohonan dengan kode " + perlalin.Kode + " telah tersedia",
				Token:  user.PushToken,
			}

			utils.SendPushNotifications(&notif)

		}
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
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": results.Error})
		return
	} else {
		var respone []models.DaftarAndalalinResponse
		for _, s := range ticket {
			var perlalin models.Perlalin

			ac.DB.First(&perlalin, "id_andalalin = ? AND id_petugas = ?", s.IdAndalalin, currentUser.ID)

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

func (ac *AndalalinController) GetPerlengkapan(ctx *gin.Context) {
	id_andalalin := ctx.Param("id_andalalin")
	id_perlengkapan := ctx.Param("id_perlengkapan")

	var perlalin models.Perlalin

	ac.DB.First(&perlalin, "id_andalalin = ?", id_andalalin)

	if perlalin.IdAndalalin != uuid.Nil {
		for _, data := range perlalin.Perlengkapan {
			if data.IdPerlengkapan == id_perlengkapan {
				perlengkapan := struct {
					IdPerlengkapan       string        `json:"id_perlengkapan,omitempty"`
					StatusPerlengkapan   string        `json:"status,omitempty"`
					KategoriUtama        string        `json:"kategori_utama,omitempty"`
					KategoriPerlengkapan string        `json:"kategori,omitempty"`
					JenisPerlengkapan    string        `json:"perlengkapan,omitempty"`
					GambarPerlengkapan   string        `json:"gambar,omitempty"`
					LokasiPemasangan     string        `json:"pemasangan,omitempty"`
					LatitudePemasangan   float64       `json:"latitude,omitempty"`
					LongitudePemasangan  float64       `json:"longitude,omitempty"`
					FotoLokasi           []models.Foto `json:"foto,omitempty"`
					Detail               *string       `json:"detail,omitempty"`
					Alasan               string        `json:"alasan,omitempty"`
					Pertimbangan         *string       `json:"pertimbangan,omitempty"`
				}{
					IdPerlengkapan:       data.IdPerlengkapan,
					StatusPerlengkapan:   data.StatusPerlengkapan,
					KategoriUtama:        data.KategoriUtama,
					KategoriPerlengkapan: data.KategoriPerlengkapan,
					JenisPerlengkapan:    data.JenisPerlengkapan,
					GambarPerlengkapan:   data.GambarPerlengkapan,
					LokasiPemasangan:     data.LokasiPemasangan,
					LatitudePemasangan:   data.LatitudePemasangan,
					LongitudePemasangan:  data.LongitudePemasangan,
					FotoLokasi:           data.FotoLokasi,
					Detail:               data.Detail,
					Alasan:               data.Alasan,
					Pertimbangan:         data.Pertimbangan,
				}

				ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": perlengkapan})
			}
		}
	}
}

func (ac *AndalalinController) IsiSurvey(ctx *gin.Context) {
	var payload *models.DataSurvey
	currentUser := ctx.MustGet("currentUser").(models.User)
	id := ctx.Param("id_andalalin")
	id_perlengkapan := ctx.Param("id_perlengkapan")

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
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Tiket tidak ditemukan"})
		return
	}

	var perlalin models.Perlalin
	resultsPerlalin := ac.DB.First(&perlalin, "id_andalalin = ?", id)

	if resultsPerlalin.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Tidak ditemukan"})
		return
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	foto := []models.Foto{}

	for _, files := range form.File {
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
			foto = append(foto, data)
		}
	}

	if perlalin.IdAndalalin != uuid.Nil {
		survey := models.Survei{
			IdAndalalin:    perlalin.IdAndalalin,
			IdTiketLevel1:  ticket1.IdTiketLevel1,
			IdTiketLevel2:  ticket2.IdTiketLevel2,
			IdPerlengkapan: id_perlengkapan,
			IdPetugas:      currentUser.ID,
			Petugas:        currentUser.Name,
			EmailPetugas:   currentUser.Email,
			Lokasi:         payload.Data.Lokasi,
			Catatan:        payload.Data.Catatan,
			Foto:           foto,
			Latitude:       payload.Data.Latitude,
			Longitude:      payload.Data.Longitude,
			TanggalSurvei:  tanggal,
			WaktuSurvei:    nowTime.Format("15:04:05"),
		}

		result := ac.DB.Create(&survey)

		if result.Error != nil {
			fmt.Println(result.Error)

			if strings.Contains(strings.ToLower(result.Error.Error()), "unique constraint") {
				ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "Data is exist"})
				return
			} else {
				ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "An error occurred on the server. Please try again later"})
				return
			}
		}

		for i, data := range perlalin.Perlengkapan {
			if data.IdPerlengkapan == id_perlengkapan {
				perlalin.Perlengkapan[i].StatusPerlengkapan = "Pengecekan"
			}
		}

		var cek []string

		for _, data := range perlalin.Perlengkapan {
			if data.StatusPerlengkapan == "Survei" {
				cek = append(cek, "Ada")
			}
		}

		if cek == nil {
			perlalin.StatusAndalalin = "Pengecekan perlengkapan"
			ac.CloseTiketLevel2(ctx, perlalin.IdAndalalin)

			var survey []models.Survei

			result := ac.DB.Find(&survey, "id_andalalin = ?", id)
			if result.Error != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": result.Error})
				return
			}

			laporan := []models.DataLaporanSurvei{}

			for _, data := range survey {
				for _, perlengkapan := range perlalin.Perlengkapan {
					if perlengkapan.IdPerlengkapan == data.IdPerlengkapan {
						laporan = append(laporan, models.DataLaporanSurvei{Perlengkapan: perlengkapan.JenisPerlengkapan, Lokasi: perlengkapan.LokasiPemasangan, Tanggal: data.TanggalSurvei, Survei: data.Lokasi, Catatan: data.Catatan, Foto: data.Foto})
					}
				}
			}

			perlengkapan := struct {
				Kode  string
				Nik   string
				Nama  string
				Email string
				Nomor string
				Data  []models.DataLaporanSurvei
			}{
				Kode:  perlalin.Kode,
				Nik:   perlalin.NikPemohon,
				Nama:  perlalin.NamaPemohon,
				Email: perlalin.EmailPemohon,
				Nomor: perlalin.NomerPemohon,
				Data:  laporan,
			}

			t, err := template.ParseFiles("templates/laporanSurvei.html")
			if err != nil {
				log.Fatal("Error reading the email template:", err)
				return
			}

			buffer := new(bytes.Buffer)
			if err = t.Execute(buffer, perlengkapan); err != nil {
				log.Fatal("Eror saat membaca template:", err)
				return
			}

			pdfContent, err := generatePDF(buffer.Bytes())
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			perlalin.BerkasPermohonan = append(perlalin.BerkasPermohonan, models.BerkasPermohonan{Status: "Selesai", Nama: "Laporan survei", Tipe: "Pdf", Berkas: pdfContent})

		}

		ac.DB.Save(&perlalin)
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success"})
}

func (ac *AndalalinController) GetSurvey(ctx *gin.Context) {
	id := ctx.Param("id_andalalin")
	id_perlengkapan := ctx.Param("id_perlengkapan")

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

	result := ac.DB.First(&survey, "id_andalalin = ? AND id_perlengkapan = ?", id, id_perlengkapan)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": result.Error})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": survey})
}

func (ac *AndalalinController) Pengecekanperlengkapan(ctx *gin.Context) {
	id := ctx.Param("id_andalalin")

	var payload *models.Pengecekan

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

	var perlalin models.Perlalin
	resultsPerlalin := ac.DB.First(&perlalin, "id_andalalin = ?", id)

	if resultsPerlalin.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Tidak ditemukan"})
		return
	}

	if perlalin.IdAndalalin != uuid.Nil {
		for _, data := range payload.Data {
			if data.Tidak == "" {
				for j, perlengkapan := range perlalin.Perlengkapan {
					if perlengkapan.IdPerlengkapan == data.ID {
						perlalin.Perlengkapan[j].StatusPerlengkapan = "Pemasangan"
						perlalin.Perlengkapan[j].Pertimbangan = data.Pertimbangan
					}
				}
			} else {
				for j, perlengkapan := range perlalin.Perlengkapan {
					if perlengkapan.IdPerlengkapan == data.ID {
						perlalin.Perlengkapan[j].StatusPerlengkapan = "Tidak disetujui"
						perlalin.Perlengkapan[j].Pertimbangan = data.Pertimbangan
					}
				}
			}
		}

		var cek []string

		for _, data := range perlalin.Perlengkapan {
			if data.StatusPerlengkapan == "Pemasangan" {
				cek = append(cek, "Ada")
			}
		}

		if cek != nil {
			perlalin.StatusAndalalin = "Pemasangan perlengkapan"
		} else {
			ac.BatalkanPermohonanForNothing(ctx, id)
		}

		ac.DB.Save(&perlalin)
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success"})
}

func (ac *AndalalinController) BatalkanPermohonanForNothing(ctx *gin.Context, id string) {
	var permohonan models.Perlalin

	result := ac.DB.First(&permohonan, "id_andalalin = ?", id)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": result.Error})
		return
	}

	permohonan.PertimbanganPembatalan = "Tidak ada perlengkapan yang disetujui"
	permohonan.StatusAndalalin = "Permohonan dibatalkan"
	ac.DB.Save(&permohonan)

	var user models.User
	resultUser := ac.DB.First(&user, "id = ?", permohonan.IdUser)
	if resultUser.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "User tidak ditemukan"})
		return
	}

	simpanNotif := models.Notifikasi{
		IdUser: user.ID,
		Title:  "Permohonan dibatalkan",
		Body:   "Permohonan anda dengan kode " + permohonan.Kode + " telah dibatalkan",
	}

	ac.DB.Create(&simpanNotif)

	if user.PushToken != "" {
		notif := utils.Notification{
			IdUser: user.ID,
			Title:  "Permohonan dibatalkan",
			Body:   "Permohonan anda dengan kode " + permohonan.Kode + " telah dibatalkan",
			Token:  user.PushToken,
		}

		utils.SendPushNotifications(&notif)

	}
}

func (ac *AndalalinController) BatalkanPermohonan(ctx *gin.Context) {
	id := ctx.Param("id_andalalin")
	var payload *models.Pertimbangan
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	var permohonan models.Perlalin

	result := ac.DB.First(&permohonan, "id_andalalin = ?", id)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": result.Error})
		return
	}

	ac.CloseTiketLevel1(ctx, permohonan.IdAndalalin)
	permohonan.PertimbanganPembatalan = payload.Pertimbangan
	permohonan.StatusAndalalin = "Permohonan dibatalkan"
	ac.DB.Save(&permohonan)

	var user models.User
	resultUser := ac.DB.First(&user, "id = ?", permohonan.IdUser)
	if resultUser.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "User tidak ditemukan"})
		return
	}

	simpanNotif := models.Notifikasi{
		IdUser: user.ID,
		Title:  "Permohonan dibatalkan",
		Body:   "Permohonan anda dengan kode " + permohonan.Kode + " telah dibatalkan",
	}

	ac.DB.Create(&simpanNotif)

	if user.PushToken != "" {
		notif := utils.Notification{
			IdUser: user.ID,
			Title:  "Permohonan dibatalkan",
			Body:   "Permohonan anda dengan kode " + permohonan.Kode + " telah dibatalkan",
			Token:  user.PushToken,
		}

		utils.SendPushNotifications(&notif)

	}
}

func (ac *AndalalinController) TundaPemasangan(ctx *gin.Context) {
	id := ctx.Param("id_andalalin")

	var payload *models.Pertimbangan
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
		perlalin.StatusAndalalin = "Pemasangan ditunda"
		perlalin.PertimbanganPenundaan = payload.Pertimbangan
		ac.DB.Save(&perlalin)

		var user models.User
		resultUser := ac.DB.First(&user, "id = ?", perlalin.IdUser)
		if resultUser.Error != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "User tidak ditemukan"})
			return
		}

		simpanNotif := models.Notifikasi{
			IdUser: user.ID,
			Title:  "Pemasangan ditunda",
			Body:   "Pemasangan perlengkapan pada permohonan anda dengan kode " + perlalin.Kode + " telah ditunda",
		}

		ac.DB.Create(&simpanNotif)

		if user.PushToken != "" {
			notif := utils.Notification{
				IdUser: user.ID,
				Title:  "Pemasangan ditunda",
				Body:   "Pemasangan perlengkapan pada permohonan anda dengan kode " + perlalin.Kode + " telah ditunda",
				Token:  user.PushToken,
			}

			utils.SendPushNotifications(&notif)

		}
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (ac *AndalalinController) LanjutkanPemasangan(ctx *gin.Context) {
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
		perlalin.StatusAndalalin = "Pemasangan perlengkapan"
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
			Title:  "Pemasangan dilanjutkan",
			Body:   "Pemasangan perlengkapan pada permohonan anda dengan kode " + perlalin.Kode + " telah dilanjutkan",
		}

		ac.DB.Create(&simpanNotif)

		if user.PushToken != "" {
			notif := utils.Notification{
				IdUser: user.ID,
				Title:  "Pemasangan dilanjutkan",
				Body:   "Pemasangan perlengkapan pada permohonan anda dengan kode " + perlalin.Kode + " telah dilanjutkan",
				Token:  user.PushToken,
			}

			utils.SendPushNotifications(&notif)

		}
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
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
		if s.StatusAndalalin == "Pemasangan perlengkapan" && s.IdPetugas == currentUser.ID {
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
	id_perlengkapan := ctx.Param("id_perlengkapan")

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
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Tiket tidak ditemukan"})
		return
	}

	var perlalin models.Perlalin
	resultsPerlalin := ac.DB.First(&perlalin, "id_andalalin = ?", id)

	if resultsPerlalin.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Tidak ditemukan"})
		return
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	foto := []models.Foto{}

	for _, files := range form.File {
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
			foto = append(foto, data)
		}
	}

	if perlalin.IdAndalalin != uuid.Nil {
		survey := models.Pemasangan{
			IdAndalalin:       perlalin.IdAndalalin,
			IdTiketLevel1:     ticket1.IdTiketLevel1,
			IdPerlengkapan:    id_perlengkapan,
			IdPetugas:         currentUser.ID,
			Petugas:           currentUser.Name,
			EmailPetugas:      currentUser.Email,
			Lokasi:            payload.Data.Lokasi,
			Catatan:           payload.Data.Catatan,
			Foto:              foto,
			Latitude:          payload.Data.Latitude,
			Longitude:         payload.Data.Longitude,
			WaktuPemasangan:   nowTime.Format("15:04:05"),
			TanggalPemasangan: tanggal,
		}

		result := ac.DB.Create(&survey)

		if result.Error != nil {
			fmt.Println(result.Error)

			if strings.Contains(strings.ToLower(result.Error.Error()), "unique constraint") {
				ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "Data is exist"})
				return
			} else {
				ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "An error occurred on the server. Please try again later"})
				return
			}
		}

		for i, data := range perlalin.Perlengkapan {
			if data.IdPerlengkapan == id_perlengkapan {
				perlalin.Perlengkapan[i].StatusPerlengkapan = "Selesai"
			}
		}

		var cek []string

		for _, data := range perlalin.Perlengkapan {
			if data.StatusPerlengkapan == "Pemasangan" {
				cek = append(cek, "Ada")
			}
		}

		if cek == nil {
			perlalin.StatusAndalalin = "Pemasangan selesai"
			ac.PemasanganSelesai(ctx, perlalin)
			ac.CloseTiketLevel1(ctx, perlalin.IdAndalalin)

			var pemasangan []models.Pemasangan

			result := ac.DB.Find(&pemasangan, "id_andalalin = ?", id)
			if result.Error != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": result.Error})
				return
			}

			laporan := []models.DataLaporanSurvei{}

			for _, data := range pemasangan {
				for _, perlengkapan := range perlalin.Perlengkapan {
					if perlengkapan.IdPerlengkapan == data.IdPerlengkapan {
						laporan = append(laporan, models.DataLaporanSurvei{Perlengkapan: perlengkapan.JenisPerlengkapan, Lokasi: perlengkapan.LokasiPemasangan, Tanggal: data.TanggalPemasangan, Survei: data.Lokasi, Catatan: data.Catatan, Foto: data.Foto})
					}
				}
			}

			perlengkapan := struct {
				Kode  string
				Nik   string
				Nama  string
				Email string
				Nomor string
				Data  []models.DataLaporanSurvei
			}{
				Kode:  perlalin.Kode,
				Nik:   perlalin.NikPemohon,
				Nama:  perlalin.NamaPemohon,
				Email: perlalin.EmailPemohon,
				Nomor: perlalin.NomerPemohon,
				Data:  laporan,
			}

			t, err := template.ParseFiles("templates/laporanPemasangan.html")
			if err != nil {
				log.Fatal("Error reading the email template:", err)
				return
			}

			buffer := new(bytes.Buffer)
			if err = t.Execute(buffer, perlengkapan); err != nil {
				log.Fatal("Eror saat membaca template:", err)
				return
			}

			pdfContent, err := generatePDF(buffer.Bytes())
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			perlalin.BerkasPermohonan = append(perlalin.BerkasPermohonan, models.BerkasPermohonan{Status: "Selesai", Nama: "Laporan pemasangan", Tipe: "Pdf", Berkas: pdfContent})

		}

		ac.DB.Save(&perlalin)
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success"})
}

func (ac *AndalalinController) PemasanganSelesai(ctx *gin.Context, permohonan models.Perlalin) {
	var user models.User
	resultUser := ac.DB.First(&user, "id = ?", permohonan.IdUser)
	if resultUser.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "User tidak ditemukan"})
		return
	}

	simpanNotif := models.Notifikasi{
		IdUser: user.ID,
		Title:  "Pemasangan selesai",
		Body:   "Permohonan anda dengan kode " + permohonan.Kode + " telah selesai",
	}

	ac.DB.Create(&simpanNotif)

	if user.PushToken != "" {
		notif := utils.Notification{
			IdUser: user.ID,
			Title:  "Pemasangan selesai",
			Body:   "Permohonan anda dengan kode " + permohonan.Kode + " telah selesai",
			Token:  user.PushToken,
		}

		utils.SendPushNotifications(&notif)

	}
}

func (ac *AndalalinController) GetPemasangan(ctx *gin.Context) {
	id := ctx.Param("id_andalalin")

	var pemasangan *models.Pemasangan

	result := ac.DB.First(&pemasangan, "id_andalalin = ?", id)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": result.Error})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": pemasangan})
}
