package utils

import "strings"

// Verilen bir metin şablonu (mainText) içerisindeki "{key}" biçimindeki
// dinamik işaretçileri, gönderilen haritadaki (keyMap) gerçek karşılıklarıyla toplu olarak değiştirir.
//
// Örnek Kullanım:
//   KeyFormat("Merhaba {user}, {guild} sunucusundasın!", map[string]string{
//       "user":  "ByteBunny",
//       "guild": "ByteBunnyDevs",
//   })
//   Çıktı → "Merhaba ByteBunny, ByteBunnyDevs sunucusundasın!"
func KeyFormat(mainText string, keyMap map[string]string) (formatted string) {

	result := mainText

	for key, value := range keyMap {
		result = strings.ReplaceAll(result, "{"+key+"}", value)
	}
	return result
}
