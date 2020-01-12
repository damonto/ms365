package config

import (
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/damonto/msonline/internal/pkg/logger"
)

// App is the RESTful API server conf
type App struct {
	ListenAddr   string `toml:"listen_addr"`
	AccessKey    string `toml:"access_key"`
	AccessSecret string `toml:"access_secret"`
}

// Microsoft it the microsoft graph api conf
type Microsoft struct {
	Domain       string
	Endpoint     string
	AuthorizeURL string `toml:"authorize_url"`
	TokenURL     string `toml:"token_url"`
	ClientID     string `toml:"client_id"`
	ClientSecret string `toml:"client_secret"`
	TenantID     string `toml:"tenant_id"`
	Scope        string `toml:"scope"`
}

// Config is the App and Microsoft conf
type Config struct {
	App       App
	Microsoft Microsoft
}

// Cfg is the unmarshal conf
var Cfg = &Config{}

// ReadConfig loads config file from .toml file.
func ReadConfig(cfgPath string) {
	buf, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		logger.Sugar.Fatalf("Unable to read config %v", err)
		os.Exit(1)
	}

	err = toml.Unmarshal(buf, Cfg)
	if err != nil {
		logger.Sugar.Fatalf("Unmarshal failed %v", err)
		os.Exit(1)
	}
}
