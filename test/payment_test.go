package payment_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/streamerd/paytr-go/config"
	"github.com/streamerd/paytr-go/domain"
	"github.com/streamerd/paytr-go/payment"
)

// HTTPClient interface
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// mockHTTPClient is a custom HTTP client for testing
type mockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

// setupTestService creates a test service with a mock HTTP client
func setupTestService(mockResponse *domain.PayTRResponse) payment.Service {
	mockClient := &mockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			responseBody, _ := json.Marshal(mockResponse)
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewBuffer(responseBody)),
			}, nil
		},
	}

	testService := payment.NewService(config.PayTRConfig{
		MerchantID:   "test_merchant",
		MerchantKey:  "test_key",
		MerchantSalt: "test_salt",
	})

	testService.SetHTTPClient(mockClient)
	return testService
}

func TestNewCardPayment(t *testing.T) {
	mockResponse := &domain.PayTRResponse{
		Status:  "success",
		Message: "Payment successful",
		Data: map[string]interface{}{
			"token": "mock_token",
		},
	}

	testService := setupTestService(mockResponse)

	req := domain.NewCardPaymentRequest{
		CommonPaymentRequest: domain.CommonPaymentRequest{
			MerchantID:    "test_merchant",
			UserIP:        "127.0.0.1",
			MerchantOid:   "test_order_123",
			Email:         "test@example.com",
			PaymentAmount: 100.00,
			PaymentType:   "card",
			Currency:      "TRY",
			TestMode:      "1",
		},
		CardOwner:   "John Doe",
		CardNumber:  "4111111111111111",
		ExpiryMonth: "12",
		ExpiryYear:  "2025",
		CVV:         "123",
	}

	resp, err := testService.NewCardPayment(req)

	if err != nil {
		t.Fatalf("NewCardPayment returned an error: %v", err)
	}

	if resp.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", resp.Status)
	}
}

func TestSavedCardPayment(t *testing.T) {
	mockResponse := &domain.PayTRResponse{
		Status:  "success",
		Message: "Saved card payment successful",
	}

	testService := setupTestService(mockResponse)

	req := domain.SavedCardPaymentRequest{
		CommonPaymentRequest: domain.CommonPaymentRequest{
			MerchantID:    "test_merchant",
			UserIP:        "127.0.0.1",
			MerchantOid:   "test_order_456",
			Email:         "test@example.com",
			PaymentAmount: 200.00,
			PaymentType:   "card",
			Currency:      "TRY",
			TestMode:      "1",
		},
		UToken: "test_utoken",
		CToken: "test_ctoken",
		CVV:    "123",
	}

	resp, err := testService.SavedCardPayment(req)

	if err != nil {
		t.Fatalf("SavedCardPayment returned an error: %v", err)
	}

	if resp.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", resp.Status)
	}
}

func TestRecurringPayment(t *testing.T) {
	mockResponse := &domain.PayTRResponse{
		Status:  "success",
		Message: "Recurring payment successful",
	}

	testService := setupTestService(mockResponse)

	req := domain.SavedCardPaymentRequest{
		CommonPaymentRequest: domain.CommonPaymentRequest{
			MerchantID:    "test_merchant",
			UserIP:        "127.0.0.1",
			MerchantOid:   "test_order_789",
			Email:         "test@example.com",
			PaymentAmount: 50.00,
			PaymentType:   "card",
			Currency:      "TRY",
			TestMode:      "1",
		},
		UToken:           "test_utoken",
		CToken:           "test_ctoken",
		CVV:              "123",
		RecurringPayment: "1",
	}

	resp, err := testService.RecurringPayment(req)

	if err != nil {
		t.Fatalf("RecurringPayment returned an error: %v", err)
	}

	if resp.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", resp.Status)
	}
}

func TestRefundPayment(t *testing.T) {
	mockResponse := &domain.PayTRResponse{
		Status:  "success",
		Message: "Refund successful",
	}

	testService := setupTestService(mockResponse)

	req := domain.RefundRequest{
		MerchantOid:  "test_order_789",
		ReturnAmount: 50.00,
	}

	resp, err := testService.RefundPayment(req)

	if err != nil {
		t.Fatalf("RefundPayment returned an error: %v", err)
	}

	if resp.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", resp.Status)
	}
}

// func TestGetTransactionDetails(t *testing.T) {
// 	mockResponse := &domain.PayTRResponse{
// 		Status:  "success",
// 		Message: "Transaction details retrieved",
// 		Data: map[string]interface{}{
// 			"status": "success",
// 			"transactions": []map[string]interface{}{
// 				{
// 					"islem_tipi":     "sale",
// 					"net_tutar":      "95.00",
// 					"kesinti_tutari": "5.00",
// 					"kesinti_orani":  "5",
// 					"islem_tutari":   "100.00",
// 					"odeme_tutari":   "100.00",
// 					"islem_tarihi":   "2023-01-01 12:00:00",
// 					"para_birimi":    "TRY",
// 					"taksit":         "1",
// 					"kart_marka":     "VISA",
// 					"kart_no":        "411111******1111",
// 					"siparis_no":     "test_order_123",
// 					"odeme_tipi":     "card",
// 				},
// 			},
// 		},
// 	}

