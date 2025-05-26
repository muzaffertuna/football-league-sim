package logger // <--- Paket adını logger olarak değiştirdik

import (
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

	// log.Logger instance'ını oluştur
	// Bayraklar: Tarih, saat ve kısa dosya yolu bilgisi eklensin
	// Prefix boş bırakıldı, çünkü Info/Error metotlarında kendimiz formatlayacağız
	stdLogger := log.New(file, "", log.Ldate|log.Ltime|log.Lshortfile)

	// Hem dosyaya hem de standart çıktıya yazmak için bir MultiWriter kullanabiliriz.
	// Ancak sizin mevcut kodunuz sadece stdout'a SetOutput yapmış, bu da log dosyasını etkisiz kılar.
	// Hem dosyaya hem de stdout'a yazmasını istiyorsanız:
	// multiWriter := io.MultiWriter(file, os.Stdout)
	// stdLogger.SetOutput(multiWriter)

	// Sizin mevcut kodunuzdaki gibi sadece stdout'a yazdırır ve dosyaya yazmaz:
	stdLogger.SetOutput(os.Stdout) // <--- Burası sadece konsola yazar, dosyaya yazmaz.
	// Eğer hem dosyaya hem konsola istiyorsanız yukarıdaki multiWriter'ı kullanın.

	return &Logger{stdLogger}
}

func (l *Logger) Info(msg string) {
	l.Printf("[INFO] %s", msg)
}

func (l *Logger) Error(msg string) {
	l.Printf("[ERROR] %s", msg)
}
