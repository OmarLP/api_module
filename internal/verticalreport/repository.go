package verticalreport

import (
	"fmt"
	"log"

	"github.com/OmarLP/api_module/internal/domain"
	"gorm.io/gorm"
)

type (
	Repository interface {
		BuscarIdCita(req domain.BuscarIdCita) ([]domain.VerticalReport, error)
	}

	repository struct {
		log *log.Logger
		db  *gorm.DB
	}
)

func NewRepository(log *log.Logger, db *gorm.DB) Repository {
	return &repository{
		log: log,
		db:  db,
	}
}

func (repo *repository) BuscarIdCita(req domain.BuscarIdCita) ([]domain.VerticalReport, error) {
	// acá voy a colocar mis resultados
	var resultados []domain.VerticalReport

	query := `SELECT * FROM fc_buscar_idcita(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	err := repo.db.Raw(query,
		req.Anio,
		req.MesInicio,
		req.MesFinal,
		req.Red,
		req.Microred,
		req.PdDigitacion,
		req.Establecimiento,
		req.CursoVida,
		req.CieCpm,
		req.Documento,
	).Scan(&resultados).Error

	if err != nil {
		return nil, fmt.Errorf("Error al ejecutar: fc_buscar_idcita: %w", err)
	}

	return resultados, nil

}
