package config

import (
	"flag"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Config struct {
	RunAddress       string `env:"RUN_ADDRESS"`
	ConnectionString string `env:"DATABASE_URI"` // строка подключения к БД
	Debug            bool   `env:"DEBUG_MODE"`   // уровень логирования
	SecretKey        string `env:"KEY"`          //JWT ключ для формирования подписи
}

func init() {
	viper.SetDefault("address", "localhost:3021")
	viper.SetDefault("conection", "defaultValue")
	viper.SetDefault("debug", true)
	viper.SetDefault("key", "SuperSecretKey")
}

func bindToEnv() {
	viper.SetEnvPrefix("GK")
	_ = viper.BindEnv("RUN_ADDRESS")
	_ = viper.BindEnv("DATABASE_URI")
	_ = viper.BindEnv("KEY")
	_ = viper.BindEnv("DEBUG_MODE")
}

func bindToFlag() {
	pflag.String("a", "localhost:3021", "Server address")
	pflag.Bool("l", true, "logger mode")
	pflag.String("d", "postgresql://postgres:postgres@localhost:5432/gomart?sslmode=disable", "Database connection string(PostgreSql)")
	pflag.String("k", "My_super_secret_KEY", "JWT Key to create signature")
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	_ = viper.BindPFlags(pflag.CommandLine)
}

func ReadConfig() (*Config, error) {
	bindToFlag()
	bindToEnv()

	result := &Config{
		RunAddress:       viper.GetString("address"),
		ConnectionString: viper.GetString("connection"),
		Debug:            viper.GetBool("debug"),
		SecretKey:        viper.GetString("key"),
	}

	return result, nil
}

func (cfs *Config) IsDebug() bool {
	return false
}

func (cfs *Config) GetDBConnection() string {
	return ""
}
