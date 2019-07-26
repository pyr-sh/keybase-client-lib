package keybase

import (
	"encoding/json"
	"os/exec"
)

// walletAPIOut sends JSON requests to the wallet API and returns its response.
func walletAPIOut(keybasePath string, w WalletAPI) (WalletAPI, error) {
	jsonBytes, _ := json.Marshal(w)

	cmd := exec.Command(keybasePath, "wallet", "api", "-m", string(jsonBytes))
	cmdOut, err := cmd.Output()
	if err != nil {
		return WalletAPI{}, err
	}

	var r WalletAPI
	json.Unmarshal(cmdOut, &r)

	return r, nil
}

// TxDetail returns details of a stellar transaction
func (k *Keybase) TxDetail(txid string) (wResult, error) {
	m := WalletAPI{}
	m.Method = "details"
	m.Params.Options.Txid = txid

	r, err := walletAPIOut(k.Path, m)
	return r.Result, err
}

// StellarAddress returns the primary stellar address of a given user
func (k *Keybase) StellarAddress(user string) (string, error) {
	m := WalletAPI{}
	m.Method = "lookup"
	m.Params.Options.Name = user

	r, err := walletAPIOut(k.Path, m)
	return r.Result.AccountID, err
}
