/*
Package keybase implements an interface for interacting with the Keybase Chat, Team, and Wallet APIs

I've tried to follow Keybase's JSON API as closely as possible, so if you're stuck on anything, or wondering
why things are organized in a certain way, it's most likely due to that. It may be helpful to look at the
Keybase JSON API docs by running some of the following commands in your terminal:
    // Chat API
    keybase chat api -h

    // Chat Message Stream
    keybase chat api-listen -h

    // Team API
    keybase team api -h

    // Wallet API
    keybase wallet api -h

The git repo for this code is hosted on Keybase. You can contact me directly (https://keybase.io/dxb),
or join the mkbot team (https://keybase.io/team/mkbot) if you need assistance, or if you'd like to contribute.

Basic Example

Here's a quick example of a bot that will attach a reaction with the sender's device name to every message sent
in @mkbot#test1:

    package main

    import (
    	"fmt"

    	"samhofi.us/x/keybase"
    )

    var k = keybase.NewKeybase()

    func main() {
    	channel := keybase.Channel{
    		Name:        "mkbot",
    		TopicName:   "test1",
    		MembersType: keybase.TEAM,
    	}
    	opts := keybase.RunOptions{
    		FilterChannel: channel,
    	}
    	fmt.Println("Running...")
    	k.Run(handler, opts)
    }

    func handler(m keybase.ChatAPI) {
	if m.ErrorListen != nil {
		fmt.Printf("Error: %s\n", *m.ErrorListen)
		return
	}

    	msgType := m.Msg.Content.Type
    	msgID := m.Msg.ID
    	deviceName := m.Msg.Sender.DeviceName

    	if msgType == "text" {
    		chat := k.NewChat(m.Msg.Channel)
    		chat.React(msgID, deviceName)
    	}
    }
*/
package keybase
