package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func init() {
	file, err := os.Open("Config.json")
	if err != nil {
		panic("HELLO, DO YOU KNOW HOW TO MOVE FILES??")
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		panic("I couldn't read the data")
	}

	json.Unmarshal([]byte(data), &Config)

	// pretty sure indexing a key value is faster than searching through an array resulting in this code
	for _, v := range Config.WhitelistIDs {
		WhitelistedIDs[v] = "https://github.com/Not-Cyrus is pretty cool"
	}
}

type (
	configData struct {
		Token              string   `json:"token"`
		Seconds            float64  `json:"Seconds"`
		Threshold          int      `json:"Threshold"`
		BanEnabled         bool     `json:"BanProtection"`
		KickEnabled        bool     `json:"KickProtection"`
		RoleSpamEnabled    bool     `json:"RoleSpamProtection"`
		RoleNukeEnabled    bool     `json:"RoleNukeProtection"`
		ChannelSpamEnabled bool     `json:"ChannelSpamProtection"`
		ChannelNukeEnabled bool     `json:"ChannelNukeProtection"`
		WhitelistIDs       []string `json:"WhitelistedIDs,omitempty"`
	}
)

var (
	Config         = configData{}
	WhitelistedIDs = map[string]string{}
)
