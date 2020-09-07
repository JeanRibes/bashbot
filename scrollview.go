package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strings"
)

const maxScrollLength = 1000

const loremIpsum = `
1 Salut 
2 Voilà la petite présentation de ma personne pour |
3 
4 Je viens de région parisienne (92 izi), je suis en 4A au département Télécom (comme la moitié du SIA presque), et j'ai fait la filière FAS (pour les bacs STI2D) lors du Premier Cycle/FIMI.|
5Ça fait un bout de temps que je fais de l'informatique, au lycée j'avais un serveur professionnel (les gros trucs lourds qui font du bruit en salle de chouille) au sous-sol de chez mes parents sur lequel j'avais un site web et tout.|
6|
7J'ai rejoint le BdE vers Octobre/Novembre 2017 quand j'étais en 1e année à l'INSA. C'était du temps où il y avait Philippe ,Gabriel, François et Hugo, et Alban commençait en tant que Resp Orgaif.|
8J'ai commencé l'année en programmant la partie de la billetterie du Gala qui communique avec la banque pour les commandes par Internet.|
9À ce moment j'étais juste 'développeur', mes contributions au code de la billetterie étaient soumises à validation et je n'avais pas accès aux données des serveurs de production (lire: je n'avais aucun accès sensible).|
10J'ai fait orga soft au Gala, notamment pour la partie Accueil où on scannait les billets des diplômés. C'était un peu pour rassurer les Resp car c'était la première année qu'on utilisait la nouvelle billetterie, et j'étais en mesure de faire des modification à la billetterie en live au cas où.|
11J'avais aussi fait entièrement un site d'inscription pour le RAID avec paiement en ligne (n'a jamais été remis en service du coup je pense pas que les orgas du Raid s'en souviennent).|
12|
13Vers la fin de la 1A, je suis devenu OrgaIF de Confiance, et j'ai reçu des accès aux serveurs et aux bases de données du BdE.|
14On a aussi créé le Adhésion qu'on a aujourd'hui (il a évolué depuis mais il vient de là), j'ai notamment fait toute la partie Wei (RIP) et paiements (qui a merdé et qu'on utilise plus).|
15|
16En 2A j’ai fait la plupart du développement tout seul, Alban s’occupait de tout ce qui est serveur/infrastructure/réseau et administratif.|
17J'ai modifié la billetterie du Gala pour la 23e édition  et je l'ai adaptée pour qu'elle serve aussi pour le Bal. C'était la première fois que le Bal utilisait la nouvelle billetterie du BdE, les commandes en ligne passaient directement par la BnP, et les personnes pouvait commander juste en présentant leur carte VA.|
18Il y avait aussi des petites modifs pour sos-laveries, le site pour signaler les problèmes des laveries.|
19La récupération des bracelets d'entrée au Bal se faisait par scan de billet ou de carte VA.|
20La plupart des OrgaIFs sont partis, il ne restait plus qu'Alban, Gabriel de temps en temps et Philippe est parti en stage.|
21C'est vers la fin de l'année que le projet du SIA a débuté, porté par Philippe Vienne et Benoît (doctorant TC et ancien Resp VA).|
22|
23|
24En 3A j'ai fait les modifications pour le Gala(rip) et le Bal (rip), les site pour l’élection des CdPs, et le site de check des cartes VA. J’ai aussi fait un nouveau site pour l’élection du Bureau (qui remplace les formulaires Limesurvey dégeus de l’INSA).|
25Je me suis formé à l'administration de notre infrastructure de nouvelle génération et à la gestion de l'authentification centralisée.|
26J'ai fait un petit peu de maintenance des PCs de la salle IF, et de la borne MA de la Coop.|
27|
28Pour le futur j'espère mettre en place l'authentification centralisée pour tous les services du BdE, et ouvrir la base de données d'Adhésion aux associations comme les 24h, tout cela en accord avec la réglementation sur la protection des données (RGPD).|
29Il faut aussi que je passe mon savoir aux nouveaux pour éviter de refaire de zéro les applications tous les 3 ans.|
30|
31|
32Jean Ribes|
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

		ts := TextScroll{
			content: loremIpsum,
		}
		//fmt.Printf("sortie: %s\n",ts.content)
		//_, e := scrollMessage(&ts, s, m.Message.ID, m.ChannelID)
		max := lastLineJump(ts.content, maxScrollLength)
		msg, e := s.ChannelMessageSend(m.ChannelID, "```"+ts.content[:max]+"```")
		ts.posLow = 0
		ts.posHigh = max
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
	ts, ok := trackedScrollText[e.MessageID]
	if ok {
		if e.Emoji.Name == emojiUp {
		}
		if e.Emoji.Name == emojiDown {
			trackedScrollText[e.MessageID] = ts
			_, err := scrollMessageDown(ts, s, e.MessageID, e.ChannelID)
			he(err)
			return
		}
	} else {
		println("not ok")
	}
	fmt.Printf("nom: %s id: %s\n", e.Emoji.Name, e.Emoji.ID)
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
