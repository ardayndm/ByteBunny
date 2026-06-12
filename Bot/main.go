package main

import (
	config "ByteBunny/Bot/Config"
	_ "ByteBunny/Bot/Core" // Komutlar için gerekli, aksi takdirde komutlar otomatik yüklenemez.
	events "ByteBunny/Bot/Core/Events"
	general "ByteBunny/Bot/Core/General"
	library "ByteBunny/Bot/Core/Library"
	preload "ByteBunny/Bot/Core/Preload"
	utils "ByteBunny/Bot/Core/Utils"
	"fmt"
	"os"
	"os/signal"
	"syscall"
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

	// ── Sisteme kayıt olmak için bekleyen Cache içindeki komutları yükle ───────────────────────────────
	events.LoadAllRegistryCommands()

	// ───────────────── Tüm komutlar sisteme kayıt olduktan sonra ─────────────────

	// ── Komutları Cache'ler ───────────────────────────────
	go preload.PreloadAllCommands()

}

func main() {
	// ── Uygulama ───────────────────────────────

	BotSession := startBot(loadAfterBot) // Botu başlatır
	defer BotSession.Close()             // Uygulama çökerse güvenli şekilde bağlantıyı kapat

	waitGracefullyShutdown() // Sisteme kapatma komutu geldiği zaman güvenle işler ve sistemi kapatır.

	/*

		? En son   Help.go yazıldı
		Todo: Help.yaml yazacaksın sonra test edeceksin, (Go tarafındaki keylere bak)
		Todo: Help komutundan sonra listcommands yazacaksın
		* Bu iki komut Core komut olacak sonra tüm komutlar
		* Modules altında olacak.

	*/

}

func startBot(onReady func(s *discordgo.Session)) *discordgo.Session {
	bot := config.AppConfig.Bot
	utils.LogToConsole(utils.INFO, fmt.Sprintf("%s başlatılıyor..", bot.Name))

	// ── Botu başlat ───────────────────────────────
	BotSession, err := discordgo.New("Bot " + bot.Token)

	if err != nil {
		utils.LogToConsole(utils.ERROR, fmt.Sprintf("Discord Bot oturumu oluşturulurken hata oluştu: %v", err))
		os.Exit(1)
	}

	// ── Bot hazır olduğu zaman çağrılır ───────────────────────────────
	BotSession.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		onReady(s)
	})

	if err := BotSession.Open(); err != nil {
		utils.LogToConsole(utils.ERROR, fmt.Sprintf("Discord Web Socket bağlantısı kurarken hata oluştu: %v", err))
		os.Exit(1)
	}

	utils.LogToConsole(utils.OK, fmt.Sprintf("%s Artık aktif.", bot.Name))
	return BotSession
}

func loadAfterBot(s *discordgo.Session) {

	utils.LogToConsole(utils.INFO, "Bot hazır, komutlar senkronize ediliyor...")

	// ── Ortak ───────────────────────────────
	botName := config.AppConfig.Bot.Name
	botPrefix := config.AppConfig.Bot.Prefix
	botLang := config.AppConfig.Bot.Lang

	// ── Bot aktivite değiştirici ───────────────────────────────
	general.StartRandomPresence(s, botName, botPrefix, botLang)

	// ── Discord Mesajı Event Bus ───────────────────────────────
	s.AddHandler(events.DiscordHandlerBus.PublishMessage)
	s.AddHandler(events.DiscordHandlerBus.PublishInteraction)

	// ── Bot Slash Komutları ───────────────────────────────
	if err := events.SyncSlashCommands(s); err != nil {
		utils.LogToConsole(utils.ERROR, fmt.Sprintf("Slash komutlar kayıt edilirken hata oluştu: %v", err))
		os.Exit(1)
	}

	utils.LogToConsole(utils.OK, "[──────────────── Sistem Tamamen Aktif ────────────────]")
}

func waitGracefullyShutdown() {
	// Ctrl+C (SIGINT) veya sistem kapatma sinyalleri (SIGTERM)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	<-stop // Kapatma sinyali gelene kadar kod bu satırda bloke olur ve bekler.

	utils.LogToConsole(utils.INFO, "[──────────────── Sistem Kapatılıyor ────────────────]")

	/*

		? Varsa veritabanı gibi güvenle kapatılması gereken işlemler buraya

	*/

	utils.LogToConsole(utils.INFO, "[──────────────── Güle Güle ────────────────]")
}
