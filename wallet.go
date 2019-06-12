package keybase

import (
	"encoding/json"
	"os/exec"
)

// ---- Struct for sending to API
type walletOut struct {
	Method string          `json:"method"`
	Params walletOutParams `json:"params"`
}
type walletOutOptions struct {
	Txid string `json:"txid"`
}
type walletOutParams struct {
	Options walletOutOptions `json:"options"`
}

// ----

// ---- Struct for data received after sending to API
type walletOutResult struct {
	Result WalletResult `json:"result"`
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
type sourceAsset struct {
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
	AccountID          string       `json:"accountID"`
	IsPrimary          bool         `json:"isPrimary"`
	Name               string       `json:"name"`
	Balance            []balance    `json:"balance"`
	ExchangeRate       exchangeRate `json:"exchangeRate"`
	AccountMode        int          `json:"accountMode"`
	TxID               string       `json:"txID"`
	Time               int64        `json:"time"`
	Status             string       `json:"status"`
	StatusDetail       string       `json:"statusDetail"`
	Amount             string       `json:"amount"`
	Asset              asset        `json:"asset"`
	DisplayAmount      string       `json:"displayAmount"`
	DisplayCurrency    string       `json:"displayCurrency"`
	SourceAmountMax    string       `json:"sourceAmountMax"`
	SourceAmountActual string       `json:"sourceAmountActual"`
	SourceAsset        sourceAsset  `json:"sourceAsset"`
	FromStellar        string       `json:"fromStellar"`
	ToStellar          string       `json:"toStellar"`
	FromUsername       string       `json:"fromUsername"`
	ToUsername         string       `json:"toUsername"`
	Note               string       `json:"note"`
	NoteErr            string       `json:"noteErr"`
	Unread             bool         `json:"unread"`
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

// TxDetail() returns details of a stellar transaction
func (k Keybase) TxDetail(txid string) (WalletResult, error) {
	m := walletOut{}
	m.Method = "details"
	m.Params.Options.Txid = txid

	r, err := walletAPIOut(k.Path, m)
	return r.Result, err
}
