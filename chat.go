package keybase

import (
	"bufio"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"time"

	"samhofi.us/x/keybase/types/chat1"
)

// Returns a string representation of a message id suitable for use in a
// pagination struct
func getID(id uint) string {
	var b []byte
	switch {
	case id < 128:
		// 7-bit int
		b = make([]byte, 1)
		b = []byte{byte(id)}

	case id <= 255:
		// uint8
		b = make([]byte, 2)
		b = []byte{204, byte(id)}

	case id > 255 && id <= 65535:
		// uint16
		b = make([]byte, 2)
		binary.BigEndian.PutUint16(b, uint16(id))
		b = append([]byte{205}, b...)

	case id > 65535 && id <= 4294967295:
		// uint32
		b = make([]byte, 4)
		binary.BigEndian.PutUint32(b, uint32(id))
		b = append([]byte{206}, b...)
	}
	return base64.StdEncoding.EncodeToString(b)
}

// Creates a string of a json-encoded channel to pass to keybase chat api-listen --filter-channel
func createFilterString(channel chat1.ChatChannel) string {
	if channel.Name == "" {
		return ""
	}
	jsonBytes, _ := json.Marshal(channel)
	return string(jsonBytes)
}

// Creates a string of json-encoded channels to pass to keybase chat api-listen --filter-channels
func createFiltersString(channels []chat1.ChatChannel) string {
	if len(channels) == 0 {
		return ""
	}
	jsonBytes, _ := json.Marshal(channels)
	return string(jsonBytes)
}

// Run `keybase chat api-listen` to get new messages coming into keybase and send them into the channel
func getNewMessages(k *Keybase, subs *SubscriptionChannels, execOptions []string) {
	execString := []string{"chat", "api-listen"}
	if len(execOptions) > 0 {
		execString = append(execString, execOptions...)
	}
	for {
		execCmd := exec.Command(k.Path, execString...)
		stdOut, _ := execCmd.StdoutPipe()
		execCmd.Start()
		scanner := bufio.NewScanner(stdOut)
		go func(scanner *bufio.Scanner, subs *SubscriptionChannels) {
			for {
				scanner.Scan()
				var subType subscriptionType
				t := scanner.Text()
				json.Unmarshal([]byte(t), &subType)
				switch subType.Type {
				case "chat":
					var notification chat1.MsgNotification
					if err := json.Unmarshal([]byte(t), &notification); err != nil {
						subs.error <- err
						break
					}
					if notification.Msg != nil {
						subscriptionMessage := SubscriptionMessage{
							Message: *notification.Msg,
							Conversation: chat1.ConvSummary{
								Id:      notification.Msg.ConvID,
								Channel: notification.Msg.Channel,
							},
						}
						subs.chat <- subscriptionMessage
					}
				case "chat_conv":
					var notification chat1.ConvNotification
					if err := json.Unmarshal([]byte(t), &notification); err != nil {
						subs.error <- err
						break
					}
					if notification.Conv != nil {
						subscriptionConv := SubscriptionConversation{
							Conversation: *notification.Conv,
						}
						subs.conversation <- subscriptionConv
					}
				case "wallet":
					var holder paymentHolder
					if err := json.Unmarshal([]byte(t), &holder); err != nil {
						subs.error <- err
						break
					}
					subscriptionPayment := SubscriptionWalletEvent(holder)
					subs.wallet <- subscriptionPayment
				default:
					continue
				}
			}
		}(scanner, subs)
		execCmd.Wait()
	}
}

