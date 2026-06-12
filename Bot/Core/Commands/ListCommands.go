package commands

import (
	config "ByteBunny/Bot/Config"
	events "ByteBunny/Bot/Core/Events"
	library "ByteBunny/Bot/Core/Library"
	preload "ByteBunny/Bot/Core/Preload"
	utils "ByteBunny/Bot/Core/Utils"
	"fmt"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/bwmarrin/discordgo"
)

// ── Tipler ───────────────────────────────────────────────────────────────────────
type ListCommandsSlash struct{ cmd *library.CommandLib }
type ListCommandsPrefix struct{ cmd *library.CommandLib }

var (
	slashOnceListCommand  sync.Once
	prefixOnceListCommand sync.Once
)

// ── Oto kayıt ───────────────────────────────────────────────────────────────────────

func init() {
	events.RegisterSlashCacheCommand(&ListCommandsSlash{})
	events.RegisterPrefixCacheCommand(&ListCommandsPrefix{})
}

// ── Yükle ────────────────────────────────────────────────────────────────────

func loadListCommandsLib() (*library.CommandLib, error) {
	var lib library.CommandLib
	path := filepath.Join("Modules", "Locales", config.AppConfig.Bot.Lang, "Commands", "Core", "listcommands.yaml")
	ok, err := utils.ReadYaml(path, &lib)
	if err != nil || !ok {
		return nil, err
	}
	return &lib, nil
}

// ── Kontrol ────────────────────────────────────────────────────────────────────

func (lcSlash *ListCommandsSlash) checkLibrary() bool {
	if lcSlash.cmd != nil {
		return true
	}
	return lcSlash.LoadCommandLib() == nil
}

func (lcPrefix *ListCommandsPrefix) checkLibrary() bool {
	if lcPrefix.cmd != nil {
		return true
	}
	return lcPrefix.LoadCommandLib() == nil
}

// ── Slash Komutu ────────────────────────────────────────────────────────────────────

func (lcSlash *ListCommandsSlash) LoadCommandLib() error {
	var loadErr error

	slashOnceListCommand.Do(func() {
		lib, err := loadListCommandsLib()
		if err != nil {
			loadErr = err
			return
		}
		lcSlash.cmd = lib
	})

	return loadErr
}

func (lcSlash *ListCommandsSlash) Name() string {
	if !lcSlash.checkLibrary() {
		return "listcommands"
	}
	return lcSlash.cmd.Name
}

func (lcSlash *ListCommandsSlash) Description() string {
	if !lcSlash.checkLibrary() {
		return "Tüm komutları listeler."
	}
	return lcSlash.cmd.Description
}

func (lcSlash *ListCommandsSlash) Options() []*discordgo.ApplicationCommandOption {
	if !lcSlash.checkLibrary() {
		return nil
	}

	var opts []*discordgo.ApplicationCommandOption
	for _, arg := range lcSlash.cmd.GetSlashArgs() {
		argType := discordgo.ApplicationCommandOptionString
		switch arg.Type {
		case "integer":
			argType = discordgo.ApplicationCommandOptionInteger
		case "boolean":
			argType = discordgo.ApplicationCommandOptionBoolean
		}

		opts = append(opts, &discordgo.ApplicationCommandOption{
			Name:        arg.Name,
			Description: arg.Description,
			Type:        argType,
			Required:    arg.Required,
		})
	}
	return opts
}

func (lcSlash *ListCommandsSlash) GetGeneralFormatKeys(extra ...map[string]string) map[string]string {
	return map[string]string{
		"prefix": "/",
	}
}

