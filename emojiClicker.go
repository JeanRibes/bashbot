package main

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"sync"
	"time"
)

var trackedEmojiClickerGames map[string]*EmojiClickerGame

var progression []EmojiClickerProgression
var emojiMap map[string]int
var maxProgression int

const refreshEmoji string = "🔄"
const pauseEmoji string = "⏸"
const playEmoji string = "▶️"

type EmojiClickerGame struct {
	clics         int
	multipliers   int
	autoClickers  int
	lock          sync.Mutex
	level         int
	closeLoop     context.CancelFunc
	_displayClics int //pour savoir s'il faut update le message
	_running      bool
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

	maxProgression = len(progression) - 1
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

func (game *EmojiClickerGame) AutoClickHandler(s *discordgo.Session, channelId string, messageId string, ctx context.Context) {
	for {
		select {
		case <-ctx.Done(): // if cancel() execute
			println("quitting")
			return
		default:
			game.lock.Lock()
			game.clics += (game.autoClickers) * game.multipliers
			game.lock.Unlock()
			if game._displayClics != game.clics {
				game._displayClics = game.clics
				hde(s.ChannelMessageEdit(channelId, messageId, game.toString()))
				//fmt.Printf("auto clicking: %d clics, ac:%d,mul:%d\n", game.clics, game.autoClickers, game.multipliers)
			}
		}
		time.Sleep(time.Second * 2)
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
	game.clics = 0
}
func (game *EmojiClickerGame) Level() *EmojiClickerProgression {
	return &progression[game.level]
}

/**

 */
func emojiClicked(s *discordgo.Session, e *discordgo.MessageReaction) {
	game := trackedEmojiClickerGames[e.MessageID]

	if e.Emoji.Name == refreshEmoji {
		repostGame(s, e, game)
		return
	}

	if e.Emoji.Name == pauseEmoji && game._running {
		pauseGame(s, e, game)
		return
	}
	if e.Emoji.Name == playEmoji && !game._running {
		unpauseGame(s, e, game)
		return
	}

	if level, ok := emojiMap[e.Emoji.Name]; ok {
		if level <= game.level {
			game.UserClick(progression[level])
			// s.ChannelMessageEdit(e.ChannelID, e.MessageID, game.toString()) en fait la MaJ est faite dans le tick d'auto-click
			if game.level < maxProgression {
				if game.clics >= progression[game.level+1].threshold {
					game.LevelUp()
					s.MessageReactionAdd(e.ChannelID, e.MessageID, game.Level().emoji)
				}
			}
		}
	}
	return
}

func unpauseGame(s *discordgo.Session, e *discordgo.MessageReaction, game *EmojiClickerGame) {
	var ctx context.Context
	ctx, game.closeLoop = context.WithCancel(context.Background())
	s.MessageReactionAdd(e.ChannelID, e.MessageID, pauseEmoji)
	he(s.MessageReactionRemove(e.ChannelID, e.MessageID, playEmoji, s.State.User.ID))
	println("unpausing")
	go game.AutoClickHandler(s, e.ChannelID, e.MessageID, ctx)
	game._running = true
}

func pauseGame(s *discordgo.Session, e *discordgo.MessageReaction, game *EmojiClickerGame) {
	s.MessageReactionAdd(e.ChannelID, e.MessageID, playEmoji)
	he(s.MessageReactionRemove(e.ChannelID, e.MessageID, pauseEmoji, s.State.User.ID))
	println("pausing")
	game.closeLoop()
	game._running = false
}

func repostGame(s *discordgo.Session, e *discordgo.MessageReaction, game *EmojiClickerGame) {
	s.ChannelMessageDelete(e.ChannelID, e.MessageID)
	game.closeLoop()

	delete(trackedEmojiClickerGames, e.MessageID)

	var ctx context.Context
	ctx, game.closeLoop = context.WithCancel(context.Background())
	msg, _ := s.ChannelMessageSend(e.ChannelID, game.toString())
	trackedEmojiClickerGames[msg.ID] = game
	s.MessageReactionAdd(msg.ChannelID, msg.ID, refreshEmoji)
	s.MessageReactionAdd(msg.ChannelID, msg.ID, pauseEmoji)

	for _, level := range progression[:game.level+1] {
		he(s.MessageReactionAdd(msg.ChannelID, msg.ID, level.emoji))
		time.Sleep(time.Millisecond * 200)
	}

	go game.AutoClickHandler(s, msg.ChannelID, msg.ID, ctx)
	game._running = true
}

func newGame(s *discordgo.Session, e *discordgo.MessageCreate) {
	ctx, cancel := context.WithCancel(context.Background())
	game := EmojiClickerGame{
		_displayClics: 0,
		clics:         0,
		autoClickers:  0,
		multipliers:   1,
		lock:          sync.Mutex{},
		level:         0,
		closeLoop:     cancel,
	}
	msg, err := s.ChannelMessageSend(e.ChannelID, game.toString())
	he(err)
	s.MessageReactionAdd(msg.ChannelID, msg.ID, "🍪")
	s.MessageReactionAdd(msg.ChannelID, msg.ID, refreshEmoji)
	s.MessageReactionAdd(msg.ChannelID, msg.ID, pauseEmoji)
	trackedEmojiClickerGames[msg.ID] = &game
	go game.AutoClickHandler(s, msg.ChannelID, msg.ID, ctx)
	game._running = true
}

func recoverGame(s *discordgo.Session, e *discordgo.MessageReaction) {

	msg, err := s.ChannelMessage(e.ChannelID, e.MessageID)
	if err == nil {
		if msg.Author.ID == s.State.User.ID {
			clics, level := reverseToString(msg.Content)
			if clics >= 0 {
				ctx, cancel := context.WithCancel(context.Background())
				game := EmojiClickerGame{
					_displayClics: clics,
					clics:         clics,
					autoClickers:  0,
					multipliers:   1,
					lock:          sync.Mutex{},
					level:         level,
					closeLoop:     cancel,
				}
				game.RecoverStats()
				trackedEmojiClickerGames[e.MessageID] = &game
				fmt.Printf("recovered with level %d\n", level)
				emojiClicked(s, e)
				go game.AutoClickHandler(s, msg.ChannelID, msg.ID, ctx)
				game._running = true
			}
		}
	}
	he(err)
}
func (game *EmojiClickerGame) RecoverStats() {
	for _, level := range progression[:game.level+1] {
		game.multipliers += level.multiplier
		game.autoClickers += level.autoClickers
	}
	fmt.Printf("%d autoclickers\n", game.autoClickers)
}
