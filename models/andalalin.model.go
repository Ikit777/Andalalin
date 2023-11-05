package models

import (
	"github.com/google/uuid"
)

type Andalalin struct {
	//Data Permohonan
	IdAndalalin       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	IdUser            uuid.UUID `gorm:"type:varchar(255);not null"`
	JenisAndalalin    string    `gorm:"type:varchar(255);not null"`
	WaktuAndalalin    string    `gorm:"not null"`
	TanggalAndalalin  string    `gorm:"not null"`
	StatusAndalalin   string    `sql:"type:enum('Cek persyaratan', 'Persyaratan tidak terpenuhi', 'Persyaratan terpenuhi', 'Berita acara pemeriksaan', 'Persetujuan dokumen', 'Pembuatan surat keputusan', 'Permohonan selesai', 'Permohonan ditolak')"`
	Bangkitan         string    `gorm:"type:varchar(255);not null"`
	Pemohon           string    `gorm:"type:varchar(255);not null"`
	Kategori          string    `gorm:"type:varchar(255);not null"`
	Jenis             string    `gorm:"type:varchar(255);not null"`
	Kode              string    `gorm:"type:varchar(255);not null"`
	LokasiPengambilan string    `gorm:"type:varchar(255);not null"`

	//Data pemohon
	NikPemohon                  string  `gorm:"type:varchar(255);not null"`
	NamaPemohon                 string  `gorm:"type:varchar(255);not null"`
	EmailPemohon                string  `gorm:"type:varchar(255);not null"`
	TempatLahirPemohon          string  `gorm:"type:varchar(255);not null"`
	TanggalLahirPemohon         string  `gorm:"type:varchar(255);not null"`
	WilayahAdministratifPemohon string  `gorm:"type:varchar(255);not null"`
	AlamatPemohon               string  `gorm:"type:varchar(255);not null"`
	JenisKelaminPemohon         string  `sql:"type:enum('Laki-laki', 'Perempuan');not null"`
	NomerPemohon                string  `gorm:"type:varchar(255);not null"`
	JabatanPemohon              *string `gorm:"type:varchar(255);"`
	NomerSertifikatPemohon      string  `gorm:"type:varchar(255);not null"`
	KlasifikasiPemohon          string  `gorm:"type:varchar(255);not null"`

	//Data Proyek
	NamaProyek                 string `gorm:"type:varchar(255);not null"`
	JenisProyek                string `gorm:"type:varchar(255);not null"`
	WilayahAdministratifProyek string `gorm:"type:varchar(255);not null"`
	AlamatProyek               string `gorm:"type:varchar(255);not null"`
	KodeJalan                  string `gorm:"type:varchar(255);not null"`
	NamaJalan                  string `gorm:"type:varchar(255);not null"`
	Kelurahan                  string `gorm:"type:varchar(255);not null"`
	Kecamatan                  string `gorm:"type:varchar(255);not null"`
	PangkalJalan               string `gorm:"type:varchar(255);not null"`
	UjungJalan                 string `gorm:"type:varchar(255);not null"`
	PanjangJalan               string `gorm:"type:varchar(255);not null"`
	LebarJalan                 string `gorm:"type:varchar(255);not null"`
	PermukaanJalan             string `gorm:"type:varchar(255);not null"`
	FungsiJalan                string `gorm:"type:varchar(255);not null"`

	//Data Perusahaan
	NamaPerusahaan                 *string `gorm:"type:varchar(255);"`
	WilayahAdministratifPerusahaan *string `gorm:"type:varchar(255);"`
	AlamatPerusahaan               *string `gorm:"type:varchar(255);"`
	NomerPerusahaan                *string `gorm:"type:varchar(255);"`
	EmailPerusahaan                *string `gorm:"type:varchar(255);"`
	NamaPimpinan                   *string `gorm:"type:varchar(255);"`
	JabatanPimpinan                *string `gorm:"type:varchar(255);"`
	JenisKelaminPimpinan           *string `sql:"type:enum('Laki-laki', 'Perempuan');"`
	WilayahAdministratifPimpinan   *string `gorm:"type:varchar(255);"`
	AlamatPimpinan                 *string `gorm:"type:varchar(255);"`

	//Data Pengembang
	NamaPengembang                         string `gorm:"type:varchar(255);not null"`
	WilayahAdministratifPengembang         string `gorm:"type:varchar(255);not null"`
	AlamatPengembang                       string `gorm:"type:varchar(255);not null"`
	NomerPengembang                        string `gorm:"type:varchar(255);not null"`
	EmailPengembang                        string `gorm:"type:varchar(255);not null"`
	NamaPimpinanPengembang                 string `gorm:"type:varchar(255);not null"`
	JabatanPimpinanPengembang              string `gorm:"type:varchar(255);not null"`
	JenisKelaminPimpinanPengembang         string `sql:"type:enum('Laki-laki', 'Perempuan');not null"`
	WilayahAdministratifPimpinanPengembang string `gorm:"type:varchar(255);not null"`
	AlamatPimpinanPengembang               string `gorm:"type:varchar(255);not null"`

	//Data Kegiatan
	Aktivitas         string  `gorm:"type:varchar(255);not null"`
	Peruntukan        string  `gorm:"type:varchar(255);not null"`
	KriteriaKhusus    *string `gorm:"type:varchar(255);"`
	NilaiKriteria     *string `gorm:"type:varchar(255);"`
	LokasiBangunan    string  `gorm:"type:varchar(255);not null"`
	LatitudeBangunan  float64
	LongitudeBangunan float64
	NomerSKRK         string `gorm:"type:varchar(255);not null"`
	TanggalSKRK       string `gorm:"type:varchar(255);not null"`
	Catatan           *string

	//Data Persyaratan
	Persyaratan []PersyaratanPermohonan `gorm:"serializer:json"`

	//Dokumen Permohonan
	Dokumen []DokumenPermohonan `gorm:"serializer:json"`

	//Petimbangan
	Pertimbangan string

	//Persyaratan tidak terpenuhi
	PersyaratanTidakSesuai []string `gorm:"serializer:json"`

	//Data Persetujuan Dokumen
	PersetujuanDokumen           string `gorm:"type:varchar(255);"`
	KeteranganPersetujuanDokumen *string

	//Data BAP
	NomerBAPDasar       string `gorm:"type:varchar(255);"`
	NomerBAPPelaksanaan string `gorm:"type:varchar(255);"`
	TanggalBAP          string `gorm:"type:varchar(255);"`
	FileBAP             []byte

	//Data SK
	FileSK []byte
}

