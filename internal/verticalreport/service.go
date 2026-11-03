package verticalreport

import (
	"log"

	"github.com/OmarLP/api_module/internal/domain"
)

type (
	Service interface {
		BuscarIdCita(req domain.BuscarIdCita) ([]domain.VerticalReport, error)
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

func (s service) BuscarIdCita(req domain.BuscarIdCita) ([]domain.VerticalReport, error) {

	// lamar al repositorio para buscar id cita
	resultados, err := s.repo.BuscarIdCita(req)
	if err != nil {
		s.log.Printf("Error al consultar el repositorio: %v", err)
		return nil, err
	}

	return resultados, nil

}