// Run runs `keybase chat api-listen`, and passes incoming messages to the message handler func
func (k *Keybase) Run(handlers Handlers, options ...RunOptions) {
	var channelCapacity = 100

	runOptions := make([]string, 0)
	if len(options) > 0 {
		if options[0].Capacity > 0 {
			channelCapacity = options[0].Capacity
		}
		if options[0].Wallet {
			runOptions = append(runOptions, "--wallet")
		}
		if options[0].Convs {
			runOptions = append(runOptions, "--convs")
		}
		if options[0].Local {
			runOptions = append(runOptions, "--local")
		}
		if options[0].HideExploding {
			runOptions = append(runOptions, "--hide-exploding")
		}
		if options[0].Dev {
			runOptions = append(runOptions, "--dev")
		}
		if len(options[0].FilterChannels) > 0 {
			runOptions = append(runOptions, "--filter-channels")
			runOptions = append(runOptions, createFiltersString(options[0].FilterChannels))

		}
		if options[0].FilterChannel.Name != "" {
			runOptions = append(runOptions, "--filter-channel")
			runOptions = append(runOptions, createFilterString(options[0].FilterChannel))
		}
	}

	chatCh := make(chan SubscriptionMessage, channelCapacity)
	convCh := make(chan SubscriptionConversation, channelCapacity)
	walletCh := make(chan SubscriptionWalletEvent, channelCapacity)
	errorCh := make(chan error, channelCapacity)

	subs := &SubscriptionChannels{
		chat:         chatCh,
		conversation: convCh,
		wallet:       walletCh,
		error:        errorCh,
	}

	defer close(subs.chat)
	defer close(subs.conversation)
	defer close(subs.wallet)
	defer close(subs.error)

	go getNewMessages(k, subs, runOptions)
	for {
		select {
		case chatMsg := <-subs.chat:
			if handlers.ChatHandler == nil {
				continue
			}
			chatHandler := *handlers.ChatHandler
			go chatHandler(chatMsg)
		case walletMsg := <-subs.wallet:
			if handlers.WalletHandler == nil {
				continue
			}
			walletHandler := *handlers.WalletHandler
			go walletHandler(walletMsg)
		case newConv := <-subs.conversation:
			if handlers.ConversationHandler == nil {
				continue
			}
			convHandler := *handlers.ConversationHandler
			go convHandler(newConv)
		case errMsg := <-subs.error:
			if handlers.ErrorHandler == nil {
				continue
			}
			errHandler := *handlers.ErrorHandler
			go errHandler(errMsg)
		}
	}
}

// chatAPIOut sends JSON requests to the chat API and returns its response.
func chatAPIOut(k *Keybase, c ChatAPI) (ChatAPI, error) {
	jsonBytes, _ := json.Marshal(c)

	cmdOut, err := k.Exec("chat", "api", "-m", string(jsonBytes))
	if err != nil {
		return ChatAPI{}, err
	}

	var r ChatAPI
	if err := json.Unmarshal(cmdOut, &r); err != nil {
		return ChatAPI{}, err
	}
	if r.ErrorRaw != nil {
		var errorRead Error
		json.Unmarshal([]byte(*r.ErrorRaw), &errorRead)
		r.ErrorRead = &errorRead
		return r, errors.New(r.ErrorRead.Message)
	}

	return r, nil
}

// SendMessage sends a chat message
func (k *Keybase) SendMessage(method string, options SendMessageOptions) (SendResponse, error) {
	var r SendResponse

	arg := newSendMessageArg(options)
	arg.Method = method

	jsonBytes, _ := json.Marshal(arg)

	cmdOut, err := k.Exec("chat", "api", "-m", string(jsonBytes))
	if err != nil {
		return r, err
	}

	err = json.Unmarshal(cmdOut, &r)
	if err != nil {
		return r, err
	}

	return r, nil
}

// SendMessageByChannel sends a chat message to a channel
func (k *Keybase) SendMessageByChannel(channel chat1.ChatChannel, message string, a ...interface{}) (SendResponse, error) {
	var r SendResponse

	opts := SendMessageOptions{
		Channel: channel,
		Message: SendMessageBody{
			Body: fmt.Sprintf(message, a...),
		},
	}

	r, err := k.SendMessage("send", opts)
	if err != nil {
		return r, err
	}

	return r, nil
}

// SendMessageByConvID sends a chat message to a conversation id
func (k *Keybase) SendMessageByConvID(convID chat1.ConvIDStr, message string, a ...interface{}) (SendResponse, error) {
	var r SendResponse

	opts := SendMessageOptions{
		ConversationID: convID,
		Message: SendMessageBody{
			Body: fmt.Sprintf(message, a...),
		},
	}

	r, err := k.SendMessage("send", opts)
	if err != nil {
		return r, err
	}

	return r, nil
}

// SendEphemeralByChannel sends an exploding chat message to a channel
func (k *Keybase) SendEphemeralByChannel(channel chat1.ChatChannel, duration time.Duration, message string, a ...interface{}) (SendResponse, error) {
	var r SendResponse

	opts := SendMessageOptions{
		Channel: channel,
		Message: SendMessageBody{
			Body: fmt.Sprintf(message, a...),
		},
		ExplodingLifetime: &ExplodingLifetime{duration},
	}

	r, err := k.SendMessage("send", opts)
	if err != nil {
		return r, err
	}

	return r, nil
}

