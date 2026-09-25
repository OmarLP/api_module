package user

import (
	"errors"
	"log"
	"time"

	"github.com/OmarLP/api_module/internal/domain"
	"github.com/OmarLP/api_module/pkg/authorization"
	"golang.org/x/crypto/bcrypt"
)

type (
	Service interface {
		CreateUser(documentNumber, email string) (*domain.User, error)
		SetFirstPassword(recoder domain.FirstLogin) error
		ResetPassword(req domain.ForgetPassword) error
		Login(email, password string) (*domain.LoginResponse, bool, error)
		RefreshToken(refreshTokenString string) (*domain.LoginResponse, error)
		Logout(refreshTokenStr string) error
		GetProfile(userID int) (*domain.UserProfileResponse, error)
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

func (s service) Login(email, password string) (*domain.LoginResponse, bool, error) {
	// validar que no haya error
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return nil, false, errors.New("invalid credentials")
	}

	// validar que su cuenta este activa
	if user.Status != 1 {
		return nil, false, errors.New("user account is inactive")
	}

	// validar que el usuario tenga contraseña
	if user.Password == nil || *user.Password == "" {
		return nil, true, nil
	}

	// comparar los passwords
	if err := bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(password)); err != nil {
		s.log.Println("invalid credentials:")
		return nil, false, errors.New("invalid credentials")
	}

	// generar Acces Token
	accessToken, err := authorization.GenerateToken(int64(user.IDUser), user.Email)
	if err != nil {
		s.log.Println("error generating token:", err)
		return nil, false, err
	}

	// generar refresh token
	refreshTokenString, err := authorization.GenerateRefreshToken()
	if err != nil {
		return nil, false, err
	}

	// persistir el refresh token en la bd
	refreshToken := &domain.RefreshToken{
		UserID:    user.IDUser,
		Token:     refreshTokenString,
		ExpiresAt: time.Now().Add(5 * time.Minute),
		Revoked:   false,
	}

	if err := s.repo.SaveRefreshToken(refreshToken); err != nil {
		return nil, false, errors.New("failed to save refresh token")
	}

	return &domain.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenString,
	}, false, nil
}

func (s service) RefreshToken(refreshTokenString string) (*domain.LoginResponse, error) {
	// borrar el token en la base de datos
	token, err := s.repo.FindRefreshToken(refreshTokenString)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// validar que no este revocado
	if token.Revoked {
		return nil, errors.New("refresh tokan has been revoked")
	}

	// validar fecha que expira
	if time.Now().After(token.ExpiresAt) {
		return nil, errors.New("refresh token has expired")
	}

	// buscar al usuario asociado y verificar el estado
	user, err := s.repo.FindByID(token.UserID)
	if err != nil || user.Status != 1 {
		return nil, errors.New("user account is unavailable")
	}

	// rotacion de tokens: revocar el refresh token usado
	if err := s.repo.RevokeRefreshToken(token.Token); err != nil {
		return nil, err
	}

	// generar un nuevo AccessToekn y RefreshToken
	newAccessToken, err := authorization.GenerateToken(int64(user.IDUser), user.Email)
	if err != nil {
		return nil, err
	}

	newRefreshTokenString, err := authorization.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	newRefreshToken := &domain.RefreshToken{
		UserID:    user.IDUser,
		Token:     newRefreshTokenString,
		ExpiresAt: time.Now().Add(5 * time.Minute),
		Revoked:   false,
	}

	if err := s.repo.SaveRefreshToken(newRefreshToken); err != nil {
		return nil, err
	}

	return &domain.LoginResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshTokenString,
	}, nil

}

// para el logout
func (s service) Logout(refreshTokenStr string) error {
	if refreshTokenStr == "" {
		return errors.New("refresh token is required")
	}

	return s.repo.RevokeRefreshToken(refreshTokenStr)
}

// para el profile
func (s service) GetProfile(userID int) (*domain.UserProfileResponse, error) {
	profile, err := s.repo.GetUserProfileByID(userID)
	if err != nil {
		return nil, err
	}

	return profile, nil
}
