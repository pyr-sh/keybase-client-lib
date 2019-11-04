package keybase

func ExampleKeybase_AdvertiseCommand() {
	var k = NewKeybase()

	// Clear out any previously advertised commands
	k.ClearCommands()

	// Create BotAdvertisement
	c := BotAdvertisement{
		Type: "public",
		BotCommands: []BotCommand{
			NewBotCommand("help", "Get help using this bot"),
			NewBotCommand("hello", "Say hello"),
		},
	}

	// Send advertisement
	k.AdvertiseCommand(c)
}
