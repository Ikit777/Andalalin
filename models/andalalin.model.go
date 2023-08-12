package models

import (
	"github.com/google/uuid"
)

type Andalalin struct {
	//Data Pemohon
	IdAndalalin            uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	IdUser                 uuid.UUID `gorm:"type:varchar(255);not null"`
	IdPetugas              uuid.UUID `gorm:"type:varchar(255);"`
	JenisAndalalin         string    `gorm:"type:varchar(255);not null"`
	NomerAndalalin         string    `gorm:"type:varchar(255);not null"`
	NikPemohon             string    `gorm:"type:varchar(255);not null"`
	NamaPemohon            string    `gorm:"type:varchar(255);not null"`
	EmailPemohon           string    `gorm:"type:varchar(255);not null"`
	TempatLahirPemohon     string    `gorm:"type:varchar(255);not null"`
	TanggalLahirPemohon    string    `gorm:"type:varchar(255);not null"`
	AlamatPemohon          string    `gorm:"type:varchar(255);not null"`
	JenisKelaminPemohon    string    `sql:"type:enum('Laki-laki', 'Perempuan');not null"`
	NomerPemohon           string    `gorm:"type:varchar(255);not null"`
	JabatanPemohon         string    `gorm:"type:varchar(255);not null"`
	LokasiPengambilan      string    `gorm:"type:varchar(255);not null"`
	WaktuAndalalin         string    `gorm:"not null"`
	StatusAndalalin        string    `sql:"type:enum('Cek persyaratan', 'Persayaratan tidak sesuai', 'Persyaratan terpenuhi', 'Survey lapangan', 'Laporan BAP', 'Pembuatan SK', 'Permohonan selesai')"`
	NamaPetugas            string    `gorm:"type:varchar(255);"`
	EmailPetugas           string    `gorm:"type:varchar(255);"`
	TandaTerimaPendaftaran []byte

	//Data Perusahaan
	NamaPerusahaan       string `gorm:"type:varchar(255);not null"`
	AlamatPerusahaan     string `gorm:"type:varchar(255);not null"`
	NomerPerusahaan      string `gorm:"type:varchar(255);not null"`
	EmailPerusahaan      string `gorm:"type:varchar(255);not null"`
	ProvinsiPerusahaan   string `gorm:"type:varchar(255);not null"`
	KabupatenPerusahaan  string `gorm:"type:varchar(255);not null"`
	KecamatanPerusahaan  string `gorm:"type:varchar(255);not null"`
	KelurahaanPerusahaan string `gorm:"type:varchar(255);not null"`
	NamaPimpinan         string `gorm:"type:varchar(255);not null"`
	JabatanPimpinan      string `gorm:"type:varchar(255);not null"`
	JenisKelaminPimpinan string `sql:"type:enum('Laki-laki', 'Perempuan');not null"`
	JenisKegiatan        string `gorm:"type:varchar(255);not null"`
	Peruntukan           string `gorm:"type:varchar(255);not null"`
	LuasLahan            string `gorm:"type:varchar(255);not null"`
	AlamatPersil         string `gorm:"type:varchar(255);not null"`
	KelurahanPersil      string `gorm:"type:varchar(255);not null"`
	NomerSKRK            string `gorm:"type:varchar(255);not null"`
	TanggalSKRK          string `gorm:"type:varchar(255);not null"`

	//Data Persyaratan
	KartuTandaPenduduk []byte
	AktaPendirianBadan []byte
	SuratKuasa         []byte

	//Data Persetujuan Dokumen
	PersetujuanDokumen           string `gorm:"type:varchar(255);"`
	KeteranganPersetujuanDokumen string

	//Data BAP
	NomerBAPDasar       string `gorm:"type:varchar(255);"`
	NomerBAPPelaksanaan string `gorm:"type:varchar(255);"`
	TanggalBAP          string `gorm:"type:varchar(255);"`
	FileBAP             []byte

	//Data SK
	FileSK []byte
}

