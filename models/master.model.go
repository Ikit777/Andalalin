package models

import "github.com/google/uuid"

type DataMaster struct {
	IdDataMaster            uuid.UUID           `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	LokasiPengambilan       Lokasi              `gorm:"serializer:json"`
	JenisRencanaPembangunan Jenis               `gorm:"serializer:json"`
	RencanaPembangunan      []Rencana           `gorm:"serializer:json"`
	KategoriPerlengkapan    Perlengkapan        `gorm:"serializer:json"`
	PerlengkapanLaluLintas  []JenisPerlengkapan `gorm:"serializer:json"`
	PersyaratanTambahan     PersyaratanTambahan `gorm:"serializer:json"`
	Provinsi                []Provinsi          `gorm:"serializer:json"`
	Kabupaten               []Kabupaten         `gorm:"serializer:json"`
	Kecamatan               []Kecamatan         `gorm:"serializer:json"`
	Kelurahan               []Kelurahan         `gorm:"serializer:json"`
	UpdatedAt               string              `gorm:"not null"`
}

type Lokasi []string

type Jenis []string

type Perlengkapan []string

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
	JenisRencana []string
}

type JenisPerlengkapan struct {
	Kategori     string
	Perlengkapan []PerlengkapanItem
}

type PerlengkapanItem struct {
	JenisPerlengkapan  string
	GambarPerlengkapan []byte
}

type PersyaratanTambahan struct {
	PersyaratanTambahanAndalalin []PersyaratanTambahanInput
	PersyaratanTambahanPerlalin  []PersyaratanTambahanInput
}

type PersyaratanTambahanInput struct {
	Persyaratan           string `json:"persyaratan" binding:"required"`
	KeteranganPersyaratan string `json:"keterangan" binding:"required"`
}
