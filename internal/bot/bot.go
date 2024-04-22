package bot

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"wjourney/internal/imagegeneration"

	"github.com/bwmarrin/discordgo"
)

var (
    appID = "1231817741669896213"
    guildID = ""
)

func Run(BotToken string, OpenAiKey string) {
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
    })
    if err != nil { fmt.Println("Bot 2"); log.Fatal(err) }

    discord.AddHandler(func (
        s *discordgo.Session,
        i * discordgo.InteractionCreate,
    ) {
        if i.Type == discordgo.InteractionApplicationCommand {
            data := i.ApplicationCommandData()
            switch data.Name {
            case "generate-image":
                if i.Interaction.Member.User.ID == s.State.User.ID { fmt.Println("Within func"); return }

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
                if prompt == "" { fmt.Println("Within func"); return }

                timeOutMsg := &discordgo.MessageSend{
                    Content: "Request timed out",
                }
                url, err := imagegeneration.GetImageUrl(prompt, OpenAiKey)
                if err != nil {
                    _, _ = s.ChannelMessageSendComplex(i.ChannelID, timeOutMsg)
                    fmt.Println("Within func"); return
                }

                msg := &discordgo.MessageSend{
                    Content: url,
                }

                msg2 := &discordgo.MessageSend{
                    Content: prompt + ":",
                }
                _, err = s.ChannelMessageSendComplex(i.ChannelID, msg2)
                _, err = s.ChannelMessageSendComplex(i.ChannelID, msg)
                if err != nil { fmt.Println("Within func"); return }
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
