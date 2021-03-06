package keybase

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// RunOptions holds a set of options to be passed to Run
type RunOptions struct {
	Capacity       int       // Channel capacity for the buffered channel that holds messages. Defaults to 100 if not set
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
	Type         string           `json:"type,omitempty"`
	Source       string           `json:"source,omitempty"`
	Msg          *msg             `json:"msg,omitempty"`
	Method       string           `json:"method,omitempty"`
	Params       *params          `json:"params,omitempty"`
	Message      string           `json:"message,omitempty"`
	ID           int              `json:"id,omitempty"`
	Ratelimits   []rateLimits     `json:"ratelimits,omitempty"`
	Notification *notification    `json:"notification,omitempty"`
	Result       *result          `json:"result,omitempty"`
	Pagination   *pagination      `json:"pagination,omitempty"`
	ErrorRaw     *json.RawMessage `json:"error,omitempty"` // Raw JSON string containing any errors returned
	ErrorRead    *Error           `json:"-"`               // Errors returned by any outgoing chat functions such as Read(), Edit(), etc
	ErrorListen  *string          `json:"-"`               // Errors returned by the api-listen command (used in the Run() function)
	keybase      Keybase          // Some methods will need this, so I'm passing it but keeping it unexported
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
	ReplyTo      int            `json:"replyTo"`
	ReplyToUID   string         `json:"replyToUID"`
	UserMentions []userMentions `json:"userMentions"`
	TeamMentions []teamMentions `json:"teamMentions"`
}

type flip struct {
	Text         string      `json:"text"`
	GameID       string      `json:"game_id"`
	FlipConvID   string      `json:"flip_conv_id"`
	UserMentions interface{} `json:"user_mentions"`
	TeamMentions interface{} `json:"team_mentions"`
}

type image struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type metadata struct {
	AssetType int   `json:"assetType"`
	Image     image `json:"image"`
}

type preview struct {
	Filename  string   `json:"filename"`
	Region    string   `json:"region"`
	Endpoint  string   `json:"endpoint"`
	Bucket    string   `json:"bucket"`
	Path      string   `json:"path"`
	Size      int      `json:"size"`
	MimeType  string   `json:"mimeType"`
	EncHash   string   `json:"encHash"`
	Key       string   `json:"key"`
	VerifyKey string   `json:"verifyKey"`
	Title     string   `json:"title"`
	Nonce     string   `json:"nonce"`
	Metadata  metadata `json:"metadata"`
	Tag       int      `json:"tag"`
}

type previews struct {
	Filename  string   `json:"filename"`
	Region    string   `json:"region"`
	Endpoint  string   `json:"endpoint"`
	Bucket    string   `json:"bucket"`
	Path      string   `json:"path"`
	Size      int      `json:"size"`
	MimeType  string   `json:"mimeType"`
	EncHash   string   `json:"encHash"`
	Key       string   `json:"key"`
	VerifyKey string   `json:"verifyKey"`
	Title     string   `json:"title"`
	Nonce     string   `json:"nonce"`
	Metadata  metadata `json:"metadata"`
	Tag       int      `json:"tag"`
}

type object struct {
	Filename  string   `json:"filename"`
	Region    string   `json:"region"`
	Endpoint  string   `json:"endpoint"`
	Bucket    string   `json:"bucket"`
	Path      string   `json:"path"`
	Size      int      `json:"size"`
	MimeType  string   `json:"mimeType"`
	EncHash   string   `json:"encHash"`
	Key       string   `json:"key"`
	VerifyKey string   `json:"verifyKey"`
	Title     string   `json:"title"`
	Nonce     string   `json:"nonce"`
	Metadata  metadata `json:"metadata"`
	Tag       int      `json:"tag"`
}

type attachment struct {
	Object   object     `json:"object"`
	Preview  preview    `json:"preview"`
	Previews []previews `json:"previews"`
	Metadata metadata   `json:"metadata"`
	Uploaded bool       `json:"uploaded"`
}

