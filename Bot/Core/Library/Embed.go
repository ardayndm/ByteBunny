package library

import "github.com/bwmarrin/discordgo"

// Bir embed mesajı oluşturmak için gerekli tüm parametreleri toplayan şablon.
type EmbedOptions struct {
	Title         string
	Description   string
	Color         int
	Fields        []*discordgo.MessageEmbedField
	FooterText    string // En altta görünecek dipnot metni
	FooterIconURL string // Footer ikonu — genellikle bot avatarı
	ThumbnailURL  string // Sağ üst köşedeki küçük kare görsel
	AuthorName    string // Sol üst köşedeki yazar adı
	AuthorIconURL string // Sol üst köşedeki yazar ikonu (AuthorName dolu olmalıdır)
	ImageURL      string // Embed altındaki tam boy görsel
	URL           string // Başlığa tıklandığında yönlendirilecek adres
}
