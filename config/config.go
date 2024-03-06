package config

type Config struct {
	MongoCfg MongoCfg
}

type MongoCfg struct {
	URI string
	DB  string
}

func InitConfig() Config {
	return Config{
		MongoCfg: MongoCfg{
			URI: "mongodb://localhost:27017",
			DB:  "dcard",
		},
	}
}
