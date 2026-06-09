package main

import (
	config "ByteBunny/Bot/Config"
	library "ByteBunny/Bot/Core/Library"
	utils "ByteBunny/Bot/Core/Utils"
	"os"
)

func init() {
	// Önyükleme noktası.

	// ── Uygulama ayarları ──────────────────────────────────────
	if config.InitConfigApp() != nil {
		utils.LogToConsole(utils.ERROR, "Config bilgileri yüklenemedi.")
		os.Exit(1)
	}

	// ── Uygulama renkleri ──────────────────────────────────────
	if library.InitColor("Config/Colors.yaml") != nil {
		utils.LogToConsole(utils.ERROR, "Config bilgileri yüklenemedi.")
		os.Exit(1)
	}

}

func main() {
	// Başlangıç noktası.
}
