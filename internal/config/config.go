package config

import "os"

type Config struct {
	PostgresURL string
	RedisAddr   string
	RedisUser   string // Can't imagine passing creds like this in a struct is the way to go, will revisit this.
	RedisPw     string
	GRPCPort    string
}

func Load() (*Config, error) {
	return &Config{
		PostgresURL: os.Getenv("POSTGRES_URL"),
		RedisAddr:   os.Getenv("REDIS_ADDR"),
		RedisUser:   os.Getenv("REDIS_USER"),
		RedisPw:     os.Getenv("REDIS_PW"),
		GRPCPort:    os.Getenv("GRPC_PORT"),
	}, nil
}
