package user

import (
	"log"

	"gorm.io/gorm"
)

type (
	Repository interface {
		CreateUser(documentNumber, email string) (string, error)
	}

	repository struct {
		log log.Logger
		db  *gorm.DB
	}
)
