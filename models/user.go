package models

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Username string  `json:"username"`
	Email    string  `json:"email"`
	Password string  `json:"-"`
	VKID     float64 `json:"vk_id"`
	GoogleID string  `json:"google_id"`
}
