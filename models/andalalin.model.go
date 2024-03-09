package models

import (
	"github.com/google/uuid"
)

type Andalalin struct {
	//Data Permohonan
	IdAndalalin       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	IdUser            uuid.UUID `gorm:"type:varchar(255);not null"`
	JenisAndalalin    string    `gorm:"type:varchar(255);not null"`
	WaktuAndalalin    string    `gorm:"not null"`
	TanggalAndalalin  string    `gorm:"not null"`
	StatusAndalalin   string    `gorm:"type:varchar(255);not null"`
	Bangkitan         string    `gorm:"type:varchar(255);not null"`
	Pemohon           string    `gorm:"type:varchar(255);not null"`
	Kategori          string    `gorm:"type:varchar(255);not null"`
	Jenis             string    `gorm:"type:varchar(255);not null"`
	Kode              string    `gorm:"type:varchar(255);not null"`
	LokasiPengambilan string    `gorm:"type:varchar(255);not null"`

	//Data pemohon
	NikPemohon          string  `gorm:"type:varchar(255);not null"`
	NamaPemohon         string  `gorm:"type:varchar(255);not null"`
	EmailPemohon        string  `gorm:"type:varchar(255);not null"`
	TempatLahirPemohon  string  `gorm:"type:varchar(255);not null"`
	TanggalLahirPemohon string  `gorm:"type:varchar(255);not null"`
	NegaraPemohon       string  `gorm:"type:varchar(255);not null"`
	ProvinsiPemohon     string  `gorm:"type:varchar(255);not null"`
	KabupatenPemohon    string  `gorm:"type:varchar(255);not null"`
	KecamatanPemohon    string  `gorm:"type:varchar(255);not null"`
	KelurahanPemohon    string  `gorm:"type:varchar(255);not null"`
	AlamatPemohon       string  `gorm:"type:varchar(255);not null"`
	JenisKelaminPemohon string  `sql:"type:enum('Laki-laki', 'Perempuan');not null"`
	NomerPemohon        string  `gorm:"type:varchar(255);not null"`
	JabatanPemohon      *string `gorm:"type:varchar(255);"`

	//Data Proyek
	NamaProyek      string `gorm:"type:varchar(255);not null"`
	JenisProyek     string `gorm:"type:varchar(255);not null"`
	NegaraProyek    string `gorm:"type:varchar(255);not null"`
	ProvinsiProyek  string `gorm:"type:varchar(255);not null"`
	KabupatenProyek string `gorm:"type:varchar(255);not null"`
	KecamatanProyek string `gorm:"type:varchar(255);not null"`
	KelurahanProyek string `gorm:"type:varchar(255);not null"`
	AlamatProyek    string `gorm:"type:varchar(255);not null"`
	KodeJalan       string `gorm:"type:varchar(255);not null"`
	KodeJalanMerge  string `gorm:"type:varchar(255);not null"`
	NamaJalan       string `gorm:"type:varchar(255);not null"`
	PangkalJalan    string `gorm:"type:varchar(255);not null"`
	UjungJalan      string `gorm:"type:varchar(255);not null"`
	PanjangJalan    string `gorm:"type:varchar(255);not null"`
	LebarJalan      string `gorm:"type:varchar(255);not null"`
	PermukaanJalan  string `gorm:"type:varchar(255);not null"`
	FungsiJalan     string `gorm:"type:varchar(255);not null"`
	StatusJalan     string `gorm:"type:varchar(255);not null"`

	//Data Perusahaan
	NamaPerusahaan              *string `gorm:"type:varchar(255);"`
	NegaraPerusahaan            string  `gorm:"type:varchar(255);"`
	ProvinsiPerusahaan          *string `gorm:"type:varchar(255);"`
	KabupatenPerusahaan         *string `gorm:"type:varchar(255);"`
	KecamatanPerusahaan         *string `gorm:"type:varchar(255);"`
	KelurahanPerusahaan         *string `gorm:"type:varchar(255);"`
	AlamatPerusahaan            *string `gorm:"type:varchar(255);"`
	NomerPerusahaan             *string `gorm:"type:varchar(255);"`
	EmailPerusahaan             *string `gorm:"type:varchar(255);"`
	NamaPimpinan                *string `gorm:"type:varchar(255);"`
	JabatanPimpinan             *string `gorm:"type:varchar(255);"`
	JenisKelaminPimpinan        *string `sql:"type:enum('Laki-laki', 'Perempuan');"`
	NegaraPimpinanPerusahaan    string  `gorm:"type:varchar(255);"`
	ProvinsiPimpinanPerusahaan  *string `gorm:"type:varchar(255);"`
	KabupatenPimpinanPerusahaan *string `gorm:"type:varchar(255);"`
	KecamatanPimpinanPerusahaan *string `gorm:"type:varchar(255);"`
	KelurahanPimpinanPerusahaan *string `gorm:"type:varchar(255);"`
	AlamatPimpinan              *string `gorm:"type:varchar(255);"`

	//Data Konsultan
	NamaKonsultan                  *string `gorm:"type:varchar(255);"`
	NegaraKonsultan                string  `gorm:"type:varchar(255);"`
	ProvinsiKonsultan              *string `gorm:"type:varchar(255);"`
	KabupatenKonsultan             *string `gorm:"type:varchar(255);"`
	KecamatanKonsultan             *string `gorm:"type:varchar(255);"`
	KelurahanKonsultan             *string `gorm:"type:varchar(255);"`
	AlamatKonsultan                *string `gorm:"type:varchar(255);"`
	NomerKonsultan                 *string `gorm:"type:varchar(255);"`
	EmailKonsultan                 *string `gorm:"type:varchar(255);"`
	NamaPenyusunDokumen            *string `gorm:"type:varchar(255);"`
	JenisKelaminPenyusunDokumen    *string `sql:"type:enum('Laki-laki', 'Perempuan');"`
	NegaraPenyusunDokumen          string  `gorm:"type:varchar(255);"`
	ProvinsiPenyusunDokumen        *string `gorm:"type:varchar(255);"`
	KabupatenPenyusunDokumen       *string `gorm:"type:varchar(255);"`
	KecamatanPenyusunDokumen       *string `gorm:"type:varchar(255);"`
	KelurahanPenyusunDokumen       *string `gorm:"type:varchar(255);"`
	AlamatPenyusunDokumen          *string `gorm:"type:varchar(255);"`
	NomerSertifikatPenyusunDokumen *string `gorm:"type:varchar(255);"`
	KlasifikasiPenyusunDokumen     *string `gorm:"type:varchar(255);"`

	//Data Kegiatan
	Aktivitas         string  `gorm:"type:varchar(255);not null"`
	Peruntukan        string  `gorm:"type:varchar(255);not null"`
	TotalLuasLahan    string  `gorm:"type:varchar(255);not null"`
	KriteriaKhusus    *string `gorm:"type:varchar(255);"`
	NilaiKriteria     *string `gorm:"type:varchar(255);"`
	Terbilang         *string `gorm:"type:varchar(255);"`
	LokasiBangunan    string  `gorm:"type:varchar(255);not null"`
	LatitudeBangunan  float64
	LongitudeBangunan float64
	NomerSKRK         string `gorm:"type:varchar(255);not null"`
	TanggalSKRK       string `gorm:"type:varchar(255);not null"`
	Catatan           *string

	//Berkas Permohonan
	BerkasPermohonan []BerkasPermohonan `gorm:"serializer:json"`

	//Status berkas permohonan (Baru atau Revisi)
	StatusBerkasPermohonan string `gorm:"type:varchar(255);not null"`

	//Data catatan asistensi dokumen andalalin
	HasilAsistensiDokumen   string
	CatatanAsistensiDokumen []CatatanAsistensi `gorm:"serializer:json"`

	//Kelengkapan Tidak Sesuai
	KelengkapanTidakSesuai []KelengkapanTidakSesuai `gorm:"serializer:json"`

	//Petimbangan
	PertimbanganPenolakan string
	PertimbanganPenundaan string

	//Persyaratan tidak terpenuhi
	PersyaratanTidakSesuai []PersayaratanTidakSesuai `gorm:"serializer:json"`

	//Pemeriksaan surat persetujuan
	HasilPemeriksaan   string `gorm:"type:varchar(255);"`
	CatatanPemeriksaan *string

	//Data Surat Permohonan
	Nomor   string `gorm:"type:varchar(255);"`
	Tanggal string
}