type InputAndalalin struct {
	JenisAndalalin       string `json:"jenis" binding:"required"`
	NikPemohon           string `json:"nik_pemohon" binding:"required"`
	NamaPemohon          string `json:"nama_pemohon" binding:"required"`
	TempatLahirPemohon   string `json:"tempat_lahir_pemohon" binding:"required"`
	TanggalLahirPemohon  string `json:"tanggal_lahir_pemohon" binding:"required"`
	AlamatPemohon        string `json:"alamat_pemohon" binding:"required"`
	JenisKelaminPemohon  string `json:"jenis_kelamin_pemohon" binding:"required"`
	NomerPemohon         string `json:"nomer_pemohon" binding:"required"`
	JabatanPemohon       string `json:"jabatan_pemohon" binding:"required"`
	LokasiPengambilan    string `json:"lokasi_pengambilan" binding:"required"`
	NamaPerusahaan       string `json:"nama_perusahaan" binding:"required"`
	AlamatPerusahaan     string `json:"alamat_perusahaan" binding:"required"`
	NomerPerusahaan      string `json:"nomer_perusahaan" binding:"required"`
	EmailPerusahaan      string `json:"email_perusahaan" binding:"required"`
	ProvinsiPerusahaan   string `json:"provinsi_perusahaan" binding:"required"`
	KabupatenPerusahaan  string `json:"kabupaten_perusahaan" binding:"required"`
	KecamatanPerusahaan  string `json:"kecamatan_perusahaan" binding:"required"`
	KelurahaanPerusahaan string `json:"kelurahan_perusahaan" binding:"required"`
	NamaPimpinan         string `json:"nama_pimpinan" binding:"required"`
	JabatanPimpinan      string `json:"jabatan_pimpinan" binding:"required"`
	JenisKelaminPimpinan string `json:"jenis_kelamin_pimpinan" binding:"required"`
	JenisKegiatan        string `json:"jenis_kegiatan" binding:"required"`
	Peruntukan           string `json:"peruntukan" binding:"required"`
	LuasLahan            string `json:"luas_lahan" binding:"required"`
	AlamatPersil         string `json:"alamat_persil" binding:"required"`
	KelurahanPersil      string `json:"kelurahan_persil" binding:"required"`
	NomerSKRK            string `json:"nomer_skrk" binding:"required"`
	TanggalSKRK          string `json:"tanggal_skrk" binding:"required"`
}

type DataAndalalin struct {
	Andalalin InputAndalalin `form:"data"`
}

type DaftarAndalalinResponse struct {
	IdAndalalin     uuid.UUID `json:"id_andalalin,omitempty"`
	NomerAndalalin  string    `json:"nomer_andalalin,omitempty"`
	WaktuAndalalin  string    `json:"waktu_andalalin,omitempty"`
	Nama            string    `json:"nama_pemohon,omitempty"`
	Alamat          string    `json:"alamat_pemohon,omitempty"`
	JenisAndalalin  string    `json:"jenis_andalalin,omitempty"`
	StatusAndalalin string    `json:"status_andalalin,omitempty"`
}

type PersayaratanRespone struct {
	IdAndalalin        uuid.UUID `json:"id_andalalin,omitempty"`
	KartuTandaPenduduk []byte    `json:"ktp,omitempty"`
	AktaPendirianBadan []byte    `json:"akta_pendirian_badan,omitempty"`
	SuratKuasa         []byte    `json:"surat_kuasa,omitempty"`
}

