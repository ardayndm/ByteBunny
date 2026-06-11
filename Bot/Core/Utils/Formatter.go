package utils

import (
	library "ByteBunny/Bot/Core/Library"
	"strings"
)

// Verilen bir metin şablonu (mainText) içerisindeki "{key}" biçimindeki
// dinamik işaretçileri, gönderilen haritadaki (keyMap) gerçek karşılıklarıyla toplu olarak değiştirir.
//
// Örnek Kullanım:
//
//	KeyFormat("Merhaba {user}, {guild} sunucusundasın!", map[string]string{
//	    "user":  "ByteBunny",
//	    "guild": "ByteBunnyDevs",
//	})
//	Çıktı → "Merhaba ByteBunny, ByteBunnyDevs sunucusundasın!"
func KeyFormat(mainText string, keyMap map[string]string) (formatted string) {

	result := mainText

	for key, value := range keyMap {
		result = strings.ReplaceAll(result, "{"+key+"}", value)
	}
	return result
}

// Verilen renk key'ini Colors.yaml'daki karşılığına çevirir.
// Bilinmeyen bir key girilirse 0 döner.
//
// Örnek Kullanım:
//
//	embed.Color = ColorFormatter("color_error")   // → 0xCF3F3F
//	embed.Color = ColorFormatter("color_success") // → 0x4FA882
//	embed.Color = ColorFormatter("color_deep")    // → 0x0D0D0D
func ColorFormat(colorKey string) int {

	colorMap := map[string]int{
		// Durum renkleri
		"color_error":   library.Colors.Status.Error,
		"color_warning": library.Colors.Status.Warn,
		"color_success": library.Colors.Status.Success,
		"color_info":    library.Colors.Status.Info,
		"color_main":    library.Colors.MainColor,
	}

	value, ok := colorMap[colorKey]
	if !ok {
		return library.Colors.MainColor // fallback
	}
	return value
}
