package api

import ()

type Chat struct {
	Type       string     `json:"type"`
	Source     string     `json:"source"`
	Msg        Msg        `json:"msg"`
	Pagination Pagination `json:"pagination"`
}
type Channel struct {
	Name        string `json:"name"`
	Public      bool   `json:"public"`
	MembersType string `json:"members_type"`
	TopicType   string `json:"topic_type"`
	TopicName   string `json:"topic_name"`
}
type Sender struct {
	UID        string `json:"uid"`
	Username   string `json:"username"`
	DeviceID   string `json:"device_id"`
	DeviceName string `json:"device_name"`
}
type Addedtoteam struct {
	Team    string   `json:"team"`
	Adder   string   `json:"adder"`
	Addee   string   `json:"addee"`
	Owners  []string `json:"owners"`
	Admins  []string `json:"admins"`
	Writers []string `json:"writers"`
	Readers []string `json:"readers"`
}
type Bulkaddtoconv struct {
	Usernames []string `json:"usernames"`
}
type System struct {
	SystemType    int           `json:"systemType"`
	Addedtoteam   Addedtoteam   `json:"addedtoteam"`
	Bulkaddtoconv Bulkaddtoconv `json:"bulkaddtoconv"`
}
type Result struct {
	ResultTyp int    `json:"resultTyp"`
	Sent      string `json:"sent"`
}
type Payments struct {
	Username    string `json:"username"`
	PaymentText string `json:"paymentText"`
	Result      Result `json:"result"`
}
type UserMentions struct {
	Text string `json:"text"`
	UID  string `json:"uid"`
}
type TeamMentions struct {
	Name    string `json:"name"`
	Channel string `json:"channel"`
}
type Reaction struct {
	M int    `json:"m"`
	B string `json:"b"`
}
type Delete struct {
	MessageIDs []int `json:"messageIDs"`
}
type Edit struct {
	MessageID    int            `json:"messageID"`
	Body         string         `json:"body"`
	Payments     []Payments     `json:"payments"`
	UserMentions []UserMentions `json:"userMentions"`
	TeamMentions []TeamMentions `json:"teamMentions"`
}
type Text struct {
	Body         string         `json:"body"`
	Payments     []Payments     `json:"payments"`
	UserMentions []UserMentions `json:"userMentions"`
	TeamMentions []TeamMentions `json:"teamMentions"`
}
type Content struct {
	Type     string   `json:"type"`
	Delete   Delete   `json:"delete"`
	Edit     Edit     `json:"edit"`
	Reaction Reaction `json:"reaction"`
	System   System   `json:"system"`
	Text     Text     `json:"text"`
}
type Msg struct {
	ID                 int         `json:"id"`
	Channel            Channel     `json:"channel"`
	Sender             Sender      `json:"sender"`
	SentAt             int         `json:"sent_at"`
	SentAtMs           int64       `json:"sent_at_ms"`
	Content            Content     `json:"content"`
	Prev               interface{} `json:"prev"`
	Unread             bool        `json:"unread"`
	AtMentionUsernames []string    `json:"at_mention_usernames"`
	IsEphemeral        bool        `json:"is_ephemeral"`
	Etime              int64       `json:"etime"`
	HasPairwiseMacs    bool        `json:"has_pairwise_macs"`
	ChannelMention     string      `json:"channel_mention"`
}
type Pagination struct {
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Num      int    `json:"num"`
	Last     bool   `json:"last"`
}
