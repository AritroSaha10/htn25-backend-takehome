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

type ScanAggregate struct {
	ActivityName     string `json:"activity_name"`
	ActivityCategory string `json:"activity_category"`
	Frequency        int    `json:"frequency"`
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

// Render implements the render.Renderer interface for ScanAggregate.
func (s *ScanAggregate) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func CreateScan(db *gorm.DB, scan *Scan) error {
	scan.ScannedAt = time.Now()
	if err := db.Create(&scan).Error; err != nil {
		return err
	}
	return nil
}

func GetScanAggregates(db *gorm.DB, activityCategory *string, minFrequency *int, maxFrequency *int) ([]ScanAggregate, error) {
	var aggregates []ScanAggregate
	tx := db.
		Model(&Scan{}).
		Select("activity_name, activity_category, COUNT(*) as frequency").
		Group("activity_name, activity_category")

	// Add all query options to tx
	if activityCategory != nil {
		tx = tx.Where("activity_category = ?", *activityCategory)
	}
	if minFrequency != nil && maxFrequency != nil {
		tx = tx.Having("frequency BETWEEN ? AND ?", *minFrequency, *maxFrequency)
	} else if minFrequency != nil {
		tx = tx.Having("frequency >= ?", *minFrequency)
	} else if maxFrequency != nil {
		tx = tx.Having("frequency <= ?", *maxFrequency)
	}

	res := tx.Find(&aggregates)
	if res.Error != nil {
		return nil, res.Error
	}
	return aggregates, nil
}
