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
)

func init() {
	// Önyükleme noktası.

	// ── Uygulama ayarları ───────────────────────────────────────
	if config.InitConfigApp() != nil {
		utils.LogToConsole(utils.ERROR, "Config bilgileri yüklenemedi.")
		os.Exit(1)
	}

	// ── Uygulama renkleri ───────────────────────────────────────
	if library.InitColor("Config/Colors.yaml") != nil {
		utils.LogToConsole(utils.ERROR, "Config bilgileri yüklenemedi.")
		os.Exit(1)
	}

	// ── Ortak kullanım içerikleri ───────────────────────────────
	if utils.InitCommonDecoder(fmt.Sprintf("Modules/Locales/%s/Common.yaml",
		config.AppConfig.Bot.Lang)) != nil {
		utils.LogToConsole(utils.ERROR, "Config bilgileri yüklenemedi.")
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
