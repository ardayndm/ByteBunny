package utils

import (
	library "ByteBunny/Bot/Core/Library"
	"os"

	"gopkg.in/yaml.v3"
)

// Belirtilen konumdaki YAML dosyasını verilen yapıya okur.
func ReadYaml(path string, to any) (ok bool, e error) {

	data, err := os.ReadFile(path)
	if err != nil {
		return false, err
	}

	if err := yaml.Unmarshal(data, to); err != nil {
		return false, err
	}

	return true, nil
}

// Belirtilen anahtarı ortam tablosunda arar.
func ReadFromEnv(key string) (value string, e error) {

	if item := os.Getenv(key); item != "" {
		return item, nil
	}

	return "", library.Err_NotFound
}
