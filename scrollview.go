package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strings"
)

const maxScrollLength = 500

const loremIpsum = `
1Lorem Ipsum dolor sit amet, consectetur adipiscing elit. Nullam est erat, tempor ut suscipit nec, lacinia vel enim. | EOL 
2Integer elementum, dui quis condimentum tincidunt, ante ipsum euismod ipsum, vitae laoreet lacus ipsum vel enim. Donec | EOL
3viverra ligula at posuere mollis. Nullam posuere orci eget gravida sodales. Aliquam erat volutpat. Nam vel libero sit amet ante dapibus luctus quis | EOL
4vel purus. Pellentesque imperdiet sem a scelerisque sodales. Fusce eget maximus mauris. Phasellus non tortor tempor, vehicula dolor ut, euismod | EOL
5sem. Phasellus id nibh nec massa tincidunt congue. Integer at turpis turpis. Integer at euismod ex. Pellentesque lectus diam, luctus iaculis turpis sit amet, porta hendrerit odio. Proin ac efficitur elit. Cras viverra quam et erat vulputate laoreet. | EOL
6Ut sit amet tortor scelerisque, rutrum tellus in, rutrum tortor. Morbi metus ligula, ultrices quis est quis, | EOL 
7pellentesque mattis libero. Etiam velit lorem, malesuada ac consequat et, sodales ut lorem. Mauris vel turpis et | EOL
8odio egestas iaculis. Nulla dictum neque elit, sed semper velit sodales sagittis. Donec consectetur orci et volutpat scelerisque. Donec posuere, nisl | EOL
9bibendum congue laoreet, nulla mi tristique lectus, ac molestie leo dolor id ex. Nulla at placerat eros, id ultricies velit. | EOL
10Nunc ornare sapien dui, a dictum lorem venenatis id. Morbi malesuada luctus commodo. Nam ultrices libero nulla, ut malesuada augue pellentesque vel. Maecenas volutpat in eros id scelerisque. | EOL
12Nam sodales sit amet enim ac eleifend. Nam cursus elementum lectus, nec commodo nisi tristique sed. Etiam tristique | EOL 
13dui eget blandit fringilla. Duis maximus egestas tortor, in tempor arcu accumsan posuere. Aenean ac leo eu purus | EOL
14vestibulum sollicitudin. Nulla faucibus dolor dolor, consectetur luctus tellus vulputate vel. Donec aliquam aliquet accumsan. | EOL
15Nullam imperdiet, urna in eleifend | EOL
16suscipit, elit turpis varius sapien, non ullamcorper nisl justo eget odio. | EOL 
17Suspendisse tempus ante leo, nec auctor nisi lacinia sed | EOL
18. Nunc a finibus ex. Aenean viverra velit blandit elit | EOL
19lacinia pretium. Aenean lobortis magna non magna mattis, non placerat tellus scelerisque. Integer vitae auctor est. Nullam fringilla mi et eros vulputate | EOL
20, sit amet malesuada diam interdum. Vestibulum vitae est molestie, egestas eros vitae, vulputate metus. Morbi volutpat enim vel semper ullamcorper | EOL
21. Aenean tristique, enim ut porta vulputate, metus lacus tincidunt dui, nec dictum est eros at nunc. Nullam sagittis nisl in turpis lacinia, molestie egestas ipsum viverra. | EOL
21Cras auctor fermentum elit fringilla efficitur. In condimentum quam sed ex suscipit, sed tempus nisl vehicula. Maecenas | EOL 
22odio felis, vestibulum vel malesuada id, sollicitudin sed purus. Nulla dictum pulvinar venenatis. Ut commodo tristique enim | EOL
23, ac vulputate lectus porta et. Pellentesque porttitor non tellus et malesuada. In a erat porttitor, porta felis ac, rhoncus justo. | EOL
24
25Nulla sed quam nec leo pellentesque lobortis interdum eget risus. Fusce sed tristique sapien, vel dignissim dolor. Etiam | EOL 
26id massa | EOL
27sit amet massa maximus eleifend sed sed nisl. Curabitur placerat sollicitudin enim, vitae commodo nulla tristique vitae | EOL
28. Aliquam erat volutpat. Aliquam vitae arcu egestas, molestie sem et, hendrerit nunc. Quisque a pellentesque ligula. Praesent vitae dignissim nulla | EOL
29, quis rutrum leo. Proin fringilla velit in mauris porta, id rhoncus libero pulvinar. Donec ullamcorper, turpis vel porta pretium, | EOL
30lorem sapien pharetra nunc, cursus venenatis eros felis sit amet neque. Suspendisse ut elit blandit, faucibus tortor sit amet, ultricies risus. In luctus enim risus. Vestibulum in auctor velit. | EOL
31Nam mauris quam, blandit at purus rhoncus, pulvinar luctus augue. Ut luctus ultrices mauris et sodales. Aenean aliquam | EOL 
32felis non elit commodo faucibus. Quisque sagittis lorem ac porta dictum. Curabitur in nisi turpis. Sed tristique, | EOL
33mauris id vehicula commodo, erat augue pretium urna, ut scelerisque odio purus ut velit. Donec gravida dolor id tincidunt malesuada. Proin ullamcorper | EOL
34facilisis nisl ac aliquam. Nulla cursus facilisis tellus nec tempor. Ut dolor orci, eleifend at sollicitudin ac, consequat quis est. | EOL
35Aliquam feugiat pulvinar pellentesque. In a elit ac dui eleifend feugiat. Nam vestibulum lorem purus, non porttitor ante semper ut. | EOL
36Vivamus ac imperdiet sapien, ac euismod dolor. Praesent mollis, ante rhoncus efficitur molestie, sapien dui bibendum | EOL 
37augue, vitae bibendum mauris enim quis elit. Quisque nec risus lacinia, semper leo non, sollicitudin sapien. | EOL
38Fusce ut orci metus. Phasellus hendrerit est ut orci viverra, sit amet finibus tortor dignissim. Sed bibendum feugiat feugiat. Sed malesuada nec | EOL
4039arcu vel molestie. Maecenas magna enim, aliquet in elementum in, blandit et enim. Quisque non leo ultricies, dictum turpis a | EOL
41, ullamcorper tellus. Nam placerat tincidunt enim, rutrum eleifend ex ullamcorper id. Nulla non molestie erat. | EOL
42Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia curae; Maecenas eleifend quam | EOL 
43quis purus volutpat facilisis. Etiam non aliquam lacus, eget vulputate nisl. Nulla facilisi. Vivamus sollicitudin ex eget | EOL
44, tempor ipsum. Phasellus malesuada felis massa, ac congue arcu tincidunt vitae. Aliquam nunc enim, congue interdum velit a, ornare | EOL
45placerat purus. Curabitur vitae augue erat. Etiam scelerisque viverra est eu sagittis. Pellentesque ac lobortis ligula. Sed pretium mollis lectus. Quisque at urna odio. Nulla et massa sapien. | EOL
46 === EOF ===
`

