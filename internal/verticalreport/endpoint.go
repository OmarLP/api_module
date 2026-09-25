package verticalreport

import (
	"net/http"

	"github.com/OmarLP/api_module/internal/domain"
)

type (
	Controller func(w http.ResponseWriter, r *http.Request)

	Endpoints struct {
		BuscarIdCita Controller
	}

	Response struct {
		Status int         `json:"status"`
		Data   interface{} `json:"data"`
		Err    string      `json:"error"`
	}
)

func MakeEndpoints(s Service) Endpoints {
	return Endpoints{
		BuscarIdCita: makeBuscarIdCitaEndpoint(s),
	}
}

func makeBuscarIdCitaEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var req domain.BuscarIdCita
	}
}
