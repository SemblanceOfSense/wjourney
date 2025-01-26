package bot

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"wjourney/internal/imagegeneration"
	"wjourney/internal/textgeneration"

	"github.com/bwmarrin/discordgo"
)

var (
    appID = "1231817741669896213"
    guildID = ""
)

func Run(BotToken string, OpenAiKey string, GrokKey string) {
    discord, err := discordgo.New(("Bot " + BotToken))
    if err != nil { fmt.Println("Bot 1"); log.Fatal(err) }

    _, err = discord.ApplicationCommandBulkOverwrite(appID, guildID, []*discordgo.ApplicationCommand {
        {
            Name: "generate-image",
            Description: "Generates an ai image from a prompt",
            Options: []*discordgo.ApplicationCommandOption {
                {
                    Type: discordgo.ApplicationCommandOptionString,
                    Name: "prompt",
                    Description: "prompt for the image",
                    Required: true,
                },
            },
        },
        {
            Name: "generate-text",
            Description: "Generates text from a prompt",
            Options: []*discordgo.ApplicationCommandOption {
                {
                    Type: discordgo.ApplicationCommandOptionString,
                    Name: "prompt",
                    Description: "prompt for the text",
                    Required: true,
                },
            },
        },
    })
    if err != nil { fmt.Println("Bot 2"); log.Fatal(err) }

    discord.AddHandler(func (
        s *discordgo.Session,
        i * discordgo.InteractionCreate,
    ) {
        var erro bool;
        if i.Type == discordgo.InteractionApplicationCommand {
            data := i.ApplicationCommandData()
            switch data.Name {
            case "generate-image":
                if i.Interaction.Member.User.ID == s.State.User.ID { fmt.Println("Within func 1"); return; }

                err = s.InteractionRespond(
                    i.Interaction,
                    &discordgo.InteractionResponse {
                        Type: discordgo.InteractionResponseChannelMessageWithSource,
                        Data: &discordgo.InteractionResponseData {
                            Content: "Wait please!",
                        },
                    },
                )
                if err != nil { fmt.Println("Bot 3"); log.Fatal(err) }

                var prompt string
                for _, v := range i.Interaction.ApplicationCommandData().Options {
                    switch v.Name {
                    case "prompt":
                        prompt = v.StringValue()
                    }
                }
                if prompt == "" { fmt.Println("Within func 2"); return; }

                timeOutMsg := &discordgo.MessageSend{
                    Content: "Request timed out",
                }
                url, err := imagegeneration.GetImageUrl(prompt, OpenAiKey)
                if err != nil {
                    _, _ = s.ChannelMessageSendComplex(i.ChannelID, timeOutMsg)
                    fmt.Println("Within func 3"); return;
                }
                if url == "rejected" { erro = true; }

                msg := &discordgo.MessageSend{
                    Content: url,
                }

                msg2 := &discordgo.MessageSend{
                    Content: prompt + ":",
                }
                errmsg := &discordgo.MessageSend{
                    Content: "Prompt Was Rejected",
                }
                if (erro == false) {
                    _, err = s.ChannelMessageSendComplex(i.ChannelID, msg2)
                    if err != nil { fmt.Println("Prompt message failed"); fmt.Println(err); return }
                    _, err = s.ChannelMessageSendComplex(i.ChannelID, msg)
                    if err != nil { fmt.Println("Response message failed"); return }
                }
                if (erro) {
                    _, _ = s.ChannelMessageSendComplex(i.ChannelID, errmsg)
                }
            case "generate-text":
                if i.Interaction.Member.User.ID == s.State.User.ID { fmt.Println("Within func 1"); return; }

                err = s.InteractionRespond(
                    i.Interaction,
                    &discordgo.InteractionResponse {
                        Type: discordgo.InteractionResponseChannelMessageWithSource,
                        Data: &discordgo.InteractionResponseData {
                            Content: "Wait please!",
                        },
                    },
                )
                if err != nil { fmt.Println("Bot 3"); log.Fatal(err) }

                var prompt string
                for _, v := range i.Interaction.ApplicationCommandData().Options {
                    switch v.Name {
                    case "prompt":
                        prompt = v.StringValue()
                    }
                }
                if prompt == "" { fmt.Println("Within func 2"); return; }

                timeOutMsg := &discordgo.MessageSend{
                    Content: "Request timed out",
                }
                resp, err := textgeneration.GetGeneratedText(prompt, GrokKey)
                if (strings.Contains(prompt, "foot") || strings.Contains(prompt, "feet")) {
                    resp = "rejected"
                }
                if err != nil {
                    _, _ = s.ChannelMessageSendComplex(i.ChannelID, timeOutMsg)
                    fmt.Println("Within func 3"); return;
                }
                if resp == "rejected" { erro = true; }

                msg := &discordgo.MessageSend{
                    Content: resp,
                }

                msg2 := &discordgo.MessageSend{
                    Content: prompt + ":",
                }
                errmsg := &discordgo.MessageSend{
                    Content: "Prompt Was Rejected",
                }
                if (erro == false) {
                    _, err = s.ChannelMessageSendComplex(i.ChannelID, msg2)
                    if err != nil { fmt.Println("Failed prompt message") }
                    _, err = s.ChannelMessageSendComplex(i.ChannelID, msg)
                    if err != nil { fmt.Println("Failed response message"); erro = true }
                }
                if (erro) {
                    _, _ = s.ChannelMessageSendComplex(i.ChannelID, errmsg)
                }
            }
        }
    })

    err = discord.Open()
    if err != nil { log.Fatal(err) }

    stop := make (chan os.Signal, 1)
    signal.Notify(stop, os.Interrupt)
    log.Println("Press Ctrl+C to Exit")
    <-stop

    err = discord.Close()
    if err != nil { log.Fatal(err) }
}
