package events

import (
	config "ByteBunny/Bot/Config"
	utils "ByteBunny/Bot/Core/Utils"
	"fmt"
	"sync"

	"github.com/bwmarrin/discordgo"
)

// Tüm komutlarının uyması gereken ortak arayüz.
type CommandShared interface {
	GetGeneralFormatKeys(extraFormat ...map[string]string) map[string]string // Komutun genel mesaj format karşılıkları (ve ek formatları)
	LoadCommandLib() error                                                   // Komutun YAML dosyasını yükler
}

// Tüm Slash komutlarının uyması gereken arayüz.
type SlashCommand interface {
	Name() string                                                 // Komut adı
	Description() string                                          // Komut açıklaması
	Options() []*discordgo.ApplicationCommandOption               // Komut ayarları
	Execute(s *discordgo.Session, i *discordgo.InteractionCreate) // Komut çalıştırma
	CommandShared
}

// Tüm Prefix komutlarının uyması gereken arayüz.
type PrefixCommand interface {
	Name() string                                                            // Komut adı
	Execute(s *discordgo.Session, m *discordgo.MessageCreate, args []string) // Komut çalıştırma
	CommandShared
}

var (
	slashRegistry           = make(map[string]SlashCommand)
	prefixRegistry          = make(map[string]PrefixCommand)
	loadCacheSlashRegistry  []SlashCommand  // Init işlemlerinden sonra registrye almak için bir cache
	loadCachePrefixRegistry []PrefixCommand // Init işlemlerinden sonra registrye almak için bir cache
	slashMutex              sync.RWMutex
	prefixMutex             sync.RWMutex
)

// Sisteme kayıtlı olan Slash komutlarını döndürür.
func GetRegisteredSlashCommands() map[string]SlashCommand {
	slashMutex.RLock()
	defer slashMutex.RUnlock()
	result := make(map[string]SlashCommand, len(slashRegistry))
	for k, v := range slashRegistry {
		result[k] = v
	}
	return result
}

// Sisteme kayıtlı olan Prefix komutlarını döndürür.
func GetRegisteredPrefixCommands() map[string]PrefixCommand {
	prefixMutex.RLock()
	defer prefixMutex.RUnlock()
	result := make(map[string]PrefixCommand, len(prefixRegistry))
	for k, v := range prefixRegistry {
		result[k] = v
	}
	return result
}

// Sisteme Slash komutu kaydeder.
func registerSlashCommand(cmd SlashCommand) {
	slashMutex.Lock()
	defer slashMutex.Unlock()
	slashRegistry[cmd.Name()] = cmd
	utils.LogToConsole(utils.DEBUG, "Slash komutu Cache'den alınıp Sisteme kaydedildi: "+"/"+cmd.Name())
}

// Sisteme Prefix komutu kaydeder.
func registerPrefixCommand(cmd PrefixCommand) {
	prefixMutex.Lock()
	defer prefixMutex.Unlock()
	prefixRegistry[cmd.Name()] = cmd

	utils.LogToConsole(utils.DEBUG, "Prefix komutu Cache'den alınıp Sisteme kaydedildi: "+config.AppConfig.Bot.Prefix+cmd.Name())
}

// Cache içindeki tüm komutları Kayıt listesine yükler.
func LoadAllRegistryCommands() {

	for _, fn := range loadCacheSlashRegistry {
		registerSlashCommand(fn)
	}

	for _, fn := range loadCachePrefixRegistry {
		registerPrefixCommand(fn)
	}
}

// Sisteme Slash komutu kaydeder.
func RegisterSlashCacheCommand(cmd SlashCommand) {
	slashMutex.Lock()
	defer slashMutex.Unlock()
	loadCacheSlashRegistry = append(loadCacheSlashRegistry, cmd)
	utils.LogToConsole(utils.DEBUG, "Slash komutu cache'e alındı")
}

// Sisteme Prefix komutu kaydeder.
func RegisterPrefixCacheCommand(cmd PrefixCommand) {
	prefixMutex.Lock()
	defer prefixMutex.Unlock()
	loadCachePrefixRegistry = append(loadCachePrefixRegistry, cmd)

	utils.LogToConsole(utils.DEBUG, "Prefix komutu cache'e alındı")
}

// Sisteme kayıt edilen Slash komutlarını discorda gönderir ve kayıt ettirir.
func SyncSlashCommands(s *discordgo.Session) error {

	slashMutex.RLock()
	var bulkCommands []*discordgo.ApplicationCommand
	for _, cmd := range slashRegistry {
		bulkCommands = append(bulkCommands, &discordgo.ApplicationCommand{
			Name:        cmd.Name(),
			Description: cmd.Description(),
			Options:     cmd.Options(),
		})
	}
	slashMutex.RUnlock() // defer yok, manuel bırak

	if len(bulkCommands) == 0 {
		utils.LogToConsole(utils.DEBUG, "Sync edilecek komut yok, eskiler siliniyor...")
		_, err := s.ApplicationCommandBulkOverwrite(s.State.User.ID, "", []*discordgo.ApplicationCommand{})
		if err != nil {
			return fmt.Errorf("mevcut komutlar temizlenemedi: %w", err)
		}
		return nil
	}

	_, err := s.ApplicationCommandBulkOverwrite(s.State.User.ID, "", bulkCommands)
	if err != nil {
		return fmt.Errorf("slash komutları sync edilirken hata: %w", err)
	}

	utils.LogToConsole(utils.OK, fmt.Sprintf("%d adet Slash komutu Discord ile senkronize edildi.", len(bulkCommands)))
	return nil
}
