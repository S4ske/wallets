package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"net/url"
)

type Config struct {
	DBHost  string
	DBPort  int
	DBUser  string
	DBPass  string
	DBName  string
	Address string
}

func LoadConfig() *Config {
	viper.SetConfigFile("config.env")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("no config.env file found")
	}

	cfg := &Config{
		DBHost:  viper.GetString("DB_HOST"),
		DBPort:  viper.GetInt("DB_PORT"),
		DBUser:  viper.GetString("DB_USER"),
		DBPass:  viper.GetString("DB_PASSWORD"),
		DBName:  viper.GetString("DB_NAME"),
		Address: viper.GetString("ADDRESS"),
	}

	return cfg
}

func (c *Config) PostgresURL() string {
	u := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(c.DBUser, c.DBPass),
		Host:   fmt.Sprintf("%s:%d", c.DBHost, c.DBPort),
		Path:   c.DBName,
	}
	q := u.Query()
	q.Set("sslmode", "disable")
	u.RawQuery = q.Encode()
	return u.String()
}
