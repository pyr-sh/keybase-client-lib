package api

import (
	"encoding/json"
	"os/exec"
	"strings"
)

// ---- Struct for sending to API
type chatOut struct {
	Method string        `json:"method"`
	Params chatOutParams `json:"params"`
}
type chatOutChannel struct {
	Name        string `json:"name"`
	MembersType string `json:"members_type"`
	TopicName   string `json:"topic_name"`
}
type chatOutMessage struct {
	Body string `json:"body"`
}
type chatOutOptions struct {
	Channel   chatOutChannel `json:"channel"`
	MessageID int            `json:"message_id"`
	Message   chatOutMessage `json:"message"`
}
type chatOutParams struct {
	Options chatOutOptions `json:"options"`
}

// ----

// ---- Struct for data received after sending to API
type chatOutResult struct {
	Result chatOutResultResult `json:"result"`
}
type chatOutResultRatelimits struct {
	Tank     string `json:"tank,omitempty"`
	Capacity int    `json:"capacity,omitempty"`
	Reset    int    `json:"reset,omitempty"`
	Gas      int    `json:"gas,omitempty"`
}
type chatOutResultChannel struct {
	Name        string `json:"name"`
	Public      bool   `json:"public"`
	MembersType string `json:"members_type"`
	TopicType   string `json:"topic_type,omitempty"`
	TopicName   string `json:"topic_name,omitempty"`
}
type chatOutResultConversations struct {
	ID           string               `json:"id"`
	Channel      chatOutResultChannel `json:"channel"`
	Unread       bool                 `json:"unread"`
	ActiveAt     int                  `json:"active_at"`
	ActiveAtMs   int64                `json:"active_at_ms"`
	MemberStatus string               `json:"member_status"`
}
type chatOutResultResult struct {
	Message       string                       `json:"message,omitempty"`
	ID            int                          `json:"id,omitempty"`
	Ratelimits    []chatOutResultRatelimits    `json:"ratelimits,omitempty"`
	Conversations []chatOutResultConversations `json:"conversations,omitempty"`
	Offline       bool                         `json:"offline,omitempty"`
}

// ----

// chatAPIOut() sends JSON requests to the chat API and returns its response.
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

// ChatSend() sends a chat message to a user.
func (k Keybase) ChatSendText(user string, message ...string) (chatOutResultResult, error) {
	m := chatOut{}
	m.Method = "send"
	m.Params.Options.Channel.Name = user
	m.Params.Options.Message.Body = strings.Join(message, " ")

	r, err := chatAPIOut(k.path, m)
	if err != nil {
		return chatOutResultResult{}, err
	}
	return r.Result, nil
}

// ChatSendTeam() sends a chat message to a team.
func (k Keybase) ChatSendTextTeam(team, channel string, message ...string) (chatOutResultResult, error) {
	m := chatOut{}
	m.Method = "send"
	m.Params.Options.Channel.Name = team
	m.Params.Options.Channel.MembersType = "team"
	m.Params.Options.Channel.TopicName = channel
	m.Params.Options.Message.Body = strings.Join(message, " ")

	r, err := chatAPIOut(k.path, m)
	if err != nil {
		return chatOutResultResult{}, err
	}
	return r.Result, nil
}

// ChatSendReaction() sends a reaction to a user's message.
func (k Keybase) ChatSendReaction(user, reaction string, messageId int) (chatOutResultResult, error) {
	m := chatOut{}
	m.Method = "reaction"
	m.Params.Options.Channel.Name = user
	m.Params.Options.MessageID = messageId
	m.Params.Options.Message.Body = reaction

	r, err := chatAPIOut(k.path, m)
	if err != nil {
		return chatOutResultResult{}, err
	}
	return r.Result, nil
}

// ChatSendReactionTeam() sends a reaction to a message on a team.
func (k Keybase) ChatSendReactionTeam(team, channel, reaction string, messageId int) (chatOutResultResult, error) {
	m := chatOut{}
	m.Method = "reaction"
	m.Params.Options.Channel.Name = team
	m.Params.Options.Channel.MembersType = "team"
	m.Params.Options.Channel.TopicName = channel
	m.Params.Options.MessageID = messageId
	m.Params.Options.Message.Body = reaction

	r, err := chatAPIOut(k.path, m)
	if err != nil {
		return chatOutResultResult{}, err
	}
	return r.Result, nil
}

// ChatList() returns a list of all conversations.
func (k Keybase) ChatList() ([]chatOutResultConversations, error) {
	m := chatOut{}
	m.Method = "list"

	r, err := chatAPIOut(k.path, m)
	return r.Result.Conversations, err
}
