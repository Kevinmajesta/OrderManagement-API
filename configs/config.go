package configs

import (
	"errors"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

type Config struct {
	Env      string         `env:"ENV"`  // Dihapus envDefault
	Port     string         `env:"PORT"` // Dihapus envDefault
	Postgres PostgresConfig `envPrefix:"POSTGRES_"`
	Redis    RedisConfig    `envPrefix:"REDIS_"`
	JWT      JwtConfig      `envPrefix:"JWT_"`
	Encrypt  EncryptConfig  `envPrefix:"ENCRYPT_"`
	SMTP     SMTPConfig     `envPrefix:"SMTP_"`
	Midtrans MidtransConfig `envPrefix:"MIDTRANS_"`
}

type SMTPConfig struct {
	Host     string `env:"HOST"` // Dihapus envDefault
	Port     string `env:"PORT"` // Dihapus envDefault
	User     string `env:"USER"`
	Password string `env:"PASSWORD"`
}

type PostgresConfig struct {
	Host     string `env:"HOST"`     // Dihapus envDefault
	Port     string `env:"PORT"`     // Dihapus envDefault
	User     string `env:"USER"`     // Dihapus envDefault
	Password string `env:"PASSWORD"` // Dihapus envDefault
	Database string `env:"DATABASE"` // Dihapus envDefault
}

type JwtConfig struct {
	SecretKey string `env:"SECRET_KEY"`
}

type RedisConfig struct {
	Host     string `env:"HOST"`     // Dihapus envDefault
	Port     string `env:"PORT"`     // Dihapus envDefault
	Password string `env:"PASSWORD"` // Dihapus envDefault
}

type EncryptConfig struct {
	SecretKey string `env:"SECRET_KEY"`
	IV        string `env:"IV"`
}

type MidtransConfig struct {
	ServerKey    string `env:"SERVER_KEY"`
	ClientKey    string `env:"CLIENT_KEY"`
	IsProduction string `env:"IS_PRODUCTION"`
}

func NewConfig(envPath string) (*Config, error) {
	// Memuat .env file. Penting: jika file tidak ada atau ada masalah,
	// godotenv.Load() akan mengembalikan error.
	err := godotenv.Load(envPath)
	if err != nil {
		// Mengembalikan error yang jelas jika gagal memuat file .env
		return nil, errors.New("failed to load .env file: " + err.Error())
	}

	cfg := new(Config)

	// Memparsing environment variables ke dalam struct Config.
	// Jika ada field tanpa envDefault yang tidak ditemukan di env vars,
	// env.Parse akan mengembalikan error.
	err = env.Parse(cfg)
	if err != nil {
		// Mengembalikan error yang jelas jika parsing gagal
		return nil, errors.New("failed to parse config from environment variables: " + err.Error())
	}

	return cfg, nil
}