type PerusahaanRespone struct {
	NamaPerusahaan       string `json:"nama_perusahaan,omitempty"`
	AlamatPerusahaan     string `json:"alamat_perusahaan,omitempty"`
	NomerPerusahaan      string `json:"nomer_perusahaan,omitempty"`
	EmailPerusahaan      string `json:"email_perusahaan,omitempty"`
	ProvinsiPerusahaan   string `json:"perovinsi_perusahaan,omitempty"`
	KabupatenPerusahaan  string `json:"kabupaten_perusahaan,omitempty"`
	KecamatanPerusahaan  string `json:"kecamatan_perusahaan,omitempty"`
	KelurahaanPerusahaan string `json:"kelurahan_perusahaan,omitempty"`
	NamaPimpinan         string `json:"nama_pimpinan,omitempty"`
	JabatanPimpinan      string `json:"jabatan_pimpinan,omitempty"`
	JenisKelaminPimpinan string `json:"jenis_kelamin,omitempty"`
	JenisKegiatan        string `json:"jenis_kegiatan,omitempty"`
	Peruntukan           string `json:"peruntukan,omitempty"`
	LuasLahan            string `json:"luas_lahan,omitempty"`
	AlamatPersil         string `json:"alamat_persil,omitempty"`
	KelurahanPersil      string `json:"kelurahan_persil,omitempty"`
	NomerSKRK            string `json:"nomer_skrk,omitempty"`
	TanggalSKRK          string `json:"tanggal_skrk,omitempty"`
}

type AndalalinResponse struct {
	//Data Pemohon
	IdAndalalin            uuid.UUID `json:"id_andalalin,omitempty"`
	JenisAndalalin         string    `json:"jenis_andalalin,omitempty"`
	NomerAndalalin         string    `json:"nomer_andalalin,omitempty"`
	NikPemohon             string    `json:"nik_pemohon,omitempty"`
	NamaPemohon            string    `json:"nama_pemohon,omitempty"`
	EmailPemohon           string    `json:"email_pemohon,omitempty"`
	TempatLahirPemohon     string    `json:"tempat_lahir_pemohon,omitempty"`
	TanggalLahirPemohon    string    `json:"tanggal_lahir_pemohon,omitempty"`
	AlamatPemohon          string    `json:"alamat_pemohon,omitempty"`
	JenisKelaminPemohon    string    `json:"jenis_kelamin_pemohon,omitempty"`
	NomerPemohon           string    `json:"nomer_pemohon,omitempty"`
	JabatanPemohon         string    `json:"jabatan_pemohon,omitempty"`
	LokasiPengambilan      string    `json:"lokasi_pengambilan,omitempty"`
	WaktuAndalalin         string    `json:"waktu_andalalin,omitempty"`
	StatusAndalalin        string    `json:"status_andalalin,omitempty"`
	TandaTerimaPendaftaran []byte    `json:"tanda_terima_pendaftaran,omitempty"`

	//Data Perusahaan
	NamaPerusahaan       string `json:"nama_perusahaan,omitempty"`
	AlamatPerusahaan     string `json:"alamat_perusahaan,omitempty"`
	NomerPerusahaan      string `json:"nomer_perusahaan,omitempty"`
	EmailPerusahaan      string `json:"email_perusahaan,omitempty"`
	ProvinsiPerusahaan   string `json:"perovinsi_perusahaan,omitempty"`
	KabupatenPerusahaan  string `json:"kabupaten_perusahaan,omitempty"`
	KecamatanPerusahaan  string `json:"kecamatan_perusahaan,omitempty"`
	KelurahaanPerusahaan string `json:"kelurahan_perusahaan,omitempty"`
	NamaPimpinan         string `json:"nama_pimpinan,omitempty"`
	JabatanPimpinan      string `json:"jabatan_pimpinan,omitempty"`
	JenisKelaminPimpinan string `json:"jenis_kelamin,omitempty"`
	JenisKegiatan        string `json:"jenis_kegiatan,omitempty"`
	Peruntukan           string `json:"peruntukan,omitempty"`
	LuasLahan            string `json:"luas_lahan,omitempty"`
	AlamatPersil         string `json:"alamat_persil,omitempty"`
	KelurahanPersil      string `json:"kelurahan_persil,omitempty"`
	NomerSKRK            string `json:"nomer_skrk,omitempty"`
	TanggalSKRK          string `json:"tanggal_skrk,omitempty"`

	//Data Persyaratan
	KartuTandaPenduduk []byte `json:"ktp,omitempty"`
	AktaPendirianBadan []byte `json:"akta_pendirian_badan,omitempty"`
	SuratKuasa         []byte `json:"surat_kuasa,omitempty"`

	//Data Petugas
	IdPetugas    uuid.UUID `json:"id_petugas,omitempty"`
	NamaPetugas  string    `json:"nama_petugas,omitempty"`
	EmailPetugas string    `json:"email_petugas,omitempty"`

	//Data Persetujuan
	PersetujuanDokumen           string `json:"persetujuan,omitempty"`
	KeteranganPersetujuanDokumen string `json:"keterangan_persetujuan,omitempty"`

	//Data BAP
	NomerBAPDasar       string `json:"nomer_bap_dasar,omitempty"`
	NomerBAPPelaksanaan string `json:"nomer_bap_pelaksanaan,omitempty"`
	TanggalBAP          string `json:"tanggal_bap,omitempty"`
	FileBAP             []byte `json:"file_bap,omitempty"`

	//Data SK
	FileSK []byte `json:"file_sk,omitempty"`
}