type InputAndalalin struct {
	//Data Permohonan
	Bangkitan                       string `json:"kategori_bangkitan" binding:"required"`
	Pemohon                         string `json:"kategori_pemohon" binding:"required"`
	KategoriJenisRencanaPembangunan string `json:"kategori" binding:"required"`
	JenisRencanaPembangunan         string `json:"jenis_rencana_pembangunan" binding:"required"`
	LokasiPengambilan               string `json:"lokasi_pengambilan" binding:"required"`

	//Data Proyek
	NamaProyek      string `json:"nama_proyek" binding:"required"`
	ProvinsiProyek  string `json:"provinsi_proyek" binding:"required"`
	KabupatenProyek string `json:"kabupaten_proyek" binding:"required"`
	KecamatanProyek string `json:"kecamatan_proyek" binding:"required"`
	KelurahanProyek string `json:"kelurahan_proyek" binding:"required"`
	JenisProyek     string `json:"jenis_proyek" binding:"required"`
	AlamatProyek    string `json:"alamat_proyek" binding:"required"`
	KodeJalan       string `json:"kode_jalan" binding:"required"`
	KodeJalanMerge  string `json:"kode_jalan_merge" binding:"required"`
	NamaJalan       string `json:"nama_jalan" binding:"required"`
	PangkalJalan    string `json:"pangkal_jalan" binding:"required"`
	UjungJalan      string `json:"ujung_jalan" binding:"required"`
	PanjangJalan    string `json:"panjang_jalan" binding:"required"`
	LebarJalan      string `json:"lebar_jalan" binding:"required"`
	PermukaanJalan  string `json:"permukaan_jalan" binding:"required"`
	FungsiJalan     string `json:"fungsi_jalan" binding:"required"`
	StatusJalan     string `json:"status_jalan" binding:"required"`

	//Data Pemohon
	NikPemohon          string  `json:"nik_pemohon" binding:"required"`
	TempatLahirPemohon  string  `json:"tempat_lahir_pemohon" binding:"required"`
	TanggalLahirPemohon string  `json:"tanggal_lahir_pemohon" binding:"required"`
	AlamatPemohon       string  `json:"alamat_pemohon" binding:"required"`
	ProvinsiPemohon     string  `json:"provinsi_pemohon" binding:"required"`
	KabupatenPemohon    string  `json:"kabupaten_pemohon" binding:"required"`
	KecamatanPemohon    string  `json:"kecamatan_pemohon" binding:"required"`
	KelurahanPemohon    string  `json:"kelurahan_pemohon" binding:"required"`
	JenisKelaminPemohon string  `json:"jenis_kelamin_pemohon" binding:"required"`
	NomerPemohon        string  `json:"nomer_pemohon" binding:"required"`
	JabatanPemohon      *string `json:"jabatan_pemohon" binding:"required"`

	//Data Perusahaan
	NamaPerusahaan              *string `json:"nama_perusahaan" binding:"required"`
	AlamatPerusahaan            *string `json:"alamat_perusahaan" binding:"required"`
	ProvinsiPerusahaan          *string `json:"provinsi_perusahaan" binding:"required"`
	KabupatenPerusahaan         *string `json:"kabupaten_perusahaan" binding:"required"`
	KecamatanPerusahaan         *string `json:"kecamatan_perusahaan" binding:"required"`
	KelurahanPerusahaan         *string `json:"kelurahan_perusahaan" binding:"required"`
	NomerPerusahaan             *string `json:"nomer_perusahaan" binding:"required"`
	EmailPerusahaan             *string `json:"email_perusahaan" binding:"required"`
	NamaPimpinan                *string `json:"nama_pimpinan" binding:"required"`
	JabatanPimpinan             *string `json:"jabatan_pimpinan" binding:"required"`
	JenisKelaminPimpinan        *string `json:"jenis_kelamin_pimpinan" binding:"required"`
	ProvinsiPimpinanPerusahaan  *string `json:"provinsi_pimpinan_perusahaan" binding:"required"`
	KabupatenPimpinanPerusahaan *string `json:"kabupaten_pimpinan_perusahaan" binding:"required"`
	KecamatanPimpinanPerusahaan *string `json:"kecamatan_pimpinan_perusahaan" binding:"required"`
	KelurahanPimpinanPerusahaan *string `json:"kelurahan_pimpinan_perusahaan" binding:"required"`
	AlamatPimpinan              *string `json:"alamat_pimpinan_perusahaan" binding:"required"`

	//Data Konsultan
	NamaKonsultan                  *string `json:"nama_konsultan" binding:"required"`
	ProvinsiKonsultan              *string `json:"provinsi_konsultan" binding:"required"`
	KabupatenKonsultan             *string `json:"kabupaten_konsultan" binding:"required"`
	KecamatanKonsultan             *string `json:"kecamatan_konsultan" binding:"required"`
	KelurahanKonsultan             *string `json:"kelurahan_konsultan" binding:"required"`
	AlamatKonsultan                *string `json:"alamat_konsultan" binding:"required"`
	NomerKonsultan                 *string `json:"nomer_konsultan" binding:"required"`
	EmailKonsultan                 *string `json:"email_konsultan" binding:"required"`
	NamaPenyusunDokumen            *string `json:"nama_penyusun_dokumen" binding:"required"`
	JenisKelaminPenyusunDokumen    *string `json:"jenis_kelamin_penyusun_dokumen" binding:"required"`
	ProvinsiPenyusunDokumen        *string `json:"provinsi_penyusun_dokumen" binding:"required"`
	KabupatenPenyusunDokumen       *string `json:"kabupaten_penyusun_dokumen" binding:"required"`
	KecamatanPenyusunDokumen       *string `json:"kecamatan_penyusun_dokumen" binding:"required"`
	KelurahanPenyusunDokumen       *string `json:"kelurahan_penyusun_dokumen" binding:"required"`
	AlamatPenyusunDokumen          *string `json:"alamat_penyusun_dokumen" binding:"required"`
	NomerSertifikatPenyusunDokumen *string `json:"nomer_sertifikat_penyusun_dokumen" binding:"required"`
	KlasifikasiPenyusunDokumen     *string `json:"klasifikasi_penyusun_dokumen" binding:"required"`

	//Data Kegiatan
	Aktivitas         string  `json:"aktivitas" binding:"required"`
	Peruntukan        string  `json:"peruntukan" binding:"required"`
	TotalLuasLahan    string  `json:"total_luas" binding:"required"`
	KriteriaKhusus    *string `json:"kriteria_khusus" binding:"required"`
	NilaiKriteria     *string `json:"nilai_kriteria" binding:"required"`
	Terbilang         *string `json:"terbilang" binding:"required"`
	LokasiBangunan    string  `json:"lokasi_bangunan" binding:"required"`
	LatitudeBangunan  float64 `protobuf:"fixed64,1,opt,name=latitude,proto3" json:"latitude" binding:"required"`
	LongitudeBangunan float64 `protobuf:"fixed64,2,opt,name=longitude,proto3" json:"longtitude" binding:"required"`
	NomerSKRK         string  `json:"nomer_skrk" binding:"required"`
	TanggalSKRK       string  `json:"tanggal_skrk" binding:"required"`
	Catatan           *string `json:"catatan" binding:"required"`
}

