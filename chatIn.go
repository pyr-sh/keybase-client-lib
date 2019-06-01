package keybase

import ()

type chatIn struct {
	Type   string    `json:"type"`
	Source string    `json:"source"`
	Msg    chatInMsg `json:"msg"`
}
type chatInChannel struct {
	Name        string `json:"name"`
	Public      bool   `json:"public"`
	MembersType string `json:"members_type"`
	TopicType   string `json:"topic_type"`
	TopicName   string `json:"topic_name"`
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
type chatInSystem struct {
	SystemType    int                 `json:"systemType"`
	Addedtoteam   chatInAddedtoteam   `json:"addedtoteam"`
	Bulkaddtoconv chatInBulkaddtoconv `json:"bulkaddtoconv"`
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
	Channel            chatInChannel `json:"channel"`
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
