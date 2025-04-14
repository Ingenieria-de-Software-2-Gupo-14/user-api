package config

import (
	"github.com/joho/godotenv"
	"os"
	"strings"
)

type Config struct {
	Host        string
	Port        string
	Environment string
}

// LoadConfig loads environment variables a Config Struct containing relevant variables
func LoadConfig() Config {
	currentWorkDirectory, _ := os.Getwd()
	split := strings.Split(currentWorkDirectory, "/")
	join := strings.Join(split[:(len(split)-2)], "/")
	err := godotenv.Load(join + "/.env")
	if err != nil {
		return Config{}
	}

	return Config{os.Getenv("HOST"), os.Getenv("PORT"), os.Getenv("ENVIROMENT")}
}