type DataAndalalin struct {
	Andalalin InputAndalalin `form:"data"`
}

type AndalalinResponse struct {
	//Data Permohonan
	IdAndalalin       uuid.UUID `json:"id_andalalin,omitempty"`
	Bangkitan         string    `json:"kategori_bangkitan,omitempty"`
	Pemohon           string    `json:"kategori_pemohon,omitempty"`
	JenisAndalalin    string    `json:"jenis_andalalin,omitempty"`
	Kategori          string    `json:"kategori,omitempty"`
	Jenis             string    `json:"jenis_rencana_pembangunan,omitempty"`
	Kode              string    `json:"kode_andalalin,omitempty"`
	WaktuAndalalin    string    `json:"waktu_andalalin,omitempty"`
	TanggalAndalalin  string    `json:"tanggal_andalalin,omitempty"`
	StatusAndalalin   string    `json:"status_andalalin,omitempty"`
	LokasiPengambilan string    `json:"lokasi_pengambilan,omitempty"`

	//Data Proyek
	NamaProyek      string `json:"nama_proyek,omitempty"`
	JenisProyek     string `json:"jenis_proyek,omitempty"`
	NegaraProyek    string `json:"negara_proyek,omitempty"`
	ProvinsiProyek  string `json:"provinsi_proyek,omitempty"`
	KabupatenProyek string `json:"kabupaten_proyek,omitempty"`
	KecamatanProyek string `json:"kecamatan_proyek,omitempty"`
	KelurahanProyek string `json:"kelurahan_proyek,omitempty"`
	AlamatProyek    string `json:"alamat_proyek,omitempty"`
	KodeJalan       string `json:"kode_jalan,omitempty"`
	KodeJalanMerge  string `json:"kode_jalan_merge,omitempty"`
	NamaJalan       string `json:"nama_jalan,omitempty"`
	PangkalJalan    string `json:"pangkal_jalan,omitempty"`
	UjungJalan      string `json:"ujung_jalan,omitempty"`
	PanjangJalan    string `json:"panjang_jalan,omitempty"`
	LebarJalan      string `json:"lebar_jalan,omitempty"`
	PermukaanJalan  string `json:"permukaan_jalan,omitempty"`
	FungsiJalan     string `json:"fungsi_jalan,omitempty"`
	StatusJalan     string `json:"status_jalan,omitempty"`

	//Data Pemohon
	NikPemohon          string  `json:"nik_pemohon,omitempty"`
	NamaPemohon         string  `json:"nama_pemohon,omitempty"`
	EmailPemohon        string  `json:"email_pemohon,omitempty"`
	TempatLahirPemohon  string  `json:"tempat_lahir_pemohon,omitempty"`
	TanggalLahirPemohon string  `json:"tanggal_lahir_pemohon,omitempty"`
	NegaraPemohon       string  `json:"negara_pemohon,omitempty"`
	ProvinsiPemohon     string  `json:"provinsi_pemohon,omitempty"`
	KabupatenPemohon    string  `json:"kabupaten_pemohon,omitempty"`
	KecamatanPemohon    string  `json:"kecamatan_pemohon,omitempty"`
	KelurahanPemohon    string  `json:"kelurahan_pemohon,omitempty"`
	AlamatPemohon       string  `json:"alamat_pemohon,omitempty"`
	JenisKelaminPemohon string  `json:"jenis_kelamin_pemohon,omitempty"`
	NomerPemohon        string  `json:"nomer_pemohon,omitempty"`
	JabatanPemohon      *string `json:"jabatan_pemohon,omitempty"`

	//Data Perusahaan
	NamaPerusahaan              *string `json:"nama_perusahaan,omitempty"`
	NegaraPerusahaan            string  `json:"negara_perusahaan,omitempty"`
	ProvinsiPerusahaan          *string `json:"provinsi_perusahaan,omitempty"`
	KabupatenPerusahaan         *string `json:"kabupaten_perusahaan,omitempty"`
	KecamatanPerusahaan         *string `json:"kecamatan_perusahaan,omitempty"`
	KelurahanPerusahaan         *string `json:"kelurahan_perusahaan,omitempty"`
	AlamatPerusahaan            *string `json:"alamat_perusahaan,omitempty"`
	NomerPerusahaan             *string `json:"nomer_perusahaan,omitempty"`
	EmailPerusahaan             *string `json:"email_perusahaan,omitempty"`
	NamaPimpinan                *string `json:"nama_pimpinan,omitempty"`
	JabatanPimpinan             *string `json:"jabatan_pimpinan,omitempty"`
	JenisKelaminPimpinan        *string `json:"jenis_kelamin,omitempty"`
	NegaraPimpinanPerusahaan    string  `json:"negara_pimpinan_perusahaan,omitempty"`
	ProvinsiPimpinanPerusahaan  *string `json:"provinsi_pimpinan_perusahaan,omitempty"`
	KabupatenPimpinanPerusahaan *string `json:"kabupaten_pimpinan_perusahaan,omitempty"`
	KecamatanPimpinanPerusahaan *string `json:"kecamatan_pimpinan_perusahaan,omitempty"`
	KelurahanPimpinanPerusahaan *string `json:"kelurahan_pimpinan_perusahaan,omitempty"`
	AlamatPimpinan              *string `json:"alamat_pimpinan_perusahaan,omitempty"`

	//Data Konsultan
	NamaKonsultan                  *string `json:"nama_konsultan,omitempty"`
	NegaraKonsultan                string  `json:"negara_konsultan,omitempty"`
	ProvinsiKonsultan              *string `json:"provinsi_konsultan,omitempty"`
	KabupatenKonsultan             *string `json:"kabupaten_konsultan,omitempty"`
	KecamatanKonsultan             *string `json:"kecamatan_konsultan,omitempty"`
	KelurahanKonsultan             *string `json:"kelurahan_konsultan,omitempty"`
	AlamatKonsultan                *string `json:"alamat_konsultan,omitempty"`
	NomerKonsultan                 *string `json:"nomer_konsultan,omitempty"`
	EmailKonsultan                 *string `json:"email_konsultan,omitempty"`
	NamaPenyusunDokumen            *string `json:"nama_penyusun_dokumen,omitempty"`
	JenisKelaminPenyusunDokumen    *string `json:"jenis_kelamin_penyusun_dokumen,omitempty"`
	NegaraPenyusunDokumen          string  `json:"negara_penyusun_dokumen,omitempty"`
	ProvinsiPenyusunDokumen        *string `json:"provinsi_penyusun_dokumen,omitempty"`
	KabupatenPenyusunDokumen       *string `json:"kabupaten_penyusun_dokumen,omitempty"`
	KecamatanPenyusunDokumen       *string `json:"kecamatan_penyusun_dokumen,omitempty"`
	KelurahanPenyusunDokumen       *string `json:"kelurahan_penyusun_dokumen,omitempty"`
	AlamatPenyusunDokumen          *string `json:"alamat_penyusun_dokumen,omitempty"`
	NomerSertifikatPenyusunDokumen *string `json:"nomer_sertifikat_penyusun_dokumen,omitempty"`
	KlasifikasiPenyusunDokumen     *string `json:"klasifikasi_penyusun_dokumen,omitempty"`

	//Data Kegiatan
	Aktivitas         string  `json:"aktivitas,omitempty"`
	Peruntukan        string  `json:"peruntukan,omitempty"`
	TotalLuasLahan    string  `json:"total_luas,omitempty"`
	KriteriaKhusus    *string `json:"kriteria_khusus,omitempty"`
	NilaiKriteria     *string `json:"nilai_kriteria,omitempty"`
	LatitudeBangunan  float64 `json:"latitude,omitempty"`
	LongitudeBangunan float64 `json:"longitude,omitempty"`
	LokasiBangunan    string  `json:"alamat_persil,omitempty"`
	NomerSKRK         string  `json:"nomer_skrk,omitempty"`
	TanggalSKRK       string  `json:"tanggal_skrk,omitempty"`
	Catatan           *string `json:"catatan,omitempty"`

	//Persyaratan tidak terpenuhi
	PersyaratanTidakSesuai []PersayaratanTidakSesuai `json:"persyaratan_tidak_sesuai,omitempty"`

	PertimbanganPenolakan string `json:"pertimbangan_penolakan,omitempty"`
	PertimbanganPenundaan string `json:"pertimbangan_penundaan,omitempty"`

	//Data Pemeriksaan Suat Persetujuan
	HasilPemeriksaan   string  `json:"hasil_pemeriksaan,omitempty"`
	CatatanPemeriksaan *string `json:"catatan_pemeriksaan,omitempty"`

	HasilAsistensiDokumen   string             `json:"hasil_asistensi,omitempty"`
	CatatanAsistensiDokumen []CatatanAsistensi `json:"catatan_asistensi,omitempty"`

	//Berkas Permohonan
	PersyaratanPermohonan []string `json:"persyaratan,omitempty"`
	BerkasPermohonan      []string `json:"berkas,omitempty"`

	//Perlengkapan Tidak Sesuai
	KelengkapanTidakSesuai []KelengkapanTidakSesuaiResponse `json:"kelengkapan,omitempty"`
}

