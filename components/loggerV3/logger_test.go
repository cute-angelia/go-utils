package loggerV3

import (
	"log"
	"testing"
)

//func TestLoggerLoad(t *testing.T) {
//	conf.LoadConfigFile("config_test.toml")
//	logger := Load("logger.test").NewLogger()
//	logger.Warn().Str("name", "kk").Msg("hello, success")
//	logger.Info().Msg("hello")
//	logger.Info().Msg("cat")
//	logger.Info().Msg("dog")
//	logger.Error().Msg("dog")
//	log.Println("heloe xxxxxxxxxxxxxxxxxxxxxx")
//}

func TestLogger(t *testing.T) {
	logger := New(
		WithIsOnline(true),
		WithProject("test.bud"),
		WithLogPath("/tmp"),
		WithLevel(1),
		WithFileJson(false),
		WithHookError(true),
	)

	log.Println("heloe ttttttttttttttttttt")

	logger.Warn().Str("name", "kk").Msg("hello, success")
	logger.Info().Msg("hello")
	logger.Info().Msg("cat")
	logger.Info().Msg("dog")
	logger.Error().Msg("dog")
}
