package keybase

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

// ChatAPI holds information about a message received by the `keybase chat api-listen` command
type ChatAPI struct {
	Type          string         `json:"type,omitempty"`
	Source        string         `json:"source,omitempty"`
	Msg           msg            `json:"msg,omitempty"`
	Method        string         `json:"method,omitempty"`
	Params        params         `json:"params,omitempty"`
	Message       string         `json:"message,omitempty"`
	ID            int            `json:"id,omitempty"`
	Ratelimits    []rateLimits   `json:"ratelimits,omitempty"`
	Conversations []conversation `json:"conversations,omitempty"`
	Offline       bool           `json:"offline,omitempty"`
	Result        result         `json:"result,omitempty"`
	keybase       Keybase        // Some methods will need this, so I'm passing it but keeping it unexported
}
type sender struct {
	UID        string `json:"uid"`
	Username   string `json:"username"`
	DeviceID   string `json:"device_id"`
	DeviceName string `json:"device_name"`
}
type addedtoteam struct {
	Team    string   `json:"team"`
	Adder   string   `json:"adder"`
	Addee   string   `json:"addee"`
	Owners  []string `json:"owners"`
	Admins  []string `json:"admins"`
	Writers []string `json:"writers"`
	Readers []string `json:"readers"`
}
type bulkaddtoconv struct {
	Usernames []string `json:"usernames"`
}
type commits struct {
	CommitHash  string `json:"commitHash"`
	Message     string `json:"message"`
	AuthorName  string `json:"authorName"`
	AuthorEmail string `json:"authorEmail"`
	Ctime       int    `json:"ctime"`
}
type refs struct {
	RefName              string    `json:"refName"`
	Commits              []commits `json:"commits"`
	MoreCommitsAvailable bool      `json:"moreCommitsAvailable"`
	IsDelete             bool      `json:"isDelete"`
}
type gitpush struct {
	Team             string `json:"team"`
	Pusher           string `json:"pusher"`
	RepoName         string `json:"repoName"`
	RepoID           string `json:"repoID"`
	Refs             []refs `json:"refs"`
	PushType         int    `json:"pushType"`
	PreviousRepoName string `json:"previousRepoName"`
}
type system struct {
	SystemType    int           `json:"systemType"`
	Addedtoteam   addedtoteam   `json:"addedtoteam"`
	Bulkaddtoconv bulkaddtoconv `json:"bulkaddtoconv"`
	Gitpush       gitpush       `json:"gitpush"`
}
type paymentsResult struct {
	ResultTyp int    `json:"resultTyp"`
	Sent      string `json:"sent"`
}
type payments struct {
	Username    string         `json:"username"`
	PaymentText string         `json:"paymentText"`
	Result      paymentsResult `json:"result"`
}
type userMentions struct {
	Text string `json:"text"`
	UID  string `json:"uid"`
}
type teamMentions struct {
	Name    string `json:"name"`
	Channel string `json:"channel"`
}
type reaction struct {
	M int    `json:"m"`
	B string `json:"b"`
}
type delete struct {
	MessageIDs []int `json:"messageIDs"`
}
type edit struct {
	MessageID    int            `json:"messageID"`
	Body         string         `json:"body"`
	Payments     []payments     `json:"payments"`
	UserMentions []userMentions `json:"userMentions"`
	TeamMentions []teamMentions `json:"teamMentions"`
}
type text struct {
	Body         string         `json:"body"`
	Payments     []payments     `json:"payments"`
	UserMentions []userMentions `json:"userMentions"`
	TeamMentions []teamMentions `json:"teamMentions"`
}
type content struct {
	Type     string   `json:"type"`
	Delete   delete   `json:"delete"`
	Edit     edit     `json:"edit"`
	Reaction reaction `json:"reaction"`
	System   system   `json:"system"`
	Text     text     `json:"text"`
}
type msg struct {
	ID                 int      `json:"id"`
	Channel            Channel  `json:"channel"`
	Sender             sender   `json:"sender"`
	SentAt             int      `json:"sent_at"`
	SentAtMs           int64    `json:"sent_at_ms"`
	Content            content  `json:"content"`
	Unread             bool     `json:"unread"`
	AtMentionUsernames []string `json:"at_mention_usernames"`
	IsEphemeral        bool     `json:"is_ephemeral"`
	Etime              int64    `json:"etime"`
	HasPairwiseMacs    bool     `json:"has_pairwise_macs"`
	ChannelMention     string   `json:"channel_mention"`
}

// Channel holds information about a conversation
type Channel struct {
	Name        string `json:"name"`
	Public      bool   `json:"public,omitempty"`
	MembersType string `json:"members_type,omitempty"`
	TopicType   string `json:"topic_type,omitempty"`
	TopicName   string `json:"topic_name,omitempty"`
}
type message struct {
	Body string `json:"body"`
}
type options struct {
	Channel    Channel    `json:"channel"`
	MessageID  int        `json:"message_id"`
	Message    message    `json:"message"`
	Pagination pagination `json:"pagination"`
}
type params struct {
	Options options `json:"options"`
}
type pagination struct {
	Next           string `json:"next"`
	Previous       string `json:"previous"`
	Num            int    `json:"num"`
	Last           bool   `json:"last,omitempty"`
	ForceFirstPage bool   `json:"forceFirstPage,omitempty"`
}
type result struct {
	Messages   []messages   `json:"messages,omitempty"`
	Pagination pagination   `json:"pagination"`
	Message    string       `json:"message"`
	ID         int          `json:"id"`
	Ratelimits []rateLimits `json:"ratelimits"`
}
type messages struct {
	Msg msg `json:"msg,omitempty"`
}
type rateLimits struct {
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

// Keybase holds basic information about the local Keybase executable
type Keybase struct {
	Path     string
	Username string
	LoggedIn bool
	Version  string
}

// Chat holds basic information about a specific conversation
type Chat struct {
	keybase *Keybase
	Channel Channel
}

type chat interface {
	Send(message ...string) (ChatAPI, error)
	Edit(messageID int, message ...string) (ChatAPI, error)
	React(messageID int, reaction string) (ChatAPI, error)
	Delete(messageID int) (ChatAPI, error)
}

type chatAPI interface {
	Next(count ...int) (*ChatAPI, error)
	Previous(count ...int) (*ChatAPI, error)
}

type keybase interface {
	NewChat(channel Channel) Chat
	Run(handler func(ChatAPI), options ...RunOptions)
	ChatList() ([]conversation, error)
	loggedIn() bool
	username() string
	version() string
}

type status struct {
	Username string `json:"Username"`
	LoggedIn bool   `json:"LoggedIn"`
}
