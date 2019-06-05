package model

import (
	"time"
)

// Account 用户表
type Account struct {
	ID           uint      `gorm:"primary_key"`
	UserID       string    `gorm:"type:varchar(255);not null;column:user_id"`
	Email        string    `gorm:"type:varchar(255);not null;column:email"`
	AccessToken  string    `gorm:"type:text;not null;column:access_token"`
	RefreshToken string    `gorm:"type:text;not null;column:refresh_token"`
	ExpiresIn    time.Time `gorm:"column:expires_in"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time `sql:"index"`
}

// TableName 表名
func (Account) TableName() string {
	return "accounts"
}
