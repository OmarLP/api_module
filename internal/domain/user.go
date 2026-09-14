package domain

import "time"

type User struct {
	IDUser        int        `gorm:"primarykey;column:id_usuario;autoIncrement" json:"id_usuario"`
	IDRegister    int64      `gorm:"column:id_registrador;unique;not null" json:"id_registrador"`
	Email         string     `gorm:"column:correo;type:varchar(100);unique;not null" json:"correo"`
	Password      *string    `gorm:"column:clave;type:varchar(255)" json:"-"` // *string permite NULL en DB
	Status        int        `gorm:"column:estado;default:1" json:"estado"`
	CreatedAt     time.Time  `gorm:"column:fecha_creacion;type:timestamptz:default:now()" json:"fecha_creación"`
	InactivatedAt *time.Time `gorm:"column:fecha_baja;type:timestamptz" json:"fecha_baja"`
}

// especificamos el nombre exacto de la tabla en la base de datos
func (User) TableName() string {
	return "usuarios"
}

// DTO's Data Transfer Objects para entradas
// Registro de usuario
type Register struct {
	DocumentNumber string `json:"numero_documento"`
	Email          string `json:"correo"`
}

// Primer inicio se sesión del usuario, para cambiar la contraseña
type FirstLogin struct {
	Email           string `json:"correo"`
	Password        string `json:"clave"`
	ConfirmPassword string `json:"confirmar_clave" validate:"eqfield=Password"`
}

// Solicita el reseteo de contraseña
type ForgetPassword struct {
	DocumentNumber  string `json:"numero_documento" validate:"required"`
	Email           string `json:"correo" validate:"required,email"`
	NewPassword     string `json:"new_password"`
	ConfirmPassword string `json:"confirm_password"`
}

// login
type Login struct {
	Email    string `json:"correo" validate:"required,email"`
	Password string `json:"clave" validate:"required"`
}
