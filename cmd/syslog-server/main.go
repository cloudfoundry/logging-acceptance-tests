package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	envstruct "code.cloudfoundry.org/go-envstruct"
)

var (
	logger *log.Logger = log.New(os.Stderr, "", log.LstdFlags)
)

func main() {
	cfg := loadConfig()

	http.ListenAndServeTLS(
		fmt.Sprint(":", cfg.HTTPPort),
		cfg.CertFile,
		cfg.KeyFile,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
	)
}

type Config struct {
	HTTPPort int    `env:"HTTP_PORT, report"`
	CertFile string `env:"CERT_FILE, required"`
	KeyFile  string `env:"KEY_FILE, required"`
}

func loadConfig() Config {
	cfg := Config{}
	if err := envstruct.Load(&cfg); err != nil {
		log.Fatal(err)
	}

	envstruct.WriteReport(&cfg)
	return cfg
}
