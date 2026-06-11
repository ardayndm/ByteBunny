package utils

import (
	library "ByteBunny/Bot/Core/Library"
	"fmt"
	"os"
)

var common *library.CommonLib

// Common.yaml dosyasını okur ve common değişkenini doldurur.
func InitCommonDecoder(path string) error {
	var lib library.CommonLib
	if ok, err := ReadYaml(path, &lib); err != nil || !ok {
		LogToConsole(ERROR, fmt.Errorf("CommonDecoder hatası: %w", err).Error())
		return err
	}
	common = &lib

	return nil
}

// Templates altındaki bir embed şablonunu döner.
// Bulunamazsa ikinci değer false döner.
//
// Örnek:
//
//	tmpl, ok := utils.GetTemplate("error")
func GetCommonTemplate(key string) (template library.EmbedTemplate, ok bool) {
	if common == nil {
		LogToConsole(ERROR, "CommonDecoder: InitCommonDecoder henüz çağrılmadı")
		os.Exit(1)
	}
	tmpl, ok := common.Templates[key]
	return tmpl, ok
}

// commands altındaki bir embed şablonunu döner.
//
// Örnek:
//
//	tmpl, ok := utils.GetCommand("cooldown")
func GetCommonCommand(key string) (template library.EmbedTemplate, ok bool) {
	if common == nil {
		LogToConsole(ERROR, "CommonDecoder: InitCommonDecoder henüz çağrılmadı")
		os.Exit(1)
	}
	tmpl, ok := common.Commands[key]
	return tmpl, ok
}

// titles altındaki varsayılan başlığı döner.
func GetCommonTitle(key string) (value string, ok bool) {
	if common == nil {
		LogToConsole(ERROR, "CommonDecoder: InitCommonDecoder henüz çağrılmadı")
		os.Exit(1)
	}
	title, ok := common.Titles[key]
	return title, ok
}

// messages altındaki varsayılan mesajı döner.
func GetCommonMessage(key string) (value string, ok bool) {
	if common == nil {
		LogToConsole(ERROR, "CommonDecoder: InitCommonDecoder henüz çağrılmadı")
		os.Exit(1)
	}
	msg, ok := common.Messages[key]
	return msg, ok
}

// icons altındaki URL'yi döner.
func GetCommonIcon(key string) (value string, ok bool) {
	if common == nil {
		LogToConsole(ERROR, "CommonDecoder: InitCommonDecoder henüz çağrılmadı")
		os.Exit(1)
	}
	url, ok := common.Icons[key]
	return url, ok
}

// duration altındaki zaman birimini döner.
func GetCommonDuration(key string) (unit library.DurationUnit, ok bool) {
	if common == nil {
		LogToConsole(ERROR, "CommonDecoder: InitCommonDecoder henüz çağrılmadı")
		os.Exit(1)
	}
	d, ok := common.Duration[key]
	return d, ok
}
