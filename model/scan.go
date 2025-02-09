package model

import (
	"fmt"
	"net/http"
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

// Bind implements the render.Binder interface for Scan.
func (s *Scan) Bind(r *http.Request) error {
	if s.ActivityName == "" || s.ActivityCategory == "" {
		return fmt.Errorf("activity_name and activity_category are required")
	}

	return nil
}

// Render implements the render.Renderer interface for Scan.
func (s *Scan) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func CreateScan(db *gorm.DB, scan *Scan) error {
	scan.ScannedAt = time.Now()
	if err := db.Create(&scan).Error; err != nil {
		return err
	}
	return nil
}
