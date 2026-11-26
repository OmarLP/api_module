package verticalreport

import (
	"log"

	"github.com/OmarLP/api_module/internal/domain"
)

type (
	Service interface {
		BuscarIdCita(req domain.BuscarIdCita, offset, limit int64) ([]domain.VerticalReport, error)
		Count(req domain.BuscarIdCita) (int64, error)

		BuscaEspecifica(req domain.BuscaEspecifica, offset, limit int64) ([]domain.VerticalReport, error)
		CountBE(req domain.BuscaEspecifica) (int64, error)
		ExportarBuscaEspecifica(req domain.BuscaEspecifica) ([]domain.VerticalReport, error)

		// filtros en cascada
		GetRedes() ([]string, error)
		GetMicroredes(red string) ([]string, error)
		GetPdDigitaciones(red, microred string) ([]string, error)
		GetEstablecimientos(red, microred, pdDigitacion string) ([]string, error)

		GetCursoVida() ([]string, error)
		ExportarBuscarIdCita(req domain.BuscarIdCita) ([]domain.VerticalReport, error)
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

func (s service) GetCursoVida() ([]string, error) {
	return s.repo.GetCursoVida()
}

func (s service) GetRedes() ([]string, error) {
	return s.repo.GetRedes()
}

func (s service) GetMicroredes(red string) ([]string, error) {
	return s.repo.GetMicroredes(red)
}

func (s service) GetPdDigitaciones(red, microred string) ([]string, error) {
	return s.repo.GetPdDigitaciones(red, microred)
}

func (s service) GetEstablecimientos(red, microred, pdDigitacion string) ([]string, error) {
	return s.repo.GetEstablecimientos(red, microred, pdDigitacion)
}

func (s service) ExportarBuscarIdCita(req domain.BuscarIdCita) ([]domain.VerticalReport, error) {
	resultados, err := s.repo.ExportarBuscarIdCita(req)
	if err != nil {
		s.log.Printf("Error al consultar el repositorio: %v", err)
		return nil, err
	}

	return resultados, nil
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

func (s service) BuscaEspecifica(req domain.BuscaEspecifica, offset, limit int64) ([]domain.VerticalReport, error) {
	// llamar al repositorio para buscar id cita
	resultados, err := s.repo.BuscaEspecifica(req, offset, limit)
	if err != nil {
		s.log.Printf("Error al consultar el repositorio: %v", err)
		return nil, err
	}

	return resultados, nil
}

func (s service) CountBE(req domain.BuscaEspecifica) (int64, error) {
	count, err := s.repo.CountBE(req)
	if err != nil {
		s.log.Printf("Error al calcular los registros en el repositorio: %v", err)
		return 0, err
	}
	return count, nil
}

func (s service) ExportarBuscaEspecifica(req domain.BuscaEspecifica) ([]domain.VerticalReport, error) {
	resultados, err := s.repo.ExportarBuscaEspecifica(req)
	if err != nil {
		s.log.Printf("Error al consultar el repositorio: %w", err)
		return nil, err
	}

	return resultados, nil
}
