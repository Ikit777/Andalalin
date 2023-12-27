package models

import "github.com/google/uuid"

type DataMaster struct {
	IdDataMaster               uuid.UUID                  `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	JenisProyek                JenisProyek                `gorm:"serializer:json"`
	LokasiPengambilan          Lokasi                     `gorm:"serializer:json"`
	KategoriRencanaPembangunan KategoriRencanaPembangunan `gorm:"serializer:json"`
	JenisRencanaPembangunan    []JenisRencanaPembangunan  `gorm:"serializer:json"`
	KategoriPerlengkapanUtama  KategoriPerlengkapanUtama  `gorm:"serializer:json"`
	KategoriPerlengkapan       []KategoriPerlengkapan     `gorm:"serializer:json"`
	PerlengkapanLaluLintas     []JenisPerlengkapan        `gorm:"serializer:json"`
	Persyaratan                Persyaratan                `gorm:"serializer:json"`
	Provinsi                   []Provinsi                 `gorm:"serializer:json"`
	Kabupaten                  []Kabupaten                `gorm:"serializer:json"`
	Kecamatan                  []Kecamatan                `gorm:"serializer:json"`
	Kelurahan                  []Kelurahan                `gorm:"serializer:json"`
	Jalan                      []Jalan                    `gorm:"serializer:json"`
	UpdatedAt                  string                     `gorm:"not null"`
}

//Lokasi pengambilan
type Lokasi []string

type LokasiInput struct {
	Lokasi string `json:"lokasi" binding:"required"`
}

type LokasiEdit struct {
	Lokasi     string `json:"lokasi" binding:"required"`
	LokasiEdit string `json:"lokasi_edit" binding:"required"`
}

//Kategori rencana pembangunan
type KategoriRencanaPembangunan []string

type KategoriRencanaInput struct {
	Kategori string `json:"kategori" binding:"required"`
}

type KategoriRencanaEdit struct {
	Kategori     string `json:"kategori" binding:"required"`
	KategoriEdit string `json:"kategori_edit" binding:"required"`
}

//Jenis rencana pembangunan
type JenisRencanaPembangunan struct {
	Kategori     string
	JenisRencana []JenisRencana
}

type JenisRencana struct {
	Jenis     string
	Kriteria  string
	Satuan    string
	Terbilang string
}

type JenisRencanaPembangunanInput struct {
	Kategori  string `json:"kategori" binding:"required"`
	Jenis     string `json:"jenis" binding:"required"`
	Kriteria  string `json:"kriteria" binding:"required"`
	Satuan    string `json:"satuan" binding:"required"`
	Terbilang string `json:"terbilang" binding:"required"`
}

type JenisRencanaPembangunanHapus struct {
	Kategori string `json:"kategori" binding:"required"`
	Jenis    string `json:"jenis" binding:"required"`
}

type JenisRencanaPembangunanEdit struct {
	Kategori  string `json:"kategori" binding:"required"`
	Jenis     string `json:"jenis" binding:"required"`
	JenisEdit string `json:"jenis_edit" binding:"required"`
	Kriteria  string `json:"kriteria" binding:"required"`
	Satuan    string `json:"satuan" binding:"required"`
	Terbilang string `json:"terbilang" binding:"required"`
}

//Kategori perlengkapan utama
type KategoriPerlengkapanUtama []string

type KategoriPerlengkapanUtamaInput struct {
	Kategori string `json:"kategori" binding:"required"`
}

type KategoriPerlengkapanUtamaEdit struct {
	Kategori     string `json:"kategori" binding:"required"`
	KategoriEdit string `json:"kategori_edit" binding:"required"`
}

//Kategori perlengkapan
type KategoriPerlengkapan struct {
	KategoriUtama string
	Kategori      string
}

type KategoriPerlengkapanInput struct {
	KategoriUtama string `json:"kategori_utama" binding:"required"`
	Kategori      string `json:"kategori" binding:"required"`
}

type KategoriPerlengkapanEdit struct {
	KategoriUtama string `json:"kategori_utama" binding:"required"`
	Kategori      string `json:"kategori" binding:"required"`
	KategoriEdit  string `json:"kategori_edit" binding:"required"`
}

//Jenis perlengkapan lalu lintas
type JenisPerlengkapan struct {
	KategoriUtama string
	Kategori      string
	Perlengkapan  []PerlengkapanItem
}

type PerlengkapanItem struct {
	JenisPerlengkapan  string
	GambarPerlengkapan []byte
}

type JenisPerlengkapanInput struct {
	KategoriUtama string `json:"kategori_utama" binding:"required"`
	Kategori      string `json:"kategori" binding:"required"`
	Jenis         string `json:"jenis" binding:"required"`
}

type DataPerlengkapan struct {
	Perlengkapan JenisPerlengkapanInput `form:"data"`
}

type JenisPerlengkapanEdit struct {
	KategoriUtama string `json:"kategori_utama" binding:"required"`
	Kategori      string `json:"kategori" binding:"required"`
	Jenis         string `json:"jenis" binding:"required"`
	JenisEdit     string `json:"jenis_edit" binding:"required"`
}

type DataPerlengkapanEdit struct {
	Perlengkapan JenisPerlengkapanEdit `form:"data"`
}

//Persyaratan
type Persyaratan struct {
	PersyaratanAndalalin []PersyaratanAndalalinInput
	PersyaratanPerlalin  []PersyaratanPerlalinInput
}

//Persyaratan andalalin
type PersyaratanAndalalinInput struct {
	Kebutuhan             string `json:"kebutuhan" binding:"required"`
	Bangkitan             string `json:"bangkitan" binding:"required"`
	Persyaratan           string `json:"persyaratan" binding:"required"`
	KeteranganPersyaratan string `json:"keterangan" binding:"required"`
}

type PersyaratanAndalalinHapus struct {
	Persyaratan string `json:"persyaratan" binding:"required"`
}

type PersyaratanAndalalinEdit struct {
	Kebutuhan             string `json:"kebutuhan" binding:"required"`
	Bangkitan             string `json:"bangkitan" binding:"required"`
	Persyaratan           string `json:"persyaratan" binding:"required"`
	PersyaratanEdit       string `json:"persyaratan_edit" binding:"required"`
	KeteranganPersyaratan string `json:"keterangan" binding:"required"`
}

//Persyaratan perlalin
type PersyaratanPerlalinInput struct {
	Kebutuhan             string `json:"kebutuhan" binding:"required"`
	Persyaratan           string `json:"persyaratan" binding:"required"`
	KeteranganPersyaratan string `json:"keterangan" binding:"required"`
}

type PersyaratanPerlalinHapus struct {
	Persyaratan string `json:"persyaratan" binding:"required"`
}

type PersyaratanPerlalinEdit struct {
	Kebutuhan             string `json:"kebutuhan" binding:"required"`
	Persyaratan           string `json:"persyaratan" binding:"required"`
	PersyaratanEdit       string `json:"persyaratan_edit" binding:"required"`
	KeteranganPersyaratan string `json:"keterangan" binding:"required"`
}

//Provinsi
type Provinsi struct {
	Id   string
	Name string
}

type ProvinsiInput struct {
	Provinsi string `json:"provinsi" binding:"required"`
}

type ProvinsiEdit struct {
	Provinsi     string `json:"provinsi" binding:"required"`
	ProvinsiEdit string `json:"provinsi_edit" binding:"required"`
}

//Kabupaten
type Kabupaten struct {
	Id         string
	IdProvinsi string
	Name       string
}

type KabupatenInput struct {
	Provinsi  string `json:"provinsi" binding:"required"`
	Kabupaten string `json:"kabupaten" binding:"required"`
}

type KabupatenHapus struct {
	Kabupaten string `json:"kabupaten" binding:"required"`
}

type KabupatenEdit struct {
	Kabupaten     string `json:"kabupaten" binding:"required"`
	KabupatenEdit string `json:"kabupaten_edit" binding:"required"`
}

//Kecamatan
type Kecamatan struct {
	Id          string
	IdKabupaten string
	Name        string
}

type KecamatanInput struct {
	Kecamatan string `json:"kecamatan" binding:"required"`
	Kabupaten string `json:"kabupaten" binding:"required"`
}

type KecamatanHapus struct {
	Kecamatan string `json:"kecamatan" binding:"required"`
}

type KecamatanEdit struct {
	Kecamatan     string `json:"kecamatan" binding:"required"`
	KecamatanEdit string `json:"kecamatan_edit" binding:"required"`
}

//Kelurahan
type Kelurahan struct {
	Id          string
	IdKecamatan string
	Name        string
}

type KelurahanInput struct {
	Kecamatan string `json:"kecamatan" binding:"required"`
	Kelurahan string `json:"kelurahan" binding:"required"`
}

type KelurahanHapus struct {
	Kelurahan string `json:"kelurahan" binding:"required"`
}

type KelurahanEdit struct {
	Kelurahan     string `json:"kelurahan" binding:"required"`
	KelurahanEdit string `json:"kelurahan_edit" binding:"required"`
}

//Jenis proyek
type JenisProyek []string

type JenisProyekInput struct {
	Jenis string `json:"jenis" binding:"required"`
}

type JenisProyekEdit struct {
	Jenis     string `json:"jenis" binding:"required"`
	JenisEdit string `json:"jenis_edit" binding:"required"`
}

//Jalan

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

type JalanInput struct {
	KodeKecamatan string `json:"kode_kecamatan" binding:"required"`
	KodeKelurahan string `json:"kode_kelurahan" binding:"required"`
	KodeJalan     string `json:"kode_jalan" binding:"required"`
	Nama          string `json:"nama" binding:"required"`
	Pangkal       string `json:"pangkal" binding:"required"`
	Ujung         string `json:"ujung" binding:"required"`
	Kelurahan     string `json:"kelurahan" binding:"required"`
	Kecamatan     string `json:"kecamatan" binding:"required"`
	Panjang       string `json:"panjang" binding:"required"`
	Lebar         string `json:"lebar" binding:"required"`
	Permukaan     string `json:"permukaan" binding:"required"`
	Fungsi        string `json:"fungsi" binding:"required"`
}

type JalanHapus struct {
	Kode string `json:"kode" binding:"required"`
}

type JalanEdit struct {
	Kode          string `json:"kode" binding:"required"`
	Jalan         string `json:"jalan" binding:"required"`
	KodeKecamatan string `json:"kode_kecamatan" binding:"required"`
	KodeKelurahan string `json:"kode_kelurahan" binding:"required"`
	KodeJalan     string `json:"kode_jalan" binding:"required"`
	Nama          string `json:"nama" binding:"required"`
	Pangkal       string `json:"pangkal" binding:"required"`
	Ujung         string `json:"ujung" binding:"required"`
	Kelurahan     string `json:"kelurahan" binding:"required"`
	Kecamatan     string `json:"kecamatan" binding:"required"`
	Panjang       string `json:"panjang" binding:"required"`
	Lebar         string `json:"lebar" binding:"required"`
	Permukaan     string `json:"permukaan" binding:"required"`
	Fungsi        string `json:"fungsi" binding:"required"`
}
