package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/Ikit777/E-Andalalin/initializers"
	"github.com/Ikit777/E-Andalalin/models"
	"github.com/Ikit777/E-Andalalin/utils"
	"github.com/tealeg/xlsx"

	_ "time/tzdata"
)

func init() {
	config, err := initializers.LoadConfig()
	if err != nil {
		log.Fatal("Could not load environment variables", err)
	}

	initializers.ConnectDB(&config)
}

func removeExtension(fileName string) string {
	return path.Base(fileName[:len(fileName)-len(path.Ext(fileName))])
}

func main() {
	initializers.DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")

	initializers.DB.Migrator().DropTable(&models.User{})
	initializers.DB.Migrator().DropTable(&models.Andalalin{})
	initializers.DB.Migrator().DropTable(&models.Perlalin{})
	initializers.DB.Migrator().DropTable(&models.Survei{})
	initializers.DB.Migrator().DropTable(&models.SurveiMandiri{})
	initializers.DB.Migrator().DropTable(&models.TiketLevel1{})
	initializers.DB.Migrator().DropTable(&models.TiketLevel2{})
	initializers.DB.Migrator().DropTable(&models.Notifikasi{})
	initializers.DB.Migrator().DropTable(&models.DataMaster{})
	initializers.DB.Migrator().DropTable(&models.UsulanPengelolaan{})
	initializers.DB.Migrator().DropTable(&models.SurveiKepuasan{})
	initializers.DB.Migrator().DropTable(&models.Pemasangan{})

	initializers.DB.AutoMigrate(&models.User{})
	initializers.DB.AutoMigrate(&models.Andalalin{})
	initializers.DB.AutoMigrate(&models.Perlalin{})
	initializers.DB.AutoMigrate(&models.Survei{})
	initializers.DB.AutoMigrate(&models.SurveiMandiri{})
	initializers.DB.AutoMigrate(&models.TiketLevel1{})
	initializers.DB.AutoMigrate(&models.TiketLevel2{})
	initializers.DB.AutoMigrate(&models.Notifikasi{})
	initializers.DB.AutoMigrate(&models.DataMaster{})
	initializers.DB.AutoMigrate(&models.UsulanPengelolaan{})
	initializers.DB.AutoMigrate(&models.SurveiKepuasan{})
	initializers.DB.AutoMigrate(&models.Pemasangan{})

	loc, _ := time.LoadLocation("Asia/Singapore")
	now := time.Now().In(loc).Format("02-01-2006")
	hashedPassword, err := utils.HashPassword("superadmin")
	if err != nil {
		return
	}

	filePath := "assets/default.png"
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal("Error reading the file:", err)
	}

	initializers.DB.Create(&models.User{
		Name:      "Super admin",
		Email:     strings.ToLower("ssuper.adm1n@gmail.com"),
		Password:  hashedPassword,
		Role:      "Super Admin",
		Photo:     fileData,
		Verified:  true,
		CreatedAt: now,
		UpdatedAt: now,
	})

	jenis_proyek := []string{"Pembangunan", "Pengembangan", "Operasional"}

	lokasi := []string{"Banjarmasin"}

	jenis_kegiatan := []string{"Pusat kegiatan", "Pemukiman", "Infrastruktur"}

	pusat_kegiatan := []models.JenisRencana{}
	pusat_kegiatan = append(pusat_kegiatan, models.JenisRencana{Jenis: "Pusat perbelanjaan atau retail", Kriteria: "Luas lantai bangunan", Satuan: "m²"})
	pusat_kegiatan = append(pusat_kegiatan, models.JenisRencana{Jenis: "Perkantoran", Kriteria: "Luas lantai bangunan", Satuan: "m²"})
	pusat_kegiatan = append(pusat_kegiatan, models.JenisRencana{Jenis: "Industri dan pergudangan", Kriteria: "Luas lantai bangunan", Satuan: "m²"})
	pusat_kegiatan = append(pusat_kegiatan, models.JenisRencana{Jenis: "Sekolah atau universitas", Kriteria: "Jumlah siswa", Satuan: "siswa"})
	pusat_kegiatan = append(pusat_kegiatan, models.JenisRencana{Jenis: "Lembaga kursus", Kriteria: "Jumlah siswa dalam bangunan", Satuan: "siswa"})
	pusat_kegiatan = append(pusat_kegiatan, models.JenisRencana{Jenis: "Rumah sakit", Kriteria: "Jumlah tempat tidur", Satuan: "tempat tidur"})
	pusat_kegiatan = append(pusat_kegiatan, models.JenisRencana{Jenis: "Klinik bersama", Kriteria: "Jumlah ruang praktek dokter", Satuan: "ruang praktek dokter"})
	pusat_kegiatan = append(pusat_kegiatan, models.JenisRencana{Jenis: "Bank", Kriteria: "Luas lantai bangunan", Satuan: "m²"})
	pusat_kegiatan = append(pusat_kegiatan, models.JenisRencana{Jenis: "Stasiun pengisin bahan bakar", Kriteria: "Jumlah dispenser", Satuan: "dispenser"})
	pusat_kegiatan = append(pusat_kegiatan, models.JenisRencana{Jenis: "Hotel", Kriteria: "Jumlah kamar", Satuan: "kamar"})
	pusat_kegiatan = append(pusat_kegiatan, models.JenisRencana{Jenis: "Gedung pertemuan", Kriteria: "Luas lantai bangunan", Satuan: "m²"})
	pusat_kegiatan = append(pusat_kegiatan, models.JenisRencana{Jenis: "Restoran", Kriteria: "Jumlah tempat duduk", Satuan: "tempat duduk"})
	pusat_kegiatan = append(pusat_kegiatan, models.JenisRencana{Jenis: "Fasilitan olah raga", Kriteria: "Jumlah kapasitas penonton", Satuan: "orang"})
	pusat_kegiatan = append(pusat_kegiatan, models.JenisRencana{Jenis: "Bengkel kendaraan bermotor", Kriteria: "Luas lantai bangunan", Satuan: "m²"})
	pusat_kegiatan = append(pusat_kegiatan, models.JenisRencana{Jenis: "Pencucian mobil", Kriteria: "Luas lantai bangunan", Satuan: "m²"})

	infrastruktur := []models.JenisRencana{}
	infrastruktur = append(infrastruktur, models.JenisRencana{Jenis: "Akses ke dan dari jalan tol", Kriteria: "", Satuan: ""})
	infrastruktur = append(infrastruktur, models.JenisRencana{Jenis: "Pelabuhan", Kriteria: "", Satuan: ""})
	infrastruktur = append(infrastruktur, models.JenisRencana{Jenis: "Bandar udara", Kriteria: "", Satuan: ""})
	infrastruktur = append(infrastruktur, models.JenisRencana{Jenis: "Terminal", Kriteria: "", Satuan: ""})
	infrastruktur = append(infrastruktur, models.JenisRencana{Jenis: "Stasiun kereta api", Kriteria: "", Satuan: ""})
	infrastruktur = append(infrastruktur, models.JenisRencana{Jenis: "Pool kendaraan", Kriteria: "", Satuan: ""})
	infrastruktur = append(infrastruktur, models.JenisRencana{Jenis: "Fasilitas parkir umum", Kriteria: "", Satuan: ""})
	infrastruktur = append(infrastruktur, models.JenisRencana{Jenis: "Flyover", Kriteria: "", Satuan: ""})
	infrastruktur = append(infrastruktur, models.JenisRencana{Jenis: "Underpass", Kriteria: "", Satuan: ""})
	infrastruktur = append(infrastruktur, models.JenisRencana{Jenis: "Terowongan", Kriteria: "", Satuan: ""})

	pemukiman := []models.JenisRencana{}
	pemukiman = append(pemukiman, models.JenisRencana{Jenis: "Perumahan sederhana", Kriteria: "Jumlah unit", Satuan: "unit"})
	pemukiman = append(pemukiman, models.JenisRencana{Jenis: "Perumahan menengan-atas", Kriteria: "Jumlah unit", Satuan: "unit"})
	pemukiman = append(pemukiman, models.JenisRencana{Jenis: "Rumah susun sederhana", Kriteria: "Jumlah unit", Satuan: "unit"})
	pemukiman = append(pemukiman, models.JenisRencana{Jenis: "Apartemen", Kriteria: "Jumlah unit", Satuan: "unit"})
	pemukiman = append(pemukiman, models.JenisRencana{Jenis: "Asrama", Kriteria: "Jumlah unit", Satuan: "unit"})
	pemukiman = append(pemukiman, models.JenisRencana{Jenis: "Ruko", Kriteria: "Luas lahan keseluruhan", Satuan: "m²"})

	rencana := []models.Rencana{}
	rencana = append(rencana, models.Rencana{Kategori: "Pusat kegiatan", JenisRencana: pusat_kegiatan})
	rencana = append(rencana, models.Rencana{Kategori: "Pemukiman", JenisRencana: pemukiman})
	rencana = append(rencana, models.Rencana{Kategori: "Infrastruktur", JenisRencana: infrastruktur})

	ketegori_perlengkapan := []string{"Rambu peringatan", "Rambu larangan", "Rambu perintah", "Rambu petunjuk"}

	andalalin := []models.PersyaratanAndalalinInput{}
	andalalin = append(andalalin, models.PersyaratanAndalalinInput{Bangkitan: "Bangkitan rendah", Persyaratan: "Surat permohonan persetujuan andalalin", KeteranganPersyaratan: "Surat permohonan persetujuan analisis dampak lalu lintas adalah surat yang digunakan untuk mengajukan permohonan kepada pihak yang berwenang, biasanya pemerintah daerah, untuk mendapatkan persetujuan atau izin terkait dengan rencana atau proyek tertentu. Surat ini harus memuat informasi yang lengkap dan jelas mengenai rencana atau proyek yang diajukan, termasuk tujuan, dampak lingkungan, serta segala persyaratan yang harus dipenuhi."})
	andalalin = append(andalalin, models.PersyaratanAndalalinInput{Bangkitan: "Bangkitan rendah", Persyaratan: "Kartu tanda penduduk atau paspor atau akta pendirian badan usaha", KeteranganPersyaratan: "Kartu tanda penduduk/Paspor/Akta pendirian badan usaha adalah dokumen identitas yang bertujuan untuk mengindentifikasi pemohon atau penanggung jawab terhadap permohonan."})
	andalalin = append(andalalin, models.PersyaratanAndalalinInput{Bangkitan: "Bangkitan rendah", Persyaratan: "NPWP", KeteranganPersyaratan: "Nomor Pokok Wajib Pajak (NPWP) adalah nomor identifikasi pajak yang diberikan kepada warga negara atau entitas yang wajib membayar pajak di Indonesia. NPWP adalah bagian penting dari administrasi pajak di Indonesia dan dikeluarkan oleh Direktorat Jenderal Pajak (DJP) yang merupakan lembaga di bawah Kementerian Keuangan Republik Indonesia. NPWP digunakan untuk mengidentifikasi wajib pajak dan melacak pembayaran pajak yang dikenakan."})
	andalalin = append(andalalin, models.PersyaratanAndalalinInput{Bangkitan: "Bangkitan rendah", Persyaratan: "Surat bukti kepemilikan atau Penguasaan lahan", KeteranganPersyaratan: "Surat bukti kepemilikan atau penguasaan lahan adalah dokumen yang dimaksudkan untuk menunjukkan hak hukum seseorang atau badan usaha atas sebidang lahan atau properti tertentu."})
	andalalin = append(andalalin, models.PersyaratanAndalalinInput{Bangkitan: "Bangkitan rendah", Persyaratan: "Surat kesesuaian tata ruang dan atau izin pemanfaatan ruang", KeteranganPersyaratan: "Surat kesesuaian tata ruang dan/atau Izin pemanfaatan ruang adalah dokumen yang diterbitkan oleh pihak berwenang, seperti pemerintah daerah atau instansi terkait, untuk memberikan izin atau persetujuan terkait dengan penggunaan dan pengembangan lahan atau ruang tertentu."})
	andalalin = append(andalalin, models.PersyaratanAndalalinInput{Bangkitan: "Bangkitan rendah", Persyaratan: "Gambar tata letak bangunan (site plan)", KeteranganPersyaratan: "Tata letak bangunan, atau dalam bahasa Inggris dikenal sebagai site plan, adalah gambar atau diagram yang menunjukkan tata letak fisik bangunan, struktur, dan fasilitas terkait dalam suatu area tertentu. Site plan biasanya digunakan dalam perencanaan konstruksi, pengembangan properti, perizinan bangunan, dan perencanaan tata ruang."})
	andalalin = append(andalalin, models.PersyaratanAndalalinInput{Bangkitan: "Bangkitan rendah", Persyaratan: "DED Bangunan yang diusulkan", KeteranganPersyaratan: "Dokumen Engineering Design (DED) adalah bagian penting dari proses perencanaan dan konstruksi bangunan. DED untuk bangunan yang diusulkan adalah dokumen yang merinci dan mendokumentasikan desain teknik dan teknis dari bangunan yang akan dibangun. DED adalah langkah selanjutnya setelah tahap perencanaan dan perancangan awal, dan sebelum memulai proses konstruksi."})
	andalalin = append(andalalin, models.PersyaratanAndalalinInput{Bangkitan: "Bangkitan rendah", Persyaratan: "Foto kondisi eksisting lapangan terkini", KeteranganPersyaratan: "Foto kondisi eksisting lapangan yang terkini adalah gambar-gambar yang menggambarkan kondisi aktual dan terbaru dari lapangan atau area tertentu pada suatu waktu. Foto-foto ini sangat berguna dalam berbagai konteks, termasuk dalam proyek konstruksi, pengembangan properti, pemantauan lingkungan, dan penelitian."})
	andalalin = append(andalalin, models.PersyaratanAndalalinInput{Bangkitan: "Bangkitan rendah", Persyaratan: "MOU Kerjsa sama", KeteranganPersyaratan: "Memorandum of Understanding (MOU), atau dalam bahasa Indonesia disebut Nota Kesepahaman, adalah dokumen tertulis yang digunakan untuk mendefinisikan kerjasama antara dua pihak atau lebih dalam suatu proyek atau inisiatif. MOU tidak memiliki kekuatan hukum yang sama dengan kontrak, tetapi berfungsi sebagai dasar kerjasama dan kerangka kerja awal."})

	andalalin = append(andalalin, models.PersyaratanAndalalinInput{Bangkitan: "Bangkitan sedang", Persyaratan: "Surat permohonan persetujuan andalalin", KeteranganPersyaratan: "Surat permohonan persetujuan analisis dampak lalu lintas adalah surat yang digunakan untuk mengajukan permohonan kepada pihak yang berwenang, biasanya pemerintah daerah, untuk mendapatkan persetujuan atau izin terkait dengan rencana atau proyek tertentu. Surat ini harus memuat informasi yang lengkap dan jelas mengenai rencana atau proyek yang diajukan, termasuk tujuan, dampak lingkungan, serta segala persyaratan yang harus dipenuhi."})
	andalalin = append(andalalin, models.PersyaratanAndalalinInput{Bangkitan: "Bangkitan sedang", Persyaratan: "Kartu tanda penduduk atau paspor atau akta pendirian badan usaha", KeteranganPersyaratan: "Kartu tanda penduduk/Paspor/Akta pendirian badan usaha adalah dokumen identitas yang bertujuan untuk mengindentifikasi pemohon atau penanggung jawab terhadap permohonan."})
	andalalin = append(andalalin, models.PersyaratanAndalalinInput{Bangkitan: "Bangkitan sedang", Persyaratan: "NPWP", KeteranganPersyaratan: "Nomor Pokok Wajib Pajak (NPWP) adalah nomor identifikasi pajak yang diberikan kepada warga negara atau entitas yang wajib membayar pajak di Indonesia. NPWP adalah bagian penting dari administrasi pajak di Indonesia dan dikeluarkan oleh Direktorat Jenderal Pajak (DJP) yang merupakan lembaga di bawah Kementerian Keuangan Republik Indonesia. NPWP digunakan untuk mengidentifikasi wajib pajak dan melacak pembayaran pajak yang dikenakan."})
	andalalin = append(andalalin, models.PersyaratanAndalalinInput{Bangkitan: "Bangkitan sedang", Persyaratan: "Sertifikat konsultan atau tenaga ahli penyusun andalalin sesuai klasifikasi", KeteranganPersyaratan: "Sertifikat konsultan atau tenaga ahli penyusun andalalin sesuai klasifikasi adalah dokumen yang menunjukkan bahwa seorang profesional atau konsultan memiliki kualifikasi, kompetensi, dan sertifikasi yang sesuai untuk menyusun dokumen Andalalin. Andalalin adalah dokumen yang berisi analisis dampak lingkungan (AMDAL) yang diperlukan dalam perencanaan dan pelaksanaan proyek-proyek yang memiliki potensi dampak terhadap lingkungan."})
	andalalin = append(andalalin, models.PersyaratanAndalalinInput{Bangkitan: "Bangkitan sedang", Persyaratan: "Surat bukti kepemilikan atau penguasaan lahan", KeteranganPersyaratan: "Surat bukti kepemilikan atau penguasaan lahan adalah dokumen yang dimaksudkan untuk menunjukkan hak hukum seseorang atau badan usaha atas sebidang lahan atau properti tertentu."})
	andalalin = append(andalalin, models.PersyaratanAndalalinInput{Bangkitan: "Bangkitan sedang", Persyaratan: "Surat kesesuaian tata ruang dan atau izin pemanfaatan ruang", KeteranganPersyaratan: "Surat kesesuaian tata ruang dan/atau Izin pemanfaatan ruang adalah dokumen yang diterbitkan oleh pihak berwenang, seperti pemerintah daerah atau instansi terkait, untuk memberikan izin atau persetujuan terkait dengan penggunaan dan pengembangan lahan atau ruang tertentu."})
	andalalin = append(andalalin, models.PersyaratanAndalalinInput{Bangkitan: "Bangkitan sedang", Persyaratan: "Gambar tata letak bangunan (site plan)", KeteranganPersyaratan: "Tata letak bangunan, atau dalam bahasa Inggris dikenal sebagai site plan, adalah gambar atau diagram yang menunjukkan tata letak fisik bangunan, struktur, dan fasilitas terkait dalam suatu area tertentu. Site plan biasanya digunakan dalam perencanaan konstruksi, pengembangan properti, perizinan bangunan, dan perencanaan tata ruang."})
	andalalin = append(andalalin, models.PersyaratanAndalalinInput{Bangkitan: "Bangkitan sedang", Persyaratan: "DED Bangunan yang diusulkan", KeteranganPersyaratan: "Dokumen Engineering Design (DED) adalah bagian penting dari proses perencanaan dan konstruksi bangunan. DED untuk bangunan yang diusulkan adalah dokumen yang merinci dan mendokumentasikan desain teknik dan teknis dari bangunan yang akan dibangun. DED adalah langkah selanjutnya setelah tahap perencanaan dan perancangan awal, dan sebelum memulai proses konstruksi."})
	andalalin = append(andalalin, models.PersyaratanAndalalinInput{Bangkitan: "Bangkitan sedang", Persyaratan: "Foto kondisi eksisting lapangan terkini", KeteranganPersyaratan: "Foto kondisi eksisting lapangan yang terkini adalah gambar-gambar yang menggambarkan kondisi aktual dan terbaru dari lapangan atau area tertentu pada suatu waktu. Foto-foto ini sangat berguna dalam berbagai konteks, termasuk dalam proyek konstruksi, pengembangan properti, pemantauan lingkungan, dan penelitian."})
	andalalin = append(andalalin, models.PersyaratanAndalalinInput{Bangkitan: "Bangkitan sedang", Persyaratan: "MOU Kerjsa sama", KeteranganPersyaratan: "Memorandum of Understanding (MOU), atau dalam bahasa Indonesia disebut Nota Kesepahaman, adalah dokumen tertulis yang digunakan untuk mendefinisikan kerjasama antara dua pihak atau lebih dalam suatu proyek atau inisiatif. MOU tidak memiliki kekuatan hukum yang sama dengan kontrak, tetapi berfungsi sebagai dasar kerjasama dan kerangka kerja awal."})

	andalalin = append(andalalin, models.PersyaratanAndalalinInput{Bangkitan: "Bangkitan tinggi", Persyaratan: "Surat permohonan persetujuan andalalin", KeteranganPersyaratan: "Surat permohonan persetujuan analisis dampak lalu lintas adalah surat yang digunakan untuk mengajukan permohonan kepada pihak yang berwenang, biasanya pemerintah daerah, untuk mendapatkan persetujuan atau izin terkait dengan rencana atau proyek tertentu. Surat ini harus memuat informasi yang lengkap dan jelas mengenai rencana atau proyek yang diajukan, termasuk tujuan, dampak lingkungan, serta segala persyaratan yang harus dipenuhi."})
	andalalin = append(andalalin, models.PersyaratanAndalalinInput{Bangkitan: "Bangkitan tinggi", Persyaratan: "Kartu tanda penduduk atau paspor atau akta pendirian badan usaha", KeteranganPersyaratan: "Kartu tanda penduduk/Paspor/Akta pendirian badan usaha adalah dokumen identitas yang bertujuan untuk mengindentifikasi pemohon atau penanggung jawab terhadap permohonan."})
	andalalin = append(andalalin, models.PersyaratanAndalalinInput{Bangkitan: "Bangkitan tinggi", Persyaratan: "NPWP", KeteranganPersyaratan: "Nomor Pokok Wajib Pajak (NPWP) adalah nomor identifikasi pajak yang diberikan kepada warga negara atau entitas yang wajib membayar pajak di Indonesia. NPWP adalah bagian penting dari administrasi pajak di Indonesia dan dikeluarkan oleh Direktorat Jenderal Pajak (DJP) yang merupakan lembaga di bawah Kementerian Keuangan Republik Indonesia. NPWP digunakan untuk mengidentifikasi wajib pajak dan melacak pembayaran pajak yang dikenakan."})
	andalalin = append(andalalin, models.PersyaratanAndalalinInput{Bangkitan: "Bangkitan tinggi", Persyaratan: "Sertifikat konsultan atau tenaga ahli penyusun andalalin sesuai klasifikasi", KeteranganPersyaratan: "Sertifikat konsultan atau tenaga ahli penyusun andalalin sesuai klasifikasi adalah dokumen yang menunjukkan bahwa seorang profesional atau konsultan memiliki kualifikasi, kompetensi, dan sertifikasi yang sesuai untuk menyusun dokumen Andalalin. Andalalin adalah dokumen yang berisi analisis dampak lingkungan (AMDAL) yang diperlukan dalam perencanaan dan pelaksanaan proyek-proyek yang memiliki potensi dampak terhadap lingkungan."})
	andalalin = append(andalalin, models.PersyaratanAndalalinInput{Bangkitan: "Bangkitan tinggi", Persyaratan: "Surat bukti kepemilikan atau penguasaan lahan", KeteranganPersyaratan: "Surat bukti kepemilikan atau penguasaan lahan adalah dokumen yang dimaksudkan untuk menunjukkan hak hukum seseorang atau badan usaha atas sebidang lahan atau properti tertentu."})
	andalalin = append(andalalin, models.PersyaratanAndalalinInput{Bangkitan: "Bangkitan tinggi", Persyaratan: "Surat kesesuaian tata ruang dan atau izin pemanfaatan ruang", KeteranganPersyaratan: "Surat kesesuaian tata ruang dan/atau Izin pemanfaatan ruang adalah dokumen yang diterbitkan oleh pihak berwenang, seperti pemerintah daerah atau instansi terkait, untuk memberikan izin atau persetujuan terkait dengan penggunaan dan pengembangan lahan atau ruang tertentu."})
	andalalin = append(andalalin, models.PersyaratanAndalalinInput{Bangkitan: "Bangkitan tinggi", Persyaratan: "Gambar tata letak bangunan (site plan)", KeteranganPersyaratan: "Tata letak bangunan, atau dalam bahasa Inggris dikenal sebagai site plan, adalah gambar atau diagram yang menunjukkan tata letak fisik bangunan, struktur, dan fasilitas terkait dalam suatu area tertentu. Site plan biasanya digunakan dalam perencanaan konstruksi, pengembangan properti, perizinan bangunan, dan perencanaan tata ruang."})
	andalalin = append(andalalin, models.PersyaratanAndalalinInput{Bangkitan: "Bangkitan tinggi", Persyaratan: "DED Bangunan yang diusulkan", KeteranganPersyaratan: "Dokumen Engineering Design (DED) adalah bagian penting dari proses perencanaan dan konstruksi bangunan. DED untuk bangunan yang diusulkan adalah dokumen yang merinci dan mendokumentasikan desain teknik dan teknis dari bangunan yang akan dibangun. DED adalah langkah selanjutnya setelah tahap perencanaan dan perancangan awal, dan sebelum memulai proses konstruksi."})
	andalalin = append(andalalin, models.PersyaratanAndalalinInput{Bangkitan: "Bangkitan tinggi", Persyaratan: "Foto kondisi eksisting lapangan terkini", KeteranganPersyaratan: "Foto kondisi eksisting lapangan yang terkini adalah gambar-gambar yang menggambarkan kondisi aktual dan terbaru dari lapangan atau area tertentu pada suatu waktu. Foto-foto ini sangat berguna dalam berbagai konteks, termasuk dalam proyek konstruksi, pengembangan properti, pemantauan lingkungan, dan penelitian."})
	andalalin = append(andalalin, models.PersyaratanAndalalinInput{Bangkitan: "Bangkitan tinggi", Persyaratan: "MOU Kerjsa sama", KeteranganPersyaratan: "Memorandum of Understanding (MOU), atau dalam bahasa Indonesia disebut Nota Kesepahaman, adalah dokumen tertulis yang digunakan untuk mendefinisikan kerjasama antara dua pihak atau lebih dalam suatu proyek atau inisiatif. MOU tidak memiliki kekuatan hukum yang sama dengan kontrak, tetapi berfungsi sebagai dasar kerjasama dan kerangka kerja awal."})

	perlalin := []models.PersyaratanPerlalinInput{}
	perlalin = append(perlalin, models.PersyaratanPerlalinInput{Persyaratan: "Kartu tanda penduduk", KeteranganPersyaratan: "Kartu tanda penduduk adalah dokumen identitas yang bertujuan untuk mengindentifikasi pemohon atau penanggung jawab terhadap permohonan."})
	perlalin = append(perlalin, models.PersyaratanPerlalinInput{Persyaratan: "Surat permohonan", KeteranganPersyaratan: "Surat permohonan berisi infromasi terkait permohonan yang akan diajukan seperti perlengkapan lalu lintas"})

	persyaratan := models.Persyaratan{
		PersyaratanAndalalin: andalalin,
		PersyaratanPerlalin:  perlalin,
	}

	perlengkapanPeringatan := []models.PerlengkapanItem{}

	folderPeringatan := "assets/Perlalin/Peringatan"

	folder1, err := os.Open(folderPeringatan)
	if err != nil {
		fmt.Println("Error opening folder:", err)
		return
	}
	defer folder1.Close()

	filePeringatan, err := folder1.Readdir(0)
	if err != nil {
		fmt.Println("Error reading folder contents:", err)
		return
	}

	for _, fileInfo := range filePeringatan {
		if fileInfo.Mode().IsRegular() {
			filePath := filepath.Join(folderPeringatan, fileInfo.Name())
			content, err := os.ReadFile(filePath)
			if err != nil {
				fmt.Printf("Error reading file %s: %v\n", fileInfo.Name(), err)
				continue
			}
			perlengkapanPeringatan = append(perlengkapanPeringatan, models.PerlengkapanItem{JenisPerlengkapan: removeExtension(fileInfo.Name()), GambarPerlengkapan: content})
		}
	}

	perlengkapanLarangan := []models.PerlengkapanItem{}

	folderLarangan := "assets/Perlalin/Larangan"

	folder2, err := os.Open(folderLarangan)
	if err != nil {
		fmt.Println("Error opening folder:", err)
		return
	}
	defer folder2.Close()

	fileLarangan, err := folder2.Readdir(0)
	if err != nil {
		fmt.Println("Error reading folder contents:", err)
		return
	}

	for _, fileInfo := range fileLarangan {
		if fileInfo.Mode().IsRegular() {
			filePath := filepath.Join(folderLarangan, fileInfo.Name())
			content, err := os.ReadFile(filePath)
			if err != nil {
				fmt.Printf("Error reading file %s: %v\n", fileInfo.Name(), err)
				continue
			}
			perlengkapanLarangan = append(perlengkapanLarangan, models.PerlengkapanItem{JenisPerlengkapan: removeExtension(fileInfo.Name()), GambarPerlengkapan: content})
		}
	}

	perlengkapanPerintah := []models.PerlengkapanItem{}

	folderPerintah := "assets/Perlalin/Perintah"

	folder3, err := os.Open(folderPerintah)
	if err != nil {
		fmt.Println("Error opening folder:", err)
		return
	}
	defer folder3.Close()

	filePerintah, err := folder3.Readdir(0)
	if err != nil {
		fmt.Println("Error reading folder contents:", err)
		return
	}

	for _, fileInfo := range filePerintah {
		if fileInfo.Mode().IsRegular() {
			filePath := filepath.Join(folderPerintah, fileInfo.Name())
			content, err := os.ReadFile(filePath)
			if err != nil {
				fmt.Printf("Error reading file %s: %v\n", fileInfo.Name(), err)
				continue
			}
			perlengkapanPerintah = append(perlengkapanPerintah, models.PerlengkapanItem{JenisPerlengkapan: removeExtension(fileInfo.Name()), GambarPerlengkapan: content})
		}
	}

	perlengkapanPetunjuk := []models.PerlengkapanItem{}

	folderPetunjuk := "assets/Perlalin/Petunjuk"

	folder4, err := os.Open(folderPetunjuk)
	if err != nil {
		fmt.Println("Error opening folder:", err)
		return
	}
	defer folder4.Close()

	filePetunjuk, err := folder4.Readdir(0)
	if err != nil {
		fmt.Println("Error reading folder contents:", err)
		return
	}

	for _, fileInfo := range filePetunjuk {
		if fileInfo.Mode().IsRegular() {
			filePath := filepath.Join(folderPetunjuk, fileInfo.Name())
			content, err := os.ReadFile(filePath)
			if err != nil {
				fmt.Printf("Error reading file %s: %v\n", fileInfo.Name(), err)
				continue
			}
			perlengkapanPetunjuk = append(perlengkapanPetunjuk, models.PerlengkapanItem{JenisPerlengkapan: removeExtension(fileInfo.Name()), GambarPerlengkapan: content})
		}
	}

	perlengkapan := []models.JenisPerlengkapan{}
	perlengkapan = append(perlengkapan, models.JenisPerlengkapan{Kategori: "Rambu peringatan", Perlengkapan: perlengkapanPeringatan})
	perlengkapan = append(perlengkapan, models.JenisPerlengkapan{Kategori: "Rambu larangan", Perlengkapan: perlengkapanLarangan})
	perlengkapan = append(perlengkapan, models.JenisPerlengkapan{Kategori: "Rambu perintah", Perlengkapan: perlengkapanPerintah})
	perlengkapan = append(perlengkapan, models.JenisPerlengkapan{Kategori: "Rambu petunjuk", Perlengkapan: perlengkapanPetunjuk})

	fileProvinsi, err := os.Open("assets/Indonesia/provinces.csv")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer fileProvinsi.Close()

	csvProvinsi := csv.NewReader(fileProvinsi)

	var provinsi []models.Provinsi

	for {
		record, err := csvProvinsi.Read()
		if err != nil {
			break // End of file
		}

		provinsi = append(provinsi, models.Provinsi{Id: record[0], Name: record[1]})
	}

	fileKabupaten, err := os.Open("assets/Indonesia/regencies.csv")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer fileKabupaten.Close()

	csvKabupaten := csv.NewReader(fileKabupaten)

	var Kabupaten []models.Kabupaten

	for {
		record, err := csvKabupaten.Read()
		if err != nil {
			break // End of file
		}

		Kabupaten = append(Kabupaten, models.Kabupaten{Id: record[0], IdProvinsi: record[1], Name: record[2]})
	}

	fileKecamatan, err := os.Open("assets/Indonesia/districts.csv")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer fileKecamatan.Close()

	csvKecamatan := csv.NewReader(fileKecamatan)

	var kecamatan []models.Kecamatan

	for {
		record, err := csvKecamatan.Read()
		if err != nil {
			break // End of file
		}

		kecamatan = append(kecamatan, models.Kecamatan{Id: record[0], IdKabupaten: record[1], Name: record[2]})
	}

	fileKelurahan, err := os.Open("assets/Indonesia/villages.csv")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer fileKelurahan.Close()

	csvKelurahan := csv.NewReader(fileKelurahan)
	csvKelurahan.Comma = ','

	var kelurahan []models.Kelurahan

	for {
		record, err := csvKelurahan.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}

		kelurahan = append(kelurahan, models.Kelurahan{Id: record[0], IdKecamatan: record[1], Name: record[2]})
	}

	jalan := []models.Jalan{}

	xlFile, err := xlsx.OpenFile("assets/Jalan/jalan.xlsx")
	if err != nil {
		log.Fatal(err)
	}

	sheetName := "Lembar1"
	sheet, ok := xlFile.Sheet[sheetName]
	if !ok {
		log.Fatalf("Sheet not found: %s", sheetName)
	}

	rowIndex1 := 1
	row1 := sheet.Rows[rowIndex1]

	for _, cell := range row1.Cells {
		jalan = append(jalan, models.Jalan{KodeProvinsi: cell.String()})
	}

	rowIndex2 := 2
	row2 := sheet.Rows[rowIndex2]

	for _, cell := range row2.Cells {
		jalan = append(jalan, models.Jalan{KodeKabupaten: cell.String()})
	}

	rowIndex3 := 3
	row3 := sheet.Rows[rowIndex3]

	for _, cell := range row3.Cells {
		jalan = append(jalan, models.Jalan{KodeKecamatan: cell.String()})
	}

	rowIndex4 := 4
	row4 := sheet.Rows[rowIndex4]

	for _, cell := range row4.Cells {
		jalan = append(jalan, models.Jalan{KodeKelurahan: cell.String()})
	}

	rowIndex5 := 5
	row5 := sheet.Rows[rowIndex5]

	for _, cell := range row5.Cells {
		jalan = append(jalan, models.Jalan{KodeJalan: cell.String()})
	}

	rowIndex6 := 6
	row6 := sheet.Rows[rowIndex6]

	for _, cell := range row6.Cells {
		jalan = append(jalan, models.Jalan{Nama: cell.String()})
	}

	rowIndex7 := 7
	row7 := sheet.Rows[rowIndex7]

	for _, cell := range row7.Cells {
		jalan = append(jalan, models.Jalan{Pangkal: cell.String()})
	}

	rowIndex8 := 8
	row8 := sheet.Rows[rowIndex8]

	for _, cell := range row8.Cells {
		jalan = append(jalan, models.Jalan{Ujung: cell.String()})
	}

	rowIndex9 := 9
	row9 := sheet.Rows[rowIndex9]

	for _, cell := range row9.Cells {
		jalan = append(jalan, models.Jalan{Kelurahan: cell.String()})
	}

	rowIndex10 := 10
	row10 := sheet.Rows[rowIndex10]

	for _, cell := range row10.Cells {
		jalan = append(jalan, models.Jalan{Kecamatan: cell.String()})
	}

	rowIndex11 := 11
	row11 := sheet.Rows[rowIndex11]

	for _, cell := range row11.Cells {
		jalan = append(jalan, models.Jalan{Panjang: cell.String()})
	}

	rowIndex12 := 12
	row12 := sheet.Rows[rowIndex12]

	for _, cell := range row12.Cells {
		jalan = append(jalan, models.Jalan{Lebar: cell.String()})
	}

	rowIndex13 := 13
	row13 := sheet.Rows[rowIndex13]

	for _, cell := range row13.Cells {
		jalan = append(jalan, models.Jalan{Permukaan: cell.String()})
	}

	rowIndex14 := 14
	row14 := sheet.Rows[rowIndex14]

	for _, cell := range row14.Cells {
		jalan = append(jalan, models.Jalan{Fungsi: cell.String()})
	}

	initializers.DB.Create(&models.DataMaster{
		JenisProyek:             jenis_proyek,
		LokasiPengambilan:       lokasi,
		JenisRencanaPembangunan: jenis_kegiatan,
		RencanaPembangunan:      rencana,
		Persyaratan:             persyaratan,
		KategoriPerlengkapan:    ketegori_perlengkapan,
		PerlengkapanLaluLintas:  perlengkapan,
		Provinsi:                provinsi,
		Kabupaten:               Kabupaten,
		Kecamatan:               kecamatan,
		Kelurahan:               kelurahan,
		Jalan:                   jalan,
		UpdatedAt:               now + " " + time.Now().In(loc).Format("15:04:05"),
	})

	fmt.Println("Migration complete")
}
