package utils

import (
	"crypto/rand"
	"io"
	"time"
)

var number = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

func Encode(s int) string {
	b := make([]byte, s)
	n, err := io.ReadAtLeast(rand.Reader, b, s)
	if n != s {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = number[int(b[i])%len(number)]
	}
	return string(b)
}

var code = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z'}

func Generate(s int) string {
	b := make([]byte, s)
	n, err := io.ReadAtLeast(rand.Reader, b, s)
	if n != s {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = code[int(b[i])%len(code)]
	}
	return string(b)
}

func Bulan(month time.Month) string {
	switch month {
	case time.January:
		return "Januari"
	case time.February:
		return "Februari"
	case time.March:
		return "Maret"
	case time.April:
		return "April"
	case time.May:
		return "Mei"
	case time.June:
		return "Juni"
	case time.July:
		return "Juli"
	case time.August:
		return "Agustus"
	case time.September:
		return "September"
	case time.October:
		return "Oktober"
	case time.November:
		return "November"
	case time.December:
		return "Desember"
	default:
		return ""
	}
}

func Month(month string) string {
	switch month {
	case "01":
		return "Januari"
	case "02":
		return "Februari"
	case "03":
		return "Maret"
	case "04":
		return "April"
	case "05":
		return "Mei"
	case "06":
		return "Juni"
	case "07":
		return "Juli"
	case "08":
		return "Agustus"
	case "09":
		return "September"
	case "10":
		return "Oktober"
	case "11":
		return "November"
	case "12":
		return "Desember"
	default:
		return ""
	}
}

func Day() string {
	waktuSekarang := time.Now()
	hariInt := waktuSekarang.Weekday()

	hariString := []string{"Minggu", "Senin", "Selasa", "Rabu", "Kamis", "Jumat", "Sabtu"}

	namaHari := hariString[hariInt]

	return namaHari
}
