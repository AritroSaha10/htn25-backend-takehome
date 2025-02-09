package model

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

// GORM model for a user. We aren't using gorm.Model so we can
// add json tags to the fields it provides.
type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Name      string         `json:"name"`
	Email     string         `json:"email" gorm:"unique;not null"`
	Phone     string         `json:"phone" gorm:"not null"`
	BadgeCode sql.NullString `json:"badge_code"`
	Scans     []Scan         `json:"scans"`
}