type Perlalin struct {
	//Data Permohonan
	IdAndalalin      uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	JenisAndalalin   string    `gorm:"type:varchar(255);not null"`
	Kategori         string    `gorm:"type:varchar(255);not null"`
	Jenis            string    `gorm:"type:varchar(255);not null"`
	Kode             string    `gorm:"type:varchar(255);not null"`
	WaktuAndalalin   string    `gorm:"not null"`
	TanggalAndalalin string    `gorm:"not null"`
	StatusAndalalin  string    `sql:"type:enum('Cek persyaratan', 'Persyaratan tidak terpenuhi', 'Persyaratan terpenuhi', 'Survei lapangan', 'Laporan survei', 'Menunggu hasil keputusan', 'Tunda pemasangan', 'Pemasangan sedang dilakukan', 'Permohonan ditolak', 'Permohonan ditunda', 'Permohonan dibatalkan', 'Pemasangan selesai')"`

	//Data Pemohon
	IdUser                      uuid.UUID `gorm:"type:varchar(255);not null"`
	NikPemohon                  string    `gorm:"type:varchar(255);not null"`
	NamaPemohon                 string    `gorm:"type:varchar(255);not null"`
	EmailPemohon                string    `gorm:"type:varchar(255);not null"`
	TempatLahirPemohon          string    `gorm:"type:varchar(255);not null"`
	TanggalLahirPemohon         string    `gorm:"type:varchar(255);not null"`
	WilayahAdministratifPemohon string    `gorm:"type:varchar(255);not null"`
	AlamatPemohon               string    `gorm:"type:varchar(255);not null"`
	JenisKelaminPemohon         string    `sql:"type:enum('Laki-laki', 'Perempuan');not null"`
	NomerPemohon                string    `gorm:"type:varchar(255);not null"`

	//Data Kegiatan
	Alasan              string `gorm:"type:varchar(255);not null"`
	Peruntukan          string `gorm:"type:varchar(255);not null"`
	LokasiPemasangan    string `gorm:"type:varchar(255);not null"`
	LatitudePemasangan  float64
	LongitudePemasangan float64
	LokasiPengambilan   string `gorm:"type:varchar(255);not null"`
	Catatan             *string

	//Data Petugas
	IdPetugas    uuid.UUID `gorm:"type:varchar(255);"`
	NamaPetugas  string    `gorm:"type:varchar(255);"`
	EmailPetugas string    `gorm:"type:varchar(255);"`

	//Tanda Terima Pendaftaran
	TandaTerimaPendaftaran []byte

	//Data Persyaratan
	Persyaratan []PersyaratanPermohonan `gorm:"serializer:json"`

	//Persyaratan tidak terpenuhi
	PersyaratanTidakSesuai []string `gorm:"serializer:json"`

	//Data Laporan Survei
	LaporanSurvei []byte

	//Data Tindakan
	Tindakan string

	//Data Pertimbangan
	PertimbanganTindakan  string
	PertimbanganPenolakan string
	PertimbanganPenundaan string
}

