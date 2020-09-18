package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

var emojis = map[int32]string{
	'a': "🇦",
	'b': "🇧",
	'c': "🇨",
	'd': "🇩",
	'e': "🇪",
	'f': "🇫",
	'g': "🇬",
	'h': "🇭",
	'i': "🇮",
	'j': "🇯",
	'k': "🇰",
	'l': "🇱",
	'm': "🇲",
	'n': "🇳",
	'o': "🇴",
	'p': "🇵",
	'q': "🇶",
	'r': "🇷",
	's': "🇸",
	't': "🇹",
	'u': "🇺",
	'v': "🇻",
	'w': "🇼",
	'x': "🇽",
	'y': "🇾",
	'z': "🇿",
}

const deleteEmoji string = "❌"

func main() {
	//nouvelle session
	//println(lastLineJump("123456789\n123456\n7891654", 20))
	//os.Exit(0)
	sess, err := discordgo.New("Bot " + os.Getenv("BOT_TOKEN"))
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}
	//jean, err := sess.User(jeanID)

	// Register the messageCreate func as a callback for MessageCreate events.
	sess.AddHandler(messageCreate)

	/*cj, err := sess.UserChannelCreate(jeanID)
	he(err)
	m, e := sess.ChannelMessageSend(cj.ID, "helo")
	he(e)
	println(m.Content)*/

	// In this example, we only care about receiving message events.
	sess.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsDirectMessages | discordgo.IntentsGuildMessages | discordgo.IntentsGuilds | discordgo.IntentsDirectMessageReactions | discordgo.IntentsGuildMessageReactions)
	// discordgo.IntentsDirectMessages pour les DM ?

	sess.AddHandler(func(s *discordgo.Session, e *discordgo.GuildDelete) {
		println("guild delete")
		fmt.Printf("on s'est fait jeter de #%s :(\n", e.Guild.ID)
	})

	sess.AddHandler(messageReactionAdd)
	sess.AddHandler(messageReactionRemove)

	err = sess.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	} else {
		fmt.Println("BashBot tourne !")
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	he(sess.Close())
}

// à chaque message que le bot peut voir (n'impore quel channel)
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	if messageCreate2(s, m) {
		return
	}

	if m.Content == "_ec" {
		newGame(s, m)
		return
	}

	if m.Content == "_progression" {
		s.ChannelMessageSend(m.ChannelID, ExplainProgression())
	}

	if amImentionned(s, m.Message) {
		hde(s.ChannelMessageSend(m.ChannelID, "Bonjour, je suis BashBot.\n"+
			"J'exécute les commandes qu'on me donne, préfixées de '$ ' \n"+
			"Exemple: ``$ echo 'hello world'``\n"+
			"\nJe dispose aussi d'autres fonctionnalités, comme un Cookie Clicker. envoie `_ec` pour lancer le jeu et `_progression` pour voir les différents niveaux."))
		return
	}
	fmt.Printf("message #%s: %s", m.ID, m.Content)

	bash(s, m)
	if m.Content == "b" {
		msg := ""
		for _, bu := range bashed {
			msg += GetNickName(s, bu.ID, m.GuildID) + " "
		}
		s.ChannelMessageSend(m.ChannelID, msg)
		println("bashed: " + msg)
	}
	if len(m.Content) < 3 {
		return
	}

	if m.Content[0] == '$' && m.Content[1] == ' ' {
		handleBash(s, m)
	}

	if m.Content[:3] == "_b " {
		addBash(s, m)
		return
	}

	if m.Content[:3] == "_t " { //réagit avec les emojis du message
		for _, lettre := range m.Content[3:] {
			emoji, ok := emojis[lettre]
			if ok {
				s.MessageReactionAdd(m.ChannelID, m.ID, emoji)
				time.Sleep(time.Millisecond * 1000)
			}
		}
		return
	}

	if m.Content[:3] == "_e " {
		p := strings.Split(m.Content, " ")
		var msgId string
		chanId := m.ChannelID
		if strings.Contains(p[1], "-") {
			f := strings.Split(p[1], "-")
			msgId = f[1]
			chanId = f[0]
		} else {
			msgId = p[1]
		}
		texte := p[2]

		for _, lettre := range texte {
			emoji, ok := emojis[lettre]
			if ok {
				println(s.MessageReactionAdd(chanId, msgId, emoji))
				time.Sleep(time.Millisecond * 100)
			}
		}
		return
	}
}

/*func execCmd(m *discordgo.Message) string {
	if _, exists := pwdChannels[m.ChannelID]; exists {

	} else {
		pwdChannels[m.ChannelID] = "/tmp\n"
	}
	fmt.Printf("pwd: %s\n", pwdChannels[m.ChannelID])

	command := "cd " + pwdChannels[m.ChannelID] + m.Content[2:] + ";echo -n 'uu' ; pwd"
	//fmt.Printf("commande: %s\n", command)
	out, _ := exec.Command("bash", "-c", command).CombinedOutput()

	comb := strings.Split(string(out), "uu")
	pwdChannels[m.ChannelID] = comb[1]
	return prompt(pwdChannels[m.ChannelID]) + formatinput(m.Content[2:]) + "\n" + comb[0]
}*/

func messageReactionAdd(s *discordgo.Session, e *discordgo.MessageReactionAdd) {
	messageReactionAddOrRemove(s, e.MessageReaction)
}
func messageReactionRemove(s *discordgo.Session, e *discordgo.MessageReactionRemove) {
	messageReactionAddOrRemove(s, e.MessageReaction)
}

func messageReactionAddOrRemove(s *discordgo.Session, e *discordgo.MessageReaction) {
	if e.UserID == s.State.User.ID {
		return
	}
	if e.Emoji.Name == deleteEmoji {
		if ogUserId, ok := userMessageMap[e.MessageID]; ok {
			if ogUserId == e.UserID {
				s.ChannelMessageDelete(e.ChannelID, e.MessageID)
			}
		}
	}
	ts, ok := trackedScrollText[e.MessageID] //dans scrollview.go
	if ok {
		if e.Emoji.Name == emojiUp {
			_, err := scrollMessageUp(ts, s, e.MessageID, e.ChannelID)
			he(err)
			return
		}
		if e.Emoji.Name == emojiDown {
			_, err := scrollMessageDown(ts, s, e.MessageID, e.ChannelID)
			he(err)
			return
		}
	}
	if _, exists := trackedEmojiClickerGames[e.MessageID]; exists {
		emojiClicked(s, e)
		return
	}
	recoverGame(s, e)
	fmt.Printf("nom: %s id: %s\n", e.Emoji.Name, e.Emoji.ID)
}
