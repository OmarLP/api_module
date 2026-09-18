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
		FindByID(userID int) (*domain.User, error)

		SaveRefreshToken(token *domain.RefreshToken) error
		FindRefreshToken(tokenString string) (*domain.RefreshToken, error)
		RevokeRefreshToken(tokenString string) error

		GetUserProfileByID(userID int) (*domain.UserProfileResponse, error)
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

// verificar si el coreo ya existe y evitar DUPLICIDAD
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

func (repo *repository) FindByID(userID int) (*domain.User, error) {
	var recoder domain.User
	err := repo.db.Where("id_usuario = ?", userID).First(&recoder).Error
	if err != nil {
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

// refresh token
// guardar token
func (repo *repository) SaveRefreshToken(token *domain.RefreshToken) error {
	return repo.db.Create(token).Error
}

// buscar token y verificar que exista
func (repo *repository) FindRefreshToken(tokenString string) (*domain.RefreshToken, error) {
	var token domain.RefreshToken
	err := repo.db.Where("token = ?", tokenString).First(&token).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("refresh token not found")
		}
		return nil, err
	}
	return &token, nil
}

// revocar un token (revoked tiene que estar en true)
func (repo *repository) RevokeRefreshToken(tokenString string) error {
	return repo.db.Model(&domain.RefreshToken{}).
		Where("token = ?", tokenString).
		Update("revoked", true).Error
}

// obtener datos de usuario para perfil
func (repo *repository) GetUserProfileByID(userID int) (*domain.UserProfileResponse, error) {
	var profile domain.UserProfileResponse

	err := repo.db.Table("usuarios u").
		Select("r.apellido_paterno_registrador, split_part(r.nombres_registrador, ' ', 1) as nombres_registrador, u.correo").
		Joins("join mstr_registrador r on(u.id_registrador = r.id_registrador)").
		Where("u.id_usuario = ?", userID).
		Scan(&profile).Error

	if err != nil {
		return nil, err
	}

	return &profile, nil
}
