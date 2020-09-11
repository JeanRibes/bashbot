package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

var trackedEmojiClickerGames map[string]EmojiClickerGame

var progression map[string]int

type EmojiClickerGame struct {
	//user  discordgo.User
	clics int
}

func (ecg *EmojiClickerGame) toString() string {
	return fmt.Sprintf("%d clics !", ecg.clics)
}
func reverseToString(str string) int {
	clics := 0
	fmt.Sscanf(str, "%d clics !", &clics)
	//he(err)
	return clics
}

func init() {
	trackedEmojiClickerGames = map[string]EmojiClickerGame{}
	progression = map[string]int{
		"🖱":  1,
		"🕹️": 2,
		"⚙️": 5,
		"🏭":  10,
		"🏙️": 50,
		"🇩🇪": 200,
		"🚀":  10000,
	}
}

/**

 */
func emojiClicked(s *discordgo.Session, e *discordgo.MessageReaction) {
	emojiGame := trackedEmojiClickerGames[e.MessageID]
	if /*e.UserID == emojiGame.user.ID*/ true {
		// comptages des clics
		inc := progression[e.Emoji.Name]
		if emojiGame.clics >= inc { //anto-cheat
			emojiGame.clics += inc
			fmt.Printf("incrementing with %d\n", inc)
		} else {
			fmt.Print("%s cheated with %s\n", e.UserID, e.Emoji.Name)
		}
		// feedback user
		trackedEmojiClickerGames[e.MessageID] = emojiGame
		s.ChannelMessageEdit(e.ChannelID, e.MessageID, emojiGame.toString())

		// progression
		for emoji, increment := range progression {
			if emojiGame.clics > increment {
				s.MessageReactionAdd(e.ChannelID, e.MessageID, emoji)
			}
		}
	}
}

func newGame(s *discordgo.Session, e *discordgo.MessageCreate) {
	emojiGame := EmojiClickerGame{
		//user:  *e.Author,
		clics: 1,
	}
	msg, err := s.ChannelMessageSend(e.ChannelID, emojiGame.toString())
	he(err)
	s.MessageReactionAdd(msg.ChannelID, msg.ID, "🖱")
	trackedEmojiClickerGames[msg.ID] = emojiGame
}

func recoverGame(s *discordgo.Session, e *discordgo.MessageReaction) {

	msg, err := s.ChannelMessage(e.ChannelID, e.MessageID)
	if err == nil {

		if msg.Author.ID == s.State.User.ID {
			clics := reverseToString(msg.Content)
			if clics > 0 {
				trackedEmojiClickerGames[e.MessageID] = EmojiClickerGame{
					clics: clics,
				}
				emojiClicked(s, e)
			}
		}
	}
	he(err)
}
