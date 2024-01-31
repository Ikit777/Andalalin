package models

import (
	"github.com/google/uuid"
)

type SurveiKepuasan struct {
	IdSurvey           uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	IdAndalalin        uuid.UUID `gorm:"type:varchar(255);not null"`
	IdUser             uuid.UUID `gorm:"type:varchar(255);not null"`
	Nama               string    `gorm:"type:varchar(255);not null"`
	Email              string    `gorm:"type:varchar(255);not null"`
	KritikSaran        *string
	TanggalPelaksanaan string
	DataSurvei         []Kepuasan `gorm:"serializer:json"`
}

type SurveiKepuasanInput struct {
	KritikSaran *string    `json:"saran" binding:"required"`
	DataSurvei  []Kepuasan `json:"data" binding:"required"`
}

type Kepuasan struct {
	Jenis string
	Nilai string
}

type SurveiMandiri struct {
	IdSurvey        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	IdPetugas       uuid.UUID `gorm:"type:varchar(255);not null"`
	Petugas         string    `gorm:"type:varchar(255);not null"`
	EmailPetugas    string    `gorm:"type:varchar(255);not null"`
	Catatan         *string
	Foto            []Foto `gorm:"serializer:json"`
	Lokasi          string
	Latitude        float64
	Longitude       float64
	WaktuSurvei     string `gorm:"not null"`
	TanggalSurvei   string `gorm:"not null"`
	StatusSurvei    string
	CatatanTindakan *string
}

type TerimaSurveiMandiri struct {
	Catatan *string `json:"catatan" binding:"required"`
}
