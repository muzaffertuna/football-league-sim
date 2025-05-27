# Football League Simulation API

Bu proje, bir futbol ligini simüle etmek ve takım şampiyonluk olasılıklarını tahmin etmek için tasarlanmış, Go ile geliştirilmiş bir REST API'sidir.

## Özellikler

* **Lig Tablosu:** Mevcut lig puan durumunu, gol farkını ve şampiyonluk tahminlerini içerir.
* **Haftalık Maç Simülasyonu:** Ligin mevcut haftasındaki maçları simüle eder ve sonuçları günceller.
* **Tüm Lig Simülasyonu:** Kalan tüm haftaları otomatik olarak simüle ederek sezonu tamamlar.
* **Lig Sıfırlama:** Takım istatistiklerini ve maç fikstürünü sıfırlayarak yeni bir sezon başlatır.
* **Şampiyonluk Tahminleri:** Monte Carlo simülasyonu algoritması kullanılarak, ligin kalanına göre takımların şampiyonluk olasılıkları hesaplanır. Bu hesaplama, Go'nun eşzamanlılık modeli (goroutine'ler) ile çoklu iş parçacıklı (multithreaded) olarak paralel çalıştırılır.
* **Gelişmiş Loglama:** Uygulama içerisinde detaylı bilgi, uyarı ve hata logları tutulur, bu da geliştirme ve hata ayıklama süreçlerini kolaylaştırır.
* **Otomatik Veritabanı Oluşturma ve Migration:** Uygulama başlatıldığında, gerekli veritabanı şeması ve tabloları otomatik olarak kontrol edilir ve yoksa oluşturulur. Ayrıca, ilk fikstür verileri otomatik olarak veritabanına eklenir.

## Ön Gereksinimler

Bu projeyi yerel ortamınızda kurup çalıştırmak için aşağıdaki yazılımlara ihtiyacınız vardır:

* **Go:** Sürüm 1.20 veya üzeri [Go Yükle](https://go.dev/doc/install)
* **Docker Desktop:** MSSQL veritabanını çalıştırmak için [Docker Desktop İndir](https://www.docker.com/products/docker-desktop/)
* **Git:** Projeyi klonlamak için [Git İndir](https://git-scm.com/downloads)

## Kurulum ve Çalıştırma

Aşağıdaki adımları takip ederek projeyi yerel ortamınızda kolayca kurup çalıştırabilirsiniz:

### 1. Depoyu Klonlama

Öncelikle proje deposunu bilgisayarınıza klonlayın:

```bash
git clone [https://github.com/muzaffertuna/football-league-sim.git](https://github.com/muzaffertuna/football-league-sim.git)
cd football-league-sim

2. Ortam Değişkenlerini Yapılandırma (.env Dosyası)
API uygulaması, veritabanı bağlantı bilgileri ve sunucu adresi gibi yapılandırma ayarlarını ortam değişkenleri aracılığıyla alır. Projenin kök dizininde (football-league-sim) .env adında yeni bir dosya oluşturun ve aşağıdaki içeriği kendi belirleyeceğiniz değerlerle ekleyin:

SA_PASSWORD="YourStrongPassword123" # <<< BURAYI KENDİ GÜÇLÜ ŞİFRENİZLE DEĞİŞTİRİN!
DB_CONNECTION_STRING="sqlserver://sa:${SA_PASSWORD}@localhost:1433?database=FootballLeagueSim&TrustServerCertificate=true"
SERVER_ADDRESS=":8080"

3. Veritabanı Kurulumu ve Başlatma (Docker Compose)

Projenin kök dizininde (football-league-sim), MSSQL veritabanı container'ını yapılandıran docker-compose.yml dosyası zaten bulunmaktadır. Bu dosya, .env dosyasındaki SA_PASSWORD değerini kullanarak MSSQL Server'ı ayağa kaldırır.

Veritabanını başlatmak için:
docker-compose up -d

4. Go Uygulamasını Çalıştırma
Veritabanı hazır olduğunda ve .env dosyanız yapılandırıldığında, API uygulamasını başlatabilirsiniz:

# Projenin ana dizininde (football-league-sim) olduğunuzdan emin olun
cd cmd/api # Go uygulamanızın main.go dosyasının bulunduğu dizine geçin

go run main.go

Uygulama başarıyla başladığında konsolda aşağıdaki gibi bir çıktı görmelisiniz:

2025/05/27 16:07:56 logger.go:39: [INFO] Successfully connected to MSSQL
2025/05/27 16:07:56 logger.go:39: [INFO] Starting server on :8080


5. API'ye Erişim ve Test Etme

Uygulama yerel olarak http://localhost:8080 adresinde çalışacaktır. API endpoint'lerine istek göndermek için iki ana yöntem bulunmaktadır:

Yöntem 1: Swagger UI Kullanımı (Önerilen)
Swagger UI, API'nin tüm endpoint'lerini görsel bir arayüzle sunar ve doğrudan tarayıcı üzerinden istek göndermenizi sağlar. Bu, API'yi keşfetmek ve test etmek için en kolay yoldur.

Tarayıcınızdan şu adresi ziyaret edin: http://localhost:8080/swagger/index.html#/

Her bir endpoint'i genişleterek detaylarını görebilir, "Try it out" düğmesine tıklayarak istek gönderebilir ve yanıtı anında inceleyebilirsiniz.

Yöntem 2: Postman, cURL veya Benzeri Bir Araç Kullanımı
Eğer tercih ediyorsanız, Postman, Insomnia veya cURL gibi popüler API test araçlarını kullanarak da endpoint'lere istek atabilirsiniz:



API Endpoint'leri ve Kullanım Sırası:
Uygulamanın doğru şekilde çalışması için ilk kullanışta ligi sıfırlamanız gerekmektedir.


POST /reset-league (Ligi Sıfırla - Yeni Bir Sezon Başlat)
Açıklama: Uygulamayı ilk kez başlattığınızda veya yeni bir sezon başlatmak istediğinizde, tüm takım istatistiklerini ve maç fikstürünü başlangıç durumuna getirir. Bu endpoint'i diğer simülasyon işlemlerinden önce çalıştırmak ZORUNLUDUR, aksi takdirde veritabanında gerekli veriler olmadığı için hatalar alabilirsiniz.(İlk kullanımdan sonra istediğiniz gibi devam edebilirsiniz.)

cURL Örneği: curl -X POST http://localhost:8080/reset-league


POST /play-week (Mevcut Haftayı Oynat)
Açıklama: Ligin mevcut haftasındaki maçları simüle eder ve takımların puanlarını günceller. Her çağrıda bir sonraki haftayı oynatır.
cURL Örneği: curl -X POST http://localhost:8080/play-week


GET /league-table (Lig Tablosunu Getir)
Açıklama: Mevcut lig puan durumunu, gol farklarını ve Monte Carlo simülasyonuna dayalı şampiyonluk tahminlerini döndürür. Bu çağrı, tahminleri yeniden hesaplamak için çoklu iş parçacıklı simülasyonları tetikler.

cURL Örneği: curl -X GET http://localhost:8080/league-table

POST /simulate-all-weeks (Tüm Kalan Haftaları Simüle Et)
Açıklama: Ligdeki kalan tüm haftaları otomatik olarak simüle eder ve sezonu tamamlar.
cURL Örneği: curl -X POST http://localhost:8080/simulate-all-weeks