type PersyaratanPermohonan struct {
	Persyaratan string
	Berkas      []byte
}

type DokumenPermohonan struct {
	Dokumen string
	Berkas  []byte
}

type InputAndalalin struct {
	//Data Permohonan
	Bangkitan                       string `json:"kategori_bangkitan" binding:"required"`
	Pemohon                         string `json:"kategori_pemohon" binding:"required"`
	KategoriJenisRencanaPembangunan string `json:"kategori" binding:"required"`
	JenisRencanaPembangunan         string `json:"jenis_rencana_pembangunan" binding:"required"`
	LokasiPengambilan               string `json:"lokasi_pengambilan" binding:"required"`

	//Data Proyek
	NamaProyek                 string `json:"nama_proyek" binding:"required"`
	JenisProyek                string `json:"jenis_proyek" binding:"required"`
	WilayahAdministratifProyek string `json:"wilayah_administratif_proyek" binding:"required"`
	AlamatProyek               string `json:"alamat_proyek" binding:"required"`
	KodeJalan                  string `json:"kode_jalan" binding:"required"`
	NamaJalan                  string `json:"nama_jalan" binding:"required"`
	Kelurahan                  string `json:"kelurahan" binding:"required"`
	Kecamatan                  string `json:"kecamatan" binding:"required"`
	PangkalJalan               string `json:"pangkal_jalan" binding:"required"`
	UjungJalan                 string `json:"ujung_jalan" binding:"required"`
	PanjangJalan               string `json:"panjang_jalan" binding:"required"`
	LebarJalan                 string `json:"lebar_jalan" binding:"required"`
	PermukaanJalan             string `json:"permukaan_jalan" binding:"required"`
	FungsiJalan                string `json:"fungsi_jalan" binding:"required"`

	//Data Pemohon
	NikPemohon                  string  `json:"nik_pemohon" binding:"required"`
	TempatLahirPemohon          string  `json:"tempat_lahir_pemohon" binding:"required"`
	TanggalLahirPemohon         string  `json:"tanggal_lahir_pemohon" binding:"required"`
	AlamatPemohon               string  `json:"alamat_pemohon" binding:"required"`
	WilayahAdministratifPemohon string  `json:"wilayah_administratif_pemohon" binding:"required"`
	JenisKelaminPemohon         string  `json:"jenis_kelamin_pemohon" binding:"required"`
	NomerPemohon                string  `json:"nomer_pemohon" binding:"required"`
	JabatanPemohon              *string `json:"jabatan_pemohon" binding:"required"`
	NomerSertifikatPemohon      string  `json:"nomer_sertifikat_pemohon" binding:"required"`
	KlasifikasiPemohon          string  `json:"klasifikasi_pemohon" binding:"required"`

	//Data Perusahaan
	NamaPerusahaan                 *string `json:"nama_perusahaan" binding:"required"`
	AlamatPerusahaan               *string `json:"alamat_perusahaan" binding:"required"`
	WilayahAdministratifPerusahaan *string `json:"wilayah_administratif_perusahaan" binding:"required"`
	NomerPerusahaan                *string `json:"nomer_perusahaan" binding:"required"`
	EmailPerusahaan                *string `json:"email_perusahaan" binding:"required"`
	NamaPimpinan                   *string `json:"nama_pimpinan" binding:"required"`
	JabatanPimpinan                *string `json:"jabatan_pimpinan" binding:"required"`
	JenisKelaminPimpinan           *string `json:"jenis_kelamin_pimpinan" binding:"required"`
	WilayahAdministratifPimpinan   *string `json:"wilayah_administratif_pimpinan" binding:"required"`
	AlamatPimpinan                 *string `json:"alamat_pimpinan" binding:"required"`

	//Data Pengembang
	NamaPengembang                         string `json:"nama_pengembang" binding:"required"`
	WilayahAdministratifPengembang         string `json:"wilayah_administratif_pengembang" binding:"required"`
	AlamatPengembang                       string `json:"alamat_pengembang" binding:"required"`
	NomerPengembang                        string `json:"nomer_pengembang" binding:"required"`
	EmailPengembang                        string `json:"email_pengembang" binding:"required"`
	NamaPimpinanPengembang                 string `json:"nama_pimpinan_pengembang" binding:"required"`
	JabatanPimpinanPengembang              string `json:"jabatan_pimpinan_pengembang" binding:"required"`
	JenisKelaminPimpinanPengembang         string `json:"jenis_kelamin_pimpinan_pengembang" binding:"required"`
	WilayahAdministratifPimpinanPengembang string `json:"wilayah_administratif_pimpinan_pengembang" binding:"required"`
	AlamatPimpinanPengembang               string `json:"alamat_pimpinan_pengembang" binding:"required"`

	//Data Kegiatan
	Aktivitas         string  `json:"aktivitas" binding:"required"`
	Peruntukan        string  `json:"peruntukan" binding:"required"`
	KriteriaKhusus    *string `json:"kriteria_khusus" binding:"required"`
	NilaiKriteria     *string `json:"nilai_kriteria" binding:"required"`
	LokasiBangunan    string  `json:"lokasi_bangunan" binding:"required"`
	LatitudeBangunan  float64 `protobuf:"fixed64,1,opt,name=latitude,proto3" json:"latitude" binding:"required"`
	LongitudeBangunan float64 `protobuf:"fixed64,2,opt,name=longitude,proto3" json:"longtitude" binding:"required"`
	NomerSKRK         string  `json:"nomer_skrk" binding:"required"`
	TanggalSKRK       string  `json:"tanggal_skrk" binding:"required"`
	Catatan           *string `json:"catatan" binding:"required"`
}

