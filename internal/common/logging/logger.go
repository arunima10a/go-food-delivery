package logging

import (
	"os"
	"time"
	"github.com/rs/zerolog"
)
func NewLogger(serviceName string ) zerolog.Logger {
	// In development, we use ConsoleWriter to make it pretty.
	// In production, we would use raw JSON for tools like ELK or Datadog.
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}

	return zerolog.New(output).
	With().
	Timestamp().
	Str("service", serviceName).
	Logger()

}
