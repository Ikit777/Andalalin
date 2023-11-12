package utils

import (
	"fmt"

	"github.com/Ikit777/E-Andalalin/repository"
)

// GetCredentialsByRole func for getting credentials from a role name.
func GetCredentialsByRole(role string) ([]string, error) {
	// Define credentials variable.
	var credentials []string

	// Switch given role.
	switch role {
	case repository.SuperAdminRoleName:
		// Super Admin credentials.
		credentials = []string{
			repository.UserAddCredential,
			repository.UserDeleteCredential,
			repository.UserGetCredential,

			repository.AndalalinGetCredential,
			repository.AndalalinUpdateCredential,
			repository.AndalalinTindakLanjut,
			repository.AndalalinSurveyCredential,
			repository.AndalalinPemasanganCredential,
			repository.AndalalinAddOfficerCredential,
			repository.AndalalinTicket2Credential,
			repository.AndalalinPersetujuanCredential,
			repository.AndalalinKelolaTiket,
			repository.AndalalinKeputusanHasil,
			repository.AndalalinSurveiKepuasan,
			repository.AndalalinDokumenCredential,

			repository.ProductAddCredential,
			repository.ProductDeleteCredential,
			repository.ProductUpdateCredential,
		}
	case repository.DinasPerhubunganRoleName:
		// Admin credentials.
		credentials = []string{
			repository.AndalalinGetCredential,

			repository.AndalalinSurveyCredential,
			repository.AndalalinSurveiKepuasan,
		}
	case repository.AdminRoleName:
		// Admin credentials.
		credentials = []string{
			repository.AndalalinGetCredential,
			repository.AndalalinUpdateCredential,

			repository.AndalalinPersetujuanCredential,
			repository.AndalalinKelolaTiket,
			repository.AndalalinTicket2Credential,
			repository.AndalalinSurveyCredential,
			repository.AndalalinKeputusanHasil,
			repository.AndalalinSurveiKepuasan,
			repository.AndalalinDokumenCredential,

			repository.ProductAddCredential,
			repository.ProductDeleteCredential,
			repository.ProductUpdateCredential,
		}
	case repository.OperatorRoleName:
		// Operator credentials.
		credentials = []string{
			repository.AndalalinGetCredential,
			repository.AndalalinUpdateCredential,
			repository.AndalalinTindakLanjut,
			repository.AndalalinAddOfficerCredential,
			repository.AndalalinSurveyCredential,
			repository.AndalalinKelolaTiket,
			repository.AndalalinSurveiKepuasan,
			repository.AndalalinDokumenCredential,
		}
	case repository.OfficerRoleName:
		// Officer credentials.
		credentials = []string{
			repository.AndalalinGetCredential,
			repository.AndalalinSurveyCredential,
			repository.AndalalinPemasanganCredential,
			repository.AndalalinTicket2Credential,
			repository.AndalalinSurveiKepuasan,
		}
	case repository.UserRoleName:
		// User credentials.
		credentials = []string{
			repository.AndalalinPengajuanCredential,
			repository.AndalalinPersyaratanredential,
			repository.AndalalinDokumenCredential,
		}
	default:
		// Return error message.
		return nil, fmt.Errorf("role '%v' does not exist", role)
	}

	return credentials, nil
}
