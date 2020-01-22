package keybase

import (
	"encoding/json"
	"errors"
)

// kvAPIOut sends a JSON request to the kvstore API and returns its response.
func kvAPIOut(k *Keybase, kv KVAPI) (KVAPI, error) {
	jsonBytes, _ := json.Marshal(kv)

	cmdOut, err := k.Exec("kvstore", "api", "-m", string(jsonBytes))
	if err != nil {
		return KVAPI{}, err
	}

	var r KVAPI
	if err := json.Unmarshal(cmdOut, &r); err != nil {
		return KVAPI{}, err
	}

	if r.Error != nil {
		return KVAPI{}, errors.New(r.Error.Message)
	}

	return r, nil
}

// Namespaces returns all namespaces for a team
func (kv KV) Namespaces() (KVAPI, error) {
	m := KVAPI{
		Params: &kvParams{},
	}
	m.Params.Options = kvOptions{
		Team: kv.Team,
	}

	m.Method = "list"

	r, err := kvAPIOut(kv.keybase, m)
	if err != nil {
		return r, err
	}
	return r, nil
}

// Keys returns all non-deleted keys for a namespace
func (kv KV) Keys(namespace string) (KVAPI, error) {
	m := KVAPI{
		Params: &kvParams{},
	}
	m.Params.Options = kvOptions{
		Team:      kv.Team,
		Namespace: namespace,
	}

	m.Method = "list"

	r, err := kvAPIOut(kv.keybase, m)
	if err != nil {
		return r, err
	}
	return r, nil
}

// Get returns an entry
func (kv KV) Get(namespace string, key string) (KVAPI, error) {
	m := KVAPI{
		Params: &kvParams{},
	}
	m.Params.Options = kvOptions{
		Team:      kv.Team,
		Namespace: namespace,
		EntryKey:  key,
	}

	m.Method = "get"

	r, err := kvAPIOut(kv.keybase, m)
	if err != nil {
		return r, err
	}
	return r, nil
}

// Put adds an entry
func (kv KV) Put(namespace string, key string, value string) (KVAPI, error) {
	m := KVAPI{
		Params: &kvParams{},
	}
	m.Params.Options = kvOptions{
		Team:       kv.Team,
		Namespace:  namespace,
		EntryKey:   key,
		EntryValue: value,
	}

	m.Method = "put"

	r, err := kvAPIOut(kv.keybase, m)
	if err != nil {
		return r, err
	}
	return r, nil
}
