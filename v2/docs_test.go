package keybase

import "samhofi.us/x/keybase/v2/types/chat1"

func ExampleKeybase_AdvertiseCommands() {
	var k = NewKeybase()

	// Clear out any previously advertised commands
	k.ClearCommands()

	// Create BotAdvertisement
	ads := AdvertiseCommandsOptions{
		Alias: "RSS Bot",
		Advertisements: []chat1.AdvertiseCommandAPIParam{
			{
				Typ: "public",
				Commands: []chat1.UserBotCommandInput{
					{
						Name:        "rss addfeed",
						Description: "Add RSS feed",
						Usage:       "<url>",
					},
					{
						Name:        "rss delfeed",
						Description: "Remove RSS feed",
						Usage:       "<url>",
					},
				},
			},
		},
	}

	// Send advertisement
	k.AdvertiseCommands(ads)
}
