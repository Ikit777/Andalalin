package controllers

import (
	"bytes"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Ikit777/E-Andalalin/initializers"
	"github.com/Ikit777/E-Andalalin/models"
	"github.com/Ikit777/E-Andalalin/repository"
	"github.com/Ikit777/E-Andalalin/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	wkhtmltopdf "github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

type AndalalinController struct {
	DB *gorm.DB
}

func NewAndalalinController(DB *gorm.DB) AndalalinController {
	return AndalalinController{DB}
}

const path = "C:/Program Files/wkhtmltopdf/bin/wkhtmltopdf" //Sesuaikan dengan path wkhtmltopdf

func (ac *AndalalinController) Pengajuan(ctx *gin.Context) {
	var payload *models.DataAndalalin
	currentUser := ctx.MustGet("currentUser").(models.User)

	config, _ := initializers.LoadConfig(".")

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
	now := time.Now().In(loc).Format("02-01-2006")
	nowTime := time.Now().In(loc)

	nomer := utils.Encode(6)

	t, err := template.ParseFiles("templates/tandaterimaTemplate.html")
	if err != nil {
		log.Fatal("Error reading the email template:", err)
		return
	}

	bukti := struct {
		Nomer  string
		Bulan  string
		Tahun  int
		Nama   string
		Alamat string
		Tlp    string
		Izin   string
		Waktu  string
	}{
		Nomer:  nomer,
		Bulan:  nowTime.Format("02"),
		Tahun:  nowTime.Year(),
		Nama:   payload.Andalalin.NamaPemohon,
		Alamat: payload.Andalalin.AlamatPemohon,
		Tlp:    payload.Andalalin.NomerPemohon,
		Izin:   payload.Andalalin.JenisAndalalin,
		Waktu:  now,
	}

	buffer := new(bytes.Buffer)
	if err = t.Execute(buffer, bukti); err != nil {
		log.Fatal("Eror saat membaca template:", err)
		return
	}

	wkhtmltopdf.SetPath(path)

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

	err = pdfg.Create()
	if err != nil {
		log.Fatal(err)
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

	permohonan := models.Andalalin{
		IdUser:                 currentUser.ID,
		JenisAndalalin:         payload.Andalalin.JenisAndalalin,
		NomerAndalalin:         nomer,
		NikPemohon:             payload.Andalalin.NikPemohon,
		NamaPemohon:            payload.Andalalin.NamaPemohon,
		EmailPemohon:           currentUser.Email,
		TempatLahirPemohon:     payload.Andalalin.TempatLahirPemohon,
		TanggalLahirPemohon:    payload.Andalalin.TanggalLahirPemohon,
		AlamatPemohon:          payload.Andalalin.AlamatPemohon,
		JenisKelaminPemohon:    payload.Andalalin.JenisKelaminPemohon,
		NomerPemohon:           payload.Andalalin.NomerPemohon,
		JabatanPemohon:         payload.Andalalin.JabatanPemohon,
		LokasiPengambilan:      payload.Andalalin.LokasiPengambilan,
		WaktuAndalalin:         now,
		StatusAndalalin:        "Cek persyaratan",
		TandaTerimaPendaftaran: pdfg.Bytes(),

		NamaPerusahaan:       payload.Andalalin.NamaPerusahaan,
		AlamatPerusahaan:     payload.Andalalin.AlamatPerusahaan,
		NomerPerusahaan:      payload.Andalalin.NomerPerusahaan,
		EmailPerusahaan:      payload.Andalalin.EmailPerusahaan,
		ProvinsiPerusahaan:   payload.Andalalin.ProvinsiPerusahaan,
		KabupatenPerusahaan:  payload.Andalalin.KabupatenPerusahaan,
		KecamatanPerusahaan:  payload.Andalalin.KecamatanPerusahaan,
		KelurahaanPerusahaan: payload.Andalalin.KelurahaanPerusahaan,
		NamaPimpinan:         payload.Andalalin.NamaPimpinan,
		JabatanPimpinan:      payload.Andalalin.JabatanPimpinan,
		JenisKelaminPimpinan: payload.Andalalin.JenisKelaminPimpinan,
		JenisKegiatan:        payload.Andalalin.JenisKegiatan,
		Peruntukan:           payload.Andalalin.Peruntukan,
		LuasLahan:            payload.Andalalin.LuasLahan + "m²",
		AlamatPersil:         payload.Andalalin.AlamatPersil,
		KelurahanPersil:      payload.Andalalin.KelurahanPersil,
		NomerSKRK:            payload.Andalalin.NomerSKRK,
		TanggalSKRK:          payload.Andalalin.TanggalSKRK,

		KartuTandaPenduduk: blobs["ktp"],
		AktaPendirianBadan: blobs["apb"],
		SuratKuasa:         blobs["sk"],
	}

	result := ac.DB.Create(&permohonan)

	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "eror saat mengirim data"})
		return
	} else {
		ac.ReleaseTicketLevel1(ctx, permohonan.IdAndalalin)
		ctx.JSON(http.StatusOK, gin.H{"status": "success"})
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

	result := ac.DB.Model(&tiket).Where("id_andalalin = ?", id).Update("status", "Tutup")
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Tiket level 1 tidak tersedia"})
		return
	}
}

