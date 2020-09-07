package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os/exec"
)

const maxScrollLength = 100

const emojiUp = "⬆️"
const emojiDown = "⬇️"

type TextScroll struct {
	content string
	scroll  int
}

var trackedScrollText map[string]TextScroll

func init() {
	trackedScrollText = map[string]TextScroll{}
}

func messageCreate2(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	if m.Content == "p" {
		t, err := exec.Command("cat", "/Users/jean/Documents/INSA/LaTeX-template/cr.tex").CombinedOutput()
		he(err)
		ts := TextScroll{
			content: string(t),
			scroll:  0,
		}
		//fmt.Printf("sortie: %s\n",ts.content)
		//_, e := scrollMessage(&ts, s, m.Message.ID, m.ChannelID)
		var e error
		var msg *discordgo.Message
		if len(ts.content) < maxScrollLength {
			msg, e = s.ChannelMessageSend(m.ChannelID, ts.content)
		} else {
			msg, e = s.ChannelMessageSend(m.ChannelID, ts.content[:maxScrollLength])
		}
		if e != nil {
			fmt.Errorf("%s", e)
		}
		trackedScrollText[msg.ID] = ts //l'ID du message qu'on a émis
		s.MessageReactionAdd(m.ChannelID, msg.ID, emojiUp)
		s.MessageReactionAdd(m.ChannelID, msg.ID, emojiDown)
		return true
	}
	if m.Content == "d" {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%d tracked scrolls", len(trackedScrollText)))
		return true
	}
	return false
}

func scrollMessage(t TextScroll, s *discordgo.Session, messageId string, channelId string) (*discordgo.Message, error) {
	max := len(t.content)
	if max < maxScrollLength {
		return nil, nil
	}
	min := maxScrollLength * t.scroll
	mid := maxScrollLength * (t.scroll + 1)
	if mid > max {
		mid = max
	}
	fmt.Printf("min: %d, mid: %d, max: %d\n", min, mid, max)
	fmt.Printf("content: %s\n", t.content[min:mid])
	return s.ChannelMessageEdit(channelId, messageId, t.content[min:mid])
}

func messageReactionAdd(s *discordgo.Session, e *discordgo.MessageReactionAdd) {
	scroll(s, e.MessageReaction)
}
func messageReactionRemove(s *discordgo.Session, e *discordgo.MessageReactionRemove) {
	scroll(s, e.MessageReaction)
}

func scroll(s *discordgo.Session, e *discordgo.MessageReaction) {
	if e.UserID == s.State.User.ID {
		return
	}
	fmt.Printf("nom: %s id: %s\n", e.Emoji.Name, e.Emoji.ID)
	ts, ok := trackedScrollText[e.MessageID]
	if ok {
		fmt.Printf("scroll avant: %d\n", ts.scroll)
		if e.Emoji.Name == emojiUp {
			ts.scroll -= 1
		}
		if e.Emoji.Name == emojiDown {
			println("inc")
			ts.scroll += 1
		}
		if ts.scroll < 0 {
			ts.scroll = 0
			println("reset 0")
		}
		fmt.Printf("scroll: après: %d\n", ts.scroll)
		trackedScrollText[e.MessageID] = ts
		_, err := scrollMessage(ts, s, e.MessageID, e.ChannelID)
		he(err)
	} else {
		println("not ok")
	}
}