type AndalalinResponseUser struct {
	//Data Pemohon
	IdAndalalin             uuid.UUID `json:"id_andalalin,omitempty"`
	JenisAndalalin          string    `json:"jenis_andalalin,omitempty"`
	Bangkitan               string    `json:"kategori_bangkitan,omitempty"`
	Pemohon                 string    `json:"kategori_pemohon,omitempty"`
	Kode                    string    `json:"kode_andalalin,omitempty"`
	WaktuAndalalin          string    `json:"waktu_andalalin,omitempty"`
	TanggalAndalalin        string    `json:"tanggal_andalalin,omitempty"`
	StatusAndalalin         string    `json:"status_andalalin,omitempty"`
	JenisRencanaPembangunan string    `json:"jenis_rencana_pembangunan,omitempty"`
	Kategori                string    `json:"kategori,omitempty"`
	LokasiPengambilan       string    `json:"lokasi_pengambilan,omitempty"`

	//Data Pemohon
	NikPemohon     string  `json:"nik_pemohon,omitempty"`
	NamaPemohon    string  `json:"nama_pemohon,omitempty"`
	JabatanPemohon *string `json:"jabatan_pemohon,omitempty"`
	EmailPemohon   string  `json:"email_pemohon,omitempty"`
	NomerPemohon   string  `json:"nomer_pemohon,omitempty"`

	//Data Proyek
	NamaProyek      string `json:"nama_proyek,omitempty"`
	JenisProyek     string `json:"jenis_proyek,omitempty"`
	NamaJalan       string `json:"nama_jalan,omitempty"`
	NegaraProyek    string `json:"negara_proyek,omitempty"`
	ProvinsiProyek  string `json:"provinsi_proyek,omitempty"`
	KabupatenProyek string `json:"kabupaten_proyek,omitempty"`
	KecamatanProyek string `json:"kecamatan_proyek,omitempty"`
	KelurahanProyek string `json:"kelurahan_proyek,omitempty"`
	AlamatProyek    string `json:"alamat_proyek,omitempty"`

	//Data Perusahaan
	NamaPerusahaan *string `json:"nama_perusahaan,omitempty"`

	//Data Konsultan
	NamaKonsultan                  *string `json:"nama_konsultan,omitempty"`
	NomerKonsultan                 *string `json:"nomer_konsultan,omitempty"`
	EmailKonsultan                 *string `json:"email_konsultan,omitempty"`
	NamaPenyusunDokumen            *string `json:"nama_penyusun_dokumen,omitempty"`
	NomerSertifikatPenyusunDokumen *string `json:"nomer_sertifikat_penyusun_dokumen,omitempty"`
	KlasifikasiPenyusunDokumen     *string `json:"klasifikasi_penyusun_dokumen,omitempty"`

	//Data Kegiatan
	Aktivitas         string  `json:"aktivitas,omitempty"`
	Peruntukan        string  `json:"peruntukan,omitempty"`
	TotalLuasLahan    string  `json:"total_luas,omitempty"`
	KriteriaKhusus    *string `json:"kriteria_khusus,omitempty"`
	NilaiKriteria     *string `json:"nilai_kriteria,omitempty"`
	LatitudeBangunan  float64 `json:"latitude,omitempty"`
	LongitudeBangunan float64 `json:"longitude,omitempty"`
	LokasiBangunan    string  `json:"alamat_persil,omitempty"`
	Catatan           *string `json:"catatan,omitempty"`

	//Persyaratan tidak terpenuhi
	PersyaratanTidakSesuai []PersayaratanTidakSesuai `json:"persyaratan_tidak_sesuai,omitempty"`

	PertimbanganPenolakan string `json:"pertimbangan_penolakan,omitempty"`
	PertimbanganPenundaan string `json:"pertimbangan_penundaan,omitempty"`

	//Berkas Permohonan
	PersyaratanPermohonan []string `json:"persyaratan,omitempty"`
	BerkasPermohonan      []string `json:"berkas,omitempty"`

	//Perlengkapan Tidak Sesuai
	KelengkapanTidakSesuai []KelengkapanTidakSesuaiResponse `json:"kelengkapan,omitempty"`
}

