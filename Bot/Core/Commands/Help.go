package commands

import (
	config "ByteBunny/Bot/Config"
	events "ByteBunny/Bot/Core/Events"
	library "ByteBunny/Bot/Core/Library"
	utils "ByteBunny/Bot/Core/Utils"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// ── Tipler ───────────────────────────────────────────────────────────────────────
type HelpCommandSlash struct{ cmd *library.CommandLib }
type HelpCommandPrefix struct{ cmd *library.CommandLib }

// ── Oto kayıt ───────────────────────────────────────────────────────────────────────

func init() {
	events.RegisterSlashCacheCommand(&HelpCommandSlash{})
	events.RegisterPrefixCacheCommand(&HelpCommandPrefix{})
}

// ── Yükle ────────────────────────────────────────────────────────────────────

func loadHelpLib() (*library.CommandLib, error) {
	var lib library.CommandLib
	path := filepath.Join("Modules", "Locales", config.AppConfig.Bot.Lang, "Commands", "Core", "help.yaml")
	ok, err := utils.ReadYaml(path, &lib)
	if err != nil || !ok {
		return nil, err
	}
	return &lib, nil
}

// ── Kontrol ────────────────────────────────────────────────────────────────────

func (helpSlash *HelpCommandSlash) checkLibrary() bool {
	if helpSlash.cmd != nil {
		return true
	}
	return helpSlash.LoadCommandLib() == nil
}

func (helpPrefix *HelpCommandPrefix) checkLibrary() bool {
	if helpPrefix.cmd != nil {
		return true
	}
	return helpPrefix.LoadCommandLib() == nil
}

// ── Slash Komutu ────────────────────────────────────────────────────────────────────

func (helpSlash *HelpCommandSlash) LoadCommandLib() error {
	lib, err := loadHelpLib()

	helpSlash.cmd = lib

	return err
}

// Komut adı
func (helpSlash *HelpCommandSlash) Name() string {
	if !helpSlash.checkLibrary() {
		return "help" // fallback
	}

	return helpSlash.cmd.Name
}

// Komut açıklaması
func (helpSlash *HelpCommandSlash) Description() string {
	if !helpSlash.checkLibrary() {
		return "Komut hakkında bilgi verir."
	}
	return helpSlash.cmd.Description
}

// Komut seçenekleri
func (helpSlash *HelpCommandSlash) Options() []*discordgo.ApplicationCommandOption {

	if !helpSlash.checkLibrary() {
		return nil
	}

	var opts []*discordgo.ApplicationCommandOption
	for _, arg := range helpSlash.cmd.GetSlashArgs() {
		opts = append(opts, &discordgo.ApplicationCommandOption{
			Name:        arg.Name,
			Description: arg.Description,
			Type:        discordgo.ApplicationCommandOptionString,
			Required:    arg.Required,
		})
	}
	return opts
}

// Placeholder karşılıkları
func (helpSlash *HelpCommandSlash) GetGeneralFormatKeys(extra ...map[string]string) map[string]string {
	return map[string]string{
		"prefix": "/",
	}
}

// Komut işleyicisi
func (helpSlash *HelpCommandSlash) Execute(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !helpSlash.checkLibrary() || helpSlash.cmd == nil {
		utils.LogToConsole(utils.ERROR, "Help: komut kütüphanesi yüklenemedi")
		return
	}

	// Komut aktif değil ise işlemi iptal et
	if enabled, ok := helpSlash.cmd.GetOptionBool("enabled"); !enabled || !ok {
		return
	}

	// Slash komutu 3 saniye içinde cevap vermek zorunda, uzun işlemler için defer gönder
	utils.SendInteractionDeferRespond(s, i.Interaction)

	target := utils.Target{Interaction: i, IsEphemeral: true, IsFollowup: true}

	// Slash argümanından komut adını al
	cmdName := ""
	if len(i.ApplicationCommandData().Options) > 0 {
		cmdName = strings.TrimSpace(i.ApplicationCommandData().Options[0].StringValue())
	}

	helpExecute(s, target, cmdName, helpSlash.cmd)
}

// ──────────────────────────────────────────────────────────────────────
// ── Prefix Komutu ────────────────────────────────────────────────────────────────────

// Komut adı
func (helpPrefix *HelpCommandPrefix) Name() string {
	if !helpPrefix.checkLibrary() {
		return "help" // fallback
	}

	return helpPrefix.cmd.Name
}

// Placeholder karşılıkları
func (helpPrefix *HelpCommandPrefix) GetGeneralFormatKeys(extra ...map[string]string) map[string]string {
	return map[string]string{
		"prefix": config.AppConfig.Bot.Prefix,
	}
}

func (helpPrefix *HelpCommandPrefix) LoadCommandLib() error {
	lib, err := loadHelpLib()
	helpPrefix.cmd = lib

	return err
}

func (helpPrefix *HelpCommandPrefix) Execute(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if !helpPrefix.checkLibrary() || helpPrefix.cmd == nil {
		utils.LogToConsole(utils.ERROR, "Help: komut kütüphanesi yüklenemedi")
		return
	}

	if enabled, ok := helpPrefix.cmd.GetOptionBool("enabled"); !enabled || !ok {
		return
	}

	target := utils.Target{Message: m}

	// args[0] komut adı, args[1] varsa hedef komut
	cmdName := ""
	if len(args) > 1 {
		cmdName = strings.TrimSpace(args[1])
	}

	helpExecute(s, target, cmdName, helpPrefix.cmd)
}

// ──────────────────────────────────────────────────────────────────────
// ── İç Mantık ────────────────────────────────────────────────────────────────────

// Verilen komut adını registry de kontrol eder , Eğer bir registry'de bile kayıtlı değilse false döner.
func isAllowedCommandName(cmdName string) bool {

	// Registry'de de kontrol et — yaml yoksa veya komut kayıtlı değilse hata göster
	_, slashExists := events.GetRegisteredSlashCommands()[cmdName]
	_, prefixExists := events.GetRegisteredPrefixCommands()[cmdName]

	if !slashExists && !prefixExists {
		return false
	}

	return true
}

// komut adı boşsa genel bilgi, doluysa o komutun detayını gönderir.
func helpExecute(s *discordgo.Session, t utils.Target, cmdName string, helpCmd *library.CommandLib) {
	prefix := config.AppConfig.Bot.Prefix
	keyMap := map[string]string{
		"prefix":   prefix,
		"bot_name": config.AppConfig.Bot.Name,
	}

	if cmdName == "" {
		sendHelpGeneral(s, t, helpCmd, keyMap)
	} else {
		sendHelpForCommand(s, t, cmdName, helpCmd, keyMap)
	}
}

// Genel help mesajı — komut adı verilmeden çağrıldığında
func sendHelpGeneral(s *discordgo.Session, t utils.Target, helpCmd *library.CommandLib, keyMap map[string]string) {
	infoIcon := helpCmd.GetIcon("info", "")

	embedOpts := utils.BuildCommandEmbedFast(
		s,
		*helpCmd,
		library.Colors.Status.Info,
		"info",
		"general_info",
		helpCmd.Footers["default"],
		t.GetGuildID(),
		infoIcon,
		keyMap,
	)

	// Açıklama içindeki {prefix} vb. yerleştir
	embedOpts.Description = utils.KeyFormat(embedOpts.Description, keyMap)

	utils.SendEmbedToChannel(s, t, embedOpts, config.AppConfig.Bot.Name)
}

// Belirli bir komutun detay mesajı
func sendHelpForCommand(s *discordgo.Session, t utils.Target, cmdName string, helpCmd *library.CommandLib, keyMap map[string]string) {
	// Komutun kendi yaml'ını yükle
	var targetLib library.CommandLib
	var isYamlReadOk bool
	var yamlReadError error

	if utils.ValidateStringAlphanumeric(cmdName) != nil {
		// Komut adı Alfanumerik formatta değil , olası saldırı.
		utils.LogToConsole(utils.DEBUG, fmt.Sprintf("Şüpheli komut adı denemesi: %s", cmdName))
		sendHelpNotFound(s, t, cmdName, helpCmd, keyMap)
		return
	}

	// Core komut yolunu dene
	// ! Komut adı ile .yaml dosya adı aynı olmalı yoksa hatalı kabul eder
	corePath, err := utils.ValidateCommandYamlFilepathStrict(cmdName, config.AppConfig.Bot.Lang, true)
	if err == nil {
		isYamlReadOk, yamlReadError = utils.ReadYaml(corePath, &targetLib)
	}

	// Core'da bulunamadıysa normal yolunu dene
	// ! Komut adı ile .yaml dosya adı aynı olmalı yoksa hatalı kabul eder
	if !isYamlReadOk {
		normalPath, err := utils.ValidateCommandYamlFilepathStrict(cmdName, config.AppConfig.Bot.Lang, false)
		if err == nil {
			isYamlReadOk, yamlReadError = utils.ReadYaml(normalPath, &targetLib)
		}
	}

	// Registry'de de kontrol et — yaml yoksa veya komut kayıtlı değilse hata göster
	if yamlReadError != nil || !isYamlReadOk || !isAllowedCommandName(cmdName) {
		sendHelpNotFound(s, t, cmdName, helpCmd, keyMap)
		return
	}

	// Field'leri oluştur
	fields := buildCommandDetailFields(&targetLib, helpCmd, keyMap)

	infoIcon := helpCmd.GetIcon("info", "")

	embedOpts := utils.BuildCommandEmbedFast(
		s,
		*helpCmd,
		library.Colors.Status.Info,
		"info",
		"info",
		config.AppConfig.Bot.Name,
		t.GetGuildID(),
		infoIcon,
		keyMap,
		fields...,
	)

	// Başlık ve açıklamayı hedef komuttan al
	embedOpts.Title = targetLib.GetTitle("info", helpCmd.GetTitle("info", "📖 Komut Bilgisi"))
	embedOpts.Description = utils.KeyFormat(targetLib.Description, keyMap)

	if t.Interaction != nil {
		utils.SendEmbedToChannel(s, t, embedOpts, config.AppConfig.Bot.Name)
	} else {
		utils.SendEmbedReplyToChannel(s, t.Message, embedOpts, config.AppConfig.Bot.Name)
	}

}

// Komut bulunamadı embed'i
func sendHelpNotFound(s *discordgo.Session, t utils.Target, cmdName string, helpCmd *library.CommandLib, keyMap map[string]string) {
	val, _ := utils.GetCommonIcon("error")
	errIcon := utils.GetOrDefault(val, "")
	keyMap["komut"] = cmdName

	embedOpts := utils.BuildCommandEmbedFast(
		s,
		*helpCmd,
		library.Colors.Status.Error,
		"error",
		"err_not_found",
		config.AppConfig.Bot.Name,
		t.GetGuildID(),
		errIcon,
		keyMap,
	)
	embedOpts.Description = utils.KeyFormat(embedOpts.Description, keyMap)

	if t.Interaction != nil {
		utils.SendEmbedToChannel(s, t, embedOpts, config.AppConfig.Bot.Name)
	} else {
		utils.SendEmbedReplyToChannel(s, t.Message, embedOpts, config.AppConfig.Bot.Name)
	}
}

// Komutun usage, examples field'lerini oluşturur
func buildCommandDetailFields(targetCmd *library.CommandLib, helpCmd *library.CommandLib, keyMap map[string]string) []*discordgo.MessageEmbedField {
	var fields []*discordgo.MessageEmbedField

	// Kullanım
	if targetCmd.Usage != "" {

		usageField := helpCmd.GetField(1) // USAGE Fieldi
		targetUsageText := utils.KeyFormat(targetCmd.Usage, keyMap)

		fields = append(fields, &discordgo.MessageEmbedField{
			Name: usageField.Name,
			Value: utils.KeyFormat(usageField.Value, map[string]string{ // "``{command_usage}``"
				"command_usage": targetUsageText,
			}),
			Inline: usageField.Inline, // false
		})
	}

	// Örnekler
	if len(targetCmd.Examples) > 0 {
		var exLines []string
		for i := 1; i <= len(targetCmd.Examples); i++ {
			key := fmt.Sprintf("%d", i)
			if ex, ok := targetCmd.Examples[key]; ok {
				exLines = append(exLines, utils.KeyFormat(ex, keyMap))
			}
		}
		if len(exLines) > 0 {
			exFields := helpCmd.GetField(2) // EXAMPLES Fieldi
			fields = append(fields, &discordgo.MessageEmbedField{
				Name: exFields.Name,
				Value: utils.KeyFormat(exFields.Value, map[string]string{ // "```{command_examples}```"
					"command_examples": strings.Join(exLines, "\n"),
				}),
				Inline: exFields.Inline, // false
			})
		}
	}

	// Cooldown
	if cooldown, ok := targetCmd.GetOptionInt("cooldown"); ok && cooldown > 0 {

		coFields := helpCmd.GetField(3) // Cooldown Fieldi
		durationUnit, _ := utils.GetCommonDuration("second")
		resultText := utils.KeyFormat(coFields.Value, map[string]string{
			"duration": fmt.Sprintf("%d %s", cooldown, durationUnit.Full),
		})

		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   coFields.Name,
			Value:  resultText,
			Inline: coFields.Inline, // true
		})
	}

	categories, _ := library.GetOptionsExtraMap[string](helpCmd, "categories")
	// Kategori
	if cat, ok := targetCmd.GetOptionString("category"); ok && cat != "" {

		cat = utils.GetOrDefault(categories[cat], "Diğer")

		catFields := helpCmd.GetField(4) // Category Fieldi

		resultText := utils.KeyFormat(catFields.Value, map[string]string{
			"category": cat,
		})

		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   catFields.Name,
			Value:  resultText,
			Inline: catFields.Inline, // true
		})
	} else {

		// Fallback kategorisi
		cat = utils.GetOrDefault(categories["other"], "Diğer")

		catFields := helpCmd.GetField(4) // Category Fieldi

		resultText := utils.KeyFormat(catFields.Value, map[string]string{
			"category": cat,
		})

		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   catFields.Name,
			Value:  resultText,
			Inline: catFields.Inline, // true
		})
	}

	return fields
}
