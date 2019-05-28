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
		c, _ := k.ChatSend(username, version)
		fmt.Println(c.Result.Message)
	} else {
		fmt.Println("Not logged in")
	}
}
