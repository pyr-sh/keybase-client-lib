package keybase

import (
	"encoding/json"
	"os/exec"
	"strings"
)

// chatAPIOut sends JSON requests to the chat API and returns its response.
func chatAPIOut(keybasePath string, c ChatAPI) (ChatAPI, error) {
	jsonBytes, _ := json.Marshal(c)

	cmd := exec.Command(keybasePath, "chat", "api", "-m", string(jsonBytes))
	cmdOut, err := cmd.Output()
	if err != nil {
		return ChatAPI{}, err
	}

	var r ChatAPI
	if err := json.Unmarshal(cmdOut, &r); err != nil {
		return ChatAPI{}, err
	}

	return r, nil
}

// Send sends a chat message
func (c Chat) Send(message ...string) (ChatAPI, error) {
	m := ChatAPI{}
	m.Method = "send"
	m.Params.Options.Channel = c.Channel
	m.Params.Options.Message.Body = strings.Join(message, " ")

	r, err := chatAPIOut(c.keybase.Path, m)
	if err != nil {
		return ChatAPI{}, err
	}
	return r, nil
}

// Edit edits a previously sent chat message
func (c Chat) Edit(messageId int, message ...string) (ChatAPI, error) {
	m := ChatAPI{}
	m.Method = "edit"
	m.Params.Options.Channel = c.Channel
	m.Params.Options.Message.Body = strings.Join(message, " ")
	m.Params.Options.MessageID = messageId

	r, err := chatAPIOut(c.keybase.Path, m)
	if err != nil {
		return ChatAPI{}, err
	}
	return r, nil
}

// React sends a reaction to a message.
func (c Chat) React(messageId int, reaction string) (ChatAPI, error) {
	m := ChatAPI{}
	m.Method = "reaction"
	m.Params.Options.Channel = c.Channel
	m.Params.Options.Message.Body = reaction
	m.Params.Options.MessageID = messageId

	r, err := chatAPIOut(c.keybase.Path, m)
	if err != nil {
		return ChatAPI{}, err
	}
	return r, nil
}

// Delete deletes a chat message
func (c Chat) Delete(messageId int) (ChatAPI, error) {
	m := ChatAPI{}
	m.Method = "delete"
	m.Params.Options.Channel = c.Channel
	m.Params.Options.MessageID = messageId

	r, err := chatAPIOut(c.keybase.Path, m)
	if err != nil {
		return ChatAPI{}, err
	}
	return r, nil
}

// ChatList returns a list of all conversations.
func (k *Keybase) ChatList() (ChatAPI, error) {
	m := ChatAPI{}
	m.Method = "list"

	r, err := chatAPIOut(k.Path, m)
	return r, err
}
