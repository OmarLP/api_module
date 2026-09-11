package user

import "log"

type (
	Service interface {
		CreateUser(documentNumber, email string) (string, error)
	}

	service struct {
		log  *log.Logger
		repo Repository
	}
)
