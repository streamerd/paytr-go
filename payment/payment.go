// Package payment provides services for interacting with the PayTR payment system,
// including handling card payments, recurring payments, refunds, and card management.
//
// This package offers a `Service` interface with methods to process various payment and card management operations
// such as new card payments, saved card payments, recurring payments, refund transactions, adding new cards, deleting, and retrieving saved cards.
//
// Example usage:
//
// 1. Define a struct to hold the configuration for the PayTR service.
//
// PayTRConfig holds the configuration necessary to interact with PayTR's API, including the merchant's credentials.
//
//		type PayTRConfig struct {
//			MerchantID   string // The merchant's unique identifier.
//			MerchantKey  string // The secret key used for HMAC generation.
//			MerchantSalt string // The salt used to enhance security for token generation.
//		}
//
//	 2. Create a new instance of the PayTR service with the configuration
//	    svc := payment.NewService(PayTRConfig{
//	    MerchantID:   "your-merchant-id",
//	    MerchantKey:  "your-merchant-key",
//	    MerchantSalt: "your-merchant-salt",
//	    })
package payment

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/streamerd/paytr-go/config"
	"github.com/streamerd/paytr-go/domain"
)

// HTTPClient interface
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Service defines the operations available for interacting with the PayTR API,
// including payment processing and card management.
type Service interface {

	// NewCardPayment processes a new card payment using the provided request data.
	// Parameters:
	//   - req: A NewCardPaymentRequest struct containing details of the card payment to be processed.
	// Returns:
	//   - A PayTRResponse containing the details of the transaction.
	//   - An error if the payment processing fails.
	NewCardPayment(req domain.NewCardPaymentRequest) (*domain.PayTRResponse, error)

	// SavedCardPayment processes a payment using a previously saved card.
	// Parameters:
	//   - req: A SavedCardPaymentRequest struct containing details of the saved card payment.
	// Returns:
	//   - A PayTRResponse containing the details of the transaction.
	//   - An error if the payment processing fails.
	SavedCardPayment(req domain.SavedCardPaymentRequest) (*domain.PayTRResponse, error)

	// RecurringPayment processes a recurring payment using a saved card.
	// Parameters:
	//   - req: A SavedCardPaymentRequest struct containing details of the recurring payment.
	// Returns:
	//   - A PayTRResponse containing the details of the transaction.
	//   - An error if the payment processing fails.
	RecurringPayment(req domain.SavedCardPaymentRequest) (*domain.PayTRResponse, error)

	// RefundPayment refunds a payment by the specified amount.
	// Parameters:
	//   - req: A RefundRequest struct containing details of the refund, including the amount to refund.
	// Returns:
	//   - A PayTRResponse containing the details of the refund transaction.
	//   - An error if the refund process fails.
	RefundPayment(req domain.RefundRequest) (*domain.PayTRResponse, error)

	// GetTransactionDetails retrieves details for a transaction within the given date range.
	// Parameters:
	//   - req: A TransactionDetailsRequest struct specifying the date range and transaction details to query.
	// Returns:
	//   - A TransactionDetailsResponse containing the transaction details.
	//   - An error if the request for transaction details fails.
	GetTransactionDetails(req domain.TransactionDetailsRequest) (*domain.TransactionDetailsResponse, error)

	// MerchantStatusInquiry inquires about the status of a merchant transaction.
	// Parameters:
	//   - req: A StatusInquiryRequest struct specifying the details of the merchant transaction to inquire about.
	// Returns:
	//   - A StatusInquiryResponse containing the status of the transaction.
	//   - An error if the status inquiry process fails.
	MerchantStatusInquiry(req domain.StatusInquiryRequest) (*domain.StatusInquiryResponse, error)

	// AddNewCard saves a new card to the user's account.
	// Parameters:
	//   - req: An AddNewCardRequest struct containing the card details to be saved.
	// Returns:
	//   - A PayTRResponse confirming the success or failure of the card saving process.
	//   - An error if the card saving process fails.
	AddNewCard(req domain.AddNewCardRequest) (*domain.PayTRResponse, error)

	// GetSavedCards retrieves the list of saved cards for a given user token.
	// Parameters:
	//   - utoken: A string representing the user's token, used to identify the user and fetch saved cards.
	// Returns:
	//   - A PayTRResponse containing the list of saved cards.
	//   - An error if the retrieval process fails.
	GetSavedCards(utoken string) (*domain.PayTRResponse, error)

	// GetBinDetails retrieves details about a BIN (Bank Identification Number).
	// Parameters:
	//   - binNumber: A string representing the BIN (first 6-8 digits of a card) to retrieve details for.
	// Returns:
	//   - A PayTRResponse containing BIN details such as the bank and card type.
	//   - An error if the BIN lookup process fails.
	GetBinDetails(binNumber string) (*domain.PayTRResponse, error)

	// DeleteSavedCard removes a saved card using the provided user and card tokens.
	// Parameters:
	//   - utoken: A string representing the user's token, used to identify the user.
	//   - ctoken: A string representing the card's token, used to identify the specific card to delete.
	// Returns:
	//   - A PayTRResponse confirming the success or failure of the card deletion process.
	//   - An error if the card deletion process fails.
	DeleteSavedCard(utoken, ctoken string) (*domain.PayTRResponse, error)
	SetHTTPClient(client HTTPClient)
}

