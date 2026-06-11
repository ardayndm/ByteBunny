package utils

import (
	"context"
	"time"
)

// Uzun sürecek işlemler için Timeout hazırlar
func GetContext(waitSecond int) (context.Context, context.CancelFunc) {
	// Belirtilen saniye kadar (Örn: 5 * time.Second) süre sınırına sahip arka plan context'i oluşturuluyor.
	// İşlem tamamlandığında veya süre bittiğinde bellekten temizlenmesi için bir CancelFunc (iptal fonksiyonu) döndürülür.
	return context.WithTimeout(context.Background(), time.Duration(waitSecond)*time.Second)
}