type content struct {
	Type           string         `json:"type"`
	Attachment     attachment     `json:"attachment"`
	Delete         delete         `json:"delete"`
	Edit           edit           `json:"edit"`
	Reaction       reaction       `json:"reaction"`
	System         system         `json:"system"`
	Text           text           `json:"text"`
	SendPayment    SendPayment    `json:"send_payment"`
	RequestPayment RequestPayment `json:"request_payment"`
	Flip           flip           `json:"flip"`
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
	Name        string `json:"name,omitempty"`
	Public      bool   `json:"public,omitempty"`
	MembersType string `json:"members_type,omitempty"`
	TopicType   string `json:"topic_type,omitempty"`
	TopicName   string `json:"topic_name,omitempty"`
}

type BotCommand struct {
	Name                string                         `json:"name"`
	Description         string                         `json:"description"`
	Usage               string                         `json:"usage"`
	ExtendedDescription *BotCommandExtendedDescription `json:"extended_description,omitempty"`
}

type BotCommandExtendedDescription struct {
	Title       string `json:"title"`
	DesktopBody string `json:"desktop_body"`
	MobileBody  string `json:"mobile_body"`
}

type BotAdvertisement struct {
	Type        string       `json:"type"`                // "public", "teamconvs", "teammembers"
	TeamName    string       `json:"team_name,omitempty"` // required if Type is not "public"
	BotCommands []BotCommand `json:"commands"`
}

type mesg struct {
	Body string `json:"body"`
}

type duration struct {
	time.Duration
}

func (d *duration) UnmarshalJSON(b []byte) (err error) {
	d.Duration, err = time.ParseDuration(strings.Trim(string(b), `"`))
	return
}

func (d *duration) MarshalJSON() (b []byte, err error) {
	return []byte(fmt.Sprintf(`"%s"`, d.String())), nil
}

