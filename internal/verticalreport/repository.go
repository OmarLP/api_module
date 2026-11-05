package verticalreport

import (
	"fmt"
	"log"

	"github.com/OmarLP/api_module/internal/domain"
	"gorm.io/gorm"
)

type (
	Repository interface {
		BuscarIdCita(req domain.BuscarIdCita, offset, limit int64) ([]domain.VerticalReport, error)
		Count(req domain.BuscarIdCita) (int64, error)
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

func (repo *repository) BuscarIdCita(req domain.BuscarIdCita, offset, limit int64) ([]domain.VerticalReport, error) {
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

func (repo *repository) Count(req domain.BuscarIdCita) (int64, error) {
	var total int64

	// Envolvemos tu función de postgres en un SELECT COUNT(*)
	query := `
        SELECT COUNT(*) 
        FROM fc_buscar_idcita(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    `

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
	).Scan(&total).Error

	if err != nil {
		return 0, fmt.Errorf("error al contar registros: %w", err)
	}

	return total, nil
}
