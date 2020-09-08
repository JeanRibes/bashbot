package main

import (
	"bufio"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io"
	"os/exec"
	"time"
)

type ChannelCommand struct {
	process *exec.Cmd
	stdin   io.WriteCloser
	stdout  *bufio.Reader
}

func (cc *ChannelCommand) Input(s string) {
	_, err := cc.stdin.Write([]byte(s + "\necho '\x03'\n"))
	he(err)
}
func (cc *ChannelCommand) GetOutput() string {
	bout, err := cc.stdout.ReadBytes('\x03')
	he(err)
	return string(bout)
}

var commandsChannels = map[string]ChannelCommand{}

func init() {
	commandsChannels = make(map[string]ChannelCommand)
}
func handleBash(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Printf("%s : %s\n", m.Author.Username, m.Message.Content)

	if _, ok := commandsChannels[m.ChannelID]; !ok {
		var newCC ChannelCommand
		newCC.process = exec.Command("bash")

		newCC.stdin, _ = newCC.process.StdinPipe()
		stdout, _ := newCC.process.StdoutPipe()
		newCC.stdout = bufio.NewReader(stdout)
		he(newCC.process.Start())
		commandsChannels[m.ChannelID] = newCC
		newCC.Input("cd /tmp\n")
	}
	cc := commandsChannels[m.ChannelID]
	cc.Input(m.Content[2:])

	out := cc.GetOutput()
	fmt.Printf("\n\n\n\nsortie: %s\n", out)
	if len(out) < maxScrollLength {
		//if strings.Count(out, "\n") < 10 {
		s.ChannelMessageSend(m.ChannelID, "```"+out+"```")
		return
		//}
	}
	createScroll(s, m, out) // dans le futur il faudrait update le message au cours de l'exécution du script
}

var bashed []*discordgo.User

func addBash(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Printf("%v+", m.Mentions)
	bashed = append(bashed, m.Mentions...)
}

/*
:troll_face:
*/
func bash(s *discordgo.Session, m *discordgo.MessageCreate) {
	for _, user := range bashed {
		if user.ID == m.Author.ID {
			for _, lettre := range "osef" {
				emoji, ok := emojis[lettre]
				if ok {
					s.MessageReactionAdd(m.ChannelID, m.ID, emoji)
					time.Sleep(time.Millisecond * 100)
				}
			}
		}
	}
}
