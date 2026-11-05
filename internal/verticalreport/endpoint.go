package verticalreport

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/OmarLP/api_module/internal/domain"
	"github.com/OmarLP/api_module/pkg/meta"
)

type (
	Controller func(w http.ResponseWriter, r *http.Request)

	Endpoints struct {
		BuscarIdCita Controller
	}

	Response struct {
		Status int         `json:"status"`
		Data   interface{} `json:"data,omitempty"`
		Err    string      `json:"error,omitempty"`
		Meta   *meta.Meta  `json:"meta,omitempty"`
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

		// leer parametros de la url (?page=2&limit=10)
		v := r.URL.Query()
		limit, _ := strconv.ParseInt(v.Get("limit"), 10, 64)
		page, _ := strconv.ParseInt(v.Get("page"), 10, 64)

		// obtener el total de registros
		count, err := s.Count(req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(&Response{Status: 500, Err: err.Error()})
			return
		}

		// inicio de metadata de paginarcion
		meta, err := meta.New(int64(page), int64(limit), int64(count))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(&Response{Status: 500, Err: err.Error()})
			return
		}

		// respuesta con req, offset y limit
		data, err := s.BuscarIdCita(req, meta.Offset(), meta.Limit())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(&Response{Status: 500, Err: err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&Response{Status: 200, Data: data, Meta: meta})

	}
}