type InputPerlalin struct {
	Kategori                    string  `json:"kategori" binding:"required"`
	Jenis                       string  `json:"jenis_perlengkapan" binding:"required"`
	NikPemohon                  string  `json:"nik_pemohon" binding:"required"`
	TempatLahirPemohon          string  `json:"tempat_lahir_pemohon" binding:"required"`
	TanggalLahirPemohon         string  `json:"tanggal_lahir_pemohon" binding:"required"`
	WilayahAdministratifPemohon string  `json:"wilayah_administratif_pemohon" binding:"required"`
	AlamatPemohon               string  `json:"alamat_pemohon" binding:"required"`
	JenisKelaminPemohon         string  `json:"jenis_kelamin_pemohon" binding:"required"`
	NomerPemohon                string  `json:"nomer_pemohon" binding:"required"`
	LokasiPengambilan           string  `json:"lokasi_pengambilan" binding:"required"`
	Alasan                      string  `json:"alasan" binding:"required"`
	Peruntukan                  string  `json:"peruntukan" binding:"required"`
	LokasiPemasangan            string  `json:"lokasi_pemasangan" binding:"required"`
	LatitudePemasangan          float64 `protobuf:"fixed64,1,opt,name=latitude,proto3" json:"latitude" binding:"required"`
	LongitudePemasangan         float64 `protobuf:"fixed64,2,opt,name=longitude,proto3" json:"longtitude" binding:"required"`
	Catatan                     *string `json:"catatan" binding:"required"`
}

type DataAndalalin struct {
	Andalalin InputAndalalin `form:"data"`
}

type DataPerlalin struct {
	Perlalin InputPerlalin `form:"data"`
}

type DaftarAndalalinResponse struct {
	IdAndalalin      uuid.UUID `json:"id_andalalin,omitempty"`
	Kode             string    `json:"kode_andalalin,omitempty"`
	TanggalAndalalin string    `json:"tanggal_andalalin,omitempty"`
	Nama             string    `json:"nama_pemohon,omitempty"`
	Alamat           string    `json:"alamat_pemohon,omitempty"`
	JenisAndalalin   string    `json:"jenis_andalalin,omitempty"`
	StatusAndalalin  string    `json:"status_andalalin,omitempty"`
}

