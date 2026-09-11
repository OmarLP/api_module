package user

import (
	"encoding/json"
	"net/http"

	"github.com/OmarLP/api_module/internal/domain"
)

type (
	Controller func(w http.ResponseWriter, r *http.Request)

	Endpoints struct {
		CreateUser Controller
	}

	Response struct {
		Status int         `json:"status"`
		Data   interface{} `json:"data"`
		Err    string      `json:"error"`
	}
)

func MakeEndpoints(s Service) Endpoints {
	return Endpoints{
		CreateUser: MakeCreateUserEndpoint(s),
	}
}

func MakeCreateUserEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		var req domain.Register
		// validación
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "invalid request format"})
			return
		}

		//
		if req.DocumentNumber == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "document number is required"})
			return
		}

	}
}
