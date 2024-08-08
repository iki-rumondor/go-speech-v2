package config

import (
	"github.com/iki-rumondor/go-speech/internal/domain/layers/handlers"
	"github.com/iki-rumondor/go-speech/internal/domain/layers/repositories"
	"github.com/iki-rumondor/go-speech/internal/domain/layers/services"
	"gorm.io/gorm"
)

type Handlers struct {
	UserHandler       *handlers.UserHandler
	MasterHandler     *handlers.MasterHandler
	FileHandler       *handlers.FileHandler
	ImportHandler     *handlers.ImportHandler
	AssignmentHandler *handlers.AssignmentHandler
}

func GetAppHandlers(db *gorm.DB) *Handlers {

	master_repo := repositories.NewMasterInterface(db)
	master_service := services.NewMasterService(master_repo)
	master_handler := handlers.NewMasterHandler(master_service)

	file_service := services.NewFileService(master_repo)
	file_handler := handlers.NewFileHandler(file_service)

	import_service := services.NewImportService(master_repo)
	import_handler := handlers.NewImportHandler(import_service)

	user_repo := repositories.NewUserInterface(db)
	user_service := services.NewUserService(user_repo)
	user_handler := handlers.NewUserHandler(user_service)

	assignment_repo := repositories.NewAssignmentInterface(db)
	assignment_service := services.NewAssignmentService(assignment_repo)
	assignment_handler := handlers.NewAssignmentHandler(assignment_service)

	return &Handlers{
		MasterHandler:     master_handler,
		UserHandler:       user_handler,
		FileHandler:       file_handler,
		ImportHandler:     import_handler,
		AssignmentHandler: assignment_handler,
	}
}
