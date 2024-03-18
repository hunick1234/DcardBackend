package config

import "time"

type Config struct {
	MongoCfg MongoCfg
}

type MongoCfg struct {
	URI            string
	DB             string
	Password       string
	Username       string
	MaxPoolSize    uint64
	MinPoolSize    uint64
	ConnectTimeout time.Duration
}

func InitConfig() Config {
	return Config{
		MongoCfg: MongoCfg{
			URI:            "mongodb://localhost:27017",
			DB:             "dcard",
			MaxPoolSize:    100,
			MinPoolSize:    0,
			ConnectTimeout: 30000 * time.Millisecond,
		},
	}
}
