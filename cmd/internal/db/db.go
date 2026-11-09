package db

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type dbConfig struct {
	connection dbConfigConnection
	options    dBConfigOptions
}

type dbConfigConnection struct {
	host    string
	port    string
	base    string
	user    string
	pass    string
	sslMode string
}

type dBConfigOptions struct {
	maxOpenConnects int
	maxIdleConnects int
}

func Instance() (*sql.DB, error) {
	config, err := getConfig(".env")
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("postgres", fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=%s", config.connection.host, config.connection.user, config.connection.pass, config.connection.base, config.connection.sslMode))
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(config.options.maxOpenConnects)
	db.SetMaxIdleConns(config.options.maxIdleConnects)

	return db, err
}

func Close(db *sql.DB) {
	_ = db.Close()
}

func getConfig(file string) (*dbConfig, error) {
	config := &dbConfig{}
	err := config.load(file)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func (c *dbConfig) load(file string) error {
	err := godotenv.Load(file)
	if err != nil {
		return errors.New("error loading DB config")
	}

	c.connection.host = os.Getenv("POSTGRES_HOST")
	c.connection.port = os.Getenv("POSTGRES_PORT")
	c.connection.base = os.Getenv("POSTGRES_BASE")
	c.connection.user = os.Getenv("POSTGRES_USER")
	c.connection.pass = os.Getenv("POSTGRES_PASS")
	c.connection.sslMode = os.Getenv("POSTGRES_SSL_MODE")

	c.options.maxOpenConnects, _ = strconv.Atoi(os.Getenv("POSTGRES_MAX_OPEN_CONNECTIONS"))
	c.options.maxOpenConnects, _ = strconv.Atoi(os.Getenv("POSTGRES_MAX_IDLE_CONNECTIONS"))

	return nil
}
