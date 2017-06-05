package main

import (
    "os"
    "fmt"
    "encoding/json"
    "github.com/jawher/mow.cli"

    "./auth"
    "./actions/query"
    "./actions/message"
)


func main () {
    app := cli.App("Locutus CLI", "A CLI for Locutus.")
    app.Spec = ""

    var profile = app.StringOpt("p profile", "", "AWS credential profile name")
    var region = app.StringOpt("r region", "us-east-1", "target AWS region")

    app.Command("self", "get instance metadata", func (cmd *cli.Cmd) {
        var ip = cmd.BoolOpt("i ip", false, "Get only the private IP")
        var az = cmd.BoolOpt("a az", false, "Get only the availability zone")
        var re = cmd.BoolOpt("r region", false, "Get only the region")
        var id = cmd.BoolOpt("d id", false, "Get only the instance ID")
        var ty = cmd.BoolOpt("t type", false, "Get only the instance type")

        cmd.Action = func () {
            session := auth.GetSession(*profile, *region)
            q := query.New(session)
            metadata := q.GetSelf()

            if *ip {
                fmt.Println("\"" + metadata.PrivateIP + "\"")
                return
            }
            if *az {
                fmt.Println("\"" + metadata.AvailabilityZone + "\"")
                return
            }
            if *re {
                fmt.Println("\"" + metadata.Region + "\"")
                return
            }
            if *id {
                fmt.Println("\"" + metadata.InstanceID + "\"")
                return
            }
            if *ty {
                fmt.Println("\"" + metadata.InstanceType + "\"")
                return
            }
            output, err := json.MarshalIndent(metadata, "", "    ")
            if err != nil {
                msg := fmt.Sprintf("Parse Error: %s", err)
                fmt.Println(msg)
                os.Exit(2)
            }
            os.Stdout.Write(output)
        }
    })

    app.Command("nodes", "find nodes by tag", func (cmd *cli.Cmd) {
        var key = cmd.StringOpt("k key", "", "the key to filter on")
        var values = cmd.StringsOpt("v value", nil, "the value to filter on")
        var ex = cmd.BoolOpt("x exclude-self", false, "exclude self")

        cmd.Action = func () {
            session := auth.GetSession(*profile, *region)
            q := query.New(session)
            ips := q.GetPrivateIPs(key, values)
            if *ex {
                i := -1
                metadata := q.GetSelf()
                for index, value := range ips {
                    if value == metadata.PrivateIP {
                        i = index
                    }
                }
                if i > -1 {
                    ips = append(ips[:i], ips[i + 1:]...)
                }
            }
            output, err := json.MarshalIndent(ips, "", "    ")
            if err != nil {
                msg := fmt.Sprintf("Parse Error: %s", err)
                fmt.Println(msg)
                os.Exit(2)
            }
            os.Stdout.Write(output)
        }
    })

    app.Command("message", "send a message to slack", func (cmd *cli.Cmd) {
        var text = cmd.StringArg("TEXT", "", "")
        var channel = cmd.StringOpt("c channel", "locutus", "slack channel")
        var note = cmd.StringOpt("n note", "", "attachment note / sub-text")
        cmd.Action = func () {
            if *text != "" {
                message.Send(channel, text, note)
            }
        }
    })

    app.Run(os.Args)
}
