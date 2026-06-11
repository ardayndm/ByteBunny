package events

import (
	config "ByteBunny/Bot/Config"
	"strings"
	"sync"

	"github.com/bwmarrin/discordgo"
)

// Eventleri yöneten merkezi birim.
type Dispatcher struct {
	mu sync.RWMutex

	messageHandlers     []MessageHandler
	interactionHandlers []InteractionHandler
}

// Global event bus.
var DiscordHandlerBus = newDispatcher()

// newDispatcher, HandleMessage ve HandleInteraction'ı otomatik register ederek başlatır.
func newDispatcher() *Dispatcher {
	d := &Dispatcher{}
	d.RegisterMessage(handlePrefixMessage)
	d.RegisterInteraction(handleInteraction)
	return d
}

// ── Mesaj Kayıt Fonksiyonları ───────────────────────────────────────────────────

// Yeni bir mesaj dinleyicisi kaydeder.
func (d *Dispatcher) RegisterMessage(handler MessageHandler) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.messageHandlers = append(d.messageHandlers, handler)
}

// Tüm mesaj dinleyicilerine mesajı iletir.
func (d *Dispatcher) PublishMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	for _, handler := range d.messageHandlers {
		go handler(s, m)
	}
}

// ── Slash Kayıt Fonksiyonları ───────────────────────────────────────────────────

// Yeni bir slash dinleyicisi kaydeder.
func (d *Dispatcher) RegisterInteraction(handler InteractionHandler) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.interactionHandlers = append(d.interactionHandlers, handler)
}

// Tüm slash dinleyicilerine interaction'ı iletir.
func (d *Dispatcher) PublishInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	for _, handler := range d.interactionHandlers {
		go handler(s, i)
	}
}

// ── Handler fonksiyonlar ───────────────────────────────────────────────────

// Gelen slash komutunu registry'den bulup çalıştırır.
func handleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	name := i.ApplicationCommandData().Name

	slashMutex.RLock()
	cmd, exists := slashRegistry[name] // Komut ismine sahip kayıtlı fonksiyonu al
	slashMutex.RUnlock()

	if !exists {
		return
	}

	cmd.Execute(s, i)
}

// Gelen mesajı prefix kontrolünden geçirip ilgili komutu çalıştırır.
func handlePrefixMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	prefix := config.AppConfig.Bot.Prefix
	if !strings.HasPrefix(m.Content, prefix) {
		return
	}

	args := strings.Fields(strings.TrimPrefix(m.Content, prefix))
	if len(args) == 0 {
		return
	}

	prefixMutex.RLock()
	cmd, exists := prefixRegistry[args[0]] // Komut ismine sahip kayıtlı fonksiyonu al
	prefixMutex.RUnlock()

	if !exists {
		return
	}

	cmd.Execute(s, m, args)
}
