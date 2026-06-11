package general

import (
	library "ByteBunny/Bot/Core/Library"
	utils "ByteBunny/Bot/Core/Utils"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
)

// (internal) | Aktivite durumunu YAML dosyasından yükler
//
//	(Hata oluşursa kendi yazdırır)
func initPresenceConfig(lang string) *library.PresenceConfig {

	var cfg library.PresenceConfig

	path := fmt.Sprintf("Core/Locale/%s/Activity.yaml", lang)

	if ok, err := utils.ReadYaml(path, &cfg); err != nil || !ok {
		utils.LogToConsole(utils.ERROR, "Activity.yaml dosyası okunamadı.")
		os.Exit(1)
	}

	return &cfg
}

// Botun discord üzerindeki aktivite durumunu rastgele şekilde ayarlara göre değiştirir.
func StartRandomPresence(s *discordgo.Session, botName, botPrefix, botLang string) {
	lang := botLang

	// Bot dili belirtilmemiş ise fallback olarak tr ayarla
	if lang == "" {
		lang = "tr"
	}

	// Durum ayarlarını okuyoruz; dosya yoksa sistemin çökmemesi için dahili varsayılanları devreye alıyoruz.
	cfg := initPresenceConfig(lang)

	if len(cfg.Activities) == 0 {
		utils.LogToConsole(utils.SKIP, "Presence aktivite havuzu boş olduğu için durum değiştirici iptal edildi, devam ediliyor..")
		return
	}

	// Aktivite isimlerindeki {bot_name} ve {prefix} şablonlarını dinamik verilerle dolduruyoruz.
	for i := range cfg.Activities {
		replacements := map[string]string{
			"bot_name": botName,
			"prefix":   botPrefix,
		}

		cfg.Activities[i].Name = utils.KeyFormat(cfg.Activities[i].Name, replacements)
		cfg.Activities[i].Details = utils.KeyFormat(cfg.Activities[i].Details, replacements)
		cfg.Activities[i].State = utils.KeyFormat(cfg.Activities[i].State, replacements)
	}

	status := cfg.Status

	// Durum belirtilmemiş ise fallback olarak online ayarla
	if status == "" {
		status = "online"
	}

	// Güncelleme süresi
	interval := time.Duration(cfg.Interval) * time.Second

	// Interval değeri 0 dan az ise fallback olarak 10 saniyede bir ayarla
	if interval <= 0 {
		interval = 10 * time.Second
	}

	ticker := time.NewTicker(interval)
	utils.LogToConsole(utils.OK, "Aktivite değiştirici başlatıldı.")

	go func() {
		// Ticker başlatıldığında ilk döngünün (0-->x) dolmasını bekler.
		// Bot açıldığı an profilin boş kalmaması için ilk güncellemeyi hemen yap

		updatePresence := func() {
			act := cfg.Activities[rand.Intn(len(cfg.Activities))]

			activityType, ok := library.ActivityTypeMap[act.Type]
			if !ok {
				activityType = discordgo.ActivityTypeGame
			}

			// Yeni Details ve State alanları Activity yapısına aktarılıyor.
			activity := &discordgo.Activity{
				Name:    act.Name,
				Type:    activityType,
				URL:     act.URL,
				Details: act.Details,
				State:   act.State,
			}

			err := s.UpdateStatusComplex(discordgo.UpdateStatusData{
				Status:     status,
				Activities: []*discordgo.Activity{activity},
			})
			if err != nil {
				utils.LogToConsole(utils.WARN, fmt.Sprintf("Bot durum güncellemesi başarısız oldu: %v", err))
			}
		}

		// Başlangıç güncellemesini yap
		updatePresence()

		// Belirlenen saniyede bir tetiklenecek olan ana döngü
		for range ticker.C {
			updatePresence()
		}
	}()
}
