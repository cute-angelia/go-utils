package loggerV3

import (
	"github.com/cute-angelia/go-utils/utils/conf"
	"log"
	"testing"
)

func TestLoggerLoad(t *testing.T) {
	conf.LoadConfigFile("config_test.toml")
	logger := Load("logger.test").NewLogger()
	log.Println("heloe")
	logger.Warn().Str("name", "kk").Msg("hello, success")
	logger.Info().Msg("hello")
	logger.Info().Msg("cat")
	logger.Info().Msg("dog")
	logger.Error().Msg("dog")
}

func TestLogger(t *testing.T) {
	logger := New(
		WithIsOnline(false),
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
