package keybase

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

// Possible MemberTypes
const (
	TEAM string = "team"
	USER string = "impteamnative"
)

// Possible TopicTypes
const (
	DEV  string = "dev"
	CHAT string = "chat"
)

// NewKeybase returns a new Keybase. Optionally, you can pass a string containing the path to the Keybase executable as the first argument.
func NewKeybase(path ...string) *Keybase {
	k := &Keybase{}
	if len(path) < 1 {
		k.Path = "keybase"
	} else {
		k.Path = path[0]
	}

	s := k.status()
	k.Version = k.version()
	k.LoggedIn = s.LoggedIn
	if k.LoggedIn {
		k.Username = s.Username
		k.Device = s.Device.Name
	}
	return k
}

// NewBotCommand returns a new BotCommand instance
func NewBotCommand(name, description, usage string, extendedDescription ...BotCommandExtendedDescription) BotCommand {
	result := BotCommand{
		Name:        name,
		Description: description,
		Usage:       usage,
	}

	if len(extendedDescription) > 0 {
		result.ExtendedDescription = &extendedDescription[0]
	}

	return result
}

// NewBotCommandExtendedDescription
func NewBotCommandExtendedDescription(title, desktopBody, mobileBody string) BotCommandExtendedDescription {
	return BotCommandExtendedDescription{
		Title:       title,
		DesktopBody: desktopBody,
		MobileBody:  mobileBody,
	}
}

// Exec executes the given Keybase command
func (k *Keybase) Exec(command ...string) ([]byte, error) {
	out, err := exec.Command(k.Path, command...).Output()
	if err != nil {
		return []byte{}, err
	}
	return out, nil
}

// NewChat returns a new Chat instance
func (k *Keybase) NewChat(channel Channel) Chat {
	return Chat{
		keybase: k,
		Channel: channel,
	}
}

// NewTeam returns a new Team instance
func (k *Keybase) NewTeam(name string) Team {
	return Team{
		keybase: k,
		Name:    name,
	}
}

// NewWallet returns a new Wallet instance
func (k *Keybase) NewWallet() Wallet {
	return Wallet{
		keybase: k,
	}
}

// status returns the results of the `keybase status` command, which includes
// information about the client, and the currently logged-in Keybase user.
func (k *Keybase) status() status {
	cmdOut, err := k.Exec("status", "-j")
	if err != nil {
		return status{}
	}

	var s status
	json.Unmarshal(cmdOut, &s)

	return s
}

// version returns the version string of the client.
func (k *Keybase) version() string {
	cmdOut, err := k.Exec("version", "-S", "-f", "s")
	if err != nil {
		return ""
	}

	return string(cmdOut)
}

// UserLookup pulls information about users.
// The following fields are currently returned: basics, profile, proofs_summary, devices -- See https://keybase.io/docs/api/1.0/call/user/lookup for more info.
func (k *Keybase) UserLookup(users ...string) (UserAPI, error) {
	var fields = []string{"basics", "profile", "proofs_summary", "devices"}

	cmdOut, err := k.Exec("apicall", "--arg", fmt.Sprintf("usernames=%s", strings.Join(users, ",")), "--arg", fmt.Sprintf("fields=%s", strings.Join(fields, ",")), "user/lookup")
	if err != nil {
		return UserAPI{}, err
	}

	var r UserAPI
	if err := json.Unmarshal(cmdOut, &r); err != nil {
		return UserAPI{}, err
	}

	return r, nil
}
