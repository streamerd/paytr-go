package config

// PayTRConfig holds the configuration necessary to interact with PayTR's API,
// including the merchant's credentials.
type PayTRConfig struct {
	MerchantID   string
	MerchantKey  string
	MerchantSalt string
}
