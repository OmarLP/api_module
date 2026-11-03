package verticalreport

import (
	"encoding/json"
	"net/http"
	"time"

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
		// validar request de acuerdo a la estructura solicitada
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()

		if err := decoder.Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "invalid request format"})
			return
		}

		//validar anio
		currentYear := time.Now().Year()

		if req.Anio < 2000 || req.Anio > currentYear {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "year must be between 2000 and 2026"})
			return
		}

		// el mes de inicio no puede ser mayor que el mes final
		if req.MesInicio > req.MesFinal {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "start month cannot be greater than end month"})
			return
		}

		data, err := s.BuscarIdCita(req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(&Response{Status: 500, Err: err.Error()})
			return
		}

		json.NewEncoder(w).Encode(&Response{Status: 200, Data: data})

	}
}
