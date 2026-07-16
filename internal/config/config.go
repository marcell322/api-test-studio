package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port       string
	DBPath     string
	JWTSecret  string
	JWTExpireH int
}

func Load() *Config {
	port := os.Getenv("PORT")
	if port == "" { port = ":8080" }
	db := os.Getenv("DB_PATH")
	if db == "" { db = "app.db" }
	secret := os.Getenv("JWT_SECRET")
	if secret == "" { secret = "change-me" }
	exp := 24
	if v := os.Getenv("JWT_EXPIRE_H"); v != "" {
		if n, err := strconv.Atoi(v); err == nil { exp = n }
	}
	return &Config{Port: port, DBPath: db, JWTSecret: secret, JWTExpireH: exp}
}
