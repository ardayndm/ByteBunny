package utils

import (
	"sync"

	"github.com/bwmarrin/discordgo"
)

var (
	botAvatarUrl string
	botAvatarMu  sync.RWMutex
)

// İlk verinin boş olma durumunu kontrol eder , eğer boş ise fallback döner.
func GetOrDefault[T comparable](value, fallback T) T {

	var zero T

	if value == zero {
		return fallback
	}
	return value
}

// Doğrudan verilen metinsel sunucu ID'sini kullanarak
// o sunucuya ait özel logonun/ikonun 128x128 boyutundaki URL adresini getirir.
func GetGuildIconFromID(s *discordgo.Session, guildId string) string {
	guild, err := GetGuild(s, guildId)
	if err != nil {
		return ""
	}
	return guild.IconURL("128")
}

// Botun kendi profil resminin (avatar) internet adresini çeker.
// Discord Embed footer alanlarında pikselleşmeyi önlemek adına 32x32 boyut şablonunu kullanır.
func GetBotAvatarURL(s *discordgo.Session) string {
	botAvatarMu.RLock()
	if botAvatarUrl != "" {
		defer botAvatarMu.RUnlock()
		return botAvatarUrl
	}
	botAvatarMu.RUnlock()

	botAvatarMu.Lock()
	defer botAvatarMu.Unlock()

	// Kilit alınana kadar başkası yazmış olabilir, tekrar kontrol et
	if botAvatarUrl != "" {
		return botAvatarUrl
	}

	user, err := s.User("@me")
	if err != nil {
		return ""
	}
	botAvatarUrl = user.AvatarURL("32")
	return botAvatarUrl
}

// Verilen sunucu (Guild) ID'sini kullanarak Discord API'sinden
// sunucuya ait tüm ham veri nesnesini güvenli bir şekilde çeker.
func GetGuild(s *discordgo.Session, guildID string) (*discordgo.Guild, error) {
	// Önce cache'e bak
	guild, err := s.State.Guild(guildID)
	if err == nil {
		return guild, nil
	}
	// Cache'de yoksa API'ye git
	return s.Guild(guildID)
}

// Sunucunun adını döndürür. Eğer işlem bir özel mesaj (DM)
// kanalında gerçekleşiyorsa sunucu ID'si boş olacağından doğrudan boş metin döner.
func GetGuildName(s *discordgo.Session, guildID string) (string, error) {
	if guildID == "" {
		return "", nil
	}

	guild, err := GetGuild(s, guildID)
	if err != nil {
		return "", err
	}
	return guild.Name, nil
}
