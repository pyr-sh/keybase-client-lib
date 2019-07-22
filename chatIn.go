package keybase

import (
	"bufio"
	"encoding/json"
	"os/exec"
	"time"
)

// RunOptions holds a set of options to be passed to Run
type RunOptions struct {
	Heartbeat      int64     // Send a heartbeat through the channel every X minutes (0 = off)
	Local          bool      // Subscribe to local messages
	HideExploding  bool      // Ignore exploding messages
	Dev            bool      // Subscribe to dev channel messages
	Wallet         bool      // Subscribe to wallet events
	FilterChannel  Channel   // Only subscribe to messages from specified channel
	FilterChannels []Channel // Only subscribe to messages from specified channels
}

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