type Perlengkapan struct {
	IdPerlengkapan       string `json:"id_perlengkapan" binding:"required"`
	StatusPerlengkapan   string
	KategoriUtama        string  `json:"kategori_utama" binding:"required"`
	KategoriPerlengkapan string  `json:"kategori" binding:"required"`
	JenisPerlengkapan    string  `json:"perlengkapan" binding:"required"`
	GambarPerlengkapan   string  `json:"gambar" binding:"required"`
	LokasiPemasangan     string  `json:"pemasangan" binding:"required"`
	LatitudePemasangan   float64 `protobuf:"fixed64,1,opt,name=latitude,proto3" json:"latitude" binding:"required"`
	LongitudePemasangan  float64 `protobuf:"fixed64,2,opt,name=longitude,proto3" json:"longtitude" binding:"required"`
	FotoLokasi           []Foto
	Detail               *string `json:"detail" binding:"required"`
	Alasan               string  `json:"alasan" binding:"required"`
	Pertimbangan         *string
}

type PerlengkapanResponse struct {
	IdPerlengkapan     string `json:"id_perlengkapan,omitempty"`
	StatusPerlengkapan string `json:"status,omitempty"`
	JenisPerlengkapan  string `json:"perlengkapan,omitempty"`
	GambarPerlengkapan string `json:"gambar,omitempty"`
	LokasiPemasangan   string `json:"pemasangan,omitempty"`
}

