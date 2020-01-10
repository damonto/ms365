package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
)

func TestReadConfig(t *testing.T) {
	dir, _ := os.Getwd()
	conf := strings.Replace(dir, "/internal/pkg/config", "/configs/config.toml", -1)
	ReadConfig(conf)

	assert.Equal(t, Cfg.App.ListenAddr, ":8088")
}
