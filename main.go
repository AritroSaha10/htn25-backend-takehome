package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/AritroSaha10/htn25-backend-takehome/lib"
	"github.com/AritroSaha10/htn25-backend-takehome/model"
	"github.com/AritroSaha10/htn25-backend-takehome/util"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"gorm.io/gorm"
)

// @title           HTN25 Backend API
// @version         0.1
// @description     Backend API for Hack the North 2025 Backend Challenge

// @host      localhost:8080
// @BasePath  /

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
	if err != nil && os.Getenv("ENV") == "development" {
		log.Warn().Err(err).Msg("failed to load environment variables")
	}

	var db *gorm.DB
	if db, err = lib.CreateNewDB(); err != nil {
		log.Fatal().Err(err).Msg("failed to create database")
	}

	// Import any new data from the initial data set
	err = importInitialData(db)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to import initial data")
	}

	// Print out all users
	users := []model.User{}
	db.Find(&users)
	for _, user := range users {
		log.Info().Any("user", user).Msg("user in db")
	}

	// Initialize and start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}
	serv = CreateNewServer(db, port, env)
	serv.MountHandlers()

	log.Info().Str("port", port).Str("env", env).Msg("starting server")
	err = http.ListenAndServe(":"+port, serv.Router)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start server")
	}
}

// importInitialData imports any new data from the initial data set.
func importInitialData(db *gorm.DB) error {
	// Import any new data from the initial data set
	initDataURL := os.Getenv("INITIAL_DATABASE_URL")
	if initDataURL == "" {
		return fmt.Errorf("INITIAL_DATABASE_URL is not set")
	}
	initData, err := util.GetJSONFromURL(initDataURL)
	if err != nil {
		return fmt.Errorf("failed to get initial data: %w", err)
	}
	err = batchAddUsersFromRaw(db, initData)
	if err != nil {
		return fmt.Errorf("failed to add users from initial data: %w", err)
	}
	return nil
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
			user.BadgeCode = userData["badge_code"].(string)

			// Parse scans from the user data
			user.Scans = []model.Scan{}
			for _, rawScan := range userData["scans"].([]interface{}) {
				scanData := rawScan.(map[string]interface{})
				scannedAt, err := time.Parse(util.ISO8601, scanData["scanned_at"].(string))
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
