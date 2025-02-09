package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/AritroSaha10/htn25-backend-takehome/model"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// Configure logging. Current logging settings are best for development, but
	// it would probably be best to log to a file in JSON format in production.
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	})

	err := godotenv.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load environment variables")
	}

	// Configure database
	db, err := gorm.Open(sqlite.Open("htn.db"), &gorm.Config{})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect database")
	}
	db.AutoMigrate(&model.User{}, &model.Scan{})

	// Import any new data from the initial data set
	initDataURL := os.Getenv("INITIAL_DATABASE_URL")
	if initDataURL == "" {
		log.Fatal().Msg("INITIAL_DATABASE_URL is not set")
	}
	initData, err := getJSONFromURL(initDataURL)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get initial data")
	}
	err = batchAddUsersFromRaw(db, initData)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to add users from initial data")
	}

	// Print out all users
	users := []model.User{}
	db.Find(&users)
	for _, user := range users {
		log.Info().Any("user", user).Msg("user in db")
	}
}

// batchAddUsersFromRaw adds all users from a raw interface without creating duplicates. Duplicates are
// defined as users with the same email address. We are assuming that the initial data will follow the
// same schema as specified in the challenge description.
func batchAddUsersFromRaw(db *gorm.DB, data interface{}) error {
	for _, rawUser := range data.([]interface{}) {
		userData := rawUser.(map[string]interface{})
		user := model.User{}
		tx := db.Limit(1).Find(&user, "email = ?", userData["email"].(string))

		// Create a new user if they don't exist
		if tx.RowsAffected == 0 {
			user.Name = userData["name"].(string)
			user.Email = userData["email"].(string)
			user.Phone = userData["phone"].(string)
			user.BadgeCode = sql.NullString{
				String: userData["badge_code"].(string),
				Valid:  userData["badge_code"].(string) != "",
			}

			// Parse scans from the user data
			user.Scans = []model.Scan{}
			for _, rawScan := range userData["scans"].([]interface{}) {
				scanData := rawScan.(map[string]interface{})
				scannedAt, err := time.Parse(ISO8601, scanData["scanned_at"].(string))
				if err != nil {
					return fmt.Errorf("failed to parse scanned_at: %w", err)
				}
				user.Scans = append(user.Scans, model.Scan{
					ActivityName:     scanData["activity_name"].(string),
					ActivityCategory: scanData["activity_category"].(string),
					ScannedAt:        scannedAt,
				})
			}

			res := db.Create(&user)
			if res.Error != nil {
				return fmt.Errorf("failed to create user: %w", res.Error)
			}
		}
	}

	return nil
}
