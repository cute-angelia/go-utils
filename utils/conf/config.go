package conf

import (
	"bytes"
	"fmt"
	"github.com/spf13/viper"
	"io"
	"log"
	"path"
)

func LoadConfigFile(cfgFile string) error {
	viper.SetConfigFile(cfgFile)
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println(fmt.Errorf("Fatal error config file: Config file not found %w \n", err))
		} else {
			log.Println(fmt.Errorf("Fatal error config file: %w \n", err))
		}
	}
	return err
}

func MustLoadConfigFile(cfgFile string) {
	if err := LoadConfigFile(cfgFile); err != nil {
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}
}

func LoadConfigByte(data []byte, filetype string) error {
	var err error

	switch filetype {
	case "yaml":
	case "yml":
	case "toml":
	case "json":
	default:
		err = fmt.Errorf("file ext not support")
	}

	if err != nil {
		return err
	}

	viper.SetConfigType(filetype)
	err = viper.ReadConfig(bytes.NewBuffer(data))

	return err
}

func MustLoadConfigByte(data []byte, filetype string) {
	if err := LoadConfigByte(data, filetype); err != nil {
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}
}

// MergeConfig 合并配置 byte
func MergeConfig(byteCfg io.Reader) error {
	return viper.MergeConfig(byteCfg)
}

// MergeConfigWithPath 合并配置 文件路径 cfgPath 文件夹路径
func MergeConfigWithPath(cfgPath string) error {
	cfgPath = path.Base(cfgPath)
	// 追加一份配置
	viper.AddConfigPath(cfgPath)
	err := viper.MergeInConfig() // Find and read the config file
	if err != nil {              // Handle errors reading the config file
		return fmt.Errorf("Fatal error config file: %w \n", err)
	} else {
		return nil
	}
}

func MustMergeConfigWithPath(cfgPath string) {
	cfgPath = path.Base(cfgPath)
	// 追加一份配置
	viper.AddConfigPath(cfgPath)
	err := viper.MergeInConfig() // Find and read the config file
	if err != nil {              // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}
}

// MergeConfigWithMap 合并配置 map
func MergeConfigWithMap(cfg map[string]interface{}) error {
	return viper.MergeConfigMap(cfg)
}

// GetEnv 获取 系统环境变量
func GetEnv(key string) interface{} {
	viper.AutomaticEnv()
	return viper.Get(key)
}
