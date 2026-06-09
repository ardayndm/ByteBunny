package library

import "errors"

var (
	Err_NotFound      error = errors.New("İçerik bulunamadı.")
	Err_FileReadError error = errors.New("Dosya okunamadı.")
)
