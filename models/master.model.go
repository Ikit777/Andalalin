package models

import "github.com/google/uuid"

type DataMaster struct {
	IdDataMaster            uuid.UUID           `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	LokasiPengambilan       Lokasi              `gorm:"serializer:json"`
	JenisRencanaPembangunan Jenis               `gorm:"serializer:json"`
	RencanaPembangunan      []Rencana           `gorm:"serializer:json"`
	PersyaratanTambahan     PersyaratanTambahan `gorm:"serializer:json"`
}

type Lokasi []string

type Jenis []string

type Rencana struct {
	Kategori     string
	JenisRencana []string
}

type PersyaratanTambahan struct {
	PersyaratanTambahanAndalalin  []PersyaratanTambahanStruct
	PersyaratanTambahanRambulalin []PersyaratanTambahanStruct
}

type PersyaratanTambahanInput struct {
	Persyaratan           string `json:"persyaratan" binding:"required"`
	KeteranganPersyaratan string `json:"keterangan" binding:"required"`
}

type PersyaratanTambahanStruct struct {
	Persyaratan           string
	KeteranganPersyaratan string
}
