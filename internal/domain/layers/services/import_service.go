package services

import (
	"log"
	"strings"

	"github.com/iki-rumondor/go-speech/internal/domain/layers/interfaces"
	"github.com/iki-rumondor/go-speech/internal/domain/structs/models"
	"github.com/iki-rumondor/go-speech/internal/domain/structs/response"
	"github.com/xuri/excelize/v2"
)

type ImportService struct {
	Repo interfaces.MasterInterface
}

func NewImportService(repo interfaces.MasterInterface) *ImportService {
	return &ImportService{
		Repo: repo,
	}
}

func (s *ImportService) SaveTeachers(xlsxPath string) ([]map[string]string, error) {
	f, err := excelize.OpenFile(xlsxPath)
	if err != nil {
		log.Println("Gagal Membuka File")
		return nil, response.SERVICE_INTERR
	}
	defer f.Close()

	rows, err := f.GetRows("Daftar Dosen")
	if err != nil {
		log.Println("Failed to get rows Mahasiswa")
		return nil, response.SERVICE_INTERR
	}
	var failedImport []map[string]string

	for i := 1; i < 10; i++ {
		cols := rows[i]
		var department models.Department
		var depCond = models.Department{
			Name: cols[9],
		}

		if err := s.Repo.FirstOrCreate(&department, depCond); err != nil {
			log.Println(err)
			failedImport = append(failedImport, map[string]string{
				"name":    cols[2],
				"message": err.Error(),
			})
			continue
		}

		nidn := strings.ReplaceAll(cols[1], "'", "")

		model := models.Teacher{
			DepartmentID: department.ID,
			Nidn:         nidn,
			User: &models.User{
				Name:     cols[2],
				Username: nidn,
				Password: nidn,
				RoleID:   2,
				Active:   true,
			},
		}

		if err := s.Repo.Create(&model); err != nil {
			log.Println(err)
			failedImport = append(failedImport, map[string]string{
				"name":    cols[2],
				"message": err.Error(),
			})
			continue
		}

	}

	return failedImport, nil
}
