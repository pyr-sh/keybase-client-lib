package keybase

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// walletAPIOut sends JSON requests to the wallet API and returns its response.
func walletAPIOut(k *Keybase, w WalletAPI) (WalletAPI, error) {
	jsonBytes, _ := json.Marshal(w)

	cmdOut, err := k.Exec("wallet", "api", "-m", string(jsonBytes))
	if err != nil {
		return WalletAPI{}, err
	}

	var r WalletAPI
	json.Unmarshal(cmdOut, &r)
	if r.Error != nil {
		return WalletAPI{}, errors.New(r.Error.Message)
	}
	return r, nil
}

// TxDetail returns details of a stellar transaction
func (w Wallet) TxDetail(txid string) (WalletAPI, error) {
	m := WalletAPI{
		Params: &wParams{},
	}
	m.Method = "details"
	m.Params.Options.Txid = txid

	r, err := walletAPIOut(w.keybase, m)
	return r, err
}

// StellarAddress returns the primary stellar address of a given user
func (w Wallet) StellarAddress(user string) (string, error) {
	m := WalletAPI{
		Params: &wParams{},
	}
	m.Method = "lookup"
	m.Params.Options.Name = user

	r, err := walletAPIOut(w.keybase, m)
	if err != nil {
		return "", err
	}
	return r.Result.AccountID, err
}

// StellarUser returns the keybase username of a given wallet address
func (w Wallet) StellarUser(wallet string) (string, error) {
	m := WalletAPI{
		Params: &wParams{},
	}
	m.Method = "lookup"
	m.Params.Options.Name = wallet

	r, err := walletAPIOut(w.keybase, m)
	if err != nil {
		return "", err
	}
	return r.Result.Username, err
}

// RequestPayment sends a request for payment to a user
func (w Wallet) RequestPayment(user string, amount float64, memo ...string) error {
	k := w.keybase
	if len(memo) > 0 {
		_, err := k.Exec("wallet", "request", user, fmt.Sprintf("%f", amount), "-m", memo[0])
		return err
	}
	_, err := k.Exec("wallet", "request", user, fmt.Sprintf("%f", amount))
	return err
}

// CancelRequest cancels a request for payment previously sent to a user
func (w Wallet) CancelRequest(requestID string) error {
	k := w.keybase
	_, err := k.Exec("wallet", "cancel-request", requestID)
	return err
}

// Send sends the specified amount of the specified currency to a user
func (w Wallet) Send(recipient string, amount string, currency string, message ...string) (WalletAPI, error) {
	m := WalletAPI{
		Params: &wParams{},
	}
	m.Method = "send"
	m.Params.Options.Recipient = recipient
	m.Params.Options.Amount = amount
	m.Params.Options.Currency = currency
	if len(message) > 0 {
		m.Params.Options.Message = strings.Join(message, " ")
	}

	r, err := walletAPIOut(w.keybase, m)
	if err != nil {
		return WalletAPI{}, err
	}
	return r, err
}

// SendXLM sends the specified amount of XLM to a user
func (w Wallet) SendXLM(recipient string, amount string, message ...string) (WalletAPI, error) {
	result, err := w.Send(recipient, amount, "XLM", message...)
	return result, err
}