type AndalalinResponse struct {
	//Data Permohonan
	IdAndalalin       uuid.UUID `json:"id_andalalin,omitempty"`
	Bangkitan         string    `json:"kategori_bangkitan,omitempty"`
	Pemohon           string    `json:"kategori_pemohon,omitempty"`
	JenisAndalalin    string    `json:"jenis_andalalin,omitempty"`
	Kategori          string    `json:"kategori,omitempty"`
	Jenis             string    `json:"jenis_rencana_pembangunan,omitempty"`
	Kode              string    `json:"kode_andalalin,omitempty"`
	WaktuAndalalin    string    `json:"waktu_andalalin,omitempty"`
	TanggalAndalalin  string    `json:"tanggal_andalalin,omitempty"`
	StatusAndalalin   string    `json:"status_andalalin,omitempty"`
	LokasiPengambilan string    `json:"lokasi_pengambilan,omitempty"`

	//Data Proyek
	NamaProyek                 string `json:"nama_proyek,omitempty"`
	JenisProyek                string `json:"jenis_proyek,omitempty"`
	WilayahAdministratifProyek string `json:"wilayah_administratif_proyek,omitempty"`
	AlamatProyek               string `json:"alamat_proyek,omitempty"`
	KodeJalan                  string `json:"kode_jalan,omitempty"`
	NamaJalan                  string `json:"nama_jalan,omitempty"`
	Kelurahan                  string `json:"kelurahan_jalan,omitempty"`
	Kecamatan                  string `json:"kecamatan_jalan,omitempty"`
	PangkalJalan               string `json:"pangkal_jalan,omitempty"`
	UjungJalan                 string `json:"ujung_jalan,omitempty"`
	PanjangJalan               string `json:"panjang_jalan,omitempty"`
	LebarJalan                 string `json:"lber_jalan,omitempty"`
	PermukaanJalan             string `json:"permukaan_jalan,omitempty"`
	FungsiJalan                string `json:"fungsi_jalan,omitempty"`

	//Data Pemohon
	NikPemohon                  string  `json:"nik_pemohon,omitempty"`
	NamaPemohon                 string  `json:"nama_pemohon,omitempty"`
	EmailPemohon                string  `json:"email_pemohon,omitempty"`
	TempatLahirPemohon          string  `json:"tempat_lahir_pemohon,omitempty"`
	TanggalLahirPemohon         string  `json:"tanggal_lahir_pemohon,omitempty"`
	WilayahAdministratifPemohon string  `json:"wilayah_administratif_pemohon,omitempty"`
	AlamatPemohon               string  `json:"alamat_pemohon,omitempty"`
	JenisKelaminPemohon         string  `json:"jenis_kelamin_pemohon,omitempty"`
	NomerPemohon                string  `json:"nomer_pemohon,omitempty"`
	JabatanPemohon              *string `json:"jabatan_pemohon,omitempty"`
	NomerSertifikatPemohon      string  `json:"nomer_sertifikat_pemohon,omitempty"`
	KlasifikasiPemohon          string  `json:"klasifikasi_pemohon,omitempty"`

	//Data Perusahaan
	NamaPerusahaan                 *string `json:"nama_perusahaan,omitempty"`
	WilayahAdministratifPerusahaan *string `json:"wilayah_administratif_perusahaan,omitempty"`
	AlamatPerusahaan               *string `json:"alamat_perusahaan,omitempty"`
	NomerPerusahaan                *string `json:"nomer_perusahaan,omitempty"`
	EmailPerusahaan                *string `json:"email_perusahaan,omitempty"`
	NamaPimpinan                   *string `json:"nama_pimpinan,omitempty"`
	JabatanPimpinan                *string `json:"jabatan_pimpinan,omitempty"`
	JenisKelaminPimpinan           *string `json:"jenis_kelamin,omitempty"`
	WilayahAdministratifPimpinan   *string `json:"wilayah_administratif_pimpinan,omitempty"`
	AlamatPimpinan                 *string `json:"alamat_pimpinan,omitempty"`

	//Data Pengembang
	NamaPengembang                         string `json:"nama_pengembang,omitempty"`
	WilayahAdministratifPengembang         string `json:"wilayah_administratif_pengembang,omitempty"`
	AlamatPengembang                       string `json:"alamat_pengembang,omitempty"`
	NomerPengembang                        string `json:"nomer_pengembang,omitempty"`
	EmailPengembang                        string `json:"email_pengembang,omitempty"`
	NamaPimpinanPengembang                 string `json:"nama_pimpinan_pengembang,omitempty"`
	JabatanPimpinanPengembang              string `json:"jabatan_pimpinan_pengembang,omitempty"`
	JenisKelaminPimpinanPengembang         string `json:"jenis_kelamin_pimpinan_pengembang,omitempty"`
	WilayahAdministratifPimpinanPengembang string `json:"wilayah_administratif_pimpinan_pengembang,omitempty"`
	AlamatPimpinanPengembang               string `json:"alamat_pimpinan_pengembang,omitempty"`

	//Data Kegiatan
	Aktivitas         string  `json:"aktivitas,omitempty"`
	Peruntukan        string  `json:"peruntukan,omitempty"`
	KriteriaKhusus    *string `json:"kriteria_khusus,omitempty"`
	NilaiKriteria     *string `json:"nilai_kriteria,omitempty"`
	LatitudeBangunan  float64 `json:"latitude,omitempty"`
	LongitudeBangunan float64 `json:"longitude,omitempty"`
	LokasiBangunan    string  `json:"alamat_persil,omitempty"`
	NomerSKRK         string  `json:"nomer_skrk,omitempty"`
	TanggalSKRK       string  `json:"tanggal_skrk,omitempty"`
	Catatan           *string `json:"catatan,omitempty"`

	//Persyaratan tidak terpenuhi
	PersyaratanTidakSesuai []string `json:"persyaratan_tidak_sesuai,omitempty"`

	Pertimbangan string `json:"pertimbangan,omitempty"`

	//Data Persetujuan
	PersetujuanDokumen           string  `json:"persetujuan,omitempty"`
	KeteranganPersetujuanDokumen *string `json:"keterangan_persetujuan,omitempty"`

	//Data BAP
	NomerBAPDasar       string `json:"nomer_bap_dasar,omitempty"`
	NomerBAPPelaksanaan string `json:"nomer_bap_pelaksanaan,omitempty"`
	TanggalBAP          string `json:"tanggal_bap,omitempty"`

	//Persyaratan Permohonan
	Persyaratan []string `json:"persyaratan,omitempty"`

	//Dokumen Permohonan
	Dokumen []string `json:"dokumen,omitempty"`
}

