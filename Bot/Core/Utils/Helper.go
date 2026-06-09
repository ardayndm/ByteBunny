package utils

// İlk verinin boş olma durumunu kontrol eder , eğer boş ise fallback döner.
func GetOrDefault[T comparable](value, fallback T) T {

	var zero T

	if value == zero {
		return fallback
	}
	return value
}
