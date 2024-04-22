package main

import (
	"flag"
	"wjourney/internal/bot"
)

var OpenaiKey, BotToken string

func init() {
    flag.StringVar(&OpenaiKey, "openaikey", "", "openai api key")
    flag.StringVar(&BotToken, "bottoken", "", "discord bot token")

    flag.Parse()
}

func main() {
    bot.Run(BotToken, OpenaiKey)
}