type PerlalinResponse struct {
	//Data Permohonan
	IdAndalalin      uuid.UUID `json:"id_andalalin,omitempty"`
	JenisAndalalin   string    `json:"jenis_andalalin,omitempty"`
	Kategori         string    `json:"kategori,omitempty"`
	Jenis            string    `json:"jenis,omitempty"`
	Kode             string    `json:"kode_andalalin,omitempty"`
	WaktuAndalalin   string    `json:"waktu_andalalin,omitempty"`
	TanggalAndalalin string    `json:"tanggal_andalalin,omitempty"`
	StatusAndalalin  string    `json:"status_andalalin,omitempty"`

	//Data Pemohon
	NikPemohon                  string `json:"nik_pemohon,omitempty"`
	NamaPemohon                 string `json:"nama_pemohon,omitempty"`
	EmailPemohon                string `json:"email_pemohon,omitempty"`
	TempatLahirPemohon          string `json:"tempat_lahir_pemohon,omitempty"`
	TanggalLahirPemohon         string `json:"tanggal_lahir_pemohon,omitempty"`
	WilayahAdministratifPemohon string `json:"wilayah_administratif_pemohon,omitempty"`
	AlamatPemohon               string `json:"alamat_pemohon,omitempty"`
	JenisKelaminPemohon         string `json:"jenis_kelamin_pemohon,omitempty"`
	NomerPemohon                string `json:"nomer_pemohon,omitempty"`
	LokasiPengambilan           string `json:"lokasi_pengambilan,omitempty"`

	//Data Kegiatan
	Alasan              string  `json:"alasan,omitempty"`
	Peruntukan          string  `json:"peruntukan,omitempty"`
	LokasiPemasangan    string  `json:"lokasi_pemasangan,omitempty"`
	LatitudePemasangan  float64 `json:"latitude,omitempty"`
	LongitudePemasangan float64 `json:"longitude,omitempty"`
	Catatan             *string `json:"catatan,omitempty"`

	//Persyaratan tidak terpenuhi
	PersyaratanTidakSesuai []string `json:"persyaratan_tidak_sesuai,omitempty"`

	//Data Petugas
	IdPetugas         uuid.UUID `json:"id_petugas,omitempty"`
	NamaPetugas       string    `json:"nama_petugas,omitempty"`
	EmailPetugas      string    `json:"email_petugas,omitempty"`
	StatusTiketLevel2 string    `json:"status_tiket,omitempty"`

	//Data Persyaratan
	Persyaratan []string `json:"persyaratan,omitempty"`

	//Data Tindakan
	Tindakan string `json:"keputusan_hasil,omitempty"`

	//Data Pertimbangan
	PertimbanganTindakan  string `json:"pertimbangan_tindakan,omitempty"`
	PertimbanganPenolakan string `json:"pertimbangan_penolakan,omitempty"`
	PertimbanganPenundaan string `json:"pertimbangan_penundaan,omitempty"`
}

