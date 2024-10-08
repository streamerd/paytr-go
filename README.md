# Payment Package for PayTR Payment System

The `payment` package provides services for interacting with the PayTR payment system. This package includes operations such as card payments, recurring payments, refund transactions, and card management.

This guide explains how to integrate with the PayTR API and how to use this package.


## Features

- New card payment transactions
- Saved card payment transactions
- Recurring payment transactions
- Refund transactions
- Adding new cards
- Viewing and deleting saved cards
- Retrieving BIN (Bank Identification Number) information


## Installation

To install the PayTR payment package, use the following `go get` command in the directory containing the go.mod file:

```bash
go get github.com/streamerd/paytr-go@v1.0.0
```

This command will download and install version 1.0.0 of the package in your Go workspace.


## Usage

### 1. Configuration Setup

First, you need to create a `PayTRConfig` structure that provides the necessary configuration information to interact with the PayTR API:

```go
type PayTRConfig struct {
    MerchantID   string // Unique identifier of the merchant.
    MerchantKey  string // Secret key used for generating HMAC.
    MerchantSalt string // Salt used to enhance security in token generation.
}
```

### 2. Creating a New Service

You can create a PayTR service using the `payment` package:

```go
svc := payment.NewService(PayTRConfig{
    MerchantID:   "your-merchant-id",
    MerchantKey:  "your-merchant-key",
    MerchantSalt: "your-merchant-salt",
})
```

This service is used to perform card transactions, refunds, and card management operations with the PayTR API.

### 3. Card Payment Transaction

To make a new payment with a card, you can use the `NewCardPayment` method:

```go
req := domain.NewCardPaymentRequest{
    // Enter payment request data
}

resp, err := svc.NewCardPayment(req)
if err != nil {
    // Error handling
}
```

### 4. Saved Card Payment Transaction

To make a payment with a previously saved card, you can use the `SavedCardPayment` method:

```go
req := domain.SavedCardPaymentRequest{
    // Enter saved card payment information
}

resp, err := svc.SavedCardPayment(req)
if err != nil {
    // Error handling
}
```

### 5. Recurring Payment

To process recurring payments, you can use the `RecurringPayment` method:

```go
req := domain.SavedCardPaymentRequest{
    // Enter recurring payment information
}

resp, err := svc.RecurringPayment(req)
if err != nil {
    // Error handling
}
```

### 6. Refund Transaction

To make a refund for a payment transaction, you can use the `RefundPayment` method:

```go
req := domain.RefundRequest{
    MerchantOid:  "transaction-id",
    ReturnAmount: 100.0, // Refund amount
}

resp, err := svc.RefundPayment(req)
if err != nil {
    // Error handling
}
```

### 7. Card Management

- Adding a new card: You can add a new card to a user account using the `AddNewCard` method.
- Listing saved cards: The `GetSavedCards` method allows you to view a user's saved cards.
- Deleting a saved card: You can delete a card using the `DeleteSavedCard` method.

```go
// Adding a new card
addCardReq := domain.AddNewCardRequest{
    // Add card information
}
resp, err := svc.AddNewCard(addCardReq)

// Listing saved cards
savedCardsResp, err := svc.GetSavedCards("user-token")

// Deleting a saved card
resp, err := svc.DeleteSavedCard("user-token", "card-token")
```

### 8. BIN (Card Number) Details

You can use the `GetBinDetails` method to get BIN details using the first 6-8 digits of a card:

```go
binDetails, err := svc.GetBinDetails("123456")
```

## HMAC Signature Generation

HMAC is used for security in requests to the PayTR API. The signature is generated by combining the request data and creating an HMAC with SHA-256. For example:

```go
hashStr := fmt.Sprintf("%s%s%.2f", config.MerchantID, req.MerchantOid, req.ReturnAmount)
hmac := hmac.New(sha256.New, []byte(config.MerchantKey))
hmac.Write([]byte(hashStr))
signature := base64.StdEncoding.EncodeToString(hmac.Sum(nil))
```

## Conclusion

This package simplifies operations such as card payments, refund transactions, and card management with the PayTR API. By following the steps explained with example codes, you can integrate secure payment transactions with PayTR into your applications.