// SendEphemeralByConvID sends an exploding chat message to a conversation id
func (k *Keybase) SendEphemeralByConvID(convID chat1.ConvIDStr, duration time.Duration, message string, a ...interface{}) (SendResponse, error) {
	var r SendResponse

	opts := SendMessageOptions{
		ConversationID: convID,
		Message: SendMessageBody{
			Body: fmt.Sprintf(message, a...),
		},
		ExplodingLifetime: &ExplodingLifetime{duration},
	}

	r, err := k.SendMessage("send", opts)
	if err != nil {
		return r, err
	}

	return r, nil
}

// ReplyByChannel sends a reply message to a channel
func (k *Keybase) ReplyByChannel(channel chat1.ChatChannel, replyTo chat1.MessageID, message string, a ...interface{}) (SendResponse, error) {
	var r SendResponse

	opts := SendMessageOptions{
		Channel: channel,
		Message: SendMessageBody{
			Body: fmt.Sprintf(message, a...),
		},
		ReplyTo: &replyTo,
	}

	r, err := k.SendMessage("send", opts)
	if err != nil {
		return r, err
	}

	return r, nil
}

// ReplyByConvID sends a reply message to a conversation id
func (k *Keybase) ReplyByConvID(convID chat1.ConvIDStr, replyTo chat1.MessageID, message string, a ...interface{}) (SendResponse, error) {
	var r SendResponse

	opts := SendMessageOptions{
		ConversationID: convID,
		Message: SendMessageBody{
			Body: fmt.Sprintf(message, a...),
		},
		ReplyTo: &replyTo,
	}

	r, err := k.SendMessage("send", opts)
	if err != nil {
		return r, err
	}

	return r, nil
}

// EditByChannel sends an edit message to a channel
func (k *Keybase) EditByChannel(channel chat1.ChatChannel, msgID chat1.MessageID, message string, a ...interface{}) (SendResponse, error) {
	var r SendResponse

	opts := SendMessageOptions{
		Channel: channel,
		Message: SendMessageBody{
			Body: fmt.Sprintf(message, a...),
		},
		MessageID: msgID,
	}

	r, err := k.SendMessage("edit", opts)
	if err != nil {
		return r, err
	}

	return r, nil
}

// EditByConvID sends an edit message to a conversation id
func (k *Keybase) EditByConvID(convID chat1.ConvIDStr, msgID chat1.MessageID, message string, a ...interface{}) (SendResponse, error) {
	var r SendResponse

	opts := SendMessageOptions{
		ConversationID: convID,
		Message: SendMessageBody{
			Body: fmt.Sprintf(message, a...),
		},
		MessageID: msgID,
	}

	r, err := k.SendMessage("edit", opts)
	if err != nil {
		return r, err
	}

	return r, nil
}

// ReactByChannel reacts to a message in a channel
func (k *Keybase) ReactByChannel(channel chat1.ChatChannel, msgID chat1.MessageID, message string, a ...interface{}) (SendResponse, error) {
	var r SendResponse

	opts := SendMessageOptions{
		Channel: channel,
		Message: SendMessageBody{
			Body: fmt.Sprintf(message, a...),
		},
		MessageID: msgID,
	}

	r, err := k.SendMessage("reaction", opts)
	if err != nil {
		return r, err
	}

	return r, nil
}

// ReactByConvID reacts to a message in a conversation id
func (k *Keybase) ReactByConvID(convID chat1.ConvIDStr, msgID chat1.MessageID, message string, a ...interface{}) (SendResponse, error) {
	var r SendResponse

	opts := SendMessageOptions{
		ConversationID: convID,
		Message: SendMessageBody{
			Body: fmt.Sprintf(message, a...),
		},
		MessageID: msgID,
	}

	r, err := k.SendMessage("reaction", opts)
	if err != nil {
		return r, err
	}

	return r, nil
}

// DeleteByChannel reacts to a message in a channel
func (k *Keybase) DeleteByChannel(channel chat1.ChatChannel, msgID chat1.MessageID) (SendResponse, error) {
	var r SendResponse

	opts := SendMessageOptions{
		Channel:   channel,
		MessageID: msgID,
	}

	r, err := k.SendMessage("delete", opts)
	if err != nil {
		return r, err
	}

	return r, nil
}

// DeleteByConvID reacts to a message in a conversation id
func (k *Keybase) DeleteByConvID(convID chat1.ConvIDStr, msgID chat1.MessageID) (SendResponse, error) {
	var r SendResponse

	opts := SendMessageOptions{
		ConversationID: convID,
		MessageID:      msgID,
	}

	r, err := k.SendMessage("delete", opts)
	if err != nil {
		return r, err
	}

	return r, nil
}