type AndalalinResponseUser struct {
	//Data Pemohon
	IdAndalalin             uuid.UUID `json:"id_andalalin,omitempty"`
	JenisAndalalin          string    `json:"jenis_andalalin,omitempty"`
	Bangkitan               string    `json:"kategori_bangkitan,omitempty"`
	Pemohon                 string    `json:"kategori_pemohon,omitempty"`
	Kode                    string    `json:"kode_andalalin,omitempty"`
	WaktuAndalalin          string    `json:"waktu_andalalin,omitempty"`
	TanggalAndalalin        string    `json:"tanggal_andalalin,omitempty"`
	StatusAndalalin         string    `json:"status_andalalin,omitempty"`
	JenisRencanaPembangunan string    `json:"jenis_rencana_pembangunan,omitempty"`
	Kategori                string    `json:"kategori,omitempty"`
	LokasiPengambilan       string    `json:"lokasi_pengambilan,omitempty"`

	//Data Pemohon
	NikPemohon             string  `json:"nik_pemohon,omitempty"`
	NamaPemohon            string  `json:"nama_pemohon,omitempty"`
	JabatanPemohon         *string `json:"jabatan_pemohon,omitempty"`
	EmailPemohon           string  `json:"email_pemohon,omitempty"`
	NomerPemohon           string  `json:"nomer_pemohon,omitempty"`
	NomerSertifikatPemohon string  `json:"nomer_sertifikat_pemohon,omitempty"`
	KlasifikasiPemohon     string  `json:"klasifikasi_pemohon,omitempty"`

	//Data Proyek
	NamaProyek                 string `json:"nama_proyek,omitempty"`
	JenisProyek                string `json:"jenis_proyek,omitempty"`
	NamaJalan                  string `json:"nama_jalan,omitempty"`
	FungsiJalan                string `json:"fungsi_jalan,omitempty"`
	WilayahAdministratifProyek string `json:"wilayah_administratif_proyek,omitempty"`

	//Data Perusahaan
	NamaPerusahaan *string `json:"nama_perusahaan,omitempty"`

	//Data Pengembang
	NamaPengembang string `json:"nama_pengembang,omitempty"`

	//Data Kegiatan
	Aktivitas         string  `json:"aktivitas,omitempty"`
	Peruntukan        string  `json:"peruntukan,omitempty"`
	KriteriaKhusus    *string `json:"kriteria_khusus,omitempty"`
	NilaiKriteria     *string `json:"nilai_kriteria,omitempty"`
	LatitudeBangunan  float64 `json:"latitude,omitempty"`
	LongitudeBangunan float64 `json:"longitude,omitempty"`
	LokasiBangunan    string  `json:"alamat_persil,omitempty"`
	Catatan           *string `json:"catatan,omitempty"`

	//Persyaratan tidak terpenuhi
	PersyaratanTidakSesuai []string `json:"persyaratan_tidak_sesuai,omitempty"`

	Pertimbangan string `json:"pertimbangan,omitempty"`

	//Dokumen Permohonan
	Dokumen []string `json:"dokumen,omitempty"`
}

type PerlalinResponseUser struct {
	//Data Permohonan
	IdAndalalin       uuid.UUID `json:"id_andalalin,omitempty"`
	JenisAndalalin    string    `json:"jenis_andalalin,omitempty"`
	Kode              string    `json:"kode_andalalin,omitempty"`
	LokasiPengambilan string    `json:"lokasi_pengambilan,omitempty"`
	WaktuAndalalin    string    `json:"waktu_andalalin,omitempty"`
	TanggalAndalalin  string    `json:"tanggal_andalalin,omitempty"`
	StatusAndalalin   string    `json:"status_andalalin,omitempty"`
	Jenis             string    `json:"jenis,omitempty"`
	Kategori          string    `json:"kategori,omitempty"`

	//Data Pemohon
	NikPemohon   string `json:"nik_pemohon,omitempty"`
	NamaPemohon  string `json:"nama_pemohon,omitempty"`
	EmailPemohon string `json:"email_pemohon,omitempty"`
	NomerPemohon string `json:"nomer_pemohon,omitempty"`

	//Data Perusahaan
	Alasan              string  `json:"alasan,omitempty"`
	Peruntukan          string  `json:"peruntukan,omitempty"`
	LokasiPemasangan    string  `json:"lokasi_pemasangan,omitempty"`
	LatitudePemasangan  float64 `json:"latitude,omitempty"`
	LongitudePemasangan float64 `json:"longitude,omitempty"`
	Catatan             *string `json:"catatan,omitempty"`

	//Persyaratan tidak terpenuhi
	PersyaratanTidakSesuai []string `json:"persyaratan_tidak_sesuai,omitempty"`

	//Data Tindakan
	Tindakan string `json:"keputusan_hasil,omitempty"`

	//Data Pertimbangan
	PertimbanganTindakan  string `json:"pertimbangan_tindakan,omitempty"`
	PertimbanganPenolakan string `json:"pertimbangan_penolakan,omitempty"`
	PertimbanganPenundaan string `json:"pertimbangan_penundaan,omitempty"`
}

