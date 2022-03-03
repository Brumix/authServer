package models

import "gorm.io/gorm"

type Hello struct {
	Msg string `json:"msg"`
}

type User struct {
	gorm.Model
	UserName string `gorm:"unique;not null"`
	Email    string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
}

type ChangePass struct {
	Email       string
	OldPassword string
	NewPassword string
}

type DeleteUser struct {
	Email    string
	Password string
}