func (lcSlash *ListCommandsSlash) Execute(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !lcSlash.checkLibrary() || lcSlash.cmd == nil {
		utils.LogToConsole(utils.ERROR, "ListCommands: komut kütüphanesi yüklenemedi")
		return
	}

	if enabled, ok := lcSlash.cmd.GetOptionBool("enabled"); !enabled || !ok {
		return
	}

	utils.SendInteractionDeferRespond(s, i.Interaction)

	target := utils.Target{Interaction: i, IsEphemeral: true, IsFollowup: true}

	page := 1

	name, ok := lcSlash.cmd.GetOptionString("name")

	if !ok {
		name = "sayfa" // fallback
	}

	for _, opt := range i.ApplicationCommandData().Options {
		if opt.Name == name {
			page = int(opt.IntValue())
		}
	}
	if page < 1 {
		page = 1
	}

	listCommandsExecute(s, target, page, lcSlash.cmd)
}

// ──────────────────────────────────────────────────────────────────────
// ── Prefix Komutu ────────────────────────────────────────────────────────────────────

func (lcPrefix *ListCommandsPrefix) LoadCommandLib() error {
	var loadErr error

	prefixOnceListCommand.Do(func() {
		lib, err := loadListCommandsLib()
		if err != nil {
			loadErr = err
			return
		}
		lcPrefix.cmd = lib
	})

	return loadErr
}

func (lcPrefix *ListCommandsPrefix) Name() string {
	if !lcPrefix.checkLibrary() {
		return "listcommands"
	}
	return lcPrefix.cmd.Name
}

func (lcPrefix *ListCommandsPrefix) GetGeneralFormatKeys(extra ...map[string]string) map[string]string {
	return map[string]string{
		"prefix": config.AppConfig.Bot.Prefix,
	}
}

func (lcPrefix *ListCommandsPrefix) Execute(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if !lcPrefix.checkLibrary() || lcPrefix.cmd == nil {
		utils.LogToConsole(utils.ERROR, "ListCommands: komut kütüphanesi yüklenemedi")
		return
	}

	if enabled, ok := lcPrefix.cmd.GetOptionBool("enabled"); !enabled || !ok {
		return
	}

	target := utils.Target{Message: m}

	page := 1
	if len(args) > 1 {
		if p, err := strconv.Atoi(strings.TrimSpace(args[1])); err == nil {
			page = p
		}
	}
	if page < 1 {
		page = 1
	}

	listCommandsExecute(s, target, page, lcPrefix.cmd)
}

// ──────────────────────────────────────────────────────────────────────
// ── İç Mantık ────────────────────────────────────────────────────────────────────

// commandEntry tek bir komutun listelenmesi için gereken minimum bilgi
type commandEntry struct {
	Name        string
	Description string
	Category    string
}

func collectCommandEntries() []commandEntry {
	// Önce preload cache'ini kontrol et
	if preload.IsPreloadDone() {
		return collectFromCache() // Cache'den al
	}

	// Preload bitmediyse eski yöntemi kullan (güvenlik)
	return collectFromDisk()
}

// Cache'den oku
func collectFromCache() []commandEntry {
	cachedCommands := preload.GetAllCachedCommands()
	entries := make([]commandEntry, 0, len(cachedCommands))

	for _, cmd := range cachedCommands {
		// Sadece aktif ve gizli olmayan komutları göster
		if !cmd.Enabled || cmd.Hidden {
			continue
		}

		entries = append(entries, commandEntry{
			Name:        cmd.Name,
			Description: cmd.Description,
			Category:    cmd.Category,
		})
	}

	return entries
}