func (ac *AndalalinController) ReleaseTicketLevel2(ctx *gin.Context, id uuid.UUID) {
	var tiket1 models.TiketLevel1
	results := ac.DB.First(&tiket1, "id_andalalin = ?", id)

	if results.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": results.Error})
		return
	}

	tiket := models.TiketLevel2{
		IdTiketLevel1: tiket1.IdTiketLevel1,
		IdAndalalin:   id,
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

	result := ac.DB.Model(&tiket).Where("id_andalalin = ?", id).Update("status", "Tutup")
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Telah terjadi sesuatu"})
		return
	}
}

func (ac *AndalalinController) GetPermohonanByIdUser(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)

	var andalalin []models.Andalalin

	results := ac.DB.Find(&andalalin, "id_user = ?", currentUser.ID)

	if results.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": results.Error})
		return
	} else {
		var respone []models.DaftarAndalalinResponse
		for _, s := range andalalin {
			respone = append(respone, models.DaftarAndalalinResponse{
				IdAndalalin:     s.IdAndalalin,
				NomerAndalalin:  s.NomerAndalalin,
				WaktuAndalalin:  s.WaktuAndalalin,
				Nama:            s.NamaPemohon,
				Alamat:          s.AlamatPemohon,
				JenisAndalalin:  s.JenisAndalalin,
				StatusAndalalin: s.StatusAndalalin,
			})
		}
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "results": len(respone), "data": respone})
	}
}

