package user

import (
	"errors"
	"log"

	"github.com/OmarLP/api_module/internal/domain"
	"gorm.io/gorm"
)

type (
	Repository interface {
		FindRecorderByDocumentNumber(documentNumber string) (*domain.Recorder, error)
		ExistsEmail(email string) (bool, error)
		CreateUser(user *domain.User) error

		FindByEmail(email string) (*domain.User, error)
		UpdatePassword(userID int, hashedPassword string) error
	}

	repository struct {
		log *log.Logger
		db  *gorm.DB
	}
)

func NewRepository(log *log.Logger, db *gorm.DB) Repository {
	return &repository{
		log: log,
		db:  db,
	}
}

// validar si el dni es usuario de hisminsa
func (repo *repository) FindRecorderByDocumentNumber(documentNumber string) (*domain.Recorder, error) {
	var recorder domain.Recorder

	err := repo.db.Where("numero_documento = ?", documentNumber).First(&recorder).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("you're not a HISMINSA user")
		}
		return nil, err
	}

	return &recorder, nil
}

// verificar si el coreo ya existe
func (repo *repository) ExistsEmail(email string) (bool, error) {
	var count int64

	err := repo.db.Model(&domain.User{}).Where("correo = ?", email).Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (repo *repository) CreateUser(user *domain.User) error {

	if err := repo.db.Create(user).Error; err != nil {
		repo.log.Println(err)
		return err
	}

	repo.log.Println("user created with email:", user.Email)
	return nil
}

// solo me valida si existe el email
func (repo *repository) FindByEmail(email string) (*domain.User, error) {
	var recoder domain.User
	err := repo.db.Where("correo = ?", email).First(&recoder).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("your email is not registered")
		}
		return nil, err
	}

	return &recoder, nil
}

// actualizar clave
func (repo *repository) UpdatePassword(userID int, hashedPassword string) error {
	return repo.db.Model(&domain.User{}).
		Where("id_usuario = ?", userID).
		Update("clave", hashedPassword).Error
}