// 	testService := setupTestService(mockResponse)

// 	req := domain.TransactionDetailsRequest{
// 		StartDate: "2023-01-01",
// 		EndDate:   "2023-12-31",
// 	}

// 	resp, err := testService.GetTransactionDetails(req)
// 	if err != nil {
// 		t.Fatalf("GetTransactionDetails returned an error: %v", err)
// 	}

// 	t.Logf("Response: %+v", resp)

// 	if resp.Status != "success" {
// 		t.Errorf("Expected status 'success', got '%s'", resp.Status)
// 	}

// 	if len(resp.Transactions) != 1 {
// 		t.Fatalf("Expected 1 transaction, got %d", len(resp.Transactions))
// 	}

// 	transaction := resp.Transactions[0]
// 	t.Logf("Transaction: %+v", transaction)

// 	expectedFields := map[string]string{
// 		"IslemTipi":     "sale",
// 		"NetTutar":      "95.00",
// 		"KesintiTutari": "5.00",
// 		"KesintiOrani":  "5",
// 		"IslemTutari":   "100.00",
// 		"OdemeTutari":   "100.00",
// 		"IslemTarihi":   "2023-01-01 12:00:00",
// 		"ParaBirimi":    "TRY",
// 		"Taksit":        "1",
// 		"KartMarka":     "VISA",
// 		"KartNo":        "411111******1111",
// 		"SiparisNo":     "test_order_123",
// 		"OdemeTipi":     "card",
// 	}

// 	for field, expected := range expectedFields {
// 		actual := reflect.ValueOf(transaction).FieldByName(field).String()
// 		if actual != expected {
// 			t.Errorf("Expected %s '%s', got '%s'", field, expected, actual)
// 		}
// 	}
// }

func TestMerchantStatusInquiry(t *testing.T) {
	mockResponse := &domain.PayTRResponse{
		Status:  "success",
		Message: "Merchant status retrieved",
		Data: map[string]interface{}{
			"status":         "success",
			"payment_amount": "100.00",
			"currency":       "TRY",
		},
	}

	testService := setupTestService(mockResponse)

	req := domain.StatusInquiryRequest{
		MerchantOid: "test_order_123",
	}

	resp, err := testService.MerchantStatusInquiry(req)

	if err != nil {
		t.Fatalf("MerchantStatusInquiry returned an error: %v", err)
	}

	if resp.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", resp.Status)
	}
}

func TestAddNewCard(t *testing.T) {
	mockResponse := &domain.PayTRResponse{
		Status:  "success",
		Message: "Card added successfully",
		Data: map[string]interface{}{
			"token": "new_card_token",
		},
	}

	testService := setupTestService(mockResponse)

	req := domain.AddNewCardRequest{
		UserID:      "test_user",
		CardOwner:   "John Doe",
		CardNumber:  "4111111111111111",
		ExpiryMonth: "12",
		ExpiryYear:  "2025",
		CVV:         "123",
	}

	resp, err := testService.AddNewCard(req)

	if err != nil {
		t.Fatalf("AddNewCard returned an error: %v", err)
	}

	if resp.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", resp.Status)
	}
}

func TestGetSavedCards(t *testing.T) {
	mockResponse := &domain.PayTRResponse{
		Status:  "success",
		Message: "Saved cards retrieved",
		Data: map[string]interface{}{
			"cards": []map[string]interface{}{
				{
					"token":     "card_token_1",
					"last_four": "1111",
					"card_type": "visa",
				},
			},
		},
	}

	testService := setupTestService(mockResponse)

	resp, err := testService.GetSavedCards("test_user_token")

	if err != nil {
		t.Fatalf("GetSavedCards returned an error: %v", err)
	}

	if resp.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", resp.Status)
	}
}

func TestGetBinDetails(t *testing.T) {
	mockResponse := &domain.PayTRResponse{
		Status:  "success",
		Message: "BIN details retrieved",
		Data: map[string]interface{}{
			"bin_brand": "VISA",
			"bin_type":  "CREDIT",
		},
	}

	testService := setupTestService(mockResponse)

	resp, err := testService.GetBinDetails("411111")

	if err != nil {
		t.Fatalf("GetBinDetails returned an error: %v", err)
	}

	if resp.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", resp.Status)
	}

	if resp.Data["bin_brand"] != "VISA" {
		t.Errorf("Expected bin_brand 'VISA', got '%s'", resp.Data["bin_brand"])
	}
}

func TestDeleteSavedCard(t *testing.T) {
	mockResponse := &domain.PayTRResponse{
		Status:  "success",
		Message: "Card deleted successfully",
	}

	testService := setupTestService(mockResponse)

	resp, err := testService.DeleteSavedCard("test_user_token", "test_card_token")

	if err != nil {
		t.Fatalf("DeleteSavedCard returned an error: %v", err)
	}

	if resp.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", resp.Status)
	}
}
