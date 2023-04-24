package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/eteissonniere/hercules/llms"

	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

var cmdReplay = cli.Command{
	Name:  "replay",
	Usage: "Replay a log file",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "file",
			Usage:    "Path to the log file",
			Required: true,
		},
	},
	Action: func(c *cli.Context) error {
		headerSystem := color.New(color.FgHiBlue, color.Bold)
		headerAssistant := color.New(color.FgHiYellow, color.Bold)
		headerUser := color.New(color.FgHiMagenta, color.Bold)
		content := color.New(color.FgHiWhite)
		command := color.New(color.FgHiGreen, color.Bold)

		file, err := os.Open(c.String("file"))
		if err != nil {
			return fmt.Errorf("failed to open file: %w", err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			var msg llms.ChatMessage
			err := json.Unmarshal(scanner.Bytes(), &msg)
			if err != nil {
				return fmt.Errorf("failed to unmarshal message: %w", err)
			}
			msg.Content = strings.TrimSpace(msg.Content)

			headerPrint := headerSystem
			switch msg.Role {
			case llms.ChatRoleSystem:
				headerPrint = headerSystem
			case llms.ChatRoleAssistant:
				headerPrint = headerAssistant
			case llms.ChatRoleUser:
				headerPrint = headerUser
			}

			headerPrint.Println(strings.ToUpper(string(msg.Role)))

			if msg.Role == llms.ChatRoleSystem {
				content.Println(msg.Content)
				continue
			}

			allLines := strings.Split(msg.Content, "\n")
			agentContent := strings.Join(allLines[0:len(allLines)-2], "\n")
			agentCommand := allLines[len(allLines)-1]

			content.Println(agentContent)
			command.Println(agentCommand)
		}

		return nil
	},
}
