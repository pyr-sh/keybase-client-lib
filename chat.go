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
	"samhofi.us/x/keybase/types/stellar1"
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
						subs.chat <- *notification.Msg
					}
				case "chat_conv":
					var notification chat1.ConvNotification
					if err := json.Unmarshal([]byte(t), &notification); err != nil {
						subs.error <- err
						break
					}
					if notification.Conv != nil {
						subs.conversation <- *notification.Conv
					}
				case "wallet":
					var holder paymentHolder
					if err := json.Unmarshal([]byte(t), &holder); err != nil {
						subs.error <- err
						break
					}
					subs.wallet <- holder.Payment
				default:
					continue
				}
			}
		}(scanner, subs)
		execCmd.Wait()
	}
}

// Run runs `keybase chat api-listen`, and passes incoming messages to the message handler func
func (k *Keybase) Run(handlers Handlers, options *RunOptions) {
	var channelCapacity = 100

	runOptions := make([]string, 0)
	if handlers.WalletHandler != nil {
		runOptions = append(runOptions, "--wallet")
	}
	if handlers.ConversationHandler != nil {
		runOptions = append(runOptions, "--convs")
	}

	if options != nil {
		if options.Capacity > 0 {
			channelCapacity = options.Capacity
		}
		if options.Local {
			runOptions = append(runOptions, "--local")
		}
		if options.HideExploding {
			runOptions = append(runOptions, "--hide-exploding")
		}
		if options.Dev {
			runOptions = append(runOptions, "--dev")
		}
		if len(options.FilterChannels) > 0 {
			runOptions = append(runOptions, "--filter-channels")
			runOptions = append(runOptions, createFiltersString(options.FilterChannels))

		}
		if options.FilterChannel.Name != "" {
			runOptions = append(runOptions, "--filter-channel")
			runOptions = append(runOptions, createFilterString(options.FilterChannel))
		}
	}

	chatCh := make(chan chat1.MsgSummary, channelCapacity)
	convCh := make(chan chat1.ConvSummary, channelCapacity)
	walletCh := make(chan stellar1.PaymentDetailsLocal, channelCapacity)
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
func (k *Keybase) SendMessage(method string, options SendMessageOptions) (chat1.SendRes, error) {
	type res struct {
		Result chat1.SendRes `json:"result"`
		Error  *Error        `json:"error,omitempty"`
	}

	var r res

	arg := newSendMessageArg(options)
	arg.Method = method

	jsonBytes, _ := json.Marshal(arg)

	cmdOut, err := k.Exec("chat", "api", "-m", string(jsonBytes))
	if err != nil {
		return r.Result, err
	}

	err = json.Unmarshal(cmdOut, &r)
	if err != nil {
		return r.Result, err
	}

	if r.Error != nil {
		return r.Result, fmt.Errorf("%v", r.Error.Message)
	}

	return r.Result, nil
}

// SendMessageByChannel sends a chat message to a channel
func (k *Keybase) SendMessageByChannel(channel chat1.ChatChannel, message string, a ...interface{}) (chat1.SendRes, error) {
	opts := SendMessageOptions{
		Channel: channel,
		Message: SendMessageBody{
			Body: fmt.Sprintf(message, a...),
		},
	}

	return k.SendMessage("send", opts)
}

// SendMessageByConvID sends a chat message to a conversation id
func (k *Keybase) SendMessageByConvID(convID chat1.ConvIDStr, message string, a ...interface{}) (chat1.SendRes, error) {
	opts := SendMessageOptions{
		ConversationID: convID,
		Message: SendMessageBody{
			Body: fmt.Sprintf(message, a...),
		},
	}

	return k.SendMessage("send", opts)
}

// SendEphemeralByChannel sends an exploding chat message to a channel
func (k *Keybase) SendEphemeralByChannel(channel chat1.ChatChannel, duration time.Duration, message string, a ...interface{}) (chat1.SendRes, error) {
	opts := SendMessageOptions{
		Channel: channel,
		Message: SendMessageBody{
			Body: fmt.Sprintf(message, a...),
		},
		ExplodingLifetime: &ExplodingLifetime{duration},
	}

	return k.SendMessage("send", opts)
}

// SendEphemeralByConvID sends an exploding chat message to a conversation id
func (k *Keybase) SendEphemeralByConvID(convID chat1.ConvIDStr, duration time.Duration, message string, a ...interface{}) (chat1.SendRes, error) {
	opts := SendMessageOptions{
		ConversationID: convID,
		Message: SendMessageBody{
			Body: fmt.Sprintf(message, a...),
		},
		ExplodingLifetime: &ExplodingLifetime{duration},
	}

	return k.SendMessage("send", opts)
}

// ReplyByChannel sends a reply message to a channel
func (k *Keybase) ReplyByChannel(channel chat1.ChatChannel, replyTo chat1.MessageID, message string, a ...interface{}) (chat1.SendRes, error) {
	opts := SendMessageOptions{
		Channel: channel,
		Message: SendMessageBody{
			Body: fmt.Sprintf(message, a...),
		},
		ReplyTo: &replyTo,
	}

	return k.SendMessage("send", opts)
}

// ReplyByConvID sends a reply message to a conversation id
func (k *Keybase) ReplyByConvID(convID chat1.ConvIDStr, replyTo chat1.MessageID, message string, a ...interface{}) (chat1.SendRes, error) {
	opts := SendMessageOptions{
		ConversationID: convID,
		Message: SendMessageBody{
			Body: fmt.Sprintf(message, a...),
		},
		ReplyTo: &replyTo,
	}

	return k.SendMessage("send", opts)
}

// EditByChannel sends an edit message to a channel
func (k *Keybase) EditByChannel(channel chat1.ChatChannel, msgID chat1.MessageID, message string, a ...interface{}) (chat1.SendRes, error) {
	opts := SendMessageOptions{
		Channel: channel,
		Message: SendMessageBody{
			Body: fmt.Sprintf(message, a...),
		},
		MessageID: msgID,
	}

	return k.SendMessage("edit", opts)
}

// EditByConvID sends an edit message to a conversation id
func (k *Keybase) EditByConvID(convID chat1.ConvIDStr, msgID chat1.MessageID, message string, a ...interface{}) (chat1.SendRes, error) {
	opts := SendMessageOptions{
		ConversationID: convID,
		Message: SendMessageBody{
			Body: fmt.Sprintf(message, a...),
		},
		MessageID: msgID,
	}

	return k.SendMessage("edit", opts)
}

// ReactByChannel reacts to a message in a channel
func (k *Keybase) ReactByChannel(channel chat1.ChatChannel, msgID chat1.MessageID, message string, a ...interface{}) (chat1.SendRes, error) {
	opts := SendMessageOptions{
		Channel: channel,
		Message: SendMessageBody{
			Body: fmt.Sprintf(message, a...),
		},
		MessageID: msgID,
	}

	return k.SendMessage("reaction", opts)
}

// ReactByConvID reacts to a message in a conversation id
func (k *Keybase) ReactByConvID(convID chat1.ConvIDStr, msgID chat1.MessageID, message string, a ...interface{}) (chat1.SendRes, error) {
	opts := SendMessageOptions{
		ConversationID: convID,
		Message: SendMessageBody{
			Body: fmt.Sprintf(message, a...),
		},
		MessageID: msgID,
	}

	return k.SendMessage("reaction", opts)
}

// DeleteByChannel reacts to a message in a channel
func (k *Keybase) DeleteByChannel(channel chat1.ChatChannel, msgID chat1.MessageID) (chat1.SendRes, error) {
	opts := SendMessageOptions{
		Channel:   channel,
		MessageID: msgID,
	}

	return k.SendMessage("delete", opts)
}

// DeleteByConvID reacts to a message in a conversation id
func (k *Keybase) DeleteByConvID(convID chat1.ConvIDStr, msgID chat1.MessageID) (chat1.SendRes, error) {
	opts := SendMessageOptions{
		ConversationID: convID,
		MessageID:      msgID,
	}

	return k.SendMessage("delete", opts)
}

// GetConversations returns a list of all conversations.
func (k *Keybase) GetConversations(unreadOnly bool) ([]chat1.ConvSummary, error) {
	type res struct {
		Result []chat1.ConvSummary `json:"result"`
		Error  *Error              `json:"error,omitempty"`
	}

	var r res

	opts := SendMessageOptions{
		UnreadOnly: unreadOnly,
	}

	arg := newSendMessageArg(opts)
	arg.Method = "list"

	jsonBytes, _ := json.Marshal(arg)

	cmdOut, err := k.Exec("chat", "api", "-m", string(jsonBytes))
	if err != nil {
		return r.Result, err
	}

	err = json.Unmarshal(cmdOut, &r)
	if err != nil {
		return r.Result, err
	}

	if r.Error != nil {
		return r.Result, fmt.Errorf("%v", r.Error.Message)
	}

	return r.Result, nil
}

// Read fetches chat messages
func (k *Keybase) Read(options ReadMessageOptions) (chat1.Thread, error) {
	type res struct {
		Result chat1.Thread `json:"result"`
		Error  *Error       `json:"error"`
	}
	var r res

	arg := newReadMessageArg(options)

	jsonBytes, _ := json.Marshal(arg)

	cmdOut, err := k.Exec("chat", "api", "-m", string(jsonBytes))
	if err != nil {
		return r.Result, err
	}

	err = json.Unmarshal(cmdOut, &r)
	if err != nil {
		return r.Result, err
	}

	if r.Error != nil {
		return r.Result, fmt.Errorf("%v", r.Error.Message)
	}

	return r.Result, nil
}

// ReadChannel fetches chat messages for a channel
func (k *Keybase) ReadChannel(channel chat1.ChatChannel) (chat1.Thread, error) {
	opts := ReadMessageOptions{
		Channel: channel,
	}
	return k.Read(opts)
}

// ReadChannelNext fetches the next page of messages for a chat channel.
func (k *Keybase) ReadChannelNext(channel chat1.ChatChannel, next []byte, num int) (chat1.Thread, error) {
	page := chat1.Pagination{
		Next: next,
		Num:  num,
	}

	opts := ReadMessageOptions{
		Channel:    channel,
		Pagination: &page,
	}
	return k.Read(opts)
}

// ReadChannelPrevious fetches the previous page of messages for a chat channel
func (k *Keybase) ReadChannelPrevious(channel chat1.ChatChannel, previous []byte, num int) (chat1.Thread, error) {
	page := chat1.Pagination{
		Previous: previous,
		Num:      num,
	}

	opts := ReadMessageOptions{
		Channel:    channel,
		Pagination: &page,
	}
	return k.Read(opts)
}

// ReadConversation fetches chat messages for a conversation
func (k *Keybase) ReadConversation(conv chat1.ConvIDStr) (chat1.Thread, error) {
	opts := ReadMessageOptions{
		ConversationID: conv,
	}
	return k.Read(opts)
}

// ReadConversationNext fetches the next page of messages for a conversation.
func (k *Keybase) ReadConversationNext(conv chat1.ConvIDStr, next []byte, num int) (chat1.Thread, error) {
	page := chat1.Pagination{
		Next: next,
		Num:  num,
	}

	opts := ReadMessageOptions{
		ConversationID: conv,
		Pagination:     &page,
	}
	return k.Read(opts)
}

// ReadConversationPrevious fetches the previous page of messages for a chat channel
func (k *Keybase) ReadConversationPrevious(conv chat1.ConvIDStr, previous []byte, num int) (chat1.Thread, error) {
	page := chat1.Pagination{
		Previous: previous,
		Num:      num,
	}

	opts := ReadMessageOptions{
		ConversationID: conv,
		Pagination:     &page,
	}
	return k.Read(opts)
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

// UploadToChannel attaches a file to a channel
// The filename must be an absolute path
func (k *Keybase) UploadToChannel(channel chat1.ChatChannel, title string, filename string) (chat1.SendRes, error) {
	opts := SendMessageOptions{
		Channel:  channel,
		Title:    title,
		Filename: filename,
	}

	return k.SendMessage("attach", opts)
}

// UploadToConversation attaches a file to a conversation
// The filename must be an absolute path
func (k *Keybase) UploadToConversation(conv chat1.ConvIDStr, title string, filename string) (chat1.SendRes, error) {
	opts := SendMessageOptions{
		ConversationID: conv,
		Title:          title,
		Filename:       filename,
	}

	return k.SendMessage("attach", opts)
}

// Download downloads a file
func (k *Keybase) Download(options DownloadOptions) error {
	type res struct {
		Error *Error `json:"error"`
	}
	var r res

	arg := newDownloadArg(options)

	jsonBytes, _ := json.Marshal(arg)

	cmdOut, err := k.Exec("chat", "api", "-m", string(jsonBytes))
	if err != nil {
		return err
	}

	err = json.Unmarshal(cmdOut, &r)
	if err != nil {
		return err
	}

	if r.Error != nil {
		return fmt.Errorf("%v", r.Error.Message)
	}

	return nil
}

// DownloadFromChannel downloads a file from a channel
func (k *Keybase) DownloadFromChannel(channel chat1.ChatChannel, msgID chat1.MessageID, output string) error {
	opts := DownloadOptions{
		Channel:   channel,
		MessageID: msgID,
		Output:    output,
	}
	return k.Download(opts)
}

// DownloadFromConversation downloads a file from a conversation
func (k *Keybase) DownloadFromConversation(conv chat1.ConvIDStr, msgID chat1.MessageID, output string) error {
	opts := DownloadOptions{
		ConversationID: conv,
		MessageID:      msgID,
		Output:         output,
	}
	return k.Download(opts)
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

// AdvertiseCommands sends bot command advertisements.
// Valid values for the `Typ` field in chat1.AdvertiseCommandAPIParam are
// "public", "teamconvs", and "teammembers"
func (k *Keybase) AdvertiseCommands(options AdvertiseCommandsOptions) error {
	type res struct {
		Error *Error `json:"error,omitempty"`
	}

	var r res

	arg := newAdvertiseCommandsArg(options)

	jsonBytes, _ := json.Marshal(arg)

	cmdOut, err := k.Exec("chat", "api", "-m", string(jsonBytes))
	if err != nil {
		return err
	}

	err = json.Unmarshal(cmdOut, &r)
	if err != nil {
		return err
	}

	if r.Error != nil {
		return fmt.Errorf("%v", r.Error.Message)
	}

	return nil
}

// ClearCommands clears bot advertisements
func (k *Keybase) ClearCommands() error {
	type res struct {
		Error *Error `json:"error,omitempty"`
	}

	var r res

	cmdOut, err := k.Exec("chat", "api", "-m", `{"method": "clearcommands"}`)
	if err != nil {
		return err
	}

	err = json.Unmarshal(cmdOut, &r)
	if err != nil {
		return err
	}

	if r.Error != nil {
		return fmt.Errorf("%v", r.Error.Message)
	}

	return nil
}
