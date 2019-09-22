package keybase

import (
	"encoding/json"
	"errors"
	"fmt"
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
func (k *Keybase) TxDetail(txid string) (WalletAPI, error) {
	m := WalletAPI{
		Params: &wParams{},
	}
	m.Method = "details"
	m.Params.Options.Txid = txid

	r, err := walletAPIOut(k, m)
	return r, err
}

// StellarAddress returns the primary stellar address of a given user
func (k *Keybase) StellarAddress(user string) (string, error) {
	m := WalletAPI{
		Params: &wParams{},
	}
	m.Method = "lookup"
	m.Params.Options.Name = user

	r, err := walletAPIOut(k, m)
	if err != nil {
		return "", err
	}
	return r.Result.AccountID, err
}

// RequestPayment sends a request for payment to a user
func (k *Keybase) RequestPayment(user string, amount float64, memo ...string) error {
	if len(memo) > 0 {
		_, err := k.Exec("wallet", "request", user, fmt.Sprintf("%f", amount), "-m", memo[0])
		return err
	}
	_, err := k.Exec("wallet", "request", user, fmt.Sprintf("%f", amount))
	return err
}

// CancelRequest cancels a request for payment previously sent to a user
func (k *Keybase) CancelRequest(requestID string) error {
	_, err := k.Exec("wallet", "cancel-request", requestID)
	return err
}
