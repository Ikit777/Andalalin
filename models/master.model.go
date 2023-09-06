package models

type DataMaster struct {
	LokasiPengambilan       Lokasi              `gorm:"serializer:json"`
	JenisRencanaPembangunan Jenis               `gorm:"serializer:json"`
	RencanaPembangunan      []Rencana           `gorm:"serializer:json"`
	PersyaratanTambahan     PersyaratanTambahan `gorm:"serializer:json"`
}

type Rencana struct {
	Kategori     string
	JenisRencana []string
}

type Lokasi []string

type Jenis []string

type PersyaratanTambahan struct {
	PersyaratanTambahanAndalalin  []string
	PersyaratanTambahanRambulalin []string
}

type PersyaratanTambahanInput struct {
	PersyaratanTambahan           string
	KeteranganPersyaratanTambahan string
}
