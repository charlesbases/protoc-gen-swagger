package conf

import (
	"path/filepath"

	"github.com/BurntSushi/toml"

	"github.com/charlesbases/protoc-gen-swagger/logger"
)

const configfile = "swagger.toml"

var conf = new(config)

// config .
type config struct {
	Host    string
	Service string
	Header  header
}

// header .
type header struct {
	Auth string
}

// Get .
func Get() *config {
	return conf
}

// Parse .
func Parse(f string) {
	if len(f) == 0 {
		f = "."
	}
	abspath, err := filepath.Abs(f)
	if err != nil {
		logger.Fatal(err)
	}
	if _, err := toml.DecodeFile(filepath.Join(abspath, configfile), conf); err != nil {
		logger.Fatal(err)
	}
}
