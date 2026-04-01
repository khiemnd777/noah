package db

type Config struct {
	Provider    string
	AutoMigrate bool
	Postgres    PostgresConfig
	MongoDB     MongoConfig
}

type PostgresConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
}

type MongoConfig struct {
	URI      string
	Database string
}

type Client interface {
	Connect() error
	Close() error
	Provider() string
}

type SQLBridge interface {
	SQLDB() any
}
