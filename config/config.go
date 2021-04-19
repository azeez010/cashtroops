package config

import "os"

type Config struct {
	Addr             string
	UsdtWalletId     string
	WalletVendorAddr string
	DatabaseUrl      string
	SessionCacheDir  string
	BlockCypherToken string
}

func New() Config {
	return Config{
		Addr:             os.Getenv("ADDR"),
		DatabaseUrl:      os.Getenv("DATABASE_URL"),
		SessionCacheDir:  os.Getenv("SESSION_CACHE"),
		BlockCypherToken: os.Getenv("BC_TOKEN"),
	}
}
