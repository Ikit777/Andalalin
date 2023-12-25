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

	jenis_kegiatan := []string{"Pusat kegiatan", "Pemukiman", "Infrastruktur", "Lainnya"}

	pusat_kegiatan := []models.JenisRencana{}
	pusat_kegiatan = append(pusat_kegiatan, models.JenisRencana{Jenis: "Kegiatan perdagangan dan perbelanjaan", Kriteria: "Luas lantai bangunan", Satuan: "m²", Terbilang: "meter persegi"})
	pusat_kegiatan = append(pusat_kegiatan, models.JenisRencana{Jenis: "Kegiatan perkantoran", Kriteria: "Luas lantai bangunan", Satuan: "m²", Terbilang: "meter persegi"})
	pusat_kegiatan = append(pusat_kegiatan, models.JenisRencana{Jenis: "Kegiatan industri", Kriteria: "Luas lantai bangunan", Satuan: "m²", Terbilang: "meter persegi"})
	pusat_kegiatan = append(pusat_kegiatan, models.JenisRencana{Jenis: "Kegiatan pergudangan", Kriteria: "Luas lantai bangunan", Satuan: "m²", Terbilang: "meter persegi"})
	pusat_kegiatan = append(pusat_kegiatan, models.JenisRencana{Jenis: "Kawasan pariwisata", Kriteria: "", Satuan: "", Terbilang: ""})
	pusat_kegiatan = append(pusat_kegiatan, models.JenisRencana{Jenis: "Tempat wisata", Kriteria: "Luas lahan", Satuan: "hektar", Terbilang: "hektar"})
	pusat_kegiatan = append(pusat_kegiatan, models.JenisRencana{Jenis: "Sekolah atau universitas", Kriteria: "Jumlah siswa", Satuan: "siswa", Terbilang: "siswa"})
	pusat_kegiatan = append(pusat_kegiatan, models.JenisRencana{Jenis: "Rumah sakit", Kriteria: "Jumlah tempat tidur", Satuan: "tempat tidur", Terbilang: "tempat tidur"})
	pusat_kegiatan = append(pusat_kegiatan, models.JenisRencana{Jenis: "Bank", Kriteria: "Luas lantai bangunan", Satuan: "m²", Terbilang: "meter persegi"})

	infrastruktur := []models.JenisRencana{}
	infrastruktur = append(infrastruktur, models.JenisRencana{Jenis: "Akses ke dan dari jalan tol", Kriteria: "", Satuan: "", Terbilang: ""})
	infrastruktur = append(infrastruktur, models.JenisRencana{Jenis: "Pelabuhan utama", Kriteria: "", Satuan: "", Terbilang: ""})
	infrastruktur = append(infrastruktur, models.JenisRencana{Jenis: "Pelabuhan pengumpul", Kriteria: "", Satuan: "", Terbilang: ""})
	infrastruktur = append(infrastruktur, models.JenisRencana{Jenis: "Pelabuhan pengumpan regional", Kriteria: "", Satuan: "", Terbilang: ""})
	infrastruktur = append(infrastruktur, models.JenisRencana{Jenis: "Pelabuhan pengumpan lokal", Kriteria: "", Satuan: "", Terbilang: ""})
	infrastruktur = append(infrastruktur, models.JenisRencana{Jenis: "Pelabuhan pengumpan khusus", Kriteria: "Luas lahan", Satuan: "m²", Terbilang: "meter persegi"})
	infrastruktur = append(infrastruktur, models.JenisRencana{Jenis: "Pelabuhan sungai, danau dan penyebrangan", Kriteria: "", Satuan: "", Terbilang: ""})
	infrastruktur = append(infrastruktur, models.JenisRencana{Jenis: "Bandar udara pengumpul skala pelayanan primer", Kriteria: "Jumlah pengguna pertahun", Satuan: "orang", Terbilang: "orang"})
	infrastruktur = append(infrastruktur, models.JenisRencana{Jenis: "Bandar udara pengumpul skala pelayanan sekunder", Kriteria: "Jumlah pengguna pertahun", Satuan: "orang", Terbilang: "orang"})
	infrastruktur = append(infrastruktur, models.JenisRencana{Jenis: "Bandar udara pengumpul skala pelayanan tersier", Kriteria: "Jumlah pengguna pertahun", Satuan: "orang", Terbilang: "orang"})
	infrastruktur = append(infrastruktur, models.JenisRencana{Jenis: "Bandar udara pengumpan (spoke)", Kriteria: "", Satuan: "", Terbilang: ""})
	infrastruktur = append(infrastruktur, models.JenisRencana{Jenis: "Terminal penumpang tipe A", Kriteria: "", Satuan: "", Terbilang: ""})
	infrastruktur = append(infrastruktur, models.JenisRencana{Jenis: "Terminal penumpang tipe B", Kriteria: "", Satuan: "", Terbilang: ""})
	infrastruktur = append(infrastruktur, models.JenisRencana{Jenis: "Terminal penumpang tipe C", Kriteria: "", Satuan: "", Terbilang: ""})
	infrastruktur = append(infrastruktur, models.JenisRencana{Jenis: "Terminal angkutan barang", Kriteria: "", Satuan: "", Terbilang: ""})
	infrastruktur = append(infrastruktur, models.JenisRencana{Jenis: "Terminal peti kemas", Kriteria: "", Satuan: "", Terbilang: ""})
	infrastruktur = append(infrastruktur, models.JenisRencana{Jenis: "Stasiun kereta api kelar besar", Kriteria: "", Satuan: "", Terbilang: ""})
	infrastruktur = append(infrastruktur, models.JenisRencana{Jenis: "Stasiun kereta api kelar sedang", Kriteria: "", Satuan: "", Terbilang: ""})
	infrastruktur = append(infrastruktur, models.JenisRencana{Jenis: "Stasiun kereta api kelar kecil", Kriteria: "", Satuan: "", Terbilang: ""})
	infrastruktur = append(infrastruktur, models.JenisRencana{Jenis: "Pool kendaraan", Kriteria: "", Satuan: "", Terbilang: ""})
	infrastruktur = append(infrastruktur, models.JenisRencana{Jenis: "Fasilitas parkir umum", Kriteria: "Besar satuan ruang parkir", Satuan: "SRP", Terbilang: "satuan ruang parkir"})

	pemukiman := []models.JenisRencana{}
	pemukiman = append(pemukiman, models.JenisRencana{Jenis: "Perumahan sederhana", Kriteria: "Jumlah unit", Satuan: "unit", Terbilang: "unit"})
	pemukiman = append(pemukiman, models.JenisRencana{Jenis: "Perumahan menengan-atas atau townhouse atau cluster", Kriteria: "Jumlah unit", Satuan: "unit", Terbilang: "unit"})
	pemukiman = append(pemukiman, models.JenisRencana{Jenis: "Rumah susun sederhana", Kriteria: "Jumlah unit", Satuan: "unit", Terbilang: "unit"})
	pemukiman = append(pemukiman, models.JenisRencana{Jenis: "Apartemen", Kriteria: "Jumlah unit", Satuan: "unit", Terbilang: "unit"})

	lainnya := []models.JenisRencana{}

	lainnya = append(lainnya, models.JenisRencana{Jenis: "Stasiun pengisin bahan bakar", Kriteria: "Jumlah dispenser", Satuan: "dispenser", Terbilang: "dispenser"})
	lainnya = append(lainnya, models.JenisRencana{Jenis: "Hotel", Kriteria: "Jumlah kamar", Satuan: "kamar", Terbilang: "kamar"})
	lainnya = append(lainnya, models.JenisRencana{Jenis: "Gedung pertemuan", Kriteria: "Luas lantai bangunan", Satuan: "m²", Terbilang: "meter persegi"})
	lainnya = append(lainnya, models.JenisRencana{Jenis: "Restaurant", Kriteria: "Jumlah tempat duduk", Satuan: "tempat duduk", Terbilang: "tempat duduk"})
	lainnya = append(lainnya, models.JenisRencana{Jenis: "Fasilitan olah raga", Kriteria: "Jumlah kapasitas penonton", Satuan: "orang", Terbilang: "orang"})
	lainnya = append(lainnya, models.JenisRencana{Jenis: "Kawasan TOD (indoor atau outdoor)", Kriteria: "Luas lantai bangunan", Satuan: "m²", Terbilang: "meter persegi"})
	lainnya = append(lainnya, models.JenisRencana{Jenis: "Asrama", Kriteria: "Jumlah unit", Satuan: "unit", Terbilang: "unit"})
	lainnya = append(lainnya, models.JenisRencana{Jenis: "Ruko", Kriteria: "Luas lahan keseluruhan", Satuan: "m²", Terbilang: "meter persegi"})
	lainnya = append(lainnya, models.JenisRencana{Jenis: "Jalan layang (flyover)", Kriteria: "", Satuan: "", Terbilang: ""})
	lainnya = append(lainnya, models.JenisRencana{Jenis: "Lintas bawas (underpass)", Kriteria: "", Satuan: "", Terbilang: ""})
	lainnya = append(lainnya, models.JenisRencana{Jenis: "Terowongan (tunnel)", Kriteria: "", Satuan: "", Terbilang: ""})
	lainnya = append(lainnya, models.JenisRencana{Jenis: "Jembatan", Kriteria: "", Satuan: "", Terbilang: ""})
	lainnya = append(lainnya, models.JenisRencana{Jenis: "Rest area tipe A", Kriteria: "", Satuan: "", Terbilang: ""})
	lainnya = append(lainnya, models.JenisRencana{Jenis: "Rest area tipe B", Kriteria: "", Satuan: "", Terbilang: ""})
	lainnya = append(lainnya, models.JenisRencana{Jenis: "Rest area tipe C", Kriteria: "", Satuan: "", Terbilang: ""})
	lainnya = append(lainnya, models.JenisRencana{Jenis: "Kegiatan yang menimbukan kepadatan 1500 kendaraan", Kriteria: "", Satuan: "", Terbilang: ""})
	lainnya = append(lainnya, models.JenisRencana{Jenis: "Kegiatan yang menimbukan kepadatan 500 kendaraan", Kriteria: "", Satuan: "", Terbilang: ""})
	lainnya = append(lainnya, models.JenisRencana{Jenis: "Kegiatan yang menimbukan kepadatan 100 kendaraan", Kriteria: "", Satuan: "", Terbilang: ""})

	rencana := []models.Rencana{}
	rencana = append(rencana, models.Rencana{Kategori: "Pusat kegiatan", JenisRencana: pusat_kegiatan})
	rencana = append(rencana, models.Rencana{Kategori: "Pemukiman", JenisRencana: pemukiman})
	rencana = append(rencana, models.Rencana{Kategori: "Infrastruktur", JenisRencana: infrastruktur})
	rencana = append(rencana, models.Rencana{Kategori: "Lainnya", JenisRencana: lainnya})

	ketegori_perlengkapan := []string{"Rambu peringatan", "Rambu larangan", "Rambu perintah", "Rambu petunjuk"}

	andalalin := []models.PersyaratanAndalalinInput{}
	andalalin = append(andalalin, models.PersyaratanAndalalinInput{Kebutuhan: "Wajib", Bangkitan: "Bangkitan rendah", Persyaratan: "Surat permohonan persetujuan andalalin", KeteranganPersyaratan: "Surat permohonan persetujuan analisis dampak lalu lintas adalah surat yang digunakan untuk mengajukan permohonan kepada pihak yang berwenang, biasanya pemerintah daerah, untuk mendapatkan persetujuan atau izin terkait dengan rencana atau proyek tertentu. Surat ini harus memuat informasi yang lengkap dan jelas mengenai rencana atau proyek yang diajukan, termasuk tujuan, dampak lingkungan, serta segala persyaratan yang harus dipenuhi."})
	andalalin = append(andalalin, models.PersyaratanAndalalinInput{Kebutuhan: "Wajib", Bangkitan: "Bangkitan rendah", Persyaratan: "Identitas pemohon atau penanggung jawab", KeteranganPersyaratan: "Identitas pemohon atau penanggung jawab terdiri dari dokumen identitas dan NPWP. dokumen identitas dapat berupa kartu tanda penduduk atau paspor atau akta pendirian badan usaha. Dokumen identitas bertujuan untuk mengindentifikasi pemohon atau penanggung jawab terhadap permohonan.\nSedangkan Nomor Pokok Wajib Pajak (NPWP) adalah nomor identifikasi pajak yang diberikan kepada warga negara atau entitas yang wajib membayar pajak di Indonesia. NPWP adalah bagian penting dari administrasi pajak di Indonesia dan dikeluarkan oleh Direktorat Jenderal Pajak (DJP) yang merupakan lembaga di bawah Kementerian Keuangan Republik Indonesia. NPWP digunakan untuk mengidentifikasi wajib pajak dan melacak pembayaran pajak yang dikenakan."})
	andalalin = append(andalalin, models.PersyaratanAndalalinInput{Kebutuhan: "Wajib", Bangkitan: "Bangkitan rendah", Persyaratan: "Surat bukti kepemilikan atau Penguasaan lahan", KeteranganPersyaratan: "Surat bukti kepemilikan atau penguasaan lahan adalah dokumen yang dimaksudkan untuk menunjukkan hak hukum seseorang atau badan usaha atas sebidang lahan atau properti tertentu."})
	andalalin = append(andalalin, models.PersyaratanAndalalinInput{Kebutuhan: "Wajib", Bangkitan: "Bangkitan rendah", Persyaratan: "Surat kesesuaian tata ruang dan atau izin pemanfaatan ruang", KeteranganPersyaratan: "Surat kesesuaian tata ruang dan/atau Izin pemanfaatan ruang adalah dokumen yang diterbitkan oleh pihak berwenang, seperti pemerintah daerah atau instansi terkait, untuk memberikan izin atau persetujuan terkait dengan penggunaan dan pengembangan lahan atau ruang tertentu."})
	andalalin = append(andalalin, models.PersyaratanAndalalinInput{Kebutuhan: "Wajib", Bangkitan: "Bangkitan rendah", Persyaratan: "Gambar tata letak bangunan (site plan) dan DED Bangunan yang diusulkan", KeteranganPersyaratan: "Tata letak bangunan, atau dalam bahasa Inggris dikenal sebagai site plan, adalah gambar atau diagram yang menunjukkan tata letak fisik bangunan, struktur, dan fasilitas terkait dalam suatu area tertentu. Site plan biasanya digunakan dalam perencanaan konstruksi, pengembangan properti, perizinan bangunan, dan perencanaan tata ruang.\nDokumen Engineering Design (DED) adalah bagian penting dari proses perencanaan dan konstruksi bangunan. DED untuk bangunan yang diusulkan adalah dokumen yang merinci dan mendokumentasikan desain teknik dan teknis dari bangunan yang akan dibangun. DED adalah langkah selanjutnya setelah tahap perencanaan dan perancangan awal, dan sebelum memulai proses konstruksi."})
	andalalin = append(andalalin, models.PersyaratanAndalalinInput{Kebutuhan: "Wajib", Bangkitan: "Bangkitan rendah", Persyaratan: "Foto kondisi eksisting lapangan terkini", KeteranganPersyaratan: "Foto kondisi eksisting lapangan yang terkini adalah gambar-gambar yang menggambarkan kondisi aktual dan terbaru dari lapangan atau area tertentu pada suatu waktu. Foto-foto ini sangat berguna dalam berbagai konteks, termasuk dalam proyek konstruksi, pengembangan properti, pemantauan lingkungan, dan penelitian."})
	andalalin = append(andalalin, models.PersyaratanAndalalinInput{Kebutuhan: "Tidak wajib", Bangkitan: "Bangkitan rendah", Persyaratan: "MOU Kerjsa sama", KeteranganPersyaratan: "Apabila ada kerja sama dengan pihak lain, semisal perjanjian sewa lahan, perjanjian penggunaan akses dsb."})

	andalalin = append(andalalin, models.PersyaratanAndalalinInput{Kebutuhan: "Wajib", Bangkitan: "Bangkitan sedang", Persyaratan: "Surat permohonan persetujuan andalalin", KeteranganPersyaratan: "Surat permohonan persetujuan analisis dampak lalu lintas adalah surat yang digunakan untuk mengajukan permohonan kepada pihak yang berwenang, biasanya pemerintah daerah, untuk mendapatkan persetujuan atau izin terkait dengan rencana atau proyek tertentu. Surat ini harus memuat informasi yang lengkap dan jelas mengenai rencana atau proyek yang diajukan, termasuk tujuan, dampak lingkungan, serta segala persyaratan yang harus dipenuhi."})
	andalalin = append(andalalin, models.PersyaratanAndalalinInput{Kebutuhan: "Wajib", Bangkitan: "Bangkitan sedang", Persyaratan: "Dokumen hasil analisis dampak lalu lintas", KeteranganPersyaratan: "Analisis dampak lalu lintas adalah proses evaluasi yang digunakan untuk memahami dan memprediksi bagaimana suatu proyek, kebijakan, atau perubahan akan memengaruhi lalu lintas jalan raya. Hasil analisis dampak lalu lintas dapat mencakup berbagai informasi yang diperlukan untuk membuat keputusan tentang pembangunan infrastruktur, rencana lalu lintas, atau perubahan aturan lalu lintas."})
	andalalin = append(andalalin, models.PersyaratanAndalalinInput{Kebutuhan: "Wajib", Bangkitan: "Bangkitan sedang", Persyaratan: "Identitas pemohon atau penanggung jawab", KeteranganPersyaratan: "Identitas pemohon atau penanggung jawab terdiri dari dokumen identitas dan NPWP. dokumen identitas dapat berupa kartu tanda penduduk atau paspor atau akta pendirian badan usaha. Dokumen identitas bertujuan untuk mengindentifikasi pemohon atau penanggung jawab terhadap permohonan.\nSedangkan Nomor Pokok Wajib Pajak (NPWP) adalah nomor identifikasi pajak yang diberikan kepada warga negara atau entitas yang wajib membayar pajak di Indonesia. NPWP adalah bagian penting dari administrasi pajak di Indonesia dan dikeluarkan oleh Direktorat Jenderal Pajak (DJP) yang merupakan lembaga di bawah Kementerian Keuangan Republik Indonesia. NPWP digunakan untuk mengidentifikasi wajib pajak dan melacak pembayaran pajak yang dikenakan."})
	andalalin = append(andalalin, models.PersyaratanAndalalinInput{Kebutuhan: "Wajib", Bangkitan: "Bangkitan sedang", Persyaratan: "Sertifikat konsultan atau tenaga ahli penyusun andalalin sesuai klasifikasi", KeteranganPersyaratan: "Sertifikat konsultan atau tenaga ahli penyusun andalalin sesuai klasifikasi adalah dokumen yang menunjukkan bahwa seorang profesional atau konsultan memiliki kualifikasi, kompetensi, dan sertifikasi yang sesuai untuk menyusun dokumen Andalalin. Andalalin adalah dokumen yang berisi analisis dampak lingkungan (AMDAL) yang diperlukan dalam perencanaan dan pelaksanaan proyek-proyek yang memiliki potensi dampak terhadap lingkungan."})
	andalalin = append(andalalin, models.PersyaratanAndalalinInput{Kebutuhan: "Wajib", Bangkitan: "Bangkitan sedang", Persyaratan: "Surat bukti kepemilikan atau penguasaan lahan", KeteranganPersyaratan: "Surat bukti kepemilikan atau penguasaan lahan adalah dokumen yang dimaksudkan untuk menunjukkan hak hukum seseorang atau badan usaha atas sebidang lahan atau properti tertentu."})
	andalalin = append(andalalin, models.PersyaratanAndalalinInput{Kebutuhan: "Wajib", Bangkitan: "Bangkitan sedang", Persyaratan: "Surat kesesuaian tata ruang dan atau izin pemanfaatan ruang", KeteranganPersyaratan: "Surat kesesuaian tata ruang dan/atau Izin pemanfaatan ruang adalah dokumen yang diterbitkan oleh pihak berwenang, seperti pemerintah daerah atau instansi terkait, untuk memberikan izin atau persetujuan terkait dengan penggunaan dan pengembangan lahan atau ruang tertentu."})
	andalalin = append(andalalin, models.PersyaratanAndalalinInput{Kebutuhan: "Wajib", Bangkitan: "Bangkitan sedang", Persyaratan: "Gambar tata letak bangunan (site plan) dan DED Bangunan yang diusulkan", KeteranganPersyaratan: "Tata letak bangunan, atau dalam bahasa Inggris dikenal sebagai site plan, adalah gambar atau diagram yang menunjukkan tata letak fisik bangunan, struktur, dan fasilitas terkait dalam suatu area tertentu. Site plan biasanya digunakan dalam perencanaan konstruksi, pengembangan properti, perizinan bangunan, dan perencanaan tata ruang.\nDokumen Engineering Design (DED) adalah bagian penting dari proses perencanaan dan konstruksi bangunan. DED untuk bangunan yang diusulkan adalah dokumen yang merinci dan mendokumentasikan desain teknik dan teknis dari bangunan yang akan dibangun. DED adalah langkah selanjutnya setelah tahap perencanaan dan perancangan awal, dan sebelum memulai proses konstruksi."})
	andalalin = append(andalalin, models.PersyaratanAndalalinInput{Kebutuhan: "Wajib", Bangkitan: "Bangkitan sedang", Persyaratan: "Foto kondisi eksisting lapangan terkini", KeteranganPersyaratan: "Foto kondisi eksisting lapangan yang terkini adalah gambar-gambar yang menggambarkan kondisi aktual dan terbaru dari lapangan atau area tertentu pada suatu waktu. Foto-foto ini sangat berguna dalam berbagai konteks, termasuk dalam proyek konstruksi, pengembangan properti, pemantauan lingkungan, dan penelitian."})
	andalalin = append(andalalin, models.PersyaratanAndalalinInput{Kebutuhan: "Tidak wajib", Bangkitan: "Bangkitan sedang", Persyaratan: "MOU Kerjsa sama", KeteranganPersyaratan: "Apabila ada kerja sama dengan pihak lain, semisal perjanjian sewa lahan, perjanjian penggunaan akses dsb."})

	andalalin = append(andalalin, models.PersyaratanAndalalinInput{Kebutuhan: "Wajib", Bangkitan: "Bangkitan tinggi", Persyaratan: "Surat permohonan persetujuan andalalin", KeteranganPersyaratan: "Surat permohonan persetujuan analisis dampak lalu lintas adalah surat yang digunakan untuk mengajukan permohonan kepada pihak yang berwenang, biasanya pemerintah daerah, untuk mendapatkan persetujuan atau izin terkait dengan rencana atau proyek tertentu. Surat ini harus memuat informasi yang lengkap dan jelas mengenai rencana atau proyek yang diajukan, termasuk tujuan, dampak lingkungan, serta segala persyaratan yang harus dipenuhi."})
	andalalin = append(andalalin, models.PersyaratanAndalalinInput{Kebutuhan: "Wajib", Bangkitan: "Bangkitan tinggi", Persyaratan: "Dokumen hasil analisis dampak lalu lintas", KeteranganPersyaratan: "Analisis dampak lalu lintas adalah proses evaluasi yang digunakan untuk memahami dan memprediksi bagaimana suatu proyek, kebijakan, atau perubahan akan memengaruhi lalu lintas jalan raya. Hasil analisis dampak lalu lintas dapat mencakup berbagai informasi yang diperlukan untuk membuat keputusan tentang pembangunan infrastruktur, rencana lalu lintas, atau perubahan aturan lalu lintas."})
	andalalin = append(andalalin, models.PersyaratanAndalalinInput{Kebutuhan: "Wajib", Bangkitan: "Bangkitan tinggi", Persyaratan: "Identitas pemohon atau penanggung jawab", KeteranganPersyaratan: "Identitas pemohon atau penanggung jawab terdiri dari dokumen identitas dan NPWP. dokumen identitas dapat berupa kartu tanda penduduk atau paspor atau akta pendirian badan usaha. Dokumen identitas bertujuan untuk mengindentifikasi pemohon atau penanggung jawab terhadap permohonan.\nSedangkan Nomor Pokok Wajib Pajak (NPWP) adalah nomor identifikasi pajak yang diberikan kepada warga negara atau entitas yang wajib membayar pajak di Indonesia. NPWP adalah bagian penting dari administrasi pajak di Indonesia dan dikeluarkan oleh Direktorat Jenderal Pajak (DJP) yang merupakan lembaga di bawah Kementerian Keuangan Republik Indonesia. NPWP digunakan untuk mengidentifikasi wajib pajak dan melacak pembayaran pajak yang dikenakan."})
	andalalin = append(andalalin, models.PersyaratanAndalalinInput{Kebutuhan: "Wajib", Bangkitan: "Bangkitan tinggi", Persyaratan: "Sertifikat konsultan atau tenaga ahli penyusun andalalin sesuai klasifikasi", KeteranganPersyaratan: "Sertifikat konsultan atau tenaga ahli penyusun andalalin sesuai klasifikasi adalah dokumen yang menunjukkan bahwa seorang profesional atau konsultan memiliki kualifikasi, kompetensi, dan sertifikasi yang sesuai untuk menyusun dokumen Andalalin. Andalalin adalah dokumen yang berisi analisis dampak lingkungan (AMDAL) yang diperlukan dalam perencanaan dan pelaksanaan proyek-proyek yang memiliki potensi dampak terhadap lingkungan."})
	andalalin = append(andalalin, models.PersyaratanAndalalinInput{Kebutuhan: "Wajib", Bangkitan: "Bangkitan tinggi", Persyaratan: "Surat bukti kepemilikan atau penguasaan lahan", KeteranganPersyaratan: "Surat bukti kepemilikan atau penguasaan lahan adalah dokumen yang dimaksudkan untuk menunjukkan hak hukum seseorang atau badan usaha atas sebidang lahan atau properti tertentu."})
	andalalin = append(andalalin, models.PersyaratanAndalalinInput{Kebutuhan: "Wajib", Bangkitan: "Bangkitan tinggi", Persyaratan: "Surat kesesuaian tata ruang dan atau izin pemanfaatan ruang", KeteranganPersyaratan: "Surat kesesuaian tata ruang dan/atau Izin pemanfaatan ruang adalah dokumen yang diterbitkan oleh pihak berwenang, seperti pemerintah daerah atau instansi terkait, untuk memberikan izin atau persetujuan terkait dengan penggunaan dan pengembangan lahan atau ruang tertentu."})
	andalalin = append(andalalin, models.PersyaratanAndalalinInput{Kebutuhan: "Wajib", Bangkitan: "Bangkitan tinggi", Persyaratan: "Gambar tata letak bangunan (site plan) dan DED Bangunan yang diusulkan", KeteranganPersyaratan: "Tata letak bangunan, atau dalam bahasa Inggris dikenal sebagai site plan, adalah gambar atau diagram yang menunjukkan tata letak fisik bangunan, struktur, dan fasilitas terkait dalam suatu area tertentu. Site plan biasanya digunakan dalam perencanaan konstruksi, pengembangan properti, perizinan bangunan, dan perencanaan tata ruang.\nDokumen Engineering Design (DED) adalah bagian penting dari proses perencanaan dan konstruksi bangunan. DED untuk bangunan yang diusulkan adalah dokumen yang merinci dan mendokumentasikan desain teknik dan teknis dari bangunan yang akan dibangun. DED adalah langkah selanjutnya setelah tahap perencanaan dan perancangan awal, dan sebelum memulai proses konstruksi."})
	andalalin = append(andalalin, models.PersyaratanAndalalinInput{Kebutuhan: "Wajib", Bangkitan: "Bangkitan tinggi", Persyaratan: "Foto kondisi eksisting lapangan terkini", KeteranganPersyaratan: "Foto kondisi eksisting lapangan yang terkini adalah gambar-gambar yang menggambarkan kondisi aktual dan terbaru dari lapangan atau area tertentu pada suatu waktu. Foto-foto ini sangat berguna dalam berbagai konteks, termasuk dalam proyek konstruksi, pengembangan properti, pemantauan lingkungan, dan penelitian."})
	andalalin = append(andalalin, models.PersyaratanAndalalinInput{Kebutuhan: "Tidak wajib", Bangkitan: "Bangkitan tinggi", Persyaratan: "MOU Kerjsa sama", KeteranganPersyaratan: "Apabila ada kerja sama dengan pihak lain, semisal perjanjian sewa lahan, perjanjian penggunaan akses dsb."})

	perlalin := []models.PersyaratanPerlalinInput{}
	perlalin = append(perlalin, models.PersyaratanPerlalinInput{Kebutuhan: "Wajib", Persyaratan: "Kartu tanda penduduk", KeteranganPersyaratan: "Kartu tanda penduduk adalah dokumen identitas yang bertujuan untuk mengindentifikasi pemohon atau penanggung jawab terhadap permohonan."})
	perlalin = append(perlalin, models.PersyaratanPerlalinInput{Kebutuhan: "Wajib", Persyaratan: "Surat permohonan", KeteranganPersyaratan: "Surat permohonan berisi infromasi terkait permohonan yang akan diajukan seperti perlengkapan lalu lintas"})

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

	fileJalan, err := os.Open("assets/Jalan/jalan.csv")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer fileJalan.Close()

	csvJalan := csv.NewReader(fileJalan)
	csvJalan.Comma = ','

	jalan := []models.Jalan{}

	for {
		record, err := csvJalan.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}

		jalan = append(jalan, models.Jalan{KodeProvinsi: record[0], KodeKabupaten: record[1], KodeKecamatan: record[2], KodeKelurahan: record[3], KodeJalan: record[4], Nama: record[5], Pangkal: record[6], Ujung: record[7], Kelurahan: record[8], Kecamatan: record[9], Panjang: record[10], Lebar: record[11], Permukaan: record[12], Fungsi: record[13]})
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