type PersayaratanTidakSesuaiInput struct {
	Persyaratan []string `json:"persyaratan" binding:"required"`
}

type Persetujuan struct {
	Persetujuan string  `json:"persetujuan" binding:"required"`
	Keterangan  *string `json:"keterangan" binding:"required"`
}

type KeputusanHasil struct {
	Keputusan    string `json:"keputusan" binding:"required"`
	Pertimbangan string `json:"pertimbangan" binding:"required"`
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

type Survei struct {
	IdSurvey      uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	IdAndalalin   uuid.UUID `gorm:"type:varchar(255);uniqueIndex;not null"`
	IdTiketLevel1 uuid.UUID `gorm:"type:varchar(255);not null"`
	IdTiketLevel2 uuid.UUID `gorm:"type:varchar(255);not null"`
	IdPetugas     uuid.UUID `gorm:"type:varchar(255);not null"`
	Petugas       string    `gorm:"type:varchar(255);not null"`
	EmailPetugas  string    `gorm:"type:varchar(255);not null"`
	Keterangan    *string
	Foto1         []byte
	Foto2         []byte
	Foto3         []byte
	Lokasi        string
	Latitude      float64
	Longitude     float64
	WaktuSurvei   string `gorm:"not null"`
	TanggalSurvei string `gorm:"not null"`
}

type SurveiMandiri struct {
	IdSurvey           uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	IdPetugas          uuid.UUID `gorm:"type:varchar(255);not null"`
	Petugas            string    `gorm:"type:varchar(255);not null"`
	EmailPetugas       string    `gorm:"type:varchar(255);not null"`
	Keterangan         *string
	Foto1              []byte
	Foto2              []byte
	Foto3              []byte
	Lokasi             string
	Latitude           float64
	Longitude          float64
	WaktuSurvei        string `gorm:"not null"`
	TanggalSurvei      string `gorm:"not null"`
	StatusSurvei       string
	KeteranganTindakan string
}

type InputSurvey struct {
	Lokasi     string  `json:"lokasi" binding:"required"`
	Keterangan *string `json:"keterangan" binding:"required"`
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
	IdPetugas     uuid.UUID `gorm:"type:varchar(255);not null"`
	Status        string    `sql:"type:enum('Buka', 'Tutup', 'Tunda', 'Batal');not null"`
}

type UsulanPengelolaan struct {
	IdUsulan                   uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	IdAndalalin                uuid.UUID `gorm:"type:varchar(255);uniqueIndex;not null"`
	IdTiketLevel1              uuid.UUID `gorm:"type:varchar(255);not null"`
	IdTiketLevel2              uuid.UUID `gorm:"type:varchar(255);not null"`
	IdPengusulTindakan         uuid.UUID `gorm:"type:varchar(255);not null"`
	NamaPengusulTindakan       string    `gorm:"type:varchar(255);not null"`
	PertimbanganUsulanTindakan string    `gorm:"type:varchar(255);not null"`
	KeteranganUsulanTindakan   *string   `gorm:"type:varchar(255);not null"`
}

type InputUsulanPengelolaan struct {
	PertimbanganUsulanTindakan string  `json:"pertimbangan" binding:"required"`
	KeteranganUsulanTindakan   *string `json:"keterangan" binding:"required"`
}

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

type Pemasangan struct {
	IdPemasangan      uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	IdAndalalin       uuid.UUID `gorm:"type:varchar(255);uniqueIndex;not null"`
	IdTiketLevel1     uuid.UUID `gorm:"type:varchar(255);not null"`
	IdPetugas         uuid.UUID `gorm:"type:varchar(255);not null"`
	Petugas           string    `gorm:"type:varchar(255);not null"`
	EmailPetugas      string    `gorm:"type:varchar(255);not null"`
	Keterangan        *string
	Foto1             []byte
	Foto2             []byte
	Foto3             []byte
	Lokasi            string
	Latitude          float64
	Longitude         float64
	WaktuPemasangan   string `gorm:"not null"`
	TanggalPemasangan string `gorm:"not null"`
}
