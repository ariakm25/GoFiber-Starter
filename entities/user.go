package entities

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
