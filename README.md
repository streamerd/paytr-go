# PayTR Ödeme Sistemi için payment Paketi

`payment` paketi, PayTR ödeme sistemi ile etkileşim sağlamak için hizmetler sunar. Bu paket, kart ödemeleri, tekrarlayan ödemeler, iade işlemleri ve kart yönetimi gibi işlemleri içerir.

Bu kılavuz, PayTR API'si ile nasıl entegre olunacağını ve bu paketin nasıl kullanılacağını açıklamaktadır.

## Özellikler

- Yeni kart ile ödeme işlemleri
- Kayıtlı kart ile ödeme işlemleri
- Tekrarlayan ödeme işlemleri
- İade işlemleri
- Yeni kart ekleme
- Kayıtlı kartları görüntüleme ve silme
- BIN (Bank Identification Number) bilgilerini alma

## Kullanım

### 1. Konfigürasyon Yapılandırması

Öncelikle, PayTR API'si ile etkileşim kurmak için gerekli yapılandırma bilgilerini sağlayan bir `PayTRConfig` yapısını oluşturmalısınız:

```go
type PayTRConfig struct {
    MerchantID   string // Satıcının benzersiz kimliği.
    MerchantKey  string // HMAC oluşturma için kullanılan gizli anahtar.
    MerchantSalt string // Token üretiminde güvenliği artırmak için kullanılan salt.
}
```

### 2. Yeni Servis Oluşturma

`payment` paketini kullanarak bir PayTR servisi oluşturabilirsiniz:

```go
svc := payment.NewService(PayTRConfig{
    MerchantID:   "your-merchant-id",
    MerchantKey:  "your-merchant-key",
    MerchantSalt: "your-merchant-salt",
})
```

Bu servis, PayTR API'si ile kart işlemleri, iade ve kart yönetimi işlemlerini gerçekleştirmek için kullanılır.

### 3. Kart ile Ödeme İşlemi

Bir kart ile yeni bir ödeme yapmak için `NewCardPayment` metodunu kullanabilirsiniz:

```go
req := domain.NewCardPaymentRequest{
    // Ödeme isteği verilerini girin
}

resp, err := svc.NewCardPayment(req)
if err != nil {
    // Hata işlemleri
}
```

### 4. Kayıtlı Kart ile Ödeme İşlemi

Daha önce kaydedilmiş bir kart ile ödeme yapmak için `SavedCardPayment` metodunu kullanabilirsiniz:

```go
req := domain.SavedCardPaymentRequest{
    // Kayıtlı kart ile ödeme bilgilerini girin
}

resp, err := svc.SavedCardPayment(req)
if err != nil {
    // Hata işlemleri
}
```

### 5. Tekrarlayan Ödeme

Tekrarlayan ödemeleri işlemek için `RecurringPayment` metodunu kullanabilirsiniz:

```go
req := domain.SavedCardPaymentRequest{
    // Tekrarlayan ödeme bilgilerini girin
}

resp, err := svc.RecurringPayment(req)
if err != nil {
    // Hata işlemleri
}
```

### 6. İade İşlemi

Bir ödeme işlemi için iade yapmak amacıyla `RefundPayment` metodunu kullanabilirsiniz:

```go
req := domain.RefundRequest{
    MerchantOid:  "işlem-id",
    ReturnAmount: 100.0, // İade miktarı
}

resp, err := svc.RefundPayment(req)
if err != nil {
    // Hata işlemleri
}
```

### 7. Kart Yönetimi

- Yeni kart eklemek: `AddNewCard` metodunu kullanarak bir kullanıcı hesabına yeni bir kart ekleyebilirsiniz.
- Kayıtlı kartları listelemek: `GetSavedCards` metodu, bir kullanıcının kayıtlı kartlarını görüntülemenize olanak tanır.
- Kayıtlı kartı silmek: `DeleteSavedCard` metodu ile bir kartı silebilirsiniz.

```go
// Yeni kart eklemek
addCardReq := domain.AddNewCardRequest{
    // Kart bilgilerini ekleyin
}
resp, err := svc.AddNewCard(addCardReq)

// Kayıtlı kartları listelemek
savedCardsResp, err := svc.GetSavedCards("user-token")

// Kayıtlı kartı silmek
resp, err := svc.DeleteSavedCard("user-token", "card-token")
```

### 8. BIN (Kart Numarası) Detayları

Bir kartın ilk 6-8 hanesini kullanarak BIN detaylarını almak için `GetBinDetails` metodunu kullanabilirsiniz:

```go
binDetails, err := svc.GetBinDetails("123456")
```

## HMAC İmza Üretimi

PayTR API'sine yapılacak isteklerde güvenlik için HMAC kullanılır. İmza, istek verilerinin birleştirilmesi ve SHA-256 ile HMAC oluşturulması yoluyla üretilir. Örneğin:

```go
hashStr := fmt.Sprintf("%s%s%.2f", config.MerchantID, req.MerchantOid, req.ReturnAmount)
hmac := hmac.New(sha256.New, []byte(config.MerchantKey))
hmac.Write([]byte(hashStr))
signature := base64.StdEncoding.EncodeToString(hmac.Sum(nil))
```

## Sonuç

Bu paket, PayTR API'si ile kart ödemeleri, iade işlemleri ve kart yönetimi gibi işlemleri kolaylaştırır. Örnek kodlarla açıklanan adımları takip ederek, PayTR ile güvenli ödeme işlemlerini uygulamalarınıza entegre edebilirsiniz.
