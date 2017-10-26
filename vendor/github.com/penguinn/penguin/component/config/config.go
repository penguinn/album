package config

import (
	"flag"
	"time"

	"github.com/spf13/viper"
)

/*
	通过-f传参的优先级最高，然后是Init函数传参次之，最后是默认的文件defaultConfigFile
*/

const (
	defaultConfigFile = "penguin.toml"
)

var (
	cf  = flag.String("f", defaultConfigFile, "Config File Path, Required")
	vip *viper.Viper
)

type ConfigComponent struct{}

func (ConfigComponent) Init(options ...interface{}) error {

	if *cf == defaultConfigFile && len(options) > 0 {
		c, ok := options[0].(string)
		if ok {
			*cf = c
		}
	}

	vip = viper.New()
	vip.SetConfigFile(*cf)
	vip.AutomaticEnv()

	return vip.ReadInConfig()
}

func Get(key string) interface{} {
	return vip.Get(key)
}

func GetBool(key string) bool {
	return vip.GetBool(key)
}

func GetDuration(key string) time.Duration {
	return vip.GetDuration(key)
}

func GetFloat64(key string) float64 {
	return vip.GetFloat64(key)
}

func GetInt(key string) int {
	return vip.GetInt(key)
}

func GetInt64(key string) int64 {
	return vip.GetInt64(key)
}

func GetSizeInBytes(key string) uint {
	return vip.GetSizeInBytes(key)
}

func GetString(key string) string {
	return vip.GetString(key)
}

func GetStringMap(key string) map[string]interface{} {
	return vip.GetStringMap(key)
}

func GetStringMapString(key string) map[string]string {
	return vip.GetStringMapString(key)
}

func GetStringMapStringSlice(key string) map[string][]string {
	return vip.GetStringMapStringSlice(key)
}

func GetStringSlice(key string) []string {
	return vip.GetStringSlice(key)
}

func GetTime(key string) time.Time {
	return vip.GetTime(key)
}

func IsSet(key string) bool {
	return vip.IsSet(key)
}
