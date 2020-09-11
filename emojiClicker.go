package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"sync"
	"time"
)

var trackedEmojiClickerGames map[string]*EmojiClickerGame

var progression []EmojiClickerProgression
var emojiMap map[string]int
var maxProgression int

type EmojiClickerGame struct {
	clics        int
	multipliers  int
	autoClickers int
	lock         sync.Mutex
	level        int
}
type EmojiClickerProgression struct {
	threshold    int // le seuil de clics à atteindre pour débloquer
	multiplier   int
	autoClickers int
	emoji        string
}

func init() {
	trackedEmojiClickerGames = map[string]*EmojiClickerGame{}
	progression = []EmojiClickerProgression{
		ecp(0, 1, 0),
		ecp(10, 2, 0),
		ecp(100, 2, 1),
		ecp(1000, 5, 2),
		ecp(2000, 10, 5),
		ecp(10000, 10, 10),
		ecp(100000, 100, 10),
		ecp(500000, 100, 100),
	}

	emojiMap = map[string]int{
		"🍪":  0,
		"🍰":  1,
		"🖱":  2,
		"🕹️": 3,
		"⚙️": 4,
		"🏭":  5,
		"🏙️": 6,
		"🇩🇪": 7,
	}
	for emoji, level := range emojiMap {
		progression[level].emoji = emoji
	}

	maxProgression = len(progression)
}
func ecp(seuil int, multiplieur int, autoClickers int) EmojiClickerProgression {
	return EmojiClickerProgression{
		threshold:    seuil,
		multiplier:   multiplieur,
		autoClickers: autoClickers,
	}
}
func (game *EmojiClickerGame) toString() string {
	return fmt.Sprintf("%d clics !\nniveau %d", game.clics, game.level)
}
func reverseToString(str string) (clics int, level int) {
	clics = -1
	level = 0
	fmt.Sscanf(str, "%d clics !\nniveau %d", &clics, &level)
	return clics, level
}

func (game *EmojiClickerGame) AutoClickHandler(s *discordgo.Session, channelId string, messageId string) {
	for {
		time.Sleep(time.Second * 2)
		game.lock.Lock()
		game.clics += (game.autoClickers) * game.multipliers
		hde(s.ChannelMessageEdit(channelId, messageId, game.toString()))
		fmt.Printf("auto clicking: %d clics, ac:%d,mul:%d\n", game.clics, game.autoClickers, game.multipliers)
		game.lock.Unlock()
	}
}
func (game *EmojiClickerGame) UserClick(level EmojiClickerProgression) {
	game.lock.Lock()
	game.clics += level.multiplier
	game.lock.Unlock()
}

func (game *EmojiClickerGame) LevelUp() {
	game.level += 1
	level := progression[game.level]
	game.autoClickers += level.autoClickers
	game.multipliers += level.autoClickers
}
func (game *EmojiClickerGame) Level() *EmojiClickerProgression {
	return &progression[game.level]
}

/**

 */
func emojiClicked(s *discordgo.Session, e *discordgo.MessageReaction) {
	game := trackedEmojiClickerGames[e.MessageID]

	if level, ok := emojiMap[e.Emoji.Name]; ok {
		if level <= game.level {
			println(game.clics)
			game.UserClick(progression[level])
			println(game.clics)
			// s.ChannelMessageEdit(e.ChannelID, e.MessageID, game.toString()) en fait la MaJ est faite dans le tick d'auto-click
			if game.clics >= progression[game.level+1].threshold && game.level < maxProgression {
				game.LevelUp()
				s.MessageReactionAdd(e.ChannelID, e.MessageID, game.Level().emoji)
			}
		}
	}
	return
}

func newGame(s *discordgo.Session, e *discordgo.MessageCreate) {
	game := EmojiClickerGame{
		clics:        0,
		autoClickers: 0,
		multipliers:  1,
		lock:         sync.Mutex{},
		level:        0,
	}
	msg, err := s.ChannelMessageSend(e.ChannelID, game.toString())
	he(err)
	s.MessageReactionAdd(msg.ChannelID, msg.ID, "🍪")
	trackedEmojiClickerGames[msg.ID] = &game
	go game.AutoClickHandler(s, msg.ChannelID, msg.ID)
}

func recoverGame(s *discordgo.Session, e *discordgo.MessageReaction) {

	msg, err := s.ChannelMessage(e.ChannelID, e.MessageID)
	if err == nil {

		if msg.Author.ID == s.State.User.ID {
			clics, level := reverseToString(msg.Content)
			if clics >= 0 {
				game := EmojiClickerGame{
					clics:        clics,
					autoClickers: 0,
					multipliers:  1,
					lock:         sync.Mutex{},
					level:        level,
				}
				game.RecoverStats()
				trackedEmojiClickerGames[e.MessageID] = &game
				fmt.Printf("recovered with level %d\n", level)
				emojiClicked(s, e)
				go game.AutoClickHandler(s, e.ChannelID, e.MessageID)
			}
		}
	}
	he(err)
}
func (game *EmojiClickerGame) RecoverStats() {
	for _, level := range progression[:game.level] {
		game.multipliers += level.multiplier
		game.autoClickers += level.autoClickers
	}
}
