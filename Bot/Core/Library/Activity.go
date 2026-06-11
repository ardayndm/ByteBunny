package library

import (
	"github.com/bwmarrin/discordgo"
)

// Aktivite tipleri şablonu.
var ActivityTypeMap = map[string]discordgo.ActivityType{
	"playing":   discordgo.ActivityTypeGame,
	"streaming": discordgo.ActivityTypeStreaming,
	"listening": discordgo.ActivityTypeListening,
	"watching":  discordgo.ActivityTypeWatching,
	"competing": discordgo.ActivityTypeCompeting,
}

// Botun genel aktivite durumu şablonu.
type PresenceConfig struct {
	Activities []struct {
		Name    string `yaml:"name"`
		Type    string `yaml:"type"`
		URL     string `yaml:"url"`
		Details string `yaml:"details"` // YENİ: Profil açıklaması
		State   string `yaml:"state"`   // YENİ: Profil durumu/notu
	} `yaml:"activities"`
	Status   string  `yaml:"status"`
	Interval float32 `yaml:"interval"` // Saniye cinsinden yenilenme sıklığı
}
