package keybase

import (
	"encoding/json"
	"os/exec"
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

// Keybase holds basic information about the local Keybase executable
type Keybase struct {
	Path     string
	Username string
	LoggedIn bool
	Version  string
}

// Channel is a map of options that can be passed to NewChat()
type Channel map[string]interface{}

// Chat holds basic information about a specific conversation
type Chat struct {
	keybase     Keybase
	Name        string
	Public      bool
	MembersType string
	TopicName   string
	TopicType   string
}

type chat interface {
	Send(message ...string) (ChatOut, error)
	React(messageId int, reaction string) (ChatOut, error)
	Delete(messageId int) (ChatOut, error)
}

type keybase interface {
	NewChat(channel map[string]interface{}) Chat
	ChatList() ([]conversation, error)
	loggedIn() bool
	username() string
	version() string
}

type status struct {
	Username string `json:"Username"`
	LoggedIn bool   `json:"LoggedIn"`
}

// New() returns a new instance of Keybase object. Optionally, you can pass a string containing the path to the Keybase executable as the first argument.
func NewKeybase(path ...string) Keybase {
	k := Keybase{}
	if len(path) < 1 {
		k.Path = "keybase"
	} else {
		k.Path = path[0]
	}
	k.Version = k.version()
	k.LoggedIn = k.loggedIn()
	if k.LoggedIn == true {
		k.Username = k.username()
	}
	return k
}

// Return a new Chat instance
func (k Keybase) NewChat(channel map[string]interface{}) Chat {
	var c Chat = Chat{}
	c.keybase = k
	if value, ok := channel["Name"].(string); ok == true {
		c.Name = value
	}
	if value, ok := channel["Public"].(bool); ok == true {
		c.Public = value
	} else {
		c.Public = false
	}
	if value, ok := channel["MembersType"].(string); ok == true {
		c.MembersType = value
	} else {
		c.MembersType = USER
	}
	if value, ok := channel["TopicName"].(string); ok == true {
		c.TopicName = value
	} else {
		if c.MembersType == TEAM {
			c.TopicName = "general"
		}
	}
	if value, ok := channel["TopicType"].(string); ok == true {
		c.TopicType = value
	} else {
		c.TopicType = CHAT
	}
	return c
}

// username() returns the username of the currently logged-in Keybase user.
func (k Keybase) username() string {
	cmd := exec.Command(k.Path, "status", "-j")
	cmdOut, err := cmd.Output()
	if err != nil {
		return ""
	}

	var s status
	json.Unmarshal(cmdOut, &s)

	return s.Username
}

// loggedIn() returns true if Keybase is currently logged in, otherwise returns false.
func (k Keybase) loggedIn() bool {
	cmd := exec.Command(k.Path, "status", "-j")
	cmdOut, err := cmd.Output()
	if err != nil {
		return false
	}

	var s status
	json.Unmarshal(cmdOut, &s)

	return s.LoggedIn
}

// version() returns the version string of the client.
func (k Keybase) version() string {
	cmd := exec.Command(k.Path, "version", "-S", "-f", "s")
	cmdOut, err := cmd.Output()
	if err != nil {
		return ""
	}

	return string(cmdOut)
}
