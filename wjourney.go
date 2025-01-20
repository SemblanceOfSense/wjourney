package main

import (
	"flag"
	"wjourney/internal/bot"
)

var AiKey, BotToken string

func init() {
    flag.StringVar(&AiKey, "aikey", "", "LLM api key")
    flag.StringVar(&BotToken, "bottoken", "", "discord bot token")

    flag.Parse()
}

func main() {
    bot.Run(BotToken, AiKey)
}
