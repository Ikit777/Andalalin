package models

type DataMaster struct {
	LokasiPengambilan       Lokasi  `gorm:"serializer:json"`
	JenisRencanaPembangunan Jenis   `gorm:"serializer:json"`
	RencanaPembangunan      Rencana `gorm:"serializer:json"`
}

type Rencana struct {
	Pusat         []string
	Pemukiman     []string
	Infrastruktur []string
}

type Lokasi []string

type Jenis []string
