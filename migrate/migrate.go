package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Ikit777/E-Andalalin/initializers"
	"github.com/Ikit777/E-Andalalin/models"
	"github.com/Ikit777/E-Andalalin/utils"

	_ "time/tzdata"
)

func init() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load environment variables", err)
	}

	initializers.ConnectDB(&config)
}

func main() {
	initializers.DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")

	initializers.DB.Migrator().DropTable(&models.User{})
	initializers.DB.Migrator().DropTable(&models.Andalalin{})
	initializers.DB.Migrator().DropTable(&models.Survey{})
	initializers.DB.Migrator().DropTable(&models.TiketLevel1{})
	initializers.DB.Migrator().DropTable(&models.TiketLevel2{})
	initializers.DB.Migrator().DropTable(&models.Notifikasi{})
	initializers.DB.Migrator().DropTable(&models.DataMaster{})

	initializers.DB.AutoMigrate(&models.User{})
	initializers.DB.AutoMigrate(&models.Andalalin{})
	initializers.DB.AutoMigrate(&models.Survey{})
	initializers.DB.AutoMigrate(&models.TiketLevel1{})
	initializers.DB.AutoMigrate(&models.TiketLevel2{})
	initializers.DB.AutoMigrate(&models.Notifikasi{})
	initializers.DB.AutoMigrate(&models.DataMaster{})

	loc, _ := time.LoadLocation("Asia/Singapore")
	now := time.Now().In(loc).Format("02-01-2006")
	hashedPassword, err := utils.HashPassword("superadmin")
	if err != nil {
		return
	}

	filePath := "assets/default.png"
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal("Error reading the file:", err)
	}

	initializers.DB.Create(&models.User{
		Name:      "super admin",
		Email:     strings.ToLower("superadmin@gmail.com"),
		Password:  hashedPassword,
		Role:      "Super Admin",
		Photo:     fileData,
		Verified:  true,
		CreatedAt: now,
		UpdatedAt: now,
	})

	lokasi := []string{"Banjarmasin"}

	jenis_kegiatan := []string{"Pusat kegiatan", "Pemukiman", "Infrastruktur", "Lainnya"}

	pusat_kegiatan := []string{"Pusat perbelanjaan atau retail", "Perkantoran", "Industri dan pergudangan", "Sekolah atau universitas",
		"Lembaga kursus", "Rumah sakit", "Klinik bersama", "Bank", "Stasiun pengisin bahan bakar", "Hotel", "Gedung pertemuan",
		"Restoran", "Fasilitan olah raga", "Bengkel kendaraan bermotor", "Pencucian mobil"}
	infrastruktur := []string{"Akses ke dan dari jalan tol", "Pelabuhan", "Bandar udara", "Terminal", "Stasiun kereta api", "Pool kendaraan", "Fasilitas parkir umum", "Flyover", "Underpass", "Terowongan"}
	pemukiman := []string{"Perumahan sederhana", "Perumahan menengan-atas", "Rumah susun sederhana", "Apartemen", "Asrama", "Ruko"}

	rencana := [][]string{pusat_kegiatan, infrastruktur, pemukiman}

	initializers.DB.Create(&models.DataMaster{
		LokasiPengambilan:       models.Lokasi{Lokasi: lokasi},
		JenisRencanaPembangunan: models.Jenis{Jenis: jenis_kegiatan},
		RencanaPembangunan:      models.Rencana{Recana: rencana},
	})

	fmt.Println("Migration complete")
}
