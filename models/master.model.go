package models

import "github.com/google/uuid"

type DataMaster struct {
	IdDataMaster            uuid.UUID           `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	JenisProyek             JenisProyek         `gorm:"serializer:json"`
	LokasiPengambilan       Lokasi              `gorm:"serializer:json"`
	JenisRencanaPembangunan Jenis               `gorm:"serializer:json"`
	RencanaPembangunan      []Rencana           `gorm:"serializer:json"`
	KategoriPerlengkapan    Perlengkapan        `gorm:"serializer:json"`
	PerlengkapanLaluLintas  []JenisPerlengkapan `gorm:"serializer:json"`
	Persyaratan             Persyaratan         `gorm:"serializer:json"`
	Provinsi                []Provinsi          `gorm:"serializer:json"`
	Kabupaten               []Kabupaten         `gorm:"serializer:json"`
	Kecamatan               []Kecamatan         `gorm:"serializer:json"`
	Kelurahan               []Kelurahan         `gorm:"serializer:json"`
	Jalan                   []Jalan             `gorm:"serializer:json"`
	UpdatedAt               string              `gorm:"not null"`
}

type Lokasi []string

type Jenis []string

type Perlengkapan []string

type JenisProyek []string

type JenisPerlengkapan struct {
	Kategori     string
	Perlengkapan []PerlengkapanItem
}

type PerlengkapanItem struct {
	JenisPerlengkapan  string
	GambarPerlengkapan []byte
}

type Provinsi struct {
	Id   string
	Name string
}

type Kabupaten struct {
	Id         string
	IdProvinsi string
	Name       string
}

type Kecamatan struct {
	Id          string
	IdKabupaten string
	Name        string
}

type Kelurahan struct {
	Id          string
	IdKecamatan string
	Name        string
}

type Rencana struct {
	Kategori     string
	JenisRencana []JenisRencana
}

type JenisRencana struct {
	Jenis    string
	Kriteria string
	Satuan   string
}

type Jalan struct {
	KodeProvinsi  string
	KodeKabupaten string
	KodeKecamatan string
	KodeKelurahan string
	KodeJalan     string
	Nama          string
	Pangkal       string
	Ujung         string
	Kelurahan     string
	Kecamatan     string
	Panjang       string
	Lebar         string
	Permukaan     string
	Fungsi        string
}

type Persyaratan struct {
	PersyaratanAndalalin []PersyaratanAndalalinInput
	PersyaratanPerlalin  []PersyaratanPerlalinInput
}

type PersyaratanAndalalinInput struct {
	Bangkitan             string `json:"bangkitan" binding:"required"`
	Persyaratan           string `json:"persyaratan" binding:"required"`
	KeteranganPersyaratan string `json:"keterangan" binding:"required"`
}

type PersyaratanPerlalinInput struct {
	Persyaratan           string `json:"persyaratan" binding:"required"`
	KeteranganPersyaratan string `json:"keterangan" binding:"required"`
}