type service struct {
	config config.PayTRConfig
	client HTTPClient
}

func (s *service) SetHTTPClient(client HTTPClient) {
	s.client = client
}

// NewService creates a new PayTR service with the provided configuration and repository.
func NewService(config config.PayTRConfig) Service {
	return &service{
		config: config,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

// PAYMENTS

// NewCardPayment processes a payment using the details from the NewCardPaymentRequest.
// The payment details are validated, and the PayTR token is generated based on the request data.
func (s *service) NewCardPayment(req domain.NewCardPaymentRequest) (*domain.PayTRResponse, error) {
	req.PayTRToken = s.generateToken(req.CommonPaymentRequest)
	return s.sendRequest(req, domain.PayTRBaseURL+"/odeme")
}

func (s *service) SavedCardPayment(req domain.SavedCardPaymentRequest) (*domain.PayTRResponse, error) {
	req.PayTRToken = s.generateToken(req.CommonPaymentRequest)
	return s.sendRequest(req, domain.PayTRBaseURL+"/odeme")
}

func (s *service) RecurringPayment(req domain.SavedCardPaymentRequest) (*domain.PayTRResponse, error) {
	req.RecurringPayment = "1"
	req.PayTRToken = s.generateToken(req.CommonPaymentRequest)
	return s.sendRequest(req, domain.PayTRBaseURL+"/odeme")
}

func (s *service) RefundPayment(req domain.RefundRequest) (*domain.PayTRResponse, error) {
	paytrReq := struct {
		MerchantID   string  `json:"merchant_id"`
		MerchantOid  string  `json:"merchant_oid"`
		ReturnAmount float64 `json:"return_amount"`
		PayTRToken   string  `json:"paytr_token"`
		ReferenceNo  string  `json:"reference_no,omitempty"`
	}{
		MerchantID:   s.config.MerchantID,
		MerchantOid:  req.MerchantOid,
		ReturnAmount: req.ReturnAmount,
		ReferenceNo:  req.ReferenceNo,
	}

	// Generate PayTR token
	hashStr := fmt.Sprintf("%s%s%.2f", s.config.MerchantID, req.MerchantOid, req.ReturnAmount)
	paytrReq.PayTRToken = s.generateSimpleToken(hashStr)

	return s.sendRequest(paytrReq, domain.PayTRBaseURL+"/odeme/iade")
}

func (s *service) MerchantStatusInquiry(req domain.StatusInquiryRequest) (*domain.StatusInquiryResponse, error) {
	paytrReq := struct {
		MerchantID  string `json:"merchant_id"`
		MerchantOid string `json:"merchant_oid"`
		PayTRToken  string `json:"paytr_token"`
	}{
		MerchantID:  s.config.MerchantID,
		MerchantOid: req.MerchantOid,
	}

	paytrReq.PayTRToken = s.generateSimpleToken(s.config.MerchantID + req.MerchantOid)

	paytrResp, err := s.sendRequest(paytrReq, domain.PayTRBaseURL+"/odeme/durum-sorgu")
	if err != nil {
		return nil, err
	}

	if paytrResp.Status != "success" {
		return nil, fmt.Errorf("PayTR error: %s", paytrResp.Message)
	}

	var result domain.StatusInquiryResponse
	err = mapstructure.Decode(paytrResp.Data, &result)
	if err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &result, nil
}

func (s *service) GetTransactionDetails(req domain.TransactionDetailsRequest) (*domain.TransactionDetailsResponse, error) {
	paytrReq := struct {
		MerchantID string `json:"merchant_id"`
		StartDate  string `json:"start_date"`
		EndDate    string `json:"end_date"`
		Dummy      int    `json:"dummy,omitempty"`
		PayTRToken string `json:"paytr_token"`
	}{
		MerchantID: s.config.MerchantID,
		StartDate:  req.StartDate,
		EndDate:    req.EndDate,
		Dummy:      req.Dummy,
	}

	paytrReq.PayTRToken = s.generateSimpleToken(s.config.MerchantID + req.StartDate + req.EndDate)

	paytrResp, err := s.sendRequest(paytrReq, domain.PayTRBaseURL+"/rapor/islem-dokumu")
	if err != nil {
		return nil, err
	}

	var result domain.TransactionDetailsResponse
	err = mapstructure.Decode(paytrResp.Data, &result)
	if err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &result, nil
}

// CARDS

func (s *service) GetBinDetails(binNumber string) (*domain.PayTRResponse, error) {
	req := struct {
		MerchantID string `json:"merchant_id"`
		BinNumber  string `json:"bin_number"`
		PayTRToken string `json:"paytr_token"`
	}{
		MerchantID: s.config.MerchantID,
		BinNumber:  binNumber,
		PayTRToken: s.generateSimpleToken(binNumber + s.config.MerchantID),
	}
	return s.sendRequest(req, domain.PayTRBaseURL+"/odeme/api/bin-detail")
}

func (s *service) GetSavedCards(utoken string) (*domain.PayTRResponse, error) {
	req := struct {
		MerchantID string `json:"merchant_id"`
		UToken     string `json:"utoken"`
		PayTRToken string `json:"paytr_token"`
	}{
		MerchantID: s.config.MerchantID,
		UToken:     utoken,
		PayTRToken: s.generateSimpleToken(utoken),
	}
	return s.sendRequest(req, domain.PayTRBaseURL+"/odeme/capi/list")
}

func (s *service) DeleteSavedCard(utoken, ctoken string) (*domain.PayTRResponse, error) {
	req := struct {
		MerchantID string `json:"merchant_id"`
		UToken     string `json:"utoken"`
		CToken     string `json:"ctoken"`
		PayTRToken string `json:"paytr_token"`
	}{
		MerchantID: s.config.MerchantID,
		UToken:     utoken,
		CToken:     ctoken,
		PayTRToken: s.generateSimpleToken(utoken + ctoken),
	}
	return s.sendRequest(req, domain.PayTRBaseURL+"/odeme/capi/delete")
}

func (s *service) AddNewCard(req domain.AddNewCardRequest) (*domain.PayTRResponse, error) {
	// Prepare the request for adding a new card
	paytrReq := domain.NewCardPaymentRequest{
		CommonPaymentRequest: domain.CommonPaymentRequest{
			MerchantID:       s.config.MerchantID,
			UserIP:           req.UserIP,
			MerchantOid:      req.MerchantOid,
			Email:            req.Email,
			PaymentAmount:    1, // Minimal amount for card validation
			PaymentType:      "card",
			Currency:         "TRY",
			TestMode:         "1",
			NonThreeD:        "0",
			MerchantOkURL:    req.MerchantOkURL,
			MerchantFailURL:  req.MerchantFailURL,
			UserName:         req.CardOwner,
			UserAddress:      req.UserAddress,
			UserPhone:        req.UserPhone,
			UserBasket:       `[["Card Validation", "1", 1]]`,
			DebugOn:          "1",
			ClientLang:       "tr",
			InstallmentCount: "0",
		},
		CardOwner:   req.CardOwner,
		CardNumber:  req.CardNumber,
		ExpiryMonth: req.ExpiryMonth,
		ExpiryYear:  req.ExpiryYear,
		CVV:         req.CVV,
		CardType:    req.CardType,
		StoreCard:   "1",
	}

	paytrReq.PayTRToken = s.generateToken(paytrReq.CommonPaymentRequest)
	return s.sendRequest(paytrReq, domain.PayTRBaseURL+"/odeme")
}

// generateToken generates an HMAC token based on the payment request and the merchant's secret key.
// Parameters:
//   - req: A CommonPaymentRequest struct containing the necessary payment details, including user IP,
//     merchant order ID, email, payment amount, payment type, installment count, currency,
//     test mode, and whether it's a non-3D payment.
//
// Returns:
//   - A base64-encoded string representing the generated HMAC token.// generateToken generates an HMAC token based on the payment request and the merchant's secret key.
//
// Parameters:
//   - req: A CommonPaymentRequest struct containing the necessary payment details, including user IP,
//     merchant order ID, email, payment amount, payment type, installment count, currency,
//     test mode, and whether it's a non-3D payment.
//
// Returns:
//   - A base64-encoded string representing the generated HMAC token.
func (s *service) generateToken(req domain.CommonPaymentRequest) string {
	hashStr := fmt.Sprintf("%s%s%s%s%s%s%s%s%s%s",
		s.config.MerchantID,
		req.UserIP,
		req.MerchantOid,
		req.Email,
		strconv.FormatFloat(req.PaymentAmount, 'f', 2, 64),
		req.PaymentType,
		req.InstallmentCount,
		req.Currency,
		req.TestMode,
		req.NonThreeD,
	)
	hmacStr := hashStr + s.config.MerchantSalt
	h := hmac.New(sha256.New, []byte(s.config.MerchantKey))
	h.Write([]byte(hmacStr))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// generateSimpleToken generates a simple HMAC-based token by concatenating the input data
// with the merchant salt and using the merchant key to hash the result. The token is then
// encoded in base64 format for secure transmission.
// Parameters:
//   - data: A string input that is concatenated with the merchant salt to form the HMAC message.
//
// Returns:
//   - A base64-encoded string that represents the generated HMAC token.
func (s *service) generateSimpleToken(data string) string {
	hmacStr := data + s.config.MerchantSalt
	h := hmac.New(sha256.New, []byte(s.config.MerchantKey))
	h.Write([]byte(hmacStr))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// sendRequest sends an HTTP POST request to the provided URL with the given request payload.
// The request is marshaled into JSON format and sent with the appropriate content type.
// It then reads and decodes the response into a PayTRResponse object.
// Parameters:
//   - req: The request payload that is marshaled into JSON and sent to the URL.
//   - url: The endpoint to which the request is sent.
//
// Returns:
//   - A pointer to PayTRResponse containing the response data from the PayTR API.
//   - An error if any issue occurs during the request or response processing.
func (s *service) sendRequest(req interface{}, url string) (*domain.PayTRResponse, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result domain.PayTRResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
