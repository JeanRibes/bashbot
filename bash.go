package main

import (
	"bufio"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io"
	"os"
	"os/exec"
	"os/user"
	"strings"
	"time"
)

var userMessageMap map[string]string //stocke les relations entre la personne qui a inité la commande
// et les résultats dans le chat. sert à éviter que n'importe qui fasse supprimer les réulstats

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
	userMessageMap = map[string]string{}
}
func handleBash(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Printf("%s : %s\n", m.Author.Username, m.Message.Content)

	if _, ok := commandsChannels[m.ChannelID]; !ok {
		var newCC ChannelCommand
		newCC.process = exec.Command("bash")
		newCC.process.Env = []string{
			"USER=bashbot",
			"PWD=/tmp",
		}

		newCC.stdin, _ = newCC.process.StdinPipe()
		stdout, _ := newCC.process.StdoutPipe()
		newCC.stdout = bufio.NewReader(stdout)
		he(newCC.process.Start())
		commandsChannels[m.ChannelID] = newCC
	}
	cc := commandsChannels[m.ChannelID]
	cc.Input(m.Content[2:])

	out := cc.GetOutput()
	fmt.Printf("\n\n\n\nsortie: %s\n", out)
	if len(out) < maxScrollLength {
		//if strings.Count(out, "\n") < 10 {
		nmsg, serr := s.ChannelMessageSend(m.ChannelID, "```"+out+"```")
		if serr == nil {
			s.MessageReactionAdd(m.ChannelID, nmsg.ID, deleteEmoji)
			userMessageMap[nmsg.ID] = m.Author.ID
		}
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
					he(s.MessageReactionAdd(m.ChannelID, m.ID, emoji))
					time.Sleep(time.Millisecond * 100)
				}
			}
		}
	}
}

func prompt(wd string) string {
	user, ue := user.Current()
	he(ue)
	hostname, herr := os.Hostname()
	he(herr)
	return user.Username + "@" + hostname + ":" + wd + "$ "
}
func formatinput(in string) string {
	return strings.ReplaceAll(in, "\n", " ; ")
}
