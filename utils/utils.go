package utils

import (
	"io/ioutil"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/valyala/fastjson"
)

func ReadFile(fileName string) string {
	file, err := os.Open(fileName)
	if err != nil {
		panic("HELLO, DO YOU KNOW HOW TO MOVE FILES??")
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		panic("I couldn't read the data")
	}
	return string(data)
}

func Writefile(filename string, data string) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	defer file.Close()
	_, err = file.WriteAt([]byte(data), 0)
	if err != nil {
		panic("Couldn't open/write to the file")
	}
}

func InArray(arrayStr string, data *fastjson.Value, m *discordgo.Message, target string) (bool, int) {
	array := data.GetArray(arrayStr)
	for index, whitelistedUser := range array {
		if string(whitelistedUser.GetStringBytes()) == target {
			return true, index
		}
	}
	return false, 0
}
