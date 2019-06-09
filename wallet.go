package keybase

import (
	"encoding/json"
	"os/exec"
)

// ---- Struct for sending to API
type walletOut struct {
	Method string `json:"method"`
}

// ----

// ---- Struct for data received after sending to API
type walletOutResult struct {
	Result []WalletResult `json:"result"`
}
type asset struct {
	Type           string `json:"type"`
	Code           string `json:"code"`
	Issuer         string `json:"issuer"`
	VerifiedDomain string `json:"verifiedDomain"`
	IssuerName     string `json:"issuerName"`
	Desc           string `json:"desc"`
	InfoURL        string `json:"infoUrl"`
}
type balance struct {
	Asset  asset  `json:"asset"`
	Amount string `json:"amount"`
	Limit  string `json:"limit"`
}
type exchangeRate struct {
	Currency string `json:"currency"`
	Rate     string `json:"rate"`
}
type WalletResult struct {
	AccountID    string       `json:"accountID"`
	IsPrimary    bool         `json:"isPrimary"`
	Name         string       `json:"name"`
	Balance      []balance    `json:"balance"`
	ExchangeRate exchangeRate `json:"exchangeRate"`
	AccountMode  int          `json:"accountMode"`
}

// ----

// walletAPIOut() sends JSON requests to the wallet API and returns its response.
func walletAPIOut(keybasePath string, w walletOut) (walletOutResult, error) {
	jsonBytes, _ := json.Marshal(w)

	cmd := exec.Command(keybasePath, "wallet", "api", "-m", string(jsonBytes))
	cmdOut, err := cmd.Output()
	if err != nil {
		return walletOutResult{}, err
	}

	var r walletOutResult
	json.Unmarshal(cmdOut, &r)

	return r, nil
}

// Balances() returns a list of all wallets on the current user's account, including their current balance.
func (k Keybase) Balances() ([]WalletResult, error) {
	m := walletOut{}
	m.Method = "balances"

	r, err := walletAPIOut(k.Path, m)
	return r.Result, err
}
