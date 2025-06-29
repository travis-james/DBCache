package config

import "os"

type Config struct {
	DatastoreDBHost string
	DatastoreDBPort string
	DatastoreDBUser string
	DatastoreDBPw   string
	DatastoreDBName string
	CacheAddr       string
	CacheUser       string // Can't imagine passing creds like this in a struct is the way to go, will revisit this.
	CachePw         string
	GRPCPort        string
}

func Load() (*Config, error) {
	return &Config{
		DatastoreDBHost: os.Getenv("POSTGRES_HOST"),
		DatastoreDBPort: os.Getenv("POSTGRES_HOST_PORT"),
		DatastoreDBUser: os.Getenv("POSTGRES_USER"),
		DatastoreDBPw:   os.Getenv("POSTGRES_PASSWORD"),
		DatastoreDBName: os.Getenv("POSTGRES_DB"),
		CacheAddr:       os.Getenv("REDIS_ADDR"),
		CacheUser:       os.Getenv("REDIS_USER"),
		CachePw:         os.Getenv("REDIS_PW"),
		GRPCPort:        os.Getenv("GRPC_SERVER_PORT"),
	}, nil
}
