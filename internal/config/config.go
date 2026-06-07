package config

import (
	"flag"
	"os"
	"sync"
)

type Config struct {
	runAddr              string
	databaseURI          string
	accrualSystenAddress string
}

var (
	config     *Config
	initConfig sync.Once
)

func GetConfig() *Config {
	initConfig.Do(func() {
		config = &Config{}
		config.parseFlags()
		config.parseEnv()
	})
	return config
}

func (c *Config) RunAddr() string {
	return c.runAddr
}

func (c *Config) DatabaseURI() string {
	return c.databaseURI
}

func (c *Config) AccrualSystenAddress() string {
	return c.accrualSystenAddress
}

func (c *Config) parseEnv() {
	envAddr := os.Getenv("ADDRESS")
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
	flag.StringVar(&c.runAddr, "a", "", "адрес и порт запуска сервиса")
	flag.StringVar(&c.databaseURI, "d", "", "адрес подключения к базе данных")
	flag.StringVar(&c.accrualSystenAddress, "r", "", "адрес системы расчёта начислений")
	flag.Parse()
}
