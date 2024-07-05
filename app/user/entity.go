package user

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint64         `json:"id" gorm:"primary_key;autoIncrement"`
	UID       string         `json:"uid" gorm:"type:varchar(40);not null;uniqueIndex:idx_user_uid" validate:"uuid"`
	Name      string         `json:"name" gorm:"type:varchar(200);not null" validate:"required,min=3,max=100"`
	Email     string         `json:"email" gorm:"type:varchar(255);not null;uniqueIndex:idx_user_email" validate:"required,email"`
	Password  string         `json:"-" gorm:"not null" validate:"required,min=6"`
	CreatedAt time.Time      `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updatedAt" gorm:"autoCreateTime;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index:idx_user_deleted_at"`
}
