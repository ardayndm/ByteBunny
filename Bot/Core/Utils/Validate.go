package utils

import (
	library "ByteBunny/Bot/Core/Library"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

// Verilen metni alfanumerik formatlarda kontrol eder.
func ValidateStringAlphanumeric(text string) error {
	if text == "" {
		return library.Err_InvalidName
	}
	// Sadece alfanumerik ve alt çizgi karakterlerine izin ver
	matched, _ := regexp.MatchString("^[a-zA-Z0-9_-]+$", text)
	if !matched {
		return library.Err_InvalidName
	}
	return nil
}

// Dil kodunu doğrular (örn: tr, en, de)
func ValidateLanguageCode(lang string) error {
	if lang == "" {
		return library.Err_InvalidName
	}
	// Sadece 2-3 karakterli harfler (tr, en, de, es, fr, etc.)
	matched, _ := regexp.MatchString("^[a-zA-Z]{2,3}$", lang)
	if !matched {
		return library.Err_InvalidName
	}
	return nil
}

// Komut YAML dosyalarının yolunu doğrular.
func ValidateCommandYamlFilepath(cmdName, botLang string, isCore bool) (resultPath string, e error) {

	// Dosya adının uzunluğunu kontrol et
	if len(cmdName) > 35 {
		return "", library.Err_InvalidCommandName
	}

	// BotLang uzunluğu kontrolü
	if len(botLang) > 3 {
		return "", library.Err_InvalidLanguage
	}

	// cmdName validasyonu
	if err := ValidateStringAlphanumeric(cmdName); err != nil {
		return "", library.Err_InvalidCommandName
	}

	// botLang validasyonu
	if err := ValidateLanguageCode(botLang); err != nil {
		return "", library.Err_InvalidLanguage
	}

	// Yol oluştur
	var path string
	if isCore {
		path = fmt.Sprintf("Modules/Commands/%s/Commands/Core/%s.yaml", botLang, cmdName)
	} else {
		path = fmt.Sprintf("Modules/Commands/%s/Commands/%s.yaml", botLang, cmdName)
	}

	// Path traversal saldırılarını temizle
	cleanPath := filepath.Clean(path)

	// Temizlenen yolda hala ../ var mı kontrol et
	if strings.Contains(cleanPath, "..") {
		return "", library.Err_SuspectedFilePath
	}

	// Base directory kontrolü
	baseDir := "Modules/Commands"
	if !strings.HasPrefix(cleanPath, baseDir+string(filepath.Separator)) {
		return "", library.Err_SuspectedFilePath
	}

	// Yolun yapısını validate et
	// Beklenen format: Modules/Commands/{lang}/Commands/[Core/]{cmdName}.yaml
	parts := strings.Split(cleanPath, string(filepath.Separator))
	if len(parts) < 4 {
		return "", library.Err_InvalidPath
	}

	// Modules/Commands/ kontrolü
	if parts[0] != "Modules" || parts[1] != "Commands" {
		return "", library.Err_InvalidPath
	}

	// Dil kodu kontrolü (tekrar)
	if err := ValidateLanguageCode(parts[2]); err != nil {
		return "", library.Err_InvalidLanguage
	}

	// Commands dizini kontrolü
	if parts[3] != "Commands" {
		return "", library.Err_InvalidPath
	}

	return cleanPath, nil
}

// Daha sıkı validasyon. (Dosya adının .yaml ile bittiğine kadar kontrol eder)
func ValidateCommandYamlFilepathStrict(cmdName, botLang string, isCore bool) (resultPath string, e error) {
	// Önce normal validasyonu yap
	cleanPath, err := ValidateCommandYamlFilepath(cmdName, botLang, isCore)
	if err != nil {
		return "", err
	}

	// Dosyanın gerçekten var olup olmadığını kontrol etme (opsiyonel)
	// Burada sadece path validasyonu yapıyoruz, dosya varlığını çağıran kontrol eder

	// Ekstra: Dosya adının .yaml ile bittiğinden emin ol
	if !strings.HasSuffix(cleanPath, ".yaml") {
		return "", library.Err_FileNotYaml
	}

	return cleanPath, nil
}