func (ac *AndalalinController) GetPermohonanByIdAndalalin(ctx *gin.Context) {
	id := ctx.Param("id_andalalin")

	var andalalin models.Andalalin

	results := ac.DB.First(&andalalin, "id_andalalin = ?", id)

	if results.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": results.Error})
		return
	}

	data := models.AndalalinResponse{
		IdAndalalin:                  andalalin.IdAndalalin,
		JenisAndalalin:               andalalin.JenisAndalalin,
		NomerAndalalin:               andalalin.NomerAndalalin,
		NikPemohon:                   andalalin.NikPemohon,
		NamaPemohon:                  andalalin.NamaPemohon,
		EmailPemohon:                 andalalin.EmailPemohon,
		TempatLahirPemohon:           andalalin.TempatLahirPemohon,
		TanggalLahirPemohon:          andalalin.TanggalLahirPemohon,
		AlamatPemohon:                andalalin.AlamatPemohon,
		JenisKelaminPemohon:          andalalin.JenisKelaminPemohon,
		NomerPemohon:                 andalalin.NomerPemohon,
		JabatanPemohon:               andalalin.JabatanPemohon,
		LokasiPengambilan:            andalalin.LokasiPengambilan,
		WaktuAndalalin:               andalalin.WaktuAndalalin,
		StatusAndalalin:              andalalin.StatusAndalalin,
		TandaTerimaPendaftaran:       andalalin.TandaTerimaPendaftaran,
		NamaPerusahaan:               andalalin.NamaPerusahaan,
		AlamatPerusahaan:             andalalin.AlamatPerusahaan,
		NomerPerusahaan:              andalalin.NomerPerusahaan,
		EmailPerusahaan:              andalalin.EmailPerusahaan,
		ProvinsiPerusahaan:           andalalin.ProvinsiPerusahaan,
		KabupatenPerusahaan:          andalalin.KabupatenPerusahaan,
		KecamatanPerusahaan:          andalalin.KecamatanPerusahaan,
		KelurahaanPerusahaan:         andalalin.KelurahaanPerusahaan,
		NamaPimpinan:                 andalalin.NamaPimpinan,
		JabatanPimpinan:              andalalin.JabatanPimpinan,
		JenisKelaminPimpinan:         andalalin.JenisKelaminPimpinan,
		JenisKegiatan:                andalalin.JenisKegiatan,
		Peruntukan:                   andalalin.Peruntukan,
		LuasLahan:                    andalalin.LuasLahan,
		AlamatPersil:                 andalalin.AlamatPersil,
		KelurahanPersil:              andalalin.KelurahanPersil,
		NomerSKRK:                    andalalin.NomerSKRK,
		TanggalSKRK:                  andalalin.TanggalSKRK,
		KartuTandaPenduduk:           andalalin.KartuTandaPenduduk,
		AktaPendirianBadan:           andalalin.AktaPendirianBadan,
		SuratKuasa:                   andalalin.SuratKuasa,
		IdPetugas:                    andalalin.IdPetugas,
		NamaPetugas:                  andalalin.NamaPetugas,
		EmailPetugas:                 andalalin.EmailPetugas,
		PersetujuanDokumen:           andalalin.PersetujuanDokumen,
		KeteranganPersetujuanDokumen: andalalin.KeteranganPersetujuanDokumen,
		NomerBAPDasar:                andalalin.NomerBAPDasar,
		NomerBAPPelaksanaan:          andalalin.NomerBAPPelaksanaan,
		TanggalBAP:                   andalalin.TanggalBAP,
		FileBAP:                      andalalin.FileBAP,
		FileSK:                       andalalin.FileSK,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": data})
}

func (ac *AndalalinController) GetPermohonanByStatus(ctx *gin.Context) {
	status := ctx.Param("status_andalalin")

	config, _ := initializers.LoadConfig(".")

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

	results := ac.DB.Find(&andalalin, "status_andalalin = ?", status)

	if results.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": results.Error})
		return
	} else {
		var respone []models.DaftarAndalalinResponse
		for _, s := range andalalin {
			respone = append(respone, models.DaftarAndalalinResponse{
				IdAndalalin:     s.IdAndalalin,
				NomerAndalalin:  s.NomerAndalalin,
				WaktuAndalalin:  s.WaktuAndalalin,
				Nama:            s.NamaPemohon,
				Alamat:          s.AlamatPemohon,
				JenisAndalalin:  s.JenisAndalalin,
				StatusAndalalin: s.StatusAndalalin,
			})
		}
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "results": len(respone), "data": respone})
	}
}

func (ac *AndalalinController) GetPermohonan(ctx *gin.Context) {
	config, _ := initializers.LoadConfig(".")

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

	results := ac.DB.Find(&andalalin)

	if results.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": results.Error})
		return
	} else {
		var respone []models.DaftarAndalalinResponse
		for _, s := range andalalin {
			respone = append(respone, models.DaftarAndalalinResponse{
				IdAndalalin:     s.IdAndalalin,
				NomerAndalalin:  s.NomerAndalalin,
				WaktuAndalalin:  s.WaktuAndalalin,
				Nama:            s.NamaPemohon,
				Alamat:          s.AlamatPemohon,
				JenisAndalalin:  s.JenisAndalalin,
				StatusAndalalin: s.StatusAndalalin,
			})
		}
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "results": len(respone), "data": respone})
	}
}

