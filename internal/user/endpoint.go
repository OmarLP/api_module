package user

import (
	"encoding/json"
	"net/http"
	"net/mail"

	"github.com/OmarLP/api_module/internal/domain"
)

type (
	Controller func(w http.ResponseWriter, r *http.Request)

	Endpoints struct {
		CreateUser     Controller
		UpdatePassword Controller
		ResetPassword  Controller
	}

	Response struct {
		Status int         `json:"status"`
		Data   interface{} `json:"data"`
		Err    string      `json:"error"`
	}
)

func MakeEndpoints(s Service) Endpoints {
	return Endpoints{
		CreateUser:     makeCreateUserEndpoint(s),
		UpdatePassword: makeUpdatePasswordEndpoint(s),
		ResetPassword:  makeResetPasswordEndpoint(s),
	}
}

func makeCreateUserEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var req domain.Register
		// validación
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "invalid request format"})
			return
		}

		//
		if req.DocumentNumber == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "document number is required"})
			return
		}

		if req.Email == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "email is required"})
			return
		}

		_, err := mail.ParseAddress(req.Email)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "invalid mail format"})
			return
		}

		user, err := s.CreateUser(req.DocumentNumber, req.Email)
		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: err.Error()})
			return
		}

		json.NewEncoder(w).Encode(&Response{Status: 200, Data: user})
	}
}

func makeUpdatePasswordEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var req domain.FirstLogin

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "invalid request format"})
			return
		}

		if req.Email == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "email is required"})
			return
		}

		_, err := mail.ParseAddress(req.Email)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "invalid mail format"})
			return
		}

		if req.Password == "" || req.ConfirmPassword == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "password and confirm password is required"})
			return
		}

		// validar coincidencia del password y su confirmación
		if req.Password != req.ConfirmPassword {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "the passwords do not match"})
			return
		}

		// pasa un struct de firstRegister para verificar en el service si existe email, clave es null y estado = 1
		if err := s.SetFirstPassword(req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&Response{Status: 200, Data: "password set successfully"})

	}
}

func makeResetPasswordEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application-json")

		var req domain.ForgetPassword
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "invalid format request"})
			return
		}

		if req.DocumentNumber == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "document number is required"})
			return
		}

		if req.Email == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "email is required"})
			return
		}

		_, err := mail.ParseAddress(req.Email)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "invalid mail format"})
			return
		}

		if req.NewPassword == "" || req.ConfirmPassword == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "new password and confirm password is required"})
			return
		}

		// validar coincidencia del password y su confirmación
		if req.NewPassword != req.ConfirmPassword {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "passwords do not match"})
			return
		}

		if err := s.ResetPassword(req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&Response{Status: 200, Data: "password successfully reset"})

	}
}
