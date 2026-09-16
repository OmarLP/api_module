package user

import (
	"encoding/json"
	"net/http"
	"net/mail"

	"github.com/OmarLP/api_module/internal/domain"
	"github.com/OmarLP/api_module/internal/middleware"
	"github.com/OmarLP/api_module/pkg/authorization"
)

type (
	Controller func(w http.ResponseWriter, r *http.Request)

	Endpoints struct {
		CreateUser     Controller
		UpdatePassword Controller
		ResetPassword  Controller
		Login          Controller
		Profile        Controller
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
		Login:          makeLoginEndpoint(s),
		Profile:        makeProfileEndpoint(s),
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
		w.Header().Set("Content-Type", "application/json")

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

func makeLoginEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var req domain.Login
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
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "invalid email format"})
			return
		}

		if req.Password == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "password is required"})
			return
		}

		token, err := s.Login(req.Email, req.Password)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(&Response{Status: 401, Err: "invalid credentials"})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&Response{Status: 200, Data: token})
	}
}

func makeProfileEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		// obtener los claims guardados en authmiddleware
		claims, ok := r.Context().Value(middleware.UserKey).(*authorization.Claim)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(&Response{Status: 500, Err: "failed to parse user context"})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&Response{Status: 200, Err: "access granted to protected route", Data: claims})
	}
}