type Perlalin struct {
	//Data Permohonan
	IdAndalalin      uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	JenisAndalalin   string    `gorm:"type:varchar(255);not null"`
	Kode             string    `gorm:"type:varchar(255);not null"`
	WaktuAndalalin   string    `gorm:"not null"`
	TanggalAndalalin string    `gorm:"not null"`
	StatusAndalalin  string    `gorm:"type:varchar(255);not null"`

	//Perlengkapan
	Perlengkapan []Perlengkapan `gorm:"serializer:json"`

	//Data Pemohon
	IdUser              uuid.UUID `gorm:"type:varchar(255);not null"`
	NikPemohon          string    `gorm:"type:varchar(255);not null"`
	NamaPemohon         string    `gorm:"type:varchar(255);not null"`
	EmailPemohon        string    `gorm:"type:varchar(255);not null"`
	TempatLahirPemohon  string    `gorm:"type:varchar(255);not null"`
	TanggalLahirPemohon string    `gorm:"type:varchar(255);not null"`
	NegaraPemohon       string    `gorm:"type:varchar(255);not null"`
	ProvinsiPemohon     string    `gorm:"type:varchar(255);not null"`
	KabupatenPemohon    string    `gorm:"type:varchar(255);not null"`
	KecamatanPemohon    string    `gorm:"type:varchar(255);not null"`
	KelurahanPemohon    string    `gorm:"type:varchar(255);not null"`
	AlamatPemohon       string    `gorm:"type:varchar(255);not null"`
	JenisKelaminPemohon string    `sql:"type:enum('Laki-laki', 'Perempuan');not null"`
	NomerPemohon        string    `gorm:"type:varchar(255);not null"`

	//Catatan
	Catatan *string

	//Data Petugas
	IdPetugas    uuid.UUID `gorm:"type:varchar(255);"`
	NamaPetugas  string    `gorm:"type:varchar(255);"`
	EmailPetugas string    `gorm:"type:varchar(255);"`

	//Data Persyaratan
	BerkasPermohonan []BerkasPermohonan `gorm:"serializer:json"`

	//Persyaratan tidak terpenuhi
	PersyaratanTidakSesuai []PersayaratanTidakSesuai `gorm:"serializer:json"`

	//Data Tindakan
	Tindakan string

	//Data Pertimbangan
	PertimbanganPembatalan string
	PertimbanganPenolakan  string
	PertimbanganPenundaan  string
}

type InputPerlalin struct {
	Perlengkapan        []Perlengkapan `json:"perlengkapan" binding:"required"`
	NikPemohon          string         `json:"nik_pemohon" binding:"required"`
	TempatLahirPemohon  string         `json:"tempat_lahir_pemohon" binding:"required"`
	TanggalLahirPemohon string         `json:"tanggal_lahir_pemohon" binding:"required"`
	ProvinsiPemohon     string         `json:"provinsi_pemohon" binding:"required"`
	KabupatenPemohon    string         `json:"kabupaten_pemohon" binding:"required"`
	KecamatanPemohon    string         `json:"kecamatan_pemohon" binding:"required"`
	KelurahanPemohon    string         `json:"kelurahan_pemohon" binding:"required"`
	AlamatPemohon       string         `json:"alamat_pemohon" binding:"required"`
	JenisKelaminPemohon string         `json:"jenis_kelamin_pemohon" binding:"required"`
	NomerPemohon        string         `json:"nomer_pemohon" binding:"required"`
	Catatan             *string        `json:"catatan" binding:"required"`
}

type DataPerlalin struct {
	Perlalin InputPerlalin `form:"data"`
}

type PerlalinResponse struct {
	//Data Permohonan
	IdAndalalin      uuid.UUID              `json:"id_andalalin,omitempty"`
	JenisAndalalin   string                 `json:"jenis_andalalin,omitempty"`
	Perlengkapan     []PerlengkapanResponse `json:"perlengkapan,omitempty"`
	Kode             string                 `json:"kode_andalalin,omitempty"`
	WaktuAndalalin   string                 `json:"waktu_andalalin,omitempty"`
	TanggalAndalalin string                 `json:"tanggal_andalalin,omitempty"`
	StatusAndalalin  string                 `json:"status_andalalin,omitempty"`

	//Data Pemohon
	NikPemohon          string `json:"nik_pemohon,omitempty"`
	NamaPemohon         string `json:"nama_pemohon,omitempty"`
	EmailPemohon        string `json:"email_pemohon,omitempty"`
	TempatLahirPemohon  string `json:"tempat_lahir_pemohon,omitempty"`
	TanggalLahirPemohon string `json:"tanggal_lahir_pemohon,omitempty"`
	NegaraPemohon       string `json:"negara_pemohon,omitempty"`
	ProvinsiPemohon     string `json:"provinsi_pemohon,omitempty"`
	KabupatenPemohon    string `json:"kabupaten_pemohon,omitempty"`
	KecamatanPemohon    string `json:"kecamatan_pemohon,omitempty"`
	KelurahanPemohon    string `json:"kelurahan_pemohon,omitempty"`
	AlamatPemohon       string `json:"alamat_pemohon,omitempty"`
	JenisKelaminPemohon string `json:"jenis_kelamin_pemohon,omitempty"`
	NomerPemohon        string `json:"nomer_pemohon,omitempty"`

	//Catatan
	Catatan *string `json:"catatan,omitempty"`

	//Persyaratan tidak terpenuhi
	PersyaratanTidakSesuai []PersayaratanTidakSesuai `json:"persyaratan_tidak_sesuai,omitempty"`

	//Data Petugas
	IdPetugas         uuid.UUID `json:"id_petugas,omitempty"`
	NamaPetugas       string    `json:"nama_petugas,omitempty"`
	EmailPetugas      string    `json:"email_petugas,omitempty"`
	StatusTiketLevel2 string    `json:"status_tiket,omitempty"`

	//Berkas Permohonan
	PersyaratanPermohonan []string `json:"persyaratan,omitempty"`
	BerkasPermohonan      []string `json:"berkas,omitempty"`

	//Data Pertimbangan
	PertimbanganPembatalan string `json:"pertimbangan_pembatalan,omitempty"`
	PertimbanganPenolakan  string `json:"pertimbangan_penolakan,omitempty"`
	PertimbanganPenundaan  string `json:"pertimbangan_penundaan,omitempty"`
}

