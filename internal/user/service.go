package user

import (
	"errors"
	"log"

	"github.com/OmarLP/api_module/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type (
	Service interface {
		CreateUser(documentNumber, email string) (*domain.User, error)
		SetFirstPassword(recoder domain.FirstLogin) error
		ResetPassword(req domain.ForgetPassword) error
	}

	service struct {
		log  *log.Logger
		repo Repository
	}
)

// constructor
func NewService(log *log.Logger, repo Repository) Service {
	return &service{
		log:  log,
		repo: repo,
	}
}

// logica del negocio para registrar usuario
func (s service) CreateUser(documentNumber, email string) (*domain.User, error) {
	// verificar si documento existe en mstr_registrador
	recorder, err := s.repo.FindRecorderByDocumentNumber(documentNumber)
	if err != nil {
		return nil, err
	}

	// verificar y evitar duplicidad de correo
	exists, err := s.repo.ExistsEmail(email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("Email already registered")
	}

	// armar la entidad user
	newUser := &domain.User{
		IDRegister: recorder.ID,
		Email:      email,
	}

	// llmamos al repo para insertar en la bd
	if err := s.repo.CreateUser(newUser); err != nil {
		return nil, err
	}

	return newUser, nil

}

func (s service) SetFirstPassword(recoder domain.FirstLogin) error {
	// buscar usuario por correo
	user, err := s.repo.FindByEmail(recoder.Email)
	if err != nil {
		return err
	}

	//validar que el usuario esté activo
	if user.Status != 1 {
		return errors.New("user account is inactive")
	}

	// validar que sea su primer inicio (clave debe ser null)
	if user.Password != nil {
		return errors.New("password has already been set, please use regular login")
	}

	// hashear la contraseña
	//hashedPassword, err := hashPassword(recoder.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(recoder.Password), 10)
	if err != nil {
		return errors.New("failed to process password")
	}

	return s.repo.UpdatePassword(user.IDUser, string(hashedPassword))

}

func (s service) ResetPassword(req domain.ForgetPassword) error {
	// buscar el documento
	recoder, err := s.repo.FindRecorderByDocumentNumber(req.DocumentNumber)
	if err != nil {
		return err
	}

	// buscar correo
	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		return err
	}

	// el id_registrador tiene que estar en user
	if recoder.ID != user.IDRegister {
		return errors.New("document number does not match registered user email")
	}

	// validar que el estado esté activo
	if user.Status != 1 {
		return errors.New("user account is inactive")
	}

	// validar que no sea el primer inicio
	if user.Password == nil {
		return errors.New("you must register your initial password first")
	}

	// hashear la contraseña
	//hashedPassword, err := hashPassword(recoder.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), 10)
	if err != nil {
		return errors.New("failed to process password")
	}

	return s.repo.UpdatePassword(user.IDUser, string(hashedPassword))
}
