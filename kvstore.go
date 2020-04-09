package keybase

import (
	"encoding/json"
	"fmt"

	"samhofi.us/x/keybase/v2/types/keybase1"
)

// KVListNamespaces returns all namespaces for a team
func (k *Keybase) KVListNamespaces(team *string) (keybase1.KVListNamespaceResult, error) {
	type res struct {
		Result keybase1.KVListNamespaceResult `json:"result"`
		Error  *Error                         `json:"error"`
	}
	var r res

	arg := newKVArg("list", KVOptions{
		Team: team,
	})

	jsonBytes, _ := json.Marshal(arg)

	cmdOut, err := k.Exec("kvstore", "api", "-m", string(jsonBytes))
	if err != nil {
		return r.Result, err
	}

	err = json.Unmarshal(cmdOut, &r)
	if err != nil {
		return r.Result, err
	}

	if r.Error != nil {
		return r.Result, fmt.Errorf("%s", r.Error.Message)
	}

	return r.Result, nil
}

// KVListKeys returns all non-deleted keys for a namespace
func (k *Keybase) KVListKeys(team *string, namespace string) (keybase1.KVListEntryResult, error) {
	type res struct {
		Result keybase1.KVListEntryResult `json:"result"`
		Error  *Error                     `json:"error"`
	}
	var r res

	arg := newKVArg("list", KVOptions{
		Team:      team,
		Namespace: &namespace,
	})

	jsonBytes, _ := json.Marshal(arg)

	cmdOut, err := k.Exec("kvstore", "api", "-m", string(jsonBytes))
	if err != nil {
		return r.Result, err
	}

	err = json.Unmarshal(cmdOut, &r)
	if err != nil {
		return r.Result, err
	}

	if r.Error != nil {
		return r.Result, fmt.Errorf("%s", r.Error.Message)
	}

	return r.Result, nil
}

// KVGet returns an entry
func (k *Keybase) KVGet(team *string, namespace string, key string) (keybase1.KVGetResult, error) {
	type res struct {
		Result keybase1.KVGetResult `json:"result"`
		Error  *Error               `json:"error"`
	}
	var r res

	arg := newKVArg("get", KVOptions{
		Team:      team,
		Namespace: &namespace,
		EntryKey:  &key,
	})

	jsonBytes, _ := json.Marshal(arg)

	cmdOut, err := k.Exec("kvstore", "api", "-m", string(jsonBytes))
	if err != nil {
		return r.Result, err
	}

	err = json.Unmarshal(cmdOut, &r)
	if err != nil {
		return r.Result, err
	}

	if r.Error != nil {
		return r.Result, fmt.Errorf("%s", r.Error.Message)
	}

	return r.Result, nil
}

// KVPutWithRevision puts an entry, specifying the revision number
func (k *Keybase) KVPutWithRevision(team *string, namespace string, key string, value string, revision int) (keybase1.KVPutResult, error) {
	type res struct {
		Result keybase1.KVPutResult `json:"result"`
		Error  *Error               `json:"error"`
	}
	var r res

	opts := KVOptions{
		Team:       team,
		Namespace:  &namespace,
		EntryKey:   &key,
		EntryValue: &value,
	}
	if revision != 0 {
		opts.Revision = &revision
	}

	arg := newKVArg("put", opts)

	jsonBytes, _ := json.Marshal(arg)

	cmdOut, err := k.Exec("kvstore", "api", "-m", string(jsonBytes))
	if err != nil {
		return r.Result, err
	}

	err = json.Unmarshal(cmdOut, &r)
	if err != nil {
		return r.Result, err
	}

	if r.Error != nil {
		return r.Result, fmt.Errorf("%s", r.Error.Message)
	}

	return r.Result, nil
}

// KVPut puts an entry
func (k *Keybase) KVPut(team *string, namespace string, key string, value string) (keybase1.KVPutResult, error) {
	return k.KVPutWithRevision(team, namespace, key, value, 0)
}

// KVDeleteWithRevision deletes an entry, specifying the revision number
func (k *Keybase) KVDeleteWithRevision(team *string, namespace string, key string, revision int) (keybase1.KVDeleteEntryResult, error) {
	type res struct {
		Result keybase1.KVDeleteEntryResult `json:"result"`
		Error  *Error                       `json:"error"`
	}
	var r res

	opts := KVOptions{
		Team:      team,
		Namespace: &namespace,
		EntryKey:  &key,
	}
	if revision != 0 {
		opts.Revision = &revision
	}

	arg := newKVArg("del", opts)

	jsonBytes, _ := json.Marshal(arg)

	cmdOut, err := k.Exec("kvstore", "api", "-m", string(jsonBytes))
	if err != nil {
		return r.Result, err
	}

	err = json.Unmarshal(cmdOut, &r)
	if err != nil {
		return r.Result, err
	}

	if r.Error != nil {
		return r.Result, fmt.Errorf("%s", r.Error.Message)
	}

	return r.Result, nil
}

// KVDelete deletes an entry
func (k *Keybase) KVDelete(team *string, namespace string, key string) (keybase1.KVDeleteEntryResult, error) {
	return k.KVDeleteWithRevision(team, namespace, key, 0)
}
