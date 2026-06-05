package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
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
		config = &Config{
			runAddr:              "localhost:8080",
			accrualSystenAddress: "localhost:8081",
		}
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

func checkRunAddr(addr string) error {
	params := strings.Split(addr, ":")
	if len(params) < 2 {
		return fmt.Errorf("формат флага адрес:порт")
	}
	_, err := strconv.Atoi(params[1])
	if err != nil {
		return fmt.Errorf("порт указан некорректно")
	}
	return nil
}

func (c *Config) parseEnv() {
	envAddr := os.Getenv("ADDRESS")
	err := checkRunAddr(envAddr)
	if err == nil {
		c.runAddr = envAddr
	}

	envdatabaseURI := os.Getenv("DATABASE_URI")
	if envdatabaseURI != "" {
		c.databaseURI = envdatabaseURI
	}

	envAccrualSystenAddress := os.Getenv("ACCRUAL_SYSTEM_ADDRESS")
	err = checkRunAddr(envAccrualSystenAddress)
	if err == nil {
		c.accrualSystenAddress = envAccrualSystenAddress
	}

}

func (c *Config) parseFlags() {
	flag.Func("a", "адрес и порт запуска сервиса", func(flagAddr string) error {
		err := checkRunAddr(flagAddr)
		if err == nil {
			c.runAddr = flagAddr

		}
		return nil
	})
	flag.StringVar(&c.databaseURI, "d", "host=localhost port=5433 user=postgres password=postgres dbname=gophermart sslmode=disable", "адрес подключения к базе данных")
	flag.Func("r", "адрес системы расчёта начислений", func(flagAddr string) error {
		err := checkRunAddr(flagAddr)
		if err == nil {
			c.accrualSystenAddress = flagAddr

		}
		return nil
	})
	flag.Parse()
}