func (ac *AndalalinController) GetAndalalinTicketLevel1(ctx *gin.Context) {
	status := ctx.Param("status")

	config, _ := initializers.LoadConfig(".")

	accessUser := ctx.MustGet("accessUser").(string)

	claim, error := utils.ValidateToken(accessUser, config.AccessTokenPublicKey)
	if error != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": error.Error()})
		return
	}

	credential := claim.Credentials[repository.AndalalinTicket1Credential]

	if !credential {
		// Return status 403 and permission denied error message.
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": true,
			"msg":   "Permission denied",
		})
		return
	}

	var ticket []models.TiketLevel1

	results := ac.DB.Find(&ticket, "status = ?", status)

	if results.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": results.Error})
		return
	} else {
		var respone []models.DaftarAndalalinResponse
		for _, s := range ticket {
			var andalalin models.Andalalin
			results := ac.DB.First(&andalalin, "id_andalalin = ?", s.IdAndalalin)

			if results.Error != nil {
				ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": results.Error})
				return
			}

			respone = append(respone, models.DaftarAndalalinResponse{
				IdAndalalin:     andalalin.IdAndalalin,
				NomerAndalalin:  andalalin.NomerAndalalin,
				WaktuAndalalin:  andalalin.WaktuAndalalin,
				Nama:            andalalin.NamaPemohon,
				Alamat:          andalalin.AlamatPemohon,
				JenisAndalalin:  andalalin.JenisAndalalin,
				StatusAndalalin: andalalin.StatusAndalalin,
			})
		}
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "results": len(respone), "data": respone})
	}

}

func (ac *AndalalinController) GetPersyaratan(ctx *gin.Context) {
	id := ctx.Param("id_andalalin")

	config, _ := initializers.LoadConfig(".")

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

	var andalalin models.Andalalin

	results := ac.DB.Find(&andalalin, "id_andalalin = ?", id)

	if results.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": results.Error})
		return
	}

	persyaratan := &models.PersayaratanRespone{
		IdAndalalin:        andalalin.IdAndalalin,
		KartuTandaPenduduk: andalalin.KartuTandaPenduduk,
		AktaPendirianBadan: andalalin.AktaPendirianBadan,
		SuratKuasa:         andalalin.SuratKuasa,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": persyaratan})
}

func (ac *AndalalinController) UpdatePersyaratan(ctx *gin.Context) {
	config, _ := initializers.LoadConfig(".")

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

	id := ctx.Param("id_andalalin")
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

	var andalalin *models.Andalalin
	results := ac.DB.Find(&andalalin, "id_andalalin = ?", id)

	if results.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": results.Error})
		return
	}

	if blobs["ktp"] == nil {
		blobs["ktp"] = andalalin.KartuTandaPenduduk
	}

	if blobs["apb"] == nil {
		blobs["apb"] = andalalin.AktaPendirianBadan
	}

	if blobs["sk"] == nil {
		blobs["sk"] = andalalin.SuratKuasa
	}

	andalalin.KartuTandaPenduduk = blobs["ktp"]
	andalalin.AktaPendirianBadan = blobs["apb"]
	andalalin.SuratKuasa = blobs["sk"]
	andalalin.StatusAndalalin = "Cek persyaratan"

	ac.DB.Save(&andalalin)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "msg": "persyaratan berhasil diupdate"})
}

