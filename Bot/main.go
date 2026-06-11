package main

import (
	config "ByteBunny/Bot/Config"
	events "ByteBunny/Bot/Core/Events"
	general "ByteBunny/Bot/Core/General"
	library "ByteBunny/Bot/Core/Library"
	utils "ByteBunny/Bot/Core/Utils"
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func init() {
	// Önyükleme noktası.

	// ── Environment yükle ───────────────────────────────────────
	if err := godotenv.Load(); err != nil {
		panic(fmt.Errorf(".env dosyası okunamadı veya bulunamadı: %w", err))
	}

	// ── Uygulama ayarları ───────────────────────────────────────
	if e := config.InitConfigApp(); e != nil {
		utils.LogToConsole(utils.ERROR, fmt.Sprintf("Config bilgileri yüklenemedi: %v", e))
		os.Exit(1)
	}

	// ── Uygulama renkleri ───────────────────────────────────────
	if e := library.InitColor("Config/Colors.yaml"); e != nil {
		utils.LogToConsole(utils.ERROR, fmt.Sprintf("Config bilgileri yüklenemedi: %v", e))
		os.Exit(1)
	}

	// ── Ortak kullanım içerikleri ───────────────────────────────
	if e := utils.InitCommonDecoder(fmt.Sprintf("Modules/Locales/%s/Common.yaml",
		config.AppConfig.Bot.Lang)); e != nil {
		utils.LogToConsole(utils.ERROR, fmt.Sprintf("Config bilgileri yüklenemedi: %v", e))
		os.Exit(1)
	}

}

func loadAfterBot(s *discordgo.Session) {

	// ── Ortak ───────────────────────────────
	botName := config.AppConfig.Bot.Name
	botPrefix := config.AppConfig.Bot.Prefix
	botLang := config.AppConfig.Bot.Lang

	// ── Bot aktivite değiştirici ───────────────────────────────
	general.StartRandomPresence(s, botName, botPrefix, botLang)

	// ── Discord Mesajı Event Bus ───────────────────────────────
	s.AddHandler(events.DiscordHandlerBus.PublishMessage)
	s.AddHandler(events.DiscordHandlerBus.PublishInteraction)
}

func main() {
	// Başlangıç noktası.

	// loadAfterBot(session)
}
