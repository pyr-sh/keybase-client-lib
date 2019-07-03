package keybase

import (
	"bufio"
	"encoding/json"
	"os/exec"
	"time"
)

// ChatIn holds information about a message received by the `keybase chat api-listen` command
type ChatIn struct {
	Type   string    `json:"type"`
	Source string    `json:"source"`
	Msg    chatInMsg `json:"msg"`
}
type chatInSender struct {
	UID        string `json:"uid"`
	Username   string `json:"username"`
	DeviceID   string `json:"device_id"`
	DeviceName string `json:"device_name"`
}
type chatInAddedtoteam struct {
	Team    string   `json:"team"`
	Adder   string   `json:"adder"`
	Addee   string   `json:"addee"`
	Owners  []string `json:"owners"`
	Admins  []string `json:"admins"`
	Writers []string `json:"writers"`
	Readers []string `json:"readers"`
}
type chatInBulkaddtoconv struct {
	Usernames []string `json:"usernames"`
}
type chatInCommits struct {
	CommitHash  string `json:"commitHash"`
	Message     string `json:"message"`
	AuthorName  string `json:"authorName"`
	AuthorEmail string `json:"authorEmail"`
	Ctime       int    `json:"ctime"`
}
type chatInRefs struct {
	RefName              string          `json:"refName"`
	Commits              []chatInCommits `json:"commits"`
	MoreCommitsAvailable bool            `json:"moreCommitsAvailable"`
	IsDelete             bool            `json:"isDelete"`
}
type chatInGitpush struct {
	Team             string       `json:"team"`
	Pusher           string       `json:"pusher"`
	RepoName         string       `json:"repoName"`
	RepoID           string       `json:"repoID"`
	Refs             []chatInRefs `json:"refs"`
	PushType         int          `json:"pushType"`
	PreviousRepoName string       `json:"previousRepoName"`
}
type chatInSystem struct {
	SystemType    int                 `json:"systemType"`
	Addedtoteam   chatInAddedtoteam   `json:"addedtoteam"`
	Bulkaddtoconv chatInBulkaddtoconv `json:"bulkaddtoconv"`
	Gitpush       chatInGitpush       `json:"gitpush"`
}
type chatInResult struct {
	ResultTyp int    `json:"resultTyp"`
	Sent      string `json:"sent"`
}
type chatInPayments struct {
	Username    string       `json:"username"`
	PaymentText string       `json:"paymentText"`
	Result      chatInResult `json:"result"`
}
type chatInUserMentions struct {
	Text string `json:"text"`
	UID  string `json:"uid"`
}
type chatInTeamMentions struct {
	Name    string `json:"name"`
	Channel string `json:"channel"`
}
type chatInReaction struct {
	M int    `json:"m"`
	B string `json:"b"`
}
type chatInDelete struct {
	MessageIDs []int `json:"messageIDs"`
}
type chatInEdit struct {
	MessageID    int                  `json:"messageID"`
	Body         string               `json:"body"`
	Payments     []chatInPayments     `json:"payments"`
	UserMentions []chatInUserMentions `json:"userMentions"`
	TeamMentions []chatInTeamMentions `json:"teamMentions"`
}
type chatInText struct {
	Body         string               `json:"body"`
	Payments     []chatInPayments     `json:"payments"`
	UserMentions []chatInUserMentions `json:"userMentions"`
	TeamMentions []chatInTeamMentions `json:"teamMentions"`
}
type chatInContent struct {
	Type     string         `json:"type"`
	Delete   chatInDelete   `json:"delete"`
	Edit     chatInEdit     `json:"edit"`
	Reaction chatInReaction `json:"reaction"`
	System   chatInSystem   `json:"system"`
	Text     chatInText     `json:"text"`
}
type chatInMsg struct {
	ID                 int           `json:"id"`
	Channel            Channel       `json:"channel"`
	Sender             chatInSender  `json:"sender"`
	SentAt             int           `json:"sent_at"`
	SentAtMs           int64         `json:"sent_at_ms"`
	Content            chatInContent `json:"content"`
	Unread             bool          `json:"unread"`
	AtMentionUsernames []string      `json:"at_mention_usernames"`
	IsEphemeral        bool          `json:"is_ephemeral"`
	Etime              int64         `json:"etime"`
	HasPairwiseMacs    bool          `json:"has_pairwise_macs"`
	ChannelMention     string        `json:"channel_mention"`
}

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
func getNewMessages(k *Keybase, c chan<- ChatIn, execOptions []string) {
	execString := []string{"chat", "api-listen"}
	if len(execOptions) > 0 {
		execString = append(execString, execOptions...)
	}
	for {
		execCmd := exec.Command(k.Path, execString...)
		stdOut, _ := execCmd.StdoutPipe()
		execCmd.Start()
		scanner := bufio.NewScanner(stdOut)
		go func(scanner *bufio.Scanner, c chan<- ChatIn) {
			var jsonData ChatIn
			for scanner.Scan() {
				json.Unmarshal([]byte(scanner.Text()), &jsonData)
				c <- jsonData
			}
		}(scanner, c)
		execCmd.Wait()
	}
}

// Run runs `keybase chat api-listen`, and passes incoming messages to the message handler func
func (k *Keybase) Run(handler func(ChatIn), options ...RunOptions) {
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
	c := make(chan ChatIn, 50)
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
func heartbeat(c chan<- ChatIn, freq time.Duration) {
	m := ChatIn{
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