func (ac *AndalalinController) GetPerusahaan(ctx *gin.Context) {
	id := ctx.Param("id_andalalin")

	config, _ := initializers.LoadConfig(".")

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

	var andalalin models.Andalalin

	results := ac.DB.Find(&andalalin, "id_andalalin = ?", id)

	if results.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": results.Error})
		return
	}

	perusahaan := &models.PerusahaanRespone{
		NamaPerusahaan:       andalalin.NamaPerusahaan,
		AlamatPerusahaan:     andalalin.AlamatPerusahaan,
		NomerPerusahaan:      andalalin.NomerPerusahaan,
		EmailPerusahaan:      andalalin.EmailPerusahaan,
		ProvinsiPerusahaan:   andalalin.ProvinsiPerusahaan,
		KabupatenPerusahaan:  andalalin.KabupatenPerusahaan,
		KecamatanPerusahaan:  andalalin.KecamatanPerusahaan,
		KelurahaanPerusahaan: andalalin.KelurahaanPerusahaan,
		NamaPimpinan:         andalalin.NamaPimpinan,
		JabatanPimpinan:      andalalin.JabatanPimpinan,
		JenisKelaminPimpinan: andalalin.JenisKelaminPimpinan,
		JenisKegiatan:        andalalin.JenisKegiatan,
		Peruntukan:           andalalin.Peruntukan,
		LuasLahan:            andalalin.LuasLahan,
		AlamatPersil:         andalalin.AlamatPersil,
		KelurahanPersil:      andalalin.KelurahanPersil,
		NomerSKRK:            andalalin.NomerSKRK,
		TanggalSKRK:          andalalin.TanggalSKRK,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": perusahaan})
}

func (ac *AndalalinController) PersyaratanTerpenuhi(ctx *gin.Context) {
	id := ctx.Param("id_andalalin")

	config, _ := initializers.LoadConfig(".")

	accessUser := ctx.MustGet("accessUser").(string)

	claim, error := utils.ValidateToken(accessUser, config.AccessTokenPublicKey)
	if error != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": error.Error()})
		return
	}

	credential := claim.Credentials[repository.AndalalinStatusCredential]

	if !credential {
		// Return status 403 and permission denied error message.
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": true,
			"msg":   "Permission denied",
		})
		return
	}

	var andalalin models.Andalalin

	result := ac.DB.Model(&andalalin).Where("id_andalalin = ?", id).Update("status_andalalin", "Persayaratan terpenuhi")
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Permohonan tidak ditemukan"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (ac *AndalalinController) PersyaratanTidakSesuai(ctx *gin.Context) {
	id := ctx.Param("id_andalalin")
	var payload *models.PersayaratanTidakSesuai

	config, _ := initializers.LoadConfig(".")

	accessUser := ctx.MustGet("accessUser").(string)

	claim, error := utils.ValidateToken(accessUser, config.AccessTokenPublicKey)
	if error != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": error.Error()})
		return
	}

	credential := claim.Credentials[repository.AndalalinStatusCredential]

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

	result := ac.DB.Model(&andalalin).Where("id_andalalin = ?", id).Update("status_andalalin", "Persayaratan tidak sesuai")
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Permohonan tidak ditemukan"})
		return
	}

	var andalalinData models.Andalalin

	resultData := ac.DB.First(&andalalinData, "id_andalalin = ?", id)
	if resultData.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Permohonan tidak ditemukan"})
		return
	}

	data := utils.PersyaratanTidakSesuai{
		Nomer:       andalalinData.NomerAndalalin,
		Nama:        andalalinData.NamaPemohon,
		Alamat:      andalalinData.AlamatPemohon,
		Tlp:         andalalinData.NomerPemohon,
		Waktu:       andalalinData.WaktuAndalalin,
		Izin:        andalalinData.JenisAndalalin,
		Status:      andalalinData.StatusAndalalin,
		Persyaratan: payload.Persyaratan,
		Subject:     "Persyaratan Tidak Sesuai",
	}

	utils.SendEmailPersyaratan(andalalinData.EmailPemohon, &data)

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (ac *AndalalinController) UpdateStatusPermohonan(ctx *gin.Context) {
	status := ctx.Param("status")
	id := ctx.Param("id_andalalin")

	config, _ := initializers.LoadConfig(".")

	accessUser := ctx.MustGet("accessUser").(string)

	claim, error := utils.ValidateToken(accessUser, config.AccessTokenPublicKey)
	if error != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": error.Error()})
		return
	}

	credential := claim.Credentials[repository.AndalalinStatusCredential]

	if !credential {
		// Return status 403 and permission denied error message.
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": true,
			"msg":   "Permission denied",
		})
		return
	}

	var andalalin models.Andalalin

	result := ac.DB.Model(&andalalin).Where("id_andalalin = ?", id).Update("status_andalalin", status)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Permohonan tidak ditemukan"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (ac *AndalalinController) TambahPetugas(ctx *gin.Context) {
	var payload *models.TambahPetugas
	id := ctx.Param("id_andalalin")

	config, _ := initializers.LoadConfig(".")

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

	var andalalin models.Andalalin
	result := ac.DB.First(&andalalin, "id_andalalin = ?", id)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Permohonan tidak ditemukan"})
		return
	}

	andalalin.IdPetugas = payload.IdPetugas
	andalalin.NamaPetugas = payload.NamaPetugas
	andalalin.EmailPetugas = payload.EmailPetugas
	andalalin.StatusAndalalin = "Survey lapangan"

	ac.DB.Save(&andalalin)

	ac.ReleaseTicketLevel2(ctx, andalalin.IdAndalalin)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Tambah petugas berhasil"})
}