// GetConversations returns a list of all conversations. Optionally, you can filter by unread
func (k *Keybase) GetConversations(unreadOnly bool) ([]chat1.ConvSummary, error) {
	var r Inbox

	opts := SendMessageOptions{
		UnreadOnly: unreadOnly,
	}

	arg := newSendMessageArg(opts)
	arg.Method = "list"

	jsonBytes, _ := json.Marshal(arg)

	cmdOut, err := k.Exec("chat", "api", "-m", string(jsonBytes))
	if err != nil {
		return []chat1.ConvSummary{}, err
	}

	err = json.Unmarshal(cmdOut, &r)
	if err != nil {
		return []chat1.ConvSummary{}, err
	}

	return r.Result.Convs, nil
}

// ChatList returns a list of all conversations.
// You can pass a Channel to use as a filter here, but you'll probably want to
// leave the TopicName empty.
func (k *Keybase) ChatList(opts ...chat1.ChatChannel) (ChatAPI, error) {
	m := ChatAPI{
		Params: &params{},
	}

	if len(opts) > 0 {
		m.Params.Options.Name = opts[0].Name
		m.Params.Options.Public = opts[0].Public
		m.Params.Options.MembersType = opts[0].MembersType
		m.Params.Options.TopicType = opts[0].TopicType
		m.Params.Options.TopicName = opts[0].TopicName
	}
	m.Method = "list"

	r, err := chatAPIOut(k, m)
	return r, err
}

// ReadMessage fetches the chat message with the specified message id from a conversation.
func (c Chat) ReadMessage(messageID int) (*ChatAPI, error) {
	m := ChatAPI{
		Params: &params{},
	}
	m.Params.Options = options{
		Pagination: &pagination{},
	}

	m.Method = "read"
	m.Params.Options.Channel = &c.Channel
	m.Params.Options.Pagination.Num = 1

	m.Params.Options.Pagination.Previous = getID(uint(messageID - 1))

	r, err := chatAPIOut(c.keybase, m)
	if err != nil {
		return &r, err
	}
	r.keybase = *c.keybase
	return &r, nil
}

// Read fetches chat messages from a conversation. By default, 10 messages will
// be fetched at a time. However, if count is passed, then that is the number of
// messages that will be fetched.
func (c Chat) Read(count ...int) (*ChatAPI, error) {
	m := ChatAPI{
		Params: &params{},
	}
	m.Params.Options = options{
		Pagination: &pagination{},
	}

	m.Method = "read"
	m.Params.Options.Channel = &c.Channel
	if len(count) == 0 {
		m.Params.Options.Pagination.Num = 10
	} else {
		m.Params.Options.Pagination.Num = count[0]
	}

	r, err := chatAPIOut(c.keybase, m)
	if err != nil {
		return &r, err
	}
	r.keybase = *c.keybase
	return &r, nil
}

// Next fetches the next page of chat messages that were fetched with Read. By
// default, Next will fetch the same amount of messages that were originally
// fetched with Read. However, if count is passed, then that is the number of
// messages that will be fetched.
func (c *ChatAPI) Next(count ...int) (*ChatAPI, error) {
	m := ChatAPI{
		Params: &params{},
	}
	m.Params.Options = options{
		Pagination: &pagination{},
	}

	m.Method = "read"
	m.Params.Options.Channel = &c.Result.Messages[0].Msg.Channel
	if len(count) == 0 {
		m.Params.Options.Pagination.Num = c.Result.Pagination.Num
	} else {
		m.Params.Options.Pagination.Num = count[0]
	}
	m.Params.Options.Pagination.Next = c.Result.Pagination.Next

	result, err := chatAPIOut(&c.keybase, m)
	if err != nil {
		return &result, err
	}
	k := c.keybase
	*c = result
	c.keybase = k
	return c, nil
}

// Previous fetches the previous page of chat messages that were fetched with Read.
// By default, Previous will fetch the same amount of messages that were
// originally fetched with Read. However, if count is passed, then that is the
// number of messages that will be fetched.
func (c *ChatAPI) Previous(count ...int) (*ChatAPI, error) {
	m := ChatAPI{
		Params: &params{},
	}
	m.Params.Options = options{
		Pagination: &pagination{},
	}

	m.Method = "read"
	m.Params.Options.Channel = &c.Result.Messages[0].Msg.Channel
	if len(count) == 0 {
		m.Params.Options.Pagination.Num = c.Result.Pagination.Num
	} else {
		m.Params.Options.Pagination.Num = count[0]
	}
	m.Params.Options.Pagination.Previous = c.Result.Pagination.Previous

	result, err := chatAPIOut(&c.keybase, m)
	if err != nil {
		return &result, err
	}
	k := c.keybase
	*c = result
	c.keybase = k
	return c, nil
}

