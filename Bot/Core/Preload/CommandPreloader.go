package preload

import (
	"ByteBunny/Bot/Config"
	"ByteBunny/Bot/Core/Events"
	"ByteBunny/Bot/Core/Library"
	utils "ByteBunny/Bot/Core/Utils"
	"sync"
)

// Cache'de tutulacak komut bilgisi
type CachedCommand struct {
	Name        string
	Description string
	Category    string
	Enabled     bool
	Hidden      bool
}

var (
	commandCache = make(map[string]*CachedCommand) // cache deposu
	cacheMu      sync.RWMutex                      // cache kilidi
	preloadDone  = false                           // preload bitti mi?
)

// Bot başlarken çağrılacak - TÜM KOMUTLARI BİR KERE OKU
func PreloadAllCommands() {
	// Tüm komutların isimlerini al
	slashCommands := events.GetRegisteredSlashCommands()
	prefixCommands := events.GetRegisteredPrefixCommands()

	// Tekrar edenleri temizle
	allCommands := make(map[string]bool)
	for name := range slashCommands {
		allCommands[name] = true
	}
	for name := range prefixCommands {
		allCommands[name] = true
	}

	utils.LogToConsole(utils.INFO, "Komut Preloaderi başladı...")

	// Her komutu tek tek oku ve cache'e kaydet
	for cmdName := range allCommands {
		loadAndCacheCommand(cmdName)
	}

	preloadDone = true
	utils.LogToConsole(utils.OK, "Komut Preloaderi tamamlandı!")
}

// Tek bir komutu okuyup cache'e kaydet
func loadAndCacheCommand(cmdName string) {
	// Geçersiz komut adlarını kontrol et
	if utils.ValidateStringAlphanumeric(cmdName) != nil {
		return
	}

	var targetLib library.CommandLib

	// Önce core komut dene
	corePath, err := utils.ValidateCommandYamlFilepathStrict(cmdName, config.AppConfig.Bot.Lang, true)
	if err == nil {
		if ok, _ := utils.ReadYaml(corePath, &targetLib); ok {
			saveToCache(cmdName, &targetLib)
			return
		}
	}

	// Normal komut dene
	normalPath, err := utils.ValidateCommandYamlFilepathStrict(cmdName, config.AppConfig.Bot.Lang, false)
	if err == nil {
		if ok, _ := utils.ReadYaml(normalPath, &targetLib); ok {
			saveToCache(cmdName, &targetLib)
			return
		}
	}
}

// Cache'e kaydet
func saveToCache(name string, lib *library.CommandLib) {
	enabled, _ := lib.GetOptionBool("enabled")
	hidden, _ := lib.GetOptionBool("hidden")
	category, _ := lib.GetOptionString("category")

	if category == "" {
		category = "other"
	}

	cacheMu.Lock()
	defer cacheMu.Unlock()

	commandCache[name] = &CachedCommand{
		Name:        lib.Name,
		Description: lib.Description,
		Category:    category,
		Enabled:     enabled,
		Hidden:      hidden,
	}
}

// Dışarıdan cache'e erişmek için - PRELOAD BİTTİYSE CACHE'DEN DÖN
func GetCachedCommand(name string) (*CachedCommand, bool) {
	cacheMu.RLock()
	defer cacheMu.RUnlock()

	cmd, ok := commandCache[name]
	return cmd, ok
}

// Tüm cache'lenmiş komutları al (listcommands için)
func GetAllCachedCommands() []*CachedCommand {
	cacheMu.RLock()
	defer cacheMu.RUnlock()

	commands := make([]*CachedCommand, 0, len(commandCache))
	for _, cmd := range commandCache {
		commands = append(commands, cmd)
	}
	return commands
}

// Preload tamamlandı mı?
func IsPreloadDone() bool {
	return preloadDone
}

// Komutun Cache içinde olup olmadığını kontrol eder
func IsCommandExists(cmdName string) bool {
	if !IsPreloadDone() {
		// Preload bitmediyse FALSE döndür (güvenli taraf)
		// veya registry'ye direkt erişim için event'e sor
		return false
	}

	cacheMu.RLock()
	defer cacheMu.RUnlock()
	_, exists := commandCache[cmdName]
	return exists
}
