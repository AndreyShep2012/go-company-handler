package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	ListenAddr               string `yaml:"listen_addr" env:"LISTEN_ADDR" env-default:"127.0.0.1:8080" env-description:"Address (IP:port pair) where server listens for the connections"`
	ApiRoot                  string `yaml:"api_root" env:"API_ROOT" env-default:"/api/v1" env-description:"Root path for the API"`
	LogLevel                 string `yaml:"log_level" env:"LOG_LEVEL" env-default:"info" env-description:"Logging level. One of following: debug, info, warn, error"`
	MongoUri                 string `yaml:"mongo_uri" env:"MONGO_URI" env-default:"mongodb://localhost:27017" env-description:"MongoDB connection URI"`
	MongoDatabaseName        string `yaml:"mongo_database_name" env:"MONGO_DATABASE_NAME" env-default:"company-handler" env-description:"MongoDB database name"`
	MongoCompaniesCollection string `yaml:"mongo_companies_collection" env:"MONGO_COMPANIES_COLLECTION" env-default:"companies" env-description:"MongoDB collection name for companies"`
	ConnectTimeoutSec        int    `yaml:"connect_timeout_sec" env:"MONGO_CONNECT_TIMEOUT_SEC" env-default:"5" env-description:"MongoDB connection timeout in seconds"`
	JWTSecretKey             string `yaml:"jwt_secret_key" env:"JWT_SECRET_KEY" env-default:"jwt_secret_key" env-description:"JWT key"`
}

func Load(path string) (cfg Config, err error) {
	if err = cleanenv.ReadConfig(path, &cfg); err != nil {
		fmt.Printf("file '%v' can not be found, using ENV/ENV-Default variables\n", path)
		if err = cleanenv.ReadEnv(&cfg); err != nil {
			return
		}
	}
	return
}
