package model

import (
	"time"

	"gorm.io/gorm"
)

// GORM model for a scan. We aren't using gorm.Model so we can
// add json tags to the fields it provides.
type Scan struct {
	ID               uint           `json:"id" gorm:"primaryKey"`
	CreatedAt        time.Time      `json:"-"`
	UpdatedAt        time.Time      `json:"-"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
	ActivityName     string         `json:"activity_name" gorm:"not null"`
	ActivityCategory string         `json:"activity_category" gorm:"not null"`
	ScannedAt        time.Time      `json:"scanned_at" gorm:"not null"`
	UserID           uint           `json:"-" gorm:"not null"`
}