// Tüm kayıtlı (slash + prefix) komutları yaml'larından okuyup commandEntry listesine çevirir.
// Aynı isimde hem slash hem prefix kaydı varsa tek seferde sayılır.
func collectFromDisk() []commandEntry {
	seen := make(map[string]bool)
	var entries []commandEntry

	addFromRegistry := func(names map[string]bool) {
		for name := range names {
			if seen[name] {
				continue
			}
			seen[name] = true

			if utils.ValidateStringAlphanumeric(name) != nil {
				continue
			}

			var targetLib library.CommandLib
			var ok bool
			var err error

			// Önce core, sonra normal yol denenir.
			corePath, perr := utils.ValidateCommandYamlFilepathStrict(name, config.AppConfig.Bot.Lang, true)
			if perr == nil {
				ok, err = utils.ReadYaml(corePath, &targetLib)
			}

			if !ok {
				normalPath, perr := utils.ValidateCommandYamlFilepathStrict(name, config.AppConfig.Bot.Lang, false)
				if perr == nil {
					ok, err = utils.ReadYaml(normalPath, &targetLib)
				}
			}

			if err != nil || !ok {
				continue
			}

			// Listede gizlenmesi istenen komutları atla (enabled: false veya hidden: true)
			if enabled, has := targetLib.GetOptionBool("enabled"); has && !enabled {
				continue
			}
			if hidden, has := targetLib.GetOptionBool("hidden"); has && hidden {
				continue
			}

			cat, ok := targetLib.GetOptionString("category")
			if !ok || cat == "" {
				cat = "other"
			}

			entries = append(entries, commandEntry{
				Name:        targetLib.Name,
				Description: targetLib.Description,
				Category:    cat,
			})
		}
	}

	slashNames := make(map[string]bool)
	for name := range events.GetRegisteredSlashCommands() {
		slashNames[name] = true
	}
	addFromRegistry(slashNames)

	prefixNames := make(map[string]bool)
	for name := range events.GetRegisteredPrefixCommands() {
		prefixNames[name] = true
	}
	addFromRegistry(prefixNames)

	return entries
}

// Komutları kategoriye göre gruplar, kategori sırasını ve içindeki komutları
// alfabetik olarak sıralar.
func groupByCategory(entries []commandEntry, categories map[string]string) ([]string, map[string][]commandEntry) {
	grouped := make(map[string][]commandEntry)

	for _, e := range entries {
		key := e.Category
		if _, ok := categories[key]; !ok {
			key = "other"
		}
		grouped[key] = append(grouped[key], e)
	}

	for key := range grouped {
		sort.Slice(grouped[key], func(i, j int) bool {
			return strings.ToLower(grouped[key][i].Name) < strings.ToLower(grouped[key][j].Name)
		})
	}

	// Kategori sırası: categories map'inde tanımlı sıra + içinde komut olanlar
	var order []string
	for key := range categories {
		if len(grouped[key]) > 0 {
			order = append(order, key)
		}
	}
	sort.Strings(order)

	return order, grouped
}

// Tek bir kategori bloğunu satır satır metne çevirir.
// ── Kategori Adı ──
//
// komutadi    açıklama
func renderCategoryBlock(categoryLabel string, entries []commandEntry, lc *library.CommandLib, keyMap map[string]string) string {
	var b strings.Builder

	headerFmt := lc.GetMessage("category_header", "── {category} ──")
	itemFmt := lc.GetMessage("command_item", "**{command}** — {description}")

	header := utils.KeyFormat(headerFmt, map[string]string{"category": categoryLabel})
	b.WriteString(header)
	b.WriteString("\n")

	for _, e := range entries {
		desc := e.Description
		if desc == "" {
			desc = lc.GetMessage("no_description", "Açıklama yok.")
		}

		line := utils.KeyFormat(itemFmt, map[string]string{
			"command":     e.Name,
			"description": desc,
			"prefix":      keyMap["prefix"],
		})
		b.WriteString(line)
		b.WriteString("\n")
	}

	return b.String()
}

// Sayfalama: kategori blokları, embed'in karakter/field limitini aşmayacak
// şekilde "sayfa"lara bölünür. Her sayfa bir veya daha fazla kategori
// içerebilir; tek bir kategori dahi sığmıyorsa kendi içinde de bölünür.
const listCommandsMaxFieldLen = 1000 // Discord embed field value limiti 1024, güvenli pay bırakıyoruz

