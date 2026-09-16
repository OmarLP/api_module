package domain

import "time"

type RefreshToken struct {
	ID        int64     `gorm:"column:id;primaryKey;autoIncrement"`
	UserID    int       `gorm:"column:id_usuario;not null"`
	Token     string    `gorm:"column:token;not null;unique"`
	ExpiresAt time.Time `gorm:"column:expires_at;type:timestamptz;not null"`
	Revoked   bool      `gorm:"column:revoked;default:false"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamptz;autoCreateTime"`
}

func (RefreshToken) TableName() string {
	return "refresh_tokens"
}
