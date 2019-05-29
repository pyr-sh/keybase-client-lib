// keybase/api implements an interface for interacting with the Keybase Chat, Team, and Wallet APIs
package api

import (
	"encoding/json"
	"os/exec"
)

type Keybase struct {
	path string
}

type keybase interface {
	ChatSend(user, message string) (chatOutResult, error)
	ChatSendTeam(team, channel, message string) (chatOutResult, error)
	ChatList() (chatOutResult, error)
	LoggedIn() bool
	Username() string
	Version() string
}

type status struct {
	Username string `json:"Username"`
	LoggedIn bool   `json:"LoggedIn"`
}

// New() returns a new instance of Keybase object. Optionally, you can pass a string containing the path to the Keybase executable as the first argument.
func New(path ...string) Keybase {
	if len(path) < 1 {
		return Keybase{path: "/usr/bin/keybase"}
	}
	return Keybase{path: path[0]}
}

// Username() returns the username of the currently logged-in Keybase user.
func (k Keybase) Username() string {
	cmd := exec.Command(k.path, "status", "-j")
	cmdOut, err := cmd.Output()
	if err != nil {
		return ""
	}

	var s status
	json.Unmarshal(cmdOut, &s)

	return s.Username
}

// LoggedIn() returns true if Keybase is currently logged in, otherwise returns false.
func (k Keybase) LoggedIn() bool {
	cmd := exec.Command(k.path, "status", "-j")
	cmdOut, err := cmd.Output()
	if err != nil {
		return false
	}

	var s status
	json.Unmarshal(cmdOut, &s)

	return s.LoggedIn
}

// Version() returns the version string of the client.
func (k Keybase) Version() string {
	cmd := exec.Command(k.path, "version", "-S", "-f", "s")
	cmdOut, err := cmd.Output()
	if err != nil {
		return ""
	}

	return string(cmdOut)
}
