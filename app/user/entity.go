package user

import (
	internal_log "GoFiber-API/internal/log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User

type User struct {
	UID       string         `json:"uid" gorm:"primary_key;type:varchar(50);not null;uniqueIndex:idx_user_uid;<-:create" validate:"uuid"`
	Name      string         `json:"name" gorm:"type:varchar(200);not null" validate:"required,min=3,max=100"`
	Email     string         `json:"email" gorm:"type:varchar(255);not null;uniqueIndex:idx_user_email" validate:"required,email"`
	Password  string         `json:"-" gorm:"not null" validate:"required,min=6"`
	CreatedAt time.Time      `json:"createdAt" gorm:"autoCreateTime;<-:create"`
	UpdatedAt time.Time      `json:"updatedAt" gorm:"autoCreateTime;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index:idx_user_deleted_at"`

	UserTokens []UserToken `json:"-" gorm:"foreignKey:UserID;references:UID"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	id, err := uuid.NewV7()
	if err != nil {
		internal_log.Logger.Error(err.Error())
		id = uuid.New()
	}
	u.UID = id.String()

	return
}

// User Token

type UserTokenType string

const (
	UserTokenTypeResetPassword UserTokenType = "reset_password"
	UserTokenTypeRefreshToken  UserTokenType = "refresh_token"
)

type UserToken struct {
	UID       string        `json:"uid" gorm:"primary_key;type:varchar(50);not null;uniqueIndex:idx_user_token_uid;<-:create" validate:"uuid"`
	UserID    string        `json:"userId" gorm:"type:varchar(50);not null;index:idx_user_token_user_id" validate:"uuid"`
	Token     string        `json:"token" gorm:"type:varchar(255);not null;uniqueIndex:idx_user_token_token" validate:"required"`
	Type      UserTokenType `json:"type" gorm:"type:varchar(50);not null" validate:"required"`
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
