package config

import (
	"errors"
	"os"
	"strconv"
)

type AppConfig struct {
	NodeId      int64
	BaseUrl     string
	ShorterPort int
}

type DbConfig struct {
	Host     string
	User     string
	Password string
	Database string
	Port     int
}

type Config struct {
	App AppConfig
	Db  DbConfig
}

func LoadFromEnv() (Config, error) {
	app, err := loadAppFromEnv()
	if err != nil {
		return Config{}, err
	}

	db, err := loadDbFromEnv()
	if err != nil {
		return Config{}, err
	}

	return Config{App: app, Db: db}, nil
}

func loadDbFromEnv() (DbConfig, error) {
	host, ok := os.LookupEnv("DB_HOST")
	if !ok {
		return DbConfig{}, errors.New("DB_HOST environment variable not set")
	}

	user, ok := os.LookupEnv("DB_USER")
	if !ok {
		return DbConfig{}, errors.New("DB_USER environment variable not set")
	}

	password, ok := os.LookupEnv("DB_PASSWORD")
	if !ok {
		return DbConfig{}, errors.New("DB_PASSWORD environment variable not set")
	}

	database, ok := os.LookupEnv("DB_DATABASE")
	if !ok {
		return DbConfig{}, errors.New("DB_DATABASE environment variable not set")
	}

	portStr, ok := os.LookupEnv("DB_PORT")
	if !ok {
		return DbConfig{}, errors.New("DB_PORT environment variable not set")
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return DbConfig{}, errors.New("DB_PORT environment variable not set")
	}

	return DbConfig{
		Host:     host,
		User:     user,
		Password: password,
		Database: database,
		Port:     port,
	}, nil
}

func loadAppFromEnv() (AppConfig, error) {
	shorterPortStr, ok := os.LookupEnv("SHORTER_LISTEN_PORT")
	if !ok {
		return AppConfig{}, errors.New("SHORTER_LISTEN_PORT must be set")
	}

	shorterPort, err := strconv.Atoi(shorterPortStr)
	if err != nil {
		return AppConfig{}, errors.New("SHORTER_LISTEN_PORT must be an integer")
	}

	nodeIdStr, ok := os.LookupEnv("SHORTER_NODE_ID")
	if !ok {
		return AppConfig{}, errors.New("SHORTER_NODE_ID must be set")
	}

	nodeId, err := strconv.ParseInt(nodeIdStr, 10, 64)
	if err != nil {
		return AppConfig{}, errors.New("SHORTER_NODE_ID must be an integer")
	}

	baseUrl, ok := os.LookupEnv("SHORTER_BASE_URL")
	if !ok {
		return AppConfig{}, errors.New("SHORTER_BASE_URL must be set")
	}

	return AppConfig{
		NodeId:      nodeId,
		BaseUrl:     baseUrl,
		ShorterPort: shorterPort,
	}, nil
}
