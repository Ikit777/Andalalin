package models

type DataMaster struct {
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
	PersyaratanTambahanAndalalin  []PersyaratanTambahanInput
	PersyaratanTambahanRambulalin []PersyaratanTambahanInput
}

type PersyaratanTambahanInput struct {
	PersyaratanTambahan           string
	KeteranganPersyaratanTambahan string
}
