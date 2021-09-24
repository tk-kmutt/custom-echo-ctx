package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID       string `json:"id" gorm:"primaryKey"`
	Name     string `json:"name" gorm:"unique"`
	Password string `json:"password"`
}

func (u *User) BeforeCreate(*gorm.DB) error {
	newUUID := uuid.New()
	u.ID = newUUID.String()
	return nil
}
