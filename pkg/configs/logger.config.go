package configs

import (
	"os"

	"github.com/rs/zerolog"
)

func SetupApplicationLogger() *zerolog.Logger {
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	return &logger
}
