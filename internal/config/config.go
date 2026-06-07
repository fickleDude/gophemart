package config

import (
	"flag"
	"fmt"
	"os"
	"sync"

	"github.com/fickleDude/gophemart/internal/logger"
	"github.com/joho/godotenv"
)

type Config struct {
	runAddr              string
	databaseURI          string
	accrualSystenAddress string
	authKey              string
	logLevel             string
}

var (
	config     *Config
	initConfig sync.Once
)

func GetConfig() *Config {
	initConfig.Do(func() {
		config = &Config{
			runAddr:              "localhost:8080",
			databaseURI:          "postgres://postgres:postgres@localhost:5433/gophermart?sslmode=disable",
			accrualSystenAddress: "http://localhost:8090",
			logLevel:             "info",
		}
		config.parseLocalEnv()
		config.parseFlags()
		config.parseEnv()
	})
	return config
}

func (c *Config) RunAddr() string {
	fmt.Println("run addr = ", c.runAddr)
	return c.runAddr
}

func (c *Config) DatabaseURI() string {
	return c.databaseURI
}

func (c *Config) AccrualSystenAddress() string {
	return c.accrualSystenAddress
}

func (c *Config) AuthKey() string {
	return c.authKey
}

func (c *Config) LogLevel() string {
	return c.logLevel
}

func (c *Config) parseLocalEnv() {
	if err := godotenv.Load(); err != nil {
		panic("No .env file found")
	}

	envAuthKey := os.Getenv("SECRET_KEY")
	if envAuthKey != "" {
		logger.Log.Error("secret key is not specified")
		c.authKey = envAuthKey
	}

	envlogLevel := os.Getenv("LOG_LEVEL")
	if envlogLevel != "" {
		c.logLevel = envlogLevel
	}

}

func (c *Config) parseEnv() {

	envAddr := os.Getenv("RUN_ADDRESS")
	if envAddr != "" {
		c.runAddr = envAddr
	}

	envdatabaseURI := os.Getenv("DATABASE_URI")
	if envdatabaseURI != "" {
		c.databaseURI = envdatabaseURI
	}

	envAccrualSystenAddress := os.Getenv("ACCRUAL_SYSTEM_ADDRESS")
	if envAccrualSystenAddress != "" {
		c.accrualSystenAddress = envAccrualSystenAddress
	}
}

func (c *Config) parseFlags() {
	flag.StringVar(&c.runAddr, "a", c.runAddr, "адрес и порт запуска сервиса")
	flag.StringVar(&c.databaseURI, "d", c.databaseURI, "адрес подключения к базе данных")
	flag.StringVar(&c.accrualSystenAddress, "r", c.accrualSystenAddress, "адрес системы расчёта начислений")
	flag.Parse()
}
