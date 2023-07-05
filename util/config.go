package util

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

const (
	DB_DRIVER                     = "DB_DRIVER"
	DB_DRIVER_DEFAULT             = "postgres"
	DB_SOURCE                     = "DB_SOURCE"
	DB_MAX_IDLE_CONNS             = "DB_MAX_IDLE_CONNS"
	DB_MAX_IDLE_CONNS_DEFAULT     = 5
	DB_MAX_OPEN_CONNS             = "DB_MAX_OPEN_CONNS"
	DB_MAX_OPEN_CONNS_DEFAULT     = 10
	DB_CONN_MAX_IDLE_TIME         = "DB_CONN_MAX_IDLE_TIME"
	DB_CONN_MAX_IDLE_TIME_DEFAULT = 1 * time.Second
	DB_CONN_MAX_LIFE_TIME         = "DB_CONN_MAX_LIFE_TIME"
	DB_CONN_MAX_LIFE_TIME_DEFAULT = 30 * time.Second
)

var envMap map[string]string

func LoadConfig() error {
	_, err := os.Stat("env")
	if !errors.Is(err, os.ErrNotExist) {
		log.Println("Loading configuration from file ./env")
		envMap, err = godotenv.Read("env")
		if err != nil {
			log.Println("Error while loading configuration from ./orders_service/env : ", err)
			return err
		}
		for k, v := range envMap {
			log.Printf("\t%v=%v\n", k, v)
		}
	} else {
		log.Println("Loading configuration from environment variables..")
	}

	required := []string{
		DB_SOURCE,
	}
	errMsg := ""
	for _, param := range required {
		config := get(param)
		if config == nil || *config == "" {
			errMsg = fmt.Sprintf("%vMissing configuration parameter '%v' \n", errMsg, param)
			continue
		}
	}
	_, err = ConfigDBMaxIdleConns()
	if err != nil {
		errMsg = fmt.Sprintf("%vInvalid configuration : %v", errMsg, err)
	}
	_, err = ConfigDBMaxOpenConns()
	if err != nil {
		errMsg = fmt.Sprintf("%vInvalid configuration : %v", errMsg, err)
	}
	_, err = ConfigDBConnMaxIdleTime()
	if err != nil {
		errMsg = fmt.Sprintf("%vInvalid configuration : %v", errMsg, err)
	}
	_, err = ConfigDBConnMaxLifeTime()
	if err != nil {
		errMsg = fmt.Sprintf("%vInvalid configuration : %v", errMsg, err)
	}

	if errMsg != "" {
		return fmt.Errorf(errMsg)
	}

	return nil
}

func ConfigDBDriver() string {
	value := get(DB_DRIVER)
	if value == nil {
		return DB_DRIVER_DEFAULT
	}
	return *value
}

func ConfigDBSource() *string {
	return get(DB_SOURCE)
}

func ConfigDBMaxIdleConns() (int, error) {
	value, err := getInt(DB_MAX_IDLE_CONNS)
	if err != nil {
		return 0, err
	}
	if value == nil {
		return DB_MAX_IDLE_CONNS_DEFAULT, nil
	}
	return *value, nil

}
func ConfigDBMaxOpenConns() (int, error) {
	value, err := getInt(DB_MAX_OPEN_CONNS)
	if err != nil {
		return 0, err
	}
	if value == nil {
		return DB_MAX_OPEN_CONNS_DEFAULT, nil
	}
	return *value, nil

}
func ConfigDBConnMaxIdleTime() (time.Duration, error) {
	value, err := getInt(DB_CONN_MAX_IDLE_TIME)
	if err != nil {
		return 0, err
	}
	if value == nil {
		return DB_CONN_MAX_IDLE_TIME_DEFAULT, nil
	}
	return time.Second * time.Duration(*value), nil

}
func ConfigDBConnMaxLifeTime() (time.Duration, error) {
	value, err := getInt(DB_CONN_MAX_LIFE_TIME)
	if err != nil {
		return 0, err
	}
	if value == nil {
		return DB_CONN_MAX_LIFE_TIME_DEFAULT, nil
	}
	return time.Second * time.Duration(*value), nil
}

func get(key string) *string {
	var value string
	if envMap != nil {
		var found bool
		value, found = envMap[key]
		if found {
			return &value
		}
	}
	value = os.Getenv(key)
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	return &value
}

func getInt(key string) (*int, error) {
	value := get(key)
	if value == nil {
		return nil, nil
	}
	i, err := strconv.Atoi(*value)
	if err != nil {
		return nil, fmt.Errorf("(%v) is not a valid value for %v", *value, key)
	}
	return &i, nil
}