// Tüm kategorileri sırayla gezip, komutları 20'şerlik sayfalara böler.
// Bir kategori sayfa sınırını aşarsa, kalan komutlar bir sonraki sayfada
// aynı kategori adıyla (devam) etiketiyle devam eder.
func buildPages(order []string, grouped map[string][]commandEntry, categories map[string]string, lc *library.CommandLib, keyMap map[string]string) [][]*library.CommandField {
	itemFmt := lc.GetMessage("command_item", "• **{command}** — {description}")
	headerSuffix := lc.GetMessage("category_continued", " (devam)")

	var pages [][]*library.CommandField
	var currentFields []*library.CommandField
	var currentLines []string
	currentCount := 0
	currentLabel := ""
	currentIsContinued := false

	flushField := func() {
		if len(currentLines) == 0 {
			return
		}
		name := currentLabel
		if currentIsContinued {
			name += headerSuffix
		}
		currentFields = append(currentFields, &library.CommandField{
			Name:   name,
			Value:  strings.Join(currentLines, "\n"),
			Inline: false,
		})
		currentLines = nil
	}

	flushPage := func() {
		flushField()
		if len(currentFields) > 0 {
			pages = append(pages, currentFields)
			currentFields = nil
		}
		currentCount = 0
	}

	maxListCommandCount, ok := lc.GetOptionInt("command_count_per_page")

	if !ok {
		maxListCommandCount = 10 // fallback
	}

	for _, catKey := range order {
		label := utils.GetOrDefault(categories[catKey], categories["other"])
		currentLabel = label
		currentIsContinued = false

		for _, e := range grouped[catKey] {
			if currentCount == maxListCommandCount {
				flushField()
				flushPage()
				currentLabel = label
				currentIsContinued = true
			}

			desc := e.Description
			if desc == "" {
				desc = lc.GetMessage("no_description", "Açıklama yok.")
			}

			line := utils.KeyFormat(itemFmt, map[string]string{
				"command":     e.Name,
				"description": desc,
				"prefix":      keyMap["prefix"],
			})

			currentLines = append(currentLines, line)
			currentCount++
		}

		// Kategori bitti, field'i kapat (yeni kategori yeni field demek)
		flushField()
	}

	flushPage()

	if len(pages) == 0 {
		pages = append(pages, []*library.CommandField{})
	}

	return pages
}

// Ana işleyici
func listCommandsExecute(s *discordgo.Session, t utils.Target, page int, lc *library.CommandLib) {
	keyMap := map[string]string{
		"prefix":   config.AppConfig.Bot.Prefix,
		"bot_name": config.AppConfig.Bot.Name,
	}

	categories, _ := library.GetOptionsExtraMap[string](lc, "categories")
	if categories == nil {
		categories = map[string]string{"other": "❓ Diğer"}
	}

	entries := collectCommandEntries()
	order, grouped := groupByCategory(entries, categories)
	pages := buildPages(order, grouped, categories, lc, keyMap)

	totalPages := len(pages)

	if page > totalPages {
		page = totalPages
	}

	embedMessageKey := "info" // Default info mesajı.

	if page < 1 {
		page = 1
	}

	pageFields := pages[page-1]

	if len(pageFields) < 1 {
		embedMessageKey = "has_no_command_info"
	}

	infoIcon := lc.GetIcon("list", "")

	footerFmt := lc.Footers["page"]
	if footerFmt == "" {
		footerFmt = "Sayfa {page}/{total_pages}"
	}
	footer := utils.KeyFormat(footerFmt, map[string]string{
		"bot_name":    config.AppConfig.Bot.Name,
		"page":        fmt.Sprintf("%d", page),
		"total_pages": fmt.Sprintf("%d", totalPages),
	})

	embedOpts := utils.BuildCommandEmbedFast(
		s,
		*lc,
		library.Colors.Status.Info,
		"info",
		embedMessageKey,
		footer,
		t.GetGuildID(),
		infoIcon,
		keyMap,
		utils.RebuildCommandFields(pageFields...)...,
	)

	embedOpts.Description = utils.KeyFormat(embedOpts.Description, keyMap)

	if t.Interaction != nil {
		utils.SendEmbedToChannel(s, t, embedOpts, config.AppConfig.Bot.Name)
	} else {
		utils.SendEmbedReplyToChannel(s, t.Message, embedOpts, config.AppConfig.Bot.Name)
	}
}
