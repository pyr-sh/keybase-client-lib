package keybase

import (
	"bufio"
	"encoding/json"
	"os/exec"
	"time"
)

// Creates a string of a json-encoded channel to pass to keybase chat api-listen --filter-channel
func createFilterString(channel Channel) string {
	if channel.Name == "" {
		return ""
	}
	jsonBytes, _ := json.Marshal(channel)
	return string(jsonBytes)
}

// Creates a string of json-encoded channels to pass to keybase chat api-listen --filter-channels
func createFiltersString(channels []Channel) string {
	if len(channels) == 0 {
		return ""
	}
	jsonBytes, _ := json.Marshal(channels)
	return string(jsonBytes)
}

// Run `keybase chat api-listen` to get new messages coming into keybase and send them into the channel
func getNewMessages(k *Keybase, c chan<- ChatAPI, execOptions []string) {
	execString := []string{"chat", "api-listen"}
	if len(execOptions) > 0 {
		execString = append(execString, execOptions...)
	}
	for {
		execCmd := exec.Command(k.Path, execString...)
		stdOut, _ := execCmd.StdoutPipe()
		execCmd.Start()
		scanner := bufio.NewScanner(stdOut)
		go func(scanner *bufio.Scanner, c chan<- ChatAPI) {
			var jsonData ChatAPI
			for scanner.Scan() {
				json.Unmarshal([]byte(scanner.Text()), &jsonData)
				c <- jsonData
			}
		}(scanner, c)
		execCmd.Wait()
	}
}

// Run runs `keybase chat api-listen`, and passes incoming messages to the message handler func
func (k *Keybase) Run(handler func(ChatAPI), options ...RunOptions) {
	var heartbeatFreq int64
	runOptions := make([]string, 0)
	if len(options) > 0 {
		if options[0].Heartbeat > 0 {
			heartbeatFreq = options[0].Heartbeat
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
	c := make(chan ChatAPI, 50)
	defer close(c)
	if heartbeatFreq > 0 {
		go heartbeat(c, time.Duration(heartbeatFreq)*time.Minute)
	}
	go getNewMessages(k, c, runOptions)
	for {
		go handler(<-c)
	}
}

// heartbeat sends a message through the channel with a message type of `heartbeat`
func heartbeat(c chan<- ChatAPI, freq time.Duration) {
	m := ChatAPI{
		Type: "heartbeat",
	}
	count := 0
	for {
		time.Sleep(freq)
		m.Msg.ID = count
		c <- m
		count++
	}
}

// Read fetches chat messages from a conversation. By default, 10 messages will
// be fetched at a time. However, if count is passed, then that is the number of
// messages that will be fetched.
func (c Chat) Read(count ...int) (*ChatAPI, error) {
	m := ChatAPI{}
	m.Method = "read"
	m.Params.Options.Channel = c.Channel
	if len(count) == 0 {
		m.Params.Options.Pagination.Num = 10
	} else {
		m.Params.Options.Pagination.Num = count[0]
	}

	r, err := chatAPIOut(c.keybase.Path, m)
	if err != nil {
		return &ChatAPI{}, err
	}
	r.keybase = *c.keybase
	return &r, nil
}

// Next fetches the next page of chat messages that were fetched with Read. By
// default, Next will fetch the same amount of messages that were originally
// fetched with Read. However, if count is passed, then that is the number of
// messages that will be fetched.
func (c *ChatAPI) Next(count ...int) (*ChatAPI, error) {
	m := ChatAPI{}
	m.Method = "read"
	m.Params.Options.Channel = c.Result.Messages[0].Msg.Channel
	if len(count) == 0 {
		m.Params.Options.Pagination.Num = c.Result.Pagination.Num
	} else {
		m.Params.Options.Pagination.Num = count[0]
	}
	m.Params.Options.Pagination.Next = c.Result.Pagination.Next

	result, err := chatAPIOut(c.keybase.Path, m)
	if err != nil {
		return &ChatAPI{}, err
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
	m := ChatAPI{}
	m.Method = "read"
	m.Params.Options.Channel = c.Result.Messages[0].Msg.Channel
	if len(count) == 0 {
		m.Params.Options.Pagination.Num = c.Result.Pagination.Num
	} else {
		m.Params.Options.Pagination.Num = count[0]
	}
	m.Params.Options.Pagination.Previous = c.Result.Pagination.Previous

	result, err := chatAPIOut(c.keybase.Path, m)
	if err != nil {
		return &ChatAPI{}, err
	}
	k := c.keybase
	*c = result
	c.keybase = k
	return c, nil
}
