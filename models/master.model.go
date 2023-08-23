package models

type DataMaster struct {
	LokasiPengambilan       Lokasi  `gorm:"type:string[]"`
	JenisRencanaPembangunan Jenis   `gorm:"type:string[]"`
	RencanaPembangunan      Rencana `gorm:"type:string[][]"`
}

type Lokasi struct {
	Lokasi []string
}

type Jenis struct {
	Jenis []string
}

type Rencana struct {
	Recana [][]string
}
