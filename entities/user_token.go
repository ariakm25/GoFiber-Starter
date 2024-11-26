package entities

import (
	internal_log "GoFiber-API/internal/log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserTokenType string

const (
	UserTokenTypeResetPassword UserTokenType = "reset_password"
	UserTokenTypeRefreshToken  UserTokenType = "refresh_token"
)

type UserToken struct {
	UID       string        `json:"uid" gorm:"primary_key;type:varchar(50);not null;uniqueIndex:idx_user_token_uid;<-:create"`
	UserID    string        `json:"userId" gorm:"type:varchar(50);not null;index:idx_user_token_user_id"`
	Token     string        `json:"token" gorm:"type:varchar(255);not null;uniqueIndex:idx_user_token_token"`
	Type      UserTokenType `json:"type" gorm:"type:varchar(50);not null"`
	ExpiredAt time.Time     `json:"expiredAt" gorm:"not null"`
	CreatedAt time.Time     `json:"createdAt" gorm:"autoCreateTime;<-:create"`
	UpdatedAt time.Time     `json:"updatedAt" gorm:"autoCreateTime;autoUpdateTime"`

	User User `json:"user" gorm:"foreignKey:UserID;references:UID"`
}

func (ut *UserToken) BeforeCreate(tx *gorm.DB) (err error) {
	id, err := uuid.NewV7()
	if err != nil {
		internal_log.Logger.Error(err.Error())
		id = uuid.New()
	}
	ut.UID = id.String()

	return
}