func (ac *AndalalinController) GantiPetugas(ctx *gin.Context) {
	var payload *models.TambahPetugas
	id := ctx.Param("id_andalalin")

	config, _ := initializers.LoadConfig(".")

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

	var andalalin models.Andalalin
	result := ac.DB.First(&andalalin, "id_andalalin = ?", id)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Permohonan tidak ditemukand"})
		return
	}

	andalalin.IdPetugas = payload.IdPetugas
	andalalin.NamaPetugas = payload.NamaPetugas
	andalalin.EmailPetugas = payload.EmailPetugas
	andalalin.StatusAndalalin = "Survey lapangan"

	ac.DB.Save(&andalalin)

	ac.CloseTiketLevel2(ctx, andalalin.IdAndalalin)

	ac.ReleaseTicketLevel2(ctx, andalalin.IdAndalalin)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Ubah petugas berhasil"})
}

func (ac *AndalalinController) GetAndalalinTicketLevel2(ctx *gin.Context) {
	status := ctx.Param("status")
	currentUser := ctx.MustGet("currentUser").(models.User)

	config, _ := initializers.LoadConfig(".")

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
			var andalalin models.Andalalin
			results := ac.DB.First(&andalalin, "id_andalalin = ? AND id_petugas = ?", s.IdAndalalin, currentUser.ID)

			if results.Error != nil {
				ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": results.Error})
				return
			}

			respone = append(respone, models.DaftarAndalalinResponse{
				IdAndalalin:     andalalin.IdAndalalin,
				NomerAndalalin:  andalalin.NomerAndalalin,
				WaktuAndalalin:  andalalin.WaktuAndalalin,
				Nama:            andalalin.NamaPemohon,
				Alamat:          andalalin.AlamatPemohon,
				JenisAndalalin:  andalalin.JenisAndalalin,
				StatusAndalalin: andalalin.StatusAndalalin,
			})
		}
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "results": len(respone), "data": respone})
	}

}

func (ac *AndalalinController) IsiSurvey(ctx *gin.Context) {
	var payload *models.DataSurvey
	currentUser := ctx.MustGet("currentUser").(models.User)
	id := ctx.Param("id_andalalin")

	config, _ := initializers.LoadConfig(".")

	accessUser := ctx.MustGet("accessUser").(string)

	claim, error := utils.ValidateToken(accessUser, config.AccessTokenPublicKey)
	if error != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": error.Error()})
		return
	}

	credential := claim.Credentials[repository.AndalalinOfficerCredential]

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

	var ticket1 models.TiketLevel1
	var ticket2 models.TiketLevel2

	resultTiket1 := ac.DB.Find(&ticket1, "id_andalalin = ?", id)
	resultTiket2 := ac.DB.Find(&ticket2, "id_andalalin = ?", id)
	if resultTiket1.Error != nil && resultTiket2.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Tiket tidak ditemukan"})
		return
	}

	var andalalin models.Andalalin

	resultAndalalin := ac.DB.First(&andalalin, "id_andalalin = ? AND id_petugas = ?", id, currentUser.ID)
	if resultAndalalin.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": resultAndalalin.Error})
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

	survey := models.Survey{
		IdAndalalin:   andalalin.IdAndalalin,
		IdTiketLevel1: ticket1.IdTiketLevel1,
		IdTiketLevel2: ticket2.IdTiketLevel2,
		IdPetugas:     currentUser.ID,
		Petugas:       currentUser.Name,
		Keterangan:    payload.Data.Keterangan,
		Foto1:         blobs["foto1"],
		Foto2:         blobs["foto2"],
		Foto3:         blobs["foto3"],
		Latitude:      payload.Data.Latitude,
		Longitude:     payload.Data.Longitude,
	}

	result := ac.DB.Create(&survey)

	if result.Error != nil && strings.Contains(result.Error.Error(), "duplicate key value violates unique") {
		ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "Data survey sudah tersedia"})
		return
	} else if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Telah terjadi sesuatu"})
		return
	}

	andalalin.StatusAndalalin = "Laporan BAP"

	ac.DB.Save(&andalalin)

	ac.CloseTiketLevel2(ctx, andalalin.IdAndalalin)

	ctx.JSON(http.StatusCreated, gin.H{"status": "success"})
}

