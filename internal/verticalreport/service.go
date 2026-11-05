package verticalreport

import (
	"log"

	"github.com/OmarLP/api_module/internal/domain"
)

type (
	Service interface {
		BuscarIdCita(req domain.BuscarIdCita, offset, limit int64) ([]domain.VerticalReport, error)
		Count(req domain.BuscarIdCita) (int64, error)
	}

	service struct {
		log  *log.Logger
		repo Repository
	}
)

// constructor
func NewService(log *log.Logger, repo Repository) Service {
	return &service{
		log:  log,
		repo: repo,
	}
}

func (s service) BuscarIdCita(req domain.BuscarIdCita, offset, limit int64) ([]domain.VerticalReport, error) {

	// lamar al repositorio para buscar id cita
	resultados, err := s.repo.BuscarIdCita(req, offset, limit)
	if err != nil {
		s.log.Printf("Error al consultar el repositorio: %v", err)
		return nil, err
	}

	return resultados, nil

}

func (s service) Count(req domain.BuscarIdCita) (int64, error) {
	count, err := s.repo.Count(req)
	if err != nil {
		s.log.Printf("Error al calcular los registros en el repositorio: %v", err)
		return 0, err
	}

	return count, nil

}
