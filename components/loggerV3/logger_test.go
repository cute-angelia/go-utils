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
	New(
		WithIsOnline(false),
		WithProject("test.bud"),
		WithLogPath("/tmp"),
		WithLevel(1),
		WithFileJson(false),
		WithHookError(true),
	)

	loggerTest := GetLogger()

	log.Println("heloe ttttttttttttttttttt")

	loggerTest.Warn().Str("name", "kk").Msg("hello, success")
	loggerTest.Info().Msg("hello")
	loggerTest.Info().Msg("cat")
	loggerTest.Info().Msg("dog")
	loggerTest.Error().Msg("dog")
	loggerTest.Error().Str("dog", "xiaokeai").Msg("")
}