type PerlalinResponseUser struct {
	//Data Permohonan
	IdAndalalin      uuid.UUID              `json:"id_andalalin,omitempty"`
	JenisAndalalin   string                 `json:"jenis_andalalin,omitempty"`
	Kode             string                 `json:"kode_andalalin,omitempty"`
	WaktuAndalalin   string                 `json:"waktu_andalalin,omitempty"`
	TanggalAndalalin string                 `json:"tanggal_andalalin,omitempty"`
	StatusAndalalin  string                 `json:"status_andalalin,omitempty"`
	Perlengkapan     []PerlengkapanResponse `json:"perlengkapan,omitempty"`

	//Data Pemohon
	NikPemohon   string `json:"nik_pemohon,omitempty"`
	NamaPemohon  string `json:"nama_pemohon,omitempty"`
	EmailPemohon string `json:"email_pemohon,omitempty"`
	NomerPemohon string `json:"nomer_pemohon,omitempty"`

	//Catatan
	Catatan *string `json:"catatan,omitempty"`

	//Persyaratan tidak terpenuhi
	PersyaratanTidakSesuai []PersayaratanTidakSesuai `json:"persyaratan_tidak_sesuai,omitempty"`

	//Berkas Permohonan
	PersyaratanPermohonan []string `json:"persyaratan,omitempty"`
	BerkasPermohonan      []string `json:"berkas,omitempty"`

	//Data Pertimbangan
	PertimbanganPembatalan string `json:"pertimbangan_pembatalan,omitempty"`
	PertimbanganPenolakan  string `json:"pertimbangan_penolakan,omitempty"`
	PertimbanganPenundaan  string `json:"pertimbangan_penundaan,omitempty"`
}

type CatatanAsistensi struct {
	Substansi string   `json:"substansi" binding:"required"`
	Catatan   []string `json:"catatan" binding:"required"`
}

type BerkasPermohonan struct {
	Nama   string
	Tipe   string
	Status string
	Berkas []byte
}

type KelengkapanTidakSesuai struct {
	Dokumen string
	Tipe    string
	Role    string
}

type KelengkapanTidakSesuaiResponse struct {
	Dokumen string `json:"dokumen,omitempty"`
	Tipe    string `json:"tipe,omitempty"`
}

type PersayaratanTidakSesuai struct {
	Persyaratan string `json:"persyaratan,omitempty"`
	Tipe        string `json:"tipe,omitempty"`
}

type DaftarAndalalinResponse struct {
	IdAndalalin      uuid.UUID `json:"id_andalalin,omitempty"`
	Kode             string    `json:"kode_andalalin,omitempty"`
	TanggalAndalalin string    `json:"tanggal_andalalin,omitempty"`
	Nama             string    `json:"nama_pemohon,omitempty"`
	Email            string    `json:"email_pemohon,omitempty"`
	Petugas          string    `json:"petugas,omitempty"`
	JenisAndalalin   string    `json:"jenis_andalalin,omitempty"`
	StatusAndalalin  string    `json:"status_andalalin,omitempty"`
}

type Pemeriksaan struct {
	Hasil   string  `json:"hasil" binding:"required"`
	Catatan *string `json:"catatan" binding:"required"`
}

type Pertimbangan struct {
	Pertimbangan string `json:"pertimbangan" binding:"required"`
}

type TambahPetugas struct {
	IdPetugas    uuid.UUID `json:"id_petugas" binding:"required"`
	NamaPetugas  string    `json:"nama_petugas" binding:"required"`
	EmailPetugas string    `json:"email_petugas" binding:"required"`
}

type Survei struct {
	IdSurvey       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	IdAndalalin    uuid.UUID `gorm:"type:varchar(255);not null"`
	IdTiketLevel1  uuid.UUID `gorm:"type:varchar(255);not null"`
	IdTiketLevel2  uuid.UUID `gorm:"type:varchar(255);not null"`
	IdPetugas      uuid.UUID `gorm:"type:varchar(255);not null"`
	IdPerlengkapan string    `gorm:"type:varchar(255);not null"`
	Petugas        string    `gorm:"type:varchar(255);not null"`
	EmailPetugas   string    `gorm:"type:varchar(255);not null"`
	Catatan        *string
	Foto           []Foto `gorm:"serializer:json"`
	Lokasi         string
	Latitude       float64
	Longitude      float64
	WaktuSurvei    string `gorm:"not null"`
	TanggalSurvei  string `gorm:"not null"`
}

type Foto []byte

type DataFoto struct {
	Id   string
	Foto []byte
}

type InputSurvey struct {
	Lokasi    string  `json:"lokasi" binding:"required"`
	Catatan   *string `json:"catatan" binding:"required"`
	Latitude  float64 `protobuf:"fixed64,1,opt,name=latitude,proto3" json:"latitude" binding:"required"`
	Longitude float64 `protobuf:"fixed64,2,opt,name=longitude,proto3" json:"longtitude" binding:"required"`
}

type DataSurvey struct {
	Data InputSurvey `form:"data"`
}

type DataLaporanSurvei struct {
	Perlengkapan string
	Lokasi       string
	Tanggal      string
	Survei       string
	Catatan      *string
	Foto         []Foto
}

type TiketLevel1 struct {
	IdTiketLevel1 uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	IdAndalalin   uuid.UUID `gorm:"type:varchar(255);unique;not null"`
	Status        string    `sql:"type:enum('Buka', 'Tutup', 'Tunda', 'Batal');not null"`
}

type TiketLevel2 struct {
	IdTiketLevel2 uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	IdTiketLevel1 uuid.UUID `gorm:"type:varchar(255);not null"`
	IdAndalalin   uuid.UUID `gorm:"type:varchar(255);not null"`
	IdPetugas     uuid.UUID `gorm:"type:varchar(255);not null"`
	Status        string    `sql:"type:enum('Buka', 'Tutup', 'Tunda', 'Batal');not null"`
}

