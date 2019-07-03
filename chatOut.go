package keybase

import (
	"encoding/json"
	"os/exec"
	"strings"
)

// chatOut holds data to be sent to the Chat API
type chatOut struct {
	Method string        `json:"method"`
	Params chatOutParams `json:"params"`
}
type Channel struct {
	Name        string `json:"name"`
	Public      bool   `json:"public,omitempty"`
	MembersType string `json:"members_type,omitempty"`
	TopicType   string `json:"topic_type,omitempty"`
	TopicName   string `json:"topic_name,omitempty"`
}
type chatOutMessage struct {
	Body string `json:"body"`
}
type chatOutOptions struct {
	Channel   Channel        `json:"channel"`
	MessageID int            `json:"message_id"`
	Message   chatOutMessage `json:"message"`
}
type chatOutParams struct {
	Options chatOutOptions `json:"options"`
}

// chatOutResult holds data received after sending to API
type chatOutResult struct {
	Result ChatOut `json:"result"`
}
type chatOutResultRatelimits struct {
	Tank     string `json:"tank,omitempty"`
	Capacity int    `json:"capacity,omitempty"`
	Reset    int    `json:"reset,omitempty"`
	Gas      int    `json:"gas,omitempty"`
}
type conversation struct {
	ID           string  `json:"id"`
	Channel      Channel `json:"channel"`
	Unread       bool    `json:"unread"`
	ActiveAt     int     `json:"active_at"`
	ActiveAtMs   int64   `json:"active_at_ms"`
	MemberStatus string  `json:"member_status"`
}
type ChatOut struct {
	Message       string                    `json:"message,omitempty"`
	ID            int                       `json:"id,omitempty"`
	Ratelimits    []chatOutResultRatelimits `json:"ratelimits,omitempty"`
	Conversations []conversation            `json:"conversations,omitempty"`
	Offline       bool                      `json:"offline,omitempty"`
}

// chatAPIOut sends JSON requests to the chat API and returns its response.
func chatAPIOut(keybasePath string, c chatOut) (chatOutResult, error) {
	jsonBytes, _ := json.Marshal(c)

	cmd := exec.Command(keybasePath, "chat", "api", "-m", string(jsonBytes))
	cmdOut, err := cmd.Output()
	if err != nil {
		return chatOutResult{}, err
	}

	var r chatOutResult
	json.Unmarshal(cmdOut, &r)

	return r, nil
}

// Send sends a chat message
func (c Chat) Send(message ...string) (ChatOut, error) {
	m := chatOut{}
	m.Method = "send"
	m.Params.Options.Channel = c.Channel
	m.Params.Options.Message.Body = strings.Join(message, " ")

	r, err := chatAPIOut(c.keybase.Path, m)
	if err != nil {
		return ChatOut{}, err
	}
	return r.Result, nil
}

// Edit edits a previously sent chat message
func (c Chat) Edit(messageId int, message ...string) (ChatOut, error) {
	m := chatOut{}
	m.Method = "edit"
	m.Params.Options.Channel = c.Channel
	m.Params.Options.Message.Body = strings.Join(message, " ")
	m.Params.Options.MessageID = messageId

	r, err := chatAPIOut(c.keybase.Path, m)
	if err != nil {
		return ChatOut{}, err
	}
	return r.Result, nil
}

// React sends a reaction to a message.
func (c Chat) React(messageId int, reaction string) (ChatOut, error) {
	m := chatOut{}
	m.Method = "reaction"
	m.Params.Options.Channel = c.Channel
	m.Params.Options.Message.Body = reaction
	m.Params.Options.MessageID = messageId

	r, err := chatAPIOut(c.keybase.Path, m)
	if err != nil {
		return ChatOut{}, err
	}
	return r.Result, nil
}

// Delete deletes a chat message
func (c Chat) Delete(messageId int) (ChatOut, error) {
	m := chatOut{}
	m.Method = "delete"
	m.Params.Options.Channel = c.Channel
	m.Params.Options.MessageID = messageId

	r, err := chatAPIOut(c.keybase.Path, m)
	if err != nil {
		return ChatOut{}, err
	}
	return r.Result, nil
}

// ChatList returns a list of all conversations.
func (k *Keybase) ChatList() ([]conversation, error) {
	m := chatOut{}
	m.Method = "list"

	r, err := chatAPIOut(k.Path, m)
	return r.Result.Conversations, err
}