const emojiUp = "⬆️"
const emojiDown = "⬇️"

type TextScroll struct {
	content string
	posLow  int
	posHigh int
}

var trackedScrollText map[string]TextScroll

func init() {
	trackedScrollText = map[string]TextScroll{}
}

func messageCreate2(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	if m.Content == "p" {
		/*t, err := exec.Command("cat", "/Users/jean/Documents/INSA/LaTeX-template/cr.tex").CombinedOutput()
		he(err)*/

		createScroll(s, m, loremIpsum)
		return true
	}
	if m.Content == "d" {
		hde(s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%d tracked scrolls", len(trackedScrollText))))
		return true
	}
	return false
}

func createScroll(s *discordgo.Session, m *discordgo.MessageCreate, content string) {
	ts := TextScroll{
		content: content,
	}
	//fmt.Printf("sortie: %s\n",ts.content)
	//_, e := scrollMessage(&ts, s, m.Message.ID, m.ChannelID)
	max := lastLineJump(ts.content, maxScrollLength)
	msg, e := s.ChannelMessageSend(m.ChannelID, "```"+ts.content[:max]+"```")
	ts.posLow = 0
	ts.posHigh = max
	if e != nil {
		fmt.Printf("%s", e)
	}
	trackedScrollText[msg.ID] = ts //l'ID du message qu'on a émis
	he(s.MessageReactionAdd(m.ChannelID, msg.ID, emojiUp))
	he(s.MessageReactionAdd(m.ChannelID, msg.ID, emojiDown))
	he(s.MessageReactionAdd(m.ChannelID, msg.ID, deleteEmoji))

	userMessageMap[msg.ID] = m.Author.ID
}

func scrollMessageDown(t TextScroll, s *discordgo.Session, messageId string, channelId string) (*discordgo.Message, error) {
	if len(t.content) < maxScrollLength {
		return nil, nil
	}
	if t.posHigh == len(t.content)-1 {
		return nil, nil
	}
	min := t.posHigh + 1
	sup := t.posHigh + maxScrollLength
	if sup > len(t.content) {
		sup = len(t.content) - 1
	} else {
		sup = min + lastLineJump(t.content[min:], maxScrollLength)
	}

	fmt.Printf("min: %d, mid: %d, max: %d, len)%d\n", min, sup, len(t.content), sup-min)
	//fmt.Printf("content: %s\n", t.content[min:sup])

	t.posLow = min
	t.posHigh = sup
	trackedScrollText[messageId] = t
	return s.ChannelMessageEdit(channelId, messageId, "```"+t.content[min:sup]+"```")
}

func scrollMessageUp(t TextScroll, s *discordgo.Session, messageId string, channelId string) (*discordgo.Message, error) {
	if len(t.content) < maxScrollLength {
		return nil, nil
	}
	if t.posLow == 0 {
		return nil, nil
	}

	inf := t.posLow - maxScrollLength
	sup := t.posHigh - maxScrollLength
	if inf < 0 {
		inf = 0
	} else { //ajustement pour coller aux sauts de ligne
		inf = inf + strings.Index(t.content[inf:sup], "\n")
	}

	if sup > len(t.content) {
		sup = len(t.content) - 1
	} else {
		sup = inf + lastLineJump(t.content[inf:], maxScrollLength)
	}

	fmt.Printf("inf: %d, sup: %d, len: %d\n", inf, sup, sup-inf)
	t.posLow = inf
	t.posHigh = sup
	trackedScrollText[messageId] = t
	return s.ChannelMessageEdit(channelId, messageId, "```"+t.content[inf:sup]+"```")
}

/*
en failt ça doit s'arrêter un saut avant la fin
*/
func lastLineJump(bigString string, max int) int {
	end := 0
	lineJumpIndex := 0
	for end < max {
		lineJumpIndex = strings.Index(bigString[end:], "\n")
		if strings.Index(bigString[end+lineJumpIndex+1:], "\n") == -1 {
			return end
		}
		if lineJumpIndex == 0 {
			end += 1
		}
		if lineJumpIndex > 0 {
			end += lineJumpIndex
		}
		if lineJumpIndex == -1 {
			println("bout !!")
			return end
		}
	}
	return end
}
