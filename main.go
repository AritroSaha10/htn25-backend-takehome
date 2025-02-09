package main

import (
	"os"
	"time"

	"github.com/AritroSaha10/htn25-backend-takehome/model"
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

	// Configure database
	db, err := gorm.Open(sqlite.Open("htn.db"), &gorm.Config{})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect database")
	}
	db.AutoMigrate(&model.User{}, &model.Scan{})
}
