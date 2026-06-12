package library

import "errors"

var (
	Err_NotFound           error = errors.New("İçerik bulunamadı.")
	Err_FileReadError      error = errors.New("Dosya okunamadı.")
	Err_InvalidName        error = errors.New("Geçersiz ad.")
	Err_InvalidPath        error = errors.New("Geçersiz dosya yolu.")
	Err_InvalidCommandName error = errors.New("Geçersiz komut adı.")
	Err_InvalidLanguage    error = errors.New("Geçersiz dil kodu.")
	Err_SuspectedFilePath  error = errors.New("Şüpheli dosya yolu.")
	Err_FileNotYaml        error = errors.New("Dosya .yaml formatında değil.")
)
