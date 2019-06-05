package model

import (
	"time"
)

// Account 用户表
type Account struct {
	ID           uint       `gorm:"primary_key" json:"-"`
	UserID       string     `gorm:"type:varchar(255);not null;column:user_id" json:"user_id"`
	Email        string     `gorm:"type:varchar(255);not null;column:email" json:"email"`
	AccessToken  string     `gorm:"type:text;not null;column:access_token" json:"-"`
	RefreshToken string     `gorm:"type:text;not null;column:refresh_token" json:"-"`
	ExpiresIn    int        `gorm:"type:int:10;not null;column:expires_in" json:"-"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `sql:"index" json:"-"`
}

// TableName 表名
func (Account) TableName() string {
	return "accounts"
}
