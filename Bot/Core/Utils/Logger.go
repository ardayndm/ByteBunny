package utils

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"strings"
)

// Log seviyesi tipi.
type LogLevel string

// Konsol mesajı rengi.
type logColor string

const (
	none  LogLevel = "NONE"  // Hiçbir mesaj yazdırılmaz.
	all   LogLevel = "ALL"   // Tüm mesajlar yazdırılır.
	INFO  LogLevel = "INFO"  // Sadece bilgi mesajları yazdırılır.
	ERROR LogLevel = "ERROR" // Sadece hata mesajları yazdırılır.
	SKIP  LogLevel = "SKIP"  // Sadece geçilen durum mesajları yazdırılır.
	OK    LogLevel = "OK"    // Sadece başarılı durum mesajları yazdırılır.
	WARN  LogLevel = "WARN"  // Sadece uyarı mesajları yazdırılır.
	DEBUG LogLevel = "DEBUG" // Tüm mesajları yazdırır. (Fonksiyon içi dahil)
)

// Uygulamanın geçerli log seviyesi. (default:ALL)
var AppLogLevel LogLevel = all

const (
	colorReset    logColor = "\033[0m"               // Mesaj rengini sıfırlar.
	colorRed      logColor = "\033[31m"              // Mesaj rengini kırmızı olarak ayarlar.
	colorGreen    logColor = "\033[32m"              // Mesaj rengini yeşil olarak ayarlar.
	colorYellow   logColor = "\033[33m"              // Mesaj rengini sarı olarak ayarlar.
	colorCyan     logColor = "\033[36m"              // Mesaj rengini camgöbeği olarak ayarlar.
	colorWhite    logColor = "\033[37m"              // Mesaj rengini beyaz olarak ayarlar.
	colorWarnGray logColor = "\x1b[38;2;117;115;69m" // Mesaj rengini sarı-gri olarak ayarlar.
	colorDebug    logColor = "\033[38;2;252;3;211m"  // Mesaj rengini mor olarak ayarlar.
)

// Log seviyesini uygulamanın log seviyesine ayarlar.
func SetAppLogLevel(level LogLevel) {

	AppLogLevel = level
}

// Verilen metne göre eşleşen Log seviyesini döndürür. (Eşleşme bulunamazsa default:ERROR)
func GetLogLevelFromString(levelText string) LogLevel {
	// Verilen log seviyesini metnine göre döndür
	switch strings.ToUpper(levelText) {
	case string(none):
		return none
	case string(all):
		return all
	case string(INFO):
		return INFO
	case string(ERROR):
		return ERROR
	case string(SKIP):
		return SKIP
	case string(OK):
		return OK
	case string(WARN):
		return WARN
	case string(DEBUG):
		return DEBUG
	default:
		return ERROR
	}
}

// Log seviyesine göre konsol rengini döndürür.
func getColor(level LogLevel) logColor {

	// Verilen log seviyesine göre rengi döndür
	switch level {
	case INFO:
		return colorCyan
	case ERROR:
		return colorRed
	case WARN:
		return colorYellow
	case OK:
		return colorGreen
	case SKIP:
		return colorWarnGray
	case DEBUG:
		return colorDebug
	default:
		return colorWhite
	}
}

// Konsola mesaj yazdırmak için log seviyesini kontrol eder.
func isLogAvailable(level LogLevel) bool {

	// Uygulamın konsol seviyesini herşeyi kapsıyor ise true dön
	if AppLogLevel == all || AppLogLevel == DEBUG {
		return true
	}

	// Uygulamanın konsol seviyesi verilen seviye ile eşleşiyor ise true dön
	if AppLogLevel == level {
		return true
	}

	// Uygulamının konsol seviyesi "Bilgi" durumunda ise ama verilen filtre "Geç" ise yine true dön.
	if AppLogLevel == INFO && level == SKIP {
		return true
	}

	// Hiçbir durum eşleşmedi
	return false
}

// Verilen log seviyesi ve mesajı alıp sistem log durumuna göre konsola yazdırır
func LogToConsole(level LogLevel, msg string) bool {

	// Seçili uygulama filtresi verilen filtreye izin vermiyorsa
	if !isLogAvailable(level) {
		return false
	}

	// Caller bilgisini al (1 = Logger fonksiyonunu projenin içinde asıl tetikleyen dosya ve satır).
	_, file, line, ok := runtime.Caller(1)
	callerInfo := ""
	logMessage := ""
	color := getColor(level)

	// Eğer bir hata (ERROR) logu basılıyorsa, müdahaleyi kolaylaştırmak için dosya adı ve satır bilgisini ekle
	if ok && level == ERROR {
		// İşletim sistemine (Windows/Linux) bağımlı kalmadan dosya adını ayıkla
		fileName := filepath.Base(file)
		callerInfo = fmt.Sprintf("[ %s:%d ]", fileName, line)

		logMessage = fmt.Sprintf("%s%s - %s >> %s %s",
			color,
			fmt.Sprintf("[%s]", level),
			callerInfo,
			colorReset,
			msg)

		log.Println(logMessage)
		return true
	}

	// Standart loglar için temiz, renkli bir çıktı hazırla
	logMessage = fmt.Sprintf("%s%s%s %s",
		color,
		fmt.Sprintf("[%s]", level),
		colorReset,
		msg)

	log.Println(logMessage)

	return true
}