type PersayaratanTidakSesuai struct {
	Persyaratan string `json:"persyaratan" binding:"required"`
}

type Persetujuan struct {
	Persetujuan string `json:"persetujuan" binding:"required"`
	Keterangan  string `json:"keterangan" binding:"required"`
}

type InputBAP struct {
	NomerBAPDasar       string `json:"nomer_dasar" binding:"required"`
	NomerBAPPelaksanaan string `json:"nomer_pelaksanaan" binding:"required"`
	TanggalBAP          string `json:"tanggal" binding:"required"`
}

type BAPData struct {
	Data InputBAP `form:"data"`
}

type TambahPetugas struct {
	IdPetugas    uuid.UUID `json:"id_petugas" binding:"required"`
	NamaPetugas  string    `json:"nama_petugas" binding:"required"`
	EmailPetugas string    `json:"email_petugas" binding:"required"`
}

type Survey struct {
	IdSurvey      uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	IdAndalalin   uuid.UUID `gorm:"type:varchar(255);uniqueIndex;not null"`
	IdTiketLevel1 uuid.UUID `gorm:"type:varchar(255);not null"`
	IdTiketLevel2 uuid.UUID `gorm:"type:varchar(255);not null"`
	IdPetugas     uuid.UUID `gorm:"type:varchar(255);not null"`
	Petugas       string    `gorm:"type:varchar(255);not null"`
	Keterangan    string
	Foto1         []byte
	Foto2         []byte
	Foto3         []byte
	Latitude      float64
	Longitude     float64
}

type InputSurvey struct {
	Keterangan string  `json:"keterangan" binding:"required"`
	Latitude   float64 `protobuf:"fixed64,1,opt,name=latitude,proto3" json:"latitude" binding:"required"`
	Longitude  float64 `protobuf:"fixed64,2,opt,name=longitude,proto3" json:"longtitude" binding:"required"`
}

type DataSurvey struct {
	Data InputSurvey `form:"data"`
}

type TiketLevel1 struct {
	IdTiketLevel1 uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	IdAndalalin   uuid.UUID `gorm:"type:varchar(255);uniqueIndex;not null"`
	Status        string    `sql:"type:enum('Buka', 'Tutup', 'Tunda', 'Batal');not null"`
}

type TiketLevel2 struct {
	IdTiketLevel2 uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	IdTiketLevel1 uuid.UUID `gorm:"type:varchar(255);not null"`
	IdAndalalin   uuid.UUID `gorm:"type:varchar(255);not null"`
	Status        string    `sql:"type:enum('Buka', 'Tutup', 'Tunda', 'Batal');not null"`
}