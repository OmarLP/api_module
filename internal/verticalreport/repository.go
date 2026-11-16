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

		// metodos para los catalogos en cascada
		GetRedes() ([]string, error)
		GetMicroredes(red string) ([]string, error)
		GetPdDigitaciones(red, microred string) ([]string, error)
		GetEstablecimientos(red, microred, pdDigitacion string) ([]string, error)

		GetCursoVida() ([]string, error)
		ExportarBuscarIdCita(req domain.BuscarIdCita) ([]domain.VerticalReport, error)
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

func (repo *repository) GetCursoVida() ([]string, error) {
	var cursoVida []string
	err := repo.db.Table("reporte_vertical").Distinct().Order("curso_vida").Pluck("curso_vida", &cursoVida).Error

	if err != nil {
		return nil, fmt.Errorf("error al obtener curso de vida: %w", err)
	}
	return cursoVida, nil
}

func (repo *repository) GetRedes() ([]string, error) {
	var redes []string
	err := repo.db.Table("reporte_vertical").Where("red is not null and red != ''").Distinct().Order("red").Pluck("red", &redes).Error

	if err != nil {
		return nil, fmt.Errorf("error al obtener redes: %w", err)
	}

	return redes, nil
}

func (repo *repository) GetMicroredes(red string) ([]string, error) {
	var microredes []string
	query := repo.db.Table("reporte_vertical").Where("microred is not null and microred != ''")

	if red != "" {
		query = query.Where("red = ?", red)
	}

	err := query.Distinct().Order("microred").Pluck("microred", &microredes).Error
	if err != nil {
		return nil, fmt.Errorf("error al obtener microredes: %w", err)
	}
	return microredes, nil
}

func (repo *repository) GetPdDigitaciones(red, microred string) ([]string, error) {
	var pdDigitacion []string
	query := repo.db.Table("reporte_vertical").Where("pd_digitacion is not null and pd_digitacion != ''")

	if red != "" {
		query = query.Where("red = ?", red)
	}

	if microred != "" {
		query = query.Where("microred = ?", microred)
	}

	err := query.Distinct().Order("pd_digitacion").Pluck("pd_digitacion", &pdDigitacion).Error
	if err != nil {
		return nil, fmt.Errorf("error al obtener pdDigitacion: %w", err)
	}
	return pdDigitacion, nil
}

func (repo *repository) GetEstablecimientos(red, microred, pdDigitacion string) ([]string, error) {
	var establecimientos []string
	query := repo.db.Table("reporte_vertical").Where("nombre_establecimiento is not null and nombre_establecimiento != ''")

	if red != "" {
		query = query.Where("red = ?", red)
	}

	if microred != "" {
		query = query.Where("microred = ?", microred)
	}

	if pdDigitacion != "" {
		query = query.Where("pd_digitacion = ?", pdDigitacion)
	}

	err := query.Distinct().Order("nombre_establecimiento").Pluck("nombre_establecimiento", &establecimientos).Error
	if err != nil {
		return nil, fmt.Errorf("error al obtener establecimientos: %w", err)
	}
	return establecimientos, nil
}

func (repo *repository) ExportarBuscarIdCita(req domain.BuscarIdCita) ([]domain.VerticalReport, error) {
	var resultados []domain.VerticalReport
	query := `SELECT * FROM fc_buscar_idcita(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 999999, 0)`

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
		return nil, fmt.Errorf("error al exportar fc_buscar_idcita: %w", err)
	}

	return resultados, nil
}

func (repo *repository) BuscarIdCita(req domain.BuscarIdCita, offset, limit int64) ([]domain.VerticalReport, error) {
	// acá voy a colocar mis resultados
	var resultados []domain.VerticalReport

	query := `SELECT * FROM fc_buscar_idcita(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

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
		limit,
		offset,
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
