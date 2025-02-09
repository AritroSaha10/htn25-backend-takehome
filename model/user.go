package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID        uint    `json:"id"`
	Name      string  `json:"name"`
	Email     string  `json:"email"`
	Phone     string  `json:"phone"`
	BadgeCode *string `json:"badge_code"` // pointer to allow no badge code
	Scans     []Scan  `json:"scans"`
}
