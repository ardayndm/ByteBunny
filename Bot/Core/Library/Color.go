package library

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Renkler şablonu.
type ColorLib struct {
	MainColor int          `yaml:"main"`
	Status    statusColors `yaml:"status"`
}

// Durum renkleri şablonu.
type statusColors struct {
	Success int `yaml:"success"`
	Error   int `yaml:"error"`
	Warn    int `yaml:"warn"`
	Info    int `yaml:"info"`
}

var Colors *ColorLib // Uygulama renkleri.

// Renkleri YAML dosyasından okur ve yükler.
func InitColor(path string) error {

	data, err := os.ReadFile(path)

	if err != nil {
		return Err_FileReadError
	}

	var lib ColorLib
	if err := yaml.Unmarshal(data, &lib); err != nil {
		return Err_NotFound
	}

	Colors = &lib
	return nil
}
