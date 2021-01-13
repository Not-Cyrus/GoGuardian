package utils

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fastjson"
)

func BanCreate(guildID string, userID string, reason string) string {
	res, stat := SendRequest("PUT", fmt.Sprintf("https://discord.com/api/v8/guilds/%s/bans/%s?reason=%s", guildID, userID, url.QueryEscape(reason)), "", nil)

	err := fastjson.Validate(res)
	if err != nil {
		return ""
	}

	parsed, err := parser.Parse(res)

	if err != nil {
		fmt.Printf("[JSON Ban error]: %s\n", err) // easier to just debug like this over getting end-users to send you the problem and leaving out crucial stuff.
		return err.Error()
	}

	message := string(parsed.GetStringBytes("message"))

	switch {
	case stat == 429:
		time.Sleep(time.Duration(parsed.GetInt("retry_after")) * time.Millisecond)
		BanCreate(guildID, userID, reason)
	case len(message) != 0:
		return message
	}

	return ""
}

func FindConfig(guildID string) (*fastjson.Value, *fastjson.Value) {
	FileContents := ReadFile("Config.json")
	parsed, err := parser.Parse(FileContents)
	if err != nil {
		SendMessage(nil, fmt.Sprintf("Error parsing json %s", err.Error()), "")
		return nil, nil
	}
	if !fastjson.Exists([]byte(FileContents), "Guilds", guildID) {
		parsed.Get("Guilds").Set(guildID, fastjson.MustParse(defaultConfig))
		SaveJSON(nil, nil, parsed, "")
	}
	guild := parsed.Get("Guilds", guildID)
	return parsed, guild
}

func GetToken() {
	fileContents := ReadFile("Config.json")
	parsed, err := parser.Parse(fileContents)
	if err != nil {
		fmt.Printf("Couldn't parse Config.json to get your Token: %s", err.Error())
		time.Sleep(5 * time.Second)
		os.Exit(0)
	}
	if fastjson.Exists([]byte(fileContents), "Token") {
		Token = string(parsed.GetStringBytes("Token"))
		return
	}
	fmt.Print("Enter your token: ")
	fmt.Scan(&Token)
}

func GetGuildOwner(s *discordgo.Session, guildID string) string {
	guild, err := s.Guild(guildID)
	if err != nil {
		fmt.Printf("Error getting guild: %s\n", err)
		return ""
	}
	return guild.OwnerID
}

func InArray(guildID string, arrayStr string, data *fastjson.Value, target string) (bool, int) {
	var array []*fastjson.Value
	switch len(arrayStr) {
	case 0:
		array = data.GetArray("Guilds")
	default:
		array = data.GetArray("Guilds", guildID, arrayStr)
	}
	for index, whitelistedUser := range array {
		if string(whitelistedUser.GetStringBytes()) == target {
			return true, index
		}
	}
	return false, 0
}

func ReadFile(fileName string) string {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("Failed to open that file: %s", err.Error())
		return ""
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Printf("I couldn't read the data: %s | Please reopen the program", err.Error())
		return ""
	}
	return string(data)
}

func SaveJSON(s *discordgo.Session, message *discordgo.Message, parsedData *fastjson.Value, sendMessage string) {
	if s != nil {
		s.ChannelMessageSend(message.ChannelID, sendMessage)
	}
	Writefile("config.json", string(parsedData.MarshalTo(nil)))
}

func SendMessage(s *discordgo.Session, message, userID string) {
	if len(userID) != 0 {
		channel, err := s.UserChannelCreate(userID)
		if err != nil {
			fmt.Printf("Couldn't make a channel on that UserID: %s\n", err.Error())
			return
		}
		s.ChannelMessageSend(channel.ID, message)
	}
	fmt.Printf("[GoGuardian]: %s\n", message)
}

func SendRequest(method, URL, ctype string, body []byte) (string, int) {

	req := fasthttp.AcquireRequest()
	res := fasthttp.AcquireResponse()

	if body != nil {
		req.SetBody(body)
	}

	req.Header.SetMethod(method)
	req.Header.SetRequestURI(URL)

	if len(ctype) != 0 {
		req.Header.Set("Content-Type", ctype)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bot %s", Token))

	err := httpClient.Do(req, res)

	if err != nil {
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(res)
		return "{}", 0
	}

	return string(res.Body()), res.StatusCode()
}

func Writefile(filename string, data string) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	defer file.Close()
	_, err = file.WriteAt([]byte(data), 0)
	if err != nil {
		fmt.Printf("Couldn't open/write to the file: %s", err)
	}
}

var (
	defaultConfig = `{"WhitelistedIDs": [],"Config": {"Threshold":2,"Seconds":2,"BanProtection":true,"KickProtection":true,"HijackProtection":true,"AntiBotProtection":true,"RoleSpamProtection":true,"RoleNukeProtection":true,"RoleUpdateProtection":true,"ChannelSpamProtection":true,"ChannelNukeProtection":true,"MemberRoleUpdateProtection":true,"WebhookProtection":true}}`
	httpClient    = &fasthttp.Client{}
	parser        fastjson.Parser
	Token         string
)
