package utils

import (
	library "ByteBunny/Bot/Core/Library"
	"fmt"
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
func GetCommonTemplate(key string) (template library.EmbedTemplate, ok bool, e error) {
	if common == nil {
		LogToConsole(ERROR, "CommonDecoder: InitCommonDecoder henüz çağrılmadı")
		return library.EmbedTemplate{}, false, fmt.Errorf("CommonDecoder: InitCommonDecoder henüz çağrılmadı")
	}
	tmpl, ok := common.Templates[key]
	return tmpl, ok, nil
}

// commands altındaki bir embed şablonunu döner.
//
// Örnek:
//
//	tmpl, ok := utils.GetCommand("cooldown")
func GetCommonCommand(key string) (template library.EmbedTemplate, ok bool, e error) {
	if common == nil {
		LogToConsole(ERROR, "CommonDecoder: InitCommonDecoder henüz çağrılmadı")
		return library.EmbedTemplate{}, false, fmt.Errorf("CommonDecoder: InitCommonDecoder henüz çağrılmadı")
	}
	tmpl, ok := common.Commands[key]
	return tmpl, ok, nil
}

// titles altındaki varsayılan başlığı döner.
func GetCommonTitle(key string) (value string, ok bool, e error) {
	if common == nil {
		LogToConsole(ERROR, "CommonDecoder: InitCommonDecoder henüz çağrılmadı")
		return "", false, fmt.Errorf("CommonDecoder: InitCommonDecoder henüz çağrılmadı")

	}
	title, ok := common.Titles[key]
	return title, ok, nil
}

// messages altındaki varsayılan mesajı döner.
func GetCommonMessage(key string) (value string, ok bool, e error) {
	if common == nil {
		LogToConsole(ERROR, "CommonDecoder: InitCommonDecoder henüz çağrılmadı")
		return "", false, fmt.Errorf("CommonDecoder: InitCommonDecoder henüz çağrılmadı")
	}
	msg, ok := common.Messages[key]
	return msg, ok, nil
}

// icons altındaki URL'yi döner.
func GetCommonIcon(key string) (value string, ok bool, e error) {
	if common == nil {
		LogToConsole(ERROR, "CommonDecoder: InitCommonDecoder henüz çağrılmadı")
		return "", false, fmt.Errorf("CommonDecoder: InitCommonDecoder henüz çağrılmadı")
	}
	url, ok := common.Icons[key]
	return url, ok, nil
}

// duration altındaki zaman birimini döner.
func GetCommonDuration(key string) (unit library.DurationUnit, ok bool, e error) {
	if common == nil {
		LogToConsole(ERROR, "CommonDecoder: InitCommonDecoder henüz çağrılmadı")
		return library.DurationUnit{}, false, fmt.Errorf("CommonDecoder: InitCommonDecoder henüz çağrılmadı")
	}
	d, ok := common.Duration[key]
	return d, ok, nil
}