func (ac *AndalalinController) GetSurvey(ctx *gin.Context) {
	id := ctx.Param("id_andalalin")

	config, _ := initializers.LoadConfig(".")

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

	var survey *models.Survey

	result := ac.DB.First(&survey, "id_andalalin = ?", id)
	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": survey})
}

func (ac *AndalalinController) LaporanBAP(ctx *gin.Context) {
	var payload *models.BAPData
	id := ctx.Param("id_andalalin")

	config, _ := initializers.LoadConfig(".")

	accessUser := ctx.MustGet("accessUser").(string)

	claim, error := utils.ValidateToken(accessUser, config.AccessTokenPublicKey)
	if error != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": error.Error()})
		return
	}

	credential := claim.Credentials[repository.AndalalinBAPCredential]

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

	file, err := ctx.FormFile("bap")
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

	var andalalin models.Andalalin

	resultAndalalin := ac.DB.First(&andalalin, "id_andalalin = ?", id)
	if resultAndalalin.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": resultAndalalin.Error})
		return
	}

	andalalin.NomerBAPDasar = payload.Data.NomerBAPDasar
	andalalin.NomerBAPPelaksanaan = payload.Data.NomerBAPPelaksanaan
	andalalin.TanggalBAP = payload.Data.TanggalBAP
	andalalin.FileBAP = data

	result := ac.DB.Save(&andalalin)

	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success"})
}

func (ac *AndalalinController) GetBAP(ctx *gin.Context) {
	id := ctx.Param("id_andalalin")

	config, _ := initializers.LoadConfig(".")

	accessUser := ctx.MustGet("accessUser").(string)

	claim, error := utils.ValidateToken(accessUser, config.AccessTokenPublicKey)
	if error != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": error.Error()})
		return
	}

	credential := claim.Credentials[repository.AndalalinBAPCredential]

	if !credential {
		// Return status 403 and permission denied error message.
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": true,
			"msg":   "Permission denied",
		})
		return
	}

	var andalalin models.Andalalin

	result := ac.DB.First(&andalalin, "id_andalalin = ?", id)
	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error})
		return
	}

	data := struct {
		NomerDasar       string `json:"nomer_dasar,omitempty"`
		NomerPelaksanaan string `json:"nomer_pelaksanaan,omitempty"`
		TanggalBAP       string `json:"tanggal,omitempty"`
		BAP              []byte `json:"bap,omitempty"`
	}{
		NomerDasar:       andalalin.NomerBAPDasar,
		NomerPelaksanaan: andalalin.NomerBAPPelaksanaan,
		TanggalBAP:       andalalin.TanggalBAP,
		BAP:              andalalin.FileBAP,
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": data})
}

