package main

import (
	"fmt"

	"keybase/api"
)

func main() {
	k := api.New()
	fmt.Printf("Username: %v\nLogged In: %v\nVersion: %v", k.Username(), k.LoggedIn(), k.Version())
}
