package api

import (
	"encoding/json"
	"os/exec"
)


// ---- Struct for sending to API
type chatOut struct {
	Method string `json:"method"`
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
	Channel chatOutChannel `json:"channel"`
	Message chatOutMessage `json:"message"`
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
	Tank     string `json:"tank"`
	Capacity int    `json:"capacity"`
	Reset    int    `json:"reset"`
	Gas      int    `json:"gas"`
}
type chatOutResultResult struct {
	Message    string       `json:"message"`
	ID         int          `json:"id"`
	Ratelimits []chatOutResultRatelimits `json:"ratelimits"`
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
func (k Keybase) ChatSend(user, message string) (chatOutResult, error) {
	m := chatOut{}
	m.Method = "send"
	m.Params.Options.Channel.Name = user
	m.Params.Options.Message.Body = message

	return chatAPIOut(k.path, m)
}

// ChatSendTeam() sends a chat message to a team.
func (k Keybase) ChatSendTeam(team, channel, message string) (chatOutResult, error) {
	m := chatOut{}
	m.Method = "send"
	m.Params.Options.Channel.Name = team
	m.Params.Options.Channel.MembersType = "team"
	m.Params.Options.Channel.TopicName = channel
	m.Params.Options.Message.Body = message

	return chatAPIOut(k.path, m)
}
