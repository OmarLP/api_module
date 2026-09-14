package domain

import "time"

type Recorder struct {
	ID             int64     `gorm:"column:id_registrador;primaryKey"`
	DocumentTypeID int       `gorm:"column:id_tipo_documento"`
	DocumentNumber string    `gorm:"column:numero_documento"`
	LastName       string    `gorm:"column:apellido_paterno_registrador"`
	MotherLastName string    `gorm:"column:apellido_materno_registrador"`
	FirstName      string    `gorm:"column:nombres_registrador"`
	BirthDate      time.Time `gorm:"column:fecha_nacimiento"`
}

func (Recorder) TableName() string {
	return "mstr_registrador"
}
