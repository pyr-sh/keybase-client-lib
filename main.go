package main

import (
	"fmt"

	"keybase/api"
)

func main() {
	// Create new keybase api instance
	k := api.New()

	// Get username, logged in status, and client version
	username := k.Username()
	loggedin := k.LoggedIn()
	version := k.Version()

	// Send current client version to self if client is logged in.
	if loggedin {
		chatList, _ := k.ChatList()
		allChats := ""
		for _, chat := range chatList {
			if chat.Channel.MembersType == "team" {
				allChats += fmt.Sprintf("%s#%s\n", chat.Channel.Name, chat.Channel.TopicName)
			} else {
				allChats += fmt.Sprintf("%s\n", chat.Channel.Name)
			}
		}
		c, _ := k.ChatSendText(username, fmt.Sprintf("Version: %s\nConversations:\n```%s```\n", version, allChats))
		fmt.Println(c.Message, "-", c.ID)
	} else {
		fmt.Println("Not logged in")
	}
}
