package entities

import (
	internal_log "GoFiber-API/internal/log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserSession struct {
	UID          string    `json:"uid" gorm:"primary_key;type:varchar(50);not null;uniqueIndex:idx_user_session_uid;<-:create"`
	UserID       string    `json:"userId" gorm:"type:varchar(50);not null;index:idx_user_session_user_id"`
	RefreshToken string    `json:"-" gorm:"type:text;not null;uniqueIndex:idx_user_session_refresh_token"`
	DeviceInfo   string    `json:"device_info" gorm:"type:text;not null"`
	ExpiredAt    time.Time `json:"expired_at" gorm:"not null"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime;<-:create"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoCreateTime;autoUpdateTime"`

	User User `json:"user" gorm:"foreignKey:UserID;references:UID"`
}

func (ut *UserSession) BeforeCreate(tx *gorm.DB) (err error) {
	id, err := uuid.NewV7()
	if err != nil {
		internal_log.Logger.Error(err.Error())
		id = uuid.New()
	}
	ut.UID = id.String()

	return
}
