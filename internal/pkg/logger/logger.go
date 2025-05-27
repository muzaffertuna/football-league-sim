package logger

import (
	"io" // io paketi eklendi
	"log"
	"os"
)

type Logger struct {
	*log.Logger
}

func NewLogger() *Logger {
	// app.log dosyasına yazma
	file, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		// Log dosyasını açamazsak kritik hata ver ve çık
		log.Fatalf("Failed to open log file: %v", err)
	}

	// Hem dosyaya (file) hem de standart çıktıya (os.Stdout) yazmak için bir MultiWriter kullanırız.
	multiWriter := io.MultiWriter(file, os.Stdout)

	// log.Logger instance'ını oluştur
	// Bayraklar: Tarih, saat ve kısa dosya yolu bilgisi eklensin
	// Prefix boş bırakıldı, çünkü Info/Error metotlarında kendimiz formatlayacağız
	stdLogger := log.New(multiWriter, "", log.Ldate|log.Ltime|log.Lshortfile)

	return &Logger{stdLogger}
}

func (l *Logger) Info(msg string) {
	l.Printf("[INFO] %s", msg)
}

func (l *Logger) Error(msg string) {
	l.Printf("[ERROR] %s", msg)
}
