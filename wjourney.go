package main

import (
	"flag"
	"wjourney/internal/bot"
)

var OpenAiKey, BotToken, GrokKey string

func init() {
    flag.StringVar(&OpenAiKey, "openaikey", "", "LLM api key")
    flag.StringVar(&BotToken, "bottoken", "", "discord bot token")
    flag.StringVar(&GrokKey, "grok key", "", "Grok api key")

    flag.Parse()
}

func main() {
    bot.Run(BotToken, OpenAiKey, GrokKey)
}
