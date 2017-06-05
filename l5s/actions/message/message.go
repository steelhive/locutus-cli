package message

import (
    "os"
    "fmt"
    "github.com/nlopes/slack"
)

func Send (channel *string, text *string, note *string) {
    token := os.Getenv("SLACK_BOT_TOKEN")
    api := slack.New(token)
    params := slack.PostMessageParameters{}
    params.AsUser = true
    if *note != "" {
        attachment := slack.Attachment{
            Text: *note,
        }
        params.Attachments = []slack.Attachment{attachment}
    }

    result, _, err := api.PostMessage(*channel, *text, params)
    if err != nil {
        fmt.Printf("%s\n", err)
        return
    }
    fmt.Printf("message sent to %s", result)
}
