package logging

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Init ZeroLog as per our standard preferences.
// If `shouldBeInDebugMode` is true, then the log level will be set to
// `zerolog.DebugLevel`, otherwise it will be set to `zerolog.InfoLevel`.
//
// Additionally, we will enable pretty printing to stderr if debug is
// turned on, otherwise we will default to JSON logging.
func Init(shouldBeInDebugMode bool) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	if shouldBeInDebugMode {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}
