package utils

import (
	library "ByteBunny/Bot/Core/Library"

	"github.com/bwmarrin/discordgo"
)

// Discord tarafından gelen tetikleme kaynağını ve cevap durumlarını barındıran taslak.
type Target struct {
	// Tetikleme tipine göre bu iki nesneden biri dolu, diğeri nil (boş) olur.
	Message     *discordgo.MessageCreate
	Interaction *discordgo.InteractionCreate
	IsEphemeral bool // Yanıtın sadece tetikleyen kişi tarafından mı görüleceğini belirler (Slash komutları için geçerlidir)
	IsFollowup  bool // Ephemeral mesajlarında ilk "Düşünüyor..." mesajından sonra mesaj göndermek için kullanılır
}

// Tetiklemeyi yapan sunucunun ID'sini döndürür
func (t *Target) GetGuildID() string {
	if t.Interaction != nil {
		// Tip Slash komutu
		return t.Interaction.GuildID
	}

	if t.Message != nil {
		// Tip Normal mesaj
		return t.Message.GuildID
	}

	return ""
}

// Tetikleme tipine ve kaynağına göre mesajı kanala gönderir
func SendEmbedToChannel(s *discordgo.Session, t Target, opts library.EmbedOptions, botName string) (err error, ok bool) {

	botAvatarUrl := getBotAvatar(s)
	if t.Interaction != nil {
		// Hedef içeriği Slash komutundan geliyor

		if t.IsFollowup {
			// İlk düşünüyor mesajından sonraki mesajı gönderir
			_, err := s.FollowupMessageCreate(t.Interaction.Interaction, t.IsEphemeral, &discordgo.WebhookParams{
				Embeds: []*discordgo.MessageEmbed{
					BuildEmbed(opts, botAvatarUrl, botName),
				},
			})
			return err, err == nil
		}

		flags := discordgo.MessageFlags(0)
		if t.IsEphemeral {
			flags = discordgo.MessageFlagsEphemeral
		}

		// Direkt Slash komutuna mesaj gönderir
		err := s.InteractionRespond(t.Interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{BuildEmbed(opts, botAvatarUrl, botName)},
				Flags:  flags,
			},
		})

		return err, err == nil

	} else if t.Message != nil {
		// Hedef içeriği Normal Mesaj tarafından geliyor.
		msg, err := s.ChannelMessageSendEmbed(t.Message.ChannelID, BuildEmbed(opts, botAvatarUrl, botName))

		return err, msg != nil
	}

	return nil, false
}

// Mesaj kaynağına göre yanıt mesajını kanala gönderir
func SendEmbedReplyToChannel(s *discordgo.Session, m *discordgo.MessageCreate, opts library.EmbedOptions, botName string) (err error, ok bool) {
	botAvatarUrl := getBotAvatar(s)

	msg, err := s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Embed:     BuildEmbed(opts, botAvatarUrl, botName),
		Reference: m.MessageReference,
	})

	return err, msg != nil
}

// Belirtilen kullanıcı kimliğine göre kullanıcıya direkt mesaj gönderir
func SendEmbedToDM(s *discordgo.Session, userId string, opts library.EmbedOptions, botName string) (err error, ok bool) {

	dm, err := s.UserChannelCreate(userId)
	if err != nil {
		return err, false
	}
	botAvatarUrl := getBotAvatar(s)
	msg, err := s.ChannelMessageSendEmbed(dm.ID, BuildEmbed(opts, botAvatarUrl, botName))

	return err, msg != nil
}

func getBotAvatar(s *discordgo.Session) string {
	return GetOrDefault(GetBotAvatarURL(s), common.Icons["bot"])
}
