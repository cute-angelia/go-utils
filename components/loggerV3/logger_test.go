package loggerV3

import (
	"log"
	"testing"
)

func TestLogger(t *testing.T) {
	// develop
	logger := Load().Build(
		WithIsOnline(true),
		WithProject("test.bud"),
		WithLevel(1),
	).NewLogger()

	log.Println("heloe")

	logger.Warn().Str("name", "kk").Msg("hello, success")
	logger.Info().Msg("hello")
	logger.Info().Msg("cat")
	logger.Info().Msg("dog")
	logger.Error().Msg("dog")
}
