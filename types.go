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
	Type         string        `json:"type,omitempty"`
	Source       string        `json:"source,omitempty"`
	Msg          *msg          `json:"msg,omitempty"`
	Method       string        `json:"method,omitempty"`
	Params       *params       `json:"params,omitempty"`
	Message      string        `json:"message,omitempty"`
	ID           int           `json:"id,omitempty"`
	Ratelimits   []rateLimits  `json:"ratelimits,omitempty"`
	Notification *notification `json:"notification"`
	Result       *result       `json:"result,omitempty"`
	Pagination   *pagination   `json:"pagination"`
	Error        *Error        `json:"error"`
	keybase      Keybase       // Some methods will need this, so I'm passing it but keeping it unexported
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
	Type        string      `json:"type"`
	Delete      delete      `json:"delete"`
	Edit        edit        `json:"edit"`
	Reaction    reaction    `json:"reaction"`
	System      system      `json:"system"`
	Text        text        `json:"text"`
	SendPayment SendPayment `json:"send_payment"`
}
type msg struct {
	ID                 int      `json:"id"`
	ConversationID     string   `json:"conversation_id"`
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
type summary struct {
	ID                  string      `json:"id"`
	TxID                string      `json:"txID"`
	Time                int64       `json:"time"`
	StatusSimplified    int         `json:"statusSimplified"`
	StatusDescription   string      `json:"statusDescription"`
	StatusDetail        string      `json:"statusDetail"`
	ShowCancel          bool        `json:"showCancel"`
	AmountDescription   string      `json:"amountDescription"`
	Delta               int         `json:"delta"`
	Worth               string      `json:"worth"`
	WorthAtSendTime     string      `json:"worthAtSendTime"`
	IssuerDescription   string      `json:"issuerDescription"`
	FromType            int         `json:"fromType"`
	ToType              int         `json:"toType"`
	AssetCode           string      `json:"assetCode"`
	FromAccountID       string      `json:"fromAccountID"`
	FromAccountName     string      `json:"fromAccountName"`
	FromUsername        string      `json:"fromUsername"`
	ToAccountID         string      `json:"toAccountID"`
	ToAccountName       string      `json:"toAccountName"`
	ToUsername          string      `json:"toUsername"`
	ToAssertion         string      `json:"toAssertion"`
	OriginalToAssertion string      `json:"originalToAssertion"`
	Note                string      `json:"note"`
	NoteErr             string      `json:"noteErr"`
	SourceAmountMax     string      `json:"sourceAmountMax"`
	SourceAmountActual  string      `json:"sourceAmountActual"`
	SourceAsset         sourceAsset `json:"sourceAsset"`
	SourceConvRate      string      `json:"sourceConvRate"`
	IsAdvanced          bool        `json:"isAdvanced"`
	SummaryAdvanced     string      `json:"summaryAdvanced"`
	Operations          interface{} `json:"operations"`
	Unread              bool        `json:"unread"`
	BatchID             string      `json:"batchID"`
	FromAirdrop         bool        `json:"fromAirdrop"`
	IsInflation         bool        `json:"isInflation"`
}
type details struct {
	PublicNote            string      `json:"publicNote"`
	PublicNoteType        string      `json:"publicNoteType"`
	ExternalTxURL         string      `json:"externalTxURL"`
	FeeChargedDescription string      `json:"feeChargedDescription"`
	PathIntermediate      interface{} `json:"pathIntermediate"`
}
type notification struct {
	Summary summary `json:"summary"`
	Details details `json:"details"`
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
	Messages      []messages     `json:"messages,omitempty"`
	Pagination    pagination     `json:"pagination"`
	Message       string         `json:"message"`
	ID            int            `json:"id"`
	Ratelimits    []rateLimits   `json:"ratelimits"`
	Conversations []conversation `json:"conversations,omitempty"`
	Offline       bool           `json:"offline,omitempty"`
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
type SendPayment struct {
	PaymentID string `json:"paymentID"`
}

// WalletAPI holds data for sending to API
type WalletAPI struct {
	Method string   `json:"method,omitempty"`
	Params *wParams `json:"params,omitempty"`
	Result *wResult `json:"result,omitempty"`
	Error  *Error   `json:"error"`
}
type wOptions struct {
	Name string `json:"name"`
	Txid string `json:"txid"`
}
type wParams struct {
	Options wOptions `json:"options"`
}
type asset struct {
	Type           string `json:"type"`
	Code           string `json:"code"`
	Issuer         string `json:"issuer"`
	VerifiedDomain string `json:"verifiedDomain"`
	IssuerName     string `json:"issuerName"`
	Desc           string `json:"desc"`
	InfoURL        string `json:"infoUrl"`
}
type sourceAsset struct {
	Type           string `json:"type"`
	Code           string `json:"code"`
	Issuer         string `json:"issuer"`
	VerifiedDomain string `json:"verifiedDomain"`
	IssuerName     string `json:"issuerName"`
	Desc           string `json:"desc"`
	InfoURL        string `json:"infoUrl"`
	InfoURLText    string `json:"infoUrlText"`
}
type balance struct {
	Asset  asset  `json:"asset"`
	Amount string `json:"amount"`
	Limit  string `json:"limit"`
}
type exchangeRate struct {
	Currency string `json:"currency"`
	Rate     string `json:"rate"`
}
type wResult struct {
	AccountID          string       `json:"accountID"`
	IsPrimary          bool         `json:"isPrimary"`
	Name               string       `json:"name"`
	Balance            []balance    `json:"balance"`
	ExchangeRate       exchangeRate `json:"exchangeRate"`
	AccountMode        int          `json:"accountMode"`
	TxID               string       `json:"txID"`
	Time               int64        `json:"time"`
	Status             string       `json:"status"`
	StatusDetail       string       `json:"statusDetail"`
	Amount             string       `json:"amount"`
	Asset              asset        `json:"asset"`
	DisplayAmount      string       `json:"displayAmount"`
	DisplayCurrency    string       `json:"displayCurrency"`
	SourceAmountMax    string       `json:"sourceAmountMax"`
	SourceAmountActual string       `json:"sourceAmountActual"`
	SourceAsset        sourceAsset  `json:"sourceAsset"`
	FromStellar        string       `json:"fromStellar"`
	ToStellar          string       `json:"toStellar"`
	FromUsername       string       `json:"fromUsername"`
	ToUsername         string       `json:"toUsername"`
	Note               string       `json:"note"`
	NoteErr            string       `json:"noteErr"`
	Unread             bool         `json:"unread"`
}

// TeamAPI holds information sent and received to/from the team api
type TeamAPI struct {
	Method string   `json:"method,omitempty"`
	Params *tParams `json:"params,omitempty"`
	Result *tResult `json:"result,omitempty"`
	Error  *Error   `json:"error"`
}
type emails struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}
type usernames struct {
	Username string `json:"username"`
	Role     string `json:"role"`
}
type user struct {
	UID      string `json:"uid"`
	Username string `json:"username"`
}
type uv struct {
	UID         string `json:"uid"`
	EldestSeqno int    `json:"eldestSeqno"`
}
type owners struct {
	Uv       uv     `json:"uv"`
	Username string `json:"username"`
	FullName string `json:"fullName"`
	NeedsPUK bool   `json:"needsPUK"`
	Status   int    `json:"status"`
}
type admins struct {
	Uv       uv     `json:"uv"`
	Username string `json:"username"`
	FullName string `json:"fullName"`
	NeedsPUK bool   `json:"needsPUK"`
	Status   int    `json:"status"`
}
type readers struct {
	Uv       uv     `json:"uv"`
	Username string `json:"username"`
	FullName string `json:"fullName"`
	NeedsPUK bool   `json:"needsPUK"`
	Status   int    `json:"status"`
}
type members struct {
	Owners  []owners      `json:"owners"`
	Admins  []admins      `json:"admins"`
	Writers []interface{} `json:"writers"`
	Readers []readers     `json:"readers"`
}
type annotatedActiveInvites struct {
}
type settings struct {
	Open   bool `json:"open"`
	JoinAs int  `json:"joinAs"`
}
type showcase struct {
	IsShowcased       bool `json:"is_showcased"`
	AnyMemberShowcase bool `json:"any_member_showcase"`
}
type tOptions struct {
	Team      string      `json:"team"`
	Emails    []emails    `json:"emails"`
	Usernames []usernames `json:"usernames"`
}
type tParams struct {
	Options tOptions `json:"options"`
}
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
type tResult struct {
	ChatSent               bool                   `json:"chatSent"`
	CreatorAdded           bool                   `json:"creatorAdded"`
	Invited                bool                   `json:"invited"`
	User                   user                   `json:"user"`
	EmailSent              bool                   `json:"emailSent"`
	ChatSending            bool                   `json:"chatSending"`
	Members                members                `json:"members"`
	KeyGeneration          int                    `json:"keyGeneration"`
	AnnotatedActiveInvites annotatedActiveInvites `json:"annotatedActiveInvites"`
	Settings               settings               `json:"settings"`
	Showcase               showcase               `json:"showcase"`
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
	Delete(messageID int) (ChatAPI, error)
	Edit(messageID int, message ...string) (ChatAPI, error)
	React(messageID int, reaction string) (ChatAPI, error)
	Send(message ...string) (ChatAPI, error)
}

type chatAPI interface {
	Next(count ...int) (*ChatAPI, error)
	Previous(count ...int) (*ChatAPI, error)
}

// Team holds basic information about a team
type Team struct {
	keybase *Keybase
	Name    string
}

type team interface {
	AddUser(user, role string) (TeamAPI, error)
	CreateSubteam(name string) (TeamAPI, error)
	MemberList() (TeamAPI, error)
}

type keybase interface {
	ChatList() (ChatAPI, error)
	CreateTeam(name string) (TeamAPI, error)
	NewChat(channel Channel) Chat
	NewTeam(name string) Team
	Run(handler func(ChatAPI), options ...RunOptions)
	loggedIn() bool
	username() string
	version() string
}

type status struct {
	Username string `json:"Username"`
	LoggedIn bool   `json:"LoggedIn"`
}
