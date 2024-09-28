package domain

import (
	"time"
)

const (
	PayTRBaseURL = "https://www.paytr.com"
)

type PayTRConfig struct {
	MerchantID   string
	MerchantKey  string
	MerchantSalt string
}

type CommonPaymentRequest struct {
	MerchantID       string  `json:"merchant_id"`
	UserIP           string  `json:"user_ip"`
	MerchantOid      string  `json:"merchant_oid"`
	Email            string  `json:"email"`
	PaymentAmount    float64 `json:"payment_amount"`
	PaymentType      string  `json:"payment_type"`
	Currency         string  `json:"currency"`
	TestMode         string  `json:"test_mode"`
	NonThreeD        string  `json:"non_3d"`
	MerchantOkURL    string  `json:"merchant_ok_url"`
	MerchantFailURL  string  `json:"merchant_fail_url"`
	UserName         string  `json:"user_name"`
	UserAddress      string  `json:"user_address"`
	UserPhone        string  `json:"user_phone"`
	UserBasket       string  `json:"user_basket"`
	DebugOn          string  `json:"debug_on"`
	ClientLang       string  `json:"client_lang"`
	PayTRToken       string  `json:"paytr_token"`
	InstallmentCount string  `json:"installment_count"`
}

type NewCardPaymentRequest struct {
	CommonPaymentRequest
	CardOwner   string `json:"cc_owner"`
	CardNumber  string `json:"card_number"`
	ExpiryMonth string `json:"expiry_month"`
	ExpiryYear  string `json:"expiry_year"`
	CVV         string `json:"cvv"`
	CardType    string `json:"card_type"`
	StoreCard   string `json:"store_card"`
}

// type PayTRResponse struct {
// 	Status  string      `json:"status"`
// 	Message string      `json:"message"`
// 	Data    interface{} `json:"data,omitempty"`
// }

type PayTRResponse struct {
	Status  string                 `json:"status"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data,omitempty"`
}

type Payment struct {
	ID            string    `bson:"_id,omitempty"`
	UserID        string    `bson:"user_id"`
	Amount        float64   `bson:"amount"`
	Currency      string    `bson:"currency"`
	Status        string    `bson:"status"`
	PaymentMethod string    `bson:"payment_method"`
	MerchantOid   string    `bson:"merchant_oid"`
	CreatedAt     time.Time `bson:"created_at"`
	UpdatedAt     time.Time `bson:"updated_at"`
}

type SavedCard struct {
	ID         string    `bson:"_id,omitempty"`
	UserID     string    `bson:"user_id"`
	UToken     string    `bson:"utoken"`
	CToken     string    `bson:"ctoken"`
	LastFour   string    `bson:"last_four"`
	CardType   string    `bson:"card_type"`
	ExpiryDate string    `bson:"expiry_date"`
	CreatedAt  time.Time `bson:"created_at"`
}

type SavedCardPaymentRequest struct {
	CommonPaymentRequest
	UToken           string `json:"utoken"`
	CToken           string `json:"ctoken"`
	CVV              string `json:"cvv"`
	RecurringPayment string `json:"recurring_payment"`
}

type AddNewCardRequest struct {
	UserID          string  `json:"user_id"`
	CardOwner       string  `json:"cc_owner"`
	CardNumber      string  `json:"card_number"`
	ExpiryMonth     string  `json:"expiry_month"`
	ExpiryYear      string  `json:"expiry_year"`
	CVV             string  `json:"cvv"`
	CardType        string  `json:"card_type"`
	UserIP          string  `json:"user_ip"`
	MerchantOid     string  `json:"merchant_oid"`
	Email           string  `json:"email"`
	UserName        string  `json:"user_name"`
	UserAddress     string  `json:"user_address"`
	UserPhone       string  `json:"user_phone"`
	MerchantOkURL   string  `json:"merchant_ok_url"`
	MerchantFailURL string  `json:"merchant_fail_url"`
	ClientLang      string  `json:"client_lang"`
	Currency        string  `json:"currency"`
	PaymentAmount   float64 `json:"payment_amount"`
}

type RefundRequest struct {
	MerchantOid  string  `json:"merchant_oid"`
	ReturnAmount float64 `json:"return_amount"`
	ReferenceNo  string  `json:"reference_no,omitempty"`
}

type StatusInquiryRequest struct {
	MerchantOid string `json:"merchant_oid"`
}

type StatusInquiryResponse struct {
	Status              string               `json:"status"`
	PaymentAmount       string               `json:"payment_amount,omitempty"`
	PaymentTotal        string               `json:"payment_total,omitempty"`
	PaymentDate         string               `json:"payment_date,omitempty"`
	Currency            string               `json:"currency,omitempty"`
	NetTutar            string               `json:"net_tutar,omitempty"`
	KesintiTutari       string               `json:"kesinti_tutari,omitempty"`
	Taksit              string               `json:"taksit,omitempty"`
	KartMarka           string               `json:"kart_marka,omitempty"`
	MaskedPan           string               `json:"masked_pan,omitempty"`
	OdemeTipi           string               `json:"odeme_tipi,omitempty"`
	TestMode            string               `json:"test_mode,omitempty"`
	Returns             string               `json:"returns,omitempty"`
	ErrNo               string               `json:"err_no,omitempty"`
	ErrMsg              string               `json:"err_msg,omitempty"`
	SubmerchantPayments []SubmerchantPayment `json:"submerchant_payments,omitempty"`
}

type SubmerchantPayment struct {
	SubmerchantId           string `json:"submerchant_id"`
	SubmerchantPrice        string `json:"submerchant_price"`
	SubmerchantPayoutRate   string `json:"submerchant_payout_rate"`
	SubmerchantPayoutAmount string `json:"submerchant_payout_amount"`
}

type TransactionDetailsRequest struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Dummy     int    `json:"dummy,omitempty"`
}

type TransactionDetailsResponse struct {
	Status       string        `json:"status"`
	Transactions []Transaction `json:"transactions,omitempty"`
	ErrMsg       string        `json:"err_msg,omitempty"`
}

type Transaction struct {
	IslemTipi     string `json:"islem_tipi"`
	NetTutar      string `json:"net_tutar"`
	KesintiTutari string `json:"kesinti_tutari"`
	KesintiOrani  string `json:"kesinti_orani"`
	IslemTutari   string `json:"islem_tutari"`
	OdemeTutari   string `json:"odeme_tutari"`
	IslemTarihi   string `json:"islem_tarihi"`
	ParaBirimi    string `json:"para_birimi"`
	Taksit        string `json:"taksit"`
	KartMarka     string `json:"kart_marka"`
	KartNo        string `json:"kart_no"`
	SiparisNo     string `json:"siparis_no"`
	OdemeTipi     string `json:"odeme_tipi"`
}
