package config

import (
	utils "ByteBunny/Bot/Core/Utils"
)

type BotConfig struct {
	Token    string
	Name     string `yaml:"name"`
	Version  string `yaml:"version"`
	AuthorID string `yaml:"author_id"`
	Prefix   string `yaml:"prefix"`
	Lang     string `yaml:"lang"`
}

// Uygulamanın ayarlar şablonu.
type Config struct {
	Bot      BotConfig
	LogLevel utils.LogLevel
}

// Uygulamanın ayarları.
var AppConfig *Config

// Uygulama ayarlarını okur ve yükler.
func InitConfigApp() error {

	conf := &Config{}

	// Bot bilgilerini yükle.
	if err := loadBot(conf); err != nil {
		return err
	}

	// Log seviye bilgilerini yükle.
	if logStr, err := utils.ReadFromEnv("LOG_LEVEL"); err == nil {
		utils.SetAppLogLevel(utils.GetLogLevelFromString(logStr))
	} else {
		return err
	}

	AppConfig = conf
	return nil
}

// Botu yükle.
func loadBot(cfg *Config) error {
	token, err := utils.ReadFromEnv("BOT_TOKEN")

	if err != nil {
		return err
	}

	if ok, err := utils.ReadYaml("Config/Bot.yaml", &cfg.Bot); err != nil || !ok {
		return err
	}

	cfg.Bot.Token = token

	// Bot ayarları içerisinde dil bulunamaz ise default olarak tr ata
	if cfg.Bot.Lang == "" {
		cfg.Bot.Lang = "tr"
		utils.LogToConsole(utils.WARN, "Bot dili bulunamadı, `tr` ile devam ediliyor...")
	}
	return nil

}