type options struct {
	Channel            *Channel           `json:"channel,omitempty"`
	MessageID          int                `json:"message_id,omitempty"`
	Message            *mesg              `json:"message,omitempty"`
	Pagination         *pagination        `json:"pagination,omitempty"`
	Filename           string             `json:"filename,omitempty,omitempty"`
	Title              string             `json:"title,omitempty,omitempty"`
	Output             string             `json:"output,omitempty,omitempty"`
	ConversationID     string             `json:"conversation_id,omitempty"`
	FlipConversationID string             `json:"flip_conversation_id,omitempty"`
	MsgID              int                `json:"msg_id,omitempty"`
	ReplyTo            int                `json:"reply_to,omitempty"`
	GameID             string             `json:"game_id,omitempty"`
	Alias              string             `json:"alias,omitempty"`
	BotAdvertisements  []BotAdvertisement `json:"advertisements,omitempty"`
	ExplodingLifetime  duration           `json:"exploding_lifetime,omitempty"`

	Name        string `json:"name,omitempty"`
	Public      bool   `json:"public,omitempty"`
	MembersType string `json:"members_type,omitempty"`
	TopicType   string `json:"topic_type,omitempty"`
	TopicName   string `json:"topic_name,omitempty"`
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

type participants struct {
	UID        string `json:"uid"`
	DeviceID   string `json:"deviceID"`
	Username   string `json:"username"`
	DeviceName string `json:"deviceName"`
	Commitment string `json:"commitment"`
	Reveal     string `json:"reveal"`
}

type dupreg struct {
	User   string `json:"user"`
	Device string `json:"device"`
}

type errorInfo struct {
	Typ    int    `json:"typ"`
	Dupreg dupreg `json:"dupreg"`
}

type resultInfo struct {
	Typ  int  `json:"typ"`
	Coin bool `json:"coin"`
}

type flipStatus struct {
	GameID                  string         `json:"gameID"`
	Phase                   int            `json:"phase"`
	ProgressText            string         `json:"progressText"`
	ResultText              string         `json:"resultText"`
	CommitmentVisualization string         `json:"commitmentVisualization"`
	RevealVisualization     string         `json:"revealVisualization"`
	Participants            []participants `json:"participants"`
	ResultInfo              *resultInfo    `json:"resultInfo"`
	ErrorInfo               *errorInfo     `json:"errorInfo"`
}

type result struct {
	Messages         []messages     `json:"messages,omitempty"`
	Pagination       pagination     `json:"pagination"`
	Message          string         `json:"message"`
	ID               int            `json:"id"`
	Ratelimits       []rateLimits   `json:"ratelimits"`
	Conversations    []conversation `json:"conversations,omitempty"`
	Offline          bool           `json:"offline,omitempty"`
	Status           flipStatus     `json:"status,omitempty"`
	IdentifyFailures interface{}    `json:"identifyFailures,omitempty"`
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

type RequestPayment struct {
	RequestID string `json:"requestID"`
	Note      string `json:"note"`
}

// WalletAPI holds data for sending to API
type WalletAPI struct {
	Method string   `json:"method,omitempty"`
	Params *wParams `json:"params,omitempty"`
	Result *wResult `json:"result,omitempty"`
	Error  *Error   `json:"error"`
}

type wOptions struct {
	Name      string `json:"name"`
	Txid      string `json:"txid"`
	Recipient string `json:"recipient"`
	Amount    string `json:"amount"`
	Currency  string `json:"currency"`
	Message   string `json:"message"`
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
	Username           string       `json:"username"`
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

type member struct {
	Uv       uv     `json:"uv"`
	Username string `json:"username"`
	FullName string `json:"fullName"`
	NeedsPUK bool   `json:"needsPUK"`
	Status   int    `json:"status"`
}

type members struct {
	Owners  []member `json:"owners"`
	Admins  []member `json:"admins"`
	Writers []member `json:"writers"`
	Readers []member `json:"readers"`
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
	Username  string      `json:"username"`
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
	Teams                  []teamInfo             `json:"teams"`
}

type implicit struct {
	Role     int    `json:"role"`
	Ancestor string `json:"ancestor"`
}

type teamInfo struct {
	UID                     string   `json:"uid"`
	TeamID                  string   `json:"team_id"`
	Username                string   `json:"username"`
	FullName                string   `json:"full_name"`
	FqName                  string   `json:"fq_name"`
	IsImplicitTeam          bool     `json:"is_implicit_team"`
	ImplicitTeamDisplayName string   `json:"implicit_team_display_name"`
	IsOpenTeam              bool     `json:"is_open_team"`
	Role                    int      `json:"role"`
	NeedsPUK                bool     `json:"needsPUK"`
	MemberCount             int      `json:"member_count"`
	MemberEldestSeqno       int      `json:"member_eldest_seqno"`
	AllowProfilePromote     bool     `json:"allow_profile_promote"`
	IsMemberShowcased       bool     `json:"is_member_showcased"`
	Status                  int      `json:"status"`
	Implicit                implicit `json:"implicit,omitempty"`
}

// KVAPI holds information sent and received to/from the kvstore api
type KVAPI struct {
	Method  string    `json:"method,omitempty"`
	Params  *kvParams `json:"params,omitempty"`
	Result  *kvResult `json:"result,omitempty"`
	Error   *Error    `json:"error"`
	keybase Keybase
}

type kvOptions struct {
	Team       string `json:"team,omitempty"`
	Namespace  string `json:"namespace,omitempty"`
	EntryKey   string `json:"entryKey,omitempty"`
	Revision   uint   `json:"revision,omitempty"`
	EntryValue string `json:"entryValue,omitempty"`
}

type kvParams struct {
	Options kvOptions `json:"options,omitempty"`
}

type entryKey struct {
	EntryKey string `json:"entryKey"`
	Revision uint   `json:"revision"`
}

type kvResult struct {
	TeamName   string     `json:"teamName"`
	Namespaces []string   `json:"namespaces"`
	EntryKeys  []entryKey `json:"entryKeys"`
	EntryKey   string     `json:"entryKey"`
	EntryValue string     `json:"entryValue"`
	Revision   uint       `json:"revision"`
}

// UserAPI holds information received from the user/lookup api
type UserAPI struct {
	Status uStatus `json:"status"`
	Them   []them  `json:"them"`
}

type uStatus struct {
	Code int    `json:"code"`
	Name string `json:"name"`
}

type basics struct {
	Ctime                int    `json:"ctime"`
	EldestSeqno          int    `json:"eldest_seqno"`
	IDVersion            int    `json:"id_version"`
	LastIDChange         int    `json:"last_id_change"`
	Mtime                int    `json:"mtime"`
	PassphraseGeneration int    `json:"passphrase_generation"`
	RandomPw             bool   `json:"random_pw"`
	Salt                 string `json:"salt"`
	Status               int    `json:"status"`
	TrackVersion         int    `json:"track_version"`
	Username             string `json:"username"`
	UsernameCased        string `json:"username_cased"`
}

type profile struct {
	Bio      string `json:"bio"`
	FullName string `json:"full_name"`
	Location string `json:"location"`
	Mtime    int    `json:"mtime"`
}

type proof struct {
	HumanURL          string `json:"human_url"`
	Nametag           string `json:"nametag"`
	PresentationGroup string `json:"presentation_group"`
	PresentationTag   string `json:"presentation_tag"`
	ProofID           string `json:"proof_id"`
	ProofType         string `json:"proof_type"`
	ProofURL          string `json:"proof_url"`
	ServiceURL        string `json:"service_url"`
	SigID             string `json:"sig_id"`
	State             int    `json:"state"`
}

type proofsSummary struct {
	All    []proof `json:"all"`
	HasWeb bool    `json:"has_web"`
}

type key struct {
	KeyRole int    `json:"key_role"`
	Kid     string `json:"kid"`
	SigID   string `json:"sig_id"`
}

type uDevice struct {
	Ctime  int    `json:"ctime"`
	Keys   []key  `json:"keys"`
	Mtime  int    `json:"mtime"`
	Name   string `json:"name"`
	Status int    `json:"status"`
	Type   string `json:"type"`
}

type them struct {
	Basics        basics             `json:"basics,omitempty"`
	ID            string             `json:"id"`
	Profile       profile            `json:"profile,omitempty"`
	ProofsSummary proofsSummary      `json:"proofs_summary"`
	Devices       map[string]uDevice `json:"devices,omitempty"`
}

// UserCardAPI holds information received from the user/card api
type UserCardAPI struct {
	AirdropRegistered bool           `json:"airdrop_registered"`
	Blocked           bool           `json:"blocked"`
	FollowSummary     followSummary  `json:"follow_summary"`
	Profile           cardProfile    `json:"profile"`
	Status            uStatus        `json:"status"`
	TeamShowcase      []teamShowcase `json:"team_showcase"`
	TheyFollowYou     bool           `json:"they_follow_you"`
	UserBlocks        userBlocks     `json:"user_blocks"`
	YouFollowThem     bool           `json:"you_follow_them"`
}

type followSummary struct {
	Followers int `json:"followers"`
	Following int `json:"following"`
}

type cardProfile struct {
	Bio                    string    `json:"bio"`
	Comment                string    `json:"comment"`
	CrimeAll               int       `json:"crime_all"`
	CrimeChat              int       `json:"crime_chat"`
	CrimeFollow            int       `json:"crime_follow"`
	CrimeIllegal           int       `json:"crime_illegal"`
	CrimeLegacyAll         int       `json:"crime_legacy_all"`
	CrimeLegacyPorn        int       `json:"crime_legacy_porn"`
	CrimeLegacyStellar     int       `json:"crime_legacy_stellar"`
	CrimePorn              int       `json:"crime_porn"`
	CrimeSmurfing          int       `json:"crime_smurfing"`
	CrimeSpacedrop         int       `json:"crime_spacedrop"`
	CrimeStellarDust       int       `json:"crime_stellar_dust"`
	CrimeStellarPaymentReq int       `json:"crime_stellar_payment_req"`
	CrimeTeam              int       `json:"crime_team"`
	Ctime                  time.Time `json:"ctime"`
	FullName               string    `json:"full_name"`
	IsAdmin                int       `json:"is_admin"`
	Location               string    `json:"location"`
	Mtime                  time.Time `json:"mtime"`
	Reporter               string    `json:"reporter"`
	Status                 int       `json:"status"`
	Twitter                string    `json:"twitter"`
	UID                    string    `json:"uid"`
	Website                string    `json:"website"`
}

type teamShowcase struct {
	Description     string   `json:"description"`
	FqName          string   `json:"fq_name"`
	NumMembers      int      `json:"num_members"`
	Open            bool     `json:"open"`
	PublicAdmins    []string `json:"public_admins"`
	Role            int      `json:"role"`
	TeamID          string   `json:"team_id"`
	TeamIsShowcased bool     `json:"team_is_showcased"`
}

type userBlocks struct {
	Chat   bool      `json:"chat"`
	Ctime  time.Time `json:"ctime"`
	Follow bool      `json:"follow"`
	Mtime  time.Time `json:"mtime"`
}

// Keybase holds basic information about the local Keybase executable
type Keybase struct {
	Path     string
	Username string
	LoggedIn bool
	Version  string
	Device   string
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
	Reply(replyTo int, message ...string) (ChatAPI, error)
	Upload(title string, filepath string) (ChatAPI, error)
	Download(messageID int, filepath string) (ChatAPI, error)
	LoadFlip(messageID int, conversationID string, flipConversationID string, gameID string) (ChatAPI, error)
	Pin(messageID int) (ChatAPI, error)
	Unpin() (ChatAPI, error)
	Mark(messageID int) (ChatAPI, error)
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
	AddAdmins(users ...string) (TeamAPI, error)
	AddOwners(users ...string) (TeamAPI, error)
	AddReaders(users ...string) (TeamAPI, error)
	AddUser(user, role string) (TeamAPI, error)
	AddWriters(users ...string) (TeamAPI, error)
	CreateSubteam(name string) (TeamAPI, error)
	MemberList() (TeamAPI, error)
}

// Wallet holds basic information about a wallet
type Wallet struct {
	keybase *Keybase
}

type wallet interface {
	CancelRequest(requestID string) error
	RequestPayment(user string, amount float64, memo ...string)
	Send(recipient string, amount string, currency string, message ...string) (WalletAPI, error)
	SendXLM(recipient string, amount string, message ...string) (WalletAPI, error)
	StellarAddress(user string) (string, error)
	TxDetail(txid string) (WalletAPI, error)
}

// KV holds basic information about a KVStore
type KV struct {
	keybase *Keybase
	Team    string
}

type kvInterface interface {
	Namespaces() (KVAPI, error)
	Keys(namespace string) (KVAPI, error)
	Get(namespace string, key string) (KVAPI, error)
	Put(namespace string, key string, value string) (KVAPI, error)
	Delete(namespace string, key string) (KVAPI, error)
}

type keybase interface {
	AdvertiseCommand(advertisement BotAdvertisement) (ChatAPI, error)
	AdvertiseCommands(advertisements []BotAdvertisement) (ChatAPI, error)
	ChatList(opts ...Channel) (ChatAPI, error)
	ClearCommands() (ChatAPI, error)
	CreateTeam(name string) (TeamAPI, error)
	NewChat(channel Channel) Chat
	NewTeam(name string) Team
	NewKV(team string) KV
	NewWallet() Wallet
	Run(handler func(ChatAPI), options ...RunOptions)
	status() status
	version() string
	UserLookup(users ...string) (UserAPI, error)
	ListUserMemberships(user string) (TeamAPI, error)
	UserCard(user string) (UserCardAPI, error)
}

type status struct {
	Username string `json:"Username"`
	LoggedIn bool   `json:"LoggedIn"`
	Device   device `json:"Device"`
}

type device struct {
	Name string `json:"name"`
}
