package utils

import (
	library "ByteBunny/Bot/Core/Library"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Doğrudan kabul edilen MessageCreate formatında Embed mesajı oluşturur.
func BuildEmbed(opts library.EmbedOptions, botAvatarUrl, botName string) *discordgo.MessageEmbed {

	// Url adreslerini kontrol et, hatalıları ayıkla
	checkUrlFormats(&opts)

	embed := &discordgo.MessageEmbed{
		Title:       opts.Title,
		Description: opts.Description,
		Color:       opts.Color,
		Fields:      opts.Fields,
		URL:         opts.URL,
	}

	// Footer (En alt sol kısım) alanını ve botun profil resmini ayarla
	embed.Footer = &discordgo.MessageEmbedFooter{
		IconURL: botAvatarUrl,
		Text:    GetOrDefault(botName, "ByteBunny"),
	}

	if opts.FooterText != "" {
		embed.Footer.Text += " • " + opts.FooterText
	}

	// Author (En üst sol kısım) alanı tanımlandıysa ayarla
	if opts.AuthorName != "" {
		embed.Author = &discordgo.MessageEmbedAuthor{
			Name: opts.AuthorName,
		}

		// AuthorName olmadan resim gözükmeyeceği için AuthorName altında kontrol ediliyor
		if opts.AuthorIconURL != "" {
			embed.Author.IconURL = opts.AuthorIconURL
		}
	}

	// Thumbnail (Sağ üst köşe görseli) tanımlandıysa ekle
	if opts.ThumbnailURL != "" {
		embed.Thumbnail = &discordgo.MessageEmbedThumbnail{
			URL: opts.ThumbnailURL,
		}
	}

	// Image (Orta büyük afiş görseli) tanımlandıysa ekle
	if opts.ImageURL != "" {
		embed.Image = &discordgo.MessageEmbedImage{
			URL: opts.ImageURL,
		}
	}

	return embed
}

// Hızlı şekilde Embed mesajı oluşturur.
//
//	Tamamen Common.yaml içindekiler ile işlemini yapar !
func BuildCommonEmbedFast(s *discordgo.Session, color int, title, message, customFooterText,
	guildId, thumbIconName string, fields ...*discordgo.MessageEmbedField) library.EmbedOptions {

	embedOpt := library.EmbedOptions{
		Title:        title,
		Description:  message,
		Color:        color,
		ThumbnailURL: GetOrDefault(common.Icons[thumbIconName], common.Icons["bug"]),
		FooterText:   customFooterText,
	}

	if len(fields) > 0 {
		embedOpt.Fields = fields
	}

	// Eğer geçerli bir sunucu ID'si varsa, Embed'ın tepesine sunucu adını yazıyoruz.
	if guildName, err := GetGuildName(s, guildId); err == nil && guildName != "" {
		embedOpt.AuthorName = guildName

		// AuthorName boş ise ikon zaten görünmez o yüzden AuthorName ile beraber atanıyor.
		// Sol üst köşedeki ikon için sunucu simgesini, yoksa dil dosyasındaki varsayılan bot ikonunu atıyoruz.
		embedOpt.AuthorIconURL = GetOrDefault(GetGuildIconFromID(s, guildId), common.Icons["bot"])
	}
	return embedOpt
}

// Hızlı şekilde Komut kütüphanesini kullanarak Embed mesajı oluşturur.
//
//	Önceliği Komut ile yapmaya çalışır , hata oluşursa fallback olarak Common'a başvurur
func BuildCommandEmbedFast(s *discordgo.Session,
	cmd library.CommandLib, color int, titleKey, messageKey,
	customFooterText, guildId, iconUrl string, keyMap map[string]string,
	fields ...*discordgo.MessageEmbedField) library.EmbedOptions {

	title := KeyFormat(GetOrDefault(cmd.Titles[titleKey], common.Titles["default"]), keyMap)
	desc := KeyFormat(GetOrDefault(cmd.Messages[messageKey], common.Messages["no_content"]), keyMap)
	embedOpt := library.EmbedOptions{
		Title:        title,
		Description:  desc,
		Color:        color,
		ThumbnailURL: GetOrDefault(iconUrl, common.Icons["bug"]),
		FooterText:   customFooterText,
	}

	if len(fields) > 0 {
		embedOpt.Fields = fields
	}

	// Eğer geçerli bir sunucu ID'si varsa, Embed'ın tepesine sunucu adını yazıyoruz.
	if guildName, err := GetGuildName(s, guildId); err == nil && guildName != "" {
		embedOpt.AuthorName = guildName

		// AuthorName boş ise ikon zaten görünmez o yüzden AuthorName ile beraber atanıyor.
		// Sol üst köşedeki ikon için sunucu simgesini, yoksa dil dosyasındaki varsayılan bot ikonunu atıyoruz.
		embedOpt.AuthorIconURL = GetOrDefault(GetGuildIconFromID(s, guildId), common.Icons["bot"])
	}
	return embedOpt
}

// Hızlı şekilde Komut için Field alanını discord'a göre düzenler
func RebuildCommandField(field *library.CommandField) *discordgo.MessageEmbedField {
	return &discordgo.MessageEmbedField{
		Name:   field.Name,
		Value:  field.Value,
		Inline: field.Inline,
	}
}

// Hızlı şekilde Komut için Field alanlarını discord'a göre düzenler
func RebuildCommandFields(fields ...*library.CommandField) []*discordgo.MessageEmbedField {
	var list []*discordgo.MessageEmbedField

	for _, f := range fields {
		list = append(list, &discordgo.MessageEmbedField{
			Name:   f.Name,
			Value:  f.Value,
			Inline: f.Inline,
		})
	}
	return list
}

// (internal) | Verilen Embed mesaj içeriğindeki URL adreslerini kontrol eder ve ayıklar
func checkUrlFormats(opts *library.EmbedOptions) {

	// Sol üst köşedeki ikon url adresini kontrol et
	// Url adresi hatalı ise discord mesajı kabul etmez ve hata oluşur
	if opts.AuthorIconURL != "" {
		// http&s:// ile başladığından emin ol
		if !strings.HasPrefix(opts.AuthorIconURL, "http://") &&
			!strings.HasPrefix(opts.AuthorIconURL, "https://") {

			opts.AuthorIconURL = ""
			LogToConsole(WARN, "Embed Mesaj | Author URL geçersiz formatta, mesajdan kaldırıldı. (http:// veya https:// ile başlamalı)")
		}
	}

	// Sağ üst köşedeki ikon url adresini kontrol et
	// Url adresi hatalı ise discord mesajı kabul etmez ve hata oluşur
	if opts.ThumbnailURL != "" {
		// http&s:// ile başladığından emin ol
		if !strings.HasPrefix(opts.ThumbnailURL, "http://") &&
			!strings.HasPrefix(opts.ThumbnailURL, "https://") {

			opts.ThumbnailURL = ""
			LogToConsole(WARN, "Embed Mesaj | Thumbnail URL geçersiz formatta, mesajdan kaldırıldı. (http:// veya https:// ile başlamalı)")
		}
	}

	// Büyük resim url adresini kontrol et
	// Url adresi hatalı ise discord mesajı kabul etmez ve hata oluşur
	if opts.ImageURL != "" {
		// http&s:// ile başladığından emin ol
		if !strings.HasPrefix(opts.ImageURL, "http://") &&
			!strings.HasPrefix(opts.ImageURL, "https://") {

			opts.ImageURL = ""
			LogToConsole(WARN, "Embed Mesaj | Image URL geçersiz formatta, mesajdan kaldırıldı. (http:// veya https:// ile başlamalı)")
		}
	}
}