// Upload attaches a file to a conversation
// The filepath must be an absolute path
func (c Chat) Upload(title string, filepath string) (ChatAPI, error) {
	m := ChatAPI{
		Params: &params{},
	}
	m.Method = "attach"
	m.Params.Options.Channel = &c.Channel
	m.Params.Options.Filename = filepath
	m.Params.Options.Title = title

	r, err := chatAPIOut(c.keybase, m)
	if err != nil {
		return r, err
	}
	return r, nil
}

// Download downloads a file from a conversation
func (c Chat) Download(messageID int, filepath string) (ChatAPI, error) {
	m := ChatAPI{
		Params: &params{},
	}
	m.Method = "download"
	m.Params.Options.Channel = &c.Channel
	m.Params.Options.Output = filepath
	m.Params.Options.MessageID = messageID

	r, err := chatAPIOut(c.keybase, m)
	if err != nil {
		return r, err
	}
	return r, nil
}

// LoadFlip returns the results of a flip
// If the flip is still in progress, this can be expected to change if called again
func (c Chat) LoadFlip(messageID int, conversationID string, flipConversationID string, gameID string) (ChatAPI, error) {
	m := ChatAPI{
		Params: &params{},
	}
	m.Method = "loadflip"
	m.Params.Options.Channel = &c.Channel
	m.Params.Options.MsgID = messageID
	m.Params.Options.ConversationID = conversationID
	m.Params.Options.FlipConversationID = flipConversationID
	m.Params.Options.GameID = gameID

	r, err := chatAPIOut(c.keybase, m)
	if err != nil {
		return r, err
	}
	return r, nil
}

// Pin pins a message to a channel
func (c Chat) Pin(messageID int) (ChatAPI, error) {
	m := ChatAPI{
		Params: &params{},
	}
	m.Method = "pin"
	m.Params.Options.Channel = &c.Channel
	m.Params.Options.MessageID = messageID

	r, err := chatAPIOut(c.keybase, m)
	if err != nil {
		return r, err
	}
	return r, nil
}

// Unpin clears any pinned messages from a channel
func (c Chat) Unpin() (ChatAPI, error) {
	m := ChatAPI{
		Params: &params{},
	}
	m.Method = "unpin"
	m.Params.Options.Channel = &c.Channel

	r, err := chatAPIOut(c.keybase, m)
	if err != nil {
		return r, err
	}
	return r, nil
}

// Mark marks a conversation as read up to a specified message
func (c Chat) Mark(messageID int) (ChatAPI, error) {
	m := ChatAPI{
		Params: &params{},
	}
	m.Method = "mark"
	m.Params.Options.Channel = &c.Channel
	m.Params.Options.MessageID = messageID

	r, err := chatAPIOut(c.keybase, m)
	if err != nil {
		return r, err
	}
	return r, nil
}

// ClearCommands clears bot advertisements
func (k *Keybase) ClearCommands() (ChatAPI, error) {
	m := ChatAPI{}
	m.Method = "clearcommands"

	r, err := chatAPIOut(k, m)
	if err != nil {
		return r, err
	}
	return r, nil
}

// AdvertiseCommands sets up bot command advertisements
// This method allows you to set up multiple different types of advertisements at once.
// Use this method if you have commands whose visibility differs from each other.
func (k *Keybase) AdvertiseCommands(advertisements []BotAdvertisement) (ChatAPI, error) {
	m := ChatAPI{
		Params: &params{},
	}
	m.Method = "advertisecommands"
	m.Params.Options.BotAdvertisements = advertisements

	r, err := chatAPIOut(k, m)
	if err != nil {
		return r, err
	}
	return r, nil
}

// AdvertiseCommand sets up bot command advertisements
// This method allows you to set up one type of advertisement.
// Use this method if you have commands whose visibility should all be the same.
func (k *Keybase) AdvertiseCommand(advertisement BotAdvertisement) (ChatAPI, error) {
	return k.AdvertiseCommands([]BotAdvertisement{
		advertisement,
	})
}
