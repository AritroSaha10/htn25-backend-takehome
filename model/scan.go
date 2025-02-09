package model

import (
	"time"

	"gorm.io/gorm"
)

type Scan struct {
	gorm.Model
	ID               uint      `json:"id"`
	ActivityName     string    `json:"activity_name"`
	ActivityCategory string    `json:"activity_category"`
	ScannedAt        time.Time `json:"scanned_at"`
	UserID           uint      `json:"user_id"` // TODO: May want to adjust JSON output later
}
