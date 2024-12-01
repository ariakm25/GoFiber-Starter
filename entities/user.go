package entities

import (
	internal_log "GoFiber-API/internal/log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	UID       string         `json:"uid" gorm:"primary_key;type:varchar(50);not null;uniqueIndex:idx_user_uid;<-:create"`
	Name      string         `json:"name" gorm:"type:varchar(200);not null"`
	Email     string         `json:"email" gorm:"type:varchar(255);not null;uniqueIndex:idx_user_email"`
	Status    string         `json:"status" gorm:"type:varchar(20);not null;default:'active'"`
	Password  string         `json:"-" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime;<-:create"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoCreateTime;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index:idx_user_deleted_at"`

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