type Pemasangan struct {
	IdPemasangan      uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	IdAndalalin       uuid.UUID `gorm:"type:varchar(255);not null"`
	IdTiketLevel1     uuid.UUID `gorm:"type:varchar(255);not null"`
	IdPetugas         uuid.UUID `gorm:"type:varchar(255);not null"`
	IdPerlengkapan    string    `gorm:"type:varchar(255);not null"`
	Petugas           string    `gorm:"type:varchar(255);not null"`
	EmailPetugas      string    `gorm:"type:varchar(255);not null"`
	Catatan           *string
	Foto              []Foto `gorm:"serializer:json"`
	Lokasi            string
	Latitude          float64
	Longitude         float64
	WaktuPemasangan   string `gorm:"not null"`
	TanggalPemasangan string `gorm:"not null"`
}

type Pengecekan struct {
	Data []DataPengecekan `json:"data" binding:"required"`
}

type DataPengecekan struct {
	ID           string  `json:"id" binding:"required"`
	Perlengkapan string  `json:"perlengkapan" binding:"required"`
	Pambar       string  `json:"gambar" binding:"required"`
	Lokasi       string  `json:"lokasi" binding:"required"`
	Setuju       string  `json:"setuju" binding:"required"`
	Tidak        string  `json:"tidak" binding:"required"`
	Pertimbangan *string `json:"pertimbangan" binding:"required"`
}

type Administrasi struct {
	NomorSurat   string             `json:"nomor" binding:"required"`
	TanggalSurat string             `json:"tanggal" binding:"required"`
	Data         []DataAdministrasi `json:"data" binding:"required"`
}

type DataAdministrasi struct {
	Persyaratan string `json:"persyaratan" binding:"required"`
	Kebutuhan   string `json:"kebutuhan" binding:"required"`
	Tipe        string `json:"tipe" binding:"required"`
	Ada         string `json:"ada" binding:"required"`
	Tidak       string `json:"tidak" binding:"required"`
	Keterangan  string `json:"keterangan" binding:"required"`
}

type AdministrasiPerlalin struct {
	Data []DataAdministrasiPerlalin `json:"data" binding:"required"`
}

type DataAdministrasiPerlalin struct {
	Persyaratan string `json:"persyaratan" binding:"required"`
	Kebutuhan   string `json:"kebutuhan" binding:"required"`
	Tipe        string `json:"tipe" binding:"required"`
	Ada         string `json:"ada" binding:"required"`
	Tidak       string `json:"tidak" binding:"required"`
}

type Kewajiban struct {
	Kewajiban []string `json:"kewajiban" binding:"required"`
}

type Keputusan struct {
	NomorKeputusan     string          `json:"nomor_keputusan" binding:"required"`
	NomorLampiran      string          `json:"nomor_lampiran" binding:"required"`
	NomorKesanggupan   string          `json:"nomor_kesanggupan" binding:"required"`
	TanggalKesanggupan string          `json:"tanggal_kesanggupan" binding:"required"`
	NamaKadis          string          `json:"nama_kadis" binding:"required"`
	NipKadis           string          `json:"nip_kadis" binding:"required"`
	Data               []DataKeputusan `json:"data" binding:"required"`
}

type DataKeputusan struct {
	Kewajiban     string          `json:"kewajiban" binding:"required"`
	DataKewajiban []DataKewajiban `json:"data_kewajiban" binding:"required"`
}

type DataKewajiban struct {
	Poin     string     `json:"poin" binding:"required"`
	DataPoin []DataPoin `json:"data_poin" binding:"required"`
}

type DataPoin struct {
	Subpoin     string   `json:"subpoin" binding:"required"`
	DataSubpoin []string `json:"data_subpoin" binding:"required"`
}

type KelengkapanAkhir struct {
	Kelengkapan []DataKelengkapanAkhir `json:"kelengkapan" binding:"required"`
}

type DataKelengkapanAkhir struct {
	Uraian     string        `json:"uraian" binding:"required"`
	Dokumen    []DataDokumen `json:"dokumen" binding:"required"`
	Role       string        `json:"role" binding:"required"`
	Ada        string        `json:"ada" binding:"required"`
	Tidak      string        `json:"tidak" binding:"required"`
	Keterangan string        `json:"keterangan" binding:"required"`
}

type DataDokumen struct {
	Dokumen string `json:"dokumen" binding:"required"`
	Tipe    string `json:"tipe" binding:"required"`
}

type PenyusunDokumen struct {
	Penyusun []DataPenyusunDokumen `json:"penyusun" binding:"required"`
}

type DataPenyusunDokumen struct {
	Substansi string       `json:"substansi" binding:"required"`
	Muatan    []DataMuatan `json:"muatan" binding:"required"`
}

type DataMuatan struct {
	Judul    string         `json:"judul" binding:"required"`
	Tambahan []DataTambahan `json:"tambahan" binding:"required"`
}

type DataTambahan struct {
	Judul string   `json:"judul" binding:"required"`
	Poin  []string `json:"poin" binding:"required"`
}

type PemeriksaanDokumen struct {
	Status      string             `json:"status" binding:"required"`
	Pemeriksaan []CatatanAsistensi `json:"pemeriksaan" binding:"required"`
}

type PerbaruiLokasi struct {
	Lokasi    string  `json:"lokasi" binding:"required"`
	Latitude  float64 `protobuf:"fixed64,1,opt,name=latitude,proto3" json:"latitude" binding:"required"`
	Longitude float64 `protobuf:"fixed64,2,opt,name=longitude,proto3" json:"longtitude" binding:"required"`
}

type DataSuratPermohonan struct {
	Bangkitan   string  `json:"bangkitan" binding:"required"`
	Pemohon     string  `json:"pemohon" binding:"required"`
	Nama        string  `json:"nama" binding:"required"`
	Jabatan     *string `json:"jabatan" binding:"required"`
	Jenis       string  `json:"jenis" binding:"required"`
	Proyek      string  `json:"proyek" binding:"required"`
	Jalan       string  `json:"jalan" binding:"required"`
	Kelurahan   string  `json:"kelurahan" binding:"required"`
	Kecamatan   string  `json:"kecamatan" binding:"required"`
	Kabupaten   string  `json:"kabupaten" binding:"required"`
	Provinsi    string  `json:"provinsi" binding:"required"`
	StatusJalan string  `json:"status" binding:"required"`
	Pengembang  string  `json:"pengembang" binding:"required"`
	Konsultan   *string `json:"konsultan" binding:"required"`
}