func (ac *AndalalinController) PersetujuanDokumen(ctx *gin.Context) {
	var payload *models.Persetujuan
	id := ctx.Param("id_andalalin")

	config, _ := initializers.LoadConfig(".")

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
	andalalin.StatusAndalalin = "Pembuatan SK"

	ac.DB.Save(&andalalin)

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (ac *AndalalinController) GetPersetujuanDokumen(ctx *gin.Context) {
	id := ctx.Param("id_andalalin")

	config, _ := initializers.LoadConfig(".")

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

	var andalalin models.Andalalin

	result := ac.DB.First(&andalalin, "id_andalalin = ?", id)
	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error})
		return
	}

	data := struct {
		PersetujuanDokumen           string `json:"persetujuan,omitempty"`
		KeteranganPersetujuanDokumen string `json:"keterangan,omitempty"`
	}{
		PersetujuanDokumen:           andalalin.PersetujuanDokumen,
		KeteranganPersetujuanDokumen: andalalin.KeteranganPersetujuanDokumen,
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": data})
}

func (ac *AndalalinController) LaporanSK(ctx *gin.Context) {
	id := ctx.Param("id_andalalin")

	config, _ := initializers.LoadConfig(".")

	accessUser := ctx.MustGet("accessUser").(string)

	claim, error := utils.ValidateToken(accessUser, config.AccessTokenPublicKey)
	if error != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": error.Error()})
		return
	}

	credential := claim.Credentials[repository.AndalalinSKCredential]

	if !credential {
		// Return status 403 and permission denied error message.
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": true,
			"msg":   "Permission denied",
		})
		return
	}

	var andalalin models.Andalalin

	result := ac.DB.First(&andalalin, "id_andalalin = ?", id)
	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error})
		return
	}

	file, err := ctx.FormFile("sk")
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

	andalalin.FileSK = data

	resultSK := ac.DB.Save(&andalalin)

	if resultSK.Error != nil && strings.Contains(resultSK.Error.Error(), "duplicate key value violates unique") {
		ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "Data SK sudah tersedia"})
		return
	} else if resultSK.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Telah terjadi sesuatu"})
		return
	}

	ac.CloseTiketLevel1(ctx, andalalin.IdAndalalin)

	ac.PermohonanSelesai(ctx, andalalin.IdAndalalin)

	ctx.JSON(http.StatusCreated, gin.H{"status": "success"})
}

func (ac *AndalalinController) GetSK(ctx *gin.Context) {
	id := ctx.Param("id_andalalin")

	config, _ := initializers.LoadConfig(".")

	accessUser := ctx.MustGet("accessUser").(string)

	claim, error := utils.ValidateToken(accessUser, config.AccessTokenPublicKey)
	if error != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": error.Error()})
		return
	}

	credential := claim.Credentials[repository.AndalalinSKCredential]

	if !credential {
		// Return status 403 and permission denied error message.
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": true,
			"msg":   "Permission denied",
		})
		return
	}

	var andalalin models.Andalalin

	result := ac.DB.First(&andalalin, "id_andalalin = ?", id)
	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error})
		return
	}

	data := struct {
		Id uuid.UUID `json:"id,omitempty"`
		SK []byte    `json:"sk,omitempty"`
	}{
		Id: andalalin.IdAndalalin,
		SK: andalalin.FileSK,
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": data})
}

func (ac *AndalalinController) PermohonanSelesai(ctx *gin.Context, id uuid.UUID) {
	var andalalin models.Andalalin

	result := ac.DB.First(&andalalin, "id_andalalin = ?", id)
	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error})
		return
	}

	andalalin.StatusAndalalin = "Permohonan Selesai"

	ac.DB.Save(&andalalin)

	data := utils.PermohonanSelesai{
		Nomer:   andalalin.NomerAndalalin,
		Nama:    andalalin.NamaPemohon,
		Alamat:  andalalin.AlamatPemohon,
		Tlp:     andalalin.NomerPemohon,
		Waktu:   andalalin.WaktuAndalalin,
		Izin:    andalalin.JenisAndalalin,
		Status:  andalalin.StatusAndalalin,
		Subject: "Permohonan telah selesai",
	}

	utils.SendEmailPermohonanSelesai(andalalin.EmailPemohon, &data)
}