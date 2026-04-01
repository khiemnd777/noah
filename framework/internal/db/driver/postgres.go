package driver

import (
	"database/sql"
	"fmt"
	"strings"

	frameworkdb "github.com/khiemnd777/noah_framework/pkg/db"
	_ "github.com/lib/pq"
)

type PostgresClient struct {
	db     *sql.DB
	config frameworkdb.PostgresConfig
}

func NewPostgresClient(cfg frameworkdb.PostgresConfig) *PostgresClient {
	return &PostgresClient{config: cfg}
}

func (p *PostgresClient) Provider() string {
	return "postgres"
}

func (p *PostgresClient) Connect() error {
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		p.config.Host, p.config.Port, p.config.User, p.config.Password, p.config.Name, p.config.SSLMode,
	)

	db, err := sql.Open("postgres", connStr)
	if err == nil && db.Ping() == nil {
		p.db = db
		return nil
	}

	defaultConnStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=postgres sslmode=%s",
		p.config.Host, p.config.Port, p.config.User, p.config.Password, p.config.SSLMode,
	)
	defaultDB, err := sql.Open("postgres", defaultConnStr)
	if err != nil {
		return fmt.Errorf("failed to connect to default postgres DB: %w", err)
	}
	defer defaultDB.Close()

	createSQL := fmt.Sprintf(`CREATE DATABASE "%s"`, p.config.Name)
	if _, err := defaultDB.Exec(createSQL); err != nil && !strings.Contains(err.Error(), "already exists") {
		return fmt.Errorf("failed to create database: %w", err)
	}

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	p.db = db
	return db.Ping()
}

func (p *PostgresClient) Close() error {
	if p.db == nil {
		return nil
	}
	return p.db.Close()
}

func (p *PostgresClient) SQLDB() *sql.DB {
	return p.db
}
