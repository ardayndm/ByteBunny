package events

import "github.com/bwmarrin/discordgo"

type (
	MessageHandler     func(s *discordgo.Session, m *discordgo.MessageCreate)     // Discord tarafından gelen ; Gönderilen mesajlar.
	InteractionHandler func(s *discordgo.Session, i *discordgo.InteractionCreate) // Discord tarafından gelen ; Kullanıcıların / komutu ile gönderdiği mesajlar.
)
